package repositories

import (
	"context"

	"millog_backend/internal/models"
)

type SupplyRequestRepository struct{}

func NewSupplyRequestRepository() *SupplyRequestRepository {
	return &SupplyRequestRepository{}
}

func (r *SupplyRequestRepository) Create(ctx context.Context, db DBExecutor, req *models.SupplyRequest) error {
	query := `INSERT INTO supply_requests (created_by, resource_id, quantity, status)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return db.QueryRow(ctx, query, req.CreatedBy, req.ResourceID, req.Quantity, req.Status).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *SupplyRequestRepository) GetByID(ctx context.Context, db DBExecutor, id string) (*models.SupplyRequest, error) {
	query := `SELECT id, created_by, resource_id, quantity, status, approved_by, approved_at, comment, created_at, updated_at
	FROM supply_requests WHERE id = $1`
	var req models.SupplyRequest
	err := db.QueryRow(ctx, query, id).Scan(
		&req.ID, &req.CreatedBy, &req.ResourceID, &req.Quantity, &req.Status,
		&req.ApprovedBy, &req.ApprovedAt, &req.Comment, &req.CreatedAt, &req.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *SupplyRequestRepository) List(ctx context.Context, db DBExecutor) ([]models.SupplyRequest, error) {
	rows, err := db.Query(ctx, `SELECT id, created_by, resource_id, quantity, status, approved_by, approved_at, comment, created_at, updated_at
	FROM supply_requests ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.SupplyRequest
	for rows.Next() {
		var req models.SupplyRequest
		if err := rows.Scan(&req.ID, &req.CreatedBy, &req.ResourceID, &req.Quantity, &req.Status,
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
	query := `UPDATE supply_requests SET status = $1, approved_by = $2, approved_at = CURRENT_TIMESTAMP, comment = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	_, err := db.Exec(ctx, query, status, approvedBy, comment, id)
	return err
}
