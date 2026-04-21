package repositories

import (
	"context"
	"fmt"

	"millog_backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitRepository struct{}

var SingleInstanceRoles = map[models.UserRole]bool{
	models.RoleRegionDirector: true,
	models.RoleBranchManager:  true,
	models.RoleDeptManager:    true,
	models.RoleTeamLead:       true,
}

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
	var query string
	var args []interface{}

	// Перевіряємо, чи є роль "унікальною" для підрозділу
	if SingleInstanceRoles[role] {
		query = `
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
		args = append(args, unitType, role)
	} else {
		// Якщо роль масова (наприклад, EMPLOYEE) - просто повертаємо підрозділи відповідного типу
		query = `
			SELECT id, name
			FROM units u
			WHERE u.unit_type = $1
			ORDER BY u.name
		`
		args = append(args, unitType)
	}

	rows, err := db.Query(ctx, query, args...)
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
            -- 1. Беремо поточний підрозділ
            SELECT id, parent_id, name, unit_type
            FROM units
            WHERE id = $1

            UNION ALL

            -- 2. Рекурсія ВНИЗ: шукаємо всі підлеглі підрозділи
            SELECT u.id, u.parent_id, u.name, u.unit_type
            FROM units u
            INNER JOIN unit_tree ut ON u.parent_id = ut.id
        ),
        parent_peek AS (
            -- 3. МАГІЯ: Підглядаємо ВГОРУ на одного безпосереднього батька
            SELECT u.id, u.parent_id, u.name, u.unit_type
            FROM units u
            WHERE u.id = (SELECT parent_id FROM units WHERE id = $1)
        )
        -- Об'єднуємо дерево вниз та одного батька зверху
        SELECT id, parent_id, name, unit_type FROM unit_tree
        UNION
        SELECT id, parent_id, name, unit_type FROM parent_peek
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

	// 🔥 МАГІЯ ТУТ: Правильно конвертуємо тип підрозділу в посаду
	var targetRole string
	switch unitType {
	case "REGION":
		targetRole = "REGION_DIRECTOR"
	case "BRANCH":
		targetRole = "BRANCH_MANAGER"
	case "DEPARTMENT":
		targetRole = "DEPT_MANAGER"
	case "TEAM":
		targetRole = "TEAM_LEAD"
	default:
		return fmt.Errorf("неможливо призначити керівника для типу підрозділу: %s", unitType)
	}

	// Звільняємо старого командира (переводимо у звичайні співробітники)
	_, err = tx.Exec(ctx, `
        UPDATE users SET unit_id = NULL, role = 'EMPLOYEE' WHERE unit_id = $1 AND role = $2
    `, unitID, targetRole)
	if err != nil {
		return fmt.Errorf("не вдалося звільнити поточного керівника: %v", err)
	}

	// Призначаємо нового керівника
	_, err = tx.Exec(ctx, `
        UPDATE users SET unit_id = $1, role = $2 WHERE id = $3
    `, unitID, targetRole, newCommanderID)
	if err != nil {
		return fmt.Errorf("не вдалося призначити нового керівника: %v", err)
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
	var query string
	var args []interface{}

	if SingleInstanceRoles[role] {
		query = `
			WITH RECURSIVE unit_tree AS (
				SELECT id, name, unit_type FROM units WHERE id = $3
				UNION ALL
				SELECT u.id, u.name, u.unit_type FROM units u
				INNER JOIN unit_tree ut ON u.parent_id = ut.id
			)
			SELECT id, name
			FROM unit_tree
			WHERE unit_type = $1
			  AND NOT EXISTS (
				  SELECT 1 FROM users usr 
				  WHERE usr.unit_id = unit_tree.id 
					AND usr.role = $2 
					AND usr.status IN ('ACTIVE', 'PENDING')
			  )
			ORDER BY name
		`
		args = append(args, unitType, role, commanderUnitID)
	} else {
		if unitType == "ANY" {
			// Якщо ANY, то НЕ фільтруємо по unit_type взагалі
			query = `
				WITH RECURSIVE unit_tree AS (
					SELECT id, name, unit_type FROM units WHERE id = $1
					UNION ALL
					SELECT u.id, u.name, u.unit_type FROM units u
					INNER JOIN unit_tree ut ON u.parent_id = ut.id
				)
				SELECT id, name
				FROM unit_tree
				ORDER BY name
			`
			args = append(args, commanderUnitID)
		} else {
			// Стандартна поведінка для конкретного типу
			query = `
				WITH RECURSIVE unit_tree AS (
					SELECT id, name, unit_type FROM units WHERE id = $2
					UNION ALL
					SELECT u.id, u.name, u.unit_type FROM units u
					INNER JOIN unit_tree ut ON u.parent_id = ut.id
				)
				SELECT id, name
				FROM unit_tree
				WHERE unit_type = $1
				ORDER BY name
			`
			args = append(args, unitType, commanderUnitID)
		}
	}

	rows, err := db.Query(ctx, query, args...)
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

// Update змінює назву підрозділу
func (r *UnitRepository) Update(ctx context.Context, db *pgxpool.Pool, id string, name string) error {
	query := `UPDATE units SET name = $1 WHERE id = $2`
	_, err := db.Exec(ctx, query, name, id)
	return err
}

// Delete повністю ліквідує підрозділ з таблиці
func (r *UnitRepository) Delete(ctx context.Context, db *pgxpool.Pool, id string) error {
	query := `DELETE FROM units WHERE id = $1`
	_, err := db.Exec(ctx, query, id)
	return err
}

// GetEffectiveTier шукає тариф орг. одиниці з урахуванням ієрархії
func (r *UnitRepository) GetEffectiveTier(ctx context.Context, db DBExecutor, unitID int) (string, error) {
	query := `
		WITH RECURSIVE unit_hierarchy AS (
			-- Початкова точка (наш відділ)
			SELECT id, parent_id, subscription_tier, 1 as depth
			FROM units
			WHERE id = $1
			
			UNION ALL
			
			-- Рекурсивно піднімаємося до батьків
			SELECT u.id, u.parent_id, u.subscription_tier, uh.depth + 1
			FROM units u
			JOIN unit_hierarchy uh ON u.id = uh.parent_id
		)
		-- Вибираємо PRO, якщо він є хоча б десь в ієрархії. 
		-- Сортуємо так, щоб PRO був пріоритетним над BASIC.
		SELECT subscription_tier 
		FROM unit_hierarchy 
		ORDER BY (CASE WHEN subscription_tier = 'PRO' THEN 1 WHEN subscription_tier = 'ENTERPRISE' THEN 0 ELSE 2 END) ASC
		LIMIT 1;
	`

	var tier string
	err := db.QueryRow(ctx, query, unitID).Scan(&tier)
	if err != nil {
		return "BASIC", nil // Якщо щось пішло не так, повертаємо базовий тариф
	}

	return tier, nil
}
