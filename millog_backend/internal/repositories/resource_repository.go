package repositories

import (
	"context"
	"errors"
	"fmt"

	"millog_backend/internal/models"
)

type ResourceRepository struct{}

func NewResourceRepository() *ResourceRepository {
	return &ResourceRepository{}
}

func (r *ResourceRepository) Create(ctx context.Context, db DBExecutor, res *models.Resource) error {
	query := `INSERT INTO resources (category_id, unit_id, name, description, quantity, unit_type, serial_number, warehouse_id, condition, min_quantity)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at, updated_at`

	// БЕЗПЕКА: Якщо вказівник існує, але вказує на порожній рядок (""),
	// робимо його nil. Тоді база даних безпечно запише SQL NULL.
	if res.WarehouseID != nil && *res.WarehouseID == "" {
		res.WarehouseID = nil
	}

	return db.QueryRow(ctx, query,
		res.CategoryID, res.UnitID, res.Name, res.Description, res.Quantity, res.UnitType, res.SerialNumber,
		res.WarehouseID, res.Condition, res.MinQuantity,
	).Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt)
}

func (r *ResourceRepository) GetByID(ctx context.Context, db DBExecutor, id string) (*models.Resource, error) {
	// ДОДАНО: LEFT JOIN users та вибірка r.assigned_to_user_id, u.full_name
	query := `
		SELECT r.id, r.category_id, COALESCE(r.unit_id, 0), r.name, COALESCE(r.description, ''), r.quantity, COALESCE(r.unit_type, 'PCS'), 
			   COALESCE(r.serial_number, ''), COALESCE(CAST(r.warehouse_id AS TEXT), ''), r.condition, r.min_quantity, r.created_at, r.updated_at,
			   r.assigned_to_user_id, u.full_name
		FROM resources r
		LEFT JOIN users u ON r.assigned_to_user_id = u.id
		WHERE r.id = $1`

	var res models.Resource
	err := db.QueryRow(ctx, query, id).Scan(
		&res.ID, &res.CategoryID, &res.UnitID, &res.Name, &res.Description, &res.Quantity, &res.UnitType,
		&res.SerialNumber, &res.WarehouseID, &res.Condition, &res.MinQuantity,
		&res.CreatedAt, &res.UpdatedAt,
		&res.AssignedToUserID, &res.AssignedToUserName, // ДОДАНО СЮДИ
	)
	if err != nil {
		fmt.Println("🚨 SCAN ERROR IN GetByID:", err)
		return nil, err
	}
	return &res, nil
}

func (r *ResourceRepository) List(ctx context.Context, db DBExecutor, unitID *int64) ([]models.Resource, error) {
	var query string
	args := []interface{}{}

	// ДОДАНО: LEFT JOIN users та вибірка r.assigned_to_user_id, u.full_name в обох варіантах запиту
	if unitID != nil {
		query = `
        WITH RECURSIVE unit_hierarchy AS (
            SELECT id FROM units WHERE id = $1
            UNION ALL
            SELECT u.id FROM units u
            JOIN unit_hierarchy uh ON u.parent_id = uh.id
        )
        SELECT r.id, r.category_id, COALESCE(r.unit_id, 0), r.name, COALESCE(r.description, ''), r.quantity, COALESCE(r.unit_type, 'PCS'),
               COALESCE(r.serial_number, ''), COALESCE(CAST(r.warehouse_id AS TEXT), ''), r.condition, r.min_quantity, r.created_at, r.updated_at,
			   r.assigned_to_user_id, u.full_name
        FROM resources r
		LEFT JOIN users u ON r.assigned_to_user_id = u.id
        WHERE r.unit_id IN (SELECT id FROM unit_hierarchy)
        ORDER BY r.name`

		args = append(args, *unitID)
	} else {
		query = `
        SELECT r.id, r.category_id, COALESCE(r.unit_id, 0), r.name, COALESCE(r.description, ''), r.quantity, COALESCE(r.unit_type, 'PCS'),
               COALESCE(r.serial_number, ''), COALESCE(CAST(r.warehouse_id AS TEXT), ''), r.condition, r.min_quantity, r.created_at, r.updated_at,
			   r.assigned_to_user_id, u.full_name
        FROM resources r
		LEFT JOIN users u ON r.assigned_to_user_id = u.id
        ORDER BY r.name`
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		fmt.Println("🚨 QUERY ERROR IN List:", err)
		return nil, err
	}
	defer rows.Close()

	var list []models.Resource
	for rows.Next() {
		var res models.Resource
		if err := rows.Scan(
			&res.ID, &res.CategoryID, &res.UnitID, &res.Name, &res.Description, &res.Quantity, &res.UnitType,
			&res.SerialNumber, &res.WarehouseID, &res.Condition, &res.MinQuantity,
			&res.CreatedAt, &res.UpdatedAt,
			&res.AssignedToUserID, &res.AssignedToUserName, // ДОДАНО СЮДИ
		); err != nil {
			fmt.Println("🚨 SCAN ERROR IN ListResources (Row Loop):", err)
			return nil, err
		}
		list = append(list, res)
	}
	return list, rows.Err()
}

