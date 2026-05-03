# 📋 Інструкція з тестування API ендпоінтів

## 🔧 Підготовка до тестування

### 1. Запуск проєкту
```bash
# З кореневої директорії проєкту
docker compose up -d
```

### 2. Заповнення тестової БД
```bash
cd scripts/seed_db
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

# Заповнити БД тестовими даними
python seed.py

# Або скинути та заповнити наново
python seed.py --reset
```

### 3. Тестові облікові записи

**Паролі всіх користувачів**: `password123`

#### Адміністратори
- `admin@Omnilog.local` - Головний адмін (ADMIN)
- `admin2@Omnilog.local` - Резервний адмін (ADMIN)

#### PRO тариф (Регіон "Захід")
- `director.west@Omnilog.local` - Директор регіону (REGION_DIRECTOR)
- `logist.west@Omnilog.local` - Логіст регіону (REGION_LOGISTICIAN)
- `storekeeper.west@Omnilog.local` - Комірник регіону (REGION_STOREKEEPER)
- `manager.lviv@Omnilog.local` - Менеджер філії Львів (BRANCH_MANAGER)
- `logist.lviv@Omnilog.local` - Логіст філії Львів (BRANCH_LOGISTICIAN)
- `employee1.lviv@Omnilog.local` - Працівник (EMPLOYEE)

#### BASIC тариф (Регіон "Центр")
- `director.center@Omnilog.local` - Директор (отримуватиме 402 на PRO-фічах)
- `logist.center@Omnilog.local` - Логіст
- `manager.kyiv@Omnilog.local` - Менеджер Київ

#### BASIC тариф - наближення до ліміту (Регіон "Схід")
- `director.east@Omnilog.local` - Директор (9/10 складів)
- `logist.east@Omnilog.local` - Логіст

#### ENTERPRISE тариф
- `director.test@Omnilog.local` - Enterprise директор
- `logist.test@Omnilog.local` - Enterprise логіст

#### Підрядники/Волонтери
- `contractor1@Omnilog.local` - Волонтер Богдан (CONTRACTOR)
- `contractor2@Omnilog.local` - Волонтер Іванна (CONTRACTOR)
- `contractor3@Omnilog.local` - Волонтер Руслан (CONTRACTOR)

#### Спеціальні статуси
- `blocked@Omnilog.local` - Заблокований користувач (BLOCKED)
- `pending@Omnilog.local` - Новий співробітник (PENDING)

---

## 🌐 Базова URL
```
http://localhost:8080/api
```

---

## 🔐 1. АУТЕНТИФІКАЦІЯ (`/api/auth`)

### 1.1. Вхід в систему
```bash
# POST /api/auth/login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin@Omnilog.local",
    "password": "password123"
  }'
```

**Відповідь:**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {
    "id": "...",
    "username": "admin",
    "email": "admin@Omnilog.local",
    "role": "ADMIN",
    "status": "ACTIVE"
  }
}
```

**Зберігаємо токен для подальших запитів:**
```bash
export TOKEN="eyJhbGc..."
```

### 1.2. Оновлення токена
```bash
# POST /api/auth/refresh
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGc..."
  }'
```

### 1.3. Реєстрація підрядника (публічна)
```bash
# POST /api/auth/register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newcontractor",
    "email": "newcontractor@example.com",
    "password": "password123",
    "full_name": "Новий Підрядник",
    "phone": "+380501234567"
  }'
```

### 1.4. Створення нового тенанта (організації)
```bash
# POST /api/auth/tenants/signup
curl -X POST http://localhost:8080/api/auth/tenants/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Нова Організація",
    "slug": "nova-org",
    "owner_email": "owner@neworg.com",
    "owner_full_name": "Власник Організації",
    "owner_password": "password123"
  }'
```

### 1.5. Отримання інформації про поточного користувача
```bash
# GET /api/auth/me
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

### 1.6. Встановлення паролю за інвайт-токеном
```bash
# POST /api/auth/setup-password
curl -X POST http://localhost:8080/api/auth/setup-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "invite-token-from-email",
    "password": "newpassword123"
  }'
```

### 1.7. Відновлення паролю
```bash
# POST /api/auth/forgot-password
curl -X POST http://localhost:8080/api/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@Omnilog.local"
  }'
```

---

## 👥 2. КОРИСТУВАЧІ (`/api/users`)

