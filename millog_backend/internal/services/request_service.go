package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"

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
	// Дізнаємось роль автора, щоб заявки адміна одразу йшли у APPROVED
	// (адмін не має над собою ієрархії, тому самопогодження — єдиний шлях).
	creator, err := s.userRepo.GetByID(ctx, s.dbPool, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load creator: %w", err)
	}

	initialStatus := models.RequestPending
	// Автоматично затверджуємо заявки від TENANT_ADMIN, SYSTEM_ADMIN або ADMIN
	if creator.Role == models.RoleTenantAdmin || creator.Role == models.RoleSystemAdmin || creator.Role == models.RoleAdmin {
		initialStatus = models.RequestApproved
	}

	sr := &models.SupplyRequest{
		CreatedBy:         userID,
		ResourceID:        req.ResourceID,
		Quantity:          req.Quantity,
		Status:            initialStatus,
		TargetWarehouseID: req.TargetWarehouseID,
	}

	if creator.Role == models.RoleTenantAdmin || creator.Role == models.RoleSystemAdmin || creator.Role == models.RoleAdmin {
		approvedBy := userID
		sr.ApprovedBy = &approvedBy
		now := time.Now()
		sr.ApprovedAt = &now
	}

	if err := s.requestRepo.Create(ctx, s.dbPool, sr); err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Для адміна одразу проставляємо approved_by/approved_at в БД.
	if creator.Role == models.RoleTenantAdmin || creator.Role == models.RoleSystemAdmin || creator.Role == models.RoleAdmin {
		if err := s.requestRepo.Approve(ctx, s.dbPool, sr.ID, userID, true, ""); err != nil {
			// Не валимо запит — заявка вже створена, просто логуємо.
			return sr, nil
		}
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

// ConfirmSmartDispatch бере затверджений розподіл (vehicle -> request_ids),
// для кожної машини формує окремий CreateShipmentRequest і викликає ту саму
// транзакційну логіку, що й ручна відправка рейсу (з резервуванням майна).
// Повертає кількість успішно створених рейсів.
func (s *RequestService) ConfirmSmartDispatch(ctx context.Context, req *models.SmartDispatchConfirmReq) (int, error) {
	if len(req.Routes) == 0 {
		return 0, errors.New("маршрути не передані")
	}
	if req.Priority == "" {
		req.Priority = "NORMAL"
	}

	// 1. Підвантажуємо всі заявки одним проходом і перевіряємо інваріанти.
	details := make(map[string]*models.SupplyRequest)
	toWarehouseID := ""
	for _, route := range req.Routes {
		for _, rid := range route.RequestIDs {
			if _, ok := details[rid]; ok {
				return 0, fmt.Errorf("заявка %s присутня в кількох маршрутах", rid)
			}
			sr, err := s.requestRepo.GetByID(ctx, s.dbPool, rid)
			if err != nil {
				return 0, fmt.Errorf("не знайдено заявку %s: %w", rid, err)
			}
			if sr.Status != models.RequestApproved && sr.Status != models.RequestPending {
				return 0, fmt.Errorf("заявка %s вже оброблена (status=%s)", rid, sr.Status)
			}
			if sr.TargetWarehouseID == "" {
				return 0, fmt.Errorf("заявка %s не має цільового складу — Smart Розподіл неможливий", rid)
			}
			if toWarehouseID == "" {
				toWarehouseID = sr.TargetWarehouseID
			} else if sr.TargetWarehouseID != toWarehouseID {
				return 0, errors.New("обрані заявки йдуть на різні цільові склади")
			}
			details[rid] = sr
		}
	}

	if toWarehouseID == "" {
		return 0, errors.New("не вдалося визначити цільовий склад")
	}
	if req.FromWarehouseID == toWarehouseID {
		return 0, errors.New("склад відправлення та отримання не можуть збігатися")
	}

	// 2. Формуємо й відправляємо рейси по одному. Якщо котрийсь рейс впаде
	// (не вистачило майна/авто зайняте) — повертаємо скільки встигли, щоб
	// оператор побачив, де саме проблема.
	successfulShipments := 0
	for idx, route := range req.Routes {
		items := make([]models.ShipmentItemRequest, 0, len(route.RequestIDs))
		for _, rid := range route.RequestIDs {
			d := details[rid]
			reqID := d.ID
			items = append(items, models.ShipmentItemRequest{
				ResourceID: d.ResourceID,
				Quantity:   d.Quantity,
				RequestID:  &reqID,
			})
		}

		shipReq := models.CreateShipmentRequest{
			FromWarehouseID: req.FromWarehouseID,
			ToWarehouseID:   toWarehouseID,
			VehicleID:       route.VehicleID,
			Priority:        req.Priority,
			Items:           items,
		}

		if err := s.resourceRepo.CreateShipment(ctx, s.dbPool, shipReq); err != nil {
			return successfulShipments, fmt.Errorf("рейс #%d (машина %s): %w", idx+1, route.VehicleID, err)
		}
		successfulShipments++
	}
	return successfulShipments, nil
}
