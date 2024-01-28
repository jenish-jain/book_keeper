package config

import (
	"book_keeper/internal/mongo"
	"fmt"
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	ServerPort      int    `mapstructure:"SERVER_PORT"`
	LogLevel        string `mapstructure:"LOG_LEVEL"`
	MongoHost       string `mapstructure:"MONGO_HOST"`
	MongoRWUser     string `mapstructure:"MONGO_RW_USERNAME"`
	MongoRWPassword string `mapstructure:"MONGO_RW_PASSWORD"`
	AssetsDb        string `mapstructure:"ASSETS_DATABASE"`
}

func InitConfig(configName string) *Config {
	viper.AutomaticEnv()
	viper.SetConfigName(configName)
	viper.SetConfigType("env")
	viper.AddConfigPath("config")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("../../config/")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using Config file:", viper.ConfigFileUsed())
	} else {
		fmt.Printf("Cannot read Config file %s.env", configName)
		panic(err)
	}

	err := viper.UnmarshalExact(&AppConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error unable to Unmarshal Config file: %s", err))
	}

	return &AppConfig
}

func (c *Config) GetLogLevel() string {
	return c.LogLevel
}

func (c *Config) GetServerPort() int {
	return c.ServerPort
}

func (c *Config) GetMongoConfig() *mongo.Config {
	return mongo.NewConfig(c.MongoHost, c.MongoRWUser, c.MongoRWPassword, c.AssetsDb)
}

func GetConfig() *Config {
	return &AppConfig
}

func GetMongoConfigToInject() *mongo.Config {
	return AppConfig.GetMongoConfig()
}
