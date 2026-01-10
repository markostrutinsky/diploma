package repositories

import (
	"context"

	"millog_backend/internal/models"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(ctx context.Context, db DBExecutor, user *models.User) error {

	query := `INSERT INTO users (
		username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
	`
	err := db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.FullName,
		user.Phone,
		user.PasswordHash,
		user.Role,
		user.Status,
		user.UnitID,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return err
	}
	return nil
}
