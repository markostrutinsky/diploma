# 📖 API Documentation - Omnilog (Omnilog)

## Базова інформація

**Base URL:** `http://localhost/api`  
**Аутентифікація:** JWT Bearer Token  
**Content-Type:** `application/json`  
**Версія API:** 1.0

---

## 📋 Зміст

- [Аутентифікація](#-аутентифікація)
- [Користувачі](#-користувачі)
- [Підрозділи](#-підрозділи)
- [Інвентар (Ресурси)](#-інвентар-ресурси)
- [Заявки на постачання](#-заявки-на-постачання)
- [Волонтерські заявки](#-волонтерські-заявки)
- [Автопарк](#-автопарк)
- [Паливо](#-паливо)
- [Склади](#-склади)
- [Аналітика](#-аналітика)
- [GPS Трекінг (PRO)](#-gps-трекінг-pro)
- [Platform Admin](#-platform-admin)
- [Audit Logs](#-audit-logs)

---

## 🔐 Аутентифікація

### Реєстрація організації (Tenant Signup)

Створює нову організацію та першого адміністратора.

**Endpoint:** `POST /api/auth/tenants/signup`  
**Auth:** ❌ Не потрібна  
**Tier:** Всі

**Request:**
```json
{
  "tenant_name": "26-та Окрема Бригада",
  "email": "admin@brigade26.mil",
  "password": "SecurePassword123!"
}
```

**Response:** `201 Created`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@brigade26.mil",
    "name": "admin@brigade26.mil",
    "role": "TENANT_ADMIN",
    "tenant_id": "660e8400-e29b-41d4-a716-446655440000",
    "tenant_name": "26-та Окрема Бригада"
  }
}
```

**Errors:**
- `400 Bad Request` - Невалідні дані
- `409 Conflict` - Email вже зайнятий

---

### Логін

**Endpoint:** `POST /api/auth/login`  
**Auth:** ❌ Не потрібна  
**Tier:** Всі

**Request:**
```json
{
  "email": "admin@brigade26.mil",
  "password": "SecurePassword123!"
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@brigade26.mil",
    "name": "Admin User",
    "role": "TENANT_ADMIN",
    "tenant_id": "660e8400-e29b-41d4-a716-446655440000",
    "unit_id": null,
    "blocked": false
  }
}
```

**Errors:**
- `401 Unauthorized` - Невірний email або пароль
- `403 Forbidden` - Акаунт заблоковано

---

### Refresh Token

Оновлення access токена за допомогою refresh токена.

**Endpoint:** `POST /api/auth/refresh`  
**Auth:** ❌ Не потрібна  
**Tier:** Всі

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc..."
}
```

---

### Реєстрація підрядника (Волонтер)

**Endpoint:** `POST /api/auth/register`  
**Auth:** ❌ Не потрібна  
**Tier:** Всі

**Request:**
```json
{
  "email": "volunteer@ngo.org",
  "password": "SecurePass123!",
  "name": "Іван Петренко",
  "organization": "NGO Допомога" // опціонально
}
```

**Response:** `201 Created`
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "email": "volunteer@ngo.org",
    "name": "Іван Петренко",
    "role": "CONTRACTOR",
    "tenant_id": null
  }
}
```

---

### Setup Password (після інвайту)

**Endpoint:** `POST /api/auth/setup-password`  
**Auth:** ❌ Не потрібна  
**Tier:** Всі

**Request:**
```json
{
  "token": "invite-token-from-email",
  "password": "NewSecurePass123!",
  "name": "Сергій Іваненко" // опціонально
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": { ... }
}
```

---

### Forgot Password

**Endpoint:** `POST /api/auth/forgot-password`  
**Auth:** ❌ Не потрібна  
**Tier:** Всі

**Request:**
```json
{
  "email": "admin@brigade26.mil"
}
```

**Response:** `200 OK`
```json
{
  "message": "Інструкції для відновлення паролю надіслано на email"
}
```

---

### Get Current User (Me)

**Endpoint:** `GET /api/auth/me`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Headers:**
```
Authorization: Bearer eyJhbGc...
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "admin@brigade26.mil",
  "name": "Admin User",
  "role": "TENANT_ADMIN",
  "tenant_id": "660e8400-e29b-41d4-a716-446655440000",
  "tenant_name": "26-та Окрема Бригада",
  "subscription_tier": "PRO",
  "unit_id": null,
  "blocked": false
}
```

---

## 👥 Користувачі

### Створити користувача (Invite)

**Endpoint:** `POST /api/admin/users`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN, TENANT_ADMIN, REGION_DIRECTOR  
**Tier:** Всі

**Request:**
```json
{
  "email": "commander@brigade26.mil",
  "role": "COMMANDER",
  "unit_id": "880e8400-e29b-41d4-a716-446655440000"
}
```

**Response:** `201 Created`
```json
{
  "message": "Користувача створено та запрошено по email",
  "user_id": "990e8400-e29b-41d4-a716-446655440000",
  "invite_token": "abc123..."
}
```

---

### Список командирів

**Endpoint:** `GET /api/users/commanders`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Іван Петренко",
    "email": "commander@brigade26.mil",
    "role": "COMMANDER",
    "unit_id": "880e8400-e29b-41d4-a716-446655440000",
    "unit_name": "1-й Батальйон"
  }
]
```

---

### Список видимих користувачів

**Endpoint:** `GET /api/users/visible`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Query Parameters:**
- `role` (optional) - фільтр по ролі
- `unit_id` (optional) - фільтр по підрозділу

**Example:** `GET /api/users/visible?role=WAREHOUSE`

**Response:** `200 OK`
```json
[
  {
    "id": "aa0e8400-e29b-41d4-a716-446655440000",
    "name": "Петро Коваль",
    "email": "warehouse@brigade26.mil",
    "role": "WAREHOUSE",
    "unit_name": "Склад №1",
    "blocked": false
  }
]
```

---

### Отримати ліміти користувача

**Endpoint:** `GET /api/users/limits`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "max_users": 50,
  "current_users": 12,
  "max_resources": 500,
  "current_resources": 234,
  "max_vehicles": 20,
  "current_vehicles": 8,
  "tier": "STANDARD"
}
```

---

### Оновити роль та підрозділ

**Endpoint:** `PUT /api/users/:id/role`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN, TENANT_ADMIN, REGION_DIRECTOR  
**Tier:** Всі

**Request:**
```json
{
  "role": "COMMANDER",
  "unit_id": "880e8400-e29b-41d4-a716-446655440000"
}
```

**Response:** `200 OK`
```json
{
  "message": "Роль та підрозділ оновлено"
}
```

---

### Заблокувати користувача

**Endpoint:** `PUT /api/users/:id/block`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN, TENANT_ADMIN  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "message": "Користувача заблоковано"
}
```

---

### Розблокувати користувача

**Endpoint:** `PUT /api/users/:id/unblock`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN, TENANT_ADMIN  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "message": "Користувача розблоковано"
}
```

---

### Оновити свій профіль

**Endpoint:** `PATCH /api/users/profile`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Request:**
```json
{
  "name": "Нове Ім'я",
  "phone": "+380501234567"
}
```

**Response:** `200 OK`
```json
{
  "message": "Профіль оновлено"
}
```

---

### Змінити свій пароль

**Endpoint:** `PATCH /api/users/password`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Request:**
```json
{
  "current_password": "OldPass123!",
  "new_password": "NewSecurePass123!"
}
```

**Response:** `200 OK`
```json
{
  "message": "Пароль змінено"
}
```

---

## 🏢 Підрозділи

### Список підрозділів

**Endpoint:** `GET /api/units`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
[
  {
    "id": "110e8400-e29b-41d4-a716-446655440000",
    "name": "26-та Окрема Бригада",
    "type": "BRIGADE",
    "parent_id": null,
    "commander_id": "220e8400-e29b-41d4-a716-446655440000",
    "commander_name": "Полковник Іваненко"
  },
  {
    "id": "330e8400-e29b-41d4-a716-446655440000",
    "name": "1-й Полк",
    "type": "REGIMENT",
    "parent_id": "110e8400-e29b-41d4-a716-446655440000",
    "commander_id": "440e8400-e29b-41d4-a716-446655440000",
    "commander_name": "Майор Петренко"
  }
]
```

---

### Створити підрозділ

**Endpoint:** `POST /api/units`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN, TENANT_ADMIN, REGION_DIRECTOR  
**Tier:** Всі

**Request:**
```json
{
  "name": "2-й Батальйон",
  "type": "BATTALION",
  "parent_id": "330e8400-e29b-41d4-a716-446655440000",
  "commander_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response:** `201 Created`
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "name": "2-й Батальйон",
  "type": "BATTALION",
  "parent_id": "330e8400-e29b-41d4-a716-446655440000",
  "commander_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

### Оновити підрозділ

**Endpoint:** `PATCH /api/units/:id`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN, TENANT_ADMIN, REGION_DIRECTOR  
**Tier:** Всі

**Request:**
```json
{
  "name": "2-й Батальйон (оновлена назва)",
  "commander_id": "770e8400-e29b-41d4-a716-446655440000"
}
```

**Response:** `200 OK`
```json
{
  "message": "Підрозділ оновлено"
}
```

---

### Видалити підрозділ

**Endpoint:** `DELETE /api/units/:id`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN, TENANT_ADMIN, REGION_DIRECTOR  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "message": "Підрозділ видалено"
}
```

---

### Отримати доступні підрозділи для ролі

**Endpoint:** `GET /api/units/available`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
[
  {
    "id": "110e8400-e29b-41d4-a716-446655440000",
    "name": "26-та Окрема Бригада",
    "type": "BRIGADE"
  }
]
```

---

### Отримати ієрархію мого підрозділу

**Endpoint:** `GET /api/units/my-hierarchy`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "current": {
    "id": "330e8400-e29b-41d4-a716-446655440000",
    "name": "1-й Полк",
    "type": "REGIMENT"
  },
  "parent": {
    "id": "110e8400-e29b-41d4-a716-446655440000",
    "name": "26-та Окрема Бригада",
    "type": "BRIGADE"
  },
  "children": [
    {
      "id": "440e8400-e29b-41d4-a716-446655440000",
      "name": "1-й Батальйон",
      "type": "BATTALION"
    }
  ]
}
```

---

### Змінити командира підрозділу

**Endpoint:** `POST /api/units/:id/change-commander`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN, TENANT_ADMIN, REGION_DIRECTOR  
**Tier:** Всі

**Request:**
```json
{
  "commander_id": "880e8400-e29b-41d4-a716-446655440000"
}
```

**Response:** `200 OK`
```json
{
  "message": "Командира змінено"
}
```

---

## 📦 Інвентар (Ресурси)

### Список категорій

**Endpoint:** `GET /api/inventory/categories`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
[
  {
    "id": "aa0e8400-e29b-41d4-a716-446655440000",
    "name": "Зброя",
    "description": "Стрілецька зброя та боєприпаси",
    "created_at": "2026-01-15T10:30:00Z"
  },
  {
    "id": "bb0e8400-e29b-41d4-a716-446655440000",
    "name": "Обмундирування",
    "description": "Військова форма та екіпіровка",
    "created_at": "2026-01-15T10:35:00Z"
  }
]
```

---

### Створити категорію

**Endpoint:** `POST /api/inventory/categories`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "name": "Медикаменти",
  "description": "Медичні препарати та засоби"
}
```

**Response:** `201 Created`
```json
{
  "id": "cc0e8400-e29b-41d4-a716-446655440000",
  "name": "Медикаменти",
  "description": "Медичні препарати та засоби",
  "created_at": "2026-04-24T01:29:00Z"
}
```

---

### Оновити категорію

**Endpoint:** `PATCH /api/inventory/categories/:id`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "name": "Медикаменти (оновлено)",
  "description": "Нова описка"
}
```

**Response:** `200 OK`
```json
{
  "message": "Категорію оновлено"
}
```

---

### Видалити категорію

**Endpoint:** `DELETE /api/inventory/categories/:id`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "message": "Категорію видалено"
}
```

**Errors:**
- `400 Bad Request` - Категорія містить ресурси

---

### Список ресурсів

**Endpoint:** `GET /api/inventory/resources`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Query Parameters:**
- `category_id` (optional) - фільтр по категорії
- `warehouse_id` (optional) - фільтр по складу
- `critical` (optional) - тільки критичні залишки (true/false)

**Example:** `GET /api/inventory/resources?critical=true`

**Response:** `200 OK`
```json
[
  {
    "id": "dd0e8400-e29b-41d4-a716-446655440000",
    "name": "АК-74М",
    "description": "Автомат Калашникова модернізований",
    "category_id": "aa0e8400-e29b-41d4-a716-446655440000",
    "category_name": "Зброя",
    "warehouse_id": "ee0e8400-e29b-41d4-a716-446655440000",
    "warehouse_name": "Склад №1",
    "quantity": 45,
    "min_quantity": 50,
    "unit": "шт",
    "is_critical": true,
    "assigned_to_user_id": null,
    "created_at": "2026-02-10T09:00:00Z",
    "updated_at": "2026-04-20T14:30:00Z"
  }
]
```

---

### Отримати ресурс за ID

**Endpoint:** `GET /api/inventory/resources/:id`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "id": "dd0e8400-e29b-41d4-a716-446655440000",
  "name": "АК-74М",
  "description": "Автомат Калашникова модернізований",
  "category_id": "aa0e8400-e29b-41d4-a716-446655440000",
  "category_name": "Зброя",
  "warehouse_id": "ee0e8400-e29b-41d4-a716-446655440000",
  "warehouse_name": "Склад №1",
  "quantity": 45,
  "min_quantity": 50,
  "unit": "шт",
  "is_critical": true,
  "price": 15000.00,
  "assigned_to_user_id": null,
  "created_at": "2026-02-10T09:00:00Z",
  "updated_at": "2026-04-20T14:30:00Z"
}
```

---

### Створити ресурс

**Endpoint:** `POST /api/inventory/resources`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "name": "Каска Kevlar",
  "description": "Бронешолом Kevlar PASGT",
  "category_id": "bb0e8400-e29b-41d4-a716-446655440000",
  "warehouse_id": "ee0e8400-e29b-41d4-a716-446655440000",
  "quantity": 120,
  "min_quantity": 30,
  "unit": "шт",
  "price": 3500.00
}
```

**Response:** `201 Created`
```json
{
  "id": "ff0e8400-e29b-41d4-a716-446655440000",
  "name": "Каска Kevlar",
  "quantity": 120,
  "is_critical": false
}
```

---

### Оновити ресурс

**Endpoint:** `PATCH /api/inventory/resources/:id`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "quantity": 150,
  "min_quantity": 40,
  "price": 3600.00
}
```

**Response:** `200 OK`
```json
{
  "message": "Ресурс оновлено"
}
```

---

### Списати ресурс

**Endpoint:** `POST /api/inventory/resources/:id/write-off`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "quantity": 5,
  "reason": "Пошкоджено під час операції"
}
```

**Response:** `200 OK`
```json
{
  "message": "Ресурс списано",
  "new_quantity": 115
}
```

---

### Видалити ресурс

**Endpoint:** `DELETE /api/inventory/resources/:id`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "message": "Ресурс видалено"
}
```

---

### Призначити ресурс користувачу

**Endpoint:** `POST /api/inventory/resources/:id/assign`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "user_id": "110e8400-e29b-41d4-a716-446655440000",
  "quantity": 1
}
```

**Response:** `200 OK`
```json
{
  "message": "Ресурс призначено користувачу"
}
```

---

### Моє обладнання

**Endpoint:** `GET /api/inventory/my-equipment`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
[
  {
    "id": "dd0e8400-e29b-41d4-a716-446655440000",
    "name": "АК-74М",
    "category_name": "Зброя",
    "quantity": 1,
    "assigned_at": "2026-03-15T10:00:00Z"
  }
]
```

---

### Видати ресурс

**Endpoint:** `POST /api/inventory/issue`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Request:**
```json
{
  "resource_id": "dd0e8400-e29b-41d4-a716-446655440000",
  "quantity": 10,
  "issued_to_user_id": "220e8400-e29b-41d4-a716-446655440000",
  "notes": "Видача для тренування"
}
```

**Response:** `200 OK`
```json
{
  "message": "Ресурс видано"
}
```

---

### Створити відправлення між складами

**Endpoint:** `POST /api/inventory/shipments`  
**Auth:** ✅ Потрібна  
**Tier:** BASIC+

**Request:**
```json
{
  "from_warehouse_id": "ee0e8400-e29b-41d4-a716-446655440000",
  "to_warehouse_id": "ff0e8400-e29b-41d4-a716-446655440000",
  "vehicle_id": "110e8400-e29b-41d4-a716-446655440000",
  "items": [
    {
      "resource_id": "dd0e8400-e29b-41d4-a716-446655440000",
      "quantity": 20
    }
  ],
  "notes": "Поповнення складу батальйону"
}
```

**Response:** `201 Created`
```json
{
  "id": "220e8400-e29b-41d4-a716-446655440000",
  "status": "PENDING",
  "created_at": "2026-04-24T01:29:00Z"
}
```

---

### Прийняти відправлення

**Endpoint:** `POST /api/inventory/shipments/:id/receive`  
**Auth:** ✅ Потрібна  
**Tier:** BASIC+

**Request:**
```json
{
  "received_items": [
    {
      "resource_id": "dd0e8400-e29b-41d4-a716-446655440000",
      "quantity_received": 20
    }
  ],
  "notes": "Отримано в повному обсязі"
}
```

**Response:** `200 OK`
```json
{
  "message": "Відправлення прийнято",
  "status": "RECEIVED"
}
```

---

### Список відправлень

**Endpoint:** `GET /api/inventory/shipments`  
**Auth:** ✅ Потрібна  
**Tier:** BASIC+

**Query Parameters:**
- `status` (optional) - PENDING, IN_TRANSIT, RECEIVED, CANCELLED

**Response:** `200 OK`
```json
[
  {
    "id": "220e8400-e29b-41d4-a716-446655440000",
    "from_warehouse_name": "Склад №1",
    "to_warehouse_name": "Склад №2",
    "vehicle_name": "КамАЗ #123",
    "status": "IN_TRANSIT",
    "items_count": 1,
    "created_at": "2026-04-24T01:29:00Z"
  }
]
```

---

### Завантажити накладну PDF

**Endpoint:** `GET /api/inventory/shipments/:id/pdf`  
**Auth:** ✅ Потрібна  
**Tier:** BASIC+

**Response:** `200 OK` (PDF file)
```
Content-Type: application/pdf
Content-Disposition: attachment; filename="shipment-220e8400.pdf"
```

---

### Завантажити QR код ресурсу

**Endpoint:** `GET /api/inventory/resources/:id/qr`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK` (PNG image)
```
Content-Type: image/png
```

---

### Завантажити шаблон Excel для імпорту

**Endpoint:** `GET /api/inventory/resources/import/template`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK` (Excel file)
```
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="import_template.xlsx"
```

---

### Імпортувати ресурси з Excel

**Endpoint:** `POST /api/inventory/resources/import`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** 🌟 **PRO**

**Request:** (multipart/form-data)
```
file: [Excel file]
warehouse_id: ee0e8400-e29b-41d4-a716-446655440000
```

**Response:** `200 OK`
```json
{
  "message": "Імпортовано 45 ресурсів",
  "imported_count": 45,
  "skipped_count": 2,
  "errors": [
    {
      "row": 12,
      "error": "Невалідна категорія"
    }
  ]
}
```

---

### Подати аудит інвентаря

**Endpoint:** `POST /api/inventory/audit`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** BASIC+

**Request:**
```json
{
  "warehouse_id": "ee0e8400-e29b-41d4-a716-446655440000",
  "resources": [
    {
      "resource_id": "dd0e8400-e29b-41d4-a716-446655440000",
      "actual_quantity": 43,
      "expected_quantity": 45
    }
  ],
  "notes": "Виявлено розбіжності"
}
```

**Response:** `201 Created`
```json
{
  "message": "Аудит збережено",
  "discrepancies_count": 1
}
```

---

### Ресурси складу

**Endpoint:** `GET /api/inventory/warehouse/:id`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
[
  {
    "id": "dd0e8400-e29b-41d4-a716-446655440000",
    "name": "АК-74М",
    "quantity": 45,
    "category_name": "Зброя"
  }
]
```

---

## 📝 Заявки на постачання

### Створити заявку

**Endpoint:** `POST /api/requests`  
**Auth:** ✅ Потрібна  
**Roles:** COMMANDER, WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "resource_id": "dd0e8400-e29b-41d4-a716-446655440000",
  "quantity": 10,
  "priority": "HIGH",
  "notes": "Термінова потреба для операції"
}
```

**Priority values:** `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`

**Response:** `201 Created`
```json
{
  "id": "330e8400-e29b-41d4-a716-446655440000",
  "resource_name": "АК-74М",
  "quantity": 10,
  "status": "PENDING",
  "priority": "HIGH",
  "created_at": "2026-04-24T01:29:00Z"
}
```

---

### Список заявок

**Endpoint:** `GET /api/requests`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Query Parameters:**
- `status` (optional) - PENDING, APPROVED, REJECTED, CANCELLED
- `priority` (optional) - LOW, MEDIUM, HIGH, CRITICAL
- `resource_id` (optional)

**Example:** `GET /api/requests?status=PENDING&priority=HIGH`

**Response:** `200 OK`
```json
[
  {
    "id": "330e8400-e29b-41d4-a716-446655440000",
    "resource_id": "dd0e8400-e29b-41d4-a716-446655440000",
    "resource_name": "АК-74М",
    "requester_name": "Іван Петренко",
    "quantity": 10,
    "status": "PENDING",
    "priority": "HIGH",
    "notes": "Термінова потреба для операції",
    "created_at": "2026-04-24T01:29:00Z",
    "sla_deadline": "2026-04-26T01:29:00Z",
    "is_sla_breached": false
  }
]
```

---

### Отримати заявку за ID

**Endpoint:** `GET /api/requests/:id`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "id": "330e8400-e29b-41d4-a716-446655440000",
  "resource_id": "dd0e8400-e29b-41d4-a716-446655440000",
  "resource_name": "АК-74М",
  "requester_id": "110e8400-e29b-41d4-a716-446655440000",
  "requester_name": "Іван Петренко",
  "quantity": 10,
  "status": "PENDING",
  "priority": "HIGH",
  "notes": "Термінова потреба для операції",
  "approved_by": null,
  "approved_at": null,
  "created_at": "2026-04-24T01:29:00Z",
  "updated_at": "2026-04-24T01:29:00Z",
  "sla_deadline": "2026-04-26T01:29:00Z",
  "is_sla_breached": false
}
```

