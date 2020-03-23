package config

import (
	"errors"
	"io/ioutil"
	"os/user"
	"strings"

	"gopkg.in/yaml.v2"
)

// DBPath stores the path db file.
var DBPath string

// User represent openvpn user info in db.
type User struct {
	Name      *string `yaml:"name" json:"name"`
	Password  *string `yaml:"password" json:"password"`
	UpdatedAt *string `yaml:"updated_at" json:"updated_at"`
	Admin     *bool   `yaml:"admin" json:"admin"`
	Ignore    *bool   `yaml:"ignore" json:"ignore"`
}

// DBConfig represent openvpn users in db.
type DBConfig struct {
	Users []*User `yaml:"users" json:"users"`
}

// NewDB create a global DB config.
func NewDB() *DBConfig{
	return &DBConfig{}
}

// LoadFromFilePath loads configuration from a specified local path for db.
// It returns error if file not found or yaml decode failed.
func (d *DBConfig) LoadFromFilePath(filePath string) error {
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

	if err := d.LoadFromContent(cYAML); err != nil {
		return err
	}

	return d.check()
}

// LoadFromContent loads configuration from a given bytes for db.
// It returns error if yaml decode failed.
func (d *DBConfig) LoadFromContent(content []byte) error {
	return yaml.Unmarshal(content, d)
}

func (d *DBConfig) check() error {
	if len(d.Users) == 0 {
		return nil
	}

	for _, user := range d.Users {
		if user.Name == nil {
			return errors.New("User name is empty in user list")
		}
		if user.Password == nil {
			return errors.New("User password is empty in user list")
		}
		if user.Admin == nil {
			return errors.New("User admin is empty in user list")
		}
		if user.Ignore == nil {
			return errors.New("User ignore is empty in user list")
		}
	}

	return nil
}
