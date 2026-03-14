package repositories

import (
	"context"

	"millog_backend/internal/models"
)

type ResourceRepository struct{}

func NewResourceRepository() *ResourceRepository {
	return &ResourceRepository{}
}

func (r *ResourceRepository) Create(ctx context.Context, db DBExecutor, res *models.Resource) error {
	query := `INSERT INTO resources (category_id, unit_id, name, description, quantity, serial_number, location, condition, min_quantity)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`
	return db.QueryRow(ctx, query,
		res.CategoryID, res.UnitID, res.Name, res.Description, res.Quantity, res.SerialNumber,
		res.Location, res.Condition, res.MinQuantity,
	).Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt)
}

func (r *ResourceRepository) GetByID(ctx context.Context, db DBExecutor, id string) (*models.Resource, error) {
	query := `SELECT id, category_id, unit_id, name, description, quantity, serial_number, location, condition, min_quantity, created_at, updated_at
	FROM resources WHERE id = $1`
	var res models.Resource
	err := db.QueryRow(ctx, query, id).Scan(
		&res.ID, &res.CategoryID, &res.UnitID, &res.Name, &res.Description, &res.Quantity,
		&res.SerialNumber, &res.Location, &res.Condition, &res.MinQuantity,
		&res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *ResourceRepository) List(ctx context.Context, db DBExecutor, unitID *int64) ([]models.Resource, error) {
	query := `SELECT id, category_id, unit_id, name, description, quantity, serial_number, location, condition, min_quantity, created_at, updated_at
	FROM resources`
	args := []interface{}{}
	if unitID != nil {
		query += ` WHERE unit_id = $1`
		args = append(args, *unitID)
	}
	query += ` ORDER BY name`
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Resource
	for rows.Next() {
		var res models.Resource
		if err := rows.Scan(&res.ID, &res.CategoryID, &res.UnitID, &res.Name, &res.Description, &res.Quantity,
			&res.SerialNumber, &res.Location, &res.Condition, &res.MinQuantity,
			&res.CreatedAt, &res.UpdatedAt); err != nil {
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
