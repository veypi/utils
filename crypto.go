package utils

//
// crypto.go
// Copyright (C) 2020 light <light@1870499383@qq.com>
//
// Distributed under terms of the MIT license.
//

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
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

// rsa

func GetRsaKey(bits int) (public *rsa.PublicKey, private *rsa.PrivateKey, err error) {
	private, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	public = &private.PublicKey
	return
}

func GetPublicStr(key *rsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}
	return string(pem.EncodeToMemory(publicBlock)), nil
}

func GetPrivateStr(key *rsa.PrivateKey) (string, error) {
	derStream := x509.MarshalPKCS1PrivateKey(key)
	priBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	return string(pem.EncodeToMemory(priBlock)), nil
}

func GetPublicFromStr(key string) (*rsa.PublicKey, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	return pub, nil
}

func GetPrivateFromStr(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func RsaEncode(msg string, key *rsa.PublicKey) (string, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		key,
		[]byte(msg),
		nil)
	if err != nil {
		return "", err
	}
	return ToBase64(encryptedBytes), nil
}

func RsaDecode(msg string, key *rsa.PrivateKey) (string, error) {
	raw, err := FromBase64(msg)
	if err != nil {
		return "", err
	}
	decryptedBytes, err := key.Decrypt(nil, raw, &rsa.OAEPOptions{Hash: crypto.SHA256})
	return string(decryptedBytes), err
}

func RsaSign(msg string, key *rsa.PrivateKey) (string, error) {
	signature, err := rsa.SignPSS(rand.Reader, key, crypto.SHA256, HashSha256Byte([]byte(msg)), nil)
	if err != nil {
		return "", err
	}
	return ToBase64(signature), nil
}

func RsaCheckSign(msg string, sign string, key *rsa.PublicKey) error {
	raw, err := FromBase64(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPSS(key, crypto.SHA256, HashSha256Byte([]byte(msg)), raw, nil)
}
