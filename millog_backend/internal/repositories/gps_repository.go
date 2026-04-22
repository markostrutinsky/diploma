package repositories

import (
	"context"
	"fmt"
	"time"

	"millog_backend/internal/models"
)

type GPSRepository struct{}

// SaveGPSLocation stores a GPS location from a vehicle
func (r *GPSRepository) SaveGPSLocation(ctx context.Context, db DBExecutor, location *models.GPSLocation) error {
	query := `
		INSERT INTO gps_locations (vehicle_id, unit_id, latitude, longitude, altitude, speed, heading, accuracy, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	return db.QueryRow(ctx, query,
		location.VehicleID, location.UnitID, location.Latitude, location.Longitude,
		location.Altitude, location.Speed, location.Heading, location.Accuracy, location.Timestamp,
	).Scan(&location.ID, &location.CreatedAt)
}

// GetLatestLocation returns the most recent location of a vehicle
func (r *GPSRepository) GetLatestLocation(ctx context.Context, db DBExecutor, vehicleID string) (*models.GPSLocation, error) {
	query := `
		SELECT id, vehicle_id, unit_id, latitude, longitude,
			COALESCE(altitude, 0), COALESCE(speed, 0), COALESCE(heading, 0), COALESCE(accuracy, 0),
			timestamp, COALESCE(created_at, NOW())
		FROM gps_locations
		WHERE vehicle_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var location models.GPSLocation
	err := db.QueryRow(ctx, query, vehicleID).Scan(
		&location.ID, &location.VehicleID, &location.UnitID,
		&location.Latitude, &location.Longitude, &location.Altitude,
		&location.Speed, &location.Heading, &location.Accuracy,
		&location.Timestamp, &location.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest location: %w", err)
	}

	return &location, nil
}

// GetVehicleLocationsHistory returns GPS history for a time range
func (r *GPSRepository) GetVehicleLocationsHistory(ctx context.Context, db DBExecutor, vehicleID string, startTime, endTime time.Time) ([]models.GPSLocation, error) {
	query := `
		SELECT id, vehicle_id, unit_id, latitude, longitude,
			COALESCE(altitude, 0), COALESCE(speed, 0), COALESCE(heading, 0), COALESCE(accuracy, 0),
			timestamp, COALESCE(created_at, NOW())
		FROM gps_locations
		WHERE vehicle_id = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp ASC
	`

	rows, err := db.Query(ctx, query, vehicleID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query locations: %w", err)
	}
	defer rows.Close()

	var locations []models.GPSLocation
	for rows.Next() {
		var loc models.GPSLocation
		if err := rows.Scan(
			&loc.ID, &loc.VehicleID, &loc.UnitID,
			&loc.Latitude, &loc.Longitude, &loc.Altitude,
			&loc.Speed, &loc.Heading, &loc.Accuracy,
			&loc.Timestamp, &loc.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, loc)
	}

	return locations, nil
}

// GetFleetLocations returns current locations for all vehicles in a unit.
// Якщо unitID <= 0 — повертає останні точки по ВСІХ машинах (режим адміна).
func (r *GPSRepository) GetFleetLocations(ctx context.Context, db DBExecutor, unitID int64) ([]models.GPSLocation, error) {
	var (
		query string
		args  []interface{}
	)
	// 🔥 Показуємо на мапі лише ті машини, в яких є активний рейс (status = 'DISPATCHED').
	// Як тільки рейс закривається статусом 'DELIVERED' або 'CANCELLED' —
	// машина зникає з fleet-map, навіть якщо симулятор продовжує слати GPS-пінги.
	if unitID > 0 {
		query = `
			WITH latest_locations AS (
				SELECT DISTINCT ON (vehicle_id)
					id, vehicle_id, unit_id, latitude, longitude, altitude, speed, heading, accuracy, timestamp, created_at
				FROM gps_locations
				WHERE unit_id = $1
				  AND timestamp > NOW() - INTERVAL '5 minutes'
				ORDER BY vehicle_id, timestamp DESC
			)
			SELECT ll.id, ll.vehicle_id, ll.unit_id, ll.latitude, ll.longitude,
				COALESCE(ll.altitude, 0), COALESCE(ll.speed, 0), COALESCE(ll.heading, 0), COALESCE(ll.accuracy, 0),
				ll.timestamp, COALESCE(ll.created_at, NOW())
			FROM latest_locations ll
			WHERE EXISTS (
				SELECT 1 FROM shipments s
				WHERE s.vehicle_id::text = ll.vehicle_id
				  AND s.status = 'DISPATCHED'
			)
		`
		args = []interface{}{unitID}
	} else {
		query = `
			WITH latest_locations AS (
				SELECT DISTINCT ON (vehicle_id)
					id, vehicle_id, unit_id, latitude, longitude, altitude, speed, heading, accuracy, timestamp, created_at
				FROM gps_locations
				WHERE timestamp > NOW() - INTERVAL '5 minutes'
				ORDER BY vehicle_id, timestamp DESC
			)
			SELECT ll.id, ll.vehicle_id, ll.unit_id, ll.latitude, ll.longitude,
				COALESCE(ll.altitude, 0), COALESCE(ll.speed, 0), COALESCE(ll.heading, 0), COALESCE(ll.accuracy, 0),
				ll.timestamp, COALESCE(ll.created_at, NOW())
			FROM latest_locations ll
			WHERE EXISTS (
				SELECT 1 FROM shipments s
				WHERE s.vehicle_id::text = ll.vehicle_id
				  AND s.status = 'DISPATCHED'
			)
		`
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query fleet locations: %w", err)
	}
	defer rows.Close()

	var locations []models.GPSLocation
	for rows.Next() {
		var loc models.GPSLocation
		if err := rows.Scan(
			&loc.ID, &loc.VehicleID, &loc.UnitID,
			&loc.Latitude, &loc.Longitude, &loc.Altitude,
			&loc.Speed, &loc.Heading, &loc.Accuracy,
			&loc.Timestamp, &loc.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, loc)
	}

	return locations, nil
}

