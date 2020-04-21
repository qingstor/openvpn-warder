package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/qingstor/openvpn-warder/check"
	"github.com/qingstor/openvpn-warder/config"
	"github.com/qingstor/openvpn-warder/constants"
	"github.com/qingstor/openvpn-warder/handlers"
)

// GetUser will send get user request to openvpn-warder server.
func GetUser(client *config.Client, userName string) {
	url := fmt.Sprintf("http://%s:%d/users/get?name=%s", *client.Host, *client.Port, userName)
	contentType := "application/json"

	password := generatePassword(*client.AuthUser, *client.AuthPassword)
	au := handlers.GetUserBody{
		AuthUser: &handlers.AuthUser{
			Name:     client.AuthUser,
			Password: &password,
		},
	}
	auJSON, err := json.Marshal(au)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}
	requestReader := strings.NewReader(string(auJSON[:]))
	resp, err := http.Post(url, contentType, requestReader)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}

	fmt.Println(string(body[:]))
}

// CreateUser will send create user request to openvpn-warder server.
func CreateUser(
	client *config.Client,
	name string,
	password string,
	admin bool,
	ignore bool) {
	url := fmt.Sprintf("http://%s:%d/users/create", *client.Host, *client.Port)
	contentType := "application/json"

	passwordEn := generatePassword(*client.AuthUser, *client.AuthPassword)
	au := handlers.CreateUserBody{
		AuthUser: &handlers.AuthUser{
			Name:     client.AuthUser,
			Password: &passwordEn,
		},
		User: &config.User{
			Name:     &name,
			Password: &password,
			Admin:    &admin,
			Ignore:   &ignore,
		},
	}
	auJSON, err := json.Marshal(au)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}
	requestReader := strings.NewReader(string(auJSON[:]))
	resp, err := http.Post(url, contentType, requestReader)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}

	fmt.Println(string(body[:]))
}

// DeleteUser will send get user request to openvpn-warder server.
func DeleteUser(client *config.Client, userName string) {
	url := fmt.Sprintf("http://%s:%d/users/delete?name=%s", *client.Host, *client.Port, userName)
	contentType := "application/json"

	password := generatePassword(*client.AuthUser, *client.AuthPassword)
	au := handlers.DeleteUserBody{
		AuthUser: &handlers.AuthUser{
			Name:     client.AuthUser,
			Password: &password,
		},
	}
	auJSON, err := json.Marshal(au)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}
	requestReader := strings.NewReader(string(auJSON[:]))
	resp, err := http.Post(url, contentType, requestReader)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}

	fmt.Println(string(body[:]))
}

// ResetUser will send reset user request to openvpn-warder server.
func ResetUser(
	client *config.Client,
	userName string,
	newName string,
	newPassword string,
	newAdmin bool,
	newIgnore bool) {
	url := fmt.Sprintf("http://%s:%d/users/reset?name=%s", *client.Host, *client.Port, userName)
	contentType := "application/json"

	passwordEn := generatePassword(*client.AuthUser, *client.AuthPassword)
	au := handlers.ResetUserBody{
		AuthUser: &handlers.AuthUser{
			Name:     client.AuthUser,
			Password: &passwordEn,
		},
		User: &config.User{
			Name:     &newName,
			Password: &newPassword,
			Admin:    &newAdmin,
			Ignore:   &newIgnore,
		},
	}
	auJSON, err := json.Marshal(au)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}
	requestReader := strings.NewReader(string(auJSON[:]))
	resp, err := http.Post(url, contentType, requestReader)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}

	fmt.Println(string(body[:]))
}

// VerifyUser will send get user request to openvpn-warder server and verify user.
func VerifyUser(client *config.Client, userName string, userPassowrd string) {
	url := fmt.Sprintf("http://%s:%d/users/get?name=%s", *client.Host, *client.Port, userName)
	contentType := "application/json"

	password := generatePassword(*client.AuthUser, *client.AuthPassword)
	au := handlers.GetUserBody{
		AuthUser: &handlers.AuthUser{
			Name:     client.AuthUser,
			Password: &password,
		},
	}
	auJSON, err := json.Marshal(au)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}
	requestReader := strings.NewReader(string(auJSON[:]))
	resp, err := http.Post(url, contentType, requestReader)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}

	defer resp.Body.Close()

	userInfo := GetUserResp{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		check.ErrorForExit(constants.Name, err)
	}

	if userInfo.Password == userPassowrd {
		os.Exit(0)
	}

	os.Exit(1)
}
