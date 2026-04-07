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

// GetVisibleUsers повертає список користувачів з урахуванням ієрархії підрозділів
func (r *UserRepository) GetVisibleUsers(ctx context.Context, db DBExecutor, requesterRole string, requesterUnitID *int64) ([]*models.User, error) {
	var query string
	var args []interface{}

	if requesterRole == "ADMIN" {
		query = `
			SELECT id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
			FROM users
			ORDER BY created_at DESC
		`
	} else {
		// Якщо це командир без підрозділу (в резерві), він бачить тільки резерв
		if requesterUnitID == nil {
			query = `
				SELECT id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
				FROM users
				WHERE unit_id IS NULL
				ORDER BY created_at DESC
			`
		} else {
			// CTE запит для пошуку всього дерева підрозділів (батальйони, роти, взводи під цією бригадою)
			query = `
				WITH RECURSIVE hierarchy AS (
					SELECT id FROM units WHERE id = $1
					UNION ALL
					SELECT u.id FROM units u
					JOIN hierarchy h ON u.parent_id = h.id
				)
				SELECT id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
				FROM users
				WHERE unit_id IN (SELECT id FROM hierarchy)
				   OR unit_id IS NULL -- Також показуємо кадровий резерв
				ORDER BY created_at DESC
			`
			args = append(args, *requesterUnitID)
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
	query := `
        SELECT id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at
        FROM users
        WHERE role IN ('BRIGADE_CMDR', 'BATTALION_CMDR', 'COMPANY_CMDR', 'PLATOON_CMDR')
    `

	rows, err := db.Query(ctx, query)
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

func (r *UserRepository) UpdateRoleAndUnit(ctx context.Context, db DBExecutor, userID string, newRole string, newUnitID *int64) error {
	query := `
		UPDATE users 
		SET role = $1, unit_id = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3
	`

	_, err := db.Exec(ctx, query, newRole, newUnitID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) BlockUser(ctx context.Context, db DBExecutor, userID string) error {
	query := `
		UPDATE users 
		SET status = $1, 
		    unit_id = NULL, 
		    updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2
	`
	_, err := db.Exec(ctx, query, models.StatusBlocked, userID)
	return err
}

func (r *UserRepository) CheckSubordination(ctx context.Context, db DBExecutor, commanderUnitID int64, targetUserID string) (bool, error) {
	query := `
		WITH RECURSIVE unit_tree AS (
			-- Беремо підрозділ командира
			SELECT id FROM units WHERE id = $1
			UNION
			-- Шукаємо всі підрозділи, які йому підпорядковуються
			SELECT u.id FROM units u
			INNER JOIN unit_tree ut ON u.parent_id = ut.id
		)
		SELECT EXISTS (
			SELECT 1 FROM users 
			WHERE id = $2 AND unit_id IN (SELECT id FROM unit_tree)
		)
	`
	var isSubordinate bool
	err := db.QueryRow(ctx, query, commanderUnitID, targetUserID).Scan(&isSubordinate)
	return isSubordinate, err
}

func (r *UserRepository) UnblockUser(ctx context.Context, db DBExecutor, userID string) error {
	query := `
		UPDATE users 
		SET status = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2
	`
	_, err := db.Exec(ctx, query, models.StatusActive, userID)
	return err
}
