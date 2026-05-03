# 🎭 Інструкція з тестування API за ролями

## 📊 Інформація про тенант

**ID тенанта:** `8cf55fa1-69a0-4a2c-b3a1-50143ec86428`  
**Назва:** Пійлівський Логістичний Центр  
**Slug:** `piylo`  
**Тарифний план:** `BASIC`  
**Власник:** markostrutinsky@gmail.com  
**Статус:** Активний

### 📈 Поточний стан даних

- **Користувачів:** 23
- **Підрозділів:** 31 (3 регіони, 5 філій, 11 відділів, 12 команд)
- **Складів:** 10/10 (ДОСЯГНУТО ЛІМІТУ BASIC)
- **Транспорту:** 5 одиниць
- **Категорій ресурсів:** 4
- **Ресурсів:** 4
- **Заявок на постачання:** 0
- **Волонтерських заявок:** 0

---

## 🔑 Тестові облікові записи за ролями

**Пароль для всіх користувачів:** `password123`

### 1️⃣ TENANT_ADMIN (Адміністратор тенанта)
```
Email: markostrutinsky@gmail.com
Username: markostrutinsky
Роль: TENANT_ADMIN
Повноваження: Найвищі права в межах тенанта
```

### 2️⃣ REGION_DIRECTOR (Директор регіону)
```
Email: d.tkachenko@logistics.ua
Username: d.tkachenko
ПІБ: Ткаченко Дмитро Олегович
Підрозділ: Центральний регіон управління постачанням (ID: 2)
```

### 3️⃣ REGION_LOGISTICIAN (Логіст регіону)
```
Email: i.melnyk@logistics.ua
Username: i.melnyk
ПІБ: Мельник Ігор Володимирович
Підрозділ: Центральний регіон управління постачанням (ID: 2)
```

### 4️⃣ REGION_STOREKEEPER (Комірник регіону)
```
Email: m.honcharenko@logistics.ua
Username: m.honcharenko
ПІБ: Гончаренко Михайло Сергійович
Підрозділ: Центральний регіон управління постачанням (ID: 2)
```

### 5️⃣ BRANCH_MANAGER (Менеджер філії)
```
Email: n.boiko@logistics.ua
Username: n.boiko
ПІБ: Бойко Наталія Сергіївна
Підрозділ: Київський головний розподільчий центр (ID: 3)
```

### 6️⃣ BRANCH_LOGISTICIAN (Логіст філії)
```
Email: m.hrytsenko@logistics.ua
Username: m.hrytsenko
ПІБ: Гриценко Максим Іванович
Підрозділ: Київський головний розподільчий центр (ID: 3)
```

### 7️⃣ BRANCH_STOREKEEPER (Комірник філії)
```
Email: o.lysenko@logistics.ua
Username: o.lysenko
ПІБ: Лисенко Олександр Петрович
Підрозділ: Київський головний розподільчий центр (ID: 3)
```

### 8️⃣ DEPT_MANAGER (Начальник відділу)
```
Email: y.kravchenko@logistics.ua
Username: y.kravchenko
ПІБ: Кравченко Юлія Анатоліївна
Підрозділ: Диспетчерський департамент (ID: 7)
```

### 9️⃣ DEPT_SUPERVISOR (Супервайзер відділу)
```
Email: s.marchenko@logistics.ua
Username: s.marchenko
ПІБ: Марченко Сергій Олександрович
Підрозділ: Диспетчерський департамент (ID: 7)
```

### 🔟 TEAM_LEAD (Керівник команди)
```
Email: o.zakharchenko@logistics.ua
Username: o.zakharchenko
ПІБ: Захарченко Ольга Павлівна
Підрозділ: Команда моніторингу GPS (ID: 8)
```

### 1️⃣1️⃣ EMPLOYEE (Працівник)
```
Email: o.kovalenko@logistics.ua
Username: o.kovalenko
ПІБ: Коваленко Олена Вікторівна
Підрозділ: Зміна денних комплектувальників (ID: 5)
```

---

## 🌐 Базова URL
```
http://localhost/api
```

---

## 📋 Матриця доступу за ролями

