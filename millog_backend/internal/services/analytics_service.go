package services

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"
)

type AnalyticsService struct {
	repo *repositories.AnalyticsRepository
	db   repositories.DBExecutor
}

func NewAnalyticsService(repo *repositories.AnalyticsRepository, db repositories.DBExecutor) *AnalyticsService {
	return &AnalyticsService{repo: repo, db: db}
}

func (s *AnalyticsService) GetDashboardAnalytics(ctx context.Context, start, end, unitID string) (*models.DashboardAnalytics, error) {
	return s.repo.GetDashboardStats(ctx, s.db, start, end, unitID)
}

// НОВА ФУНКЦІЯ: Приймає налаштування замовлення з фронтенду та ID користувача
func (s *AnalyticsService) RunSmartReplenish(ctx context.Context, req models.SmartReplenishRequest, userID string) (int, error) {
	return s.repo.ProcessSmartReplenish(ctx, s.db, req, userID)
}

// GenerateInventoryExcel формує XLSX з поточними залишками
func (s *AnalyticsService) GenerateInventoryExcel(ctx context.Context, unitID *int) ([]byte, error) {
	data, err := s.repo.GetInventoryForExport(ctx, s.db, unitID)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання даних: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Залишки на складах"
	f.SetSheetName("Sheet1", sheetName)

	// Стиль для шапки (жирний шрифт, сірий фон)
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
	})

	headers := []string{"Філія / Підрозділ", "Склад", "Категорія", "Найменування майна", "Кількість", "Од. вим.", "Стан"}
	for col, header := range headers {
		colName, _ := excelize.ColumnNumberToName(col + 1)
		cell := fmt.Sprintf("%s1", colName)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
		f.SetColWidth(sheetName, colName, colName, 20) // Задаємо ширину колонок
	}
	f.SetColWidth(sheetName, "D", "D", 40) // Назва майна ширша

	// Заповнюємо дані
	for i, row := range data {
		rowIndex := i + 2 // Рядок 1 - це шапка

		// Переклад значень для бухгалтерії
		unitType := row.UnitType
		switch unitType {
		case "PCS":
			unitType = "шт"
		case "KIT":
			unitType = "компл"
		case "KG":
			unitType = "кг"
		case "L":
			unitType = "л"
		}

		condition := "Нове"
		if row.Condition == "USED" {
			condition = "Вживане"
		} else if row.Condition == "WRITTEN_OFF" {
			condition = "Списано"
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), row.UnitName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), row.Warehouse)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), row.Category)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), row.ItemName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), row.Quantity)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIndex), unitType)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIndex), condition)
	}

	// Додаємо автофільтр на всі дані
	f.AutoFilter(sheetName, fmt.Sprintf("A1:G%d", len(data)+1), nil)

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateFuelExcel формує XLSX зі звітом по пальному
func (s *AnalyticsService) GenerateFuelExcel(ctx context.Context, startDate, endDate time.Time) ([]byte, error) {
	data, err := s.repo.GetFuelForExport(ctx, s.db, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("помилка отримання даних: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Звіт по пальному"
	f.SetSheetName("Sheet1", sheetName)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E26B0A"}, Pattern: 1},
	})

	headers := []string{"Дата та Час", "Транспортний засіб", "Держ. номер", "Тип операції", "Літри", "Відповідальний"}
	for col, header := range headers {
		colName, _ := excelize.ColumnNumberToName(col + 1)
		cell := fmt.Sprintf("%s1", colName)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
		f.SetColWidth(sheetName, colName, colName, 22)
	}

	loc, _ := time.LoadLocation("Europe/Kyiv")

	for i, row := range data {
		rowIndex := i + 2

		opType := "Витрата (Списання)"
		if row.RecordType == "REFUEL" {
			opType = "Заправка (Надходження)"
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), row.Date.In(loc).Format("02.01.2006 15:04"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), row.Vehicle)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), row.Plate)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), opType)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), row.Liters)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIndex), row.Driver)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
