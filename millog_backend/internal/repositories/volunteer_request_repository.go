package repositories

import (
	"context"
	"fmt"

	"Omnilog_backend/internal/models"
)

type ContractorRequestRepository struct{}

func NewContractorRequestRepository() *ContractorRequestRepository {
	return &ContractorRequestRepository{}
}

func (r *ContractorRequestRepository) Create(ctx context.Context, db DBExecutor, vr *models.ContractorRequest) error {
	tid := TenantFromCtx(ctx)
	if tid == "" {
		query := `
		INSERT INTO contractor_requests (created_by, unit_id, target_warehouse_id, title, description, status)
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at`
		return db.QueryRow(ctx, query, vr.CreatedBy, vr.UnitID, vr.TargetWarehouseID, vr.Title, vr.Description, vr.Status).Scan(&vr.ID, &vr.CreatedAt)
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

// GetTenantID повертає tenant_id організації-замовника заявки.
// Потрібно, щоб перевірити, чи має підрядник схвалення саме від цієї організації.
func (r *ContractorRequestRepository) GetTenantID(ctx context.Context, db DBExecutor, requestID string) (string, error) {
	var tid *string
	err := db.QueryRow(ctx, `SELECT tenant_id FROM contractor_requests WHERE id = $1`, requestID).Scan(&tid)
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

	switch newStatus {
	case models.ContractorTaken:
		query += fmt.Sprintf(", taken_by = $%d, taken_at = CURRENT_TIMESTAMP", paramID)
		args = append(args, userID)
		paramID++
	case models.ContractorAccepted, models.ContractorRejected, models.ContractorCanceled:
		query += ", completed_at = CURRENT_TIMESTAMP"
	}

	query += fmt.Sprintf(" WHERE id = $%d", paramID)
	args = append(args, requestID)

	cmdTag, err := db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("заявку не знайдено або статус не змінено")
	}

	return nil
}
