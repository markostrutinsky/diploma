package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type InventoryService struct {
	categoryRepo        *repositories.CategoryRepository
	resourceRepo        *repositories.ResourceRepository
	dbPool              *pgxpool.Pool
	userRepo            *repositories.UserRepository
	notificationService *NotificationService
}

func NewInventoryService(catRepo *repositories.CategoryRepository, resRepo *repositories.ResourceRepository, userRepo *repositories.UserRepository, db *pgxpool.Pool, notifService *NotificationService) *InventoryService {
	return &InventoryService{categoryRepo: catRepo, resourceRepo: resRepo, userRepo: userRepo, dbPool: db, notificationService: notifService}
}

// GetDB повертає пул підключень до бази даних для використання в хендлерах.
func (s *InventoryService) GetDB() *pgxpool.Pool {
	return s.dbPool
}

// ReportEquipment — подає рапорт про майно (зламано/втрачено).
// Знаходить керівника tenant'у (TENANT_ADMIN або найвищий unit_id=NULL) і надсилає нотифікацію.
// Якщо заявник сам є TENANT_ADMIN — нотифікація йде всім логістам у tenant'і.
func (s *InventoryService) ReportEquipment(ctx context.Context, requesterID string, assignmentID string, reason string) error {
	// 1. Беремо дані заявника та назву майна
	var requesterName, resourceName, tenantID string
	var requesterRole string
	err := s.dbPool.QueryRow(ctx, `
		SELECT u.full_name, r.name, u.tenant_id, u.role
		FROM resource_assignments ra
		JOIN resources r ON r.id = ra.resource_id
		JOIN users u ON u.id = $1
		WHERE ra.id = $2 AND ra.user_id = $1
	`, requesterID, assignmentID).Scan(&requesterName, &resourceName, &tenantID, &requesterRole)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("майно не знайдено або не належить вам")
		}
		return fmt.Errorf("помилка отримання даних: %w", err)
	}

	reasonLabel := map[string]string{
		"BROKEN":   "Зламано / Пошкоджено",
		"LOST":     "Втрачено",
		"WORN_OUT": "Зношено (термін експлуатації)",
	}[reason]
	if reasonLabel == "" {
		reasonLabel = reason
	}

	title := fmt.Sprintf("⚠️ Рапорт щодо майна: %s", resourceName)
	message := fmt.Sprintf("%s подав(ла) рапорт на списання: %s. Причина: %s.", requesterName, resourceName, reasonLabel)

	// 2. Отримуємо список отримувачів нотифікації
	// Якщо заявник — TENANT_ADMIN: сповіщаємо всіх логістів у tenant'і
	// Інакше: сповіщаємо TENANT_ADMIN та регіональних логістів
	var recipientRows interface{ Scan(...any) error }
	var rows interface {
		Next() bool
		Scan(...any) error
		Close()
	}

	if requesterRole == "TENANT_ADMIN" || requesterRole == "ADMIN" {
		// Адмін подає рапорт — сповіщаємо логістів
		r, err2 := s.dbPool.Query(ctx, `
			SELECT id FROM users
			WHERE tenant_id = $1
			AND role IN ('REGION_LOGISTICIAN','BRANCH_LOGISTICIAN','REGION_STOREKEEPER','BRANCH_STOREKEEPER')
			AND status = 'ACTIVE'
			LIMIT 10
		`, tenantID)
		if err2 != nil {
			return err2
		}
		rows = r
		_ = recipientRows
	} else {
		// Звичайний користувач — сповіщаємо TENANT_ADMIN / адміністраторів
		r, err2 := s.dbPool.Query(ctx, `
			SELECT id FROM users
			WHERE tenant_id = $1
			AND role IN ('TENANT_ADMIN','ADMIN')
			AND status = 'ACTIVE'
			LIMIT 5
		`, tenantID)
		if err2 != nil {
			return err2
		}
		rows = r
		_ = recipientRows
	}
	defer rows.Close()

	sent := 0
	for rows.Next() {
		var recipientID string
		if err2 := rows.Scan(&recipientID); err2 != nil {
			continue
		}
		// Не надсилаємо самому собі
		if recipientID == requesterID {
			continue
		}
		if s.notificationService != nil {
			_ = s.notificationService.CreateNotification(ctx, &models.CreateNotificationRequest{
				UserID:    recipientID,
				Type:      models.NotificationEquipmentReport,
				Title:     title,
				Message:   message,
				RelatedID: &assignmentID,
			})
			sent++
		}
	}

	// Якщо нема кому відправити — всеодно вважаємо успіхом (нема логістів ще в системі)
	_ = sent
	return nil
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
		WeightKg:     req.WeightKg,
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

