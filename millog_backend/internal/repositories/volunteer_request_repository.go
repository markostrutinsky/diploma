package repositories

import (
	"context"
	"fmt"
	"strings"

	"Omnilog_backend/internal/models"
)

type ContractorRequestRepository struct{}

func NewContractorRequestRepository() *ContractorRequestRepository {
	return &ContractorRequestRepository{}
}

// ValidateCreateScope гарантує, що підрозділ і склад, передані з body, належать
// поточній організації. FK перевіряє лише існування UUID, але не tenant ownership.
func (r *ContractorRequestRepository) ValidateCreateScope(ctx context.Context, db DBExecutor, unitID *int64, warehouseID *string) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for contractor requests")
	}

	if unitID != nil {
		var ok bool
		if err := db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM units WHERE id = $1 AND tenant_id = $2
			)
		`, *unitID, tid).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("підрозділ не знайдено у вашій організації")
		}
	}

	if warehouseID != nil && *warehouseID != "" {
		var ok bool
		if err := db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM warehouses WHERE id = $1 AND tenant_id = $2
			)
		`, *warehouseID, tid).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("склад не знайдено у вашій організації")
		}
	}

	return nil
}

func (r *ContractorRequestRepository) ValidateCategoryScope(ctx context.Context, db DBExecutor, categoryID string) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for creating resources")
	}

	var ok bool
	if err := db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM resource_categories WHERE id = $1 AND tenant_id = $2
		)
	`, categoryID, tid).Scan(&ok); err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("категорію не знайдено у вашій організації")
	}
	return nil
}

func (r *ContractorRequestRepository) ResolveCategoryID(ctx context.Context, db DBExecutor, categoryID, categoryName string) (string, error) {
	categoryID = strings.TrimSpace(categoryID)
	categoryName = strings.TrimSpace(categoryName)

	if categoryID != "" {
		if err := r.ValidateCategoryScope(ctx, db, categoryID); err != nil {
			return "", err
		}
		return categoryID, nil
	}

	if categoryName == "" {
		return "", fmt.Errorf("оберіть категорію або введіть назву нової категорії")
	}

	tid := TenantFromCtx(ctx)
	if tid == "" {
		return "", fmt.Errorf("tenant_id is required for creating categories")
	}

	var id string
	if err := db.QueryRow(ctx, `
		INSERT INTO resource_categories (name, description, tenant_id)
		VALUES ($1, '', $2)
		ON CONFLICT (tenant_id, name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, categoryName, tid).Scan(&id); err != nil {
		return "", fmt.Errorf("не вдалося створити категорію: %w", err)
	}

	return id, nil
}

type ContractorAcceptanceTarget struct {
	RequestUnitID         *int64
	TargetWarehouseID     *string
	TargetWarehouseUnitID *int64
}

func (r *ContractorRequestRepository) GetAcceptanceTarget(ctx context.Context, db DBExecutor, requestID string) (*ContractorAcceptanceTarget, error) {
	args := []any{requestID}
	tFilter := tenantFilter(ctx, "cr", "AND", &args)
	query := `
		SELECT cr.unit_id, cr.target_warehouse_id, w.unit_id
		FROM contractor_requests cr
		LEFT JOIN warehouses w ON w.id = cr.target_warehouse_id
		WHERE cr.id = $1` + tFilter + `
	`

	var target ContractorAcceptanceTarget
	if err := db.QueryRow(ctx, query, args...).Scan(
		&target.RequestUnitID,
		&target.TargetWarehouseID,
		&target.TargetWarehouseUnitID,
	); err != nil {
		return nil, err
	}

	return &target, nil
}

func (r *ContractorRequestRepository) Create(ctx context.Context, db DBExecutor, vr *models.ContractorRequest) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		return fmt.Errorf("tenant_id is required for contractor requests")
	}
	query := `
		INSERT INTO contractor_requests (created_by, unit_id, target_warehouse_id, title, description, status, tenant_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, created_at`
	return db.QueryRow(ctx, query, vr.CreatedBy, vr.UnitID, vr.TargetWarehouseID, vr.Title, vr.Description, vr.Status, tid).Scan(&vr.ID, &vr.CreatedAt)
}

