package app

import (
	"context"
	"time"

	"github.com/eatdetey/letterboxd-replica/source/user-service/internal/config"
	"go.uber.org/zap"
)

type ServerApp struct {
	cfg config.ServerConfig
	log *zap.SugaredLogger
}

func NewServerApp() (*ServerApp, error) {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		return nil, err
	}

	return &ServerApp{
		cfg: *cfg,
	}, nil
}

func (a *ServerApp) Run(ctx context.Context) error {
	return nil
}

func (a *ServerApp) Stop(ctx context.Context) {
	_, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Shutdown.ShutdownTimeout)*time.Second)
	defer cancel()
}
