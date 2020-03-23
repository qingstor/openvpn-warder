package handlers

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/qingstor/openvpn-warder/config"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/yaml.v2"
)

// WriteDBLock for write db file.
var WriteDBLock sync.Mutex

// AuthUser represent auth user parsed from body.
type AuthUser struct {
	Name     *string `form:"name" json:"name" xml:"name" binding:"required"`
	Password *[]byte `form:"password" xml:"password" json:"password" binding:"required"`
}

// GetUserBody represent body of GetUser method.
type GetUserBody struct {
	AuthUser *AuthUser `form:"auth_user" json:"auth_user" xml:"auth_user" binding:"required"`
}

// CreateUserBody represent body of AddUser method.
type CreateUserBody struct {
	AuthUser *AuthUser    `form:"auth_user" json:"auth_user" xml:"auth_user" binding:"required"`
	User     *config.User `form:"user" json:"user" xml:"user" binding:"required"`
}

// DeleteUserBody represent body of AddUser method.
type DeleteUserBody struct {
	AuthUser *AuthUser `form:"auth_user" json:"auth_user" xml:"auth_user" binding:"required"`
}

// ResetUserBody represent body of AddUser method.
type ResetUserBody struct {
	AuthUser *AuthUser    `form:"auth_user" json:"auth_user" xml:"auth_user" binding:"required"`
	User     *config.User `form:"user" json:"user" xml:"user" binding:"required"`
}

// HTTPError represent http error code and msg.
type HTTPError struct {
	Code int
	Msg  gin.H
}

func checkPassword(salt string, rawPasswd string, encrptedPasswd []byte) error {
	passEn := pbkdf2.Key([]byte(rawPasswd), []byte(salt), 4096, 32, sha256.New)
	if !bytes.Equal(passEn, encrptedPasswd) {
		return errors.New("Password not matched")
	}

	return nil
}

func parseAuthUser(au *AuthUser, dbPath string) (u *config.User, err *HTTPError) {
	// Parse query
	if au.Name == nil {
		return nil, &HTTPError{
			Code: http.StatusBadRequest,
			Msg:  gin.H{"msg": "Auth user name in body is empty"},
		}
	}
	if au.Password == nil {
		return nil, &HTTPError{
			Code: http.StatusBadRequest,
			Msg:  gin.H{"msg": "Auth user password in body is empty"},
		}
	}

	//  db user info
	db := config.NewDB()
	e := db.LoadFromFilePath(dbPath)
	if e != nil {
		return nil, &HTTPError{
			Code: http.StatusInternalServerError,
			Msg:  gin.H{"msg": e.Error()},
		}
	}

	for _, user := range db.Users {
		if *user.Name == *au.Name {
			e := checkPassword(*au.Name, *user.Password, *au.Password)
			if e != nil {
				return nil, &HTTPError{
					Code: http.StatusForbidden,
					Msg:  gin.H{"msg": e.Error()},
				}
			}

			return user, nil
		}
	}

	return nil, &HTTPError{
		Code: http.StatusNotFound,
		Msg:  gin.H{"msg": "User not found"},
	}
}

func getUserFromDB(userName string) (user *config.User, err error) {
	db := config.NewDB()
	err = db.LoadFromFilePath(config.DBPath)
	if err != nil {
		return nil, err
	}
	for _, u := range db.Users {
		if *u.Name == userName {
			return u, nil
		}
	}

	return
}

func readEntireDB() (db *config.DBConfig, err error) {
	db = config.NewDB()
	err = db.LoadFromFilePath(config.DBPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func writeEntireDB(db *config.DBConfig) (err error) {
	yamlData, err := yaml.Marshal(db)
	if err != nil {
		return
	}

	WriteDBLock.Lock()
	err = ioutil.WriteFile(config.DBPath, yamlData, 0777)
	if err != nil {
		return
	}
	WriteDBLock.Unlock()

	return nil
}
