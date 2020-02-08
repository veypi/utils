package utils

import (
	"runtime"
	"sync/atomic"
)

type FastLocker struct {
	mu uint32
}

func (mu *FastLocker) Lock() {
	for !atomic.CompareAndSwapUint32(&mu.mu, 0, 1) {
		runtime.Gosched()
	}
}

func (mu *FastLocker) Unlock() {
	atomic.StoreUint32(&mu.mu, 0)
}
