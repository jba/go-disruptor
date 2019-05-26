package implementations

import (
	"runtime"
	"sync/atomic"
	"time"
)

const cpuCacheLinePadding = 7

type LockFreeQueue struct {
	items   []int64
	lastPut int64 // atomic
	padding [cpuCacheLinePadding]int64
	lastGot int64 // atomic
	mask    int64
}

func NewLockFreeQueue(capacity int) *LockFreeQueue {
	return &LockFreeQueue{
		items:   make([]int64, capacity),
		lastPut: -1,
		lastGot: -1,
		mask:    int64(capacity - 1), // assumes capacity is a power of 2
	}
}

func (q *LockFreeQueue) Put(x int64) {
	// Assume a single writer (this goroutine).
	lp := atomic.LoadInt64(&q.lastPut)
	lg := atomic.LoadInt64(&q.lastGot)
	l := int64(len(q.items))
	for lp-lg == l {
		runtime.Gosched()
		lg = atomic.LoadInt64(&q.lastGot)
	}
	q.items[int((lp+1)&q.mask)] = x
	atomic.AddInt64(&q.lastPut, 1)
}

func (q *LockFreeQueue) Get() int64 {
	// Assume a single reader (this goroutine).
	lp := atomic.LoadInt64(&q.lastPut)
	lg := atomic.LoadInt64(&q.lastGot)
	for lp-lg == 0 {
		time.Sleep(time.Millisecond)
		lp = atomic.LoadInt64(&q.lastPut)
	}
	x := q.items[int((lg+1)&q.mask)]
	atomic.AddInt64(&q.lastGot, 1)
	return x
}
