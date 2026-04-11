package services

import (
	"context"
	"fmt"
	"millog_backend/internal/models"
	"millog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WarehouseService struct {
	repo   *repositories.WarehouseRepository
	dbPool *pgxpool.Pool
}

func NewWarehouseService(repo *repositories.WarehouseRepository, db *pgxpool.Pool) *WarehouseService {
	return &WarehouseService{repo: repo, dbPool: db}
}

func (s *WarehouseService) CreateWarehouse(ctx context.Context, req *models.CreateWarehouseRequest) (*models.Warehouse, error) {
	w := &models.Warehouse{
		UnitID:       req.UnitID,
		Name:         req.Name,
		LocationType: req.LocationType,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
	}

	if err := s.repo.Create(ctx, s.dbPool, w); err != nil {
		return nil, fmt.Errorf("failed to create warehouse: %w", err)
	}

	return w, nil
}

func (s *WarehouseService) ListWarehouses(ctx context.Context, unitID int64) ([]models.Warehouse, error) {
	return s.repo.ListByUnit(ctx, s.dbPool, unitID)
}

// Додай цей метод до WarehouseService
func (s *WarehouseService) UpdateLocation(ctx context.Context, warehouseID string, lat, lng float64) error {
	return s.repo.UpdateLocation(ctx, s.dbPool, warehouseID, lat, lng)
}
