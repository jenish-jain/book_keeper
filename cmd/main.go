package main

import (
	"book_keeper/internal/config"
	"book_keeper/internal/logger"
)

func main() {
	configStore := config.InitConfig("production")
	logger.Init(configStore.GetLogLevel())
	serverDependencies, _ := InitDependencies()
	serverDependencies.server.Run(serverDependencies.handlers)
}
