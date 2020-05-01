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

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var size = int32(len(letters))
var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandSeq produce random string seq
func RandSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[seed.Int31n(size)]
	}
	return *(*string)(unsafe.Pointer(&b))
}

func Rand(n int) []byte {
	b := make([]byte, n)
	seed.Read(b)
	return b
}
