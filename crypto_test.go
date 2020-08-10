package utils

//
// crypto_test.go
// Copyright (C) 2020 light <light@1870499383@qq.com>
//
// Distributed under terms of the MIT license.
//

import (
	"fmt"
	"github.com/veypi/utils/log"
	"testing"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

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

func TestGetRsaKey(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			log.Error().Err(nil).Msgf("%v", e)
		}
	}()
	base := "pZUTCEBr4FhPb/7OemgBWkcBWsTMSELRFzvKAW6FDMcozQcQwo9yI2Sq2S//90vTkahPQKBWRYM1zvTnEIJy28oS1nNUJiykOA0U7Ozbue8fHbi8QeyegtvkVlMNch39TcDRh9NFI72LZE8FJCvYt5WhPmIFuqscjw0H0oI1DmY="
	s, err := FromBase64(base)
	if err != nil {
		t.Error(err)
		return
	}
	log.Warn().Msg(string(s))

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