---

### Затвердити заявку

**Endpoint:** `POST /api/requests/:id/approve`  
**Auth:** ✅ Потрібна  
**Roles:** COMMANDER, WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "notes": "Затверджено, видача через склад №1"
}
```

**Response:** `200 OK`
```json
{
  "message": "Заявку затверджено",
  "status": "APPROVED"
}
```

---

### Відхилити заявку

**Endpoint:** `POST /api/requests/:id/reject`  
**Auth:** ✅ Потрібна  
**Roles:** COMMANDER, WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "reason": "Недостатньо ресурсів на складі"
}
```

**Response:** `200 OK`
```json
{
  "message": "Заявку відхилено",
  "status": "REJECTED"
}
```

---

### Скасувати заявку

**Endpoint:** `POST /api/requests/:id/cancel`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "message": "Заявку скасовано",
  "status": "CANCELLED"
}
```

---

### Smart Dispatch - Попередній перегляд

**Endpoint:** `POST /api/requests/smart-dispatch-preview`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Request:**
```json
{
  "request_ids": [
    "330e8400-e29b-41d4-a716-446655440000",
    "440e8400-e29b-41d4-a716-446655440000"
  ]
}
```

**Response:** `200 OK`
```json
{
  "total_requests": 2,
  "warehouses": [
    {
      "warehouse_id": "ee0e8400-e29b-41d4-a716-446655440000",
      "warehouse_name": "Склад №1",
      "requests_count": 2,
      "distance_km": 15.3,
      "estimated_time_hours": 1.5
    }
  ],
  "total_distance_km": 15.3,
  "estimated_total_time_hours": 1.5
}
```

---

### Smart Dispatch - Підтвердження

**Endpoint:** `POST /api/requests/smart-dispatch-confirm`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Request:**
```json
{
  "request_ids": [
    "330e8400-e29b-41d4-a716-446655440000",
    "440e8400-e29b-41d4-a716-446655440000"
  ],
  "vehicle_id": "110e8400-e29b-41d4-a716-446655440000"
}
```

**Response:** `200 OK`
```json
{
  "message": "Smart Dispatch виконано",
  "shipment_id": "550e8400-e29b-41d4-a716-446655440000",
  "approved_requests": 2
}
```

---

## 🤝 Волонтерські заявки

### Список волонтерських заявок

**Endpoint:** `GET /api/contractor-requests`  
**Auth:** ✅ Потрібна  
**Tier:** BASIC+

**Query Parameters:**
- `status` (optional) - OPEN, IN_PROGRESS, DELIVERED, ACCEPTED, REJECTED, CANCELLED

**Response:** `200 OK`
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "title": "Медикаменти для госпіталю",
    "description": "Потрібні антибіотики, бинти, кровоспинні",
    "category": "medical",
    "quantity": "50 одиниць",
    "priority": "HIGH",
    "status": "OPEN",
    "requester_name": "Майор Іваненко",
    "requester_unit": "26-та Бригада",
    "contractor_id": null,
    "contractor_name": null,
    "created_at": "2026-04-20T10:00:00Z"
  }
]
```

