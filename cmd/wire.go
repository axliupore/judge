//go:build wireinject
// +build wireinject

package main

import (
	"github.com/axliupore/judge/config"
	"github.com/axliupore/judge/internal/http"
	"github.com/axliupore/judge/internal/nsq"
	"github.com/axliupore/judge/internal/ws"
	"github.com/axliupore/judge/pkg/cache"
	"github.com/google/wire"
)

func InitializeApp(conf *config.Config, c *cache.Cache) (*App, func(), error) {
	wire.Build(
		wire.Struct(new(App), "*"),
		wire.FieldsOf(new(*config.Config), "Http", "Ws", "Nsq"),
		provideHttpServer,
		provideWsServer,
		provideNsqServer,
	)
	return &App{}, nil, nil
}

func provideHttpServer(conf *config.Http, c *cache.Cache) *http.Server {
	return http.NewServer(conf, c)
}

func provideWsServer(conf *config.Ws, c *cache.Cache) *ws.Server {
	return ws.NewServer(conf, c)
}

func provideNsqServer(conf *config.Nsq, c *cache.Cache) *nsq.Server {
	return nsq.NewServer(conf, c)
}
