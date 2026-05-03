<div align="center">

# 📦 Omnilog (Omnilog)

### Інтелектуальна платформа управління логістикою організацій

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18.2+-61DAFB?style=for-the-badge&logo=react&logoColor=black)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

**[Демо](http://localhost)** | **[Документація](docs/)** | **[API Docs](#-api-документація)** | **[Архітектура](ARCHITECTURE.md)**

</div>

---

## 📋 Зміст

- [Про проєкт](#-про-проєкт)
- [Ключові особливості](#-ключові-особливості)
- [Технологічний стек](#️-технологічний-стек)
- [Швидкий старт](#-швидкий-старт)
- [Доступи до системи](#-доступи-до-системи)
- [Ролі користувачів](#-ролі-користувачів)
- [Тарифні плани](#-тарифні-плани)
- [API документація](#-api-документація)
- [Розробка](#-розробка)
- [Архітектура](#-архітектура)

---

## 🎯 Про проєкт

**Omnilog (Omnilog)** — це сучасна SaaS-платформа для управління логістикою організацій з підтримкою multi-tenant архітектури, GPS-трекінгу, AI-аналітики та системи запобігання шахрайству.

### Основні можливості:

- 📦 **Управління ресурсами** — облік категорій, складів, автоматичне відстеження критичних залишків
- 🚚 **Заявки на постачання** — повний життєвий цикл від створення до виконання з SLA-моніторингом
- 🤝 **Інтеграція з підрядниками** — координація зовнішніх постачань через контрагентів
- 🚙 **Автопарк** — GPS-трекінг, геозони, прогнозне обслуговування техніки
- 📊 **Аналітика та KPI** — dashboard з метриками SLA, TCO, Risk Score, Demand Forecasting
- 🛡️ **Anti-Fraud система** — детекція аномалій у витраті палива з AI-алгоритмами
- 👥 **Multi-tenant** — підтримка необмеженої кількості організацій з ізоляцією даних
- 🔐 **Безпека** — JWT авторизація, RBAC (7 ролей), audit logging всіх операцій

---

## ✨ Ключові особливості

### � Для організацій з розподіленою структурою
- Ієрархічна структура: Компанія → Регіональний відділ → Філія → Локальний підрозділ
- Контроль руху ресурсів між підрозділами
- Автоматичне оновлення балансу після виконання заявок
- Excel-імпорт для масового завантаження (PRO)

### 🚀 Premium Features (PRO tier)
1. **Advanced KPI Dashboard** — SLA%, TCO, Risk Score, Depletion Forecast
2. **Demand Forecasting** — AI-прогнозування потреб з 3 сценаріями (Low/Medium/High)
3. **Predictive Maintenance** — автоматичний розрахунок графіку ТО транспорту
4. **Fuel Anti-Fraud Detection** — детекція підозрілих заправок і аномалій споживання
5. **Real-Time GPS Tracking** — карта флоту, траєкторії, геозони з алертами

### 🔒 Безпека та контроль
- JWT токени (access + refresh)
- Role-Based Access Control (7 ролей)
- Multi-tenant ізоляція на рівні БД
- Повний audit trail усіх операцій
- Subscription-based feature gating

---

## 🛠️ Технологічний стек

### Backend
- **Go 1.21+** — високопродуктивний API сервер
- **Gin Framework** — швидкий HTTP router
- **PostgreSQL 15** — реляційна БД з повною підтримкою JSON
- **pgx/v5** — нативний PostgreSQL драйвер
- **JWT** — авторизація та refresh tokens

### Frontend
- **React 18.2** — UI бібліотека
- **TypeScript 5.0** — типізований JavaScript
- **Vite** — швидкий bundler
- **CSS Modules** — ізольовані стилі

### Infrastructure
- **Docker Compose** — оркестрація контейнерів
- **Caddy** — reverse proxy з автоматичним HTTPS
- **SMTP** — Email нотифікації (Gmail/SendGrid)

---

## 🚀 Швидкий старт

### Вимоги
- Docker 20.10+
- Docker Compose 2.0+
- mkcert (для локального HTTPS)
- (Опціонально) Go 1.21+ та Node.js 18+ для локальної розробки

### Запуск (Docker)

```bash
# 1. Клонуйте репозиторій
git clone https://github.com/markostrutinsky/diploma.git
cd diploma

# 2. Налаштуйте змінні середовища
cp .env.example .env
# Відредагуйте .env (SMTP credentials, JWT secret)

# 3. Налаштуйте HTTPS (один раз)
./setup-certs.sh
# Встановить локальний CA та згенерує довірені сертифікати

# 4. Запустіть всі сервіси
docker compose up -d --build

# 5. Перевірте статус
docker compose ps
```

**Додаток буде доступний:**
- 🌐 **Frontend (HTTPS):** https://localhost
- 🔌 **API:** https://localhost/api
- 💾 **PostgreSQL:** localhost:5432

> 🔒 **HTTPS:** Проект використовує mkcert для довірених локальних сертифікатів. Браузери автоматично довіряють з'єднанню без попереджень! Детальніше: [HTTPS_SETUP.md](HTTPS_SETUP.md)

---

## 🔑 Доступи до системи

### 1️⃣ Системний адміністратор (Platform Owner)

**Для управління всіма організаціями та тарифами:**

```
URL:      https://localhost/platform
Email:    platform@omnilog.system
Password: AdminSystem2024!
```

**Можливості:**
- ✅ Крос-тенантний доступ до всіх організацій
- ✅ Управління тарифними планами (FREE → PRO)
- ✅ Блокування/активація організацій
- ✅ Перегляд глобальної статистики платформи

> 💡 Створюється автоматично при першому запуску системи

---

### 2️⃣ Нова організація / бізнес

**Крок 1:** Відкрийте http://localhost/signup

**Крок 2:** Заповніть форму реєстрації:
- Назва організації
- Email (стане логіном першого адміна)
- Пароль

**Крок 3:** Автоматично отримуєте:
- 🎟️ Тариф: **FREE** (базовий функціонал)
- 👤 Роль: **TENANT_ADMIN** (власник організації)
- 🏢 Власний tenant (ізольовані дані)

**Крок 4:** Оновлення тарифу:
- Перейдіть на `/billing`
- Надішліть запит на оновлення до STANDARD або PRO
- Очікуйте схвалення від Platform Admin

---

### 3️⃣ Підрядник (Постачальник)

**Крок 1:** Відкрийте http://localhost/register

**Крок 2:** Зареєструйтесь:
- ПІБ
- Email
- Пароль
- Організація (опціонально)

**Крок 3:** Автоматично отримуєте:
- 👷 Роль: **CONTRACTOR**
- 📋 Доступ до заявок на постачання
- 🚚 Можливість брати заявки в роботу

---

## 👥 Ролі користувачів

| Роль | Рівень | Опис | Основні можливості |
|------|--------|------|-------------------|
| **SYSTEM_ADMIN** | 🌐 Платформа | Власник платформи | Управління всіма tenant'ами, тарифами, глобальна аналітика |
| **TENANT_ADMIN** | 🏢 Організація | Власник організації | Повні права в межах свого tenant, запит оновлення тарифу, управління користувачами |
| **REGION_DIRECTOR** | 🗺️ Регіон | Директор регіону | Управління регіональними підрозділами, погодження заявок |
| **COMMANDER** | 📊 Підрозділ | Менеджер відділу | Створення заявок, перегляд ресурсів, затвердження заявок підлеглих |
| **WAREHOUSE** | 📦 Склад | Комірник | Управління категоріями, ресурсами, складами, видача товарів |
| **LOGISTICS** | 🚚 Логістика | Логіст | Затвердження заявок, управління поставками, аналітика |
| **CONTRACTOR** | 🤝 Зовнішній | Підрядник/Постачальник | Перегляд заявок, взяття в роботу, доставка |

### Матриця доступу

| Функція | SYSTEM_ADMIN | TENANT_ADMIN | COMMANDER | WAREHOUSE | CONTRACTOR |
|---------|:------------:|:------------:|:---------:|:---------:|:----------:|
| Крос-тенантний доступ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Управління користувачами | ✅ | ✅ | ❌ | ❌ | ❌ |
| Створення заявок | ✅ | ✅ | ✅ | ❌ | ❌ |
| Затвердження заявок | ✅ | ✅ | ✅ | ❌ | ❌ |
| Управління складами | ✅ | ✅ | ❌ | ✅ | ❌ |
| Управління ресурсами | ✅ | ✅ | ❌ | ✅ | ❌ |
| GPS трекінг (PRO) | ✅ | ✅ | ✅ | ❌ | ❌ |
| Заявки підрядникам | ✅ | ✅ | ✅ | ❌ | ✅ |
| Advanced Analytics (PRO) | ✅ | ✅ | ❌ | ❌ | ❌ |

---

## 💎 Тарифні плани

| Функція | FREE | BASIC | STANDARD | PRO |
|---------|:----:|:-----:|:--------:|:---:|
| **Базовий функціонал** | | | | |
| Управління ресурсами | ✅ | ✅ | ✅ | ✅ |
| Заявки на постачання | ✅ | ✅ | ✅ | ✅ |
| Складський облік | ✅ | ✅ | ✅ | ✅ |
| Управління користувачами | ✅ | ✅ | ✅ | ✅ |
| Базова аналітика | ✅ | ✅ | ✅ | ✅ |
| **Розширені можливості** | | | | |
| Автопарк і паливо | ❌ | ✅ | ✅ | ✅ |
| Заявки для підрядників | ❌ | ✅ | ✅ | ✅ |
| SLA моніторинг | ❌ | ✅ | ✅ | ✅ |
| Audit logging | ❌ | ✅ | ✅ | ✅ |
| Smart Dispatch | ❌ | ❌ | ✅ | ✅ |
| Excel імпорт/експорт | ❌ | ❌ | ✅ | ✅ |
| **Premium Features** | | | | |
| Advanced KPI Dashboard | ❌ | ❌ | ❌ | ✅ |
| Demand Forecasting (AI) | ❌ | ❌ | ❌ | ✅ |
| Predictive Maintenance | ❌ | ❌ | ❌ | ✅ |
| Fuel Anti-Fraud Detection | ❌ | ❌ | ❌ | ✅ |
| Real-Time GPS Tracking | ❌ | ❌ | ❌ | ✅ |
| Geofencing & Alerts | ❌ | ❌ | ❌ | ✅ |
| **Обмеження** | | | | |
| Кількість користувачів | 10 | 50 | 200 | ∞ |
| Кількість ресурсів | 100 | 500 | 2000 | ∞ |
| Кількість транспорту | 0 | 20 | 100 | ∞ |
| Підтримка | Email | Email | Priority | 24/7 |

---

## 📚 API документація

Детальну документацію API дивіться у [API_DOCUMENTATION.md](API_DOCUMENTATION.md)

### Базова URL
```
http://localhost/api
```

### Аутентифікація
Використовується JWT Bearer токен:
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost/api/inventory/resources
```

### Швидкий старт API

#### 1. Реєстрація організації
```bash
POST /api/auth/tenants/signup
Content-Type: application/json

{
  "tenant_name": "26-та Окрема Бригада",
  "email": "admin@brigade26.mil",
  "password": "SecurePass123!"
}
```

#### 2. Логін
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "admin@brigade26.mil",
  "password": "SecurePass123!"
}

# Response:
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": { ... }
}
```

#### 3. Отримання ресурсів
```bash
GET /api/inventory/resources
Authorization: Bearer eyJhbGc...
```

### Основні групи endpoints

- 🔐 `/api/auth/*` — Аутентифікація (login, register, refresh)
- 👥 `/api/users/*` — Управління користувачами
- 🏢 `/api/units/*` — Підрозділи та організаційна структура
- 📦 `/api/inventory/*` — Ресурси, категорії, склади
- 📝 `/api/requests/*` — Заявки на постачання
- 🤝 `/api/contractor-requests/*` — Волонтерські заявки
- 🚙 `/api/vehicles/*` — Автопарк та паливо
- 📊 `/api/analytics/*` — Аналітика та звіти
- 🌍 `/api/gps/*` — GPS трекінг (PRO)
- 🔧 `/api/admin/*` — Адміністрування
- 🌐 `/api/platform/*` — Platform Admin (SYSTEM_ADMIN only)

**Повний список з прикладами:** [API_DOCUMENTATION.md](API_DOCUMENTATION.md)

---

## 🔧 Розробка

### Локальний запуск (без Docker)

#### Backend

```bash

```bash
# 1. Запустіть PostgreSQL (або використовуйте Docker)
docker run -d \
  --name Omnilog-postgres \
  -e POSTGRES_DB=Omnilog \
  -e POSTGRES_USER=Omnilog \
  -e POSTGRES_PASSWORD=Omnilog123 \
  -p 5432:5432 \
  postgres:15-alpine

# 2. Налаштуйте .env
cat > .env << EOF
DATABASE_URL=postgresql://Omnilog:Omnilog123@localhost:5432/Omnilog?sslmode=disable
JWT_SECRET=dev-secret-change-in-production
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your@gmail.com
SMTP_PASSWORD=your-app-password
FRONTEND_URL=http://localhost:5173
PORT=8080
ALLOW_BOOTSTRAP_OVERRIDE=true
EOF

# 3. Запустіть backend
cd Omnilog_backend
go mod download
go run main.go

# Backend доступний на http://localhost:8080
```

#### Frontend
```bash
# 1. Встановіть залежності
cd Omnilog_frontend
npm install

# 2. Запустіть dev сервер
npm run dev

# Frontend доступний на http://localhost:5173
```

### Структура проєкту

```
diploma/
├── Omnilog_backend/          # Go API сервер
│   ├── internal/
│   │   ├── handlers/        # HTTP обробники (Controllers)
│   │   ├── services/        # Бізнес-логіка
│   │   ├── repositories/    # Робота з БД (Data Access Layer)
│   │   ├── models/          # Структури даних
│   │   ├── middleware/      # Auth, RBAC, Subscription
│   │   ├── database/        # Міграції
│   │   └── tokens/          # JWT utilities
│   ├── main.go              # Entry point
│   └── go.mod
├── Omnilog_frontend/         # React UI
│   ├── src/
│   │   ├── components/      # React компоненти
│   │   ├── pages/           # Сторінки
│   │   ├── contexts/        # React Context (Auth)
│   │   ├── api/             # API client
│   │   └── hooks/           # Custom hooks
│   ├── package.json
│   └── vite.config.ts
├── docker/                  # Docker конфігурація
│   └── postgres/
│       └── init.sql         # Початкові дані
├── scripts/                 # Утиліти
│   ├── gps_simulator/       # Симулятор GPS даних
│   └── seed_db/             # Заповнення БД тестовими даними
├── docker-compose.yml       # Оркестрація
├── Caddyfile               # Reverse proxy
└── README.md

```

### Змінні середовища (.env)

```env
# Database
DATABASE_URL=postgresql://Omnilog:Omnilog123@postgres:5432/Omnilog?sslmode=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Email (Gmail example)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your-email@gmail.com
SMTP_PASSWORD=your-app-password  # Use App Password, not regular password

# Frontend
FRONTEND_URL=http://localhost

# Backend
PORT=8080

# Bootstrap (для dev/testing)
ALLOW_BOOTSTRAP_OVERRIDE=true  # Дозволити створення admin'а навіть якщо користувачі вже є
```

### Корисні команди

```bash
# Перезапустити всі сервіси
docker compose restart

# Переглянути логи
docker compose logs -f backend
docker compose logs -f frontend

# Зупинити та видалити всі дані
docker compose down -v

# Перебудувати тільки backend
docker compose up -d --build backend

# Виконати міграції вручну
docker compose exec backend go run main.go migrate

# Підключитись до БД
docker compose exec postgres psql -U Omnilog -d Omnilog

# Backup БД
docker compose exec postgres pg_dump -U Omnilog Omnilog > backup.sql

# Restore БД
docker compose exec -T postgres psql -U Omnilog Omnilog < backup.sql
```

---

## 🏗️ Архітектура

### Високорівнева діаграма

```
┌─────────────────────────────────────────────────────────┐
│                    Internet / Users                      │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
            ┌────────────────┐
            │  Caddy Proxy   │  (HTTPS, Reverse Proxy)
            │  Port 80/443   │
            └────────┬───────┘
                     │
        ┌────────────┴────────────┐
        │                         │
        ▼                         ▼
┌───────────────┐        ┌────────────────┐
│   Frontend    │        │    Backend     │
│  React + TS   │◄──────►│   Go + Gin     │
│  Port 5173    │  HTTP  │   Port 8080    │
└───────────────┘        └────────┬───────┘
                                  │
                                  ▼
                         ┌────────────────┐
                         │  PostgreSQL    │
                         │   Port 5432    │
                         └────────────────┘
```

### Backend архітектура (Clean Architecture)

```
HTTP Request
    │
    ▼
┌─────────────────────────────────────────┐
│         Middleware Layer                │
│  • AuthMiddleware (JWT)                 │
│  • TenantMiddleware (Multi-tenant)      │
│  • SubscriptionMiddleware (Tier check)  │
│  • RequireRoleMiddleware (RBAC)         │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Handlers (Controllers)          │
│  • Parse request                        │
│  • Validate input                       │
│  • Call service                         │
│  • Return response                      │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Services (Business Logic)       │
│  • Inventory, Requests, Analytics       │
│  • GPS Tracking, Fuel, Vehicles         │
│  • Auth, Email, Audit                   │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Repositories (Data Access)      │
│  • SQL queries                          │
│  • Transaction management               │
│  • Tenant isolation                     │
└──────────────┬──────────────────────────┘
               │
               ▼
         PostgreSQL Database
```

### Multi-Tenant ізоляція

```sql
-- Кожна таблиця має tenant_id
CREATE TABLE resources (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,  -- 👈 Ізоляція
    name TEXT,
    quantity INTEGER,
    ...
);

-- Middleware автоматично додає WHERE tenant_id = ?
SELECT * FROM resources WHERE tenant_id = 'current-tenant-id';
```

### Subscription Tiers (Feature Gating)

```go
// Middleware перевіряє tier перед доступом
middleware.RequireSubscriptionTier("PRO", dbPool)

// Приклад:
router.GET("/analytics/kpi", 
    middleware.AuthMiddleware(),
    middleware.RequireSubscriptionTier("PRO"),  // 👈 Доступ тільки PRO
    handler.GetKPI
)
```

Детальна архітектура описана у [ARCHITECTURE.md](ARCHITECTURE.md)

---

## 🧪 Тестування

### GPS Simulator (для демо)
```bash
cd scripts/gps_simulator
pip install -r requirements.txt
python simulate.py --vehicles 5 --duration 3600
```

### Database Seeder (тестові дані)
```bash
cd scripts/seed_db
pip install -r requirements.txt
python seed.py
```

---

## 📄 Ліцензія

Цей проєкт розповсюджується під ліцензією MIT. Деталі у файлі [LICENSE](LICENSE).

---

## 👨‍💻 Автор

**Марко Струтинський**  
🎓 Дипломний проєкт • 2026  
📧 Email: [markostrutinsky@example.com](mailto:markostrutinsky@example.com)  
💼 GitHub: [@markostrutinsky](https://github.com/markostrutinsky)

---

## 🙏 Подяки

- Go Community за чудову документацію
- React Team за потужну бібліотеку
- PostgreSQL за надійну БД
- Всім open-source contributors

---

<div align="center">

**[⬆ Повернутись до змісту](#-зміст)**

Made with ❤️ in Ukraine 🇺🇦

</div>
````
```

**Примітка:** Якщо БД вже існує, нові таблиці (inventory, requests, vehicles) можуть не з’явитися. Для чистого старту або якщо міграції не пройшли:
```bash
docker compose down -v
docker compose up -d --build
```
