package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	MongoConfig mongoConfig
}

type mongoConfig struct {
	Host       string `mapstructure:"MONGO_HOST"`
	RWUser     string `mapstructure:"MONGO_RW_USERNAME"`
	RWPassword string `mapstructure:"MONGO_RW_PASSWORD"`
}

func InitConfig(configName string) *Config {
	viper.AutomaticEnv()
	viper.SetConfigName(configName)
	viper.SetConfigType("env")
	viper.AddConfigPath("config")
	viper.AddConfigPath("../config/")

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