### 2.1. Список командирів/менеджерів
```bash
# GET /api/users/commanders
curl -X GET http://localhost:8080/api/users/commanders \
  -H "Authorization: Bearer $TOKEN"
```

### 2.2. Видимі користувачі (в межах підрозділу)
```bash
# GET /api/users/visible
curl -X GET http://localhost:8080/api/users/visible \
  -H "Authorization: Bearer $TOKEN"
```

### 2.3. Отримання лімітів користувача
```bash
# GET /api/users/limits
curl -X GET http://localhost:8080/api/users/limits \
  -H "Authorization: Bearer $TOKEN"
```

### 2.4. Оновлення ролі та підрозділу користувача
```bash
# PUT /api/users/:id/role
curl -X PUT http://localhost:8080/api/users/{user_id}/role \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "BRANCH_LOGISTICIAN",
    "unit_id": 123
  }'
```

### 2.5. Блокування користувача
```bash
# PUT /api/users/:id/block
curl -X PUT http://localhost:8080/api/users/{user_id}/block \
  -H "Authorization: Bearer $TOKEN"
```

### 2.6. Розблокування користувача
```bash
# PUT /api/users/:id/unblock
curl -X PUT http://localhost:8080/api/users/{user_id}/unblock \
  -H "Authorization: Bearer $TOKEN"
```

### 2.7. Оновлення свого профілю
```bash
# PATCH /api/users/profile
curl -X PATCH http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Оновлене Ім'\''я",
    "phone": "+380501234567"
  }'
```

### 2.8. Зміна свого паролю
```bash
# PATCH /api/users/password
curl -X PATCH http://localhost:8080/api/users/password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword456"
  }'
```

---

## 🛠️ 3. АДМІНІСТРУВАННЯ (`/api/admin`)

**Вимагає ролі:** ADMIN, REGION_DIRECTOR, BRANCH_MANAGER та інші керівні ролі

### 3.1. Створення нового користувача
```bash
# POST /api/admin/users
curl -X POST http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@omnilog.local",
    "full_name": "Новий Користувач",
    "role": "EMPLOYEE",
    "unit_id": 123,
    "phone": "+380501234567"
  }'
```

### 3.2. Отримання аудит-логів
```bash
# GET /api/admin/audit-logs
curl -X GET "http://localhost:8080/api/admin/audit-logs?limit=50&offset=0" \
  -H "Authorization: Bearer $TOKEN"
```

### 3.3. Запуск перевірки SLA
```bash
# POST /api/admin/sla/trigger
curl -X POST http://localhost:8080/api/admin/sla/trigger \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🏢 4. ПЛАТФОРМНЕ УПРАВЛІННЯ (`/api/platform`)

**Вимагає ролі:** SYSTEM_ADMIN (крос-тенантні операції)

### 4.1. Статистика платформи
```bash
# GET /api/platform/stats
curl -X GET http://localhost:8080/api/platform/stats \
  -H "Authorization: Bearer $TOKEN"
```

### 4.2. Список всіх тенантів
```bash
# GET /api/platform/tenants
curl -X GET http://localhost:8080/api/platform/tenants \
  -H "Authorization: Bearer $TOKEN"
```

### 4.3. Отримання інформації про тенанта
```bash
# GET /api/platform/tenants/:id
curl -X GET http://localhost:8080/api/platform/tenants/{tenant_id} \
  -H "Authorization: Bearer $TOKEN"
```

### 4.4. Оновлення тарифного плану тенанта
```bash
# PATCH /api/platform/tenants/:id/tier
curl -X PATCH http://localhost:8080/api/platform/tenants/{tenant_id}/tier \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "subscription_tier": "PRO",
    "subscription_expires_at": "2026-12-31T23:59:59Z"
  }'
```

### 4.5. Активація/деактивація тенанта
```bash
# PATCH /api/platform/tenants/:id/active
curl -X PATCH http://localhost:8080/api/platform/tenants/{tenant_id}/active \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_active": false
  }'
```

### 4.6. Видалення тенанта
```bash
# DELETE /api/platform/tenants/:id
curl -X DELETE http://localhost:8080/api/platform/tenants/{tenant_id} \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🏗️ 5. ПІДРОЗДІЛИ (`/api/units`)

### 5.1. Список підрозділів
```bash
# GET /api/units
curl -X GET http://localhost:8080/api/units \
  -H "Authorization: Bearer $TOKEN"
```

