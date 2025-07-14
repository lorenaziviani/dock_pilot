package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type ServiceLogger struct {
	serviceName string
	file        *os.File
	logger      *log.Logger
}

func sanitizeName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return re.ReplaceAllString(name, "_")
}

func NewServiceLogger(serviceName string, logDir string) (*ServiceLogger, error) {
	if err := os.MkdirAll(logDir, 0750); err != nil {
		return nil, err
	}
	safeName := sanitizeName(serviceName)
	logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", safeName))
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	logger := log.New(file, "", log.LstdFlags|log.Lmicroseconds)
	return &ServiceLogger{serviceName, file, logger}, nil
}

func (l *ServiceLogger) Log(msg string) {
	l.logger.Printf("[%s] %s", l.serviceName, msg)
}

func (l *ServiceLogger) Close() error {
	return l.file.Close()
}
