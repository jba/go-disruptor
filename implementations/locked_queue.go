package implementations

import "sync"

type LockedQueue struct {
	mu                sync.Mutex
	notEmpty, notFull *sync.Cond
	items             []int64
	lastPut, lastGot  int
}

func NewLockedQueue(capacity int) *LockedQueue {
	q := &LockedQueue{
		items:   make([]int64, capacity),
		lastPut: -1,
		lastGot: -1,
	}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

func (q *LockedQueue) Put(x int64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for q.size() == len(q.items) {
		q.notFull.Wait()
	}
	q.lastPut++
	q.items[int(q.lastPut%len(q.items))] = x
	if q.size() == 1 {
		q.notEmpty.Broadcast()
	}
}

func (q *LockedQueue) Get() int64 {
	q.mu.Lock()
	defer q.mu.Unlock()
	for q.size() == 0 {
		q.notEmpty.Wait()
	}
	q.lastGot++
	if q.size() == len(q.items)-1 {
		q.notFull.Broadcast()
	}
	return q.items[q.lastGot%len(q.items)]
}

func (q *LockedQueue) size() int { return q.lastPut - q.lastGot }
