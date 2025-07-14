package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type ServiceLogger struct {
	serviceName string
	file        *os.File
	logger      *log.Logger
}

func NewServiceLogger(serviceName string, logDir string) (*ServiceLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}
	logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", serviceName))
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
