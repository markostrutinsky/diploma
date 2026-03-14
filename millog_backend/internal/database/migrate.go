package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Migrate виконує міграції. Працює з порожньою БД (init.sql не спрацював)
// або доповнює існуючу схему.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	migrations := []string{
		// 1. Користувачі (якщо init.sql не виконався)
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username VARCHAR(100) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			full_name VARCHAR(255),
			phone VARCHAR(50),
			password_hash TEXT,
			role VARCHAR(30) NOT NULL DEFAULT 'VOLUNTEER',
			status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
			unit_id BIGINT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`,
		`CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)`,

		// 2. Invite-токени
		`CREATE TABLE IF NOT EXISTS invite_tokens (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
			token_hash VARCHAR(255) NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			used_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_invite_tokens_hash ON invite_tokens(token_hash)`,

		// 3. Підрозділи
		`CREATE TABLE IF NOT EXISTS units (
			id BIGSERIAL PRIMARY KEY,
			parent_id BIGINT REFERENCES units(id) ON DELETE SET NULL,
			name VARCHAR(255) NOT NULL,
			unit_type VARCHAR(20) NOT NULL
		)`,

		// 4. Категорії ресурсів
		`CREATE TABLE IF NOT EXISTS resource_categories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) NOT NULL UNIQUE,
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,

		// 5. Ресурси (базова таблиця без unit_id для старих БД)
		`CREATE TABLE IF NOT EXISTS resources (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			category_id UUID NOT NULL REFERENCES resource_categories(id) ON DELETE RESTRICT,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			quantity INTEGER NOT NULL DEFAULT 0,
			serial_number VARCHAR(100),
			location VARCHAR(255),
			condition VARCHAR(20) NOT NULL DEFAULT 'NEW',
			min_quantity INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_resources_category ON resources(category_id)`,

		// 6. Додати unit_id до resources (якщо немає)
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='resources' AND column_name='unit_id') THEN
				ALTER TABLE resources ADD COLUMN unit_id BIGINT REFERENCES units(id) ON DELETE SET NULL;
			END IF;
		END $$`,
		`CREATE INDEX IF NOT EXISTS idx_resources_unit ON resources(unit_id)`,

		// 7. Заявки на постачання
		`CREATE TABLE IF NOT EXISTS supply_requests (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
			quantity INTEGER NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
			approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
			approved_at TIMESTAMP WITH TIME ZONE,
			comment TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_supply_requests_status ON supply_requests(status)`,

		// 8. Транспорт
		`CREATE TABLE IF NOT EXISTS vehicles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			brand VARCHAR(100) NOT NULL,
			model VARCHAR(100),
			plate_number VARCHAR(20) NOT NULL UNIQUE,
			status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
			driver_id UUID REFERENCES users(id) ON DELETE SET NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS fuel_records (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
			liters DECIMAL(10,2) NOT NULL,
			odometer_km INTEGER,
			record_type VARCHAR(20) NOT NULL,
			created_by UUID REFERENCES users(id) ON DELETE SET NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_fuel_records_vehicle ON fuel_records(vehicle_id)`,

		// 9. Заявки для волонтерів
		`CREATE TABLE IF NOT EXISTS volunteer_requests (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
			taken_by UUID REFERENCES users(id) ON DELETE SET NULL,
			taken_at TIMESTAMP WITH TIME ZONE,
			completed_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_volunteer_requests_status ON volunteer_requests(status)`,

		// 10. Refresh-токени
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token_hash VARCHAR(255) NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			revoked_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens(token_hash)`,

		// 11. Розширити role для нових ролей
		`ALTER TABLE users ALTER COLUMN role TYPE VARCHAR(30)`,
	}

	for i, m := range migrations {
		if _, err := pool.Exec(ctx, m); err != nil {
			return fmt.Errorf("migration %d: %w", i+1, err)
		}
	}
	return nil
}
