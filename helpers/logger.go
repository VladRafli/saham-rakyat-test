package helpers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func InitLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to determine working directory: %s", err)
	}

	f, err := os.OpenFile(filepath.Join(cwd, "logs/access.log"), os.O_APPEND|os.O_CREATE, 0755)

	if err != nil {
		logger.Fatal(err)
	}

	logger.SetOutput(io.MultiWriter(os.Stdout, f))
	return logger
}

func ApacheFormatLogger(method string, url string, host string, ip string, ua string, time string) string {
	return fmt.Sprintf("%s - [%s] %s %s %s %s %s", method, time, ip, host, url, ua, "\n")
}