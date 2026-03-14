package services

import (
	"context"
	"fmt"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitService struct {
	repo   *repositories.UnitRepository
	dbPool *pgxpool.Pool
}

func NewUnitService(repo *repositories.UnitRepository, db *pgxpool.Pool) *UnitService {
	return &UnitService{repo: repo, dbPool: db}
}

func (s *UnitService) Create(ctx context.Context, req *models.CreateUnitRequest) (*models.Unit, error) {
	u := &models.Unit{
		ParentID: req.ParentID,
		Name:     req.Name,
		UnitType: req.UnitType,
	}
	if err := s.repo.Create(ctx, s.dbPool, u); err != nil {
		return nil, fmt.Errorf("failed to create unit: %w", err)
	}
	return u, nil
}

func (s *UnitService) List(ctx context.Context) ([]models.Unit, error) {
	return s.repo.List(ctx, s.dbPool)
}
