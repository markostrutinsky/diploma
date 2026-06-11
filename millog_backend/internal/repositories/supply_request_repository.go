package repositories

import (
	"context"
	"fmt"
	"time"

	"Omnilog_backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplyRequestRepository struct{}

func NewSupplyRequestRepository() *SupplyRequestRepository {
	return &SupplyRequestRepository{}
}

func (r *SupplyRequestRepository) ValidateCreateScope(ctx context.Context, db DBExecutor, req *models.SupplyRequest) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for supply requests")
	}

	if req.ResourceID != nil && *req.ResourceID != "" {
		var ok bool
		if err := db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM resources WHERE id = $1 AND tenant_id = $2
			)
		`, *req.ResourceID, tid).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("ресурс не знайдено у вашій організації")
		}
	}

	if req.ResourceCategoryID != nil && *req.ResourceCategoryID != "" {
		var ok bool
		if err := db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM resource_categories WHERE id = $1 AND tenant_id = $2
			)
		`, *req.ResourceCategoryID, tid).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("категорію не знайдено у вашій організації")
		}
	}

	if req.TargetWarehouseID != "" {
		var ok bool
		if err := db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM warehouses WHERE id = $1 AND tenant_id = $2
			)
		`, req.TargetWarehouseID, tid).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("склад призначення не знайдено у вашій організації")
		}
	}

	return nil
}

func (r *SupplyRequestRepository) Create(ctx context.Context, db DBExecutor, req *models.SupplyRequest) error {
	var targetWH interface{}
	if req.TargetWarehouseID == "" {
		targetWH = nil
	} else {
		targetWH = req.TargetWarehouseID
	}

	var resID interface{}
	if req.ResourceID == nil || *req.ResourceID == "" {
		resID = nil
	} else {
		resID = *req.ResourceID
	}

	var catID interface{}
	if req.ResourceCategoryID == nil || *req.ResourceCategoryID == "" {
		catID = nil
	} else {
		catID = *req.ResourceCategoryID
	}

	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for supply requests")
	}
	if err := r.ValidateCreateScope(ctx, db, req); err != nil {
		return err
	}
	query := `INSERT INTO supply_requests (created_by, resource_id, resource_name, resource_category_id, quantity, status, target_warehouse_id, tenant_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`
	return db.QueryRow(ctx, query, req.CreatedBy, resID, req.ResourceName, catID, req.Quantity, req.Status, targetWH, tid).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *SupplyRequestRepository) GetByID(ctx context.Context, db DBExecutor, id string) (*models.SupplyRequest, error) {
	args := []any{id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `SELECT id, created_by, resource_id, COALESCE(resource_name, ''), resource_category_id, quantity, status, COALESCE(target_warehouse_id::text, ''), approved_by, approved_at, comment, created_at, updated_at
	FROM supply_requests WHERE id = $1` + tFilter

	var req models.SupplyRequest
	err := db.QueryRow(ctx, query, args...).Scan(
		&req.ID, &req.CreatedBy, &req.ResourceID, &req.ResourceName, &req.ResourceCategoryID, &req.Quantity, &req.Status, &req.TargetWarehouseID,
		&req.ApprovedBy, &req.ApprovedAt, &req.Comment, &req.CreatedAt, &req.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *SupplyRequestRepository) List(ctx context.Context, db DBExecutor, userRole string, userUnitID *int64) ([]models.SupplyRequest, error) {
	var rows pgx.Rows
	var err error

	if userRole == "ADMIN" || userRole == "TENANT_ADMIN" || userRole == "SYSTEM_ADMIN" || userRole == "CONTRACTOR" {
		args := []any{}
		tFilter := tenantFilter(ctx, "", "WHERE", &args)
		query := `SELECT id, created_by, resource_id, COALESCE(resource_name, ''), resource_category_id, quantity, status, COALESCE(target_warehouse_id::text, ''), approved_by, approved_at, comment, created_at, updated_at
				  FROM supply_requests` + tFilter + ` ORDER BY created_at DESC`
		rows, err = db.Query(ctx, query, args...)
	} else {
		if userUnitID == nil {
			return []models.SupplyRequest{}, nil
		}
		args := []any{*userUnitID}
		tFilter := tenantFilter(ctx, "sr", "AND", &args)
		query := `
			WITH RECURSIVE unit_tree AS (
				SELECT id FROM units WHERE id = $1
				UNION
				SELECT u.id FROM units u
				INNER JOIN unit_tree ut ON u.parent_id = ut.id
			)
			SELECT sr.id, sr.created_by, sr.resource_id, COALESCE(sr.resource_name, ''), sr.resource_category_id, sr.quantity, sr.status, COALESCE(sr.target_warehouse_id::text, ''), sr.approved_by, sr.approved_at, sr.comment, sr.created_at, sr.updated_at
			FROM supply_requests sr
			JOIN users u ON sr.created_by = u.id
			WHERE u.unit_id IN (SELECT id FROM unit_tree)` + tFilter + `
			ORDER BY sr.created_at DESC
		`
		rows, err = db.Query(ctx, query, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.SupplyRequest
	for rows.Next() {
		var req models.SupplyRequest
		if err := rows.Scan(&req.ID, &req.CreatedBy, &req.ResourceID, &req.ResourceName, &req.ResourceCategoryID, &req.Quantity, &req.Status, &req.TargetWarehouseID,
			&req.ApprovedBy, &req.ApprovedAt, &req.Comment, &req.CreatedAt, &req.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, req)
	}
	return list, rows.Err()
}

func (r *SupplyRequestRepository) Approve(ctx context.Context, db DBExecutor, id, approvedBy string, approved bool, comment string) error {
	status := models.RequestRejected
	if approved {
		status = models.RequestApproved
	}

	var validApprovedBy interface{}
	if approvedBy == "" {
		validApprovedBy = nil
	} else {
		validApprovedBy = approvedBy
	}

	args := []any{status, validApprovedBy, comment, id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `UPDATE supply_requests SET status = $1, approved_by = $2, approved_at = CURRENT_TIMESTAMP, comment = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4` + tFilter

	_, err := db.Exec(ctx, query, args...)
	return err
}

// UpdateStatus змінює статус та додає коментар (наприклад, причину відмови)
func (r *SupplyRequestRepository) UpdateStatus(ctx context.Context, db DBExecutor, id string, status string, comment string) error {
	args := []any{status, comment, id}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `UPDATE supply_requests SET status = $1, comment = $2 WHERE id = $3` + tFilter
	_, err := db.Exec(ctx, query, args...)
	return err
}

// EscalateStatus змінює статус на ESCALATED ТІЛЬКИ якщо поточний статус — PENDING.
// Захищає від гонки стану: якщо заявку скасували між вибіркою і оновленням — вона НЕ буде ескальована.
func (r *SupplyRequestRepository) EscalateStatus(ctx context.Context, db DBExecutor, id string, comment string) (bool, error) {
	args := []any{"ESCALATED", comment, id, "PENDING"}
	tFilter := tenantFilter(ctx, "", "AND", &args)
	query := `UPDATE supply_requests SET status = $1, comment = $2 WHERE id = $3 AND status = $4` + tFilter
	result, err := db.Exec(ctx, query, args...)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}

// GetRequestsForDispatch витягує вибрані заявки, їхню вагу та цільовий склад
func (r *SupplyRequestRepository) GetRequestsForDispatch(ctx context.Context, db DBExecutor, reqIDs []string) ([]models.RequestItem, error) {
	args := []any{reqIDs}
	tFilter := tenantFilter(ctx, "sr", "AND", &args)
	query := `
		SELECT sr.id::text,
		       COALESCE(res.name, sr.resource_name, 'Невідомий ресурс') AS name,
		       (COALESCE(res.weight_kg, 1.0) * sr.quantity) AS weight_kg,
		       COALESCE(sr.target_warehouse_id::text, '') AS target_warehouse_id
		FROM supply_requests sr
		LEFT JOIN resources res ON sr.resource_id = res.id
		WHERE sr.id = ANY($1) AND sr.status IN ('PENDING', 'APPROVED')` + tFilter + `
	`
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.RequestItem
	for rows.Next() {
		var item models.RequestItem
		if err := rows.Scan(&item.ID, &item.Name, &item.WeightKg, &item.TargetWarehouseID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// GetAvailableVehicles повертає вільні авто для Smart Розподілу.
// Повертає тільки транспорт, що фізично знаходиться на складі відправника
// або складі отримувача. Ніякої ієрархії — лише ці два склади.
// Якщо fromWarehouseID порожній — шукає тільки на складі отримувача.
func (r *SupplyRequestRepository) GetAvailableVehicles(ctx context.Context, db DBExecutor, fromWarehouseID string, targetWarehouseID string) ([]models.VehicleBin, error) {
	var args []any
	var warehouseFilter string

	if fromWarehouseID != "" {
		args = []any{fromWarehouseID, targetWarehouseID}
		warehouseFilter = `v.current_warehouse_id IN ($1, $2)`
	} else {
		args = []any{targetWarehouseID}
		warehouseFilter = `v.current_warehouse_id = $1`
	}

	tFilter := tenantFilter(ctx, "v", "AND", &args)

	query := `
		SELECT v.id::text, v.brand || ' ' || COALESCE(v.plate_number, ''), v.capacity_kg,
		       v.fuel_norm, v.tank_capacity,
		       GREATEST(0, COALESCE((
		         SELECT SUM(CASE WHEN record_type = 'REFUEL' THEN liters ELSE -liters END)
		         FROM fuel_records WHERE vehicle_id = v.id
		       ), 0)) AS current_fuel_liters
		FROM vehicles v
		WHERE v.status = 'ACTIVE'
		  AND v.driver_id IS NOT NULL
		  AND v.type IN ('VAN', 'TRUCK', 'PICKUP')
		  AND ` + warehouseFilter + tFilter + `
		ORDER BY v.capacity_kg ASC
	`
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []models.VehicleBin
	for rows.Next() {
		var v models.VehicleBin
		if err := rows.Scan(&v.ID, &v.Name, &v.MaxWeight, &v.FuelNorm, &v.TankCapacity, &v.FuelLiters); err != nil {
			return nil, err
		}
		v.UsedWeight = 0
		v.Items = make([]models.RequestItem, 0)
		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

func (r *SupplyRequestRepository) GetOverdueRequests(ctx context.Context, db *pgxpool.Pool, status string, hoursLimit int) ([]models.SupplyRequest, error) {
	// Визначаємо часову межу: поточний час мінус ліміт годин
	threshold := time.Now().Add(-time.Duration(hoursLimit) * time.Hour)

	query := `
		SELECT id, created_by, COALESCE(target_warehouse_id::text, ''), status, created_at 
		FROM supply_requests 
		WHERE status = $1 AND created_at < $2
	`

	rows, err := db.Query(ctx, query, status, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overdue []models.SupplyRequest
	for rows.Next() {
		var req models.SupplyRequest
		if err := rows.Scan(&req.ID, &req.CreatedBy, &req.TargetWarehouseID, &req.Status, &req.CreatedAt); err == nil {
			overdue = append(overdue, req)
		}
	}
	return overdue, nil
}

func (r *SupplyRequestRepository) GetEscalatedOverdueRequests(ctx context.Context, db *pgxpool.Pool, status string, hoursLimit int) ([]models.SupplyRequest, error) {
	threshold := time.Now().Add(-time.Duration(hoursLimit) * time.Hour)

	// 🔥 МАГІЯ SQL: Динамічно шукаємо email локального керівника
	// База перевірить, чи є в цьому підрозділі директор, керівник філії або відділу.
	query := `
		SELECT 
			sr.id, sr.created_by, COALESCE(sr.target_warehouse_id::text, ''), sr.status, sr.created_at,
			COALESCE(
				(SELECT u.email FROM users u 
				 JOIN warehouses w ON u.unit_id = w.unit_id 
				 WHERE w.id = sr.target_warehouse_id 
				   AND u.role IN ('BRANCH_MANAGER', 'REGION_DIRECTOR', 'DEPT_MANAGER') 
				 LIMIT 1), 
				'admin@omnilog.app' -- Резервний імейл, якщо підрозділ тимчасово без керівника
			) AS manager_email
		FROM supply_requests sr
		WHERE sr.status = $1 AND sr.created_at < $2
	`

	rows, err := db.Query(ctx, query, status, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overdue []models.SupplyRequest
	for rows.Next() {
		var req models.SupplyRequest
		// Скануємо нове поле &req.ManagerEmail (яке ти вже додав у models/requests.go)
		if err := rows.Scan(
			&req.ID, &req.CreatedBy, &req.TargetWarehouseID,
			&req.Status, &req.CreatedAt, &req.ManagerEmail,
		); err == nil {
			overdue = append(overdue, req)
		}
	}
	return overdue, nil
}
