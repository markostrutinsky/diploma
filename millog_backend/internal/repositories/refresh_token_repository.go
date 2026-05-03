package repositories

import (
	"context"
	"time"

	"Omnilog_backend/internal/models"
)

type RefreshTokenRepository struct{}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, db DBExecutor, rt *models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(ctx, query, rt.ID, rt.UserID, rt.TokenHash, rt.ExpiresAt, rt.CreatedAt)
	return err
}

func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, db DBExecutor, hash string) (*models.RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
	FROM refresh_tokens WHERE token_hash = $1 AND revoked_at IS NULL`
	var rt models.RefreshToken
	err := db.QueryRow(ctx, query, hash).Scan(
		&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.RevokedAt, &rt.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, db DBExecutor, id string) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE id = $2`
	_, err := db.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, db DBExecutor, userID string) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`
	_, err := db.Exec(ctx, query, time.Now(), userID)
	return err
}
