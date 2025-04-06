package main

import "sync"

// We can make use of channels to manage concurrent requests without blocking
// We can use Mutex, RWMutex, Atomic (sync/atomic)
// sync.Mutex for cases that require critical selection (we use full lock, that means one goroutine at a time)
// sync.RWMutex for cases that allows multiplr reads and rare writes, so we allow simulatenous reads and we lock writes
// sync/atomic for cases with lightweight, fast lock-free counters

type Counter struct {
	mu    sync.RWMutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Get() int {
	c.mu.RUnlock()
	defer c.mu.RUnlock()
	return c.value
}
