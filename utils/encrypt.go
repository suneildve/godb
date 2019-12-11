package utils

import (
	"golang.org/x/crypto/scrypt"
	"encoding/base64"
	"log"
)

func Encrypt(str string) string{
	salt := []byte{0xc8, 0x28, 0xf2, 0x58, 0xa7, 0x6a, 0xad, 0x7b}
	dk , err := scrypt.Key([]byte(str), []byte(salt), 16384, 8, 1, 32)
// 	dk, err := scrypt.Key([]byte("some password"), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(dk)
}