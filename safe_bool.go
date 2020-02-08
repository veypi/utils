package utils

import "sync/atomic"

// 线程安全的标志位 set 方法保证置换操作的原子性
type SafeBool uint32

func (b *SafeBool) SetTrue() bool {
	return atomic.CompareAndSwapUint32((*uint32)(b), 0, 1)
}

func (b *SafeBool) SetFalse() bool {
	return atomic.CompareAndSwapUint32((*uint32)(b), 1, 0)
}

func (b *SafeBool) ForceSetTrue() {
	atomic.StoreUint32((*uint32)(b), 1)
}
func (b *SafeBool) ForceSetFalse() {
	atomic.StoreUint32((*uint32)(b), 0)
}

func (b *SafeBool) IfTrue() bool {
	return atomic.LoadUint32((*uint32)(b)) == 1
}
