package model_test

import (
	"encoding/json"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/pkg/model"
)

func TestConfigLogLevel(t *testing.T) {
	var s struct {
		LogLevel model.ConfigLogLevel `json:"logLevel"`
	}

	j := `{ "logLevel": "info" }`

	err := json.Unmarshal([]byte(j), &s)
	if err != nil {
		t.Fatal(err)
	}

	if s.LogLevel.Level != logrus.InfoLevel {
		t.Fatalf("expected info but got %v", s.LogLevel.Level)
	}
}
