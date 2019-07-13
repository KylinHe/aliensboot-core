package pool

import "sync"


var DefaultSyncPool *SyncPool

func NewPool() *SyncPool {
	DefaultSyncPool = NewSyncPool(
		5,
		30000,
		2,
	)
	return DefaultSyncPool
}

func Alloc(size int) []int64 {
	return DefaultSyncPool.Alloc(size)
}

func Free(mem []int64) {
	DefaultSyncPool.Free(mem)
}

// SyncPool is a sync.Pool base slab allocation memory pool
type SyncPool struct {
	classes     []sync.Pool
	classesSize []int
	minSize     int
	maxSize     int
}

func NewSyncPool(minSize, maxSize, factor int) *SyncPool {
	n := 0
	for chunkSize := minSize; chunkSize <= maxSize; chunkSize *= factor {
		n++
	}
	pool := &SyncPool{
		make([]sync.Pool, n),
		make([]int, n),
		minSize, maxSize,
	}
	n = 0
	for chunkSize := minSize; chunkSize <= maxSize; chunkSize *= factor {
		pool.classesSize[n] = chunkSize
		pool.classes[n].New = func(size int) func() interface{} {
			return func() interface{} {
				buf := make([]int64, size)
				return &buf
			}
		}(chunkSize)
		n++
	}
	return pool
}

func (pool *SyncPool) Alloc(size int) []int64 {
	if size <= pool.maxSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				mem := pool.classes[i].Get().(*[]int64)
				// return (*mem)[:size]
				return (*mem)[:0]
			}
		}
	}
	return make([]int64, 0, size)
}

func (pool *SyncPool) Free(mem []int64) {
	if size := cap(mem); size <= pool.maxSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				pool.classes[i].Put(&mem)
				return
			}
		}
	}
}