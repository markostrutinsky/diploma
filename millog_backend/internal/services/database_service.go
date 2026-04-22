package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"millog_backend/internal/repositories"
)

func NewPostgresDB(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// RLS: перед видачею connection у pool — виставляємо app.tenant_id із request-контексту.
	// Якщо tenant порожній (SYSTEM_ADMIN або фон/міграції) — reset, щоб RLS пропускав усе.
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, "SELECT set_config('app.tenant_id', '', false)")
		return err
	}
	config.AfterRelease = func(conn *pgx.Conn) bool {
		// скидаємо GUC, щоб наступний споживач не успадкував чужий tenant
		_, err := conn.Exec(context.Background(), "SELECT set_config('app.tenant_id', '', false)")
		return err == nil
	}
	config.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		tid := repositories.TenantFromCtx(ctx)
		if _, err := conn.Exec(ctx, "SELECT set_config('app.tenant_id', $1, false)", tid); err != nil {
			slog.Warn("failed to set app.tenant_id on connection", "err", err)
			return false
		}
		return true
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Successfully connected to PostgreSQL via pgxpool")

	return pool, nil
}
