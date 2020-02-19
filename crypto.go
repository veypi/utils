package utils

//
// crypto.go
// Copyright (C) 2020 light <light@1870499383@qq.com>
//
// Distributed under terms of the MIT license.
//

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// PKCS5Padding  密文填充
func PKCS5Padding(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	paddingSuffix := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, paddingSuffix...)
}

// PKCS5UnPadding  取消填充
func PKCS5UnPadding(origData []byte) ([]byte, bool) {
	length := len(origData)
	padding := int(origData[length-1])
	if padding >= length {
		return nil, false
	}
	return origData[:(length - padding)], true
}

// AesEncrypt aes 加密
// key: 16, 24, or 32 bytes to select
func AesEncrypt(orig string, key []byte) (string, error) {
	key = PKCS5Padding(key, 32)[:32]
	origData := []byte(orig)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize() // 16
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// AesDecrypt aes解密
// key: 16, 24, or 32 bytes to select
func AesDecrypt(encrypted string, key []byte) (string, error) {
	key = PKCS5Padding(key, 32)[:32]
	cryptData, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cryptData))
	blockMode.CryptBlocks(origData, cryptData)
	if origData, isRight := PKCS5UnPadding(origData); isRight {
		return string(origData), nil
	}
	return "", errors.New("invalid key")
}