func (r *ResourceRepository) UpdateQuantity(ctx context.Context, db DBExecutor, id string, quantity int) error {
	_, err := db.Exec(ctx, `UPDATE resources SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, quantity, id)
	return err
}

func (r *ResourceRepository) WriteOff(ctx context.Context, db DBExecutor, id string, quantity int) error {
	orig, err := r.GetByID(ctx, db, id)
	if err != nil {
		return err
	}

	if quantity > orig.Quantity {
		return errors.New("спроба списати більше, ніж є в наявності")
	}

	if quantity == orig.Quantity {
		// Повне списання: просто змінюємо статус поточного запису
		_, err = db.Exec(ctx, `UPDATE resources SET condition = 'WRITTEN_OFF', updated_at = CURRENT_TIMESTAMP WHERE id = $1`, id)
		return err
	}

	// 1. Часткове списання: віднімаємо від поточного активного запису
	_, err = db.Exec(ctx, `UPDATE resources SET quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, quantity, id)
	if err != nil {
		return err
	}

	// БЕЗПЕКА ДЛЯ UUID
	var wID interface{} = orig.WarehouseID
	if orig.WarehouseID != nil && *orig.WarehouseID == "" {
		wID = nil
	}

	// 2. РОЗУМНЕ ОБ'ЄДНАННЯ: Шукаємо, чи є вже такий самий списаний ресурс
	updateWOQuery := `
        UPDATE resources 
        SET quantity = quantity + $1, updated_at = CURRENT_TIMESTAMP 
        WHERE category_id = $2 AND unit_id = $3 AND name = $4 AND condition = 'WRITTEN_OFF'
        AND warehouse_id IS NOT DISTINCT FROM $5
    `
	result, err := db.Exec(ctx, updateWOQuery, quantity, orig.CategoryID, orig.UnitID, orig.Name, wID)
	if err != nil {
		return err
	}

	// 3. Якщо такого списаного ресурсу ще не було (база оновила 0 рядків), створюємо новий
	if result.RowsAffected() == 0 {
		insertQuery := `
            INSERT INTO resources (category_id, unit_id, name, description, quantity, unit_type, serial_number, warehouse_id, condition, min_quantity)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'WRITTEN_OFF', $9)`

		_, err = db.Exec(ctx, insertQuery,
			orig.CategoryID, orig.UnitID, orig.Name, orig.Description, quantity, orig.UnitType, orig.SerialNumber,
			wID, orig.MinQuantity,
		)
		return err
	}

	return nil
}

func (r *ResourceRepository) Update(ctx context.Context, db DBExecutor, id string, req models.UpdateResourceRequest) error {
	query := "UPDATE resources SET updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}
	argID := 1

	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argID)
		args = append(args, *req.Name)
		argID++
	}
	if req.Description != nil {
		query += fmt.Sprintf(", description = $%d", argID)
		args = append(args, *req.Description)
		argID++
	}
	if req.Quantity != nil {
		query += fmt.Sprintf(", quantity = $%d", argID)
		args = append(args, *req.Quantity)
		argID++
	}
	if req.SerialNumber != nil {
		query += fmt.Sprintf(", serial_number = $%d", argID)
		args = append(args, *req.SerialNumber)
		argID++
	}

	if req.WarehouseID != nil {
		query += fmt.Sprintf(", warehouse_id = $%d", argID)
		if *req.WarehouseID == "" {
			args = append(args, nil)
		} else {
			args = append(args, *req.WarehouseID)
		}
		argID++
	}

	if req.Condition != nil {
		query += fmt.Sprintf(", condition = $%d", argID)
		args = append(args, *req.Condition)
		argID++
	}
	if req.MinQuantity != nil {
		query += fmt.Sprintf(", min_quantity = $%d", argID)
		args = append(args, *req.MinQuantity)
		argID++
	}

	if len(args) == 0 {
		return nil
	}

	query += fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, id)

	result, err := db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("database exec error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("resource not found")
	}

	return nil
}

