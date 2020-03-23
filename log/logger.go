package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is a instance.
var Logger = logrus.New()

// InitLogger will init Logger.
func InitLogger(logPath string) (err error) {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}

	Logger.Out = file

	return
}
