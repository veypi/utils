package utils

//
// crypto_test.go
// Copyright (C) 2020 light <light@1870499383@qq.com>
//
// Distributed under terms of the MIT license.
//

import (
	"encoding/hex"
	"testing"

	"github.com/veypi/utils/logv"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestAes(t *testing.T) {
	// text := RandSeq(32)
	text := "9a57abf6444d45b7ab1bd5d357d90e86"
	key, _ := hex.DecodeString("dcaece4ab04454cdf20208d9ff537aea16d140ac3ae57edc865582d707306b41")
	iv, _ := hex.DecodeString("c611c0b0bb274e6d9577777f37420240")
	xText, err := AesEncrypt([]byte(text), key[:32], iv)
	if err != nil {
		t.Errorf(err.Error())
	}
	nText, err := AesDecrypt([]byte("SNF3gMXzrVLEwyI1O97GUbFjJ5AW05lL7xdsP3OWq1dk4QA2yrtIEYKas0sTppQ3"), key[:32], iv)
	if err != nil {
		t.Errorf(err.Error())
	}
	if text != nText {
		t.Errorf("aes is not right.")
	} else {
		t.Logf("\ntext(%d) %s;\nxtext(%d) %s;\nntext(%d) %s;\nkey(%d) %v \niv(%d) %v",
			len(text), text, len(xText), xText, len(nText), nText, len(key), key, len(iv), iv)
	}
}

func TestGetRsaKey(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			logv.Error().Err(nil).Msgf("%v", e)
		}
	}()
	base := "pZUTCEBr4FhPb/7OemgBWkcBWsTMSELRFzvKAW6FDMcozQcQwo9yI2Sq2S//90vTkahPQKBWRYM1zvTnEIJy28oS1nNUJiykOA0U7Ozbue8fHbi8QeyegtvkVlMNch39TcDRh9NFI72LZE8FJCvYt5WhPmIFuqscjw0H0oI1DmY="
	_, err := FromBase64(base)
	if err != nil {
		t.Error(err)
		return
	}

	msg := "msg 123 111@#-()'\"         "
	pub, pri, err := GetRsaKey(1024)
	panicErr(err)
	sPub, err := GetPublicStr(pub)
	panicErr(err)
	sPri, err := GetPrivateStr(pri)
	panicErr(err)
	t.Log(sPub, sPri)
	nPub, err := GetPublicFromStr(sPub)
	panicErr(err)
	nPri, err := GetPrivateFromStr(sPri)
	panicErr(err)

	cMsg, err := RsaEncode(msg, pub)
	panicErr(err)
	cMsgr, err := RsaDecode(cMsg, pri)
	panicErr(err)
	nMsg, err := RsaEncode(msg, nPub)
	panicErr(err)
	nMsgr, err := RsaDecode(nMsg, nPri)
	panicErr(err)
	t.Log(cMsgr)
	t.Log(nMsgr)
	if cMsgr != msg || nMsgr != msg {
		t.Error("decode or encode failed")
	}
	sign, err := RsaSign(msg, pri)
	panicErr(err)
	err = RsaCheckSign(msg, sign, pub)
	if err != nil {
		t.Error(err)
	}
	t.Log(sign)
}
