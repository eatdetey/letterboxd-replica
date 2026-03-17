package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"time"

	commonjwt "github.com/eatdetey/letterboxd-replica/source/go-common/pkg/jwt"
	userpb "github.com/eatdetey/letterboxd-replica/source/user-service/gen/go/user/v1"
	"github.com/eatdetey/letterboxd-replica/source/user-service/internal/config/settings"
	"github.com/eatdetey/letterboxd-replica/source/user-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultUserStatus = "active"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer

	repo repository.UserRepository
	log  *zap.SugaredLogger
	auth settings.AuthSettings
}

func NewUserService(repo repository.UserRepository, log *zap.SugaredLogger, auth settings.AuthSettings) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
		auth: auth,
	}
}

func (s *UserService) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "username, password and email are required")
	}
	if err := validatePassword(req.Password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	hash := hashPassword(req.Password)

	user, err := s.repo.CreateUser(ctx, repository.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		Status:       defaultUserStatus,
		Role:         repository.RoleUser,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	accessToken, refreshToken, err := s.generateTokenPair(*user)
	if err != nil {
		s.log.Errorw("user_service.generate_tokens_failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	if _, err := s.repo.CreateSession(ctx, repository.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().UTC().Add(time.Duration(s.auth.RefreshTokenTTLMin) * time.Minute),
	}); err != nil {
		s.log.Errorw("user_service.create_session_failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to save session")
	}

	pbUser := convertUser(*user)

	return &userpb.RegisterResponse{
		User: pbUser,
		Tokens: &userpb.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if user.PasswordHash != hashPassword(req.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	accessToken, refreshToken, err := s.generateTokenPair(*user)
	if err != nil {
		s.log.Errorw("user_service.generate_tokens_failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	if _, err := s.repo.CreateSession(ctx, repository.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().UTC().Add(time.Duration(s.auth.RefreshTokenTTLMin) * time.Minute),
	}); err != nil {
		s.log.Errorw("user_service.create_session_failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to save session")
	}

	pbUser := convertUser(*user)

	return &userpb.LoginResponse{
		User: pbUser,
		Tokens: &userpb.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *UserService) Refresh(ctx context.Context, req *userpb.RefreshRequest) (*userpb.RefreshResponse, error) {
	if req == nil || req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	claims, err := parseToken(req.RefreshToken, []byte(s.auth.RefreshSecret))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	session, err := s.repo.GetSessionByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	now := time.Now().UTC()
	if session.ExpiresAt.Before(now) {
		_ = s.repo.DeleteSessionByToken(ctx, req.RefreshToken)
		return nil, status.Error(codes.Unauthenticated, "session expired")
	}

	if session.UserID != claims.Id {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	users, err := s.repo.GetByIDs(ctx, []int64{session.UserID}, nil, 1, 0)
	if err != nil || len(users) == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	user := users[0]

	accessToken, refreshToken, err := s.generateTokenPair(user)
	if err != nil {
		s.log.Errorw("user_service.generate_tokens_failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	if err := s.repo.DeleteSessionByToken(ctx, req.RefreshToken); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		s.log.Errorw("user_service.delete_session_failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to update session")
	}

	if _, err := s.repo.CreateSession(ctx, repository.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(time.Duration(s.auth.RefreshTokenTTLMin) * time.Minute),
	}); err != nil {
		s.log.Errorw("user_service.create_session_failed", "err", err)
		return nil, status.Error(codes.Internal, "failed to save session")
	}

	return &userpb.RefreshResponse{
		Tokens: &userpb.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *UserService) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if len(req.Ids) == 0 && len(req.Usernames) == 0 {
		return nil, status.Error(codes.InvalidArgument, "ids or usernames must be provided")
	}

	users, err := s.repo.GetByIDs(ctx, req.Ids, req.Usernames, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	respUsers := make([]*userpb.User, 0, len(users))
	seen := make(map[int64]struct{})
	for _, u := range users {
		if _, exists := seen[u.ID]; exists {
			continue
		}
		seen[u.ID] = struct{}{}
		respUsers = append(respUsers, convertUser(u))
	}

	return &userpb.GetUsersResponse{
		Users: respUsers,
	}, nil
}

func (s *UserService) generateTokenPair(user repository.User) (string, string, error) {
	now := time.Now().UTC()

	accessToken, err := buildToken(user, []byte(s.auth.AccessSecret), time.Duration(s.auth.AccessTokenTTLMin)*time.Minute, now)
	if err != nil {
		s.log.Errorw("user_service.generate_access_token_failed", "err", err)
		return "", "", err
	}

	refreshToken, err := buildToken(user, []byte(s.auth.RefreshSecret), time.Duration(s.auth.RefreshTokenTTLMin)*time.Minute, now)
	if err != nil {
		s.log.Errorw("user_service.generate_refresh_token_failed", "err", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func buildToken(user repository.User, secret []byte, ttl time.Duration, now time.Time) (string, error) {
	claims := commonjwt.Claims{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
		Roles:    []string{string(user.Role)},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func parseToken(tokenStr string, secret []byte) (*commonjwt.Claims, error) {
	return commonjwt.ParseToken(tokenStr, secret)
}

func convertUser(user repository.User) *userpb.User {
	return &userpb.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Bio:       valueOrEmpty(user.Bio),
		AvatarUrl: valueOrEmpty(user.AvatarURL),
		Status:    user.Status,
		Role:      convertRole(user.Role),
	}
}

func convertRole(role repository.Role) userpb.Role {
	switch role {
	case repository.RoleAdmin:
		return userpb.Role_ROLE_ADMIN
	case repository.RoleUser:
		return userpb.Role_ROLE_USER
	default:
		return userpb.Role_ROLE_UNSPECIFIED
	}
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", sum[:])
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	hasLetter := false
	hasDigit := false
	for _, r := range password {
		switch {
		case ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z'):
			hasLetter = true
		case '0' <= r && r <= '9':
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return fmt.Errorf("password must contain letters and digits")
	}
	return nil
}