---

### Створити волонтерську заявку

**Endpoint:** `POST /api/contractor-requests`  
**Auth:** ✅ Потрібна  
**Roles:** COMMANDER, WAREHOUSE, TENANT_ADMIN (не CONTRACTOR)  
**Tier:** BASIC+

**Request:**
```json
{
  "title": "Термоодяг для підрозділу",
  "description": "Потрібні зимові куртки та термобілизна",
  "category": "clothing",
  "quantity": "100 комплектів",
  "priority": "MEDIUM",
  "delivery_address": "26-та Бригада, Склад №1",
  "contact_phone": "+380501234567"
}
```

**Categories:** `medical`, `clothing`, `equipment`, `food`, `other`

**Response:** `201 Created`
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440000",
  "status": "OPEN",
  "created_at": "2026-04-24T01:29:00Z"
}
```

---

### Взяти заявку в роботу (Волонтер)

**Endpoint:** `POST /api/contractor-requests/:id/take`  
**Auth:** ✅ Потрібна  
**Roles:** CONTRACTOR  
**Tier:** BASIC+

**Response:** `200 OK`
```json
{
  "message": "Заявку взято в роботу",
  "status": "IN_PROGRESS"
}
```

---

### Позначити як доставлено (Волонтер)

**Endpoint:** `POST /api/contractor-requests/:id/deliver`  
**Auth:** ✅ Потрібна  
**Roles:** CONTRACTOR  
**Tier:** BASIC+

**Request:**
```json
{
  "delivery_notes": "Доставлено в повному обсязі 24.04.2026",
  "delivery_photo_url": "https://example.com/proof.jpg" // опціонально
}
```

**Response:** `200 OK`
```json
{
  "message": "Заявку позначено як доставлено",
  "status": "DELIVERED"
}
```

---

### Прийняти на баланс (Військовий)

**Endpoint:** `POST /api/contractor-requests/:id/accept`  
**Auth:** ✅ Потрібна  
**Roles:** COMMANDER, WAREHOUSE, TENANT_ADMIN  
**Tier:** BASIC+

**Request:**
```json
{
  "notes": "Прийнято на баланс, дякуємо волонтерам!"
}
```

**Response:** `200 OK`
```json
{
  "message": "Заявку прийнято на баланс",
  "status": "ACCEPTED"
}
```

---

### Відхилити доставку (Військовий)

**Endpoint:** `POST /api/contractor-requests/:id/reject`  
**Auth:** ✅ Потрібна  
**Roles:** COMMANDER, WAREHOUSE, TENANT_ADMIN  
**Tier:** BASIC+

**Request:**
```json
{
  "reason": "Неповний комплект, бракує 20 одиниць"
}
```

**Response:** `200 OK`
```json
{
  "message": "Доставку відхилено",
  "status": "REJECTED"
}
```

---

### Скасувати заявку (Військовий)

**Endpoint:** `POST /api/contractor-requests/:id/cancel`  
**Auth:** ✅ Потрібна  
**Roles:** COMMANDER, WAREHOUSE, TENANT_ADMIN  
**Tier:** BASIC+

**Request:**
```json
{
  "reason": "Потреба відпала"
}
```

**Response:** `200 OK`
```json
{
  "message": "Заявку скасовано",
  "status": "CANCELLED"
}
```

---

## 🚙 Автопарк

### Список транспорту

**Endpoint:** `GET /api/vehicles`  
**Auth:** ✅ Потрібна  
**Tier:** BASIC+

**Query Parameters:**
- `status` (optional) - ACTIVE, MAINTENANCE, RETIRED
- `type` (optional) - TRUCK, CAR, ARMORED, OTHER

**Response:** `200 OK`
```json
[
  {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "name": "КамАЗ #123",
    "type": "TRUCK",
    "license_plate": "АА 1234 ВС",
    "vin": "XTA12345678901234",
    "status": "ACTIVE",
    "fuel_type": "DIESEL",
    "fuel_capacity_liters": 350,
    "current_driver_name": "Сергій Коваль",
    "current_mileage_km": 125000,
    "created_at": "2025-01-15T10:00:00Z"
  }
]
```

---

### Створити транспорт

**Endpoint:** `POST /api/vehicles`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** BASIC+

**Request:**
```json
{
  "name": "УАЗ #456",
  "type": "CAR",
  "license_plate": "ВВ 5678 ХТ",
  "vin": "XTA98765432109876",
  "fuel_type": "GASOLINE",
  "fuel_capacity_liters": 70,
  "manufacturer": "УАЗ",
  "model": "Патриот",
  "year": 2020
}
```

**Vehicle Types:** `TRUCK`, `CAR`, `ARMORED`, `OTHER`  
**Fuel Types:** `GASOLINE`, `DIESEL`, `ELECTRIC`, `HYBRID`

**Response:** `201 Created`
```json
{
  "id": "990e8400-e29b-41d4-a716-446655440000",
  "name": "УАЗ #456",
  "status": "ACTIVE",
  "created_at": "2026-04-24T01:29:00Z"
}
```

---

### Отримати транспорт за ID

**Endpoint:** `GET /api/vehicles/:id`  
**Auth:** ✅ Потрібна  
**Tier:** BASIC+

**Response:** `200 OK`
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "name": "КамАЗ #123",
  "type": "TRUCK",
  "license_plate": "АА 1234 ВС",
  "vin": "XTA12345678901234",
  "status": "ACTIVE",
  "fuel_type": "DIESEL",
  "fuel_capacity_liters": 350,
  "current_driver_id": "aa0e8400-e29b-41d4-a716-446655440000",
  "current_driver_name": "Сергій Коваль",
  "current_mileage_km": 125000,
  "manufacturer": "КамАЗ",
  "model": "5350",
  "year": 2018,
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2026-04-20T12:00:00Z"
}
```

