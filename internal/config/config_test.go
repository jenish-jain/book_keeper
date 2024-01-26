package config_test

import (
	"book_keeper/internal/config"
	"book_keeper/internal/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
	config *config.Config
}

func (s *ConfigTestSuite) SetupTest() {
	s.config = config.InitConfig("test")
}

func (s *ConfigTestSuite) TestGetLogLevel() {
	assert.Equal(s.T(), "debug", s.config.GetLogLevel())
}

func (s *ConfigTestSuite) TestGetServerPort() {
	assert.Equal(s.T(), 8080, s.config.GetServerPort())
}

func (s *ConfigTestSuite) TestGetMongoConfig() {
	expectedMongoConfig := mongo.NewConfig("mongo.dev:70217", "read_user", "secret", "")
	assert.Equal(s.T(), expectedMongoConfig, s.config.GetMongoConfig())
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
