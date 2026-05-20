package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"Omnilog_backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResourceRepository struct{}

func NewResourceRepository() *ResourceRepository {
	return &ResourceRepository{}
}

func (r *ResourceRepository) Create(ctx context.Context, db DBExecutor, res *models.Resource) error {
	if res.WarehouseID != nil && *res.WarehouseID == "" {
		res.WarehouseID = nil
	}

	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for creating resources")
	}

	query := `INSERT INTO resources (category_id, unit_id, name, description, quantity, unit_type, serial_number, barcode, warehouse_id, condition, min_quantity, weight_kg, unit_price, tenant_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, created_at, updated_at`
	return db.QueryRow(ctx, query,
		res.CategoryID, res.UnitID, res.Name, res.Description, res.Quantity, res.UnitType, res.SerialNumber,
		res.Barcode, res.WarehouseID, res.Condition, res.MinQuantity, res.WeightKg, res.UnitPrice, tid,
	).Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt)
}

func (r *ResourceRepository) GetByID(ctx context.Context, db DBExecutor, id string) (*models.Resource, error) {
	args := []any{id}
	tFilter := tenantFilter(ctx, "r", "AND", &args)
	query := `
        SELECT 
            r.id, r.category_id, COALESCE(r.unit_id, 0), r.name, COALESCE(r.description, ''), 
            r.quantity, COALESCE(SUM(ra.quantity) FILTER (WHERE ra.status = 'ACTIVE'), 0) as issued_quantity, 
            COALESCE(r.unit_type, 'PCS'), COALESCE(r.serial_number, ''), 
            COALESCE(r.barcode, ''),
            COALESCE(CAST(r.warehouse_id AS TEXT), ''),
            COALESCE(w.name, 'Без складу') as warehouse_name,
            r.condition, r.min_quantity, 
			r.weight_kg, r.unit_price,
            r.created_at, r.updated_at
        FROM resources r
        LEFT JOIN resource_assignments ra ON r.id = ra.resource_id
        LEFT JOIN warehouses w ON r.warehouse_id = w.id
        WHERE r.id = $1` + tFilter + `
        GROUP BY r.id, w.name`

	var res models.Resource
	err := db.QueryRow(ctx, query, args...).Scan(
		&res.ID, &res.CategoryID, &res.UnitID, &res.Name, &res.Description,
		&res.Quantity, &res.IssuedQuantity, &res.UnitType, &res.SerialNumber,
		&res.Barcode,
		&res.WarehouseID, &res.WarehouseName, &res.Condition, &res.MinQuantity, &res.WeightKg, &res.UnitPrice, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *ResourceRepository) List(ctx context.Context, db DBExecutor, unitID *int64) ([]models.Resource, error) {
	var query string
	args := []interface{}{}

	if unitID != nil {
		args = append(args, *unitID)
		tFilter := tenantFilter(ctx, "r", "AND", &args)
		query = `
        WITH RECURSIVE unit_hierarchy AS (
            SELECT id FROM units WHERE id = $1
            UNION ALL
            SELECT u.id FROM units u
            JOIN unit_hierarchy uh ON u.parent_id = uh.id
        )
        SELECT 
            r.id, r.category_id, COALESCE(r.unit_id, 0), r.name, COALESCE(r.description, ''), 
            r.quantity,
            COALESCE(SUM(ra.quantity) FILTER (WHERE ra.status = 'ACTIVE'), 0) as issued_quantity,
            COALESCE(r.unit_type, 'PCS'), COALESCE(r.serial_number, ''), COALESCE(r.barcode, ''), 
            COALESCE(CAST(r.warehouse_id AS TEXT), ''),
            COALESCE(w.name, 'Без складу') as warehouse_name,
            r.condition, r.min_quantity, 
			r.weight_kg, r.unit_price,
            r.created_at, r.updated_at
        FROM resources r
        LEFT JOIN resource_assignments ra ON r.id = ra.resource_id
        LEFT JOIN warehouses w ON r.warehouse_id = w.id
        WHERE r.unit_id IN (SELECT id FROM unit_hierarchy)` + tFilter + `
        GROUP BY r.id, w.name
        ORDER BY r.name`
	} else {
		tFilter := tenantFilter(ctx, "r", "WHERE", &args)
		query = `
        SELECT 
            r.id, r.category_id, COALESCE(r.unit_id, 0), r.name, COALESCE(r.description, ''), 
            r.quantity,
            COALESCE(SUM(ra.quantity) FILTER (WHERE ra.status = 'ACTIVE'), 0) as issued_quantity,
            COALESCE(r.unit_type, 'PCS'), COALESCE(r.serial_number, ''), COALESCE(r.barcode, ''),
            COALESCE(CAST(r.warehouse_id AS TEXT), ''),
            COALESCE(w.name, 'Без складу') as warehouse_name,
            r.condition, r.min_quantity, 
			r.weight_kg, r.unit_price,
            r.created_at, r.updated_at
        FROM resources r
        LEFT JOIN resource_assignments ra ON r.id = ra.resource_id
        LEFT JOIN warehouses w ON r.warehouse_id = w.id` + tFilter + `
        GROUP BY r.id, w.name
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
			&res.ID, &res.CategoryID, &res.UnitID, &res.Name, &res.Description,
			&res.Quantity,       // Це тепер складський залишок
			&res.IssuedQuantity, // НОВЕ: Це скільки на руках
			&res.UnitType, &res.SerialNumber, &res.Barcode, &res.WarehouseID, &res.WarehouseName, &res.Condition, &res.MinQuantity, &res.WeightKg, &res.UnitPrice,
			&res.CreatedAt, &res.UpdatedAt,
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
		tid := TenantFromCtx(ctx)
		if tid == "" {
			return fmt.Errorf("tenant_id is required for creating written-off resources")
		}

		insertQuery := `
            INSERT INTO resources (category_id, unit_id, name, description, quantity, unit_type, serial_number, warehouse_id, condition, min_quantity, tenant_id)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'WRITTEN_OFF', $9, $10)`

		_, err = db.Exec(ctx, insertQuery,
			orig.CategoryID, orig.UnitID, orig.Name, orig.Description, quantity, orig.UnitType, orig.SerialNumber,
			wID, orig.MinQuantity, tid,
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

	if req.WeightKg != nil {
		query += fmt.Sprintf(", weight_kg = $%d", argID)
		args = append(args, *req.WeightKg)
		argID++
	}

	if req.UnitPrice != nil {
		query += fmt.Sprintf(", unit_price = $%d", argID)
		args = append(args, *req.UnitPrice)
		argID++
	}

	if len(args) == 0 {
		return nil
	}

	query += fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, id)
	argID++
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query += tFilter
	_ = argID

	result, err := db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("database exec error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("resource not found")
	}

	return nil
}

func (r *ResourceRepository) Delete(ctx context.Context, db DBExecutor, id string) error {
	args := []any{id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `DELETE FROM resources WHERE id = $1` + tFilter
	result, err := db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("ресурс не знайдено")
	}
	return nil
}

func (r *ResourceRepository) AssignResource(ctx context.Context, db *pgxpool.Pool, resourceID string, userID string, quantity int) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("не вдалося почати транзакцію: %w", err)
	}
	defer tx.Rollback(ctx)

	// =====================================================================
	// 🛡️ ЗОЛОТЕ ПРАВИЛО ЛОГІСТИКИ: Перевірка приналежності до одного підрозділу
	// =====================================================================
	var resourceUnitID, targetUserUnitID int

	// 1. Дізнаємось, якому підрозділу належить майно
	err = tx.QueryRow(ctx, "SELECT COALESCE(unit_id, 0) FROM resources WHERE id = $1", resourceID).Scan(&resourceUnitID)
	if err != nil {
		return errors.New("ресурс не знайдено")
	}

	// 2. Дізнаємось, у якому підрозділі служить боєць
	err = tx.QueryRow(ctx, "SELECT COALESCE(unit_id, 0) FROM users WHERE id = $1", userID).Scan(&targetUserUnitID)
	if err != nil {
		return errors.New("співробітника не знайдено")
	}

	// 3. Блокуємо видачу, якщо підрозділи різні (або хтось із них "висить у повітрі" без підрозділу)
	if resourceUnitID == 0 || targetUserUnitID == 0 || resourceUnitID != targetUserUnitID {
		return errors.New("порушення субординації: неможливо видати майно бійцю з іншого підрозділу. Використовуйте 'Трансфер' або сформуйте рейс.")
	}
	// =====================================================================

	// Далі йде стандартна логіка списання зі складу та запису
	updateQuery := `
		UPDATE resources 
		SET quantity = quantity - $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND quantity >= $1
	`

	result, err := tx.Exec(ctx, updateQuery, quantity, resourceID)
	if err != nil {
		return fmt.Errorf("помилка оновлення залишків на складі: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("недостатньо майна на складі або ресурс не знайдено")
	}

	insertQuery := `
		INSERT INTO resource_assignments (resource_id, user_id, quantity, status)
		VALUES ($1, $2, $3, 'ACTIVE')
	`

	_, err = tx.Exec(ctx, insertQuery, resourceID, userID, quantity)
	if err != nil {
		return fmt.Errorf("помилка створення запису про видачу: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("помилка фіксації транзакції: %w", err)
	}

	return nil
}

func (r *ResourceRepository) GetMyEquipment(ctx context.Context, db DBExecutor, userID string) ([]models.MyEquipmentItem, error) {
	query := `
		SELECT 
			ra.id as assignment_id,
			r.id as resource_id,
			r.name as resource_name,
			ra.quantity,
			COALESCE(r.unit_type, 'PCS') as unit_type,
			ra.issued_at,
			ra.status
		FROM resource_assignments ra
		JOIN resources r ON ra.resource_id = r.id
		WHERE ra.user_id = $1 AND ra.status = 'ACTIVE'
		ORDER BY ra.issued_at DESC
	`

	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MyEquipmentItem
	for rows.Next() {
		var item models.MyEquipmentItem
		if err := rows.Scan(
			&item.AssignmentID, &item.ResourceID, &item.ResourceName,
			&item.Quantity, &item.UnitType, &item.IssuedAt, &item.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// IssueToUser видає ресурс зі складу підрозділу командира конкретному солдату
func (r *ResourceRepository) IssueToUser(ctx context.Context, db *pgxpool.Pool, commanderUnitID int64, resourceID string, targetUserID string, quantity int, notes string) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("помилка старту транзакції: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Шукаємо склад, який належить підрозділу командира
	var warehouseID string
	err = tx.QueryRow(ctx, `SELECT id FROM warehouses WHERE unit_id = $1`, commanderUnitID).Scan(&warehouseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("у вашого підрозділу немає власного складу")
		}
		return err
	}

	// 2. Блокуємо ресурс для оновлення (FOR UPDATE) і перевіряємо чи він є на ЦЬОМУ складі
	var currentQty int
	err = tx.QueryRow(ctx, `
		SELECT quantity 
		FROM resources 
		WHERE id = $1 AND warehouse_id = $2 
		FOR UPDATE
	`, resourceID, warehouseID).Scan(&currentQty)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("цього ресурсу немає на складі вашого підрозділу")
		}
		return err
	}

	// 3. Перевіряємо залишки
	if currentQty < quantity {
		return fmt.Errorf("недостатньо на складі: є %d, запитується %d", currentQty, quantity)
	}

	// 4. Списуємо зі складу
	_, err = tx.Exec(ctx, `UPDATE resources SET quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, quantity, resourceID)
	if err != nil {
		return err
	}

	// 5. Записуємо факт видачі
	_, err = tx.Exec(ctx, `
		INSERT INTO resource_assignments (resource_id, user_id, quantity, status, notes)
		VALUES ($1, $2, $3, 'ACTIVE', $4)
	`, resourceID, targetUserID, quantity, notes)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// CreateShipment формує рейс (віднімає зі складу відправника і блокує авто)
// CreateShipment створює рейс, списує майно і блокує машину в одній транзакції
func (r *ResourceRepository) CreateShipment(ctx context.Context, db *pgxpool.Pool, req models.CreateShipmentRequest) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("помилка старту транзакції: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Перевіряємо машину та водія...
	var vStatus string
	var driverID *string
	var currentWarehouseID *string
	err = tx.QueryRow(ctx, "SELECT status, driver_id, current_warehouse_id FROM vehicles WHERE id = $1 FOR UPDATE", req.VehicleID).Scan(&vStatus, &driverID, &currentWarehouseID)
	if err != nil {
		return errors.New("транспорт не знайдено")
	}
	if vStatus != "ACTIVE" {
		return errors.New("цей транспорт недоступний або вже перебуває в рейсі")
	}

	// Перевіряємо, що машина знаходиться на складі відправника або отримувача
	// Якщо current_warehouse_id не встановлений (NULL) - дозволяємо (legacy машини)
	if currentWarehouseID != nil && *currentWarehouseID != "" {
		if *currentWarehouseID != req.FromWarehouseID && *currentWarehouseID != req.ToWarehouseID {
			return fmt.Errorf("машина зараз знаходиться на іншому складі. Для створення рейсу вона повинна бути на складі відправника або отримувача")
		}
	}

	// Перевіряємо чи водій не зайнятий іншим активним рейсом
	if driverID != nil && *driverID != "" {
		var activeShipmentCount int
		err = tx.QueryRow(ctx, `
			SELECT COUNT(*) 
			FROM shipments s 
			JOIN vehicles v ON s.vehicle_id = v.id 
			WHERE v.driver_id = $1 AND s.status IN ('PENDING', 'IN_TRANSIT')
		`, *driverID).Scan(&activeShipmentCount)

		if err == nil && activeShipmentCount > 0 {
			return errors.New("цей водій вже має активний рейс. Він зможе взяти новий після завершення поточного")
		}
	}

	// 2. Відправляємо машину...
	_, err = tx.Exec(ctx, "UPDATE vehicles SET status = 'ON_MISSION', updated_at = CURRENT_TIMESTAMP WHERE id = $1", req.VehicleID)
	if err != nil {
		return err
	}

	// 3. Отримуємо tenant_id та unit_id зі складу відправника та отримувача
	var tenantID *string
	var fromUnitID, toUnitID int64
	var fromWarehouseName, toWarehouseName string

	err = tx.QueryRow(ctx, "SELECT tenant_id, unit_id, name FROM warehouses WHERE id = $1", req.FromWarehouseID).Scan(&tenantID, &fromUnitID, &fromWarehouseName)
	if err != nil {
		return fmt.Errorf("не вдалося знайти склад відправника: %w", err)
	}

	err = tx.QueryRow(ctx, "SELECT unit_id, name FROM warehouses WHERE id = $1", req.ToWarehouseID).Scan(&toUnitID, &toWarehouseName)
	if err != nil {
		return fmt.Errorf("не вдалося знайти склад отримувача: %w", err)
	}

	// 🔄 Визначаємо напрямок руху по ієрархії (для логування та аналітики)
	// Система дозволяє рух в обидва боки: "вгору" (консолідація/повернення) та "вниз" (розподіл)
	var direction string
	if fromUnitID == toUnitID {
		direction = "LATERAL" // Переміщення між складами одного підрозділу
	} else {
		// Перевіряємо, чи є toUnit нащадком fromUnit
		var isDownward bool
		err = tx.QueryRow(ctx, `
			WITH RECURSIVE hierarchy AS (
				SELECT id FROM units WHERE id = $1
				UNION ALL
				SELECT u.id FROM units u
				JOIN hierarchy h ON u.parent_id = h.id
			)
			SELECT EXISTS(SELECT 1 FROM hierarchy WHERE id = $2)
		`, fromUnitID, toUnitID).Scan(&isDownward)

		if err == nil && isDownward {
			direction = "DOWNSTREAM" // Рух вниз по ієрархії (розподіл)
		} else {
			direction = "UPSTREAM" // Рух вгору по ієрархії (консолідація/повернення)
		}
	}

	// 4. Створюємо запис про рейс з напрямком руху...
	var shipmentID string
	err = tx.QueryRow(ctx, `
        INSERT INTO shipments (from_warehouse_id, to_warehouse_id, vehicle_id, priority, status, tenant_id, direction)
        VALUES ($1, $2, $3, $4, 'PENDING', $5, $6) RETURNING id
    `, req.FromWarehouseID, req.ToWarehouseID, req.VehicleID, req.Priority, tenantID, direction).Scan(&shipmentID)
	if err != nil {
		return err
	}

	// Логуємо створення рейсу з інформацією про напрямок
	fmt.Printf("✅ Створено рейс %s: %s (%d) → %s (%d) [%s]\n",
		shipmentID, fromWarehouseName, fromUnitID, toWarehouseName, toUnitID, direction)

	// 4. Списуємо майно зі складу відправника та додаємо в маніфест...
	for _, item := range req.Items {
		// Крок 4.1. Дізнаємося, ЩО САМЕ просив боєць (назву ресурсу)
		var resName string
		err = tx.QueryRow(ctx, "SELECT name FROM resources WHERE id = $1", item.ResourceID).Scan(&resName)
		if err != nil {
			return fmt.Errorf("ресурс із заявки не знайдено в базі (ID: %s)", item.ResourceID)
		}

		// Крок 4.2. Шукаємо такий самий ресурс НА СКЛАДІ ВІДПРАВНИКА
		var sourceResID string
		var currentQty int
		err = tx.QueryRow(ctx, "SELECT id, quantity FROM resources WHERE name = $1 AND warehouse_id = $2 FOR UPDATE", resName, req.FromWarehouseID).Scan(&sourceResID, &currentQty)
		if err != nil {
			// Якщо ресурсу з такою назвою взагалі немає на цьому складі
			return fmt.Errorf("На обраному складі відправника взагалі немає ресурсу: %s", resName)
		}

		// Крок 4.3. Перевіряємо, чи вистачає кількості
		if currentQty < item.Quantity {
			return fmt.Errorf("Недостатньо майна '%s' на складі відправника (треба: %d, є: %d)", resName, item.Quantity, currentQty)
		}

		// Крок 4.4. Списуємо майно САМЕ ЗІ СКЛАДУ ВІДПРАВНИКА
		_, err = tx.Exec(ctx, "UPDATE resources SET quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", item.Quantity, sourceResID)
		if err != nil {
			return err
		}

		// Крок 4.5. Записуємо в маніфест (вказуємо реальний ID ресурсу, який поїхав зі складу)
		_, err = tx.Exec(ctx, `
            INSERT INTO shipment_items (shipment_id, resource_id, quantity, request_id)
            VALUES ($1, $2, $3, $4)
        `, shipmentID, sourceResID, item.Quantity, item.RequestID)
		if err != nil {
			return err
		}
	}

	// 🔥 5. Переводимо заявки в DISPATCHED одразу при формуванні рейсу.
	// Для заявок з явним request_id — пряме оновлення.
	// Для рейсів без request_id (legacy або ручне створення) — шукаємо по resource_name + to_warehouse.
	var requestIDs []string
	for _, item := range req.Items {
		if item.RequestID != nil && *item.RequestID != "" {
			requestIDs = append(requestIDs, *item.RequestID)
		}
	}

	if len(requestIDs) > 0 {
		// Пряме оновлення по request_id
		_, err = tx.Exec(ctx, `
			UPDATE supply_requests
			SET status = 'DISPATCHED', updated_at = CURRENT_TIMESTAMP
			WHERE id = ANY($1) AND status IN ('PENDING', 'APPROVED')
		`, requestIDs)
		if err != nil {
			return fmt.Errorf("помилка оновлення статусу заявок: %w", err)
		}
	} else {
		// Legacy: рейс без request_id — знаходимо заявки по resource_name + to_warehouse
		_, err = tx.Exec(ctx, `
			UPDATE supply_requests sr
			SET status = 'DISPATCHED', updated_at = CURRENT_TIMESTAMP
			FROM shipment_items si
			JOIN resources r ON r.id = si.resource_id
			WHERE si.shipment_id = $1
			  AND sr.resource_name = r.name
			  AND sr.target_warehouse_id = $2
			  AND sr.status IN ('PENDING', 'APPROVED')
			  AND sr.tenant_id = $3
		`, shipmentID, req.ToWarehouseID, tenantID)
		if err != nil {
			return fmt.Errorf("помилка оновлення статусу заявок (legacy): %w", err)
		}
	}

	// Якщо всі кроки пройшли успішно — зберігаємо зміни в базу
	return tx.Commit(ctx)
}

// GetByWarehouse повертає всі доступні товари на конкретному складі
func (r *ResourceRepository) GetByWarehouse(ctx context.Context, db DBExecutor, warehouseID string) ([]models.InventoryItem, error) {
	args := []any{warehouseID}
	tFilter := tenantFilter(ctx, "r", "AND", &args)
	query := `
		SELECT 
			r.id, 
			CAST(r.warehouse_id AS TEXT),
			r.name, 
			COALESCE(c.name, 'Без категорії') as category, 
			r.quantity, 
			r.weight_kg
		FROM resources r
		LEFT JOIN resource_categories c ON r.category_id = c.id
		WHERE r.warehouse_id = $1 AND r.quantity > 0` + tFilter + `
		ORDER BY r.name ASC
	`

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання залишків складу: %w", err)
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.ID, &item.WarehouseID, &item.Name, &item.Category, &item.Available, &item.WeightKg); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

type ShipmentItem struct {
	ResourceName string  `json:"resource_name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
}

type ShipmentRecord struct {
	ID                string         `json:"id"`
	VehicleID         string         `json:"vehicle_id,omitempty"`
	VehiclePlate      string         `json:"vehicle_plate,omitempty"`
	FromWarehouseID   string         `json:"from_warehouse_id,omitempty"`
	FromWarehouse     string         `json:"from_warehouse,omitempty"`
	FromWarehouseName string         `json:"from_warehouse_name,omitempty"`
	ToWarehouseID     string         `json:"to_warehouse_id,omitempty"`
	ToWarehouse       string         `json:"to_warehouse,omitempty"`
	ToWarehouseName   string         `json:"to_warehouse_name,omitempty"`
	Vehicle           string         `json:"vehicle"`
	Priority          string         `json:"priority"`
	Status            string         `json:"status"`
	Direction         string         `json:"direction,omitempty"` // DOWNSTREAM, UPSTREAM, LATERAL
	CreatedAt         time.Time      `json:"created_at"`
	StartedAt         *time.Time     `json:"started_at,omitempty"`
	DeliveredAt       *time.Time     `json:"delivered_at,omitempty"`
	Items             []ShipmentItem `json:"items,omitempty"`
}

func (r *ResourceRepository) ListShipments(ctx context.Context, db DBExecutor) ([]ShipmentRecord, error) {
	var args []any
	tFilter := tenantFilter(ctx, "s", "WHERE", &args)
	query := fmt.Sprintf(`
		SELECT 
			s.id, w1.name as from_warehouse, w2.name as to_warehouse, 
			v.brand || ' (' || v.plate_number || ')' as vehicle, 
			s.priority, s.status, COALESCE(s.direction, '') as direction, s.created_at
		FROM shipments s
		JOIN warehouses w1 ON s.from_warehouse_id = w1.id
		JOIN warehouses w2 ON s.to_warehouse_id = w2.id
		JOIN vehicles v ON s.vehicle_id = v.id
		%s
		ORDER BY s.created_at DESC
	`, tFilter)
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ShipmentRecord
	for rows.Next() {
		var s ShipmentRecord
		if err := rows.Scan(&s.ID, &s.FromWarehouse, &s.ToWarehouse, &s.Vehicle, &s.Priority, &s.Status, &s.Direction, &s.CreatedAt); err == nil {
			list = append(list, s)
		}
	}
	if list == nil {
		list = []ShipmentRecord{}
	}
	return list, nil
}

func (r *ResourceRepository) ListMyShipments(ctx context.Context, db DBExecutor, userID string) ([]ShipmentRecord, error) {
	query := `
		SELECT 
			s.id, 
			s.vehicle_id,
			v.plate_number as vehicle_plate,
			s.from_warehouse_id,
			w1.name as from_warehouse_name, 
			s.to_warehouse_id,
			w2.name as to_warehouse_name, 
			s.status,
			COALESCE(s.direction, '') as direction, 
			s.created_at,
			s.started_at,
			s.delivered_at
		FROM shipments s
		JOIN warehouses w1 ON s.from_warehouse_id = w1.id
		JOIN warehouses w2 ON s.to_warehouse_id = w2.id
		JOIN vehicles v ON s.vehicle_id = v.id
		WHERE v.driver_id = $1
		ORDER BY 
			CASE s.status 
				WHEN 'PENDING' THEN 1
				WHEN 'IN_TRANSIT' THEN 2
				WHEN 'DELIVERED' THEN 3
			END,
			s.created_at DESC
	`
	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ShipmentRecord
	for rows.Next() {
		var s ShipmentRecord
		if err := rows.Scan(
			&s.ID,
			&s.VehicleID,
			&s.VehiclePlate,
			&s.FromWarehouseID,
			&s.FromWarehouseName,
			&s.ToWarehouseID,
			&s.ToWarehouseName,
			&s.Status,
			&s.Direction,
			&s.CreatedAt,
			&s.StartedAt,
			&s.DeliveredAt,
		); err == nil {
			// Завантажуємо items для рейсу
			itemsQuery := `
				SELECT 
					r.name as resource_name,
					si.quantity,
					r.unit
				FROM shipment_items si
				JOIN resources r ON si.resource_id = r.id
				WHERE si.shipment_id = $1
			`
			itemRows, itemErr := db.Query(ctx, itemsQuery, s.ID)
			if itemErr == nil {
				var items []ShipmentItem
				for itemRows.Next() {
					var item ShipmentItem
					if scanErr := itemRows.Scan(&item.ResourceName, &item.Quantity, &item.Unit); scanErr == nil {
						items = append(items, item)
					}
				}
				itemRows.Close()
				s.Items = items
			}

			list = append(list, s)
		}
	}
	if list == nil {
		list = []ShipmentRecord{}
	}
	return list, nil
}

func (r *ResourceRepository) ReceiveShipment(ctx context.Context, db *pgxpool.Pool, shipmentID string) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var targetWarehouseID, vehicleID, status string
	err = tx.QueryRow(ctx, `SELECT to_warehouse_id, vehicle_id, status FROM shipments WHERE id = $1 FOR UPDATE`, shipmentID).Scan(&targetWarehouseID, &vehicleID, &status)
	if err != nil {
		return errors.New("рейс не знайдено")
	}
	if status == "DELIVERED" {
		return errors.New("цей рейс вже було прийнято")
	}

	var targetUnitID int64
	tx.QueryRow(ctx, `SELECT unit_id FROM warehouses WHERE id = $1`, targetWarehouseID).Scan(&targetUnitID)

	// 1. Звільняємо машину, оновлюємо її локацію та закриваємо рейс
	_, err = tx.Exec(ctx, `UPDATE vehicles SET status = 'ACTIVE', current_warehouse_id = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, targetWarehouseID, vehicleID)
	if err != nil {
		return fmt.Errorf("помилка звільнення авто: %w", err)
	}

	_, err = tx.Exec(ctx, `UPDATE shipments SET status = 'DELIVERED', delivered_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, shipmentID)
	if err != nil {
		return fmt.Errorf("помилка закриття рейсу: %w", err)
	}

	// 2. АВТОМАТИЧНЕ ЗАКРИТТЯ ЗАЯВОК
	_, err = tx.Exec(ctx, `
        UPDATE supply_requests 
        SET status = 'COMPLETED', updated_at = CURRENT_TIMESTAMP 
        WHERE id IN (
            SELECT request_id FROM shipment_items WHERE shipment_id = $1 AND request_id IS NOT NULL
        )
    `, shipmentID)
	if err != nil {
		return fmt.Errorf("помилка оновлення заявок: %w", err)
	}

	// 3. Розподіляємо товари на новому складі (з COALESCE для захисту від NULL)
	rows, err := tx.Query(ctx, `
        SELECT si.quantity, r.category_id, r.name, COALESCE(r.description, ''), COALESCE(r.unit_type, 'PCS'), COALESCE(r.weight_kg, 1), COALESCE(r.min_quantity, 0)
        FROM shipment_items si JOIN resources r ON si.resource_id = r.id WHERE si.shipment_id = $1
    `, shipmentID)
	if err != nil {
		return fmt.Errorf("помилка читання маніфесту: %w", err)
	}

	type itemData struct {
		Qty    int
		CatID  string
		Name   string
		Desc   string
		UType  string
		Weight float64
		MinQty int
	}
	var items []itemData
	for rows.Next() {
		var i itemData
		if err := rows.Scan(&i.Qty, &i.CatID, &i.Name, &i.Desc, &i.UType, &i.Weight, &i.MinQty); err != nil {
			rows.Close()
			return fmt.Errorf("помилка сканування маніфесту: %w", err)
		}
		items = append(items, i)
	}
	rows.Close()

	// 4. Оновлюємо або створюємо ресурси на складі-одержувачі
	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for shipment receiving")
	}

	for _, item := range items {
		var existingResID string
		err := tx.QueryRow(ctx, `SELECT id FROM resources WHERE warehouse_id = $1 AND name = $2 AND condition != 'WRITTEN_OFF'`, targetWarehouseID, item.Name).Scan(&existingResID)

		if err == nil {
			// Якщо такий товар вже є — просто додаємо кількість
			_, updateErr := tx.Exec(ctx, `UPDATE resources SET quantity = quantity + $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, item.Qty, existingResID)
			if updateErr != nil {
				return fmt.Errorf("помилка оновлення залишку: %w", updateErr)
			}
		} else {
			// Якщо товару немає — створюємо новий запис (ДОДАНО condition = 'NEW' та tenant_id)
			_, insertErr := tx.Exec(ctx, `
                INSERT INTO resources (category_id, unit_id, warehouse_id, name, description, quantity, unit_type, weight_kg, min_quantity, condition, tenant_id) 
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'NEW', $10)
            `, item.CatID, targetUnitID, targetWarehouseID, item.Name, item.Desc, item.Qty, item.UType, item.Weight, item.MinQty, tid)
			if insertErr != nil {
				return fmt.Errorf("помилка створення нового ресурсу на складі: %w", insertErr)
			}
		}
	}

	return tx.Commit(ctx)
}

