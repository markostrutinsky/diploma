# 🔄 Рефакторинг тарифної системи: FREE → BASIC

## 📋 Зміни

### Проблема
Наявність двох безкоштовних тарифів (FREE і BASIC) створювала плутанину:
- FREE: 1 склад, 20 товарів, 5 користувачів
- BASIC: 10 складів, 100 товарів, 50 користувачів
- Обидва безкоштовні → немає сенсу в FREE

### Рішення
**Видалено тариф FREE**, тепер система має **3 чіткі рівні**:

1. **BASIC** (безкоштовний) — стартовий тариф
2. **PRO** (4999 грн/міс) — для бізнесу з преміум-фічами
3. **ENTERPRISE** (індивідуально) — корпоративний з безлімітами

---

## 🛠️ Технічні зміни

### Backend

#### 1. База даних (`migrate.go`)
```sql
-- Змінено DEFAULT для нових тенантів
subscription_tier VARCHAR(30) NOT NULL DEFAULT 'BASIC'  -- було 'FREE'

-- Міграція існуючих даних
UPDATE tenants SET subscription_tier = 'BASIC' WHERE subscription_tier = 'FREE';
UPDATE units SET subscription_tier = 'BASIC' WHERE subscription_tier = 'FREE';
```

#### 2. Сервіси

**`auth_service.go`**
```go
// Нові організації створюються з BASIC
insertTenant := `INSERT INTO tenants (name, slug, subscription_tier, owner_email, is_active)
    VALUES ($1, $2, 'BASIC', $3, TRUE)  -- було 'FREE'
    RETURNING id`
```

**`limitation_service.go`**
```go
// Видалено FREE з лімітів
var LimitsByTier = map[string]SubscriptionLimits{
    // "FREE": {...},  // видалено
    "BASIC": {
        MaxWarehouses: 10,
        MaxResources:  100,
        MaxUsers:      50,
        MaxVehicles:   5,
    },
    ...
}
```

**`middleware/subscription.go`**
```go
// Оновлено ієрархію тарифів
var SubscriptionTierWeight = map[string]int{
    // "FREE": 0,  // видалено
    "BASIC":      1,
    "PRO":        2,
    "ENTERPRISE": 3,
}

// Fallback значення змінено на BASIC
return "BASIC", nil  // було "FREE", nil
```

#### 3. Моделі

**`models/tenant.go`**
```go
const (
    // TierFree SubscriptionTier = "FREE"  // видалено
    TierBasic      SubscriptionTier = "BASIC"
    TierPro        SubscriptionTier = "PRO"
    TierEnterprise SubscriptionTier = "ENTERPRISE"
)
```

#### 4. Репозиторії

**`user_repository.go`**
```sql
COALESCE(t.subscription_tier, 'BASIC') AS effective_tier  -- було 'FREE'
```

**`unit_repository.go`**
```go
return "BASIC", nil  // було "FREE", nil (2 місця)
```

---

### Frontend

#### 1. Типи (`api/client.ts`)
```typescript
// Вже було правильно без FREE
export type SubscriptionTier = 'BASIC' | 'PRO' | 'ENTERPRISE';
```

#### 2. Константи (`constants/roles.ts`)
```typescript
// Вже було правильно без FREE
const TIER_WEIGHT: Record<Tier, number> = {
  BASIC: 0,
  PRO: 1,
  ENTERPRISE: 2,
}
```

#### 3. Компоненти

**`PlatformAdmin.tsx`**
```typescript
type Tenant = {
  subscription_tier: 'BASIC' | 'PRO' | 'ENTERPRISE'  // видалено 'FREE'
}

const TIERS = ['BASIC', 'PRO', 'ENTERPRISE'] as const  // видалено 'FREE'
```

**`Billing.tsx`**
Повністю переписано з 3 картками:
- ✅ **BASIC** — безкоштовний план (10 складів, 100 товарів)
- 💎 **PRO** — 4999 грн/міс (Smart Dispatch, GPS, Analytics)
- 🏢 **ENTERPRISE** — індивідуально (безліміти + 24/7 підтримка)

**`Billing.css`**
```css
/* Додано стилі для третьої картки */
.plan-card.enterprise {
  border: 2px solid #f59e0b;
  box-shadow: 0 10px 25px -5px rgba(245, 158, 11, 0.2);
}

/* Змінено grid на 3 колонки */
.plans-grid {
  grid-template-columns: repeat(3, 1fr);
}
```

