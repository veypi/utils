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
		fnMap: make(map[int]func(T)),
	}
}

type Bus[T any] struct {
	utils.FastLocker
	count int
	fnMap map[int]func(T)
}

func (b *Bus[T]) On(fn func(T)) func() {
	idx := b.count
	b.count++
	b.fnMap[idx] = fn
	return func() {
		delete(b.fnMap, idx)
	}
}

func (b *Bus[T]) Emit(data T) {
	idx := 0
	for idx < b.count {
		fn, ok := b.fnMap[idx]
		idx++
		if !ok {
			continue
		}
		func(d T) {
			defer func() {
				if err := recover(); err != nil {
					logv.Warn().Msgf("bus emit error: %v", err)
				}
			}()
			fn(d)
		}(data)
	}
}
