# 🔐 Демо-доступи до системи

## Системні ролі

### 1️⃣ SYSTEM_ADMIN (Власник платформи)

**Credentials:**
```
Email:    platform@omnilog.system
Password: AdminSystem2024!
```

**Можливості:**
- ✅ Крос-тенантний доступ до всіх організацій
- ✅ Управління тарифами через `/platform`
- ✅ Перегляд статистики по всіх tenant'ах
- ✅ Bypass всіх платних обмежень

**Сторінка входу:** http://localhost/login  
**Адмін-панель:** http://localhost/platform

---

### 2️⃣ TENANT_ADMIN (Власник організації)

**Як створити:**
1. Перейти на http://localhost/signup
2. Заповнити форму створення організації:
   - Назва організації
   - Email власника
   - Пароль
3. Перший користувач автоматично отримує роль **TENANT_ADMIN**

**Можливості:**
- ✅ Повні права в межах своєї організації
- ✅ Запит оновлення тарифу на `/billing`
- ✅ Створення підрозділів та користувачів
- ⛔ Немає доступу до інших організацій
- ⛔ Немає доступу до `/platform`

**Тариф за замовчуванням:** BASIC → можна оновити до PRO/ENTERPRISE через SYSTEM_ADMIN

---

### 3️⃣ CONTRACTOR (Підрядник)

**Як створити:**
1. Перейти на http://localhost/register
2. Заповнити форму реєстрації підрядника
3. Роль: **CONTRACTOR**

**Можливості:**
- ✅ Перегляд відкритих заявок від організацій
- ✅ Прийом заявок в роботу
- ✅ Відмітка про доставку

---

## Архітектура Multi-Tenant

```
┌─────────────────────────────────────────────────┐
│         SYSTEM_ADMIN (Platform Owner)           │
│  platform@omnilog.system | AdminSystem2024!     │
│                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌────────┐│
│  │  Tenant A    │  │  Tenant B    │  │ Tenant C││
│  │  (BASIC)     │  │  (PRO)       │  │ (ENTER.)││
│  └──────────────┘  └──────────────┘  └────────┘│
└─────────────────────────────────────────────────┘
           │                 │               │
           ▼                 ▼               ▼
    TENANT_ADMIN_A    TENANT_ADMIN_B   TENANT_ADMIN_C
    (власник орг-А)   (власник орг-Б)  (власник орг-С)
```

---

## Тарифні плани

| Тариф | Ціна | Основні можливості |
|-------|------|------------|
| **BASIC** | 0 грн/міс | 10 складів, 100 ресурсів, 50 юзерів, 5 авто — базовий облік |
| **PRO** | 4999 грн/міс | 100 складів, 1000 ресурсів, 500 юзерів, 50 авто + Smart Dispatch, GPS Tracking, Predictive Maintenance, Advanced Analytics, Fuel Anti-Fraud |
| **ENTERPRISE** | Custom | Безлімітні ресурси, підтримка 24/7, SLA гарантії, персональний менеджер |

**Детальне порівняння:**

### 🆓 BASIC (Безкоштовний)
- **Ліміти:** 10 складів | 100 товарів | 50 користувачів | 5 авто
- **Функції:** 
  - ✅ Базовий облік інвентаря
  - ✅ Ручне формування рейсів
  - ✅ Журнал аудиту (30 днів)
  - ✅ Базова звітність (PDF/Excel)

### 💎 PRO (4999 грн/міс)
- **Ліміти:** 100 складів | 1000 товарів | 500 користувачів | 50 авто
- **Функції BASIC +**
  - ✨ **Smart Dispatch** — автоматична оптимізація маршрутів
  - 🛰️ **GPS Tracking & Geofencing** — відстеження транспорту в реальному часі
  - 📊 **Advanced Analytics** — KPI dashboard (SLA%, TCO, ризик-метрики)
  - 🔮 **Demand Forecasting** — прогнозування попиту на 3 місяці
  - 🔧 **Predictive Maintenance** — прогноз ТО на основі пробігу
  - ⛽ **Fuel Anti-Fraud Detection** — виявлення аномалій заправок
  - 📥 **Excel Import/Export** — масовий імпорт/експорт
  - 🏪 **Smart Warehouse Replenishment** — автоматичне поповнення
  - 👨‍💻 Пріоритетна підтримка

### 🏢 ENTERPRISE (Індивідуальна ціна)
- **Ліміти:** ♾️ Безлімітні ресурси (склади, товари, користувачі, транспорт)
- **Функції PRO +**
  - 🆘 **Підтримка 24/7** з гарантованим часом відповіді
  - 📜 **SLA гарантії** з фінансовими компенсаціями
  - 👨‍💼 Персональний менеджер проєкту
  - 🔧 Кастомні інтеграції та доробки
  - 🎓 Навчання команди on-site
  - 🌍 Мульти-регіональна підтримка

**Симуляція оплати:**
- TENANT_ADMIN бачить тарифи на `/billing`
- Кнопка "Запросити оновлення" (в продакшені — інтеграція з LiqPay/Stripe)
- SYSTEM_ADMIN змінює тарифи через `/platform` → `Change Tier`

---

## RLS (Row-Level Security)

Кожна організація (tenant) ізольована на рівні БД:

```sql
-- Політика RLS на всіх tenant-таблицях:
CREATE POLICY tenant_isolation ON users
USING (
  current_setting('app.tenant_id', true) = ''  -- SYSTEM_ADMIN: повний доступ
  OR tenant_id::text = current_setting('app.tenant_id', true)  -- tenant бачить тільки свої дані
);
```

**Застосовується на таблицях:**
- `users`, `units`, `resources`, `categories`, `warehouses`
- `vehicles`, `fuel_records`, `supply_requests`, `contractor_requests`
- `shipments`, `audit_logs`, `gps_locations`, `geofences`

---

## Швидкий старт

### Сценарій 1: Демо для інвестора (Platform Admin)

```bash
# 1. Вхід як SYSTEM_ADMIN
http://localhost/login
Email: platform@omnilog.system
Password: AdminSystem2024!

# 2. Перегляд панелі /platform
- Список всіх організацій
- Зміна тарифів
- Статистика по tenant'ах
```

### Сценарій 2: Демо для клієнта (New Business)

```bash
# 1. Створення нової організації
http://localhost/signup
Назва: "Тестова Логістика"
Email: test@example.com
Password: TestPass123!

# 2. Вхід як TENANT_ADMIN
- Створення підрозділів
- Додавання користувачів (invite)
- Управління ресурсами

# 3. Запит оновлення тарифу
/billing → "Запросити оновлення до PRO"
(SYSTEM_ADMIN змінює через /platform)
```

---

## Примітки безпеки

⚠️ **Для production:**
1. Змінити `SYSTEM_ADMIN` credentials через ENV:
   ```bash
   SYSTEM_ADMIN_EMAIL=your-email@company.com
   SYSTEM_ADMIN_PASSWORD=your-bcrypt-hash
   ```
2. Вимкнути auto-creation SYSTEM_ADMIN у міграції
3. Додати 2FA для SYSTEM_ADMIN
4. Інтегрувати платіжну систему (LiqPay/Stripe) замість симуляції

---

**Створено:** 23.04.2026  
**Версія:** v2.0 (Multi-Tenant SaaS)