### 5.2. Доступні підрозділи для ролі
```bash
# GET /api/units/available
curl -X GET http://localhost:8080/api/units/available \
  -H "Authorization: Bearer $TOKEN"
```

### 5.3. Моя ієрархія підрозділів
```bash
# GET /api/units/my-hierarchy
curl -X GET http://localhost:8080/api/units/my-hierarchy \
  -H "Authorization: Bearer $TOKEN"
```

### 5.4. Створення підрозділу
```bash
# POST /api/units
curl -X POST http://localhost:8080/api/units \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Новий підрозділ",
    "unit_type": "DEPARTMENT",
    "parent_id": 123,
    "subscription_tier": "BASIC"
  }'
```

### 5.5. Зміна керівника підрозділу
```bash
# POST /api/units/:id/change-commander
curl -X POST http://localhost:8080/api/units/{unit_id}/change-commander \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "commander_id": "user-uuid"
  }'
```

### 5.6. Оновлення підрозділу
```bash
# PATCH /api/units/:id
curl -X PATCH http://localhost:8080/api/units/{unit_id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Оновлена назва",
    "subscription_tier": "PRO"
  }'
```

### 5.7. Видалення підрозділу
```bash
# DELETE /api/units/:id
curl -X DELETE http://localhost:8080/api/units/{unit_id} \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📦 6. ІНВЕНТАР (`/api/inventory`)

### 6.1. Категорії ресурсів

#### Список категорій
```bash
# GET /api/inventory/categories
curl -X GET http://localhost:8080/api/inventory/categories \
  -H "Authorization: Bearer $TOKEN"
```

#### Створення категорії
```bash
# POST /api/inventory/categories
curl -X POST http://localhost:8080/api/inventory/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Нова категорія",
    "description": "Опис категорії"
  }'
```

#### Оновлення категорії
```bash
# PATCH /api/inventory/categories/:id
curl -X PATCH http://localhost:8080/api/inventory/categories/{category_id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Оновлена назва",
    "description": "Новий опис"
  }'
```

#### Видалення категорії
```bash
# DELETE /api/inventory/categories/:id
curl -X DELETE http://localhost:8080/api/inventory/categories/{category_id} \
  -H "Authorization: Bearer $TOKEN"
```

### 6.2. Ресурси

#### Список ресурсів
```bash
# GET /api/inventory/resources
curl -X GET "http://localhost:8080/api/inventory/resources?category_id=xxx&warehouse_id=xxx" \
  -H "Authorization: Bearer $TOKEN"
```

#### Отримання конкретного ресурсу
```bash
# GET /api/inventory/resources/:id
curl -X GET http://localhost:8080/api/inventory/resources/{resource_id} \
  -H "Authorization: Bearer $TOKEN"
```

#### Створення ресурсу
```bash
# POST /api/inventory/resources
curl -X POST http://localhost:8080/api/inventory/resources \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": "category-uuid",
    "warehouse_id": "warehouse-uuid",
    "unit_id": 123,
    "name": "Ноутбук Lenovo",
    "description": "ThinkPad T14",
    "quantity": 5,
    "min_quantity": 2,
    "unit_type": "PCS",
    "serial_number": "SN-12345",
    "condition": "NEW",
    "weight_kg": 1.5
  }'
```

#### Оновлення ресурсу
```bash
# PATCH /api/inventory/resources/:id
curl -X PATCH http://localhost:8080/api/inventory/resources/{resource_id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 10,
    "location": "Склад А, Полиця 3",
    "condition": "GOOD"
  }'
```

#### Списання ресурсу
```bash
# POST /api/inventory/resources/:id/write-off
curl -X POST http://localhost:8080/api/inventory/resources/{resource_id}/write-off \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 2,
    "reason": "Пошкоджено при транспортуванні"
  }'
```

#### Видалення ресурсу
```bash
# DELETE /api/inventory/resources/:id
curl -X DELETE http://localhost:8080/api/inventory/resources/{resource_id} \
  -H "Authorization: Bearer $TOKEN"
```

#### Призначення ресурсу
```bash
# POST /api/inventory/resources/:id/assign
curl -X POST http://localhost:8080/api/inventory/resources/{resource_id}/assign \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-uuid",
    "quantity": 1
  }'