func (r *ContractorRequestRepository) List(ctx context.Context, db DBExecutor, status models.ContractorRequestStatus, isContractor bool, contractorID string) ([]models.ContractorRequest, error) {
	query := `
		SELECT 
			vr.id, vr.created_by, vr.unit_id, u.name as unit_name, 
			vr.target_warehouse_id, w.name as warehouse_name,
			vr.title, vr.description, 
			vr.status, vr.taken_by, vr.taken_at, vr.completed_at, vr.created_at,
			vr.tenant_id, t.name as tenant_name
		FROM contractor_requests vr
		LEFT JOIN units u ON vr.unit_id = u.id
		LEFT JOIN warehouses w ON vr.target_warehouse_id = w.id
		LEFT JOIN tenants t ON vr.tenant_id = t.id
		WHERE 1=1
	`

	var args []interface{}

	if status != "" {
		args = append(args, status)
		query += fmt.Sprintf(" AND vr.status = $%d", len(args))
	}

	if isContractor {
		// Marketplace для підрядника: бачить усі ВІДКРИТІ завдання будь-якої організації
		// плюс ті, які він особисто взяв (його активна/історична робота). Tenant-фільтр
		// свідомо НЕ застосовуємо — підрядник глобальний і працює крос-tenant.
		query += " AND t.is_active = TRUE"
		args = append(args, contractorID)
		query += fmt.Sprintf(" AND (vr.status = 'OPEN' OR vr.taken_by = $%d)", len(args))
	} else {
		// Бізнес-користувач бачить лише заявки своєї організації.
		query += tenantFilter(ctx, "vr", "AND", &args)
	}

	query += " ORDER BY vr.created_at DESC"

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ContractorRequest
	for rows.Next() {
		var vr models.ContractorRequest
		if err := rows.Scan(
			&vr.ID,
			&vr.CreatedBy,
			&vr.UnitID,
			&vr.UnitName,
			&vr.TargetWarehouseID,
			&vr.WarehouseName,
			&vr.Title,
			&vr.Description,
			&vr.Status,
			&vr.TakenBy,
			&vr.TakenAt,
			&vr.CompletedAt,
			&vr.CreatedAt,
			&vr.TenantID,
			&vr.TenantName,
		); err != nil {
			return nil, err
		}
		list = append(list, vr)
	}
	return list, rows.Err()
}

// GetTenantID повертає tenant_id організації-замовника відкритої заявки.
// Потрібно, щоб перевірити, чи має підрядник схвалення саме від цієї організації.
func (r *ContractorRequestRepository) GetTenantID(ctx context.Context, db DBExecutor, requestID string) (string, error) {
	var tid *string
	err := db.QueryRow(ctx, `
		SELECT vr.tenant_id
		FROM contractor_requests vr
		JOIN tenants t ON t.id = vr.tenant_id
		WHERE vr.id = $1 AND vr.status = 'OPEN' AND t.is_active = TRUE
	`, requestID).Scan(&tid)
	if err != nil {
		return "", err
	}
	if tid == nil {
		return "", nil
	}
	return *tid, nil
}

func (r *ContractorRequestRepository) UpdateStatus(ctx context.Context, db DBExecutor, requestID string, userID string, newStatus models.ContractorRequestStatus) error {

	query := `UPDATE contractor_requests SET status = $1`
	args := []interface{}{newStatus}
	paramID := 2
	var where string

	switch newStatus {
	case models.ContractorTaken:
		query += fmt.Sprintf(", taken_by = $%d, taken_at = CURRENT_TIMESTAMP", paramID)
		args = append(args, userID)
		paramID++
		where = fmt.Sprintf(" WHERE id = $%d AND status = 'OPEN'", paramID)
		args = append(args, requestID)
		where += " AND EXISTS (SELECT 1 FROM tenants t WHERE t.id = contractor_requests.tenant_id AND t.is_active = TRUE)"
	case models.ContractorAccepted, models.ContractorRejected, models.ContractorCanceled:
		if TenantFromCtx(ctx) == "" {
			return fmt.Errorf("tenant_id is required for this operation")
		}
		query += ", completed_at = CURRENT_TIMESTAMP"
		where = fmt.Sprintf(" WHERE id = $%d", paramID)
		args = append(args, requestID)
		if newStatus == models.ContractorCanceled {
			where += " AND status IN ('OPEN', 'TAKEN')"
		} else {
			where += " AND status = 'DELIVERED'"
		}
		where += tenantFilter(ctx, "", "AND", &args)
	case models.ContractorDelivered:
		where = fmt.Sprintf(" WHERE id = $%d AND status = 'TAKEN'", paramID)
		args = append(args, requestID)
		paramID++
		where += fmt.Sprintf(" AND taken_by = $%d", paramID)
		args = append(args, userID)
		where += " AND EXISTS (SELECT 1 FROM tenants t WHERE t.id = contractor_requests.tenant_id AND t.is_active = TRUE)"
	default:
		return fmt.Errorf("некоректний статус заявки")
	}

	query += where

	cmdTag, err := db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("заявку не знайдено або статус не змінено")
	}

	return nil
}
