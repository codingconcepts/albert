package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/codingconcepts/pocs/chaos/harambe/pkg/model"
)

func TestConfigDuration(t *testing.T) {
	var s struct {
		Duration model.ConfigDuration `json:"duration"`
	}

	j := `{ "duration": "1h30m40s" }`

	err := json.Unmarshal([]byte(j), &s)
	if err != nil {
		t.Fatal(err)
	}

	if s.Duration.Duration != time.Hour+time.Minute*30+time.Second*40 {
		t.Fatalf("expected 1h30m40s but got %v", s.Duration)
	}
}