```

#### Моє обладнання
```bash
# GET /api/inventory/my-equipment
curl -X GET http://localhost:8080/api/inventory/my-equipment \
  -H "Authorization: Bearer $TOKEN"
```

#### Видача ресурсу
```bash
# POST /api/inventory/issue
curl -X POST http://localhost:8080/api/inventory/issue \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "resource_id": "resource-uuid",
    "recipient_id": "user-uuid",
    "quantity": 2
  }'
```

### 6.3. Переміщення (Shipments)

#### Створення переміщення
```bash
# POST /api/inventory/shipments
curl -X POST http://localhost:8080/api/inventory/shipments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_warehouse_id": "warehouse-uuid-1",
    "to_warehouse_id": "warehouse-uuid-2",
    "vehicle_id": "vehicle-uuid",
    "items": [
      {
        "resource_id": "resource-uuid-1",
        "quantity": 5
      },
      {
        "resource_id": "resource-uuid-2",
        "quantity": 3
      }
    ]
  }'
```

#### Список переміщень
```bash
# GET /api/inventory/shipments
curl -X GET "http://localhost:8080/api/inventory/shipments?status=IN_TRANSIT" \
  -H "Authorization: Bearer $TOKEN"
```

#### Прийняття переміщення
```bash
# POST /api/inventory/shipments/:id/receive
curl -X POST http://localhost:8080/api/inventory/shipments/{shipment_id}/receive \
  -H "Authorization: Bearer $TOKEN"
```

#### Завантаження PDF переміщення
```bash
# GET /api/inventory/shipments/:id/pdf
curl -X GET http://localhost:8080/api/inventory/shipments/{shipment_id}/pdf \
  -H "Authorization: Bearer $TOKEN" \
  --output shipment.pdf
```

### 6.4. Ресурси по складу
```bash
# GET /api/inventory/warehouse/:id
curl -X GET http://localhost:8080/api/inventory/warehouse/{warehouse_id} \
  -H "Authorization: Bearer $TOKEN"
```

### 6.5. QR-код ресурсу
```bash
# GET /api/inventory/resources/:id/qr
curl -X GET http://localhost:8080/api/inventory/resources/{resource_id}/qr \
  -H "Authorization: Bearer $TOKEN" \
  --output resource-qr.png
```

### 6.6. Інвентаризація

#### Проведення інвентаризації
```bash
# POST /api/inventory/audit
curl -X POST http://localhost:8080/api/inventory/audit \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "warehouse_id": "warehouse-uuid",
    "items": [
      {
        "resource_id": "resource-uuid-1",
        "actual_quantity": 45,
        "notes": "Все в порядку"
      }
    ]
  }'
```

### 6.7. 🚀 PRO: Імпорт Excel

#### Завантаження шаблону
```bash
# GET /api/inventory/resources/import/template
curl -X GET http://localhost:8080/api/inventory/resources/import/template \
  -H "Authorization: Bearer $TOKEN" \
  --output template.xlsx
```

#### Імпорт ресурсів з Excel (PRO тільки)
```bash
# POST /api/inventory/resources/import
curl -X POST http://localhost:8080/api/inventory/resources/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@resources.xlsx"
```

---

## 📝 7. ЗАЯВКИ НА ПОСТАЧАННЯ (`/api/requests`)

### 7.1. Створення заявки
```bash
# POST /api/requests
curl -X POST http://localhost:8080/api/requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "resource_id": "resource-uuid",
    "quantity": 10,
    "comment": "Терміново необхідно"
  }'
```

### 7.2. Список заявок
```bash
# GET /api/requests
curl -X GET "http://localhost:8080/api/requests?status=PENDING&limit=20" \
  -H "Authorization: Bearer $TOKEN"
```

### 7.3. Отримання заявки за ID
```bash
# GET /api/requests/:id
curl -X GET http://localhost:8080/api/requests/{request_id} \
  -H "Authorization: Bearer $TOKEN"
```

### 7.4. Затвердження заявки
```bash
# POST /api/requests/:id/approve
curl -X POST http://localhost:8080/api/requests/{request_id}/approve \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "comment": "Затверджено"
  }'
```

### 7.5. Відхилення заявки
```bash
# POST /api/requests/:id/reject
curl -X POST http://localhost:8080/api/requests/{request_id}/reject \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "comment": "Недостатньо ресурсів на складі"
  }'
