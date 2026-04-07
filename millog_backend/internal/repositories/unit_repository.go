package repositories

import (
	"context"
	"fmt"

	"millog_backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitRepository struct{}

func NewUnitRepository() *UnitRepository {
	return &UnitRepository{}
}

func (r *UnitRepository) Create(ctx context.Context, db DBExecutor, u *models.Unit) error {
	query := `INSERT INTO units (parent_id, name, unit_type) VALUES ($1, $2, $3) RETURNING id`
	return db.QueryRow(ctx, query, u.ParentID, u.Name, u.UnitType).Scan(&u.ID)
}

func (r *UnitRepository) List(ctx context.Context, db DBExecutor) ([]models.Unit, error) {
	rows, err := db.Query(ctx, `SELECT id, parent_id, name, unit_type FROM units ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Unit
	for rows.Next() {
		var u models.Unit
		if err := rows.Scan(&u.ID, &u.ParentID, &u.Name, &u.UnitType); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, rows.Err()
}

func (r *UnitRepository) GetAvailableUnitsForRole(ctx context.Context, db *pgxpool.Pool, unitType string, role models.UserRole) ([]models.Unit, error) {
	query := `
        SELECT id, name
        FROM units u
        WHERE u.unit_type = $1
          AND NOT EXISTS (
              SELECT 1 
              FROM users usr 
              WHERE usr.unit_id = u.id 
                AND usr.role = $2 
                AND usr.status IN ('ACTIVE', 'PENDING')
          )
        ORDER BY u.name
    `

	rows, err := db.Query(ctx, query, unitType, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []models.Unit
	for rows.Next() {
		var u models.Unit
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		units = append(units, u)
	}

	return units, nil
}

func (r *UnitRepository) GetUnitsHierarchy(ctx context.Context, db *pgxpool.Pool, rootUnitID int64) ([]models.Unit, error) {
	query := `
        WITH RECURSIVE unit_tree AS (
            -- Базовий рівень: беремо підрозділ самого користувача
            SELECT id, parent_id, name, unit_type
            FROM units
            WHERE id = $1

            UNION ALL

            -- Рекурсія: шукаємо всі підрозділи, які підпорядковуються знайденим раніше
            SELECT u.id, u.parent_id, u.name, u.unit_type
            FROM units u
            INNER JOIN unit_tree ut ON u.parent_id = ut.id
        )
        SELECT id, parent_id, name, unit_type FROM unit_tree
        ORDER BY parent_id NULLS FIRST, name;
    `

	rows, err := db.Query(ctx, query, rootUnitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []models.Unit
	for rows.Next() {
		var u models.Unit
		if err := rows.Scan(&u.ID, &u.ParentID, &u.Name, &u.UnitType); err != nil {
			return nil, err
		}
		units = append(units, u)
	}

	return units, nil
}

func (r *UnitRepository) GetByID(ctx context.Context, db *pgxpool.Pool, id int64) (*models.Unit, error) {
	var u models.Unit

	query := `
        SELECT id, parent_id, name, unit_type 
        FROM units 
        WHERE id = $1
    `

	err := db.QueryRow(ctx, query, id).Scan(&u.ID, &u.ParentID, &u.Name, &u.UnitType)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("підрозділ з ID %d не знайдено", id)
		}
		return nil, fmt.Errorf("помилка бази даних при пошуку підрозділу: %w", err)
	}

	return &u, nil
}

func (r *UnitRepository) CheckHierarchy(ctx context.Context, db *pgxpool.Pool, parentUnitID, targetUnitID int64) (bool, error) {
	if parentUnitID == targetUnitID {
		return false, nil
	}

	query := `
    WITH RECURSIVE hierarchy AS (
        SELECT id FROM units WHERE id = $1
        UNION ALL
        SELECT u.id FROM units u
        JOIN hierarchy h ON u.parent_id = h.id
    )
    SELECT EXISTS(SELECT 1 FROM hierarchy WHERE id = $2)`

	var exists bool
	err := db.QueryRow(ctx, query, parentUnitID, targetUnitID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UnitRepository) ChangeCommanderTx(ctx context.Context, db *pgxpool.Pool, unitID int64, newCommanderID string) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var unitType string
	err = tx.QueryRow(ctx, "SELECT unit_type FROM units WHERE id = $1", unitID).Scan(&unitType)
	if err != nil {
		return fmt.Errorf("підрозділ не знайдено: %v", err)
	}

	targetRole := unitType + "_CMDR"

	_, err = tx.Exec(ctx, `
        UPDATE users SET unit_id = NULL, role = 'USER' WHERE unit_id = $1 AND role = $2
    `, unitID, targetRole)
	if err != nil {
		return fmt.Errorf("не вдалося звільнити поточного командира: %v", err)
	}

	_, err = tx.Exec(ctx, `
        UPDATE users SET unit_id = $1, role = $2 WHERE id = $3
    `, unitID, targetRole, newCommanderID)
	if err != nil {
		return fmt.Errorf("не вдалося призначити нового командира: %v", err)
	}

	return tx.Commit(ctx)
}

func (r *UnitRepository) GetSubordinateUnitIDs(ctx context.Context, db DBExecutor, parentUnitID int64) ([]int64, error) {
	query := `
        WITH RECURSIVE unit_tree AS (
            SELECT id FROM units WHERE id = $1
            UNION ALL
            SELECT u.id FROM units u
            INNER JOIN unit_tree ut ON u.parent_id = ut.id
        )
        SELECT id FROM unit_tree;
    `

	rows, err := db.Query(ctx, query, parentUnitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *UnitRepository) GetAvailableUnitsInHierarchy(ctx context.Context, db *pgxpool.Pool, unitType string, role models.UserRole, commanderUnitID int64) ([]models.Unit, error) {
	query := `
        WITH RECURSIVE unit_tree AS (
            -- Починаємо з поточного підрозділу командира
            SELECT id, name, unit_type FROM units WHERE id = $3
            UNION ALL
            -- Шукаємо всі підлеглі підрозділи
            SELECT u.id, u.name, u.unit_type FROM units u
            INNER JOIN unit_tree ut ON u.parent_id = ut.id
        )
        SELECT id, name
        FROM unit_tree
        WHERE unit_type = $1
          AND NOT EXISTS (
              -- Перевіряємо, чи немає вже активного користувача на цій посаді
              SELECT 1 FROM users usr 
              WHERE usr.unit_id = unit_tree.id 
                AND usr.role = $2 
                AND usr.status IN ('ACTIVE', 'PENDING')
          )
        ORDER BY name
    `

	rows, err := db.Query(ctx, query, unitType, role, commanderUnitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []models.Unit
	for rows.Next() {
		var u models.Unit
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		units = append(units, u)
	}

	return units, nil
}