| Ендпоінт | TENANT_ADMIN | REGION_DIRECTOR | REGION_LOGISTICIAN | REGION_STOREKEEPER | BRANCH_MANAGER | BRANCH_LOGISTICIAN | BRANCH_STOREKEEPER | DEPT_MANAGER | DEPT_SUPERVISOR | TEAM_LEAD | EMPLOYEE |
|----------|--------------|-----------------|--------------------|--------------------|----------------|--------------------|--------------------|--------------|-----------------|-----------|----------|
| **AUTH** |
| POST /auth/login | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| POST /auth/refresh | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| GET /auth/me | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **USERS** |
| GET /users/commanders | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| GET /users/visible | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| GET /users/limits | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| PUT /users/:id/role | ✅ | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| PUT /users/:id/block | ✅ | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| PUT /users/:id/unblock | ✅ | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| PATCH /users/profile | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| PATCH /users/password | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **ADMIN** |
| POST /admin/users | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| GET /admin/audit-logs | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| POST /admin/sla/trigger | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| **UNITS** |
| GET /units | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| GET /units/available | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| GET /units/my-hierarchy | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| POST /units | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| POST /units/:id/change-commander | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| PATCH /units/:id | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| DELETE /units/:id | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| **INVENTORY** |
| GET /inventory/categories | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| GET /inventory/resources | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| GET /inventory/resources/:id | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| POST /inventory/categories | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| POST /inventory/resources | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| PATCH /inventory/resources/:id | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| DELETE /inventory/resources/:id | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| POST /inventory/resources/:id/assign | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| GET /inventory/my-equipment | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| POST /inventory/shipments | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| GET /inventory/shipments | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **REQUESTS** |
| POST /requests | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| GET /requests | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| POST /requests/:id/approve | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| POST /requests/:id/reject | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **VEHICLES** |
| POST /vehicles | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| GET /vehicles | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| PATCH /vehicles/:id | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| POST /vehicles/:id/fuel | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **WAREHOUSES** |
| GET /warehouses | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| POST /warehouses | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| PATCH /warehouses/:id | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| DELETE /warehouses/:id | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| **ANALYTICS** |
| GET /analytics/dashboard | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| GET /analytics/export/* | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

---

## 🧪 Сценарії тестування за ролями

### 🔐 Загальна підготовка

#### Отримання токена для будь-якої ролі:
```bash
# Заміни email на потрібну роль
TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "email@logistics.ua",
    "password": "password123"
  }' | jq -r '.access_token')

echo $TOKEN
```

---

## 1️⃣ TENANT_ADMIN - Повне тестування

**Користувач:** `markostrutinsky@gmail.com`

### 1.1. Аутентифікація
```bash
# Вхід
curl -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "markostrutinsky@gmail.com",
    "password": "password123"
  }' | jq .

# Зберегти токен
ADMIN_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"markostrutinsky@gmail.com","password":"password123"}' \
  | jq -r '.access_token')

# Перевірка профілю
curl -X GET http://localhost/api/auth/me \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
```

### 1.2. Управління користувачами
```bash
# Список командирів
curl -X GET http://localhost/api/users/commanders \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Видимі користувачі
curl -X GET http://localhost/api/users/visible \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Ліміти користувача
curl -X GET http://localhost/api/users/limits \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Створення нового користувача
curl -X POST http://localhost/api/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@logistics.ua",
    "full_name": "Новий Працівник",
    "role": "EMPLOYEE",
    "unit_id": 5,
    "phone": "+380501234567"
  }' | jq .
```

### 1.3. Управління підрозділами
```bash
# Список підрозділів
curl -X GET http://localhost/api/units \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Доступні підрозділи
curl -X GET http://localhost/api/units/available \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Моя ієрархія
curl -X GET http://localhost/api/units/my-hierarchy \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Створення нового підрозділу (НЕ МОЖНА - досягнуто ліміту BASIC)
curl -X POST http://localhost/api/units \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Новий відділ",
    "unit_type": "DEPARTMENT",
    "parent_id": 3
  }' | jq .
```

### 1.4. Управління інвентарем
```bash
# Категорії
curl -X GET http://localhost/api/inventory/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Створення категорії
curl -X POST http://localhost/api/inventory/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Офісна техніка",
    "description": "Принтери, сканери, МФУ"
  }' | jq .

# Список ресурсів
curl -X GET http://localhost/api/inventory/resources \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Створення ресурсу (потрібен warehouse_id)
# Спочатку отримуємо склади
curl -X GET http://localhost/api/warehouses \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.[] | {id, name, unit_id}'

# Потім створюємо ресурс (підставте реальний warehouse_id та category_id)
curl -X POST http://localhost/api/inventory/resources \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": "CATEGORY_UUID",
    "warehouse_id": "WAREHOUSE_UUID",
    "unit_id": 3,
    "name": "Ноутбук Dell Latitude 5520",
    "description": "i5-1145G7, 16GB RAM, 512GB SSD",
    "quantity": 10,
    "min_quantity": 2,
    "unit_type": "PCS",
    "serial_number": "DELL-2024-001",
    "condition": "NEW",
    "weight_kg": 1.8
  }' | jq .
```

### 1.5. Управління складами
```bash
# Список складів
curl -X GET http://localhost/api/warehouses \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Створення нового складу (НЕ МОЖНА - досягнуто ліміту 10/10)
curl -X POST http://localhost/api/warehouses \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "unit_id": 3,
    "name": "Тестовий склад",
    "location_type": "STATIONARY",
    "latitude": 50.4501,
    "longitude": 30.5234,
    "address": "вул. Тестова, 1"
  }' | jq .
# Очікується помилка: 422 - досягнуто ліміту складів для тарифу BASIC
```

### 1.6. Управління транспортом
```bash
# Список транспорту
curl -X GET http://localhost/api/vehicles \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Створення нового транспорту
curl -X POST http://localhost/api/vehicles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "brand": "Ford",
    "model": "Transit 350",
    "plate_number": "XY9876ZZ",
    "type": "VAN",
    "capacity_kg": 1500,
    "tank_capacity": 80,
    "fuel_norm": 11.0
  }' | jq .

# Отримання транспорту за ID
VEHICLE_ID=$(curl -s -X GET http://localhost/api/vehicles \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.[0].id')

curl -X GET http://localhost/api/vehicles/$VEHICLE_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Додавання запису про пальне
curl -X POST http://localhost/api/vehicles/$VEHICLE_ID/fuel \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "liters": 60.5,
    "odometer_km": 125000,
    "record_type": "REFUEL"
  }' | jq .

# Історія пального
curl -X GET http://localhost/api/vehicles/$VEHICLE_ID/fuel \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
```

### 1.7. Заявки на постачання
```bash
# Список заявок
curl -X GET http://localhost/api/requests \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Створення заявки (потрібен resource_id)
RESOURCE_ID=$(curl -s -X GET http://localhost/api/inventory/resources \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.[0].id')

curl -X POST http://localhost/api/requests \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"resource_id\": \"$RESOURCE_ID\",
    \"quantity\": 5,
    \"comment\": \"Терміново необхідно для проєкту\"
  }" | jq .

# Затвердження заявки
REQUEST_ID=$(curl -s -X GET http://localhost/api/requests \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.[0].id')

curl -X POST http://localhost/api/requests/$REQUEST_ID/approve \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"comment": "Затверджено"}' | jq .
```

### 1.8. Аналітика
```bash
# Дашборд
curl -X GET http://localhost/api/analytics/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Експорт інвентаря
curl -X GET http://localhost/api/analytics/export/inventory \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  --output inventory-export.xlsx

# Експорт пального
curl -X GET http://localhost/api/analytics/export/fuel \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  --output fuel-export.xlsx
```

### 1.9. PRO фічі (НЕ ДОСТУПНІ для BASIC)
```bash
# GPS - має повернути 402 Payment Required
curl -X GET http://localhost/api/gps/fleet-map \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Smart Dispatch - має повернути 402
curl -X POST http://localhost/api/requests/smart-dispatch-preview \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"request_ids": []}' | jq .

# Excel Import - має повернути 402
curl -X POST http://localhost/api/inventory/resources/import \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -F "file=@test.xlsx" | jq .

# Advanced KPI - має повернути 402
curl -X GET http://localhost/api/analytics/kpi \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
```

---

## 2️⃣ REGION_DIRECTOR - Управління регіоном

**Користувач:** `d.tkachenko@logistics.ua`

### 2.1. Отримання токена
```bash
DIRECTOR_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"d.tkachenko@logistics.ua","password":"password123"}' \
  | jq -r '.access_token')
```

### 2.2. Що може директор регіону

#### ✅ ДОЗВОЛЕНО:
```bash
# Створення користувачів у своєму регіоні
curl -X POST http://localhost/api/admin/users \
  -H "Authorization: Bearer $DIRECTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "regional_employee",
    "email": "regional_employee@logistics.ua",
    "full_name": "Регіональний Працівник",
    "role": "EMPLOYEE",
    "unit_id": 3
  }' | jq .

# Зміна ролі користувачів
curl -X PUT http://localhost/api/users/USER_ID/role \
  -H "Authorization: Bearer $DIRECTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "TEAM_LEAD", "unit_id": 5}' | jq .

# Управління підрозділами
curl -X POST http://localhost/api/units \
  -H "Authorization: Bearer $DIRECTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Новий департамент",
    "unit_type": "DEPARTMENT",
    "parent_id": 3
  }' | jq .

# Створення та управління заявками
curl -X POST http://localhost/api/requests \
  -H "Authorization: Bearer $DIRECTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "resource_id": "RESOURCE_UUID",
    "quantity": 10,
    "comment": "Для регіонального офісу"
  }' | jq .

# Затвердження заявок
curl -X POST http://localhost/api/requests/REQUEST_ID/approve \
  -H "Authorization: Bearer $DIRECTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"comment": "Схвалено директором"}' | jq .

# Управління транспортом
curl -X POST http://localhost/api/vehicles \
  -H "Authorization: Bearer $DIRECTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "brand": "Volkswagen",
    "model": "Crafter",
    "plate_number": "AA1111BB",
    "type": "VAN"
  }' | jq .

# Перегляд аналітики
curl -X GET http://localhost/api/analytics/dashboard \
  -H "Authorization: Bearer $DIRECTOR_TOKEN" | jq .
```

#### ❌ ЗАБОРОНЕНО:
```bash
# Створення складів (тільки STOREKEEPER)
curl -X POST http://localhost/api/warehouses \
  -H "Authorization: Bearer $DIRECTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Новий склад"}' | jq .
# Очікується: 403 Forbidden

# Створення категорій ресурсів (тільки STOREKEEPER)
curl -X POST http://localhost/api/inventory/categories \
  -H "Authorization: Bearer $DIRECTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Нова категорія"}' | jq .
# Очікується: 403 Forbidden
```

---

## 3️⃣ REGION_LOGISTICIAN - Логістика регіону

**Користувач:** `i.melnyk@logistics.ua`

### 3.1. Отримання токена
```bash
LOGIST_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"i.melnyk@logistics.ua","password":"password123"}' \
  | jq -r '.access_token')
```

### 3.2. Основні операції логіста

#### ✅ ДОЗВОЛЕНО:
```bash
# Створення заявок на постачання
curl -X POST http://localhost/api/requests \
  -H "Authorization: Bearer $LOGIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "resource_id": "RESOURCE_UUID",
    "quantity": 15,
    "comment": "Для поповнення складу"
  }' | jq .

# Затвердження заявок
curl -X POST http://localhost/api/requests/REQUEST_ID/approve \
  -H "Authorization: Bearer $LOGIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"comment": "Схвалено логістом"}' | jq .

# Відхилення заявок
curl -X POST http://localhost/api/requests/REQUEST_ID/reject \
  -H "Authorization: Bearer $LOGIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"comment": "Недостатньо запасів"}' | jq .

# Управління транспортом
curl -X GET http://localhost/api/vehicles \
  -H "Authorization: Bearer $LOGIST_TOKEN" | jq .

curl -X POST http://localhost/api/vehicles/VEHICLE_ID/fuel \
  -H "Authorization: Bearer $LOGIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "liters": 55.0,
    "odometer_km": 130000,
    "record_type": "REFUEL"
  }' | jq .

# Створення переміщень
curl -X POST http://localhost/api/inventory/shipments \
  -H "Authorization: Bearer $LOGIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_warehouse_id": "WAREHOUSE_UUID_1",
    "to_warehouse_id": "WAREHOUSE_UUID_2",
    "vehicle_id": "VEHICLE_UUID",
    "items": [
      {"resource_id": "RESOURCE_UUID", "quantity": 5}
    ]
  }' | jq .

# Перегляд аналітики
curl -X GET http://localhost/api/analytics/dashboard \
  -H "Authorization: Bearer $LOGIST_TOKEN" | jq .
```

#### ❌ ЗАБОРОНЕНО:
```bash
# Створення користувачів
curl -X POST http://localhost/api/admin/users \
  -H "Authorization: Bearer $LOGIST_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden

# Управління підрозділами
curl -X POST http://localhost/api/units \
  -H "Authorization: Bearer $LOGIST_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden

# Управління складами
curl -X POST http://localhost/api/warehouses \
  -H "Authorization: Bearer $LOGIST_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden
```

---

## 4️⃣ REGION_STOREKEEPER - Управління складами

**Користувач:** `m.honcharenko@logistics.ua`

### 4.1. Отримання токена
```bash
STOREKEEPER_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"m.honcharenko@logistics.ua","password":"password123"}' \
  | jq -r '.access_token')
```

### 4.2. Основні операції комірника

#### ✅ ДОЗВОЛЕНО:
```bash
# Управління складами
curl -X GET http://localhost/api/warehouses \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" | jq .

curl -X POST http://localhost/api/warehouses \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "unit_id": 2,
    "name": "Регіональний склад №11",
    "location_type": "STATIONARY",
    "latitude": 50.45,
    "longitude": 30.52
  }' | jq .

# Управління категоріями
curl -X POST http://localhost/api/inventory/categories \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Спецодяг",
    "description": "Робочий одяг та засоби захисту"
  }' | jq .

# Управління ресурсами
curl -X POST http://localhost/api/inventory/resources \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": "CATEGORY_UUID",
    "warehouse_id": "WAREHOUSE_UUID",
    "unit_id": 2,
    "name": "Каска захисна",
    "quantity": 100,
    "min_quantity": 20
  }' | jq .

# Оновлення ресурсів
curl -X PATCH http://localhost/api/inventory/resources/RESOURCE_ID \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 95,
    "condition": "GOOD"
  }' | jq .

# Списання ресурсів
curl -X POST http://localhost/api/inventory/resources/RESOURCE_ID/write-off \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 5,
    "reason": "Пошкоджено при транспортуванні"
  }' | jq .

# Призначення ресурсів
curl -X POST http://localhost/api/inventory/resources/RESOURCE_ID/assign \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "USER_UUID",
    "quantity": 1
  }' | jq .

# Інвентаризація
curl -X POST http://localhost/api/inventory/audit \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "warehouse_id": "WAREHOUSE_UUID",
    "items": [
      {
        "resource_id": "RESOURCE_UUID",
        "actual_quantity": 95,
        "notes": "Перевірено"
      }
    ]
  }' | jq .
```

#### ❌ ЗАБОРОНЕНО:
```bash
# Створення заявок на постачання
curl -X POST http://localhost/api/requests \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden

# Управління транспортом
curl -X POST http://localhost/api/vehicles \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden
```

---

## 5️⃣ BRANCH_MANAGER - Управління філією

**Користувач:** `n.boiko@logistics.ua`

### 5.1. Отримання токена
```bash
BRANCH_MGR_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"n.boiko@logistics.ua","password":"password123"}' \
  | jq -r '.access_token')
```

### 5.2. Операції менеджера філії

#### ✅ ДОЗВОЛЕНО:
```bash
# Створення користувачів у своїй філії
curl -X POST http://localhost/api/admin/users \
  -H "Authorization: Bearer $BRANCH_MGR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "branch_worker",
    "email": "branch_worker@logistics.ua",
    "full_name": "Працівник Філії",
    "role": "EMPLOYEE",
    "unit_id": 4
  }' | jq .

# Управління підрозділами філії
curl -X POST http://localhost/api/units \
  -H "Authorization: Bearer $BRANCH_MGR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Нова команда",
    "unit_type": "TEAM",
    "parent_id": 4
  }' | jq .

# Зміна керівника підрозділу
curl -X POST http://localhost/api/units/UNIT_ID/change-commander \
  -H "Authorization: Bearer $BRANCH_MGR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"commander_id": "USER_UUID"}' | jq .

# Перегляд ієрархії
curl -X GET http://localhost/api/units/my-hierarchy \
  -H "Authorization: Bearer $BRANCH_MGR_TOKEN" | jq .

# Перегляд підлеглих користувачів
curl -X GET http://localhost/api/users/visible \
  -H "Authorization: Bearer $BRANCH_MGR_TOKEN" | jq .

# Блокування користувачів
curl -X PUT http://localhost/api/users/USER_ID/block \
  -H "Authorization: Bearer $BRANCH_MGR_TOKEN" | jq .

# Перегляд аудит-логів
curl -X GET http://localhost/api/admin/audit-logs \
  -H "Authorization: Bearer $BRANCH_MGR_TOKEN" | jq .
```

#### ❌ ЗАБОРОНЕНО:
```bash
# Створення заявок (тільки логісти)
curl -X POST http://localhost/api/requests \
  -H "Authorization: Bearer $BRANCH_MGR_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden

# Управління транспортом
curl -X GET http://localhost/api/vehicles \
  -H "Authorization: Bearer $BRANCH_MGR_TOKEN" | jq .
# Очікується: 403 Forbidden
```

---

## 6️⃣ BRANCH_LOGISTICIAN - Логістика філії

**Користувач:** `m.hrytsenko@logistics.ua`

### 6.1. Отримання токена
```bash
BRANCH_LOG_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"m.hrytsenko@logistics.ua","password":"password123"}' \
  | jq -r '.access_token')
```

### 6.2. Операції логіста філії

#### ✅ ДОЗВОЛЕНО:
```bash
# Всі операції як у REGION_LOGISTICIAN, але в межах своєї філії
curl -X POST http://localhost/api/requests \
  -H "Authorization: Bearer $BRANCH_LOG_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "resource_id": "RESOURCE_UUID",
    "quantity": 10,
    "comment": "Для філії"
  }' | jq .

curl -X GET http://localhost/api/vehicles \
  -H "Authorization: Bearer $BRANCH_LOG_TOKEN" | jq .

curl -X POST http://localhost/api/inventory/shipments \
  -H "Authorization: Bearer $BRANCH_LOG_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_warehouse_id": "WH_UUID_1",
    "to_warehouse_id": "WH_UUID_2",
    "items": []
  }' | jq .
```

---

## 7️⃣ DEPT_MANAGER - Управління відділом

**Користувач:** `y.kravchenko@logistics.ua`

### 7.1. Отримання токена
```bash
DEPT_MGR_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"y.kravchenko@logistics.ua","password":"password123"}' \
  | jq -r '.access_token')
```

### 7.2. Операції начальника відділу

#### ✅ ДОЗВОЛЕНО:
```bash
# Створення користувачів у своєму відділі
curl -X POST http://localhost/api/admin/users \
  -H "Authorization: Bearer $DEPT_MGR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "dept_employee",
    "email": "dept_employee@logistics.ua",
    "full_name": "Працівник Відділу",
    "role": "EMPLOYEE",
    "unit_id": 8
  }' | jq .

# Створення команд у своєму відділі
curl -X POST http://localhost/api/units \
  -H "Authorization: Bearer $DEPT_MGR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Нова команда відділу",
    "unit_type": "TEAM",
    "parent_id": 7
  }' | jq .

# Перегляд своєї ієрархії
curl -X GET http://localhost/api/units/my-hierarchy \
  -H "Authorization: Bearer $DEPT_MGR_TOKEN" | jq .

# Перегляд інвентаря
curl -X GET http://localhost/api/inventory/resources \
  -H "Authorization: Bearer $DEPT_MGR_TOKEN" | jq .
```

#### ❌ ЗАБОРОНЕНО:
```bash
# Управління складами
curl -X POST http://localhost/api/warehouses \
  -H "Authorization: Bearer $DEPT_MGR_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden

# Створення заявок
curl -X POST http://localhost/api/requests \
  -H "Authorization: Bearer $DEPT_MGR_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden

# Управління транспортом
curl -X GET http://localhost/api/vehicles \
  -H "Authorization: Bearer $DEPT_MGR_TOKEN" | jq .
# Очікується: 403 Forbidden
```

---

## 8️⃣ TEAM_LEAD - Керівник команди

**Користувач:** `o.zakharchenko@logistics.ua`

### 8.1. Отримання токена
```bash
TEAM_LEAD_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"o.zakharchenko@logistics.ua","password":"password123"}' \
  | jq -r '.access_token')
```

### 8.2. Операції керівника команди

#### ✅ ДОЗВОЛЕНО:
```bash
# Перегляд інформації
curl -X GET http://localhost/api/auth/me \
  -H "Authorization: Bearer $TEAM_LEAD_TOKEN" | jq .

# Перегляд своєї команди
curl -X GET http://localhost/api/users/visible \
  -H "Authorization: Bearer $TEAM_LEAD_TOKEN" | jq .

# Перегляд ієрархії
curl -X GET http://localhost/api/units/my-hierarchy \
  -H "Authorization: Bearer $TEAM_LEAD_TOKEN" | jq .

# Перегляд інвентаря
curl -X GET http://localhost/api/inventory/resources \
  -H "Authorization: Bearer $TEAM_LEAD_TOKEN" | jq .

# Моє обладнання
curl -X GET http://localhost/api/inventory/my-equipment \
  -H "Authorization: Bearer $TEAM_LEAD_TOKEN" | jq .

# Перегляд заявок
curl -X GET http://localhost/api/requests \
  -H "Authorization: Bearer $TEAM_LEAD_TOKEN" | jq .

# Створення переміщень
curl -X POST http://localhost/api/inventory/shipments \
  -H "Authorization: Bearer $TEAM_LEAD_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_warehouse_id": "WH_UUID_1",
    "to_warehouse_id": "WH_UUID_2",
    "items": []
  }' | jq .

# Оновлення свого профілю
curl -X PATCH http://localhost/api/users/profile \
  -H "Authorization: Bearer $TEAM_LEAD_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Захарченко Ольга Павлівна (оновлено)",
    "phone": "+380501112233"
  }' | jq .
```

#### ❌ ЗАБОРОНЕНО:
```bash
# Створення користувачів
curl -X POST http://localhost/api/admin/users \
  -H "Authorization: Bearer $TEAM_LEAD_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden

# Управління ролями
curl -X PUT http://localhost/api/users/USER_ID/role \
  -H "Authorization: Bearer $TEAM_LEAD_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden

# Створення заявок
curl -X POST http://localhost/api/requests \
  -H "Authorization: Bearer $TEAM_LEAD_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden
```

---

## 9️⃣ EMPLOYEE - Працівник

**Користувач:** `o.kovalenko@logistics.ua`

### 9.1. Отримання токена
```bash
EMPLOYEE_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"o.kovalenko@logistics.ua","password":"password123"}' \
  | jq -r '.access_token')
```

### 9.2. Операції працівника

#### ✅ ДОЗВОЛЕНО:
```bash
# Базова інформація
curl -X GET http://localhost/api/auth/me \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" | jq .

# Моє обладнання
curl -X GET http://localhost/api/inventory/my-equipment \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" | jq .

# Перегляд ресурсів (тільки читання)
curl -X GET http://localhost/api/inventory/resources \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" | jq .

# Перегляд категорій
curl -X GET http://localhost/api/inventory/categories \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" | jq .

# Перегляд складів
curl -X GET http://localhost/api/warehouses \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" | jq .

# Перегляд заявок
curl -X GET http://localhost/api/requests \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" | jq .

# Перегляд переміщень
curl -X GET http://localhost/api/inventory/shipments \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" | jq .

# Створення переміщень (якщо має доступ до складів)
curl -X POST http://localhost/api/inventory/shipments \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_warehouse_id": "WH_UUID",
    "to_warehouse_id": "WH_UUID_2",
    "items": []
  }' | jq .

# Оновлення свого профілю
curl -X PATCH http://localhost/api/users/profile \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"phone": "+380509998877"}' | jq .

# Зміна свого паролю
curl -X PATCH http://localhost/api/users/password \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword456"
  }' | jq .

# Перегляд дашборду
curl -X GET http://localhost/api/analytics/dashboard \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" | jq .
```

#### ❌ ЗАБОРОНЕНО (майже все керівне):
```bash
# Створення будь-чого
curl -X POST http://localhost/api/inventory/resources \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden

curl -X POST http://localhost/api/requests \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden

curl -X POST http://localhost/api/admin/users \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" \
  -d '{}' | jq .
# Очікується: 403 Forbidden
```

---

## 🎯 Швидкі сценарії тестування

### Сценарій 1: Повний цикл заявки на постачання
```bash
# 1. Логіст створює заявку
LOGIST_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"i.melnyk@logistics.ua","password":"password123"}' \
  | jq -r '.access_token')

RESOURCE_ID=$(curl -s -X GET http://localhost/api/inventory/resources \
  -H "Authorization: Bearer $LOGIST_TOKEN" | jq -r '.[0].id')

REQUEST_ID=$(curl -s -X POST http://localhost/api/requests \
  -H "Authorization: Bearer $LOGIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"resource_id\":\"$RESOURCE_ID\",\"quantity\":5,\"comment\":\"Тест\"}" \
  | jq -r '.id')

# 2. Директор затверджує
DIRECTOR_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"d.tkachenko@logistics.ua","password":"password123"}' \
  | jq -r '.access_token')

curl -X POST http://localhost/api/requests/$REQUEST_ID/approve \
  -H "Authorization: Bearer $DIRECTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"comment":"Затверджено"}' | jq .

# 3. Перевірка статусу
curl -X GET http://localhost/api/requests/$REQUEST_ID \
  -H "Authorization: Bearer $LOGIST_TOKEN" | jq .
```

### Сценарій 2: Управління інвентарем
```bash
# 1. Комірник створює категорію
STOREKEEPER_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"m.honcharenko@logistics.ua","password":"password123"}' \
  | jq -r '.access_token')

CATEGORY_ID=$(curl -s -X POST http://localhost/api/inventory/categories \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Тестова категорія","description":"Для тестування"}' \
  | jq -r '.id')

# 2. Отримує список складів
WAREHOUSE_ID=$(curl -s -X GET http://localhost/api/warehouses \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" | jq -r '.[0].id')

# 3. Створює ресурс
curl -X POST http://localhost/api/inventory/resources \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"category_id\":\"$CATEGORY_ID\",
    \"warehouse_id\":\"$WAREHOUSE_ID\",
    \"unit_id\":2,
    \"name\":\"Тестовий ресурс\",
    \"quantity\":100,
    \"min_quantity\":10
  }" | jq .
```

### Сценарій 3: Перевірка лімітів BASIC
```bash
ADMIN_TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"markostrutinsky@gmail.com","password":"password123"}' \
  | jq -r '.access_token')

# Спроба створити 11-й склад (має не вдатися)
curl -X POST http://localhost/api/warehouses \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "unit_id":2,
    "name":"Склад №11 (понад ліміт)",
    "location_type":"STATIONARY",
    "latitude":50.45,
    "longitude":30.52
  }' | jq .
# Очікується: 422 Unprocessable Entity або подібна помилка з лімітом

# Спроба доступу до PRO фічі
curl -X GET http://localhost/api/gps/fleet-map \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
# Очікується: 402 Payment Required
```

---

## 📝 Корисні команди для тестування

### Масовий експорт даних для всіх ролей
```bash
#!/bin/bash

ROLES=(
  "markostrutinsky@gmail.com:ADMIN"
  "d.tkachenko@logistics.ua:DIRECTOR"
  "i.melnyk@logistics.ua:LOGIST"
  "m.honcharenko@logistics.ua:STOREKEEPER"
  "n.boiko@logistics.ua:BRANCH_MGR"
  "o.kovalenko@logistics.ua:EMPLOYEE"
)

for role in "${ROLES[@]}"; do
  email="${role%%:*}"
  name="${role##*:}"
  
  echo "Testing $name ($email)..."
  
  TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$email\",\"password\":\"password123\"}" \
    | jq -r '.access_token')
  
  echo "  - Auth: OK"
  
  curl -s -X GET http://localhost/api/auth/me \
    -H "Authorization: Bearer $TOKEN" > "test_${name}_me.json"
  
  curl -s -X GET http://localhost/api/units \
    -H "Authorization: Bearer $TOKEN" > "test_${name}_units.json"
  
  echo "  - Data exported"
done
```

### Перевірка всіх ендпоінтів для ролі
```bash
#!/bin/bash

EMAIL="i.melnyk@logistics.ua"
TOKEN=$(curl -s -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$EMAIL\",\"password\":\"password123\"}" \
  | jq -r '.access_token')

ENDPOINTS=(
  "/api/auth/me"
  "/api/users/commanders"
  "/api/units"
  "/api/inventory/categories"
  "/api/inventory/resources"
  "/api/warehouses"
  "/api/vehicles"
  "/api/requests"
  "/api/analytics/dashboard"
)

for endpoint in "${ENDPOINTS[@]}"; do
  echo "Testing $endpoint..."
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "http://localhost$endpoint")
  echo "  Status: $STATUS"
done
```

---

## ⚠️ Важливі примітки

1. **Тариф BASIC** - досягнуто ліміту складів (10/10), тому створення нових складів не працюватиме
2. **PRO фічі недоступні** - GPS, Smart Dispatch, Excel Import, Advanced Analytics - повертають 402
3. **Ієрархія доступу** - кожна роль бачить тільки свій рівень та нижчі підрозділи
4. **Паролі** - всі користувачі мають пароль `password123`
5. **Email формат** - username може бути або email, або частина до @

---

**Версія документа:** 1.0  
**Дата створення:** 2 травня 2026  
**Тенант:** Пійлівський Логістичний Центр (piylo)  
**Тариф:** BASIC
