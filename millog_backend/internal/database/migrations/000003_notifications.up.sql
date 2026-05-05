-- Таблиця сповіщень для користувачів
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- SHIPMENT_ASSIGNED, REQUEST_APPROVED, REQUEST_REJECTED тощо
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    related_id UUID, -- ID пов'язаної сутності (shipment_id, request_id тощо)
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP
);

-- Індекси для швидкого пошуку
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_tenant_id ON notifications(tenant_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);

-- RLS політики
ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;

CREATE POLICY notifications_tenant_isolation ON notifications
    USING (tenant_id = current_setting('app.current_tenant_id', true)::uuid);

CREATE POLICY notifications_user_access ON notifications
    USING (user_id = current_setting('app.current_user_id', true)::uuid);

COMMENT ON TABLE notifications IS 'Системні сповіщення для користувачів';
COMMENT ON COLUMN notifications.type IS 'Тип сповіщення для фільтрації та іконок';
COMMENT ON COLUMN notifications.related_id IS 'UUID пов''язаної сутності для навігації';
