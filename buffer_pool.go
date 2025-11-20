package syndrdbsimd

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	// metricsGranularitySec defines the granularity of metrics collection in seconds.
	// Metrics are aggregated into buckets of this size.
	// Can be adjusted for finer or coarser granularity as the project matures.
	metricsGranularitySec = 1

	// metricsWindowMinutes defines how many minutes of metrics history to keep.
	// TODO: Make metrics window duration configurable in future.
	metricsWindowMinutes = 10

	// maxBufferSize is the maximum size of pooled buffers (64KB).
	// Buffers larger than this are allocated on the heap and not pooled.
	maxBufferSize = 64 * 1024

	// metricsRingSize is the number of metric entries in the ring buffer.
	metricsRingSize = (metricsWindowMinutes * 60) / metricsGranularitySec
)

// Buffer pool sizes for pre-warming.
// TODO: Make warmup sizes configurable in future.
var poolWarmupSizes = []int{
	256,      // 256 bytes
	1024,     // 1 KB
	4096,     // 4 KB
	16384,    // 16 KB
	65536,    // 64 KB
}

// bufferPool is a sync.Pool for reusing byte slices.
var bufferPool = sync.Pool{
	New: func() interface{} {
		// Return nil - we'll allocate on-demand in getPooledBuffer
		return nil
	},
}

// metricsEntry holds metrics for a single time bucket.
type metricsEntry struct {
	timestamp    int64  // Unix timestamp in seconds
	poolHits     uint64 // Number of successful pool retrievals
	poolMisses   uint64 // Number of pool misses (new allocations)
	heapFallbacks uint64 // Number of heap allocations (>64KB)
	valid        bool   // Whether this entry contains valid data
}

// metricsRing holds the circular buffer of metrics.
var metricsRing [metricsRingSize]metricsEntry
var metricsRingIndex atomic.Uint64
var metricsRingMutex sync.Mutex

// BufferPoolStats contains aggregated statistics about buffer pool usage.
type BufferPoolStats struct {
	PoolHits       uint64    // Total successful pool retrievals in window
	PoolMisses     uint64    // Total pool misses (new allocations) in window
	HeapFallbacks  uint64    // Total heap allocations for >64KB buffers in window
	WindowStartTime time.Time // Start of the metrics window
	WindowEndTime   time.Time // End of the metrics window
}

func init() {
	// Pre-warm the buffer pool with common sizes.
	// Creates 4 buffers of each size to reduce first-call latency.
	for _, size := range poolWarmupSizes {
		for i := 0; i < 4; i++ {
			buf := make([]byte, size)
			bufferPool.Put(&buf)
		}
	}
}

// getPooledBuffer returns a byte slice of at least the requested size.
// For sizes <= 64KB, it attempts to reuse a buffer from the pool.
// For sizes > 64KB, it allocates on the heap and records a heap fallback metric.
func getPooledBuffer(size int) []byte {
	if size > maxBufferSize {
		// Too large for pool - allocate on heap
		recordMetric(false, false, true)
		return make([]byte, size)
	}

	// Try to get from pool
	if bufInterface := bufferPool.Get(); bufInterface != nil {
		buf := bufInterface.(*[]byte)
		// Check if buffer is large enough
		if cap(*buf) >= size {
			recordMetric(true, false, false)
			return (*buf)[:size]
		}
		// Buffer too small - return to pool and allocate new
		bufferPool.Put(bufInterface)
	}

	// Pool miss - allocate new buffer
	recordMetric(false, true, false)
	buf := make([]byte, size)
	return buf
}

// returnPooledBuffer returns a buffer to the pool for reuse.
// Only buffers <= 64KB are pooled; larger buffers are discarded for GC.
func returnPooledBuffer(buf []byte) {
	if cap(buf) <= maxBufferSize {
		// Reset length to capacity for reuse
		buf = buf[:cap(buf)]
		bufferPool.Put(&buf)
	}
	// Buffers > 64KB are not pooled, let GC handle them
	// TODO: Consider explicit cleanup for long-running idle processes.
}

// recordMetric records a buffer pool event in the metrics ring buffer.
func recordMetric(hit, miss, heapFallback bool) {
	now := time.Now().Unix()
	bucket := now / metricsGranularitySec

	metricsRingMutex.Lock()
	defer metricsRingMutex.Unlock()

	// Find or create entry for current bucket
	idx := metricsRingIndex.Load() % metricsRingSize
	entry := &metricsRing[idx]

	// Check if we need a new bucket
	if !entry.valid || entry.timestamp != bucket {
		// Advance to next slot
		idx = metricsRingIndex.Add(1) % metricsRingSize
		entry = &metricsRing[idx]
		entry.timestamp = bucket
		entry.poolHits = 0
		entry.poolMisses = 0
		entry.heapFallbacks = 0
		entry.valid = true
	}

	// Update counters
	if hit {
		entry.poolHits++
	}
	if miss {
		entry.poolMisses++
	}
	if heapFallback {
		entry.heapFallbacks++
	}
}

// GetBufferPoolStats returns aggregated statistics for the metrics window.
// It sums up all valid entries in the ring buffer within the time window.
//
// TODO: Expose via monitoring endpoint in future.
func GetBufferPoolStats() BufferPoolStats {
	metricsRingMutex.Lock()
	defer metricsRingMutex.Unlock()

	now := time.Now().Unix()
	windowStart := now - (metricsWindowMinutes * 60)

	var stats BufferPoolStats
	var oldestTime, newestTime int64

	// Iterate through ring buffer and aggregate valid entries
	for i := 0; i < metricsRingSize; i++ {
		entry := &metricsRing[i]
		if !entry.valid {
			continue
		}

		entryTime := entry.timestamp * metricsGranularitySec
		if entryTime >= windowStart && entryTime <= now {
			stats.PoolHits += entry.poolHits
			stats.PoolMisses += entry.poolMisses
			stats.HeapFallbacks += entry.heapFallbacks

			if oldestTime == 0 || entryTime < oldestTime {
				oldestTime = entryTime
			}
			if entryTime > newestTime {
				newestTime = entryTime
			}
		} else if entryTime < windowStart {
			// Mark old entries as invalid
			entry.valid = false
		}
	}

	if oldestTime > 0 {
		stats.WindowStartTime = time.Unix(oldestTime, 0)
	} else {
		stats.WindowStartTime = time.Unix(windowStart, 0)
	}

	if newestTime > 0 {
		stats.WindowEndTime = time.Unix(newestTime, 0)
	} else {
		stats.WindowEndTime = time.Now()
	}

	return stats
}