// IsUnitInSubtree перевіряє, чи targetUnitID є самим userUnitID або одним з його нащадків
// (дочірні, онучаті підрозділи тощо). Повертає true якщо так.
// Це запобігає ситуації, коли керівник філії може видавати/списувати
// ресурси, що належать регіональному чи тенантному рівню.
func (s *InventoryService) IsUnitInSubtree(ctx context.Context, userUnitID, targetUnitID int64) (bool, error) {
	if userUnitID == targetUnitID {
		return true, nil
	}
	var count int
	err := s.dbPool.QueryRow(ctx, `
		WITH RECURSIVE subtree AS (
			SELECT id FROM units WHERE id = $1
			UNION ALL
			SELECT u.id FROM units u JOIN subtree st ON u.parent_id = st.id
		)
		SELECT COUNT(*) FROM subtree WHERE id = $2
	`, userUnitID, targetUnitID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
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

func (s *InventoryService) IssueResource(ctx context.Context, commanderUnitID *int64, role string, req models.IssueResourceRequest) error {
	// Адміністраторські ролі можуть видавати без прив'язки до підрозділу
	isPrivileged := role == "TENANT_ADMIN" || role == "ADMIN" || role == "SYSTEM_ADMIN"

	if !isPrivileged {
		if commanderUnitID == nil {
			return errors.New("ви не прив'язані до жодного підрозділу і не маєте доступу до складів")
		}

		// Перевіряємо, чи солдат реально підпорядковується цьому командиру (рекурсивна перевірка)
		isSubordinate, err := s.userRepo.CheckSubordination(ctx, nil, *commanderUnitID, req.UserID)
		if err != nil {
			return errors.New("помилка перевірки підпорядкування")
		}
		if !isSubordinate {
			return errors.New("ви не можете видати майно цьому військовослужбовцю, він не у вашому підпорядкуванні")
		}
	}

	var unitID int64
	if commanderUnitID != nil {
		unitID = *commanderUnitID
	}

	// Викликаємо репозиторій
	return s.resourceRepo.IssueToUser(ctx, s.dbPool, unitID, req.ResourceID, req.UserID, req.Quantity, req.Notes, req.WarehouseID)
}

func (s *InventoryService) CreateShipment(ctx context.Context, req models.CreateShipmentRequest) error {
	// Якщо distance_km не передано з фронтенду — розрахуємо через OSRM
	if req.DistanceKm == 0 {
		req.DistanceKm = s.calcOSRMDistance(ctx, req.FromWarehouseID, req.ToWarehouseID)
	}
	// Передаємо весь об'єкт req напряму в репозиторій для виконання транзакції
	err := s.resourceRepo.CreateShipment(ctx, s.dbPool, req)
	if err != nil {
		return err
	}

	// Отримуємо інформацію про водія машини для сповіщення
	var driverID *string
	err = s.dbPool.QueryRow(ctx, "SELECT driver_id FROM vehicles WHERE id = $1", req.VehicleID).Scan(&driverID)

	log.Printf("DEBUG: CreateShipment - vehicle_id=%s, err=%v, driverID=%v", req.VehicleID, err, driverID)

	// Якщо є водій, створюємо сповіщення
	if err == nil && driverID != nil && *driverID != "" {
		// Отримуємо tenant_id водія для коректного створення повідомлення
		var driverTenantID string
		err = s.dbPool.QueryRow(ctx, "SELECT tenant_id FROM users WHERE id = $1", *driverID).Scan(&driverTenantID)

		log.Printf("DEBUG: Driver info - driverID=%s, tenantID=%s, err=%v", *driverID, driverTenantID, err)

		if err == nil && driverTenantID != "" {
			// Отримуємо назви складів для повідомлення
			var fromWarehouse, toWarehouse string
			_ = s.dbPool.QueryRow(ctx, "SELECT name FROM warehouses WHERE id = $1", req.FromWarehouseID).Scan(&fromWarehouse)
			_ = s.dbPool.QueryRow(ctx, "SELECT name FROM warehouses WHERE id = $1", req.ToWarehouseID).Scan(&toWarehouse)

			// Отримуємо ID щойно створеного shipment (останній створений для цього транспорту)
			var shipmentID string
			_ = s.dbPool.QueryRow(ctx, "SELECT id FROM shipments WHERE vehicle_id = $1 ORDER BY created_at DESC LIMIT 1", req.VehicleID).Scan(&shipmentID)

			log.Printf("DEBUG: NotifyDriver - driverID=%s, shipmentID=%s, from=%s, to=%s", *driverID, shipmentID, fromWarehouse, toWarehouse)

			// Створюємо новий контекст з tenant_id водія
			if s.notificationService != nil && shipmentID != "" {
				// Імпортуємо функцію WithTenant з repositories
				driverCtx := context.WithValue(ctx, "tenant_id", driverTenantID)
				notifErr := s.notificationService.NotifyDriverAboutShipment(driverCtx, *driverID, shipmentID, fromWarehouse, toWarehouse)
				log.Printf("DEBUG: Notification result - err=%v", notifErr)
			} else {
				log.Printf("DEBUG: Skipping notification - service=%v, shipmentID=%s", s.notificationService != nil, shipmentID)
			}
		}
	}

	return nil
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

func (s *InventoryService) ListMyShipments(ctx context.Context, userID string) ([]repositories.ShipmentRecord, error) {
	return s.resourceRepo.ListMyShipments(ctx, s.dbPool, userID)
}

func (s *InventoryService) ListShipmentsByVehicle(ctx context.Context, vehicleID string) ([]repositories.ShipmentRecord, error) {
	return s.resourceRepo.ListShipmentsByVehicle(ctx, s.dbPool, vehicleID)
}

func (s *InventoryService) ReceiveShipment(ctx context.Context, shipmentID string, actualKm float64) error {
	return s.resourceRepo.ReceiveShipment(ctx, s.dbPool, shipmentID, actualKm)
}

func (s *InventoryService) StartShipment(ctx context.Context, shipmentID string) error {
	return s.resourceRepo.StartShipment(ctx, s.dbPool, shipmentID)
}

// LogShipmentRefuel — реєструє дозаправку під час рейсу.
// Одночасно пише в shipment_refuels І в fuel_records як REFUEL.
func (s *InventoryService) LogShipmentRefuel(ctx context.Context, shipmentID string, userID string, req *models.LogShipmentRefuelRequest) (*models.ShipmentRefuel, error) {
	// Отримуємо vehicle_id і статус рейсу
	var vehicleID string
	var status string
	err := s.dbPool.QueryRow(ctx,
		`SELECT vehicle_id, status FROM shipments WHERE id = $1`, shipmentID,
	).Scan(&vehicleID, &status)
	if err != nil {
		return nil, fmt.Errorf("рейс не знайдено")
	}
	if status != "IN_TRANSIT" {
		return nil, fmt.Errorf("дозаправку можна додати лише під час активного рейсу (статус: %s)", status)
	}

	// Перевіряємо місткість баку
	var tankCapacity float64
	var currentBalance float64
	_ = s.dbPool.QueryRow(ctx, `SELECT tank_capacity FROM vehicles WHERE id = $1`, vehicleID).Scan(&tankCapacity)
	_ = s.dbPool.QueryRow(ctx, `
		SELECT GREATEST(0, COALESCE(SUM(CASE WHEN record_type = 'REFUEL' THEN liters ELSE -liters END), 0))
		FROM fuel_records WHERE vehicle_id = $1`, vehicleID,
	).Scan(&currentBalance)

	liters := req.Liters
	if tankCapacity > 0 {
		maxRefuel := tankCapacity - currentBalance
		if maxRefuel <= 0 {
			return nil, fmt.Errorf("бак повний (%.1f / %.1f л)", currentBalance, tankCapacity)
		}
		if liters > maxRefuel {
			liters = maxRefuel // тихо кліпуємо до максимуму
		}
	}

	tenantID := ""
	_ = s.dbPool.QueryRow(ctx, `SELECT tenant_id FROM shipments WHERE id = $1`, shipmentID).Scan(&tenantID)

	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 1. Запис у shipment_refuels
	refuel := &models.ShipmentRefuel{}
	var uid *string
	if userID != "" {
		uid = &userID
	}
	var tid *string
	if tenantID != "" {
		tid = &tenantID
	}
	err = tx.QueryRow(ctx, `
		INSERT INTO shipment_refuels (shipment_id, vehicle_id, liters, odometer_km, station_name, cost_uah, created_by, tenant_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8::uuid)
		RETURNING id, shipment_id, vehicle_id, liters, odometer_km, station_name, cost_uah, created_by, tenant_id, created_at
	`, shipmentID, vehicleID, liters, req.OdometerKm, req.StationName, req.CostUAH, uid, tid,
	).Scan(&refuel.ID, &refuel.ShipmentID, &refuel.VehicleID, &refuel.Liters,
		&refuel.OdometerKm, &refuel.StationName, &refuel.CostUAH, &refuel.CreatedBy, &refuel.TenantID, &refuel.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("помилка збереження дозаправки: %w", err)
	}

	// 2. Запис у fuel_records як REFUEL (для балансу пального)
	if tenantID != "" {
		_, err = tx.Exec(ctx, `
			INSERT INTO fuel_records (vehicle_id, liters, odometer_km, record_type, created_by, is_anomaly, tenant_id)
			VALUES ($1, $2, $3, 'REFUEL', $4, false, $5::uuid)
		`, vehicleID, liters, req.OdometerKm, uid, tenantID)
	} else {
		_, err = tx.Exec(ctx, `
			INSERT INTO fuel_records (vehicle_id, liters, odometer_km, record_type, created_by, is_anomaly)
			VALUES ($1, $2, $3, 'REFUEL', $4, false)
		`, vehicleID, liters, req.OdometerKm, uid)
	}
	if err != nil {
		return nil, fmt.Errorf("помилка запису в fuel_records: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return refuel, nil
}

// GetShipmentRefuels — повертає всі дозаправки для конкретного рейсу.
func (s *InventoryService) GetShipmentRefuels(ctx context.Context, shipmentID string) ([]*models.ShipmentRefuel, error) {
	rows, err := s.dbPool.Query(ctx, `
		SELECT id, shipment_id, vehicle_id, liters, odometer_km, station_name, cost_uah, created_by, tenant_id, created_at
		FROM shipment_refuels
		WHERE shipment_id = $1
		ORDER BY created_at ASC
	`, shipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.ShipmentRefuel
	for rows.Next() {
		r := &models.ShipmentRefuel{}
		if err := rows.Scan(&r.ID, &r.ShipmentID, &r.VehicleID, &r.Liters,
			&r.OdometerKm, &r.StationName, &r.CostUAH, &r.CreatedBy, &r.TenantID, &r.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}
func (s *InventoryService) GenerateShipmentPDF(ctx context.Context, shipmentID string) ([]byte, error) {
	info, err := s.resourceRepo.GetShipmentInfo(ctx, s.dbPool, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання інфо про рейс: %w", err)
	}

	items, err := s.resourceRepo.GetShipmentItems(ctx, s.dbPool, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання маніфесту: %w", err)
	}

	// ФІКС ЧАСОВОГО ПОЯСУ: Переводимо час з UTC у Київський
	loc, errLoc := time.LoadLocation("Europe/Kyiv")
	if errLoc == nil {
		info.CreatedAt = info.CreatedAt.In(loc)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("Roboto", "", "fonts/Roboto-Regular.ttf")
	pdf.AddPage()

	pdf.SetFont("Roboto", "", 16)
	pdf.CellFormat(190, 10, "ТОВАРНО-ТРАНСПОРТНА НАКЛАДНА (МАРШРУТНИЙ ЛИСТ)", "0", 1, "C", false, 0, "")

	pdf.SetFont("Roboto", "", 11)
	pdf.SetTextColor(100, 100, 100)

	shortID := strings.ToUpper(strings.Split(shipmentID, "-")[0])
	pdf.CellFormat(190, 6, fmt.Sprintf("Рейс №: %s", shortID), "0", 1, "C", false, 0, "")

	// ЗАЛИШАЄМО ТІЛЬКИ ОДНУ ДАТУ (Дату фактичного відправлення рейсу)
	pdf.CellFormat(190, 6, fmt.Sprintf("Дата відправлення: %s", info.CreatedAt.Format("02.01.2006 15:04")), "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Roboto", "", 12)
	pdf.CellFormat(190, 8, fmt.Sprintf("ВІДПРАВНИК (Склад): %s", info.FromWarehouse), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("ОДЕРЖУВАЧ (Склад): %s", info.ToWarehouse), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("ТРАНСПОРТНИЙ ЗАСІБ: %s", info.Vehicle), "0", 1, "L", false, 0, "")
	pdf.Ln(8)

	pdf.SetFillColor(230, 230, 230)
	pdf.SetFont("Roboto", "", 11)
	pdf.CellFormat(15, 10, "№", "1", 0, "C", true, 0, "")
	pdf.CellFormat(115, 10, "Найменування майна", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Кількість", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Од. вим.", "1", 1, "C", true, 0, "")

	pdf.SetFont("Roboto", "", 11)
	for idx, item := range items {
		// ПЕРЕКЛАДАЄМО ОДИНИЦІ ВИМІРУ
		unitTranslated := item.Unit
		switch item.Unit {
		case "PCS":
			unitTranslated = "шт"
		case "KIT":
			unitTranslated = "компл"
		case "KG":
			unitTranslated = "кг"
		case "L":
			unitTranslated = "л"
		}

		pdf.CellFormat(15, 10, fmt.Sprintf("%d", idx+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(115, 10, "  "+item.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", item.Qty), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, unitTranslated, "1", 1, "C", false, 0, "")
	}
	pdf.Ln(25)

	pdf.SetFont("Roboto", "", 12)
	pdf.CellFormat(95, 10, "ЗДАВ (Відправник): ___________________", "0", 0, "L", false, 0, "")
	pdf.CellFormat(95, 10, "ПРИЙНЯВ (Одержувач): ___________________", "0", 1, "L", false, 0, "")
	pdf.Ln(10)
	pdf.CellFormat(190, 10, "ВОДІЙ-ЕКСПЕДИТОР: ___________________", "0", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("помилка рендеру PDF: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *InventoryService) GenerateResourceQR(resourceID string) ([]byte, error) {
	content := fmt.Sprintf("Omnilog-resource:%s", resourceID)

	pngImage, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("помилка генерації QR-коду: %w", err)
	}

	return pngImage, nil
}

// UpdateCategory передає дані в репозиторій
func (s *InventoryService) UpdateCategory(ctx context.Context, id, name, description string) error {
	return s.categoryRepo.Update(ctx, s.dbPool, id, name, description)
}

// DeleteCategory перевіряє, чи порожня категорія, і лише тоді видаляє
func (s *InventoryService) DeleteCategory(ctx context.Context, id string) error {
	var count int
	query := `SELECT count(*) FROM resources WHERE category_id = $1`

	err := s.dbPool.QueryRow(ctx, query, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("категорія не порожня")
	}

	return s.categoryRepo.Delete(ctx, s.dbPool, id)
}

func (s *InventoryService) SubmitInventoryAudit(ctx context.Context, userID string, req models.SubmitAuditRequest) error {
	// Можна додати додаткові перевірки, наприклад, чи існує такий склад
	if req.WarehouseID == "" {
		return fmt.Errorf("warehouse ID is required")
	}

	return s.resourceRepo.SubmitInventoryAudit(ctx, s.dbPool, userID, req)
}

func normalizeText(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

type categoryItem struct {
	ID   int64
	Name string
}

func findBestCategoryMatch(importName string, categories []categoryItem) *int64 {
	cleanImport := normalizeText(importName)

	// 1. Точний збіг
	for _, cat := range categories {
		if normalizeText(cat.Name) == cleanImport {
			return &cat.ID
		}
	}

	// 2. Частковий збіг
	for _, cat := range categories {
		cleanDB := normalizeText(cat.Name)
		if strings.Contains(cleanDB, cleanImport) || strings.Contains(cleanImport, cleanDB) {
			return &cat.ID
		}
	}

	return nil
}

// ---------------------------------------------------------
// ГОЛОВНА ФУНКЦІЯ ІМПОРТУ
// ---------------------------------------------------------

func (s *RequestService) ProcessBulkImport(ctx context.Context, rows []models.ImportResourceRow) (*models.BulkImportResponse, error) {
	resp := &models.BulkImportResponse{
		TotalProcessed: len(rows),
		Errors:         []string{},
	}

	// 1. ЗАВАНТАЖУЄМО ДОВІДНИКИ В ПАМ'ЯТЬ (щоб не смикати БД на кожному рядку)

	// Підрозділи (Ключ - нормалізована назва, Значення - ID)
	unitMap := make(map[string]int64)
	uRows, err := s.dbPool.Query(ctx, "SELECT id, name FROM units")
	if err == nil {
		for uRows.Next() {
			var id int64
			var name string
			uRows.Scan(&id, &name)
			unitMap[normalizeText(name)] = id
		}
		uRows.Close()
	}

	// Склади (Ключ - "unitID_normalizedWarehouseName", Значення - ID)
	warehouseMap := make(map[string]int64)
	wRows, err := s.dbPool.Query(ctx, "SELECT id, unit_id, name FROM warehouses")
	if err == nil {
		for wRows.Next() {
			var id, unitID int64
			var name string
			wRows.Scan(&id, &unitID, &name)
			key := fmt.Sprintf("%d_%s", unitID, normalizeText(name))
			warehouseMap[key] = id
		}
		wRows.Close()
	}

	// Категорії (Масив для розумного пошуку)
	var categories []categoryItem
	cRows, err := s.dbPool.Query(ctx, "SELECT id, name FROM categories")
	if err == nil {
		for cRows.Next() {
			var cat categoryItem
			cRows.Scan(&cat.ID, &cat.Name)
			categories = append(categories, cat)
		}
		cRows.Close()
	}

	// 2. ПІДГОТОВКА ПАКЕТНОГО ЗАПИТУ (Batch)
	batch := &pgx.Batch{}
	var validRowsCount int

	for i, row := range rows {
		rowNumber := i + 2 // +2 бо індекс з 0, і перший рядок це зазвичай заголовки в CSV

		// -- Перевірка Організації --
		unitID, unitExists := unitMap[normalizeText(row.UnitName)]
		if !unitExists {
			resp.Errors = append(resp.Errors, fmt.Sprintf("Рядок %d: Організацію '%s' не знайдено", rowNumber, row.UnitName))
			resp.Failed++
			continue
		}

		// -- Перевірка Складу --
		wKey := fmt.Sprintf("%d_%s", unitID, normalizeText(row.WarehouseName))
		warehouseID, warehouseExists := warehouseMap[wKey]
		if !warehouseExists {
			resp.Errors = append(resp.Errors, fmt.Sprintf("Рядок %d: Склад '%s' не знайдено в організації '%s'", rowNumber, row.WarehouseName, row.UnitName))
			resp.Failed++
			continue
		}

		// -- Розумна обробка Категорії --
		var finalCategoryID int64
		matchedCatID := findBestCategoryMatch(row.CategoryName, categories)

		if matchedCatID != nil {
			finalCategoryID = *matchedCatID
		} else {
			// Створюємо нову категорію прямо зараз, щоб наступні рядки її вже бачили
			err := s.dbPool.QueryRow(ctx, "INSERT INTO categories (name) VALUES ($1) RETURNING id", row.CategoryName).Scan(&finalCategoryID)
			if err != nil {
				resp.Errors = append(resp.Errors, fmt.Sprintf("Рядок %d: Помилка створення категорії '%s'", rowNumber, row.CategoryName))
				resp.Failed++
				continue
			}
			// Додаємо в пам'ять
			categories = append(categories, categoryItem{ID: finalCategoryID, Name: row.CategoryName})
		}

		// -- Додаємо валідний запис у Batch --
		tid := repositories.TenantFromCtx(ctx)
		if tid == "" {
			resp.Errors = append(resp.Errors, fmt.Sprintf("Рядок %d: tenant_id відсутній у контексті", rowNumber))
			resp.Failed++
			continue
		}

		query := `
			INSERT INTO resources (
				category_id, unit_id, warehouse_id, name, description, 
				quantity, unit_type, weight_kg, condition, serial_number, barcode, tenant_id, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW()
			)
		`
		batch.Queue(query,
			finalCategoryID, unitID, warehouseID, row.ResourceName, row.Description,
			row.Quantity, row.UnitType, row.WeightKg, row.Condition, row.SerialNumber, row.Barcode, tid,
		)
		validRowsCount++
	}

	// 3. ВИКОНАННЯ ПАКЕТНОГО ЗАПИТУ (Bulk Insert)
	if validRowsCount > 0 {
		br := s.dbPool.SendBatch(ctx, batch)
		defer br.Close()

		for i := 0; i < validRowsCount; i++ {
			_, err := br.Exec()
			if err != nil {
				resp.Errors = append(resp.Errors, fmt.Sprintf("Помилка БД під час збереження запису: %v", err))
				resp.Failed++
			} else {
				resp.Successfully++
			}
		}
	}

	return resp, nil
}

func (s *InventoryService) GetAllCategories(ctx context.Context) ([]models.ResourceCategory, error) {
	categories, err := s.categoryRepo.GetAll(ctx, s.dbPool)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання категорій: %w", err)
	}
	return categories, nil
}

func (s *InventoryService) GetResourceByNameAndWarehouse(ctx context.Context, name string, warehouseID string) (*models.Resource, error) {
	return s.resourceRepo.GetByNameAndWarehouse(ctx, s.dbPool, name, warehouseID)
}

func (s *InventoryService) UpdateResourceQuantity(ctx context.Context, id string, newQuantity int) error {
	return s.resourceRepo.UpdateQuantity(ctx, s.dbPool, id, newQuantity)
}

// calcOSRMDistance — розраховує відстань між двома складами через публічний OSRM API.
func (s *InventoryService) calcOSRMDistance(ctx context.Context, fromWarehouseID, toWarehouseID string) float64 {
	var fromLat, fromLon, toLat, toLon *float64
	s.dbPool.QueryRow(ctx, "SELECT latitude, longitude FROM warehouses WHERE id = $1", fromWarehouseID).Scan(&fromLat, &fromLon)
	s.dbPool.QueryRow(ctx, "SELECT latitude, longitude FROM warehouses WHERE id = $1", toWarehouseID).Scan(&toLat, &toLon)

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
