package app

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/eatdetey/letterboxd-replica/source/go-common/pkg/logger"
	grpcm "github.com/eatdetey/letterboxd-replica/source/go-common/pkg/middleware/grpc/server"
	movielistpb "github.com/eatdetey/letterboxd-replica/source/movie-list-service/gen/go/movielist/v1"
	"github.com/eatdetey/letterboxd-replica/source/movie-list-service/internal/config"
	"github.com/eatdetey/letterboxd-replica/source/movie-list-service/internal/handler"
	"github.com/eatdetey/letterboxd-replica/source/movie-list-service/internal/repository"
	"github.com/eatdetey/letterboxd-replica/source/movie-list-service/internal/service"
	"github.com/eatdetey/letterboxd-replica/source/movie-list-service/internal/transport/postgres"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type ServerApp struct {
	cfg config.ServerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient
	grpcServer     *grpc.Server
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

	if err := a.initGRPCServer(ctx); err != nil {
		return err
	}

	a.log.Infow("app.started")
	return nil
}

func (a *ServerApp) Stop(ctx context.Context) {
	_, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Shutdown.ShutdownTimeout)*time.Second)
	defer cancel()

	a.postgresClient.Close()
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}

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

func (a *ServerApp) initGRPCServer(ctx context.Context) error {
	protectedMethods := func(string) bool { return true }

	a.log.Infow("DEBUG: auth_access_secret", "secret", a.cfg.Auth.AccessSecret)

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpcm.RequestIdMiddleware(),
			grpcm.AuthMiddleware([]byte(a.cfg.Auth.AccessSecret), protectedMethods),
			grpcm.RecoveryInterceptor(a.log),
			grpcm.LogMiddleware(a.log),
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     time.Duration(a.cfg.GRPCServer.MaxConnectionIdle) * time.Second,
			MaxConnectionAge:      time.Duration(a.cfg.GRPCServer.MaxConnectionAge) * time.Second,
			MaxConnectionAgeGrace: time.Duration(a.cfg.GRPCServer.MaxConnectionAgeGrace) * time.Second,
			Time:                  time.Duration(a.cfg.GRPCServer.KeepaliveTime) * time.Second,
			Timeout:               time.Duration(a.cfg.GRPCServer.KeepaliveTimeout) * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: a.cfg.GRPCServer.PermitWithoutStream,
		}),
	}

	a.grpcServer = grpc.NewServer(opts...)

	repo := repository.NewPlaylistRepository(a.postgresClient)
	movielistService := service.NewMovieListService(repo, a.log)
	handler := handler.NewMovieListHandler(movielistService)
	movielistpb.RegisterMovieListServiceServer(a.grpcServer, handler)

	lis, err := net.Listen("tcp", a.cfg.GRPCServer.Port)
	if err != nil {
		return status.Errorf(codes.Internal, "listen: %v", err)
	}

	go func() {
		if err := a.grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			a.log.Errorw("app.grpc_serve_failed", "err", err)
		}
	}()

	a.log.Infow("app.grpc_started", "port", a.cfg.GRPCServer.Port)
	return nil
}
