//
// slice.go
// Copyright (C) 2024 veypi <i@veypi.com>
// 2024-09-09 11:07
// Distributed under terms of the MIT license.
//

package utils

import "sort"

func SortItems(x any, less func(i int, j int, sid func(string) int) bool, items ...string) {
	fc := func(s string) int {
		for i, ii := range items {
			if s == ii {
				return i
			}
		}
		return len(items)
	}
	sort.Slice(x, func(i, j int) bool {
		return less(i, j, fc)
	})
}

func SliceGet[T any](slice []T, fc func(T) bool) *T {
	for _, s := range slice {
		if fc(s) {
			return &s
		}
	}
	return nil
}

func ForEach[T any](slice []T, fn func(int, T) error) error {
	for i, s := range slice {
		err := fn(i, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertAt[T any](slice []T, index int, value T) []T {
	// 检查插入位置是否有效
	if index > len(slice) || len(slice) == 0 {
		return append(slice, value)
	} else if index < 0 {
		index = len(slice) + index
	}

	// 使用 append 和 切片操作符来插入元素
	return append(slice[:index], append([]T{value}, slice[index:]...)...)
}

// InList 判断列表是否含有某元素
func InList(str string, list []string) bool {
	for _, temp := range list {
		if str == temp {
			return true
		}
	}
	return false
}

// RemoveRep 通过map主键唯一的特性过滤重复元素
func RemoveRep(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}
