package app

import (
	"context"
	"fmt"
	"time"

	userpb "github.com/eatdetey/letterboxd-replica/source/api-gateway/gen/go/user/v1"
	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/config"
	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/handler"
	"github.com/eatdetey/letterboxd-replica/source/api-gateway/internal/router"
	"github.com/eatdetey/letterboxd-replica/source/go-common/pkg/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	cfg *config.Config
	log *zap.SugaredLogger

	app        *fiber.App
	userClient userpb.UserServiceClient
	userConn   *grpc.ClientConn
}

func New() (*Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log, err := logger.New()
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	s := &Server{
		cfg: cfg,
		log: log,
	}

	if err := s.initUserClient(context.Background()); err != nil {
		return nil, err
	}
	s.initFiber()

	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		if err := s.app.Listen(s.cfg.HTTP.Addr); err != nil {
			s.log.Errorw("gateway.listen_failed", "err", err)
		}
	}()

	s.log.Infow("gateway.started", "addr", s.cfg.HTTP.Addr)

	<-ctx.Done()
	return s.Stop(context.Background())
}

func (s *Server) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if s.app != nil {
		if err := s.app.ShutdownWithContext(shutdownCtx); err != nil {
			s.log.Errorw("gateway.shutdown_failed", "err", err)
		}
	}

	if s.userConn != nil {
		s.userConn.Close()
	}

	s.log.Infow("gateway.stopped")
	return nil
}

func (s *Server) initUserClient(ctx context.Context) error {
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		dialCtx,
		s.cfg.UserService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)),
	)
	if err != nil {
		return fmt.Errorf("dial user service: %w", err)
	}

	s.userConn = conn
	s.userClient = userpb.NewUserServiceClient(conn)
	return nil
}

func (s *Server) initFiber() {
	app := fiber.New(fiber.Config{
		Immutable: true,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{"error": e.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
		},
	})

	userHandler := handler.NewUserHandler(s.userClient)
	router.Setup(app, userHandler)

	s.app = app
}
