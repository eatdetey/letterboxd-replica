package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	pgclient "github.com/eatdetey/letterboxd-replica/source/user-service/internal/transport/postgres"
	"github.com/jackc/pgx/v5"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"

	defaultUsersLimit int32 = 100
)

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	Bio          *string
	AvatarURL    *string
	Status       string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Session struct {
	ID           int64
	UserID       int64
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByIDs(ctx context.Context, ids []int64, usernames []string, limit, offset int32) ([]User, error)
	CreateSession(ctx context.Context, session Session) (*Session, error)
	GetSessionByToken(ctx context.Context, refreshToken string) (*Session, error)
	DeleteSessionByToken(ctx context.Context, refreshToken string) error
}

type userRepository struct {
	pg *pgclient.PostgresClient
}

func NewUserRepository(pg *pgclient.PostgresClient) UserRepository {
	return &userRepository{pg: pg}
}

func (r *userRepository) CreateUser(ctx context.Context, user User) (*User, error) {
	if user.Username == "" {
		return nil, errors.New("username is required")
	}
	if user.Email == "" {
		return nil, errors.New("email is required")
	}
	if user.PasswordHash == "" {
		return nil, errors.New("password hash is required")
	}
	if user.Status == "" {
		return nil, errors.New("status is required")
	}
	if user.Role == "" {
		user.Role = RoleUser
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	now := time.Now().UTC()

	row := conn.QueryRow(ctx, createUserQuery,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Bio,
		user.AvatarURL,
		now,
		user.Status,
		string(user.Role),
	)

	created, err := scanUserRow(row)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return created, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, getUserByUsernameQuery, username)

	user, err := scanUserRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("query user by username: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByIDs(ctx context.Context, ids []int64, usernames []string, limit, offset int32) ([]User, error) {
	if len(ids) == 0 && len(usernames) == 0 {
		return nil, errors.New("ids or usernames must be provided")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	if limit <= 0 {
		limit = defaultUsersLimit
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := conn.Query(ctx, getUsersByIDsQuery, ids, usernames, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		user, err := scanUserRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, *user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *userRepository) CreateSession(ctx context.Context, session Session) (*Session, error) {
	if session.UserID == 0 {
		return nil, errors.New("user id is required")
	}
	if session.RefreshToken == "" {
		return nil, errors.New("refresh token is required")
	}
	if session.ExpiresAt.IsZero() {
		return nil, errors.New("expires_at is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	now := time.Now().UTC()

	row := conn.QueryRow(ctx, createSessionQuery, session.UserID, session.RefreshToken, session.ExpiresAt, now)

	var created Session
	if err := row.Scan(&created.ID, &created.UserID, &created.RefreshToken, &created.ExpiresAt, &created.CreatedAt); err != nil {
		return nil, fmt.Errorf("insert session: %w", err)
	}

	return &created, nil
}

func (r *userRepository) GetSessionByToken(ctx context.Context, refreshToken string) (*Session, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, getSessionByTokenQuery, refreshToken)

	var session Session
	if err := row.Scan(&session.ID, &session.UserID, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *userRepository) DeleteSessionByToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.New("refresh token is required")
	}

	conn, err := r.pg.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	tag, err := conn.Exec(ctx, deleteSessionByTokenQuery, refreshToken)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func scanUserRow(row pgx.Row) (*User, error) {
	var user User
	var bio sql.NullString
	var avatar sql.NullString

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&bio,
		&avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Status,
		&user.Role,
	); err != nil {
		return nil, err
	}

	user.Bio = nullStringToPtr(bio)
	user.AvatarURL = nullStringToPtr(avatar)

	return &user, nil
}

func nullStringToPtr(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}
	return nil
}

const createUserQuery = `
INSERT INTO users (username, email, password_hash, bio, avatar_url, created_at, updated_at, status, role)
VALUES ($1, $2, $3, $4, $5, $6, $6, $7, $8)
RETURNING id, username, email, password_hash, bio, avatar_url, created_at, updated_at, status, role;
`

const getUserByUsernameQuery = `
SELECT id, username, email, password_hash, bio, avatar_url, created_at, updated_at, status, role
FROM users
WHERE username = $1;
`

const getUsersByIDsQuery = `
SELECT id, username, email, password_hash, bio, avatar_url, created_at, updated_at, status, role
FROM users
WHERE (cardinality($1::bigint[]) > 0 AND id = ANY($1))
   OR (cardinality($2::text[]) > 0 AND username = ANY($2))
ORDER BY id
LIMIT $3
OFFSET $4;
`

const createSessionQuery = `
INSERT INTO user_sessions (user_id, refresh_token, expires_at, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, refresh_token, expires_at, created_at;
`

const getSessionByTokenQuery = `
SELECT id, user_id, refresh_token, expires_at, created_at
FROM user_sessions
WHERE refresh_token = $1;
`

const deleteSessionByTokenQuery = `
DELETE FROM user_sessions WHERE refresh_token = $1;
`
