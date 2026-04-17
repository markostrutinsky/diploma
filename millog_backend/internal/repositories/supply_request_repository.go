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
	// ДОДАНО: target_warehouse_id у запит
	query := `INSERT INTO supply_requests (created_by, resource_id, quantity, status, target_warehouse_id)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	return db.QueryRow(ctx, query, req.CreatedBy, req.ResourceID, req.Quantity, req.Status, req.TargetWarehouseID).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *SupplyRequestRepository) GetByID(ctx context.Context, db DBExecutor, id string) (*models.SupplyRequest, error) {
	// ДОДАНО: target_warehouse_id у SELECT
	query := `SELECT id, created_by, resource_id, quantity, status, target_warehouse_id, approved_by, approved_at, comment, created_at, updated_at
	FROM supply_requests WHERE id = $1`

	var req models.SupplyRequest
	err := db.QueryRow(ctx, query, id).Scan(
		&req.ID, &req.CreatedBy, &req.ResourceID, &req.Quantity, &req.Status, &req.TargetWarehouseID,
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

	if userRole == "ADMIN" || userRole == "CONTRACTOR" {
		// ДОДАНО: target_warehouse_id у SELECT
		query := `SELECT id, created_by, resource_id, quantity, status, target_warehouse_id, approved_by, approved_at, comment, created_at, updated_at
				  FROM supply_requests ORDER BY created_at DESC`
		rows, err = db.Query(ctx, query)
	} else {
		if userUnitID == nil {
			return []models.SupplyRequest{}, nil
		}

		// ДОДАНО: sr.target_warehouse_id у SELECT
		query := `
			WITH RECURSIVE unit_tree AS (
				SELECT id FROM units WHERE id = $1
				UNION
				SELECT u.id FROM units u
				INNER JOIN unit_tree ut ON u.parent_id = ut.id
			)
			SELECT sr.id, sr.created_by, sr.resource_id, sr.quantity, sr.status, sr.target_warehouse_id, sr.approved_by, sr.approved_at, sr.comment, sr.created_at, sr.updated_at
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
		// ДОДАНО: &req.TargetWarehouseID у Scan
		if err := rows.Scan(&req.ID, &req.CreatedBy, &req.ResourceID, &req.Quantity, &req.Status, &req.TargetWarehouseID,
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

	query := `UPDATE supply_requests SET status = $1, approved_by = $2, approved_at = CURRENT_TIMESTAMP, comment = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`

	_, err := db.Exec(ctx, query, status, validApprovedBy, comment, id)
	return err
}

// UpdateStatus змінює статус та додає коментар (наприклад, причину відмови)
func (r *SupplyRequestRepository) UpdateStatus(ctx context.Context, db DBExecutor, id string, status string, comment string) error {
	query := `UPDATE supply_requests SET status = $1, comment = $2 WHERE id = $3`
	_, err := db.Exec(ctx, query, status, comment, id)
	return err
}
