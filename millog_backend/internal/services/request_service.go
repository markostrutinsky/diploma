package services

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CanApproveRequest(creatorRole models.UserRole, approverRole models.UserRole) bool {
	if approverRole == models.RoleAdmin {
		return true
	}
	allowedApprovers, exists := models.ApprovalMatrix[creatorRole]
	if !exists {
		return false
	}
	for _, role := range allowedApprovers {
		if role == approverRole {
			return true
		}
	}
	return false
}

type RequestService struct {
	requestRepo  *repositories.SupplyRequestRepository
	resourceRepo *repositories.ResourceRepository
	userRepo     *repositories.UserRepository
	dbPool       *pgxpool.Pool
}

func NewRequestService(reqRepo *repositories.SupplyRequestRepository, resRepo *repositories.ResourceRepository, userRepo *repositories.UserRepository, db *pgxpool.Pool) *RequestService {
	return &RequestService{requestRepo: reqRepo, resourceRepo: resRepo, userRepo: userRepo, dbPool: db}
}

func (s *RequestService) Create(ctx context.Context, userID string, req *models.CreateSupplyRequest) (*models.SupplyRequest, error) {
	sr := &models.SupplyRequest{
		CreatedBy:         userID,
		ResourceID:        req.ResourceID,
		Quantity:          req.Quantity,
		Status:            models.RequestPending,
		TargetWarehouseID: req.TargetWarehouseID,
	}

	if err := s.requestRepo.Create(ctx, s.dbPool, sr); err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return sr, nil
}

// Додаємо аргументи userRole та userUnitID
func (s *RequestService) List(ctx context.Context, userRole string, userUnitID *int64) ([]models.SupplyRequest, error) {
	return s.requestRepo.List(ctx, s.dbPool, userRole, userUnitID)
}

func (s *RequestService) Approve(ctx context.Context, requestID, approverID string, approverRole models.UserRole, approved bool, comment string) error {
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	req, err := s.requestRepo.GetByID(ctx, tx, requestID)
	if err != nil {
		return fmt.Errorf("request not found")
	}
	if req.Status != models.RequestPending {
		return fmt.Errorf("request already processed")
	}

	if req.CreatedBy == approverID {
		return errors.New("неможливо погодити власну заявку (конфлікт інтересів)")
	}

	creator, err := s.userRepo.GetByID(ctx, tx, req.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to get creator details: %w", err)
	}

	if !CanApproveRequest(models.UserRole(creator.Role), approverRole) {
		return errors.New("недостатньо прав для погодження заявки цього рівня (порушення субординації)")
	}

	// Просто оновлюємо статус заявки в базі
	if err := s.requestRepo.Approve(ctx, tx, requestID, approverID, approved, comment); err != nil {
		return err
	}

	// ✅ Більше ніяких махінацій з ресурсами! Майно додасться тільки тоді,
	// коли фура приїде і комірник натисне "Прийняти вантаж" (статус COMPLETED).

	return tx.Commit(ctx)
}

func (s *RequestService) GetByID(ctx context.Context, id string) (*models.SupplyRequest, error) {
	return s.requestRepo.GetByID(ctx, s.dbPool, id)
}

func (s *RequestService) Reject(ctx context.Context, id string, comment string) error {
	// Додаємо префікс до коментаря, щоб було зрозуміло, звідки він взявся
	finalComment := "Відхилено: " + comment
	return s.requestRepo.UpdateStatus(ctx, s.dbPool, id, "REJECTED", finalComment)
}

func (s *RequestService) Cancel(ctx context.Context, id string, userID string) error {
	req, err := s.requestRepo.GetByID(ctx, s.dbPool, id)
	if err != nil {
		return err
	}

	// Перевірка безпеки: скасувати може лише той, хто створив!
	if req.CreatedBy != userID {
		return errors.New("ви не можете скасувати чужу заявку")
	}
	if req.Status != "PENDING" {
		return errors.New("неможливо скасувати заявку, яка вже в обробці")
	}

	return s.requestRepo.UpdateStatus(ctx, s.dbPool, id, "CANCELLED", "Скасовано ініціатором")
}

// GetSmartDispatchPreview виконує розрахунок пакування
func (s *RequestService) GetSmartDispatchPreview(ctx context.Context, reqIDs []string) (*models.SmartDispatchResult, error) {
	// 1. Отримуємо заявки з БД
	requests, err := s.requestRepo.GetRequestsForDispatch(ctx, s.dbPool, reqIDs)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання заявок: %v", err)
	}
	if len(requests) == 0 {
		return nil, fmt.Errorf("не знайдено валідних заявок для обробки")
	}

	// 2. Отримуємо доступні авто з БД
	vehicles, err := s.requestRepo.GetAvailableVehicles(ctx, s.dbPool)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання автопарку: %v", err)
	}
	if len(vehicles) == 0 {
		return nil, fmt.Errorf("наразі немає вільних вантажних автомобілів")
	}

	// 3. АЛГОРИТМ First-Fit Decreasing

	// Сортуємо заявки за вагою (від найважчих до найлегших)
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].WeightKg > requests[j].WeightKg
	})

	// Сортуємо авто за вантажопідйомністю (від найбільших до найменших)
	sort.Slice(vehicles, func(i, j int) bool {
		return vehicles[i].MaxWeight > vehicles[j].MaxWeight
	})

	var unassigned []models.RequestItem

	// Розподіляємо вантаж
	for _, req := range requests {
		placed := false
		for i := range vehicles {
			if vehicles[i].UsedWeight+req.WeightKg <= vehicles[i].MaxWeight {
				vehicles[i].Items = append(vehicles[i].Items, req)
				vehicles[i].UsedWeight += req.WeightKg
				placed = true
				break
			}
		}
		if !placed {
			unassigned = append(unassigned, req)
		}
	}

	// 4. Фільтруємо авто, щоб повернути тільки ті, що отримали вантаж
	var activeRoutes []models.VehicleBin
	for _, v := range vehicles {
		if v.UsedWeight > 0 {
			activeRoutes = append(activeRoutes, v)
		}
	}

	// Якщо жодне авто не задіяно (наприклад, вантажі занадто важкі)
	if len(activeRoutes) == 0 {
		return nil, fmt.Errorf("жодна з машин не може вмістити обраний вантаж (перевищення лімітів)")
	}

	result := models.SmartDispatchResult{
		OptimizedRoutes: activeRoutes,
		Unassigned:      unassigned,
	}

	return &result, nil
}
