package repositories

import (
	"context"
	"fmt"

	"millog_backend/internal/models"
)

type CONTRACTORRequestRepository struct{}

func NewCONTRACTORRequestRepository() *CONTRACTORRequestRepository {
	return &CONTRACTORRequestRepository{}
}

func (r *CONTRACTORRequestRepository) Create(ctx context.Context, db DBExecutor, vr *models.CONTRACTORRequest) error {
	query := `
		INSERT INTO CONTRACTOR_requests (created_by, unit_id, title, description, status)
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at
	`
	return db.QueryRow(ctx, query, vr.CreatedBy, vr.UnitID, vr.Title, vr.Description, vr.Status).Scan(&vr.ID, &vr.CreatedAt)
}

func (r *CONTRACTORRequestRepository) List(ctx context.Context, db DBExecutor, status models.CONTRACTORRequestStatus) ([]models.CONTRACTORRequest, error) {

	query := `
		SELECT 
			vr.id, 
			vr.created_by, 
			vr.unit_id, 
			u.name as unit_name, 
			vr.title, 
			vr.description, 
			vr.status, 
			vr.taken_by, 
			vr.taken_at, 
			vr.completed_at, 
			vr.created_at
		FROM CONTRACTOR_requests vr
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

	query += " ORDER BY vr.created_at DESC"

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.CONTRACTORRequest
	for rows.Next() {
		var vr models.CONTRACTORRequest
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

func (r *CONTRACTORRequestRepository) UpdateStatus(ctx context.Context, db DBExecutor, requestID string, userID string, newStatus models.CONTRACTORRequestStatus) error {

	query := `UPDATE CONTRACTOR_requests SET status = $1`
	args := []interface{}{newStatus}
	paramID := 2

	switch newStatus {
	case models.CONTRACTORTaken:
		query += fmt.Sprintf(", taken_by = $%d, taken_at = CURRENT_TIMESTAMP", paramID)
		args = append(args, userID)
		paramID++
	case models.CONTRACTORAccepted, models.CONTRACTORRejected, models.CONTRACTORCanceled:
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
