package config

import (
	"errors"
	"io/ioutil"
	"os/user"
	"strings"

	"gopkg.in/yaml.v2"
)

// Email represent server email config.
type Email struct {
	Host       *string `yaml:"host"`
	Port       *int    `yaml:"port"`
	Encryption *string `yaml:"encryption"`
	Auth       *string `yaml:"auth"`
	User       *string `yaml:"user"`
	Password   *string `yaml:"password"`
}

// WarderServer represent server config.
type WarderServer struct {
	Port        *int    `yaml:"port"`
	VPNName     *string `yaml:"vpn_name"`
	LogPath     *string `yaml:"log_path"`
	UpdateCycle *int    `yaml:"update_cycle"`
	DBPath      *string `yaml:"db_path"`
	Mail        *Email  `yaml:"email"`
}

// NewWarderServer create a DB config.
// The config need PORT and CYCLE env varible or will set it default 80 and 90.
func NewWarderServer() *WarderServer {
	return &WarderServer{}
}

// LoadFromFilePath loads configuration from a specified local path for warder server.
// It returns error if file not found or yaml decode failed.
func (w *WarderServer) LoadFromFilePath(filePath string) error {
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

	if err := w.LoadFromContent(cYAML); err != nil {
		return err
	}

	return w.check()
}

// LoadFromContent loads configuration from a given bytes for warder server.
// It returns error if yaml decode failed.
func (w *WarderServer) LoadFromContent(content []byte) error {
	return yaml.Unmarshal(content, w)
}

func (w *WarderServer) check() error {
	if w.Port == nil {
		return errors.New("Warder server port is not defined")
	}
	if w.LogPath == nil {
		return errors.New("Warder server log_path is not defined")
	}
	if w.VPNName == nil {
		return errors.New("Warder server vpn_name is not defined")
	}
	if w.UpdateCycle == nil {
		return errors.New("Warder server update_cycle is not defined")
	}
	if w.DBPath == nil {
		return errors.New("Warder server db_path is not defined")
	}
	if w.Mail == nil {
		return errors.New("Warder server mail is not defined")
	}
	if w.Mail.Host == nil {
		return errors.New("Warder server mail.host is not defined")
	}
	if w.Mail.Port == nil {
		return errors.New("Warder server mail.port is not defined")
	}
	if w.Mail.Encryption == nil {
		return errors.New("Warder server mail.encryption is not defined")
	}
	if w.Mail.Auth == nil {
		return errors.New("Warder server mail.auth is not defined")
	}
	if w.Mail.User == nil {
		return errors.New("Warder server mail.user is not defined")
	}
	if w.Mail.Password == nil {
		return errors.New("Warder server mail.password is not defined")
	}

	return nil
}
