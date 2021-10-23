package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var ExpDelta int64 = 60 * 60 * 24

var (
	InvalidToken = errors.New("invalid token")
	ExpiredToken = errors.New("expired token")
)

type Payload struct {
	Iat int64 `json:"iat"` //token time
	Exp int64 `json:"exp"`
}

func (p *Payload) SetIat(t int64) {
	p.Iat = t
}

func (p *Payload) GetIat() int64 {
	return p.Iat
}
func (p *Payload) SetExp() {
	p.Exp = p.Iat + ExpDelta
}

func (p *Payload) GetExp() int64 {
	return p.Exp
}

func (p *Payload) IsExpired() bool {
	if time.Now().Unix() > p.Exp {
		return true
	}
	return false
}

type PayloadInterface interface {
	SetIat(int64)
	GetIat() int64
	SetExp()
	GetExp() int64
	IsExpired() bool
}

func GetToken(payload PayloadInterface, key []byte) (string, error) {
	header := map[string]string{
		"typ": "JWT",
		"alg": "HS256",
	}
	//header := "{\"typ\": \"JWT\", \"alg\": \"HS256\"}"
	now := time.Now().Unix()
	payload.SetIat(now)
	payload.SetExp()
	a, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	A := base64.StdEncoding.EncodeToString(a)
	B := base64.StdEncoding.EncodeToString(b)
	hmacCipher := hmac.New(sha256.New, key)
	hmacCipher.Write([]byte(A + "." + B))
	C := hmacCipher.Sum(nil)
	return A + "." + B + "." + base64.StdEncoding.EncodeToString(C), nil
}

func ParseToken(token string, payload PayloadInterface, key []byte) (bool, error) {
	var A, B, C string
	if seqs := strings.Split(token, "."); len(seqs) == 3 {
		A, B, C = seqs[0], seqs[1], seqs[2]
	} else {
		return false, InvalidToken
	}
	tempPayload, err := base64.StdEncoding.DecodeString(B)
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(tempPayload, payload); err != nil {
		return false, err
	}
	hmacCipher := hmac.New(sha256.New, key)
	hmacCipher.Write([]byte(A + "." + B))
	tempC := hmacCipher.Sum(nil)
	if !hmac.Equal([]byte(C), []byte(base64.StdEncoding.EncodeToString(tempC))) {
		return false, nil
	}
	if payload.IsExpired() {
		return false, ExpiredToken
	}
	return true, nil
}
