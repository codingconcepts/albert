package model

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/test"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(os.Stdout, logrus.WarnLevel)
	test.Equals(t, logrus.WarnLevel, logger.Level)
	test.Equals(t, os.Stdout, logger.Out)

	utcFormatter, ok := logger.Formatter.(UTCFormatter)
	test.Assert(t, ok)

	innerFormatter, ok := utcFormatter.Formatter.(*logrus.TextFormatter)
	test.Assert(t, ok)

	test.Assert(t, innerFormatter.FullTimestamp)
	test.Assert(t, innerFormatter.DisableColors)
}

func TestFormatProvidesUTCTime(t *testing.T) {
	formatter := UTCFormatter{
		Formatter: &logrus.TextFormatter{},
	}

	london, err := time.LoadLocation("Europe/London")
	test.ErrorNil(t, err)

	bstTime, err := time.ParseInLocation(time.RFC1123, "Tue, 28 Apr 2017 07:33:01 BST", london)
	test.ErrorNil(t, err)

	logBytes, err := formatter.Format(&logrus.Entry{Time: bstTime})
	test.ErrorNil(t, err)

	log.Println(string(logBytes))
	test.Assert(t, strings.Contains(string(logBytes), `time="2017-04-28T06:33:01Z"`))
}
