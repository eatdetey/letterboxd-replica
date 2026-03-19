package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, connString string) error {
	migrationsDir := "./migrations"

	// Try to find migrations directory in different locations
	possiblePaths := []string{
		migrationsDir,
		filepath.Join("movie-list-service", migrationsDir),
		filepath.Join("source", "movie-list-service", migrationsDir),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			migrationsDir = path
			break
		}
	}

	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return fmt.Errorf("open db connection: %w", err)
	}
	defer db.Close()

	if err := goose.UpContext(ctx, db, migrationsDir); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
