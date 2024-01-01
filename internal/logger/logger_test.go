package logger_test

import (
	"book_keeper/internal/logger"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LoggerTestSuite struct {
	suite.Suite
}

func (k *LoggerTestSuite) SetupTest() {
}

func (k *LoggerTestSuite) TestShouldPrintLogs() {
	log := logger.Init("debug")

	msg := "I am debug"
	fmt.Println("---Debug---")
	log.Debug(msg)
	log.Info(msg)
	log.Warn(msg)
	log.Error(msg)

	msg = "I am info"
	log = logger.Init("info")
	fmt.Println("---Info---")
	log.Debug(msg)
	log.Info(msg)
	log.Warn(msg)
	log.Error(msg)

	msg = "I am warn"
	log = logger.Init("warn")
	fmt.Println("---Warn---")
	log.Debug(msg)
	log.Info(msg)
	log.Warn(msg)
	log.Error(msg)

	msg = "I am error"
	log = logger.Init("error")
	fmt.Println("---Error---")
	log.Debug(msg)
	log.Info(msg)
	log.Warn(msg)
	log.Error(msg)
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
