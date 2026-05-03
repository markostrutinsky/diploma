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
		INSERT INTO contractor_requests (created_by, unit_id, title, description, status)
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at`
		return db.QueryRow(ctx, query, vr.CreatedBy, vr.UnitID, vr.Title, vr.Description, vr.Status).Scan(&vr.ID, &vr.CreatedAt)
	}
	query := `
		INSERT INTO contractor_requests (created_by, unit_id, title, description, status, tenant_id)
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at`
	return db.QueryRow(ctx, query, vr.CreatedBy, vr.UnitID, vr.Title, vr.Description, vr.Status, tid).Scan(&vr.ID, &vr.CreatedAt)
}

func (r *ContractorRequestRepository) List(ctx context.Context, db DBExecutor, status models.ContractorRequestStatus) ([]models.ContractorRequest, error) {
	query := `
		SELECT 
			vr.id, vr.created_by, vr.unit_id, u.name as unit_name, vr.title, vr.description, 
			vr.status, vr.taken_by, vr.taken_at, vr.completed_at, vr.created_at
		FROM contractor_requests vr
		LEFT JOIN units u ON vr.unit_id = u.id
		WHERE 1=1
	`

	var args []interface{}
	paramID := 1

	if status != "" {
		query += fmt.Sprintf(" AND vr.status = $%d", paramID)
		args = append(args, status)
		paramID++
	}

	// Volunteer-маркетплейс: зазвичай показуємо крос-tenant заявки волонтерам (CONTRACTOR).
	// Тому tenant фільтр — опціональний. Якщо у контексті є tenant — бізнес-користувач бачить лише свої.
	// CONTRACTOR'и не мають tenant_id → бачать всі.
	tFilter := tenantFilter(ctx, "vr", "AND", &args)
	query += tFilter
	_ = paramID

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
			&vr.Title,
			&vr.Description,
			&vr.Status,
			&vr.TakenBy,
			&vr.TakenAt,
			&vr.CompletedAt,
			&vr.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, vr)
	}
	return list, rows.Err()
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
