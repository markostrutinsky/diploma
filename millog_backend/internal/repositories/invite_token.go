package repositories

import (
	"context"
	"millog_backend/internal/models"
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
