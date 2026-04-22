package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"millog_backend/internal/models"
	"millog_backend/internal/repositories"
)

type GPSTrackingService struct {
	gpsRepo *repositories.GPSRepository
	db      repositories.DBExecutor
}

func NewGPSTrackingService(db repositories.DBExecutor) *GPSTrackingService {
	return &GPSTrackingService{
		gpsRepo: &repositories.GPSRepository{},
		db:      db,
	}
}

// RecordVehicleLocation saves a GPS update from a vehicle
func (s *GPSTrackingService) RecordVehicleLocation(ctx context.Context, location *models.GPSLocation) error {
	if err := s.gpsRepo.SaveGPSLocation(ctx, s.db, location); err != nil {
		return fmt.Errorf("failed to save GPS location: %w", err)
	}

	// Check for geofence breaches
	if err := s.checkGeofenceBreach(ctx, location); err != nil {
		// Log error but don't fail the request
		fmt.Printf("geofence check error: %v\n", err)
	}

	return nil
}

// GetVehicleLocation returns current position of a vehicle
func (s *GPSTrackingService) GetVehicleLocation(ctx context.Context, vehicleID string) (*models.GPSLocation, error) {
	location, err := s.gpsRepo.GetLatestLocation(ctx, s.db, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle location: %w", err)
	}
	return location, nil
}

// GetFleetLocations returns current positions of all vehicles in a unit (real-time map)
func (s *GPSTrackingService) GetFleetLocations(ctx context.Context, unitID int64) ([]models.GPSLocation, error) {
	locations, err := s.gpsRepo.GetFleetLocations(ctx, s.db, unitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fleet locations: %w", err)
	}
	return locations, nil
}

// GetVehiclePlates returns a map vehicleID(UUID) -> plate_number for all active vehicles.
// Використовується у fleet-map, щоб показати читабельну назву.
func (s *GPSTrackingService) GetVehiclePlates(ctx context.Context) (map[string]string, error) {
	result := make(map[string]string)
	rows, err := s.db.Query(ctx, `SELECT id, plate_number FROM vehicles`)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, plate string
		if err := rows.Scan(&id, &plate); err != nil {
			continue
		}
		result[id] = plate
	}
	return result, nil
}

// GetVehicleTrajectory returns the path a vehicle traveled in a time window
func (s *GPSTrackingService) GetVehicleTrajectory(ctx context.Context, vehicleID string, startTime, endTime time.Time) ([]models.GPSLocation, error) {
	locations, err := s.gpsRepo.GetVehicleLocationsHistory(ctx, s.db, vehicleID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get trajectory: %w", err)
	}
	return locations, nil
}

// CalculateDistance calculates the distance traveled based on GPS points (Haversine formula)
func (s *GPSTrackingService) CalculateDistance(locations []models.GPSLocation) float64 {
	if len(locations) < 2 {
		return 0
	}

	const earthRadiusKm = 6371.0
	totalDistance := 0.0

	for i := 0; i < len(locations)-1; i++ {
		lat1 := locations[i].Latitude * math.Pi / 180
		lon1 := locations[i].Longitude * math.Pi / 180
		lat2 := locations[i+1].Latitude * math.Pi / 180
		lon2 := locations[i+1].Longitude * math.Pi / 180

		dlat := lat2 - lat1
		dlon := lon2 - lon1

		a := math.Sin(dlat/2)*math.Sin(dlat/2) +
			math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
		c := 2 * math.Asin(math.Sqrt(a))
		distance := earthRadiusKm * c

		totalDistance += distance
	}

	return totalDistance
}

// CreateGeofence creates a boundary alert zone
func (s *GPSTrackingService) CreateGeofence(ctx context.Context, geofence *models.Geofence) error {
	if err := s.gpsRepo.CreateGeofence(ctx, s.db, geofence); err != nil {
		return fmt.Errorf("failed to create geofence: %w", err)
	}
	return nil
}

// GetGeofences returns all geofences for unit
func (s *GPSTrackingService) GetGeofences(ctx context.Context, unitID int64) ([]models.Geofence, error) {
	geofences, err := s.gpsRepo.GetGeofences(ctx, s.db, unitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get geofences: %w", err)
	}
	return geofences, nil
}

// GetGeofenceAlerts returns recent geofence breach events
func (s *GPSTrackingService) GetGeofenceAlerts(ctx context.Context, unitID int64, hours int) ([]models.GeofenceAlert, error) {
	alerts, err := s.gpsRepo.GetGeofenceAlerts(ctx, s.db, unitID, hours)
	if err != nil {
		return nil, fmt.Errorf("failed to get geofence alerts: %w", err)
	}
	return alerts, nil
}

// checkGeofenceBreach checks if a vehicle location violates any geofence and records alerts
func (s *GPSTrackingService) checkGeofenceBreach(ctx context.Context, location *models.GPSLocation) error {
	geofences, err := s.gpsRepo.GetGeofences(ctx, s.db, location.UnitID)
	if err != nil {
		return err
	}

	for _, gf := range geofences {
		if !gf.Active {
			continue
		}

		// Calculate distance from vehicle to geofence center
		distance := s.calculateHaversine(
			location.Latitude, location.Longitude,
			gf.Latitude, gf.Longitude,
		)

		// Check if vehicle is within geofence radius
		isInside := distance <= gf.Radius/1000 // Convert meters to km

		// Only log ENTER events (we could enhance this to track state changes)
		if isInside {
			alert := &models.GeofenceAlert{
				VehicleID:  location.VehicleID,
				GeofenceID: gf.ID,
				EventType:  "ENTER",
				Latitude:   location.Latitude,
				Longitude:  location.Longitude,
				Timestamp:  location.Timestamp,
			}

			if err := s.gpsRepo.RecordGeofenceAlert(ctx, s.db, alert); err != nil {
				fmt.Printf("failed to record geofence alert: %v\n", err)
			}
		}
	}

	return nil
}

// calculateHaversine calculates distance in km between two GPS coordinates
func (s *GPSTrackingService) calculateHaversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Asin(math.Sqrt(a))

	return earthRadiusKm * c
}

// GetDetailedFleetStatus returns comprehensive fleet information with location, speed, fuel
func (s *GPSTrackingService) GetDetailedFleetStatus(ctx context.Context, unitID int64) (map[string]interface{}, error) {
	locations, err := s.GetFleetLocations(ctx, unitID)
	if err != nil {
		return nil, err
	}

	status := map[string]interface{}{
		"timestamp": time.Now(),
		"vehicles":  make([]map[string]interface{}, 0),
		"count":     len(locations),
	}

	vehicles := status["vehicles"].([]map[string]interface{})
	for _, loc := range locations {
		vehicle := map[string]interface{}{
			"vehicle_id": loc.VehicleID,
			"latitude":   loc.Latitude,
			"longitude":  loc.Longitude,
			"speed":      loc.Speed,
			"heading":    loc.Heading,
			"timestamp":  loc.Timestamp,
			"updated":    time.Since(loc.Timestamp).Seconds(), // seconds since last update
		}
		vehicles = append(vehicles, vehicle)
	}
	status["vehicles"] = vehicles

	return status, nil
}
