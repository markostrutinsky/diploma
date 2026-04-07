package repositories

import (
	"context"

	"millog_backend/internal/models"

	"github.com/jackc/pgx/v5"
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

func (r *SupplyRequestRepository) List(ctx context.Context, db DBExecutor, userRole string, userUnitID *int64) ([]models.SupplyRequest, error) {
	var rows pgx.Rows
	var err error

	if userRole == "ADMIN" || userRole == "VOLUNTEER" {
		query := `SELECT id, created_by, resource_id, quantity, status, approved_by, approved_at, comment, created_at, updated_at
				  FROM supply_requests ORDER BY created_at DESC`
		rows, err = db.Query(ctx, query)
	} else {
		if userUnitID == nil {
			return []models.SupplyRequest{}, nil
		}

		query := `
			WITH RECURSIVE unit_tree AS (
				-- Беремо підрозділ поточного користувача
				SELECT id FROM units WHERE id = $1
				UNION
				-- Шукаємо всі підрозділи, які йому підпорядковуються (через parent_id)
				SELECT u.id FROM units u
				INNER JOIN unit_tree ut ON u.parent_id = ut.id
			)
			SELECT sr.id, sr.created_by, sr.resource_id, sr.quantity, sr.status, sr.approved_by, sr.approved_at, sr.comment, sr.created_at, sr.updated_at
			FROM supply_requests sr
			JOIN users u ON sr.created_by = u.id
			WHERE u.unit_id IN (SELECT id FROM unit_tree)
			ORDER BY sr.created_at DESC
		`
		rows, err = db.Query(ctx, query, *userUnitID)
	}

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