---

### Документація

**`DEMO_ACCESS.md`**
```markdown
## Тарифні плани

| Тариф | Ціна | Основні можливості |
|-------|------|------------|
| **BASIC** | 0 грн/міс | 10 складів, 100 ресурсів, 50 юзерів, 5 авто |
| **PRO** | 4999 грн/міс | + Smart Dispatch, GPS, Analytics, Fuel Anti-Fraud |
| **ENTERPRISE** | Custom | Безліміти + 24/7 + SLA |

### 🆓 BASIC (Безкоштовний)
- **Ліміти:** 10 складів | 100 товарів | 50 користувачів | 5 авто
- **Функції:** Базовий облік, ручні рейси, аудит 30 днів

### 💎 PRO (4999 грн/міс)
- **Ліміти:** 100 складів | 1000 товарів | 500 користувачів | 50 авто
- **Функції:** Smart Dispatch, GPS Tracking, Advanced Analytics, 
  Predictive Maintenance, Fuel Anti-Fraud, Excel import/export

### 🏢 ENTERPRISE (Індивідуально)
- **Ліміти:** ♾️ Безлімітні
- **Функції:** Все з PRO + підтримка 24/7 + SLA + персональний менеджер
```

---

## ✅ Результат

### Бізнес-логіка
```
┌─────────────────────────────────────┐
│  1. Нова організація реєструється   │
│     └─> Автоматично BASIC (було FREE)│
│                                      │
│  2. Досягає лімітів (10 складів)    │
│     └─> Запитує апгрейд             │
│                                      │
│  3. Варіанти оновлення:             │
│     • PRO: 4999 грн/міс (преміум)   │
│     • ENTERPRISE: custom (безліміт) │
└─────────────────────────────────────┘
```

### Порівняння до/після

| Аспект | До (FREE + BASIC) | Після (тільки BASIC) |
|--------|-------------------|----------------------|
| **Безкоштовних планів** | 2 (заплутано) | 1 (чітко) |
| **Стартовий тариф** | FREE (1 склад!) | BASIC (10 складів) |
| **Конверсія в платників** | Складна | Проста (BASIC → PRO) |
| **Логіка signup** | FREE за замовчуванням | BASIC за замовчуванням |
| **Зрозумілість для клієнтів** | ⚠️ Низька | ✅ Висока |

---

## 🧪 Тестування

### Перевірити:
1. ✅ **Міграція БД** — всі FREE тенанти стали BASIC
2. ✅ **Signup** — нові організації отримують BASIC
3. ✅ **Ліміти** — працюють для BASIC (10 складів, 100 товарів)
4. ✅ **UI** — `/billing` показує 3 картки (BASIC/PRO/ENTERPRISE)
5. ✅ **PlatformAdmin** — dropdown має тільки BASIC/PRO/ENTERPRISE
6. ✅ **Middleware** — перевірка тарифів працює без FREE

### Команди для тестування:
```bash
# 1. Перезапустити БД (міграція автоматично виконається)
docker-compose down
docker-compose up -d postgres

# 2. Перевірити дані
docker exec -it diploma-postgres-1 psql -U postgres -d Omnilog
\c Omnilog
SELECT name, subscription_tier FROM tenants;
SELECT subscription_tier, COUNT(*) FROM units GROUP BY subscription_tier;

# 3. Тестовий signup
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "organization_name": "Test Org",
    "slug": "testorg",
    "owner_email": "test@example.com",
    "owner_password": "Test123!",
    "owner_full_name": "Test User"
  }'

# 4. Перевірити тариф
curl http://localhost:8080/api/platform/tenants \
  -H "Authorization: Bearer <SYSTEM_ADMIN_TOKEN>"
```

---

## 📝 Висновок

Система тепер має **чітку тарифну градацію**:
- 🆓 **BASIC** — для старту (безкоштовно, але обмежено)
- 💎 **PRO** — для бізнесу (платно, з AI-фічами)
- 🏢 **ENTERPRISE** — для корпорацій (індивідуально, безліміт)

**Переваги:**
- ✅ Простіше для клієнтів (немає плутанини з FREE/BASIC)
- ✅ Логічна структура (freemium → premium → enterprise)
- ✅ Кращий user experience на сторінці `/billing`
- ✅ Відповідає стандартам SaaS-бізнесу

**Створено:** 23.04.2026  
**Автор:** GitHub Copilot  
**Версія системи:** v2.1
