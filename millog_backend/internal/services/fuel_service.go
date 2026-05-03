package services

import (
	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/repositories"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FuelService struct {
	fuelRepo *repositories.FuelRepository
	Pool     *pgxpool.Pool
}

func NewFuelService(fuelRepo *repositories.FuelRepository, pool *pgxpool.Pool) *FuelService {
	return &FuelService{fuelRepo: fuelRepo, Pool: pool}
}

func (s *FuelService) AddFuelRecord(ctx context.Context, record *models.FuelRecord) error {
	if record.RecordType == models.FuelExpense {
		var currentBalance float64

		query := `
			SELECT COALESCE(
				SUM(CASE WHEN record_type = 'REFUEL' THEN liters ELSE -liters END), 0
			) FROM fuel_records 
			WHERE vehicle_id = $1
		`
		err := s.Pool.QueryRow(ctx, query, record.VehicleID).Scan(&currentBalance)
		if err != nil {
			return fmt.Errorf("помилка перевірки балансу пального: %w", err)
		}

		if currentBalance < record.Liters {
			return fmt.Errorf("недостатньо пального. У баку: %.1f л, спроба списати: %.1f л", currentBalance, record.Liters)
		}
	}

	return s.fuelRepo.CreateFuelRecord(ctx, record, s.Pool)
}

func (s *FuelService) GetVehicleFuelHistory(ctx context.Context, vehicleID string) ([]*models.FuelRecord, error) {
	return s.fuelRepo.GetRecordsByVehicleID(ctx, vehicleID, s.Pool)
}