// StartShipment - водій або логіст підтверджує початок рейсу (PENDING → IN_TRANSIT)
func (r *ResourceRepository) StartShipment(ctx context.Context, db *pgxpool.Pool, shipmentID string) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("помилка старту транзакції: %w", err)
	}
	defer tx.Rollback(ctx)

	// Перевіряємо поточний статус
	var currentStatus string
	err = tx.QueryRow(ctx, `SELECT status FROM shipments WHERE id = $1 FOR UPDATE`, shipmentID).Scan(&currentStatus)
	if err != nil {
		return fmt.Errorf("рейс не знайдено: %w", err)
	}

	if currentStatus != "PENDING" {
		return fmt.Errorf("рейс вже в статусі '%s', неможливо почати", currentStatus)
	}

	// Оновлюємо статус на IN_TRANSIT та фіксуємо час початку
	_, err = tx.Exec(ctx, `
		UPDATE shipments 
		SET status = 'IN_TRANSIT', started_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1
	`, shipmentID)
	if err != nil {
		return fmt.Errorf("помилка оновлення статусу рейсу: %w", err)
	}

	// Оновлюємо статус пов'язаних заявок на DISPATCHED (тепер реально в дорозі).
	// Спочатку — заявки з явним request_id у маніфесті,
	// потім — заявки без прямого зв'язку, але які відповідають ресурсам рейсу
	// (для shipment-ів створених до додавання request_id або через legacy-потік).
	_, err = tx.Exec(ctx, `
		UPDATE supply_requests
		SET status = 'DISPATCHED', updated_at = CURRENT_TIMESTAMP
		WHERE status = 'APPROVED'
		  AND (
		    -- Пряме посилання через shipment_items.request_id
		    id IN (
		      SELECT request_id FROM shipment_items
		      WHERE shipment_id = $1 AND request_id IS NOT NULL
		    )
		    OR
		    -- Непряме: збіг за resource_name та target_warehouse = to_warehouse рейсу
		    id IN (
		      SELECT sr.id
		      FROM supply_requests sr
		      JOIN shipments s ON s.id = $1
		      JOIN shipment_items si ON si.shipment_id = s.id
		      JOIN resources r ON r.id = si.resource_id
		      WHERE sr.target_warehouse_id = s.to_warehouse_id
		        AND sr.resource_name = r.name
		        AND sr.tenant_id = s.tenant_id
		    )
		  )
	`, shipmentID)
	if err != nil {
		return fmt.Errorf("помилка оновлення статусу заявок: %w", err)
	}

	return tx.Commit(ctx)
}