```

### 7.6. Скасування заявки
```bash
# POST /api/requests/:id/cancel
curl -X POST http://localhost:8080/api/requests/{request_id}/cancel \
  -H "Authorization: Bearer $TOKEN"
```

### 7.7. 🚀 PRO: Smart Dispatch

#### Попередній перегляд
```bash
# POST /api/requests/smart-dispatch-preview
curl -X POST http://localhost:8080/api/requests/smart-dispatch-preview \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "request_ids": ["request-uuid-1", "request-uuid-2"]
  }'
```

#### Підтвердження розподілу
```bash
# POST /api/requests/smart-dispatch-confirm
curl -X POST http://localhost:8080/api/requests/smart-dispatch-confirm \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "dispatch_plan": {
      "request-uuid-1": {
        "warehouse_id": "warehouse-uuid",
        "vehicle_id": "vehicle-uuid"
      }
    }
  }'
```

---

## 🤝 8. ВОЛОНТЕРСЬКІ ЗАЯВКИ (`/api/contractor-requests`)

### 8.1. Список заявок
```bash
# GET /api/contractor-requests
curl -X GET "http://localhost:8080/api/contractor-requests?status=OPEN" \
  -H "Authorization: Bearer $TOKEN"
```

### 8.2. Створення заявки (військові)
```bash
# POST /api/contractor-requests
curl -X POST http://localhost:8080/api/contractor-requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Потрібен генератор 5 кВт",
    "description": "Для забезпечення електропостачання",
    "unit_id": 123
  }'
```

### 8.3. Взяти заявку в роботу (волонтер)
```bash
# POST /api/contractor-requests/:id/take
curl -X POST http://localhost:8080/api/contractor-requests/{request_id}/take \
  -H "Authorization: Bearer $TOKEN"
```

### 8.4. Доставлено (волонтер)
```bash
# POST /api/contractor-requests/:id/deliver
curl -X POST http://localhost:8080/api/contractor-requests/{request_id}/deliver \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "delivery_comment": "Доставлено в повному обсязі"
  }'
```

### 8.5. Прийняти на баланс (військові)
```bash
# POST /api/contractor-requests/:id/accept
curl -X POST http://localhost:8080/api/contractor-requests/{request_id}/accept \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "warehouse_id": "warehouse-uuid",
    "resource_mapping": [
      {
        "category_id": "category-uuid",
        "name": "Генератор 5кВт",
        "quantity": 1,
        "unit_type": "PCS"
      }
    ]
  }'
```

### 8.6. Відхилити (військові)
```bash
# POST /api/contractor-requests/:id/reject
curl -X POST http://localhost:8080/api/contractor-requests/{request_id}/reject \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Не відповідає вимогам"
  }'
```

### 8.7. Скасувати (військові)
```bash
# POST /api/contractor-requests/:id/cancel
curl -X POST http://localhost:8080/api/contractor-requests/{request_id}/cancel \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🚗 9. ТРАНСПОРТ (`/api/vehicles`)

### 9.1. Створення транспорту
```bash
# POST /api/vehicles
curl -X POST http://localhost:8080/api/vehicles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "brand": "Mercedes",
    "model": "Sprinter 2021",
    "plate_number": "AA1234BB",
    "type": "VAN",
    "capacity_kg": 2500,
    "tank_capacity": 75,
    "fuel_norm": 10.5,
    "driver_id": "user-uuid"
  }'
```

### 9.2. Список транспорту
```bash
# GET /api/vehicles
curl -X GET http://localhost:8080/api/vehicles \
  -H "Authorization: Bearer $TOKEN"
```

### 9.3. Отримання транспорту за ID
```bash
# GET /api/vehicles/:id
curl -X GET http://localhost:8080/api/vehicles/{vehicle_id} \
  -H "Authorization: Bearer $TOKEN"
```

### 9.4. Оновлення транспорту
```bash
# PATCH /api/vehicles/:id
curl -X PATCH http://localhost:8080/api/vehicles/{vehicle_id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "ACTIVE",
    "capacity_kg": 3000
  }'
```

### 9.5. Оновлення статусу
```bash
# PATCH /api/vehicles/:id/status
curl -X PATCH http://localhost:8080/api/vehicles/{vehicle_id}/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "MAINTENANCE"
  }'
```

