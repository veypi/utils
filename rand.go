package utils

//
// rand.go
// Copyright (C) 2020 light <light@1870499383@qq.com>
//
// Distributed under terms of the MIT license.
//

import (
	"math/rand"
	"time"
	"unsafe"
)

// 字符列表
var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var size = int32(len(letters))
var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandSeq 构造指定长度的随机字符串
// 字符包括 数字 和 大小写字母
func RandSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[seed.Int31n(size)]
	}
	return *(*string)(unsafe.Pointer(&b))
}

// Rand 构造指定长度的随机 Byte 序列
func Rand(n int) []byte {
	b := make([]byte, n)
	seed.Read(b)
	return b
}
