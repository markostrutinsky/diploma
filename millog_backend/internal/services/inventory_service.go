package services

import (
	"context"
	"fmt"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryService struct {
	categoryRepo *repositories.CategoryRepository
	resourceRepo *repositories.ResourceRepository
	dbPool       *pgxpool.Pool
}

func NewInventoryService(catRepo *repositories.CategoryRepository, resRepo *repositories.ResourceRepository, db *pgxpool.Pool) *InventoryService {
	return &InventoryService{categoryRepo: catRepo, resourceRepo: resRepo, dbPool: db}
}

func (s *InventoryService) CreateCategory(ctx context.Context, req *models.CreateCategoryRequest) (*models.ResourceCategory, error) {
	c := &models.ResourceCategory{Name: req.Name, Description: req.Description}
	if err := s.categoryRepo.Create(ctx, s.dbPool, c); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return c, nil
}

func (s *InventoryService) ListCategories(ctx context.Context) ([]models.ResourceCategory, error) {
	return s.categoryRepo.List(ctx, s.dbPool)
}

func (s *InventoryService) CreateResource(ctx context.Context, req *models.CreateResourceRequest) (*models.Resource, error) {
	cond := req.Condition
	if cond == "" {
		cond = models.ConditionNew
	}
	res := &models.Resource{
		CategoryID:   req.CategoryID,
		UnitID:       req.UnitID,
		Name:         req.Name,
		Description:  req.Description,
		Quantity:     req.Quantity,
		SerialNumber: req.SerialNumber,
		Location:     req.Location,
		Condition:    cond,
		MinQuantity:  req.MinQuantity,
	}
	if err := s.resourceRepo.Create(ctx, s.dbPool, res); err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	return res, nil
}

func (s *InventoryService) ListResources(ctx context.Context, unitID *int64) ([]models.Resource, error) {
	return s.resourceRepo.List(ctx, s.dbPool, unitID)
}

func (s *InventoryService) GetResource(ctx context.Context, id string) (*models.Resource, error) {
	return s.resourceRepo.GetByID(ctx, s.dbPool, id)
}
