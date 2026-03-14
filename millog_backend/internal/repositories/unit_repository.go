package repositories

import (
	"context"

	"millog_backend/internal/models"
)

type UnitRepository struct{}

func NewUnitRepository() *UnitRepository {
	return &UnitRepository{}
}

func (r *UnitRepository) Create(ctx context.Context, db DBExecutor, u *models.Unit) error {
	query := `INSERT INTO units (parent_id, name, unit_type) VALUES ($1, $2, $3) RETURNING id`
	return db.QueryRow(ctx, query, u.ParentID, u.Name, u.UnitType).Scan(&u.ID)
}

func (r *UnitRepository) List(ctx context.Context, db DBExecutor) ([]models.Unit, error) {
	rows, err := db.Query(ctx, `SELECT id, parent_id, name, unit_type FROM units ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Unit
	for rows.Next() {
		var u models.Unit
		if err := rows.Scan(&u.ID, &u.ParentID, &u.Name, &u.UnitType); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, rows.Err()
}
