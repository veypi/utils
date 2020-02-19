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
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var size = len(letters)

// RandSeq produce random string seq
func RandSeq(n int) string {
	b := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[r.Intn(size)]
	}
	return string(b)
}
