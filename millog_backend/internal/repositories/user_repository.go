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

func (r *UserRepository) GetByEmail(ctx context.Context, db DBExecutor, email string) (*models.User, error) {
	query := `SELECT id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
	FROM users WHERE email = $1`
	var u models.User
	err := db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.FullName, &u.Phone, &u.PasswordHash,
		&u.Role, &u.Status, &u.UnitID, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, db DBExecutor, id string) (*models.User, error) {
	query := `SELECT id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
	FROM users WHERE id = $1`
	var u models.User
	err := db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.FullName, &u.Phone, &u.PasswordHash,
		&u.Role, &u.Status, &u.UnitID, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, db DBExecutor, userID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.Exec(ctx, query, passwordHash, userID)
	return err
}

func (r *UserRepository) UpdateStatus(ctx context.Context, db DBExecutor, userID string, status models.UserStatus) error {
	query := `UPDATE users SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.Exec(ctx, query, status, userID)
	return err
}

func (r *UserRepository) Count(ctx context.Context, db DBExecutor) (int, error) {
	var n int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&n)
	return n, err
}
