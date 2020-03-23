package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
            Name: client.AuthUser,
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