---

### Оновити транспорт

**Endpoint:** `PATCH /api/vehicles/:id`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** BASIC+

**Request:**
```json
{
  "current_mileage_km": 126500,
  "status": "MAINTENANCE"
}
```

**Response:** `200 OK`
```json
{
  "message": "Транспорт оновлено"
}
```

---

### Оновити статус транспорту

**Endpoint:** `PATCH /api/vehicles/:id/status`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** BASIC+

**Request:**
```json
{
  "status": "MAINTENANCE",
  "notes": "Планове ТО"
}
```

**Statuses:** `ACTIVE`, `MAINTENANCE`, `RETIRED`

**Response:** `200 OK`
```json
{
  "message": "Статус оновлено"
}
```

---

### Видалити транспорт

**Endpoint:** `DELETE /api/vehicles/:id`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** BASIC+

**Response:** `200 OK`
```json
{
  "message": "Транспорт видалено"
}
```

---

### Призначити водія

**Endpoint:** `PATCH /api/vehicles/:id/driver`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** BASIC+

**Request:**
```json
{
  "driver_id": "bb0e8400-e29b-41d4-a716-446655440000"
}
```

**Response:** `200 OK`
```json
{
  "message": "Водія призначено"
}
```

---

### Історія водіїв

**Endpoint:** `GET /api/vehicles/:id/drivers`  
**Auth:** ✅ Потрібна  
**Tier:** BASIC+

