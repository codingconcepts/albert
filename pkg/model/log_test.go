package model_test

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	logrustest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/codingconcepts/pocs/chaos/harambe/pkg/model"
	"github.com/codingconcepts/pocs/chaos/harambe/test"
)

var (
	logLevel = logrus.WarnLevel
	logger   *logrus.Logger
	hook     *logrustest.Hook
)

func TestMain(m *testing.M) {
	logger = model.NewLogger(os.Stdout, logLevel)
	hook = logrustest.NewLocal(logger)

	os.Exit(m.Run())
}

func TestInit(t *testing.T) {
	test.Equals(t, logLevel, logger.Level)
}