// CreateGeofence creates a new geofence boundary
func (r *GPSRepository) CreateGeofence(ctx context.Context, db DBExecutor, geofence *models.Geofence) error {
	query := `
		INSERT INTO geofences (unit_id, name, latitude, longitude, radius, type, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	return db.QueryRow(ctx, query,
		geofence.UnitID, geofence.Name, geofence.Latitude, geofence.Longitude,
		geofence.Radius, geofence.Type, geofence.Active,
	).Scan(&geofence.ID, &geofence.CreatedAt, &geofence.UpdatedAt)
}

// GetGeofences returns all geofences for a unit
func (r *GPSRepository) GetGeofences(ctx context.Context, db DBExecutor, unitID int64) ([]models.Geofence, error) {
	query := `
		SELECT id, unit_id, name, latitude, longitude, radius, type, active, created_at, updated_at
		FROM geofences
		WHERE unit_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, query, unitID)
	if err != nil {
		return nil, fmt.Errorf("failed to query geofences: %w", err)
	}
	defer rows.Close()

	var geofences []models.Geofence
	for rows.Next() {
		var g models.Geofence
		if err := rows.Scan(
			&g.ID, &g.UnitID, &g.Name, &g.Latitude, &g.Longitude,
			&g.Radius, &g.Type, &g.Active, &g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan geofence: %w", err)
		}
		geofences = append(geofences, g)
	}

	return geofences, nil
}

// GetGeofenceAlerts returns recent geofence breach events
func (r *GPSRepository) GetGeofenceAlerts(ctx context.Context, db DBExecutor, unitID int64, hours int) ([]models.GeofenceAlert, error) {
	query := `
		SELECT 
			ga.id, ga.vehicle_id, ga.geofence_id, ga.event_type, ga.latitude, ga.longitude, ga.timestamp, ga.created_at
		FROM geofence_alerts ga
		JOIN geofences g ON ga.geofence_id = g.id
		WHERE g.unit_id = $1 AND ga.created_at > NOW() - INTERVAL '1 hour' * $2
		ORDER BY ga.created_at DESC
	`

	rows, err := db.Query(ctx, query, unitID, hours)
	if err != nil {
		return nil, fmt.Errorf("failed to query geofence alerts: %w", err)
	}
	defer rows.Close()

	var alerts []models.GeofenceAlert
	for rows.Next() {
		var alert models.GeofenceAlert
		if err := rows.Scan(
			&alert.ID, &alert.VehicleID, &alert.GeofenceID, &alert.EventType,
			&alert.Latitude, &alert.Longitude, &alert.Timestamp, &alert.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// RecordGeofenceAlert logs a geofence breach event
func (r *GPSRepository) RecordGeofenceAlert(ctx context.Context, db DBExecutor, alert *models.GeofenceAlert) error {
	query := `
		INSERT INTO geofence_alerts (vehicle_id, geofence_id, event_type, latitude, longitude, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	return db.QueryRow(ctx, query,
		alert.VehicleID, alert.GeofenceID, alert.EventType,
		alert.Latitude, alert.Longitude, alert.Timestamp,
	).Scan(&alert.ID, &alert.CreatedAt)
}
