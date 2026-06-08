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
            role VARCHAR(30) NOT NULL DEFAULT 'CONTRACTOR',
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
		`CREATE TABLE IF NOT EXISTS CONTRACTOR_requests (
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
		`CREATE INDEX IF NOT EXISTS idx_CONTRACTOR_requests_status ON CONTRACTOR_requests(status)`,

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

		`ALTER TABLE CONTRACTOR_requests ADD COLUMN IF NOT EXISTS unit_id BIGINT REFERENCES units(id) ON DELETE SET NULL`,
		`ALTER TABLE contractor_requests ADD COLUMN IF NOT EXISTS deadline TIMESTAMP WITH TIME ZONE;`,
		`ALTER TABLE contractor_requests ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;`,
		`ALTER TABLE resources ADD COLUMN IF NOT EXISTS unit_type VARCHAR(50) DEFAULT 'PCS';`,

		// 12. Оновлення таблиці vehicles для існуючих БД
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS tank_capacity DECIMAL(10,2) NOT NULL DEFAULT 0;`,
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS fuel_norm DECIMAL(10,2) NOT NULL DEFAULT 0;`,

		// 13. Детектор аномалій пального
		`ALTER TABLE fuel_records ADD COLUMN IF NOT EXISTS is_anomaly BOOLEAN NOT NULL DEFAULT FALSE;`,
		`ALTER TABLE fuel_records ADD COLUMN IF NOT EXISTS anomaly_reason TEXT;`,
		// 13a. Обсяг «зайвого» пального запису (перевитрата понад норму / витрата без руху).
		//      Потрібно антифрод-системі для точного розрахунку грошових втрат.
		`ALTER TABLE fuel_records ADD COLUMN IF NOT EXISTS anomaly_excess_liters DECIMAL(10,2) NOT NULL DEFAULT 0;`,

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

		// 25a. Поточна локація машини (на якому складі зараз знаходиться)
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS current_warehouse_id UUID REFERENCES warehouses(id) ON DELETE SET NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_vehicles_warehouse ON vehicles(current_warehouse_id);`,

		// 25b. Базовий (домашній) склад машини — постійна приписка, не змінюється після рейсів
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS home_warehouse_id UUID REFERENCES warehouses(id) ON DELETE SET NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_vehicles_home_warehouse ON vehicles(home_warehouse_id);`,
		// Заповнюємо home_warehouse_id з current_warehouse_id для існуючих записів
		`UPDATE vehicles SET home_warehouse_id = current_warehouse_id WHERE home_warehouse_id IS NULL AND current_warehouse_id IS NOT NULL;`,

		// 26. Вага ресурсу/товару (САМЕ ЦЬОГО НЕ ВИСТАЧАЛО!)
		`ALTER TABLE resources ADD COLUMN IF NOT EXISTS weight_kg DECIMAL(10,2) NOT NULL DEFAULT 1.00;`,

		// 26b. Ціна одиниці ресурсу (для фінансової аналітики)
		`ALTER TABLE resources ADD COLUMN IF NOT EXISTS unit_price NUMERIC(12,2) NOT NULL DEFAULT 0;`,

		// 27. Рейси (Накладні / Відправки)
		`CREATE TABLE IF NOT EXISTS shipments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			from_warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
			to_warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
			vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE RESTRICT,
			priority VARCHAR(20) NOT NULL DEFAULT 'NORMAL',
			status VARCHAR(30) NOT NULL DEFAULT 'PENDING', -- PENDING, IN_TRANSIT, DELIVERED, CANCELLED
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			started_at TIMESTAMP WITH TIME ZONE, -- Коли водій почав рейс
			delivered_at TIMESTAMP WITH TIME ZONE -- Коли доставлено
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

		// Додаємо поля для збереження назви ресурсу (щоб не прив'язуватись до конкретного екземпляру на складі)
		`ALTER TABLE supply_requests ADD COLUMN IF NOT EXISTS resource_name VARCHAR(255);`,
		`ALTER TABLE supply_requests ADD COLUMN IF NOT EXISTS resource_category_id UUID REFERENCES resource_categories(id) ON DELETE SET NULL;`,

		// Робимо resource_id nullable (тепер можна створити заявку без прив'язки до конкретного складу)
		`ALTER TABLE supply_requests ALTER COLUMN resource_id DROP NOT NULL;`,

		// Додаємо нові колонки для відстеження етапів рейсу
		`ALTER TABLE shipments ADD COLUMN IF NOT EXISTS started_at TIMESTAMP WITH TIME ZONE;`,
		`ALTER TABLE shipments ADD COLUMN IF NOT EXISTS delivered_at TIMESTAMP WITH TIME ZONE;`,

		// Додаємо поле для відстеження напрямку руху по ієрархії підрозділів
		// DOWNSTREAM - розподіл (від вищого до нижчого), UPSTREAM - консолідація/повернення (від нижчого до вищого), LATERAL - в межах одного підрозділу
		`ALTER TABLE shipments ADD COLUMN IF NOT EXISTS direction VARCHAR(20);`,
		`CREATE INDEX IF NOT EXISTS idx_shipments_direction ON shipments(direction);`,

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

		// 29. Таблиці для професійної інвентаризації
		`CREATE TABLE IF NOT EXISTS inventory_checks (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
            created_by UUID NOT NULL REFERENCES users(id),
            status VARCHAR(30) NOT NULL DEFAULT 'IN_PROGRESS', -- IN_PROGRESS, COMPLETED, CANCELLED
            started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            completed_at TIMESTAMP WITH TIME ZONE,
            notes TEXT
        );`,

		`CREATE TABLE IF NOT EXISTS inventory_check_items (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            check_id UUID NOT NULL REFERENCES inventory_checks(id) ON DELETE CASCADE,
            resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
            book_quantity INTEGER NOT NULL,  -- Скільки було в базі на момент старту
            actual_quantity INTEGER,         -- Скільки нарахували по факту
            difference INTEGER GENERATED ALWAYS AS (actual_quantity - book_quantity) STORED,
            verified_at TIMESTAMP WITH TIME ZONE
        );`,

		`ALTER TABLE units ADD COLUMN IF NOT EXISTS subscription_tier VARCHAR(20) DEFAULT 'BASIC';`,

		`ALTER TABLE resources ADD COLUMN IF NOT EXISTS barcode VARCHAR(255) DEFAULT '';`,

		// ==========================================
		// GPS TRACKING & GEOFENCING (PRO FEATURE #5)
		// ==========================================

		// Real-time GPS locations from vehicles
		// ВАЖЛИВО: vehicle_id — UUID (vehicles.id теж UUID).
		`CREATE TABLE IF NOT EXISTS gps_locations (
			id BIGSERIAL PRIMARY KEY,
			vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
			unit_id BIGINT NOT NULL REFERENCES units(id) ON DELETE CASCADE,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			altitude DOUBLE PRECISION,
			speed DOUBLE PRECISION,
			heading DOUBLE PRECISION,
			accuracy DOUBLE PRECISION,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		// Переводимо числові колонки з NUMERIC/DECIMAL у DOUBLE PRECISION,
		// щоб pgx міг сканувати їх прямо у float64 без pgtype.Numeric.
		`ALTER TABLE gps_locations
			ALTER COLUMN latitude  TYPE DOUBLE PRECISION USING latitude::double precision,
			ALTER COLUMN longitude TYPE DOUBLE PRECISION USING longitude::double precision,
			ALTER COLUMN altitude  TYPE DOUBLE PRECISION USING altitude::double precision,
			ALTER COLUMN speed     TYPE DOUBLE PRECISION USING speed::double precision,
			ALTER COLUMN heading   TYPE DOUBLE PRECISION USING heading::double precision,
			ALTER COLUMN accuracy  TYPE DOUBLE PRECISION USING accuracy::double precision;`,

		// Fix для існуючих БД, де gps_locations.vehicle_id було створено як BIGINT.
		// Таблиця в реальному житті ще не була заповнена (бачили це на сіді),
		// тому безпечно TRUNCATE + перестворити колонку.
		`DO $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_schema = 'public'
				  AND table_name = 'gps_locations'
				  AND column_name = 'vehicle_id'
				  AND data_type = 'bigint'
			) THEN
				TRUNCATE TABLE gps_locations;
				ALTER TABLE gps_locations DROP COLUMN vehicle_id;
				ALTER TABLE gps_locations
					ADD COLUMN vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE;
			END IF;
		END $$;`,

		`CREATE INDEX IF NOT EXISTS idx_gps_locations_vehicle ON gps_locations(vehicle_id);`,
		`CREATE INDEX IF NOT EXISTS idx_gps_locations_unit ON gps_locations(unit_id);`,
		`CREATE INDEX IF NOT EXISTS idx_gps_locations_timestamp ON gps_locations(timestamp DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_gps_locations_vehicle_time ON gps_locations(vehicle_id, timestamp DESC);`,

		// Geofences (alert zones)
		`CREATE TABLE IF NOT EXISTS geofences (
			id BIGSERIAL PRIMARY KEY,
			unit_id BIGINT NOT NULL REFERENCES units(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			radius DOUBLE PRECISION NOT NULL,
			type VARCHAR(50) NOT NULL,
			active BOOLEAN DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		`ALTER TABLE geofences
			ALTER COLUMN latitude  TYPE DOUBLE PRECISION USING latitude::double precision,
			ALTER COLUMN longitude TYPE DOUBLE PRECISION USING longitude::double precision,
			ALTER COLUMN radius    TYPE DOUBLE PRECISION USING radius::double precision;`,

		`CREATE INDEX IF NOT EXISTS idx_geofences_unit ON geofences(unit_id);`,

		// Geofence breach alerts
		// ВАЖЛИВО: vehicle_id — UUID (vehicles.id теж UUID).
		`CREATE TABLE IF NOT EXISTS geofence_alerts (
			id BIGSERIAL PRIMARY KEY,
			vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
			geofence_id BIGINT NOT NULL REFERENCES geofences(id) ON DELETE CASCADE,
			event_type VARCHAR(20) NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,

		`ALTER TABLE geofence_alerts
			ALTER COLUMN latitude  TYPE DOUBLE PRECISION USING latitude::double precision,
			ALTER COLUMN longitude TYPE DOUBLE PRECISION USING longitude::double precision;`,

		// Fix типу vehicle_id у geofence_alerts для існуючих БД.
		`DO $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_schema = 'public'
				  AND table_name = 'geofence_alerts'
				  AND column_name = 'vehicle_id'
				  AND data_type = 'bigint'
			) THEN
				TRUNCATE TABLE geofence_alerts;
				ALTER TABLE geofence_alerts DROP COLUMN vehicle_id;
				ALTER TABLE geofence_alerts
					ADD COLUMN vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE;
			END IF;
		END $$;`,

		`CREATE INDEX IF NOT EXISTS idx_geofence_alerts_vehicle ON geofence_alerts(vehicle_id);`,
		`CREATE INDEX IF NOT EXISTS idx_geofence_alerts_geofence ON geofence_alerts(geofence_id);`,
		`CREATE INDEX IF NOT EXISTS idx_geofence_alerts_created ON geofence_alerts(created_at DESC);`,

		// ==========================================
		// MULTI-TENANT (tenants + tenant_id скрізь)
		// ==========================================

		// T1. Таблиця tenants
		`CREATE TABLE IF NOT EXISTS tenants (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL UNIQUE,
		slug VARCHAR(100) NOT NULL UNIQUE,
		subscription_tier VARCHAR(30) NOT NULL DEFAULT 'BASIC',
		subscription_expires_at TIMESTAMP WITH TIME ZONE,
		owner_email VARCHAR(255),
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`,
		`CREATE INDEX IF NOT EXISTS idx_tenants_slug ON tenants(slug);`,

		// Міграція FREE → BASIC для існуючих тенантів
		`UPDATE tenants SET subscription_tier = 'BASIC' WHERE subscription_tier = 'FREE';`,
		`UPDATE units SET subscription_tier = 'BASIC' WHERE subscription_tier = 'FREE';`,

		// T2. Створити дефолтний tenant для існуючих даних (якщо даних немає — створиться пустий)
		`INSERT INTO tenants (name, slug, subscription_tier, is_active)
			SELECT 'Default Organization', 'default', 'ENTERPRISE', TRUE
			WHERE NOT EXISTS (SELECT 1 FROM tenants WHERE slug = 'default');`,

		// T3. Додати tenant_id до users
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		// Підрядники (CONTRACTOR) — глобальні учасники marketplace і НЕ належать жодній організації,
		// тому НЕ призначаємо їм дефолтний tenant (інакше RLS обмежить їх однією організацією
		// і вони перестануть бачити крос-tenant дошку відкритих завдань).
		`UPDATE users SET tenant_id = (SELECT id FROM tenants WHERE slug = 'default') WHERE tenant_id IS NULL AND role <> 'CONTRACTOR';`,
		`CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id);`,

		// T4. Перебудувати UNIQUE на users: глобальні → per-tenant
		`ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;`,
		`ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;`,
		`DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'users_tenant_username_unique'
			) THEN
				ALTER TABLE users ADD CONSTRAINT users_tenant_username_unique UNIQUE (tenant_id, username);
			END IF;
		END $$;`,
		`DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'users_tenant_email_unique'
			) THEN
				ALTER TABLE users ADD CONSTRAINT users_tenant_email_unique UNIQUE (tenant_id, email);
			END IF;
		END $$;`,

		// T5. Додати tenant_id до units
		`ALTER TABLE units ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		`UPDATE units SET tenant_id = (SELECT id FROM tenants WHERE slug = 'default') WHERE tenant_id IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_units_tenant ON units(tenant_id);`,

		// T6. Додати tenant_id до resource_categories (+ UNIQUE per tenant)
		`ALTER TABLE resource_categories ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		`UPDATE resource_categories SET tenant_id = (SELECT id FROM tenants WHERE slug = 'default') WHERE tenant_id IS NULL;`,
		`ALTER TABLE resource_categories DROP CONSTRAINT IF EXISTS resource_categories_name_key;`,
		`DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'resource_categories_tenant_name_unique'
			) THEN
				ALTER TABLE resource_categories ADD CONSTRAINT resource_categories_tenant_name_unique UNIQUE (tenant_id, name);
			END IF;
		END $$;`,
		`CREATE INDEX IF NOT EXISTS idx_resource_categories_tenant ON resource_categories(tenant_id);`,

		// T7. Додати tenant_id до resources
		`ALTER TABLE resources ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		`UPDATE resources SET tenant_id = (SELECT id FROM tenants WHERE slug = 'default') WHERE tenant_id IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_resources_tenant ON resources(tenant_id);`,

		// T8. Додати tenant_id до supply_requests
		`ALTER TABLE supply_requests ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		`UPDATE supply_requests SET tenant_id = (SELECT id FROM tenants WHERE slug = 'default') WHERE tenant_id IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_supply_requests_tenant ON supply_requests(tenant_id);`,

		// T9. Додати tenant_id до vehicles (+ UNIQUE per tenant)
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		`UPDATE vehicles SET tenant_id = (SELECT id FROM tenants WHERE slug = 'default') WHERE tenant_id IS NULL;`,
		`ALTER TABLE vehicles DROP CONSTRAINT IF EXISTS vehicles_plate_number_key;`,
		`DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'vehicles_tenant_plate_unique'
			) THEN
				ALTER TABLE vehicles ADD CONSTRAINT vehicles_tenant_plate_unique UNIQUE (tenant_id, plate_number);
			END IF;
		END $$;`,
		`CREATE INDEX IF NOT EXISTS idx_vehicles_tenant ON vehicles(tenant_id);`,

		// T10. Додати tenant_id до fuel_records
		`ALTER TABLE fuel_records ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		`UPDATE fuel_records SET tenant_id = (SELECT id FROM tenants WHERE slug = 'default') WHERE tenant_id IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_fuel_records_tenant ON fuel_records(tenant_id);`,

		// T11. Додати tenant_id до contractor_requests
		`ALTER TABLE contractor_requests ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		`UPDATE contractor_requests SET tenant_id = (SELECT id FROM tenants WHERE slug = 'default') WHERE tenant_id IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_contractor_requests_tenant ON contractor_requests(tenant_id);`,

		// T11b. Додати target_warehouse_id до contractor_requests (склад призначення для підрядника)
		`ALTER TABLE contractor_requests ADD COLUMN IF NOT EXISTS target_warehouse_id UUID REFERENCES warehouses(id) ON DELETE SET NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_contractor_requests_warehouse ON contractor_requests(target_warehouse_id);`,

		// T11c. Підрядницькі членства (marketplace + per-tenant вето).
		// Підрядник реєструється глобально (tenant_id IS NULL) і бачить крос-tenant дошку
		// відкритих завдань. Але щоб ВЗЯТИ завдання конкретної організації, він має бути
		// схвалений цією організацією. Один підрядник може співпрацювати з кількома tenant-ами.
		`CREATE TABLE IF NOT EXISTS contractor_memberships (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            contractor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
            status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
            note TEXT,
            requested_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            decided_at TIMESTAMP WITH TIME ZONE,
            decided_by UUID REFERENCES users(id) ON DELETE SET NULL,
            UNIQUE (contractor_id, tenant_id)
        )`,
		`CREATE INDEX IF NOT EXISTS idx_contractor_memberships_tenant ON contractor_memberships(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_contractor_memberships_contractor ON contractor_memberships(contractor_id)`,
		`CREATE INDEX IF NOT EXISTS idx_contractor_memberships_status ON contractor_memberships(status)`,

		// T11d. Виправлення: попередні версії бекфілу могли призначити підрядникам дефолтний
		// tenant. Повертаємо їх до глобального стану (tenant_id IS NULL), щоб RLS-контекст був
		// порожнім і підрядник знову бачив крос-tenant дошку відкритих завдань.
		`UPDATE users SET tenant_id = NULL WHERE role = 'CONTRACTOR' AND tenant_id IS NOT NULL;`,

		// T12. Додати tenant_id до warehouses
		`ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		`UPDATE warehouses SET tenant_id = (SELECT id FROM tenants WHERE slug = 'default') WHERE tenant_id IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_warehouses_tenant ON warehouses(tenant_id);`,

		// T13. Додати tenant_id до shipments
		`ALTER TABLE shipments ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		`UPDATE shipments SET tenant_id = (SELECT id FROM tenants WHERE slug = 'default') WHERE tenant_id IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_shipments_tenant ON shipments(tenant_id);`,

		// T14. Додати tenant_id до audit_logs
		`ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE SET NULL;`,
		`UPDATE audit_logs SET tenant_id = (SELECT id FROM tenants WHERE slug = 'default') WHERE tenant_id IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant ON audit_logs(tenant_id);`,

		// T15. Перенести subscription_tier з units на tenants.
		// Беремо найвищий тариф серед units tenant-у (якщо був) і ставимо на tenant.
		// Умова: оновлюємо тільки якщо tenant ще має BASIC (не перезаписуємо вручну встановлений тариф).
		`DO $$
		DECLARE r RECORD;
		BEGIN
			FOR r IN SELECT t.id AS tenant_id,
					COALESCE(MAX(CASE u.subscription_tier
						WHEN 'ENTERPRISE' THEN 3
						WHEN 'PRO' THEN 2
						WHEN 'BASIC' THEN 1
						ELSE 0 END), 0) AS lvl
				FROM tenants t
				LEFT JOIN units u ON u.tenant_id = t.id
				WHERE t.subscription_tier = 'BASIC'
				GROUP BY t.id
			LOOP
				UPDATE tenants SET subscription_tier = CASE r.lvl
					WHEN 3 THEN 'ENTERPRISE'
					WHEN 2 THEN 'PRO'
					WHEN 1 THEN 'BASIC'
					ELSE subscription_tier END
				WHERE id = r.tenant_id AND r.lvl > 1;
			END LOOP;
		END $$;`,

		// T16. Додати tenant_id у GPS таблиці (опціонально, лишаємо nullable — працює через units)
		`ALTER TABLE gps_locations ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		`UPDATE gps_locations gl SET tenant_id = u.tenant_id FROM units u WHERE gl.unit_id = u.id AND gl.tenant_id IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_gps_locations_tenant ON gps_locations(tenant_id);`,

		`ALTER TABLE geofences ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;`,
		`UPDATE geofences g SET tenant_id = u.tenant_id FROM units u WHERE g.unit_id = u.id AND g.tenant_id IS NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_geofences_tenant ON geofences(tenant_id);`,

		// T16b. Зберігати відстань рейсу (планова) та фактичний пробіг після завершення
		`ALTER TABLE shipments ADD COLUMN IF NOT EXISTS distance_km FLOAT DEFAULT 0;`,
		`ALTER TABLE shipments ADD COLUMN IF NOT EXISTS actual_km FLOAT DEFAULT 0;`,

		// T16a. Додати tenant_id до notifications
		`CREATE TABLE IF NOT EXISTS notifications (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			type VARCHAR(50) NOT NULL,
			title VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			related_id UUID,
			is_read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			read_at TIMESTAMP WITH TIME ZONE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id ON notifications(tenant_id);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);`,

		// ============================================================
		// T17. RLS — Row-Level Security (defense-in-depth).
		// Політика: якщо app.tenant_id не встановлено / порожній — повний доступ
		// (це режим SYSTEM_ADMIN/міграцій). Інакше рядки фільтруються за tenant_id.
		// FORCE потрібен, бо застосунок конектиться від owner-а таблиць.
		// ============================================================
		`DO $$
		DECLARE
			t TEXT;
			tables TEXT[] := ARRAY[
				'users','units','resources','categories','warehouses',
				'vehicles','fuel_records','supply_requests','contractor_requests',
				'shipments','audit_logs','gps_locations','geofences','notifications'
			];
		BEGIN
			FOREACH t IN ARRAY tables LOOP
				IF EXISTS (SELECT 1 FROM information_schema.tables
						WHERE table_schema = 'public' AND table_name = t) THEN
					EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY', t);
					EXECUTE format('ALTER TABLE %I FORCE ROW LEVEL SECURITY', t);
					EXECUTE format('DROP POLICY IF EXISTS tenant_isolation ON %I', t);
					EXECUTE format($f$
						CREATE POLICY tenant_isolation ON %I
						USING (
							current_setting('app.tenant_id', true) IS NULL
							OR current_setting('app.tenant_id', true) = ''
							OR tenant_id::text = current_setting('app.tenant_id', true)
						)
						WITH CHECK (
							current_setting('app.tenant_id', true) IS NULL
							OR current_setting('app.tenant_id', true) = ''
							OR tenant_id::text = current_setting('app.tenant_id', true)
						)
					$f$, t);
				END IF;
			END LOOP;
		END $$;`,

		// ============================================================
		// T18. Журнал дозаправок під час рейсу.
		// Кожен запис = одна заправка водія в дорозі.
		// Одночасно пишеться в fuel_records як REFUEL (через сервіс).
		// ============================================================
		`CREATE TABLE IF NOT EXISTS shipment_refuels (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			shipment_id UUID NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
			vehicle_id  UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
			liters      NUMERIC(10,2) NOT NULL CHECK (liters > 0),
			odometer_km INT,
			station_name VARCHAR(200),
			cost_uah    NUMERIC(10,2),
			created_by  UUID REFERENCES users(id) ON DELETE SET NULL,
			tenant_id   UUID REFERENCES tenants(id) ON DELETE CASCADE,
			created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_shipment_refuels_shipment ON shipment_refuels(shipment_id)`,
		`CREATE INDEX IF NOT EXISTS idx_shipment_refuels_vehicle  ON shipment_refuels(vehicle_id)`,
		// RLS для shipment_refuels
		`DO $$ BEGIN
			IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema='public' AND table_name='shipment_refuels') THEN
				ALTER TABLE shipment_refuels ENABLE ROW LEVEL SECURITY;
				ALTER TABLE shipment_refuels FORCE ROW LEVEL SECURITY;
				DROP POLICY IF EXISTS tenant_isolation ON shipment_refuels;
				CREATE POLICY tenant_isolation ON shipment_refuels
					USING (current_setting('app.tenant_id',true) IS NULL OR current_setting('app.tenant_id',true)='' OR tenant_id::text=current_setting('app.tenant_id',true))
					WITH CHECK (current_setting('app.tenant_id',true) IS NULL OR current_setting('app.tenant_id',true)='' OR tenant_id::text=current_setting('app.tenant_id',true));
			END IF;
		END $$`,

		// ============================================================
		// T18. Створення SYSTEM_ADMIN (platform owner).
		// Credentials беруться з ENV або fallback на дефолтні для демо.
		// SYSTEM_ADMIN не прив'язаний до tenant (tenant_id = NULL) — крос-тенантний доступ.
		// ============================================================
		`DO $$
		DECLARE
			admin_email TEXT := COALESCE(current_setting('myapp.system_admin_email', true), 'platform@omnilog.system');
			admin_password TEXT := COALESCE(current_setting('myapp.system_admin_password', true), '$2b$12$2N2cQVoJVoY8Zp23weSPGup24cChzzXDk90dcOiZOcgUo1hQgNCFS'); -- AdminSystem2024!
		BEGIN
			-- Перевіряємо чи вже існує SYSTEM_ADMIN
			IF NOT EXISTS (SELECT 1 FROM users WHERE role = 'SYSTEM_ADMIN' LIMIT 1) THEN
				INSERT INTO users (
					id, tenant_id, username, email, full_name, 
					password_hash, role, status, unit_id
				) VALUES (
					gen_random_uuid(),
					NULL, -- крос-тенантний
					'system_admin',
					admin_email,
					'Platform Administrator',
					admin_password,
					'SYSTEM_ADMIN',
					'ACTIVE',
					NULL
				);
				RAISE NOTICE 'SYSTEM_ADMIN created: %', admin_email;
			ELSE
				RAISE NOTICE 'SYSTEM_ADMIN already exists, skipping';
			END IF;
		END $$;`,
	}

	for i, m := range migrations {
		if _, err := pool.Exec(ctx, m); err != nil {
			return fmt.Errorf("migration %d: %w", i+1, err)
		}
	}
	return nil
}