### 9.6. Видалення транспорту
```bash
# DELETE /api/vehicles/:id
curl -X DELETE http://localhost:8080/api/vehicles/{vehicle_id} \
  -H "Authorization: Bearer $TOKEN"
```

### 9.7. Призначення водія
```bash
# PATCH /api/vehicles/:id/driver
curl -X PATCH http://localhost:8080/api/vehicles/{vehicle_id}/driver \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "driver_id": "user-uuid"
  }'
```

### 9.8. Історія водіїв
```bash
# GET /api/vehicles/:id/drivers
curl -X GET http://localhost:8080/api/vehicles/{vehicle_id}/drivers \
  -H "Authorization: Bearer $TOKEN"
```

### 9.9. Доступні для маршруту
```bash
# GET /api/vehicles/available-for-route
curl -X GET http://localhost:8080/api/vehicles/available-for-route \
  -H "Authorization: Bearer $TOKEN"
```

---

## ⛽ 10. ПАЛИВО (`/api/vehicles/:id/fuel`)

### 10.1. Додавання запису про паливо
```bash
# POST /api/vehicles/:id/fuel
curl -X POST http://localhost:8080/api/vehicles/{vehicle_id}/fuel \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "liters": 45.5,
    "odometer_km": 125000,
    "record_type": "REFUEL"
  }'
```

### 10.2. Історія пального
```bash
# GET /api/vehicles/:id/fuel
curl -X GET http://localhost:8080/api/vehicles/{vehicle_id}/fuel \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🔧 11. ТЕХНІЧНЕ ОБСЛУГОВУВАННЯ

### 11.1. Виконання ТО
```bash
# POST /api/vehicles/:id/maintenance
curl -X POST http://localhost:8080/api/vehicles/{vehicle_id}/maintenance \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "odometer_km": 130000,
    "description": "Заміна оливи та фільтрів",
    "performed_by": "СТО Центральне",
    "cost_amount": 2500.00
  }'
```

### 11.2. Історія ТО
```bash
# GET /api/vehicles/:id/maintenance
curl -X GET http://localhost:8080/api/vehicles/{vehicle_id}/maintenance \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🏪 12. СКЛАДИ (`/api/warehouses`)

### 12.1. Список складів
```bash
# GET /api/warehouses
curl -X GET http://localhost:8080/api/warehouses \
  -H "Authorization: Bearer $TOKEN"
```

### 12.2. Створення складу
```bash
# POST /api/warehouses
curl -X POST http://localhost:8080/api/warehouses \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "unit_id": 123,
    "name": "Новий склад",
    "location_type": "STATIONARY",
    "latitude": 49.8397,
    "longitude": 24.0297,
    "address": "вул. Шевченка, 1"
  }'
```

### 12.3. Оновлення розташування
```bash
# PATCH /api/warehouses/:id/location
curl -X PATCH http://localhost:8080/api/warehouses/{warehouse_id}/location \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 49.8400,
    "longitude": 24.0300
  }'
```

### 12.4. Оновлення складу
```bash
# PATCH /api/warehouses/:id
curl -X PATCH http://localhost:8080/api/warehouses/{warehouse_id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Оновлена назва",
    "address": "Нова адреса"
  }'
```

### 12.5. Видалення складу
```bash
# DELETE /api/warehouses/:id
curl -X DELETE http://localhost:8080/api/warehouses/{warehouse_id} \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📊 13. АНАЛІТИКА (`/api/analytics`)

### 13.1. Дашборд (доступно на всіх тарифах)
```bash
# GET /api/analytics/dashboard
curl -X GET http://localhost:8080/api/analytics/dashboard \
  -H "Authorization: Bearer $TOKEN"
```

### 13.2. Експорт інвентаря
```bash
# GET /api/analytics/export/inventory
curl -X GET http://localhost:8080/api/analytics/export/inventory \
  -H "Authorization: Bearer $TOKEN" \
  --output inventory-export.xlsx
```

### 13.3. Експорт пального
```bash
# GET /api/analytics/export/fuel
curl -X GET http://localhost:8080/api/analytics/export/fuel \
  -H "Authorization: Bearer $TOKEN" \
  --output fuel-export.xlsx
```

### 13.4. 🚀 PRO: Автоматичне поповнення
```bash
# POST /api/analytics/auto-replenish
curl -X POST http://localhost:8080/api/analytics/auto-replenish \
  -H "Authorization: Bearer $TOKEN"
