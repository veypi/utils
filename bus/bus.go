//
// bus.go
// Copyright (C) 2025 veypi <i@veypi.com>
// 2025-03-05 17:14
// Distributed under terms of the MIT license.
//

package bus

import (
	"github.com/veypi/utils"
	"github.com/veypi/utils/logv"
)

func New[T any]() *Bus[T] {
	return &Bus[T]{
		capacity: 10,
		cache:    make([]*T, 0, 10),
	}
}

type Bus[T any] struct {
	utils.FastLocker
	fns      []func(T)
	cache    []*T
	capacity int
}

func (b *Bus[T]) On(fn func(T)) func() {
	b.fns = append(b.fns, fn)
	idx := len(b.fns) - 1
	for _, v := range b.cache {
		fn(*v)
	}
	return func() {
		b.fns = append(b.fns[:idx], b.fns[idx+1:]...)
	}
}

func (b *Bus[T]) Emit(data T) {
	idx := 0
	if len(b.cache) < b.capacity {
		b.cache = append(b.cache, &data)
	} else {
		b.cache = append(b.cache[1:], &data)
	}
	for idx < len(b.fns) {
		fn := b.fns[idx]
		idx++
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logv.Warn().Msgf("bus emit error: %v", err)
				}
				fn(data)
			}()
		}()
	}
}
