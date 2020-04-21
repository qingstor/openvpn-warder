package client

import (
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

// GetUserResp represent data reponse from server get user method.
type GetUserResp struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	UpdatedAt string `json:"updated_at"`
	Admin     string `json:"admin"`
	Ignore    string `json:"ignore"`
}

func generatePassword(salt, rawPasswd string) []byte {
	passEn := pbkdf2.Key([]byte(rawPasswd), []byte(salt), 4096, 32, sha256.New)

	return passEn
}
