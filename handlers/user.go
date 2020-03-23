package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qingstor/openvpn-warder/config"
)

// HandleGetUser will handle get user API.
func HandleGetUser(c *gin.Context) {
	getUserBody := GetUserBody{}
	err := c.ShouldBind(&getUserBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("User body error %s", err.Error())})
		return
	}

	au, e := parseAuthUser(getUserBody.AuthUser, config.DBPath)
	if e != nil {
		c.JSON(e.Code, e.Msg)
		return
	}
	if !*au.Admin {
		c.JSON(http.StatusForbidden, gin.H{"msg": "You do not have enough permission to access the source"})
		return
	}

	userName, exist := c.GetQuery("name")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Name in query is not defined"})
		return
	}

	user, err := getUserFromDB(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("User %s is not found", userName)})
		return
	}

	c.JSON(http.StatusOK, user)
	return
}

// HandleCreateUser will handle create user API.
func HandleCreateUser(c *gin.Context) {
	createUserBody := CreateUserBody{}
	err := c.ShouldBind(&createUserBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("User body error %s", err.Error())})
		return
	}

	au, e := parseAuthUser(createUserBody.AuthUser, config.DBPath)
	if e != nil {
		c.JSON(e.Code, e.Msg)
		return
	}
	if !*au.Admin {
		c.JSON(http.StatusForbidden, gin.H{"msg": "You do not have enough permission to access the source"})
		return
	}

	u := createUserBody.User
	if u.Name == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "User name is emtpy"})
		return
	}
	if u.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "User password is emtpy"})
		return
	}
	if u.Ignore == nil {
		ignore := false
		u.Ignore = &ignore
	}
	if u.Admin == nil {
		admin := false
		u.Admin = &admin
	}

	user, err := getUserFromDB(*u.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if user != nil {
		c.JSON(http.StatusConflict, gin.H{"msg": "User already existed"})
		return
	}

	timeNow := time.Now()
	timeNowStr := timeNow.Format("2006-01-02")
	u.UpdatedAt = &timeNowStr

	db, err := readEntireDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	db.Users = append(db.Users, u)
	err = writeEntireDB(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"msg":  fmt.Sprintf("User %s created", *u.Name),
		"user": *u})
	return
}

// HandleDeleteUser will handle delete user API.
func HandleDeleteUser(c *gin.Context) {
	deleteUserBody := DeleteUserBody{}
	err := c.ShouldBind(&deleteUserBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("User body error %s", err.Error())})
		return
	}

	au, e := parseAuthUser(deleteUserBody.AuthUser, config.DBPath)
	if e != nil {
		c.JSON(e.Code, e.Msg)
		return
	}
	if !*au.Admin {
		c.JSON(http.StatusForbidden, gin.H{"msg": "You do not have enough permission to access the source"})
		return
	}

	userName, exist := c.GetQuery("name")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Name in query is not defined"})
		return
	}

	db, err := readEntireDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	for index, user := range db.Users {
		if *user.Name == userName {
			db.Users = append(db.Users[:index], db.Users[index+1:]...)
			err = writeEntireDB(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"msg": fmt.Sprintf("User %s is deleted", *user.Name)})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("User %s is not found", userName)})
	return
}

// HandleResetUser will handle getuser API.
func HandleResetUser(c *gin.Context) {
	resetUserBody := ResetUserBody{}
	err := c.ShouldBind(&resetUserBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("User body error %s", err.Error())})
		return
	}

	au, e := parseAuthUser(resetUserBody.AuthUser, config.DBPath)
	if e != nil {
		c.JSON(e.Code, e.Msg)
		return
	}
	if !*au.Admin {
		c.JSON(http.StatusForbidden, gin.H{"msg": "You do not have enough permission to access the source"})
		return
	}

	userName, exist := c.GetQuery("name")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Name in query is not defined"})
		return
	}
	u := resetUserBody.User

	db, err := readEntireDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	for index, user := range db.Users {
		if *user.Name == userName {
			if *u.Name != "" {
				db.Users[index].Name = u.Name
			}
			if *u.Password != "" {
				timeNow := time.Now()
				timeNowStr := timeNow.Format("2006-01-02")
				db.Users[index].UpdatedAt = &timeNowStr
				db.Users[index].Password = u.Password
			}
			db.Users[index].Admin = u.Admin
			db.Users[index].Ignore = u.Admin

			err = writeEntireDB(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"msg":  fmt.Sprintf("User %s reseted", *u.Name),
				"user": *db.Users[index]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("User %s is not found", userName)})
	return
}