**Response:** `200 OK`
```json
[
  {
    "driver_id": "aa0e8400-e29b-41d4-a716-446655440000",
    "driver_name": "Сергій Коваль",
    "assigned_at": "2026-01-10T09:00:00Z",
    "unassigned_at": null
  },
  {
    "driver_id": "bb0e8400-e29b-41d4-a716-446655440000",
    "driver_name": "Петро Шевченко",
    "assigned_at": "2025-06-15T10:00:00Z",
    "unassigned_at": "2026-01-10T09:00:00Z"
  }
]
```

---

### Виконати технічне обслуговування

**Endpoint:** `POST /api/vehicles/:id/maintenance`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** BASIC+

**Request:** (multipart/form-data)
```
type: OIL_CHANGE
description: Заміна моторного масла
cost: 1500.00
mileage_km: 126500
performed_by: Механік Іванов
files: [photo1.jpg, photo2.jpg] // опціонально
```

**Maintenance Types:**
- `OIL_CHANGE` - Заміна масла
- `TIRE_ROTATION` - Ротація шин
- `FILTER_REPLACEMENT` - Заміна фільтрів
- `INSPECTION` - Інспекція
- `REPAIR` - Ремонт
- `OTHER` - Інше

**Response:** `201 Created`
```json
{
  "id": "cc0e8400-e29b-41d4-a716-446655440000",
  "message": "ТО зареєстровано",
  "performed_at": "2026-04-24T01:29:00Z"
}
```

---

### Історія ТО

**Endpoint:** `GET /api/vehicles/:id/maintenance`  
**Auth:** ✅ Потрібна  
**Tier:** BASIC+

**Response:** `200 OK`
```json
[
  {
    "id": "cc0e8400-e29b-41d4-a716-446655440000",
    "type": "OIL_CHANGE",
    "description": "Заміна моторного масла",
    "cost": 1500.00,
    "mileage_km": 126500,
    "performed_by": "Механік Іванов",
    "performed_at": "2026-04-24T01:29:00Z",
    "attachments": [
      "/uploads/maintenance/photo1.jpg"
    ]
  }
]
```

---

### Доступний транспорт для маршруту

**Endpoint:** `GET /api/vehicles/available-for-route`  
**Auth:** ✅ Потрібна  
**Tier:** BASIC+

**Response:** `200 OK`
```json
[
  {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "name": "КамАЗ #123",
    "type": "TRUCK",
    "status": "ACTIVE",
    "current_driver_name": "Сергій Коваль"
  }
]
```

---

## ⛽ Паливо

### Додати запис про заправку

**Endpoint:** `POST /api/vehicles/:id/fuel`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN, COMMANDER  
**Tier:** BASIC+

**Request:**
```json
{
  "liters": 180,
  "price_per_liter": 45.50,
  "mileage_km": 125500,
  "location": "АЗС Київ",
  "notes": "Повна заправка"
}
```

