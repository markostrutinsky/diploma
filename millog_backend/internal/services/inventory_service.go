package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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

// --- НОВИЙ МЕТОД ГЕНЕРАЦІЇ PDF ---
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
		query := `
			INSERT INTO resources (
				category_id, unit_id, warehouse_id, name, description, 
				quantity, unit_type, weight_kg, condition, serial_number, barcode, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW()
			)
		`
		batch.Queue(query,
			finalCategoryID, unitID, warehouseID, row.ResourceName, row.Description,
			row.Quantity, row.UnitType, row.WeightKg, row.Condition, row.SerialNumber, row.Barcode,
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
