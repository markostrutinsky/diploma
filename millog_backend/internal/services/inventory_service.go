package services

import (
	"context"
	"errors"
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
		WarehouseID:  req.WarehouseID, // <--- ЗАМІНИЛИ Location на WarehouseID
		Name:         req.Name,
		Description:  req.Description,
		Quantity:     req.Quantity,
		UnitType:     req.UnitType, // <--- ДОДАЛИ одиниці виміру
		SerialNumber: req.SerialNumber,
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

func (s *InventoryService) WriteOff(ctx context.Context, id string, req models.WriteOffResourceRequest) error {
	if req.Quantity <= 0 {
		return errors.New("кількість для списання має бути більшою за нуль")
	}
	return s.resourceRepo.WriteOff(ctx, s.dbPool, id, req.Quantity)
}

// UpdateResource обгортає процес оновлення ресурсу в транзакцію
func (s *InventoryService) UpdateResource(ctx context.Context, id string, req models.UpdateResourceRequest) error {
	// 1. Відкриваємо транзакцію
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// 2. Гарантуємо, що транзакція відкотиться у разі паніки або помилки
	defer tx.Rollback(ctx)

	// 3. Викликаємо репозиторій, ПЕРЕДАЮЧИ ЙОМУ ТРАНЗАКЦІЮ (tx), а не звичайний пул
	err = s.resourceRepo.Update(ctx, tx, id, req)
	if err != nil {
		return fmt.Errorf("repository update failed: %w", err)
	}

	// 4. Якщо все чудово, комітимо зміни в базу
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *InventoryService) Transfer(ctx context.Context, id string, req models.TransferResourceRequest) error {
	// Базова валідація бізнес-логіки
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// 2. Гарантуємо, що транзакція відкотиться у разі паніки або помилки
	defer tx.Rollback(ctx)
	if req.Quantity <= 0 {
		return errors.New("кількість для переміщення має бути більшою за нуль")
	}

	if id == "" {
		return errors.New("не вказано ID ресурсу")
	}

	// Викликаємо репозиторій.
	// s.db - це твій підключений пул бази даних (DBExecutor), який лежить у структурі сервісу.
	err = s.resourceRepo.Transfer(ctx, tx, id, req)
	if err != nil {
		return fmt.Errorf("не вдалося виконати переміщення: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *InventoryService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("не вказано ID ресурсу")
	}
	return s.resourceRepo.Delete(ctx, s.dbPool, id)
}

func (s *InventoryService) Assign(ctx context.Context, id string, req models.AssignResourceRequest) error {
	if req.Quantity <= 0 {
		return errors.New("кількість для видачі має бути більшою за нуль")
	}
	if req.UserID == "" {
		return errors.New("не вказано користувача")
	}
	return s.resourceRepo.Assign(ctx, s.dbPool, id, req.UserID, req.Quantity)
}