**Response:** `201 Created`
```json
{
  "id": "dd0e8400-e29b-41d4-a716-446655440000",
  "total_cost": 8190.00,
  "recorded_at": "2026-04-24T01:29:00Z"
}
```

---

### Історія заправок

**Endpoint:** `GET /api/vehicles/:id/fuel`  
**Auth:** ✅ Потрібна  
**Tier:** BASIC+

**Query Parameters:**
- `from_date` (optional) - YYYY-MM-DD
- `to_date` (optional) - YYYY-MM-DD

**Example:** `GET /api/vehicles/:id/fuel?from_date=2026-01-01&to_date=2026-04-24`

**Response:** `200 OK`
```json
[
  {
    "id": "dd0e8400-e29b-41d4-a716-446655440000",
    "liters": 180,
    "price_per_liter": 45.50,
    "total_cost": 8190.00,
    "mileage_km": 125500,
    "location": "АЗС Київ",
    "recorded_by_name": "Сергій Коваль",
    "recorded_at": "2026-04-24T01:29:00Z",
    "consumption_per_100km": 28.5, // розраховується автоматично
    "notes": "Повна заправка"
  }
]
```

---

## 🏢 Склади

### Список складів

**Endpoint:** `GET /api/warehouses`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
[
  {
    "id": "ee0e8400-e29b-41d4-a716-446655440000",
    "name": "Склад №1",
    "location": "Київ, вул. Військова 15",
    "unit_id": "330e8400-e29b-41d4-a716-446655440000",
    "unit_name": "1-й Полк",
    "manager_name": "Петро Складський",
    "capacity": 1000,
    "current_load": 650,
    "created_at": "2025-01-10T10:00:00Z"
  }
]
```

---

### Створити склад

**Endpoint:** `POST /api/warehouses`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "name": "Склад №2",
  "location": "Львів, вул. Героїв 20",
  "unit_id": "440e8400-e29b-41d4-a716-446655440000",
  "capacity": 500
}
```

**Response:** `201 Created`
```json
{
  "id": "ff0e8400-e29b-41d4-a716-446655440000",
  "name": "Склад №2",
  "created_at": "2026-04-24T01:29:00Z"
}
```

---

### Оновити локацію складу

**Endpoint:** `PATCH /api/warehouses/:id/location`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "location": "Львів, нова адреса вул. Перемоги 5"
}
```

**Response:** `200 OK`
```json
{
  "message": "Локацію оновлено"
}
```

---

### Оновити склад

**Endpoint:** `PATCH /api/warehouses/:id`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Request:**
```json
{
  "name": "Центральний склад",
  "capacity": 1500
}
```

**Response:** `200 OK`
```json
{
  "message": "Склад оновлено"
}
```

---

### Видалити склад

**Endpoint:** `DELETE /api/warehouses/:id`  
**Auth:** ✅ Потрібна  
**Roles:** WAREHOUSE, TENANT_ADMIN  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "message": "Склад видалено"
}
```

**Errors:**
- `400 Bad Request` - Склад містить ресурси

---

## 📊 Аналітика

### Dashboard (базові метрики)

**Endpoint:** `GET /api/analytics/dashboard`  
**Auth:** ✅ Потрібна  
**Tier:** Всі

**Response:** `200 OK`
```json
{
  "total_resources": 234,
  "critical_resources": 12,
  "pending_requests": 5,
  "total_vehicles": 8,
  "active_vehicles": 6,
  "sla_compliance": 94.5,
  "total_cost_last_month": 125000.00
}
```

---

### Auto Replenish (авто-поповнення)

**Endpoint:** `POST /api/analytics/auto-replenish`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Request:**
```json
{
  "threshold_percentage": 80
}
```

**Response:** `200 OK`
```json
{
  "message": "Створено 3 автоматичні заявки",
  "requests_created": 3,
  "resources": [
    {
      "resource_name": "АК-74М",
      "current": 45,
      "min": 50,
      "requested": 10
    }
  ]
}
```

---

### Advanced KPI Dashboard

**Endpoint:** `GET /api/analytics/kpi`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Response:** `200 OK`
```json
{
  "sla_percentage": 94.5,
  "total_cost_of_ownership": 2500000.00,
  "risk_score": 32,
  "depletion_forecast_days": 45,
  "inventory_turnover_ratio": 3.2,
  "fulfillment_rate": 97.8,
  "average_response_time_hours": 4.5
}
```

---

### Demand Forecasting (AI прогноз)

**Endpoint:** `GET /api/analytics/forecast`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Query Parameters:**
- `resource_id` (optional) - прогноз для конкретного ресурсу
- `days` (optional, default: 30) - період прогнозу

**Example:** `GET /api/analytics/forecast?resource_id=dd0e8400-e29b-41d4-a716-446655440000&days=60`

**Response:** `200 OK`
```json
{
  "resource_id": "dd0e8400-e29b-41d4-a716-446655440000",
  "resource_name": "АК-74М",
  "current_quantity": 45,
  "forecast_period_days": 60,
  "scenarios": {
    "low_demand": {
      "daily_consumption": 0.5,
      "depletion_date": "2026-08-22",
      "recommended_order_quantity": 15
    },
    "medium_demand": {
      "daily_consumption": 0.75,
      "depletion_date": "2026-07-23",
      "recommended_order_quantity": 30
    },
    "high_demand": {
      "daily_consumption": 1.2,
      "depletion_date": "2026-06-13",
      "recommended_order_quantity": 50
    }
  },
  "confidence": 87.5,
  "trend": "STABLE"
}
```

---

### Predictive Maintenance Schedule

**Endpoint:** `GET /api/analytics/maintenance`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Query Parameters:**
- `vehicle_id` (optional) - графік для конкретного транспорту
- `priority` (optional) - LOW, MEDIUM, HIGH

**Response:** `200 OK`
```json
[
  {
    "vehicle_id": "880e8400-e29b-41d4-a716-446655440000",
    "vehicle_name": "КамАЗ #123",
    "maintenance_type": "OIL_CHANGE",
    "current_mileage_km": 125000,
    "next_maintenance_mileage_km": 130000,
    "miles_remaining": 5000,
    "estimated_days": 15,
    "priority": "MEDIUM",
    "estimated_cost": 1500.00
  },
  {
    "vehicle_id": "990e8400-e29b-41d4-a716-446655440000",
    "vehicle_name": "УАЗ #456",
    "maintenance_type": "TIRE_ROTATION",
    "current_mileage_km": 48000,
    "next_maintenance_mileage_km": 50000,
    "miles_remaining": 2000,
    "estimated_days": 6,
    "priority": "HIGH",
    "estimated_cost": 800.00
  }
]
```

---

### Fuel Anti-Fraud Detection

**Endpoint:** `GET /api/analytics/fuel-anomalies`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Query Parameters:**
- `vehicle_id` (optional)
- `from_date` (optional) - YYYY-MM-DD
- `to_date` (optional) - YYYY-MM-DD

