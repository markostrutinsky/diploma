package repositories

import (
	"context"
	"fmt"

	"Omnilog_backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitRepository struct{}

var SingleInstanceRoles = map[models.UserRole]bool{
	models.RoleRegionDirector:    true,
	models.RoleRegionLogistician: true,
	models.RoleRegionStorekeeper: true,
	models.RoleBranchManager:     true,
	models.RoleBranchLogistician: true,
	models.RoleBranchStorekeeper: true,
	models.RoleDeptManager:       true,
	models.RoleDeptSupervisor:    true,
	models.RoleTeamLead:          true,
}

func NewUnitRepository() *UnitRepository {
	return &UnitRepository{}
}

// tenantFilter повертає SQL фрагмент "AND u.tenant_id = $N" / "WHERE u.tenant_id = $N"
// і нарощує args, якщо tenant є в контексті. Для SYSTEM_ADMIN (tenant == "") фільтр не додається.
func tenantFilter(ctx context.Context, alias, joiner string, args *[]any) string {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		return ""
	}
	*args = append(*args, tid)
	col := "tenant_id"
	if alias != "" {
		col = alias + ".tenant_id"
	}
	return fmt.Sprintf(" %s %s = $%d", joiner, col, len(*args))
}

func (r *UnitRepository) Create(ctx context.Context, db DBExecutor, u *models.Unit) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		// SYSTEM_ADMIN або нетенантний запит — беремо default tenant
		query := `INSERT INTO units (parent_id, name, unit_type) VALUES ($1, $2, $3) RETURNING id`
		return db.QueryRow(ctx, query, u.ParentID, u.Name, u.UnitType).Scan(&u.ID)
	}
	query := `INSERT INTO units (parent_id, name, unit_type, tenant_id) VALUES ($1, $2, $3, $4) RETURNING id`
	return db.QueryRow(ctx, query, u.ParentID, u.Name, u.UnitType, tid).Scan(&u.ID)
}

