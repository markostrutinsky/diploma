package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"Omnilog_backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(ctx context.Context, db DBExecutor, user *models.User) error {

	query := `INSERT INTO users (
		tenant_id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id
	`
	err := db.QueryRow(ctx, query,
		user.TenantID,
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

// GetVisibleUsers повертає список користувачів з урахуванням ієрархії підрозділів
func (r *UserRepository) GetVisibleUsers(ctx context.Context, db DBExecutor, requesterRole string, requesterUnitID *int64) ([]*models.User, error) {
	var query string
	var args []interface{}

	if requesterRole == "SYSTEM_ADMIN" {
		// Платформний адмін бачить всіх
		query = `SELECT id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
			FROM users ORDER BY created_at DESC`
	} else if requesterRole == "ADMIN" || requesterRole == "TENANT_ADMIN" {
		tFilter := tenantFilter(ctx, "", "WHERE", &args)
		query = `
			SELECT id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
			FROM users` + tFilter + `
			ORDER BY created_at DESC
		`
	} else {
		if requesterUnitID == nil {
			tFilter := tenantFilter(ctx, "", "AND", &args)
			query = `
				SELECT id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
				FROM users
				WHERE unit_id IS NULL` + tFilter + `
				ORDER BY created_at DESC
			`
		} else {
			args = append(args, *requesterUnitID)
			tFilter := tenantFilter(ctx, "users", "AND", &args)
			query = `
				WITH RECURSIVE hierarchy AS (
					SELECT id FROM units WHERE id = $1
					UNION ALL
					SELECT u.id FROM units u
					JOIN hierarchy h ON u.parent_id = h.id
				)
				SELECT id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
				FROM users
				WHERE (unit_id IN (SELECT id FROM hierarchy) OR unit_id IS NULL)` + tFilter + `
				ORDER BY created_at DESC
			`
		}
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.FullName, &u.Phone, &u.PasswordHash,
			&u.Role, &u.Status, &u.UnitID, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, rows.Err()
}

func (r *UserRepository) GetCommanders(ctx context.Context, db DBExecutor) ([]*models.User, error) {
	args := []any{}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `
        SELECT id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
        FROM users
        WHERE role IN ('REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 'TEAM_LEAD')` + tFilter + `
    `

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commanders []*models.User

	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.FullName, &u.Phone, &u.PasswordHash,
			&u.Role, &u.Status, &u.UnitID, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		commanders = append(commanders, &u)
	}

	return commanders, rows.Err()
}

func (r *UserRepository) GetByEmail(ctx context.Context, db DBExecutor, email string) (*models.User, error) {
	query := `SELECT id, tenant_id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
	FROM users WHERE email = $1`
	var u models.User
	err := db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.TenantID, &u.Username, &u.Email, &u.FullName, &u.Phone, &u.PasswordHash,
		&u.Role, &u.Status, &u.UnitID, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, db DBExecutor, id string) (*models.User, error) {
	// Тариф беремо з tenants.subscription_tier. SYSTEM_ADMIN завжди ENTERPRISE.
	query := `
		SELECT u.id, u.tenant_id, u.username, u.email, u.full_name, u.phone, u.password_hash,
			u.role, u.status, u.unit_id, u.created_at, u.updated_at,
			COALESCE(t.subscription_tier, 'BASIC') AS effective_tier
		FROM users u
		LEFT JOIN tenants t ON t.id = u.tenant_id
		WHERE u.id = $1`
	var u models.User
	err := db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.TenantID, &u.Username, &u.Email, &u.FullName, &u.Phone, &u.PasswordHash,
		&u.Role, &u.Status, &u.UnitID, &u.CreatedAt, &u.UpdatedAt, &u.EffectiveSubscriptionTier,
	)
	if err != nil {
		return nil, err
	}
	if u.Role == models.RoleSystemAdmin {
		u.EffectiveSubscriptionTier = "ENTERPRISE"
	}
	return &u, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, db DBExecutor, userID, passwordHash string) error {
	args := []any{passwordHash, userID}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

func (r *UserRepository) UpdateStatus(ctx context.Context, db DBExecutor, userID string, status models.UserStatus) error {
	args := []any{status, userID}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `UPDATE users SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

func (r *UserRepository) Count(ctx context.Context, db DBExecutor) (int, error) {
	args := []any{}
	tFilter := tenantFilter(ctx, "", "WHERE", &args)
	var n int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM users"+tFilter, args...).Scan(&n)
	return n, err
}

func (r *UserRepository) UpdateRoleAndUnit(ctx context.Context, db DBExecutor, userID string, newRole string, newUnitID *int64) error {
	args := []any{newRole, newUnitID, userID}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `
		UPDATE users 
		SET role = $1, unit_id = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

func (r *UserRepository) BlockUser(ctx context.Context, db DBExecutor, userID string) error {
	args := []any{models.StatusBlocked, userID}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `
		UPDATE users 
		SET status = $1, 
		    unit_id = NULL, 
		    updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

func (r *UserRepository) CheckSubordination(ctx context.Context, db DBExecutor, commanderUnitID int64, targetUserID string) (bool, error) {
	args := []any{commanderUnitID, targetUserID}
	tidBase := tenantFilter(ctx, "", "AND", &args)
	tidUsers := tenantFilter(ctx, "", "AND", &args)
	query := `
		WITH RECURSIVE unit_tree AS (
			SELECT id FROM units WHERE id = $1` + tidBase + `
			UNION
			SELECT u.id FROM units u
			INNER JOIN unit_tree ut ON u.parent_id = ut.id
		)
		SELECT EXISTS (
			SELECT 1 FROM users 
			WHERE id = $2 AND unit_id IN (SELECT id FROM unit_tree)` + tidUsers + `
		)
	`
	var isSubordinate bool
	err := db.QueryRow(ctx, query, args...).Scan(&isSubordinate)
	return isSubordinate, err
}

func (r *UserRepository) UnblockUser(ctx context.Context, db DBExecutor, userID string) error {
	args := []any{models.StatusActive, userID}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `
		UPDATE users 
		SET status = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

func (r *UserRepository) UpdateProfile(ctx context.Context, db *pgxpool.Pool, userID string, req models.UpdateProfileRequest) error {
	query := "UPDATE users SET updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}
	argID := 1

	if req.FullName != nil {
		query += fmt.Sprintf(", full_name = $%d", argID)
		args = append(args, *req.FullName)
		argID++
	}
	if req.Phone != nil {
		query += fmt.Sprintf(", phone = $%d", argID)
		args = append(args, *req.Phone)
		argID++
	}
	if req.Username != nil {
		query += fmt.Sprintf(", username = $%d", argID)
		args = append(args, *req.Username)
		argID++
	}
	if req.Email != nil {
		query += fmt.Sprintf(", email = $%d", argID)
		args = append(args, *req.Email)
		argID++
	}

	if len(args) == 0 {
		return nil
	}

	query += fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, userID)
	argID++
	if tid := TenantFromCtx(ctx); tid != "" {
		query += fmt.Sprintf(" AND tenant_id = $%d", argID)
		args = append(args, tid)
	}

	result, err := db.Exec(ctx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "users_username_key") {
			return errors.New("Цей логін вже зайнятий іншим користувачем")
		}
		if strings.Contains(err.Error(), "users_email_key") {
			return errors.New("Ця пошта вже використовується")
		}
		return fmt.Errorf("помилка оновлення профілю: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("користувача не знайдено")
	}

	return nil
}
