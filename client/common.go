package client

import (
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

func generatePassword(salt, rawPasswd string) []byte {
    passEn := pbkdf2.Key([]byte(rawPasswd), []byte(salt), 4096, 32, sha256.New)

    return passEn
}
