package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

// calcOSRMDistance розраховує відстань у км між двома складами через OSRM.
// Повертає 0 якщо координати відсутні або сервіс недоступний.
func calcOSRMDistance(ctx context.Context, db *pgxpool.Pool, fromWarehouseID, toWarehouseID string) float64 {
	var fromLat, fromLon, toLat, toLon *float64
	db.QueryRow(ctx, "SELECT latitude, longitude FROM warehouses WHERE id = $1", fromWarehouseID).Scan(&fromLat, &fromLon)
	db.QueryRow(ctx, "SELECT latitude, longitude FROM warehouses WHERE id = $1", toWarehouseID).Scan(&toLat, &toLon)

	if fromLat == nil || fromLon == nil || toLat == nil || toLon == nil {
		return 0
	}

	url := fmt.Sprintf("https://router.project-osrm.org/route/v1/driving/%.6f,%.6f;%.6f,%.6f?overview=false",
		*fromLon, *fromLat, *toLon, *toLat)

	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return 0
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	var result struct {
		Routes []struct {
			Distance float64 `json:"distance"`
		} `json:"routes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || len(result.Routes) == 0 {
		return 0
	}
	return result.Routes[0].Distance / 1000.0 // метри → км
}

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
	requestRepo         *repositories.SupplyRequestRepository
	resourceRepo        *repositories.ResourceRepository
	userRepo            *repositories.UserRepository
	dbPool              *pgxpool.Pool
	notificationService *NotificationService
}

func NewRequestService(reqRepo *repositories.SupplyRequestRepository, resRepo *repositories.ResourceRepository, userRepo *repositories.UserRepository, db *pgxpool.Pool, notifService *NotificationService) *RequestService {
	return &RequestService{requestRepo: reqRepo, resourceRepo: resRepo, userRepo: userRepo, dbPool: db, notificationService: notifService}
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
		CreatedBy:          userID,
		ResourceID:         req.ResourceID,
		ResourceName:       req.ResourceName,
		ResourceCategoryID: req.ResourceCategoryID,
		Quantity:           req.Quantity,
		Status:             initialStatus,
		TargetWarehouseID:  req.TargetWarehouseID,
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
func (s *RequestService) GetSmartDispatchPreview(ctx context.Context, reqIDs []string, fromWarehouseID string) (*models.SmartDispatchResult, error) {
	// 1. Отримуємо заявки з БД
	requests, err := s.requestRepo.GetRequestsForDispatch(ctx, s.dbPool, reqIDs)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання заявок: %v", err)
	}
	if len(requests) == 0 {
		return nil, fmt.Errorf("не знайдено валідних заявок для обробки")
	}

	// Визначаємо цільовий склад (усі заявки мають однаковий, фронт це гарантує).
	targetWarehouseID := requests[0].TargetWarehouseID

	// 2. Отримуємо доступні авто з БД — тільки ті, що фізично знаходяться
	// на складі відправника або складі отримувача (без ієрархії).
	vehicles, err := s.requestRepo.GetAvailableVehicles(ctx, s.dbPool, fromWarehouseID, targetWarehouseID)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання автопарку: %v", err)
	}
	if len(vehicles) == 0 {
		if fromWarehouseID != "" {
			return nil, fmt.Errorf("на складах відправника та отримувача немає вільних вантажних автомобілів. Переконайтесь, що транспорт зареєстрований на одному з цих складів")
		}
		return nil, fmt.Errorf("на складі отримувача немає вільних вантажних автомобілів")
	}

	// 3. АЛГОРИТМ First-Fit Decreasing

	// Сортуємо заявки за вагою (від найважчих до найлегших)
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].WeightKg > requests[j].WeightKg
	})

	// Сортуємо авто за вантажопідйомністю від НАЙМЕНШИХ до НАЙБІЛЬШИХ.
	// Це гарантує, що легкий вантаж отримає найменшу відповідну машину,
	// а не одразу 22-тонну фуру (FFD з "best-fit" вибором біна).
	sort.Slice(vehicles, func(i, j int) bool {
		return vehicles[i].MaxWeight < vehicles[j].MaxWeight
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
			if sr.Status != models.RequestApproved && sr.Status != models.RequestPending && sr.Status != models.RequestLoading {
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

			// Якщо в заявці не вказано конкретний resource_id (новий підхід - тільки назва),
			// знаходимо цей ресурс на складі відправника
			var actualResourceID string
			if d.ResourceID != nil && *d.ResourceID != "" {
				// Старий підхід - є конкретний resource_id
				actualResourceID = *d.ResourceID
			} else {
				// Новий підхід - шукаємо ресурс за назвою на складі відправника
				foundRes, err := s.resourceRepo.GetByNameAndWarehouse(ctx, s.dbPool, d.ResourceName, req.FromWarehouseID)
				if err != nil {
					return successfulShipments, fmt.Errorf("не знайдено ресурс '%s' на складі відправника: %w", d.ResourceName, err)
				}
				if foundRes.Quantity < d.Quantity {
					return successfulShipments, fmt.Errorf("на складі недостатньо '%s': потрібно %d, є %d", d.ResourceName, d.Quantity, foundRes.Quantity)
				}
				actualResourceID = foundRes.ID
			}

			items = append(items, models.ShipmentItemRequest{
				ResourceID: actualResourceID,
				Quantity:   d.Quantity,
				RequestID:  &reqID,
			})
		}

		// Розраховуємо відстань між складами (OSRM), якщо не передана у payload
		distanceKm := req.DistanceKm
		if distanceKm == 0 {
			distanceKm = calcOSRMDistance(ctx, s.dbPool, req.FromWarehouseID, toWarehouseID)
		}

		shipReq := models.CreateShipmentRequest{
			FromWarehouseID: req.FromWarehouseID,
			ToWarehouseID:   toWarehouseID,
			VehicleID:       route.VehicleID,
			Priority:        req.Priority,
			Items:           items,
			DistanceKm:      distanceKm,
		}

		if err := s.resourceRepo.CreateShipment(ctx, s.dbPool, shipReq); err != nil {
			return successfulShipments, fmt.Errorf("рейс #%d (машина %s): %w", idx+1, route.VehicleID, err)
		}
		successfulShipments++

		// Відправляємо повідомлення водію про новий рейс
		if s.notificationService != nil {
			// Отримуємо інформацію про водія машини
			var driverID *string
			err := s.dbPool.QueryRow(ctx, "SELECT driver_id FROM vehicles WHERE id = $1", route.VehicleID).Scan(&driverID)

			if err == nil && driverID != nil && *driverID != "" {
				// Отримуємо tenant_id водія
				var driverTenantID string
				err = s.dbPool.QueryRow(ctx, "SELECT tenant_id FROM users WHERE id = $1", *driverID).Scan(&driverTenantID)

				log.Printf("DEBUG: Smart Dispatch - driverID=%s, tenantID=%s, err=%v", *driverID, driverTenantID, err)

				if err == nil && driverTenantID != "" {
					// Отримуємо назви складів
					var fromWarehouse, toWarehouse string
					_ = s.dbPool.QueryRow(ctx, "SELECT name FROM warehouses WHERE id = $1", req.FromWarehouseID).Scan(&fromWarehouse)
					_ = s.dbPool.QueryRow(ctx, "SELECT name FROM warehouses WHERE id = $1", toWarehouseID).Scan(&toWarehouse)

					// Отримуємо ID щойно створеного shipment
					var shipmentID string
					_ = s.dbPool.QueryRow(ctx, "SELECT id FROM shipments WHERE vehicle_id = $1 ORDER BY created_at DESC LIMIT 1", route.VehicleID).Scan(&shipmentID)

					if shipmentID != "" {
						// Створюємо контекст з tenant_id водія
						driverCtx := context.WithValue(ctx, "tenant_id", driverTenantID)
						notifErr := s.notificationService.NotifyDriverAboutShipment(driverCtx, *driverID, shipmentID, fromWarehouse, toWarehouse)
						log.Printf("DEBUG: Smart Dispatch notification result for driver %s - err=%v", *driverID, notifErr)
					}
				}
			}
		}
	}
	return successfulShipments, nil
}
