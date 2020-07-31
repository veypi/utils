package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func HashMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func HashSha256Byte(msg []byte) []byte {
	msgHash := sha256.New()
	_, err := msgHash.Write(msg)
	if err != nil {
		panic(err)
	}
	msgHashSum := msgHash.Sum(nil)
	return msgHashSum
}

func HashSha256(msg string) string {
	return hex.EncodeToString(HashSha256Byte([]byte(msg)))
}

func ToBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func FromBase64(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}
