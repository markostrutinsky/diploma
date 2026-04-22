-- MilLog: повна схема БД (виконується при першому запуску PostgreSQL)

-- 0. Тенанти (організації)
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    subscription_tier VARCHAR(30) NOT NULL DEFAULT 'FREE',
    subscription_expires_at TIMESTAMP WITH TIME ZONE,
    owner_email VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_tenants_slug ON tenants(slug);

-- 1. Користувачі
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    phone VARCHAR(50),
    password_hash TEXT,
    role VARCHAR(30) NOT NULL DEFAULT 'CONTRACTOR',
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    unit_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, username),
    UNIQUE (tenant_id, email)
);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id);

-- 2. Invite-токени для встановлення паролю
CREATE TABLE IF NOT EXISTS invite_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_invite_tokens_hash ON invite_tokens(token_hash);

-- 3. Підрозділи (units має бути до resources через unit_id)
CREATE TABLE IF NOT EXISTS units (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES units(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    unit_type VARCHAR(20) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_units_tenant ON units(tenant_id);

-- 4. Категорії ресурсів
CREATE TABLE IF NOT EXISTS resource_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, name)
);
CREATE INDEX IF NOT EXISTS idx_resource_categories_tenant ON resource_categories(tenant_id);

-- 5. Ресурси (зі складом unit_id)
CREATE TABLE IF NOT EXISTS resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES resource_categories(id) ON DELETE RESTRICT,
    unit_id BIGINT REFERENCES units(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    quantity INTEGER NOT NULL DEFAULT 0,
    serial_number VARCHAR(100),
    location VARCHAR(255),
    condition VARCHAR(20) NOT NULL DEFAULT 'NEW',
    min_quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_resources_category ON resources(category_id);
CREATE INDEX IF NOT EXISTS idx_resources_unit ON resources(unit_id);
CREATE INDEX IF NOT EXISTS idx_resources_tenant ON resources(tenant_id);

-- 6. Заявки на постачання
CREATE TABLE IF NOT EXISTS supply_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP WITH TIME ZONE,
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_supply_requests_status ON supply_requests(status);
CREATE INDEX IF NOT EXISTS idx_supply_requests_created_by ON supply_requests(created_by);
CREATE INDEX IF NOT EXISTS idx_supply_requests_tenant ON supply_requests(tenant_id);

-- 7. Транспорт
CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100),
    plate_number VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    driver_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, plate_number)
);
CREATE INDEX IF NOT EXISTS idx_vehicles_tenant ON vehicles(tenant_id);

CREATE TABLE IF NOT EXISTS fuel_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    liters DECIMAL(10,2) NOT NULL,
    odometer_km INTEGER,
    record_type VARCHAR(20) NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_fuel_records_vehicle ON fuel_records(vehicle_id);
CREATE INDEX IF NOT EXISTS idx_fuel_records_tenant ON fuel_records(tenant_id);

-- 8. Заявки для волонтерів
CREATE TABLE IF NOT EXISTS contractor_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    unit_id BIGINT REFERENCES units(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    taken_by UUID REFERENCES users(id) ON DELETE SET NULL,
    taken_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_contractor_requests_status ON contractor_requests(status);
CREATE INDEX IF NOT EXISTS idx_contractor_requests_tenant ON contractor_requests(tenant_id);

-- 9. Refresh-токени
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens(token_hash);

-- ============================================================
-- Роль для застосунку (не-суперюзер, щоб RLS реально працювала).
-- Суперюзер (postgres) BYPASSRLS за замовчуванням, тому застосунок
-- має підключатися як omnilog_app.
-- ============================================================
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'omnilog_app') THEN
        CREATE ROLE omnilog_app LOGIN PASSWORD 'omnilog_app' NOSUPERUSER NOBYPASSRLS;
    END IF;
END $$;

GRANT CONNECT ON DATABASE omnilog TO omnilog_app;
GRANT USAGE, CREATE ON SCHEMA public TO omnilog_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO omnilog_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO omnilog_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO omnilog_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT USAGE, SELECT ON SEQUENCES TO omnilog_app;

-- Переводимо всі існуючі таблиці та послідовності у власність omnilog_app,
-- щоб застосунок міг виконувати ALTER/RLS міграції і щоб FORCE RLS
-- поширювався на нього (superuser bypass-ить RLS).
DO $$
DECLARE r RECORD;
BEGIN
    FOR r IN SELECT tablename FROM pg_tables WHERE schemaname='public' LOOP
        EXECUTE format('ALTER TABLE public.%I OWNER TO omnilog_app', r.tablename);
    END LOOP;
    FOR r IN SELECT sequence_name FROM information_schema.sequences WHERE sequence_schema='public' LOOP
        EXECUTE format('ALTER SEQUENCE public.%I OWNER TO omnilog_app', r.sequence_name);
    END LOOP;
END $$;
