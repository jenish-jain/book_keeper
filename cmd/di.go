//go:build wireinject
// +build wireinject

package main

import (
	"book_keeper/internal/config"
	"book_keeper/internal/file"
	"book_keeper/internal/health"
	"book_keeper/internal/mongo"
	"book_keeper/internal/server"
	"github.com/google/wire"
)

type ServerDependencies struct {
	config   *config.Config
	server   *server.Server
	handlers server.Handlers
}

func InitDependencies() (ServerDependencies, error) {
	wire.Build(
		wire.Struct(new(ServerDependencies), "*"),
		wire.Struct(new(server.Handlers), "*"),
		config.GetConfig,
		config.GetMongoConfigToInject,
		server.WireSet,
		health.WireSet,
		file.WireSet,
		mongo.WireSet,
	)

	return ServerDependencies{}, nil
}
