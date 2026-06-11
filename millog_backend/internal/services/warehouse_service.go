package services

import (
	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"
	"context"
	"errors"
	"fmt"

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

func (s *WarehouseService) UpdateWarehouse(ctx context.Context, id, name, capacityLevel, zoneType string) error {
	return s.repo.Update(ctx, s.dbPool, id, name, capacityLevel, zoneType)
}

func (s *WarehouseService) DeleteWarehouse(ctx context.Context, id string) error {
	var resourceCount int
	tid := repositories.TenantFromCtx(ctx)
	if tid == "" {
		return errors.New("tenant_id is required for warehouses")
	}
	err := s.dbPool.QueryRow(ctx, `
		SELECT count(*)
		FROM resources r
		WHERE r.warehouse_id = $1 AND r.tenant_id = $2::uuid
	`, id, tid).Scan(&resourceCount)
	if err != nil {
		return err
	}

	if resourceCount > 0 {
		return errors.New("неможливо видалити: на складі ще є майно (або картки майна)")
	}

	return s.repo.Delete(ctx, s.dbPool, id)
}
