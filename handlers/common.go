package handlers

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qingstor/openvpn-warder/config"
	"golang.org/x/crypto/pbkdf2"
)

// AuthUser represent auth user parsed from body.
type AuthUser struct {
    Name *string `form:"name" json:"name" xml:"name" binding:"required"`
    Password *[]byte `form:"password" xml:"password" json:"password" binding:"required"`
}

// GetUserBody represent body of GetUser method.
type GetUserBody struct {
    AuthUser *AuthUser `form:"auth_user" json:"auth_user" xml:"auth_user" binding:"required"`
}

// HTTPError represent http error code and msg.
type HTTPError struct {
    Code int
    Msg  gin.H
}

func checkPassword(salt string, rawPasswd string, encrptedPasswd []byte) error{
    passEn := pbkdf2.Key([]byte(rawPasswd), []byte(salt), 4096, 32, sha256.New)
    if ! bytes.Equal(passEn, encrptedPasswd) {
        return errors.New("Password not matched")
    }

    return nil
}

func parseAuthUser(au *AuthUser, dbPath string) (u *config.User, err *HTTPError) {    
    // Parse query
    if au.Name == nil {
        return nil, &HTTPError{
            Code: http.StatusBadRequest,
            Msg: gin.H{"msg": "Auth user name in body is empty"},
        }
    }
    if au.Password == nil {
        return nil, &HTTPError{
            Code: http.StatusBadRequest,
            Msg: gin.H{"msg": "Auth user password in body is empty"},
        }
    }

    //  db user info
    db := config.NewDB()
    e := db.LoadFromFilePath(dbPath)
    if e != nil {
        return nil, &HTTPError{
            Code: http.StatusInternalServerError,
            Msg: gin.H{"msg": e.Error()},
        }
    }

    for _, user := range db.Users {
        if *user.Name == *au.Name {
            e := checkPassword(*au.Name, *user.Password, *au.Password)
            if e != nil {
                return nil, &HTTPError{
                    Code: http.StatusForbidden,
                    Msg: gin.H{"msg": e.Error()},
                }
            }

            return user, nil
        }
    }

    return nil, &HTTPError{
        Code: http.StatusNotFound,
        Msg: gin.H{"msg": "User not found"},
    }
}
