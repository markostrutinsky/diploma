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

		`ALTER TABLE volunteer_requests ADD COLUMN IF NOT EXISTS unit_id BIGINT REFERENCES units(id) ON DELETE SET NULL`,
		`ALTER TABLE resources ADD COLUMN IF NOT EXISTS unit_type VARCHAR(50) DEFAULT 'PCS';`,

		// 12. Оновлення таблиці vehicles для існуючих БД
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS tank_capacity DECIMAL(10,2) NOT NULL DEFAULT 0;`,
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS fuel_norm DECIMAL(10,2) NOT NULL DEFAULT 0;`,

		// 13. Детектор аномалій пального
		`ALTER TABLE fuel_records ADD COLUMN IF NOT EXISTS is_anomaly BOOLEAN NOT NULL DEFAULT FALSE;`,
		`ALTER TABLE fuel_records ADD COLUMN IF NOT EXISTS anomaly_reason TEXT;`,

		// 14. Поля для ТО та ремонту
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS maintenance_interval_km INT NOT NULL DEFAULT 10000;`,
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS last_maintenance_odometer INT NOT NULL DEFAULT 0;`,

		// 15. Історія технічного обслуговування (ТО)
		`CREATE TABLE IF NOT EXISTS maintenance_records (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
            odometer_km INT NOT NULL,
            description TEXT NOT NULL,
            performed_by VARCHAR(255),
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE INDEX IF NOT EXISTS idx_maintenance_records_vehicle_id ON maintenance_records(vehicle_id)`,

		// 16. Вартість ремонту
		`ALTER TABLE maintenance_records ADD COLUMN IF NOT EXISTS cost_amount DECIMAL(10,2) NOT NULL DEFAULT 0;`,

		// 17. Причина зміни статусу (для ремонту/списання)
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS status_reason TEXT;`,

		// 18. Посилання на файл акту виконаних робіт
		`ALTER TABLE maintenance_records ADD COLUMN IF NOT EXISTS document_url TEXT;`,

		// 19. Історія водіїв
		`CREATE TABLE IF NOT EXISTS vehicle_driver_history (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            vehicle_id UUID REFERENCES vehicles(id) ON DELETE CASCADE,
            driver_id UUID REFERENCES users(id) ON DELETE SET NULL,
            assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );`,
		`CREATE INDEX IF NOT EXISTS idx_vehicle_driver_history_vehicle_id ON vehicle_driver_history(vehicle_id);`,
		`CREATE INDEX IF NOT EXISTS idx_vehicle_driver_history_driver_id ON vehicle_driver_history(driver_id);`,
		`CREATE INDEX IF NOT EXISTS idx_vehicle_driver_history_assigned_at ON vehicle_driver_history(assigned_at);`,

		`ALTER TABLE maintenance_records ADD COLUMN IF NOT EXISTS driver_id UUID REFERENCES users(id) ON DELETE SET NULL;`,

		// ==========================================
		// НОВИЙ ФУНДАМЕНТ ДЛЯ ЛОГІСТИКИ (СКЛАДИ)
		// ==========================================

		// 20. Створюємо таблицю складів / локацій
		`CREATE TABLE IF NOT EXISTS warehouses (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			unit_id BIGINT REFERENCES units(id) ON DELETE CASCADE,
            name VARCHAR(255) NOT NULL,
            location_type VARCHAR(50) DEFAULT 'STATIONARY',
            latitude DOUBLE PRECISION,
            longitude DOUBLE PRECISION,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );`,

		// 21. Додаємо зв'язок ресурсів зі складами
		`ALTER TABLE resources ADD COLUMN IF NOT EXISTS warehouse_id UUID REFERENCES warehouses(id) ON DELETE SET NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_resources_warehouse_id ON resources(warehouse_id);`,

		// 22. Видаляємо старі координати з ресурсів (якщо вони там були)
		`ALTER TABLE resources DROP COLUMN IF EXISTS latitude;`,
		`ALTER TABLE resources DROP COLUMN IF EXISTS longitude;`,

		// 23. Додаємо поле для відповідального користувача за ресурс
		`ALTER TABLE resources ADD COLUMN IF NOT EXISTS assigned_to_user_id UUID REFERENCES users(id) ON DELETE SET NULL;`,

		// 24. Таблиця для відстеження видачі ресурсів
		`CREATE TABLE IF NOT EXISTS resource_assignments (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            quantity INT NOT NULL DEFAULT 1,
            status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE', -- Статуси: ACTIVE (на руках), RETURNED (здано), LOST (втрачено), WRITTEN_OFF (списано)
            issued_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            returned_at TIMESTAMP,
            notes TEXT
        );`,

		`CREATE INDEX IF NOT EXISTS idx_resource_assignments_user ON resource_assignments(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_resource_assignments_resource ON resource_assignments(resource_id);`,

		`ALTER TABLE resources DROP COLUMN IF EXISTS assigned_to_user_name;`,
		`ALTER TABLE resources DROP COLUMN IF EXISTS assigned_to_user_id;`,

		// ==========================================
		// ВАНТАЖОПЕРЕВЕЗЕННЯ ТА КОНТРОЛЬ ВАГИ
		// ==========================================

		// 25. Тип та вантажопідйомність авто
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS type VARCHAR(50) DEFAULT 'VAN';`,
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS capacity_kg DECIMAL(10,2) NOT NULL DEFAULT 1500.00;`,

		// 26. Вага ресурсу/товару (САМЕ ЦЬОГО НЕ ВИСТАЧАЛО!)
		`ALTER TABLE resources ADD COLUMN IF NOT EXISTS weight_kg DECIMAL(10,2) NOT NULL DEFAULT 1.00;`,

		// 27. Рейси (Накладні / Відправки)
		`CREATE TABLE IF NOT EXISTS shipments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			from_warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
			to_warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
			vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE RESTRICT,
			priority VARCHAR(20) NOT NULL DEFAULT 'NORMAL',
			status VARCHAR(30) NOT NULL DEFAULT 'DISPATCHED', -- DISPATCHED, DELIVERED, CANCELLED
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// 28. Вміст рейсу (що саме везуть)
		`CREATE TABLE IF NOT EXISTS shipment_items (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			shipment_id UUID NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
			resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE RESTRICT,
			quantity INT NOT NULL CHECK (quantity > 0)
		);`,

		`CREATE INDEX IF NOT EXISTS idx_shipments_from ON shipments(from_warehouse_id);`,
		`CREATE INDEX IF NOT EXISTS idx_shipments_to ON shipments(to_warehouse_id);`,
		`CREATE INDEX IF NOT EXISTS idx_shipment_items_shipment ON shipment_items(shipment_id);`,

		`CREATE TABLE IF NOT EXISTS shipments (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            from_warehouse_id UUID REFERENCES warehouses(id) ON DELETE CASCADE,
            to_warehouse_id UUID REFERENCES warehouses(id) ON DELETE CASCADE,
            vehicle_id UUID REFERENCES vehicles(id) ON DELETE SET NULL,
            priority VARCHAR(50) DEFAULT 'NORMAL', -- NORMAL або URGENT
            status VARCHAR(50) DEFAULT 'DISPATCHED', -- DISPATCHED або DELIVERED
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );`,

		`CREATE TABLE IF NOT EXISTS shipment_items (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            shipment_id UUID REFERENCES shipments(id) ON DELETE CASCADE,
            resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
            quantity INT NOT NULL CHECK (quantity > 0)
        );`,

		`CREATE INDEX IF NOT EXISTS idx_shipments_from ON shipments(from_warehouse_id);`,
		`CREATE INDEX IF NOT EXISTS idx_shipments_to ON shipments(to_warehouse_id);`,

		`ALTER TABLE shipment_items ADD COLUMN IF NOT EXISTS request_id UUID REFERENCES supply_requests(id) ON DELETE SET NULL;`,
		`ALTER TABLE supply_requests ADD COLUMN IF NOT EXISTS target_warehouse_id UUID REFERENCES warehouses(id) ON DELETE SET NULL;`,

		`ALTER TABLE shipments ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;`,

		`CREATE TABLE IF NOT EXISTS audit_logs (
            id SERIAL PRIMARY KEY,
            user_id UUID REFERENCES users(id) ON DELETE SET NULL,
            action_type VARCHAR(50) NOT NULL,  -- 'DELETE', 'WRITE_OFF', 'UPDATE_ROLE'
            entity_type VARCHAR(50) NOT NULL,  -- 'RESOURCE', 'USER', 'WAREHOUSE'
            entity_id VARCHAR(50),             -- ID того, що змінили
            details TEXT,                      -- Детальний опис ("Списано 5 шт бронежилетів")
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );`,
	}

	for i, m := range migrations {
		if _, err := pool.Exec(ctx, m); err != nil {
			return fmt.Errorf("migration %d: %w", i+1, err)
		}
	}
	return nil
}
