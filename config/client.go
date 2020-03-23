package config

import (
	"errors"
	"io/ioutil"
	"os/user"
	"strings"

	"gopkg.in/yaml.v2"
)

// Client represent openvpn-warder client info.
type Client struct {
	Host         *string `yaml:"host" json:"host"`
	Port         *int    `yaml:"port" json:"port"`
	AuthUser     *string `yaml:"auth_name" json:"auth_name"`
	AuthPassword *string `yaml:"auth_password" json:"auth_password"`
}

// NewClient create a global DB config.
func NewClient() *Client {
	return &Client{}
}

// LoadFromFilePath loads configuration from a specified local path for db.
// It returns error if file not found or yaml decode failed.
func (c *Client) LoadFromFilePath(filePath string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	if strings.Index(filePath, "~/") == 0 {
		filePath = strings.Replace(filePath, "~/", usr.HomeDir+"/", 1)
	}

	cYAML, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := c.LoadFromContent(cYAML); err != nil {
		return err
	}

	return c.check()
}

// LoadFromContent loads configuration from a given bytes for db.
// It returns error if yaml decode failed.
func (c *Client) LoadFromContent(content []byte) error {
	return yaml.Unmarshal(content, c)
}

func (c *Client) check() error {
	if c.Host == nil {
		return errors.New("Client host is not defined")
	}
	if c.Port == nil {
		return errors.New("Client port is not defined")
	}
	if c.AuthUser == nil {
		return errors.New("Client auth_user is not defined")
	}
	if c.AuthPassword == nil {
		return errors.New("Client auth_password is not defined")
	}

	return nil
}