**Response:** `200 OK`
```json
[
  {
    "id": "110e8400-e29b-41d4-a716-446655440000",
    "vehicle_id": "880e8400-e29b-41d4-a716-446655440000",
    "vehicle_name": "КамАЗ #123",
    "anomaly_type": "EXTREME_REFILL",
    "severity": "HIGH",
    "risk_score": 85,
    "details": {
      "liters": 500,
      "tank_capacity": 350,
      "overfill_liters": 150
    },
    "description": "Заправлено більше ніж місткість баку",
    "recorded_at": "2026-04-20T14:30:00Z",
    "estimated_loss": 6825.00,
    "investigation_status": "PENDING"
  },
  {
    "id": "220e8400-e29b-41d4-a716-446655440000",
    "vehicle_id": "990e8400-e29b-41d4-a716-446655440000",
    "vehicle_name": "УАЗ #456",
    "anomaly_type": "FREQUENT_SMALL_REFILLS",
    "severity": "MEDIUM",
    "risk_score": 65,
    "details": {
      "refills_count": 5,
      "average_liters": 12,
      "period_days": 3
    },
    "description": "Підозріла серія малих заправок",
    "recorded_at": "2026-04-22T10:00:00Z",
    "estimated_loss": 2000.00,
    "investigation_status": "PENDING"
  }
]
```

**Anomaly Types:**
- `EXTREME_REFILL` - Перевищення місткості баку
- `FREQUENT_SMALL_REFILLS` - Серія підозрілих малих заправок
- `PRICE_ANOMALY` - Аномальна ціна за літр
- `ABNORMAL_CONSUMPTION` - Аномальне споживання палива

---

### Експорт інвентаря (Excel)

**Endpoint:** `GET /api/analytics/export/inventory`  
**Auth:** ✅ Потрібна  
**Tier:** STANDARD+

**Query Parameters:**
- `format` (optional, default: xlsx) - xlsx, csv
- `category_id` (optional)
- `warehouse_id` (optional)

**Response:** `200 OK` (Excel file)
```
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="inventory_export_2026-04-24.xlsx"
```

---

### Експорт палива (Excel)

**Endpoint:** `GET /api/analytics/export/fuel`  
**Auth:** ✅ Потрібна  
**Tier:** STANDARD+

**Query Parameters:**
- `format` (optional, default: xlsx) - xlsx, csv
- `vehicle_id` (optional)
- `from_date` (optional) - YYYY-MM-DD
- `to_date` (optional) - YYYY-MM-DD

**Response:** `200 OK` (Excel file)
```
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="fuel_export_2026-04-24.xlsx"
```

---

## 🌍 GPS Трекінг (PRO)

### Записати локацію транспорту

**Endpoint:** `POST /api/gps/locations`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Request:**
```json
{
  "vehicle_id": "880e8400-e29b-41d4-a716-446655440000",
  "latitude": 50.4501,
  "longitude": 30.5234,
  "speed_kmh": 65,
  "heading_degrees": 180,
  "altitude_meters": 150
}
```

**Response:** `201 Created`
```json
{
  "id": "330e8400-e29b-41d4-a716-446655440000",
  "recorded_at": "2026-04-24T01:29:00Z"
}
```

---

### Карта флоту (всі транспорти)

**Endpoint:** `GET /api/gps/fleet-map`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Response:** `200 OK`
```json
[
  {
    "vehicle_id": "880e8400-e29b-41d4-a716-446655440000",
    "vehicle_name": "КамАЗ #123",
    "vehicle_type": "TRUCK",
    "last_location": {
      "latitude": 50.4501,
      "longitude": 30.5234,
      "speed_kmh": 65,
      "heading_degrees": 180,
      "recorded_at": "2026-04-24T01:29:00Z"
    },
    "status": "MOVING",
    "driver_name": "Сергій Коваль"
  }
]
```

---

### Траєкторія транспорту

**Endpoint:** `GET /api/gps/trajectory`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Query Parameters:**
- `vehicle_id` (required)
- `from_date` (required) - YYYY-MM-DD
- `to_date` (required) - YYYY-MM-DD

**Example:** `GET /api/gps/trajectory?vehicle_id=880e8400-e29b-41d4-a716-446655440000&from_date=2026-04-23&to_date=2026-04-24`

**Response:** `200 OK`
```json
{
  "vehicle_id": "880e8400-e29b-41d4-a716-446655440000",
  "vehicle_name": "КамАЗ #123",
  "period": {
    "from": "2026-04-23T00:00:00Z",
    "to": "2026-04-24T23:59:59Z"
  },
  "points": [
    {
      "latitude": 50.4501,
      "longitude": 30.5234,
      "speed_kmh": 65,
      "recorded_at": "2026-04-24T08:15:00Z"
    },
    {
      "latitude": 50.4520,
      "longitude": 30.5280,
      "speed_kmh": 70,
      "recorded_at": "2026-04-24T08:20:00Z"
    }
  ],
  "total_distance_km": 145.3,
  "average_speed_kmh": 58.5
}
```

---

### Створити геозону

**Endpoint:** `POST /api/gps/geofences`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Request:**
```json
{
  "name": "Зона бази",
  "center_latitude": 50.4501,
  "center_longitude": 30.5234,
  "radius_meters": 500,
  "type": "BASE",
  "alert_on_enter": true,
  "alert_on_exit": true
}
```

**Geofence Types:**
- `BASE` - База/склад
- `RESTRICTED` - Заборонена зона
- `SAFE_ZONE` - Безпечна зона
- `CHECKPOINT` - Контрольна точка

**Response:** `201 Created`
```json
{
  "id": "440e8400-e29b-41d4-a716-446655440000",
  "name": "Зона бази",
  "created_at": "2026-04-24T01:29:00Z"
}
```

---

### Список геозон

**Endpoint:** `GET /api/gps/geofences`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Response:** `200 OK`
```json
[
  {
    "id": "440e8400-e29b-41d4-a716-446655440000",
    "name": "Зона бази",
    "center_latitude": 50.4501,
    "center_longitude": 30.5234,
    "radius_meters": 500,
    "type": "BASE",
    "alert_on_enter": true,
    "alert_on_exit": true,
    "active": true,
    "created_at": "2026-04-24T01:29:00Z"
  }
]
```

---

### Алерти геозон

