package services

import (
	"os"
	"testing"
)

func TestServiceLogger(t *testing.T) {
	logDir := "./testlogs"
	logger, err := NewServiceLogger("test-service", logDir)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	logger.Log("test message")
	logger.Close()

	logPath := logDir + "/test-service.log"
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if string(data) == "" {
		t.Errorf("log file is empty")
	}
	os.RemoveAll(logDir)
}