func (r *UnitRepository) List(ctx context.Context, db DBExecutor) ([]models.Unit, error) {
	args := []any{}
	where := tenantFilter(ctx, "", "WHERE", &args)

	// 🔍 Тимчасовий дебаг
	tid := TenantFromCtx(ctx)
	fmt.Printf("🔍 UnitRepository.List: tenant_id from context = '%s', where clause = '%s'\n", tid, where)

	q := `SELECT id, parent_id, name, unit_type FROM units` + where + ` ORDER BY id`
	rows, err := db.Query(ctx, q, args...)
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
	args := []any{unitType}
	tFilter := tenantFilter(ctx, "u", "AND", &args)

	var query string
	if SingleInstanceRoles[role] {
		args = append(args, role)
		query = fmt.Sprintf(`
			SELECT u.id, u.name
			FROM units u
			WHERE u.unit_type = $1 %s
			  AND NOT EXISTS (
				  SELECT 1 FROM users usr 
				  WHERE usr.unit_id = u.id 
					AND usr.role = $%d 
					AND usr.status IN ('ACTIVE', 'PENDING')
			  )
			ORDER BY u.name`, tFilter, len(args))
	} else {
		query = fmt.Sprintf(`
			SELECT u.id, u.name
			FROM units u
			WHERE u.unit_type = $1 %s
			ORDER BY u.name`, tFilter)
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
	args := []any{rootUnitID}
	tFilter := tenantFilter(ctx, "u", "AND", &args)
	tFilterRoot := ""
	if tid := TenantFromCtx(ctx); tid != "" {
		// root-перевірка робиться у WHERE id = $1 — рекурсія вже обмежена по parent
		tFilterRoot = fmt.Sprintf(" AND tenant_id = $%d", len(args))
	}

	query := fmt.Sprintf(`
        WITH RECURSIVE unit_tree AS (
            SELECT id, parent_id, name, unit_type
            FROM units
            WHERE id = $1 %s
            UNION ALL
            SELECT u.id, u.parent_id, u.name, u.unit_type
            FROM units u
            INNER JOIN unit_tree ut ON u.parent_id = ut.id
            WHERE 1=1 %s
        ),
        parent_peek AS (
            SELECT u.id, u.parent_id, u.name, u.unit_type
            FROM units u
            WHERE u.id = (SELECT parent_id FROM units WHERE id = $1) %s
        )
        SELECT id, parent_id, name, unit_type FROM unit_tree
        UNION
        SELECT id, parent_id, name, unit_type FROM parent_peek
        ORDER BY parent_id NULLS FIRST, name;`, tFilterRoot, tFilter, tFilter)

	rows, err := db.Query(ctx, query, args...)
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
	args := []any{id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `SELECT id, parent_id, name, unit_type FROM units WHERE id = $1` + tFilter

	var u models.Unit
	err := db.QueryRow(ctx, query, args...).Scan(&u.ID, &u.ParentID, &u.Name, &u.UnitType)
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
	args := []any{parentUnitID, targetUnitID}
	tFilter := tenantFilter(ctx, "u", "AND", &args)
	tFilterRoot := ""
	if tid := TenantFromCtx(ctx); tid != "" {
		tFilterRoot = fmt.Sprintf(" AND tenant_id = $%d", len(args))
	}
	query := fmt.Sprintf(`
    WITH RECURSIVE hierarchy AS (
        SELECT id FROM units WHERE id = $1 %s
        UNION ALL
        SELECT u.id FROM units u
        JOIN hierarchy h ON u.parent_id = h.id
        WHERE 1=1 %s
    )
    SELECT EXISTS(SELECT 1 FROM hierarchy WHERE id = $2)`, tFilterRoot, tFilter)

	var exists bool
	if err := db.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
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

	// Отримуємо unit з урахуванням tenant scoping
	args := []any{unitID}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	var unitType string
	err = tx.QueryRow(ctx, "SELECT unit_type FROM units WHERE id = $1"+tFilter, args...).Scan(&unitType)
	if err != nil {
		return fmt.Errorf("підрозділ не знайдено: %v", err)
	}

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

	// Звільняємо старого командира
	args2 := []any{unitID, targetRole}
	tFilter2 := tenantFilter(ctx, "", "AND", &args2)
	_, err = tx.Exec(ctx, `UPDATE users SET unit_id = NULL, role = 'EMPLOYEE' WHERE unit_id = $1 AND role = $2`+tFilter2, args2...)
	if err != nil {
		return fmt.Errorf("не вдалося звільнити поточного керівника: %v", err)
	}

	// Призначаємо нового керівника (перевіряємо його тенант)
	args3 := []any{unitID, targetRole, newCommanderID}
	tFilter3 := tenantFilter(ctx, "", "AND", &args3)
	_, err = tx.Exec(ctx, `UPDATE users SET unit_id = $1, role = $2 WHERE id = $3`+tFilter3, args3...)
	if err != nil {
		return fmt.Errorf("не вдалося призначити нового керівника: %v", err)
	}

	return tx.Commit(ctx)
}

func (r *UnitRepository) GetSubordinateUnitIDs(ctx context.Context, db DBExecutor, parentUnitID int64) ([]int64, error) {
	args := []any{parentUnitID}
	tFilter := tenantFilter(ctx, "u", "AND", &args)
	tFilterRoot := ""
	if tid := TenantFromCtx(ctx); tid != "" {
		tFilterRoot = fmt.Sprintf(" AND tenant_id = $%d", len(args))
	}
	query := fmt.Sprintf(`
        WITH RECURSIVE unit_tree AS (
            SELECT id FROM units WHERE id = $1 %s
            UNION ALL
            SELECT u.id FROM units u
            INNER JOIN unit_tree ut ON u.parent_id = ut.id
            WHERE 1=1 %s
        )
        SELECT id FROM unit_tree;`, tFilterRoot, tFilter)

	rows, err := db.Query(ctx, query, args...)
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
	// tenant фільтр застосовується і в root (WHERE id = $), і у рекурсивному join
	baseArgs := []any{}
	placeholder := func(v any) string {
		baseArgs = append(baseArgs, v)
		return fmt.Sprintf("$%d", len(baseArgs))
	}

	tid := TenantFromCtx(ctx)
	tRootCond := func() string {
		if tid == "" {
			return ""
		}
		return " AND tenant_id = " + placeholder(tid)
	}
	tJoinCond := func() string {
		if tid == "" {
			return ""
		}
		return " AND u.tenant_id = " + placeholder(tid)
	}

	var query string
	if SingleInstanceRoles[role] {
		typeP := placeholder(unitType)
		rootP := placeholder(commanderUnitID)
		rootT := tRootCond()
		joinT := tJoinCond()
		roleP := placeholder(role)
		query = fmt.Sprintf(`
			WITH RECURSIVE unit_tree AS (
				SELECT id, name, unit_type FROM units WHERE id = %s%s
				UNION ALL
				SELECT u.id, u.name, u.unit_type FROM units u
				INNER JOIN unit_tree ut ON u.parent_id = ut.id
				WHERE 1=1%s
			)
			SELECT id, name FROM unit_tree
			WHERE unit_type = %s
			  AND NOT EXISTS (
				  SELECT 1 FROM users usr 
				  WHERE usr.unit_id = unit_tree.id 
					AND usr.role = %s 
					AND usr.status IN ('ACTIVE', 'PENDING')
			  )
			ORDER BY name`, rootP, rootT, joinT, typeP, roleP)
	} else if unitType == "ANY" {
		rootP := placeholder(commanderUnitID)
		rootT := tRootCond()
		joinT := tJoinCond()
		query = fmt.Sprintf(`
			WITH RECURSIVE unit_tree AS (
				SELECT id, name, unit_type FROM units WHERE id = %s%s
				UNION ALL
				SELECT u.id, u.name, u.unit_type FROM units u
				INNER JOIN unit_tree ut ON u.parent_id = ut.id
				WHERE 1=1%s
			)
			SELECT id, name FROM unit_tree ORDER BY name`, rootP, rootT, joinT)
	} else {
		typeP := placeholder(unitType)
		rootP := placeholder(commanderUnitID)
		rootT := tRootCond()
		joinT := tJoinCond()
		query = fmt.Sprintf(`
			WITH RECURSIVE unit_tree AS (
				SELECT id, name, unit_type FROM units WHERE id = %s%s
				UNION ALL
				SELECT u.id, u.name, u.unit_type FROM units u
				INNER JOIN unit_tree ut ON u.parent_id = ut.id
				WHERE 1=1%s
			)
			SELECT id, name FROM unit_tree WHERE unit_type = %s ORDER BY name`, rootP, rootT, joinT, typeP)
	}

	rows, err := db.Query(ctx, query, baseArgs...)
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
	args := []any{name, id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `UPDATE units SET name = $1 WHERE id = $2` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

// Delete повністю ліквідує підрозділ з таблиці
func (r *UnitRepository) Delete(ctx context.Context, db *pgxpool.Pool, id string) error {
	args := []any{id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `DELETE FROM units WHERE id = $1` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

// GetEffectiveTier шукає тариф з таблиці tenants (якщо є tenant scope) або з units (legacy)
func (r *UnitRepository) GetEffectiveTier(ctx context.Context, db DBExecutor, unitID int) (string, error) {
	// Якщо в контексті є tenant — беремо напряму з tenants
	if tid := TenantFromCtx(ctx); tid != "" {
		var tier string
		err := db.QueryRow(ctx, `SELECT subscription_tier FROM tenants WHERE id = $1`, tid).Scan(&tier)
		if err != nil {
			return "BASIC", nil
		}
		return tier, nil
	}

	// Fallback для легасі-даних без tenant_id — піднімаємось по units
	query := `
		WITH RECURSIVE unit_hierarchy AS (
			SELECT id, parent_id, subscription_tier, 1 as depth
			FROM units WHERE id = $1
			UNION ALL
			SELECT u.id, u.parent_id, u.subscription_tier, uh.depth + 1
			FROM units u
			JOIN unit_hierarchy uh ON u.id = uh.parent_id
		)
		SELECT subscription_tier FROM unit_hierarchy 
		ORDER BY (CASE WHEN subscription_tier = 'ENTERPRISE' THEN 0 WHEN subscription_tier = 'PRO' THEN 1 WHEN subscription_tier = 'BASIC' THEN 2 ELSE 3 END) ASC
		LIMIT 1;`

	var tier string
	err := db.QueryRow(ctx, query, unitID).Scan(&tier)
	if err != nil {
		return "BASIC", nil
	}
	return tier, nil
}
