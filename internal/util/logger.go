package util

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logPath := filepath.Join(".verifier", "logs", "verifier.log")
	os.MkdirAll(filepath.Dir(logPath), 0755)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}

	// Also log to stdout for CLI feedback
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
}
