package repositories

import (
	"Omnilog_backend/internal/models"
	"context"
)

type InviteTokenRepository struct {
}

func NewInviteTokenRepository() *InviteTokenRepository {
	return &InviteTokenRepository{}
}

func (r *InviteTokenRepository) CreateInviteToken(ctx context.Context, db DBExecutor, token *models.InviteToken) error {
	query := `INSERT INTO invite_tokens (user_id, token_hash, expires_at, used_at, created_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := db.QueryRow(ctx, query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.UsedAt,
		token.CreatedAt,
	).Scan(&token.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *InviteTokenRepository) FindByTokenHash(ctx context.Context, db DBExecutor, tokenHash string) (*models.InviteToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, used_at, created_at
	FROM invite_tokens WHERE token_hash = $1`
	var t models.InviteToken
	err := db.QueryRow(ctx, query, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *InviteTokenRepository) MarkAsUsed(ctx context.Context, db DBExecutor, tokenID string) error {
	query := `UPDATE invite_tokens SET used_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := db.Exec(ctx, query, tokenID)
	return err
}