func (r *ResourceRepository) Transfer(ctx context.Context, db DBExecutor, resourceID string, req models.TransferResourceRequest) error {
	orig, err := r.GetByID(ctx, db, resourceID)
	if err != nil {
		return fmt.Errorf("помилка отримання ресурсу: %w", err)
	}

	if orig.Quantity < req.Quantity {
		return errors.New("недостатньо майна для переміщення")
	}

	newOrigQty := orig.Quantity - req.Quantity
	updateQuery := `UPDATE resources SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err = db.Exec(ctx, updateQuery, newOrigQty, resourceID)
	if err != nil {
		return fmt.Errorf("помилка оновлення залишку: %w", err)
	}

	finalUnitID := orig.UnitID
	if req.TargetUnitID != nil && *req.TargetUnitID != 0 {
		finalUnitID = *req.TargetUnitID
	}

	var finalWarehouseID *string = req.TargetWarehouseID
	if finalWarehouseID != nil && *finalWarehouseID == "" {
		finalWarehouseID = nil
	}

	insertQuery := `
        INSERT INTO resources (category_id, unit_id, name, description, quantity, unit_type, serial_number, warehouse_id, condition, min_quantity)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = db.Exec(ctx, insertQuery,
		orig.CategoryID, finalUnitID, orig.Name, orig.Description, req.Quantity, orig.UnitType, orig.SerialNumber,
		finalWarehouseID, orig.Condition, orig.MinQuantity,
	)
	if err != nil {
		return fmt.Errorf("помилка створення переміщеного запису: %w", err)
	}

	return nil
}

func (r *ResourceRepository) Delete(ctx context.Context, db DBExecutor, id string) error {
	query := `DELETE FROM resources WHERE id = $1`
	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("ресурс не знайдено")
	}
	return nil
}

// ==========================================
// НОВИЙ МЕТОД: ВИДАЧА МАЙНА КОРИСТУВАЧУ
// ==========================================
func (r *ResourceRepository) Assign(ctx context.Context, db DBExecutor, id string, userID string, quantity int) error {
	orig, err := r.GetByID(ctx, db, id)
	if err != nil {
		return err
	}

	if quantity > orig.Quantity {
		return errors.New("спроба видати більше, ніж є в наявності")
	}

	if quantity == orig.Quantity {
		// Видаємо всю кількість: змінюємо власника і забираємо зі складу (warehouse_id = NULL)
		_, err = db.Exec(ctx, `UPDATE resources SET assigned_to_user_id = $1, warehouse_id = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, userID, id)
		return err
	}

	// 1. Часткова видача: віднімаємо від поточного активного запису
	_, err = db.Exec(ctx, `UPDATE resources SET quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, quantity, id)
	if err != nil {
		return err
	}

	// 2. РОЗУМНЕ ОБ'ЄДНАННЯ: Шукаємо, чи вже є таке майно, видане ЦІЙ ЖЕ людині
	updateQuery := `
		UPDATE resources 
		SET quantity = quantity + $1, updated_at = CURRENT_TIMESTAMP 
		WHERE category_id = $2 
		  AND name = $3 
		  AND condition = $4
		  AND unit_id IS NOT DISTINCT FROM $5
		  AND serial_number IS NOT DISTINCT FROM $6
		  AND unit_type IS NOT DISTINCT FROM $7
		  AND assigned_to_user_id = $8
	`
	result, err := db.Exec(ctx, updateQuery,
		quantity, orig.CategoryID, orig.Name, orig.Condition, orig.UnitID, orig.SerialNumber, orig.UnitType, userID,
	)
	if err != nil {
		return err
	}

	// 3. Якщо у цієї людини такого ще немає — створюємо новий рядок
	if result.RowsAffected() == 0 {
		insertQuery := `
			INSERT INTO resources (category_id, unit_id, name, description, quantity, unit_type, serial_number, warehouse_id, condition, min_quantity, assigned_to_user_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NULL, $8, $9, $10)`

		_, err = db.Exec(ctx, insertQuery,
			orig.CategoryID, orig.UnitID, orig.Name, orig.Description, quantity, orig.UnitType, orig.SerialNumber,
			orig.Condition, orig.MinQuantity, userID,
		)
		return err
	}

	return nil
}
