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
	userRepo     *repositories.UserRepository
}

func NewInventoryService(catRepo *repositories.CategoryRepository, resRepo *repositories.ResourceRepository, userRepo *repositories.UserRepository, db *pgxpool.Pool) *InventoryService {
	return &InventoryService{categoryRepo: catRepo, resourceRepo: resRepo, userRepo: userRepo, dbPool: db}
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
	return s.resourceRepo.AssignResource(ctx, s.dbPool, id, req.UserID, req.Quantity)
}

func (s *InventoryService) GetMyEquipment(ctx context.Context, userID string) ([]models.MyEquipmentItem, error) {
	// Якщо в тебе транзакції обробляються на рівні сервісу, передавай s.db,
	// або якщо репо сам знає свою БД, то без s.db
	return s.resourceRepo.GetMyEquipment(ctx, s.dbPool, userID)
}

func (s *InventoryService) IssueResource(ctx context.Context, commanderUnitID *int64, req models.IssueResourceRequest) error {
	if commanderUnitID == nil {
		return errors.New("ви не прив'язані до жодного підрозділу і не маєте доступу до складів")
	}

	// Перевіряємо, чи солдат реально підпорядковується цьому командиру (рекурсивна перевірка)
	// Використовуємо твій існуючий метод CheckSubordination
	isSubordinate, err := s.userRepo.CheckSubordination(ctx, nil, *commanderUnitID, req.UserID)
	if err != nil {
		return errors.New("помилка перевірки підпорядкування")
	}
	if !isSubordinate {
		return errors.New("ви не можете видати майно цьому військовослужбовцю, він не у вашому підпорядкуванні")
	}

	// Викликаємо репозиторій
	return s.resourceRepo.IssueToUser(ctx, s.dbPool, *commanderUnitID, req.ResourceID, req.UserID, req.Quantity, req.Notes)
}

func (s *InventoryService) CreateShipment(ctx context.Context, req models.CreateShipmentRequest) error {
	// Передаємо весь об'єкт req напряму в репозиторій для виконання транзакції
	return s.resourceRepo.CreateShipment(ctx, s.dbPool, req)
}

func (s *InventoryService) GetByWarehouse(ctx context.Context, warehouseID string) ([]models.InventoryItem, error) {
	if warehouseID == "" {
		return nil, errors.New("не вказано ID складу")
	}
	return s.resourceRepo.GetByWarehouse(ctx, s.dbPool, warehouseID)
}

func (s *InventoryService) ListShipments(ctx context.Context) ([]repositories.ShipmentRecord, error) {
	return s.resourceRepo.ListShipments(ctx, s.dbPool)
}

func (s *InventoryService) ReceiveShipment(ctx context.Context, shipmentID string) error {
	return s.resourceRepo.ReceiveShipment(ctx, s.dbPool, shipmentID)
}