```

### 13.5. 🚀 PRO: Розширені KPI
```bash
# GET /api/analytics/kpi
curl -X GET http://localhost:8080/api/analytics/kpi \
  -H "Authorization: Bearer $TOKEN"
```

### 13.6. 🚀 PRO: Прогноз попиту
```bash
# GET /api/analytics/forecast
curl -X GET "http://localhost:8080/api/analytics/forecast?resource_id=xxx&days=30" \
  -H "Authorization: Bearer $TOKEN"
```

### 13.7. 🚀 PRO: Прогнозне ТО
```bash
# GET /api/analytics/maintenance
curl -X GET http://localhost:8080/api/analytics/maintenance \
  -H "Authorization: Bearer $TOKEN"
```

### 13.8. 🚀 PRO: Виявлення аномалій пального
```bash
# GET /api/analytics/fuel-anomalies
curl -X GET http://localhost:8080/api/analytics/fuel-anomalies \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📍 14. GPS ТРЕКІНГ (`/api/gps`) - 🚀 PRO FEATURE

**Примітка:** Всі GPS ендпоінти доступні тільки на тарифі PRO та ENTERPRISE

### 14.1. Запис геолокації транспорту
```bash
# POST /api/gps/locations
curl -X POST http://localhost:8080/api/gps/locations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "vehicle_id": "vehicle-uuid",
    "latitude": 49.8397,
    "longitude": 24.0297,
    "speed": 65.5,
    "heading": 180,
    "accuracy": 5.2
  }'
```

### 14.2. Карта флоту
```bash
# GET /api/gps/fleet-map
curl -X GET http://localhost:8080/api/gps/fleet-map \
  -H "Authorization: Bearer $TOKEN"
```

### 14.3. Траєкторія транспорту
```bash
# GET /api/gps/trajectory
curl -X GET "http://localhost:8080/api/gps/trajectory?vehicle_id=xxx&from=2026-05-01T00:00:00Z&to=2026-05-02T00:00:00Z" \
  -H "Authorization: Bearer $TOKEN"
```

### 14.4. Створення геозони
```bash
# POST /api/gps/geofences
curl -X POST http://localhost:8080/api/gps/geofences \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "unit_id": 123,
    "name": "Безпечна зона - штаб",
    "latitude": 49.8397,
    "longitude": 24.0297,
    "radius": 500,
    "type": "SAFE"
  }'
```

### 14.5. Список геозон
```bash
# GET /api/gps/geofences
curl -X GET http://localhost:8080/api/gps/geofences \
  -H "Authorization: Bearer $TOKEN"
```

### 14.6. Оповіщення про геозони
```bash
# GET /api/gps/geofence-alerts
curl -X GET "http://localhost:8080/api/gps/geofence-alerts?limit=50" \
  -H "Authorization: Bearer $TOKEN"
```

### 14.7. Статус флоту
```bash
# GET /api/gps/fleet-status
curl -X GET http://localhost:8080/api/gps/fleet-status \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🎯 15. BOOTSTRAP АДМІНА

**Важливо:** Цей ендпоінт використовується тільки один раз для створення першого адміністратора

```bash
# POST /api/bootstrap
curl -X POST http://localhost:8080/api/bootstrap \
  -H "Content-Type: application/json" \
  -d '{
    "username": "superadmin",
    "email": "superadmin@omnilog.local",
    "password": "securepassword123",
    "full_name": "Супер Адмін"
  }'
