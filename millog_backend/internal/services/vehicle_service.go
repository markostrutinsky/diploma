package services

import (
	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VehicleService struct {
	repo *repositories.VehicleRepository
	Pool *pgxpool.Pool
}

func NewVehicleService(repo *repositories.VehicleRepository, pool *pgxpool.Pool) *VehicleService {
	return &VehicleService{repo: repo, Pool: pool}
}

func (s *VehicleService) CreateVehicle(ctx context.Context, v *models.Vehicle) error {
	v.PlateNumber = strings.ToUpper(strings.TrimSpace(v.PlateNumber))
	v.Brand = strings.TrimSpace(v.Brand)

	// ВИПРАВЛЕНО: Просто працюємо з рядком (string), без кастомних типів
	v.Type = strings.ToUpper(strings.TrimSpace(v.Type))

	if v.Brand == "" {
		return errors.New("марка автомобіля є обов'язковою")
	}
	if v.PlateNumber == "" {
		return errors.New("номерний знак є обов'язковим")
	}
	if v.Type == "" {
		return errors.New("тип автомобіля (PICKUP, VAN, TRUCK) є обов'язковим")
	}
	if v.CapacityKg <= 0 {
		return errors.New("вантажопідйомність повинна бути більшою за нуль")
	}
	if v.TankCapacity <= 0 {
		return errors.New("об'єм бака повинен бути більшим за нуль")
	}

	if v.Status == "" {
		v.Status = models.VehicleActive
	}

	err := s.repo.Create(ctx, s.Pool, v)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "unique constraint") {
			return errors.New("автомобіль з таким номерним знаком вже існує у базі")
		}
		return err
	}

	return nil
}

func (s *VehicleService) GetAllVehicles(ctx context.Context) ([]models.Vehicle, error) {
	return s.repo.GetAll(ctx, s.Pool)
}

func (s *VehicleService) GetVehicleByID(ctx context.Context, id string) (*models.Vehicle, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("ID автомобіля не вказано")
	}
	return s.repo.GetByID(ctx, id, s.Pool)
}

func (s *VehicleService) UpdateStatus(ctx context.Context, vehicleID string, req *models.VehicleStatusUpdate) error {
	return s.repo.UpdateStatus(ctx, vehicleID, req.Status, req.Reason, s.Pool)
}

func (s *VehicleService) PerformMaintenance(ctx context.Context, record *models.MaintenanceRecord) error {
	return s.repo.PerformMaintenance(ctx, record, s.Pool)
}

func (s *VehicleService) ScheduleMaintenance(ctx context.Context, record *models.MaintenanceRecord) error {
	return s.repo.ScheduleMaintenance(ctx, record, s.Pool)
}

func (s *VehicleService) GetMaintenanceHistory(ctx context.Context, vehicleID string) ([]*models.MaintenanceRecord, error) {
	return s.repo.GetMaintenanceHistory(ctx, vehicleID, s.Pool)
}

func (s *VehicleService) AssignDriver(ctx context.Context, vehicleID string, driverID *string) error {
	return s.repo.AssignDriver(ctx, vehicleID, driverID, s.Pool)
}

func (s *VehicleService) GetDriverHistory(ctx context.Context, vehicleID string) ([]models.VehicleDriverHistory, error) {
	return s.repo.GetDriverHistory(ctx, vehicleID, s.Pool)
}

func (s *VehicleService) UpdateVehicle(ctx context.Context, id, brand, model, plateNumber string, capacityKg float64) error {
	return s.repo.Update(ctx, s.Pool, id, brand, model, plateNumber, capacityKg)
}

func (s *VehicleService) DeleteVehicle(ctx context.Context, id string) error {
	// 1. Отримуємо поточний стан авто (наприклад, щоб перевірити чи не в рейсі)
	// Припускаю, що у тебе є метод GetByID
	vehicle, err := s.repo.GetByID(ctx, id, s.Pool)
	if err != nil {
		return errors.New("автомобіль не знайдено")
	}

	// 2. Блокуємо видалення, якщо машина зараз виконує рейс
	if vehicle.Status == "ON_MISSION" {
		return errors.New("неможливо списати автомобіль, який зараз перебуває у рейсі")
	}

	// 3. Видаляємо (якщо є історія заправок, база даних може видати помилку зовнішнього ключа,
	// тоді можна змінити статус на WRITTEN_OFF замість фізичного DELETE)
	return s.repo.Delete(ctx, s.Pool, id)
}

// GetAvailableForRoute повертає список вільних машин, доступних для рейсу між двома складами.
// Фільтрує по current_warehouse_id (де зараз фізично), а якщо current = NULL — по home_warehouse_id.
func (s *VehicleService) GetAvailableForRoute(ctx context.Context, fromWarehouseID string, toWarehouseID string) ([]models.Vehicle, error) {
	return s.repo.GetAvailableForRoute(ctx, s.Pool, fromWarehouseID, toWarehouseID)
}
