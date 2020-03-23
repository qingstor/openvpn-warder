package config

import (
	"errors"
	"io/ioutil"
	"os/user"
	"strings"

	"gopkg.in/yaml.v2"
)

// WarderServer represent server config.
type WarderServer struct {
	Port        *int    `yaml:"port"`
	UpdateCycle *int    `yaml:"update_cycle"`
	DBPath      *string `yaml:"db_path"`
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
	if w.UpdateCycle == nil {
		return errors.New("Warder server update_cycle is not defined")
	}
	if w.DBPath == nil {
		return errors.New("Warder server db_path is not defined")
	}

	return nil
}