**Endpoint:** `GET /api/gps/geofence-alerts`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Query Parameters:**
- `vehicle_id` (optional)
- `geofence_id` (optional)
- `from_date` (optional) - YYYY-MM-DD
- `severity` (optional) - LOW, MEDIUM, HIGH

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "vehicle_id": "880e8400-e29b-41d4-a716-446655440000",
    "vehicle_name": "КамАЗ #123",
    "geofence_id": "440e8400-e29b-41d4-a716-446655440000",
    "geofence_name": "Зона бази",
    "alert_type": "EXIT",
    "severity": "MEDIUM",
    "location": {
      "latitude": 50.4550,
      "longitude": 30.5280
    },
    "distance_from_center_meters": 520,
    "triggered_at": "2026-04-24T10:15:00Z",
    "acknowledged": false
  }
]
```

---

### Статус флоту

**Endpoint:** `GET /api/gps/fleet-status`  
**Auth:** ✅ Потрібна  
**Tier:** 🌟 **PRO**

**Response:** `200 OK`
```json
{
  "total_vehicles": 8,
  "active_vehicles": 6,
  "statuses": {
    "MOVING": 3,
    "IDLE": 2,
    "PARKED": 1,
    "UNKNOWN": 2
  },
  "in_geofence": 4,
  "out_of_geofence": 2,
  "recent_alerts_count": 1
}
```

---

## 🌐 Platform Admin

### Статистика платформи

**Endpoint:** `GET /api/platform/stats`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN  
**Tier:** N/A

**Response:** `200 OK`
```json
{
  "total_tenants": 45,
  "active_tenants": 42,
  "total_users": 1250,
  "tier_distribution": {
    "FREE": 20,
    "BASIC": 15,
    "STANDARD": 8,
    "PRO": 2
  },
  "total_resources": 15000,
  "total_requests": 8500,
  "total_vehicles": 340
}
```

---

### Список організацій (tenants)

**Endpoint:** `GET /api/platform/tenants`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN  
**Tier:** N/A

**Query Parameters:**
- `tier` (optional) - FREE, BASIC, STANDARD, PRO
- `active` (optional) - true, false

**Response:** `200 OK`
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "name": "26-та Окрема Бригада",
    "subscription_tier": "PRO",
    "active": true,
    "users_count": 45,
    "resources_count": 340,
    "created_at": "2025-06-15T10:00:00Z",
    "tier_updated_at": "2026-01-10T12:00:00Z"
  }
]
```

---

### Отримати організацію за ID

**Endpoint:** `GET /api/platform/tenants/:id`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN  
**Tier:** N/A

**Response:** `200 OK`
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "name": "26-та Окрема Бригада",
  "subscription_tier": "PRO",
  "active": true,
  "users_count": 45,
  "resources_count": 340,
  "vehicles_count": 12,
  "warehouses_count": 3,
  "pending_requests_count": 5,
  "created_at": "2025-06-15T10:00:00Z",
  "tier_updated_at": "2026-01-10T12:00:00Z",
  "owner": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "name": "Полковник Іваненко",
    "email": "admin@brigade26.mil"
  }
}
```

---

### Оновити тариф організації

**Endpoint:** `PATCH /api/platform/tenants/:id/tier`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN  
**Tier:** N/A

**Request:**
```json
{
  "subscription_tier": "PRO"
}
```

**Tiers:** `FREE`, `BASIC`, `STANDARD`, `PRO`

**Response:** `200 OK`
```json
{
  "message": "Тариф оновлено",
  "new_tier": "PRO"
}
```

---

### Активувати/деактивувати організацію

**Endpoint:** `PATCH /api/platform/tenants/:id/active`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN  
**Tier:** N/A

**Request:**
```json
{
  "active": false
}
```

**Response:** `200 OK`
```json
{
  "message": "Статус організації оновлено",
  "active": false
}
```

---

### Видалити організацію

**Endpoint:** `DELETE /api/platform/tenants/:id`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN  
**Tier:** N/A

**Response:** `200 OK`
```json
{
  "message": "Організацію видалено"
}
```

**Warning:** Видаляє всі дані організації (користувачів, ресурси, заявки, тощо). Незворотна операція!

---

## 🔍 Audit Logs

### Отримати логи аудиту

**Endpoint:** `GET /api/admin/audit-logs`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN, TENANT_ADMIN  
**Tier:** BASIC+

**Query Parameters:**
- `entity_type` (optional) - resource, request, vehicle, user, etc.
- `action` (optional) - CREATE, UPDATE, DELETE, APPROVE, REJECT, etc.
- `user_id` (optional)
- `from_date` (optional) - YYYY-MM-DD
- `to_date` (optional) - YYYY-MM-DD
- `limit` (optional, default: 100)

**Example:** `GET /api/admin/audit-logs?entity_type=resource&action=DELETE&limit=50`

**Response:** `200 OK`
```json
[
  {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "user_id": "990e8400-e29b-41d4-a716-446655440000",
    "user_name": "Петро Складський",
    "user_email": "warehouse@brigade26.mil",
    "action": "DELETE",
    "entity_type": "resource",
    "entity_id": "aa0e8400-e29b-41d4-a716-446655440000",
    "entity_name": "Старі каски",
    "details": {
      "reason": "Застарілі, списано"
    },
    "ip_address": "192.168.1.10",
    "user_agent": "Mozilla/5.0...",
    "created_at": "2026-04-24T01:29:00Z"
  }
]
```

---

### Trigger SLA Check (manual)

**Endpoint:** `POST /api/admin/sla/trigger`  
**Auth:** ✅ Потрібна  
**Roles:** SYSTEM_ADMIN, TENANT_ADMIN  
**Tier:** BASIC+

**Response:** `200 OK`
```json
{
  "message": "SLA перевірку запущено",
  "breached_requests": 2,
  "notifications_sent": 2
}
```

---

## 🔒 Помилки (Error Codes)

### Стандартні HTTP коди

| Код | Значення | Опис |
|-----|----------|------|
| `200` | OK | Успішний запит |
| `201` | Created | Ресурс створено |
| `400` | Bad Request | Невалідні дані |
| `401` | Unauthorized | Не авторизовано (немає токена) |
| `403` | Forbidden | Немає прав доступу |
| `404` | Not Found | Ресурс не знайдено |
| `409` | Conflict | Конфлікт (наприклад, email вже існує) |
| `422` | Unprocessable Entity | Бізнес-логічна помилка |
| `429` | Too Many Requests | Rate limit |
| `500` | Internal Server Error | Серверна помилка |
| `503` | Service Unavailable | Сервіс недоступний |

### Формат помилки

```json
{
  "error": "Unauthorized",
  "message": "JWT token is missing or invalid",
  "code": "AUTH_001",
  "timestamp": "2026-04-24T01:29:00Z"
}
```

### Типові помилки

**Auth Errors:**
- `AUTH_001` - Token missing or invalid
- `AUTH_002` - Token expired
- `AUTH_003` - Insufficient permissions
- `AUTH_004` - Account blocked

**Validation Errors:**
- `VAL_001` - Required field missing
- `VAL_002` - Invalid format
- `VAL_003` - Value out of range

**Business Logic Errors:**
- `BIZ_001` - Insufficient quantity
- `BIZ_002` - Resource already assigned
- `BIZ_003` - Subscription tier required
- `BIZ_004` - Limit reached (max users/resources/vehicles)

---

## 📌 Rate Limiting

**Ліміти:**
- Неавторизовані запити: 100 req/min
- Авторизовані запити: 1000 req/min
- Heavy operations (imports, exports): 10 req/min

**Headers:**
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 950
X-RateLimit-Reset: 1714780140
```

---

## 🔧 Webhook Events (Future)

Система підтримуватиме webhooks для наступних подій:

- `request.created` - Нова заявка
- `request.approved` - Заявка затверджена
- `request.rejected` - Заявка відхилена
- `resource.critical` - Критичний залишок ресурсу
- `vehicle.maintenance_due` - Час для ТО
- `geofence.alert` - Порушення геозони
- `sla.breached` - Порушення SLA

---

## 📞 Підтримка

- 📧 Email: support@omnilog.system
- 📖 Документація: https://docs.omnilog.system
- 🐛 Bug Reports: https://github.com/markostrutinsky/diploma/issues

---

<div align="center">

**Made with ❤️ in Ukraine 🇺🇦**

</div>
