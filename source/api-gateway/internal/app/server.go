package app

import (
	"context"
	"fmt"
	"time"

	moviepb "github.com/eatdetey/letterboxd-replica/source/api-gateway/gen/go/movie/v1"
	reviewpb "github.com/eatdetey/letterboxd-replica/source/api-gateway/gen/go/review/v1"
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

	app          *fiber.App
	userClient   userpb.UserServiceClient
	movieClient  moviepb.MovieServiceClient
	reviewClient reviewpb.ReviewServiceClient
	userConn     *grpc.ClientConn
	movieConn    *grpc.ClientConn
	reviewConn   *grpc.ClientConn
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

	if err := s.initUserClient(); err != nil {
		return nil, err
	}
	if err := s.initMovieClient(); err != nil {
		return nil, err
	}
	if err := s.initReviewClient(); err != nil {
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
	if s.movieConn != nil {
		s.movieConn.Close()
	}
	if s.reviewConn != nil {
		s.reviewConn.Close()
	}

	s.log.Infow("gateway.stopped")
	return nil
}

func (s *Server) initUserClient() error {
	conn, err := s.dialService(s.cfg.UserService.Address)
	if err != nil {
		return fmt.Errorf("dial user service: %w", err)
	}

	s.userConn = conn
	s.userClient = userpb.NewUserServiceClient(conn)
	return nil
}

func (s *Server) initMovieClient() error {
	conn, err := s.dialService(s.cfg.MovieService.Address)
	if err != nil {
		return fmt.Errorf("dial movie service: %w", err)
	}

	s.movieConn = conn
	s.movieClient = moviepb.NewMovieServiceClient(conn)
	return nil
}

func (s *Server) initReviewClient() error {
	conn, err := s.dialService(s.cfg.ReviewService.Address)
	if err != nil {
		return fmt.Errorf("dial review service: %w", err)
	}

	s.reviewConn = conn
	s.reviewClient = reviewpb.NewReviewServiceClient(conn)
	return nil
}

func (s *Server) dialService(address string) (*grpc.ClientConn, error) {
	// Non-blocking dial: API gateway should start even if downstream services
	// are temporarily unavailable, so health/swagger remain accessible.
	return grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)),
	)
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
	movieHandler := handler.NewMovieHandler(s.movieClient)
	reviewHandler := handler.NewReviewHandler(s.reviewClient)
	router.Setup(app, userHandler, movieHandler, reviewHandler)

	s.app = app
}
