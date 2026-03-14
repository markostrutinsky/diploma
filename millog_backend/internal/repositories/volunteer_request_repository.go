package repositories

import (
	"context"

	"millog_backend/internal/models"
)

type VolunteerRequestRepository struct{}

func NewVolunteerRequestRepository() *VolunteerRequestRepository {
	return &VolunteerRequestRepository{}
}

func (r *VolunteerRequestRepository) Create(ctx context.Context, db DBExecutor, vr *models.VolunteerRequest) error {
	query := `INSERT INTO volunteer_requests (created_by, title, description, status)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return db.QueryRow(ctx, query, vr.CreatedBy, vr.Title, vr.Description, vr.Status).Scan(&vr.ID, &vr.CreatedAt)
}

func (r *VolunteerRequestRepository) List(ctx context.Context, db DBExecutor) ([]models.VolunteerRequest, error) {
	rows, err := db.Query(ctx, `SELECT id, created_by, title, description, status, taken_by, taken_at, completed_at, created_at
		FROM volunteer_requests ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.VolunteerRequest
	for rows.Next() {
		var vr models.VolunteerRequest
		if err := rows.Scan(&vr.ID, &vr.CreatedBy, &vr.Title, &vr.Description, &vr.Status, &vr.TakenBy, &vr.TakenAt, &vr.CompletedAt, &vr.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, vr)
	}
	return list, rows.Err()
}

func (r *VolunteerRequestRepository) Take(ctx context.Context, db DBExecutor, id, userID string) error {
	query := `UPDATE volunteer_requests SET status = $1, taken_by = $2, taken_at = CURRENT_TIMESTAMP WHERE id = $3 AND status = $4`
	_, err := db.Exec(ctx, query, models.VolunteerTaken, userID, id, models.VolunteerOpen)
	return err
}

func (r *VolunteerRequestRepository) Complete(ctx context.Context, db DBExecutor, id, userID string) error {
	query := `UPDATE volunteer_requests SET status = $1, completed_at = CURRENT_TIMESTAMP WHERE id = $2 AND taken_by = $3 AND status = $4`
	_, err := db.Exec(ctx, query, models.VolunteerCompleted, id, userID, models.VolunteerTaken)
	return err
}
