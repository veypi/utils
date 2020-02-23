package utils

//
// crypto_test.go
// Copyright (C) 2020 light <light@1870499383@qq.com>
//
// Distributed under terms of the MIT license.
//

import (
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {
	text := RandSeq(32)
	fmt.Println(len(text))
	key := []byte("123456")
	xText, err := AesEncrypt(text, key)
	if err != nil {
		t.Errorf(err.Error())
	}
	nText, err := AesDecrypt(xText, key)
	if err != nil {
		t.Errorf(err.Error())
	}
	if text != nText {
		t.Errorf("aes is not right.")
	} else {
		t.Logf("\ntext(%d) %s;\nxtext(%d) %s;\nntext(%d) %s;\nkey(%d) %s",
			len(text), text, len(xText), xText, len(nText), nText, len(key), key)
	}
}