```

---

## 📂 Тестові дані в БД

Після запуску `seed.py` в базі даних будуть наступні дані:

### Підрозділи (Units)
- **Регіон "Захід"** (PRO) - 4 склади, 6 автомобілів, ~72 ресурси
  - Філія Львів
  - Філія Івано-Франківськ
  - Відділ логістики Львів
  - Команда кур'єрів Львів

- **Регіон "Центр"** (BASIC) - 2 склади, 3 автомобілі, ~20 ресурсів
  - Філія Київ
  - Відділ Київ

- **Регіон "Схід"** (BASIC, наближення до ліміту) - 9/10 складів, 4 автомобілі, ~72 ресурси
  - Філія Харків
  - Відділ Харків

- **Регіон "Тест-ENTERPRISE"** (ENTERPRISE) - 1 склад, 2 автомобілі, ~5 ресурсів

### Категорії ресурсів
- Канцелярія
- Електроніка
- Інструмент
- Медикаменти
- Продукти
- Одяг
- Паливо-мастильні

### Транспорт
- Різні моделі: Renault Master, Ford Transit, Mercedes Sprinter, Toyota Hilux, Volkswagen Crafter, MAN TGE
- Записи про заправки за останні 30 днів
- Історія технічного обслуговування
- GPS-траєкторії (тільки для PRO регіонів)

### Заявки
- Supply requests (заявки на постачання) в різних статусах: PENDING, APPROVED, REJECTED, COMPLETED
- Contractor requests (волонтерські заявки) в статусах: OPEN, IN_PROGRESS, DELIVERED, COMPLETED

### GPS (тільки PRO регіони)
- Геозони (SAFE, FORBIDDEN)
- GPS-локації транспорту (20 точок на кожне авто)
- Геофенс алерти

---

## 🔍 Тестування тарифних планів

### FREE тариф
- Базові операції з інвентарем
- Обмежений функціонал

### BASIC тариф
- Ліміти:
  - Максимум 10 складів
  - Максимум 50 користувачів
  - Максимум 1000 ресурсів
  - Максимум 20 автомобілів
- Тестуйте з `director.east@Omnilog.local` (9/10 складів)

### PRO тариф
- Всі фічі BASIC +
- GPS трекінг і геозони
- Excel імпорт
- Smart Dispatch
- Розширена аналітика (KPI, Forecast)
- Прогнозне ТО
- Виявлення аномалій пального
- Тестуйте з `director.west@Omnilog.local`

### ENTERPRISE тариф
- Необмежені ресурси
- Всі PRO фічі
- Тестуйте з `director.test@Omnilog.local`

### Перевірка обмежень
Спробуйте під BASIC акаунтом:
```bash
# Це має повернути 402 Payment Required
curl -X GET http://localhost:8080/api/gps/fleet-map \
  -H "Authorization: Bearer $BASIC_TOKEN"
```

---

## 📝 Корисні поради

### 1. Збереження токена в змінну
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin@Omnilog.local","password":"password123"}' \
  | jq -r '.access_token')
```

### 2. Форматований вивід з jq
```bash
curl -X GET http://localhost:8080/api/inventory/resources \
  -H "Authorization: Bearer $TOKEN" | jq .
```

### 3. Перевірка помилок
```bash
# Спроба доступу без токена (має повернути 401)
curl -X GET http://localhost:8080/api/inventory/resources

# Спроба доступу до PRO фічі з BASIC акаунтом (має повернути 402)
curl -X GET http://localhost:8080/api/gps/fleet-map \
  -H "Authorization: Bearer $BASIC_TOKEN"

# Спроба доступу до чужих даних (має повернути 403)
curl -X GET http://localhost:8080/api/inventory/resources \
  -H "Authorization: Bearer $OTHER_UNIT_TOKEN"
```

### 4. Використання Postman
Імпортуйте всі запити в Postman, створивши collection з Environment variables:
- `base_url`: `http://localhost:8080/api`
- `token`: `{{access_token}}`

### 5. Автоматичне тестування
Створіть bash-скрипт для тестування всіх ендпоінтів:
```bash
#!/bin/bash
source test-all-endpoints.sh
```

---

## 🐛 Очікувані поведінки

### Успішні відповіді
- `200 OK` - Успішна операція
- `201 Created` - Ресурс створено
- `204 No Content` - Операція виконана, відповідь порожня

### Помилки
- `400 Bad Request` - Невалідні дані
- `401 Unauthorized` - Токен відсутній або невалідний
- `402 Payment Required` - Потрібен вищий тарифний план
- `403 Forbidden` - Недостатньо прав або заборонений доступ
- `404 Not Found` - Ресурс не знайдено
- `409 Conflict` - Конфлікт даних (напр., дублікат)
- `422 Unprocessable Entity` - Бізнес-логіка не дозволяє операцію
- `500 Internal Server Error` - Серверна помилка

---

## 📧 Контакти та підтримка

При виникненні проблем:
1. Перевірте логи backend: `docker compose logs backend`
2. Перевірте стан БД: `docker compose logs postgres`
3. Перезапустіть seed: `python seed.py --reset`
4. Перезапустіть контейнери: `docker compose restart`

---

**Версія документа:** 1.0  
**Дата оновлення:** 2 травня 2026  
**Проєкт:** Omnilog / MilLog
