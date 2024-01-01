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
	logger.Init("debug")

	msg := "I am debug"
	fmt.Println("---Debug---")
	logger.Debug(msg)
	logger.Info(msg)
	logger.Warn(msg)
	logger.Error(msg)

	msg = "I am info"
	logger.Init("info")
	fmt.Println("---Info---")
	logger.Debug(msg)
	logger.Info(msg)
	logger.Warn(msg)
	logger.Error(msg)

	msg = "I am warn"
	logger.Init("warn")
	fmt.Println("---Warn---")
	logger.Debug(msg)
	logger.Info(msg)
	logger.Warn(msg)
	logger.Error(msg)

	msg = "I am error"
	logger.Init("error")
	fmt.Println("---Error---")
	logger.Debug(msg)
	logger.Info(msg)
	logger.Warn(msg)
	logger.Error(msg)
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
