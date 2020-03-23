package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qingstor/openvpn-warder/config"
)

// HandleGetUser will handle getuser API.
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

	db := config.NewDB()
	err = db.LoadFromFilePath(config.DBPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	for _, user := range db.Users {
		if *user.Name == userName {
			c.JSON(http.StatusOK, user)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("User %s is not found", userName)})
	return
}