// --- СТРУКТУРИ ДЛЯ ТТН (PDF) ---

type ShipmentInfo struct {
	FromWarehouse string
	ToWarehouse   string
	Vehicle       string
	Status        string
	CreatedAt     time.Time
}

type ShipmentItemInfo struct {
	Name string
	Qty  int
	Unit string
}

// GetShipmentInfo отримує базову інформацію про рейс для друку ТТН
func (r *ResourceRepository) GetShipmentInfo(ctx context.Context, db DBExecutor, shipmentID string) (*ShipmentInfo, error) {
	var info ShipmentInfo
	query := `
		SELECT w1.name, w2.name, v.brand || ' (' || v.plate_number || ')', s.status, s.created_at
		FROM shipments s
		JOIN warehouses w1 ON s.from_warehouse_id = w1.id
		JOIN warehouses w2 ON s.to_warehouse_id = w2.id
		JOIN vehicles v ON s.vehicle_id = v.id
		WHERE s.id = $1
	`
	err := db.QueryRow(ctx, query, shipmentID).Scan(
		&info.FromWarehouse,
		&info.ToWarehouse,
		&info.Vehicle,
		&info.Status,
		&info.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// GetShipmentItems отримує список майна у рейсі для друку ТТН
func (r *ResourceRepository) GetShipmentItems(ctx context.Context, db DBExecutor, shipmentID string) ([]ShipmentItemInfo, error) {
	query := `
		SELECT r.name, si.quantity, COALESCE(r.unit_type, 'шт')
		FROM shipment_items si
		JOIN resources r ON si.resource_id = r.id
		WHERE si.shipment_id = $1
	`
	rows, err := db.Query(ctx, query, shipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ShipmentItemInfo
	for rows.Next() {
		var i ShipmentItemInfo
		if err := rows.Scan(&i.Name, &i.Qty, &i.Unit); err == nil {
			items = append(items, i)
		}
	}
	return items, rows.Err()
}

func (r *ResourceRepository) SubmitInventoryAudit(ctx context.Context, db DBExecutor, userID string, req models.SubmitAuditRequest) error {
	// Починаємо транзакцію
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Створюємо сесію переобліку (Акт)
	var checkID string
	err = tx.QueryRow(ctx, `
		INSERT INTO inventory_checks (warehouse_id, created_by, status, completed_at)
		VALUES ($1, $2, 'COMPLETED', CURRENT_TIMESTAMP)
		RETURNING id
	`, req.WarehouseID, userID).Scan(&checkID)
	if err != nil {
		return fmt.Errorf("failed to create inventory check: %w", err)
	}

	// 2. Проходимося по всіх розбіжностях (якщо вони є)
	for _, item := range req.Discrepancies {
		// Записуємо факт перевірки конкретної позиції
		// (поле difference вираховується в БД автоматично через GENERATED ALWAYS)
		_, err = tx.Exec(ctx, `
			INSERT INTO inventory_check_items (check_id, resource_id, book_quantity, actual_quantity, verified_at)
			VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		`, checkID, item.ResourceID, item.BookQuantity, item.ActualQuantity)
		if err != nil {
			return fmt.Errorf("failed to insert check item for resource %s: %w", item.ResourceID, err)
		}

		// Якщо є розбіжність - оновлюємо фактичний залишок у таблиці resources
		if item.Difference != 0 {
			_, err = tx.Exec(ctx, `
				UPDATE resources
				SET quantity = $1, updated_at = CURRENT_TIMESTAMP
				WHERE id = $2
			`, item.ActualQuantity, item.ResourceID)
			if err != nil {
				return fmt.Errorf("failed to update resource %s quantity: %w", item.ResourceID, err)
			}
		}
	}

	// Зберігаємо всі зміни
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *CategoryRepository) GetAll(ctx context.Context, db DBExecutor) ([]models.ResourceCategory, error) {
	args := []any{}
	tFilter := tenantFilter(ctx, "", "WHERE", &args)
	query := `SELECT id, name FROM resource_categories` + tFilter
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.ResourceCategory
	for rows.Next() {
		var c models.ResourceCategory
		// Якщо в БД ID категорії це string або int64, адаптуй під свою модель
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *ResourceRepository) GetByNameAndWarehouse(ctx context.Context, db DBExecutor, name string, warehouseID string) (*models.Resource, error) {
	args := []any{name, warehouseID}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `SELECT id, quantity FROM resources WHERE name = $1 AND warehouse_id = $2` + tFilter + ` LIMIT 1`
	var res models.Resource
	err := db.QueryRow(ctx, query, args...).Scan(&res.ID, &res.Quantity)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
