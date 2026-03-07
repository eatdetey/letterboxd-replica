package app

import (
	"context"
	"time"

	"github.com/eatdetey/letterboxd-replica/source/go-common/pkg/logger"
	"github.com/eatdetey/letterboxd-replica/source/user-service/internal/config"
	"github.com/eatdetey/letterboxd-replica/source/user-service/internal/transport/postgres"
	"go.uber.org/zap"
)

type ServerApp struct {
	cfg config.ServerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient
}

func NewServerApp() (*ServerApp, error) {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		return nil, err
	}

	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	return &ServerApp{
		cfg: *cfg,
		log: log,
	}, nil
}

func (a *ServerApp) Run(ctx context.Context) error {
	if err := a.initPostgresClient(ctx); err != nil {
		return err
	}

	a.log.Infow("app.started")
	return nil
}

func (a *ServerApp) Stop(ctx context.Context) {
	_, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Shutdown.ShutdownTimeout)*time.Second)
	defer cancel()

	a.postgresClient.Close()

	a.log.Infow("app.stopped")
}

func (a *ServerApp) initPostgresClient(ctx context.Context) error {
	if a.cfg.Migrate.NeedToMigrate {
		err := postgres.Migrate(ctx, a.cfg.DB.ConnectionString)
		if err != nil {
			a.log.Errorw("app.postgres_migrate_failed", "err", err)
			return err
		}
	} else {
		a.log.Infow("app.postgres_migrate_skipped")
	}

	pgClient, err := postgres.New(ctx, a.cfg.DB.ConnectionString)
	if err != nil {
		a.log.Errorw("app.postgres_connect_failed", "err", err)
		return err
	}
	a.postgresClient = pgClient

	a.log.Infow("app.postgres_connected")
	return nil
}
