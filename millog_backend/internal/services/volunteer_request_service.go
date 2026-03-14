package services

import (
	"context"
	"fmt"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VolunteerRequestService struct {
	repo   *repositories.VolunteerRequestRepository
	dbPool *pgxpool.Pool
}

func NewVolunteerRequestService(repo *repositories.VolunteerRequestRepository, db *pgxpool.Pool) *VolunteerRequestService {
	return &VolunteerRequestService{repo: repo, dbPool: db}
}

func (s *VolunteerRequestService) Create(ctx context.Context, userID string, req *models.CreateVolunteerRequest) (*models.VolunteerRequest, error) {
	vr := &models.VolunteerRequest{
		CreatedBy:   userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      models.VolunteerOpen,
	}
	if err := s.repo.Create(ctx, s.dbPool, vr); err != nil {
		return nil, fmt.Errorf("failed to create volunteer request: %w", err)
	}
	return vr, nil
}

func (s *VolunteerRequestService) List(ctx context.Context) ([]models.VolunteerRequest, error) {
	return s.repo.List(ctx, s.dbPool)
}

func (s *VolunteerRequestService) Take(ctx context.Context, id, userID string) error {
	return s.repo.Take(ctx, s.dbPool, id, userID)
}

func (s *VolunteerRequestService) Complete(ctx context.Context, id, userID string) error {
	return s.repo.Complete(ctx, s.dbPool, id, userID)
}
