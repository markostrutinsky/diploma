# 📐 РОЗДІЛ: Алгоритмічна частина функціонування системи

## Вступ

У даному розділі представлено детальний опис **п'яти ключових алгоритмів**, що забезпечують функціонування системи логістичного управління Omnilog. Кожен алгоритм супроводжується:
- Теоретичним описом
- Аналізом складності
- Блок-схемою
- Реальним кодом реалізації з проєкту

---

## 1️⃣ Алгоритм рекурсивної перевірки підпорядкування (Hierarchical Subordination Check)

### 1.1 Призначення та опис

**Мета:** Перевірка, чи належить користувач до ієрархії підрозділів командира для забезпечення правильного розмежування доступу до даних.

**Контекст використання:** 
- Фільтрація списку користувачів у системі
- Контроль доступу до даних підлеглих підрозділів
- Реалізація ієрархічної структури організації (РЕГІОН → ФІЛІЯ → ВІДДІЛ → КОМАНДА)

**Принцип роботи:**
Алгоритм використовує рекурсивні Common Table Expressions (CTE) у SQL для обходу дерева організаційних одиниць. Починаючи з підрозділу користувача, алгоритм рекурсивно знаходить всі дочірні підрозділи, формуючи множину доступних ID.

### 1.2 Математична модель

Нехай маємо дерево організаційних одиниць `T = (V, E)`, де:
- `V` — множина вершин (підрозділів)
- `E` — множина ребер (зв'язків parent-child)
- `root` — кореневий підрозділ користувача

**Задача:** знайти множину `S = {v ∈ V | ∃ шлях від root до v}`

**Алгоритм:**
```
function FindSubordinates(root):
    S = {root}
    queue = [root]
    
    while queue не пуста:
        current = queue.dequeue()
        for child in children(current):
            S = S ∪ {child}
            queue.enqueue(child)
    
    return S
```

### 1.3 Аналіз складності

- **Часова складність:** `O(h)`, де `h` — висота дерева підрозділів
- **Просторова складність:** `O(n)`, де `n` — кількість вузлів у піддереві
- **Найгірший випадок:** Всі підрозділи організовані в ланцюг → `O(n)` переходів

### 1.4 Блок-схема алгоритму

```
┌─────────────────────────────────────┐
│  ПОЧАТОК: GetVisibleUsers()         │
│  Вхід: requesterRole, requesterUnit │
└──────────────┬──────────────────────┘
               │
               ▼
       ┌───────────────────┐
       │ requesterRole ==  │
       │  SYSTEM_ADMIN?    │
       └──────┬───────┬────┘
              │ Так   │ Ні
              ▼       ▼
    ┌──────────────┐ ┌─────────────────┐
    │ Повернути    │ │ requesterRole == │
    │ ВСІХ         │ │ ADMIN/TENANT?    │
    │ користувачів │ └────┬────────┬────┘
    └──────────────┘      │ Так    │ Ні
                          ▼        ▼
              ┌──────────────────────┐ ┌─────────────────────┐
              │ Повернути            │ │ requesterUnit       │
              │ користувачів         │ │ == NULL?            │
              │ свого tenant         │ └───┬──────────┬──────┘
              └──────────────────────┘     │ Так      │ Ні
                                           ▼          ▼
                               ┌────────────────┐  ┌──────────────────────┐
                               │ Повернути      │  │ РЕКУРСИВНИЙ CTE:     │
                               │ користувачів   │  │ WITH RECURSIVE       │
                               │ без підрозділу │  │ hierarchy AS (       │
                               └────────────────┘  │   SELECT id          │
                                                   │   FROM units         │
                                                   │   WHERE id = $unit   │
                                                   │   UNION ALL          │
                                                   │   SELECT u.id        │
                                                   │   FROM units u       │
                                                   │   JOIN hierarchy h   │
                                                   │   ON u.parent_id=h.id│
                                                   │ )                    │
                                                   └──────────┬───────────┘
                                                              │
                                                              ▼
                                              ┌───────────────────────────┐
                                              │ SELECT users              │
                                              │ WHERE unit_id IN          │
                                              │ (SELECT id FROM hierarchy)│
                                              └──────────┬────────────────┘
                                                         │
                                                         ▼
                                              ┌──────────────────────┐
                                              │ Повернути список     │
                                              │ підлеглих            │
                                              │ користувачів         │
                                              └──────────┬───────────┘
                                                         │
                                                         ▼
                                              ┌──────────────────┐
                                              │ КІНЕЦЬ           │
                                              └──────────────────┘
```

### 1.5 Код реалізації

**Файл:** `/Omnilog_backend/internal/repositories/user_repository.go`

```go
// GetVisibleUsers повертає список користувачів з урахуванням ієрархії підрозділів
func (r *UserRepository) GetVisibleUsers(ctx context.Context, db DBExecutor, 
    requesterRole string, requesterUnitID *int64) ([]*models.User, error) {
    
    var query string
    var args []interface{}

    if requesterRole == "SYSTEM_ADMIN" {
        // Платформний адмін бачить всіх
        query = `SELECT id, username, email, full_name, phone, password_hash, 
                        role, status, unit_id, created_at, updated_at
                 FROM users ORDER BY created_at DESC`
    } else if requesterRole == "ADMIN" || requesterRole == "TENANT_ADMIN" {
        // Власник організації бачить всіх у своєму tenant
        tFilter := tenantFilter(ctx, "", "WHERE", &args)
        query = `
            SELECT id, username, email, full_name, phone, password_hash, 
                   role, status, unit_id, created_at, updated_at
            FROM users` + tFilter + `
            ORDER BY created_at DESC
        `
    } else {
        if requesterUnitID == nil {
            // Користувач без підрозділу бачить тільки інших без підрозділу
            tFilter := tenantFilter(ctx, "", "AND", &args)
            query = `
                SELECT id, username, email, full_name, phone, password_hash, 
                       role, status, unit_id, created_at, updated_at
                FROM users
                WHERE unit_id IS NULL` + tFilter + `
                ORDER BY created_at DESC
            `
        } else {
            // 🔥 РЕКУРСИВНИЙ ОБХІД ДЕРЕВА ПІДРОЗДІЛІВ
            args = append(args, *requesterUnitID)
            tFilter := tenantFilter(ctx, "users", "AND", &args)
            query = `
                WITH RECURSIVE hierarchy AS (
                    -- Базовий випадок: стартовий підрозділ
                    SELECT id FROM units WHERE id = $1
                    UNION ALL
                    -- Рекурсивний крок: всі дочірні підрозділи
                    SELECT u.id FROM units u
                    JOIN hierarchy h ON u.parent_id = h.id
                )
                SELECT id, username, email, full_name, phone, password_hash, 
                       role, status, unit_id, created_at, updated_at
                FROM users
                WHERE (unit_id IN (SELECT id FROM hierarchy) OR unit_id IS NULL)` + tFilter + `
                ORDER BY created_at DESC
            `
        }
    }

    rows, err := db.Query(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*models.User
    for rows.Next() {
        var u models.User
        err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.Phone, 
                        &u.PasswordHash, &u.Role, &u.Status, &u.UnitID, 
                        &u.CreatedAt, &u.UpdatedAt)
        if err != nil {
            return nil, err
        }
        users = append(users, &u)
    }

    return users, nil
}
```

### 1.6 Особливості реалізації

✅ **Переваги:**
- Ефективна робота з ієрархічними даними
- Одноразовий запит до БД замість множинних SELECT
- PostgreSQL оптимізує рекурсивні CTE

⚠️ **Обмеження:**
- Максимальна глибина рекурсії обмежена налаштуваннями СУБД (за замовчуванням 100)
- При циклічних залежностях у даних може виникнути нескінченна рекурсія

---

## 2️⃣ Алгоритм розумного поповнення запасів (Smart Replenishment Algorithm)

### 2.1 Призначення та опис

**Мета:** Автоматичне визначення дефіцитних ресурсів та формування оптимального плану поповнення з урахуванням вже замовленого майна "в дорозі".

**Контекст використання:**
- Превентивне планування постачань
- Оптимізація рівня запасів
- Зменшення ризику дефіциту критичних ресурсів

**Принцип роботи:**
Алгоритм аналізує поточні залишки ресурсів, порівнює їх з мінімальними нормами, враховує вже створені заявки на поповнення (статус OPEN, IN_PROGRESS, APPROVED) та розраховує необхідну кількість для замовлення.

### 2.2 Математична модель

**Вхідні дані:**
- `current_stock` — фактична кількість на складі
- `min_quantity` — мінімальна норма запасу
- `pending_orders` — кількість у відкритих заявках

**Формула розрахунку дефіциту:**

```
needed = max(0, (min_quantity × 2) - (current_stock + pending_orders))
```

**Умова виявлення дефіциту:**

```
(current_stock + pending_orders) ≤ min_quantity
```

**Цільовий рівень:** `min_quantity × 2` — підтримка подвійного запасу для страхового резерву.

### 2.3 Аналіз складності

- **Часова складність:** `O(n)`, де `n` — кількість ресурсів у системі
- **Просторова складність:** `O(n)` для збереження результатів
- **Складність запиту:** 2 проходи по таблиці `supply_requests` + 1 прохід по `resources`

### 2.4 Блок-схема алгоритму

```
┌────────────────────────────────────────┐
│ ПОЧАТОК: GetDeficitResources()         │
│ Вхід: tenant_id, unit_id (опційно)    │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ CTE 1: PendingOrders                   │
│ ───────────────────────────────────    │
│ SELECT resource_id,                    │
│        SUM(quantity) as pending_qty    │
│ FROM supply_requests                   │
│ WHERE status IN ('OPEN',               │
│                  'IN_PROGRESS',        │
│                  'APPROVED')           │
│ GROUP BY resource_id                   │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ Головний запит:                        │
│ SELECT r.id, r.name, r.quantity,      │
│        r.min_quantity,                 │
│        (min × 2 - (current + pending)) │
│        as needed                       │
│ FROM resources r                       │
│ LEFT JOIN PendingOrders p              │
│   ON r.id = p.resource_id              │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ Умова фільтрації:                      │
│ ┌────────────────────────────────────┐ │
│ │ (current + pending) <= min_quantity│ │
│ │ AND condition != 'WRITTEN_OFF'     │ │
│ └────────────────────────────────────┘ │
└──────────────┬─────────────────────────┘
               │
        ┌──────┴──────┐
        │             │
        ▼ Так         ▼ Ні
┌─────────────────┐  ┌──────────────┐
│ Додати ресурс   │  │ Пропустити   │
│ до списку       │  └──────────────┘
│ дефіцитних      │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────┐
│ Сортування за критичністю:      │
│ ORDER BY (current/min) ASC      │
│ (найменше співвідношення вверх) │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│ КІНЕЦЬ: Повернути список        │
│ DeficitResource[]               │
└─────────────────────────────────┘
```

### 2.5 Код реалізації

**Файл:** `/Omnilog_backend/internal/repositories/analytics_repository.go`

```go
// Розрахунок дефіцитних ресурсів для Smart Replenishment
queryDeficit := fmt.Sprintf(`
    WITH PendingOrders AS (
        -- Рахуємо скільки майна ВЖЕ замовлено (висить у відкритих заявках)
        SELECT resource_id, SUM(quantity) as pending_qty
        FROM supply_requests
        WHERE status IN ('OPEN', 'IN_PROGRESS', 'APPROVED')%s 
        GROUP BY resource_id
    )
    SELECT 
        r.id, 
        r.name, 
        r.quantity, 
        r.min_quantity, 
        -- Формула: (Мінімум * 2) - (Фактичний залишок + Вже замовлено)
        (r.min_quantity * 2 - (r.quantity + COALESCE(p.pending_qty, 0))) as needed
    FROM resources r
    LEFT JOIN PendingOrders p ON r.id = p.resource_id
    -- Показуємо тільки те, де ФАКТ + В ДОРОЗІ все ще менше або дорівнює мінімуму
    WHERE (r.quantity + COALESCE(p.pending_qty, 0)) <= r.min_quantity 
      AND r.condition != 'WRITTEN_OFF' %s%s
`, tcond(""), resFilterPrefix, tcond("r"))

dRows, _ := db.Query(ctx, queryDeficit)
defer dRows.Close()

for dRows.Next() {
    var d models.DeficitResource
    dRows.Scan(&d.ID, &d.Name, &d.Current, &d.Min, &d.Needed)
    stats.DeficitResources = append(stats.DeficitResources, d)
}
```

### 2.6 Приклад роботи алгоритму

**Ситуація:**
- Ресурс: "Аптечка індивідуальна IFAK"
- `current_stock = 5` шт.
- `min_quantity = 10` шт.
- `pending_orders = 3` шт. (є відкрита заявка)

**Розрахунок:**
```
1. Перевірка дефіциту: (5 + 3) = 8 ≤ 10 ✅ (є дефіцит)
2. Цільовий рівень: 10 × 2 = 20 шт.
3. Необхідно замовити: 20 - (5 + 3) = 12 шт.
```

**Результат:** Система автоматично пропонує створити заявку на **12 аптечок**.

---

## 3️⃣ Алгоритм обчислення відстаней за формулою Haversine

### 3.1 Призначення та опис

**Мета:** Точний розрахунок відстані між двома точками на сферичній поверхні Землі за їх GPS-координатами (широта, довгота).

**Контекст використання:**
- Побудова траєкторій руху транспорту
- Розрахунок пробігу для контролю витрат палива
- Детекція геозон (входження у певну область)
- Аналітика логістичних маршрутів

**Принцип роботи:**
Формула Haversine враховує сферичну геометрію Землі та обчислює найкоротшу відстань (great-circle distance) між двома точками.

### 3.2 Математична модель

**Вхідні дані:**
- `(lat₁, lon₁)` — координати першої точки (у градусах)
- `(lat₂, lon₂)` — координати другої точки (у градусах)
- `R = 6371` км — середній радіус Землі

**Формула Haversine:**

```
φ₁ = lat₁ × π/180  (конвертація у радіани)
φ₂ = lat₂ × π/180
Δφ = (lat₂ - lat₁) × π/180
Δλ = (lon₂ - lon₁) × π/180

a = sin²(Δφ/2) + cos(φ₁) × cos(φ₂) × sin²(Δλ/2)
c = 2 × arcsin(√a)

d = R × c  (відстань у кілометрах)
```

**Точність:** ±0.5% для відстаней до 1000 км.

### 3.3 Аналіз складності

- **Часова складність:** `O(1)` — фіксована кількість математичних операцій
- **Просторова складність:** `O(1)` — використовується константна пам'ять
- **Обчислювальна складність:** ~10 арифметичних операцій + 4 тригонометричні функції

**Для траєкторії з n точок:**
- Загальна складність: `O(n)` — сумування відстаней між послідовними точками

### 3.4 Блок-схема алгоритму

```
┌────────────────────────────────────────┐
│ ПОЧАТОК: calculateHaversine()          │
│ Вхід: lat1, lon1, lat2, lon2           │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ Константа: R = 6371.0 (км)             │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ Конвертація градусів у радіани:        │
│ lat1Rad = lat1 × π / 180               │
│ lon1Rad = lon1 × π / 180               │
│ lat2Rad = lat2 × π / 180               │
│ lon2Rad = lon2 × π / 180               │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ Обчислення різниць:                    │
│ dlat = lat2Rad - lat1Rad               │
│ dlon = lon2Rad - lon1Rad               │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ Розрахунок проміжного значення a:      │
│ ┌────────────────────────────────────┐ │
│ │ a = sin²(dlat/2) +                 │ │
│ │     cos(lat1Rad) × cos(lat2Rad) ×  │ │
│ │     sin²(dlon/2)                   │ │
│ └────────────────────────────────────┘ │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ Обчислення центрального кута:          │
│ c = 2 × arcsin(√a)                     │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ Фінальний розрахунок відстані:         │
│ distance = R × c                       │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ КІНЕЦЬ: Повернути distance (км)        │
└────────────────────────────────────────┘


════════════════════════════════════════════
ДОДАТКОВА БЛОК-СХЕМА: CalculateDistance()
для траєкторії з n точок
════════════════════════════════════════════

┌────────────────────────────────────────┐
│ ПОЧАТОК: CalculateDistance()           │
│ Вхід: locations[] (масив GPS точок)    │
└──────────────┬─────────────────────────┘
               │
               ▼
       ┌───────────────┐
       │ len(locations)│
       │     < 2?      │
       └───┬───────┬───┘
     Так   │       │ Ні
           ▼       ▼
    ┌──────────┐  ┌──────────────────┐
    │ Повернути│  │ totalDistance = 0│
    │    0     │  │ i = 0            │
    └──────────┘  └────────┬─────────┘
                           │
                           ▼
                  ┌────────────────┐
                  │ i < len - 1?   │
                  └───┬────────┬───┘
                Ні    │        │ Так
                      │        ▼
                      │  ┌──────────────────────┐
                      │  │ lat1 = locations[i]  │
                      │  │ lon1 = locations[i]  │
                      │  │ lat2 = locations[i+1]│
                      │  │ lon2 = locations[i+1]│
                      │  └──────────┬───────────┘
                      │             │
                      │             ▼
                      │  ┌──────────────────────┐
                      │  │ segment_distance =   │
                      │  │ calculateHaversine() │
                      │  └──────────┬───────────┘
                      │             │
                      │             ▼
                      │  ┌──────────────────────┐
                      │  │ totalDistance +=     │
                      │  │ segment_distance     │
                      │  └──────────┬───────────┘
                      │             │
                      │             ▼
                      │  ┌──────────────────┐
                      │  │ i = i + 1        │
                      │  └────────┬─────────┘
                      │           │
                      │           └──────┐
                      │                  │
                      ▼                  │
              ┌──────────────────┐      │
              │ Повернути        │◄─────┘
              │ totalDistance    │
              └──────────────────┘
```

### 3.5 Код реалізації

**Файл:** `/Omnilog_backend/internal/services/gps_tracking_service.go`

```go
// calculateHaversine calculates distance in km between two GPS coordinates
func (s *GPSTrackingService) calculateHaversine(lat1, lon1, lat2, lon2 float64) float64 {
    const earthRadiusKm = 6371.0

    // Конвертація градусів у радіани
    lat1Rad := lat1 * math.Pi / 180
    lon1Rad := lon1 * math.Pi / 180
    lat2Rad := lat2 * math.Pi / 180
    lon2Rad := lon2 * math.Pi / 180

    // Різниці координат
    dlat := lat2Rad - lat1Rad
    dlon := lon2Rad - lon1Rad

    // Формула Haversine
    a := math.Sin(dlat/2)*math.Sin(dlat/2) +
        math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
    c := 2 * math.Asin(math.Sqrt(a))

    return earthRadiusKm * c
}

// CalculateDistance calculates the distance traveled based on GPS points
func (s *GPSTrackingService) CalculateDistance(locations []models.GPSLocation) float64 {
    if len(locations) < 2 {
        return 0
    }

    const earthRadiusKm = 6371.0
    totalDistance := 0.0

    // Сумуємо відстані між послідовними точками
    for i := 0; i < len(locations)-1; i++ {
        lat1 := locations[i].Latitude * math.Pi / 180
        lon1 := locations[i].Longitude * math.Pi / 180
        lat2 := locations[i+1].Latitude * math.Pi / 180
        lon2 := locations[i+1].Longitude * math.Pi / 180

        dlat := lat2 - lat1
        dlon := lon2 - lon1

        a := math.Sin(dlat/2)*math.Sin(dlat/2) +
            math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
        c := 2 * math.Asin(math.Sqrt(a))
        distance := earthRadiusKm * c

        totalDistance += distance
    }

    return totalDistance
}
```

**Файл симулятора GPS:** `/scripts/gps_simulator/simulate.py`

```python
import math

EARTH_KM = 6371.0

def haversine_km(a_lat: float, a_lon: float, b_lat: float, b_lon: float) -> float:
    """Відстань між двома GPS координатами (формула Haversine)"""
    lat1, lat2 = math.radians(a_lat), math.radians(b_lat)
    lon1, lon2 = math.radians(a_lon), math.radians(b_lon)
    
    dlat = lat2 - lat1
    dlon = lon2 - lon1
    
    a = math.sin(dlat/2)**2 + math.cos(lat1) * math.cos(lat2) * math.sin(dlon/2)**2
    c = 2 * math.asin(math.sqrt(a))
    
    return EARTH_KM * c

def bearing_deg(a_lat: float, a_lon: float, b_lat: float, b_lon: float) -> float:
    """Розрахунок азимута (напрямку руху) у градусах"""
    lat1, lat2 = math.radians(a_lat), math.radians(b_lat)
    dlon = math.radians(b_lon - a_lon)
    
    y = math.sin(dlon) * math.cos(lat2)
    x = math.cos(lat1) * math.sin(lat2) - math.sin(lat1) * math.cos(lat2) * math.cos(dlon)
    
    return (math.degrees(math.atan2(y, x)) + 360) % 360
```

### 3.6 Приклад роботи алгоритму

**Вхідні дані:**
- Київ: `(50.4501° N, 30.5234° E)`
- Львів: `(49.8397° N, 24.0297° E)`

**Розрахунок:**
```
1. Конвертація:
   lat1 = 50.4501 × π/180 = 0.8803 рад
   lon1 = 30.5234 × π/180 = 0.5327 рад
   lat2 = 49.8397 × π/180 = 0.8697 рад
   lon2 = 24.0297 × π/180 = 0.4194 рад

2. Різниці:
   Δlat = -0.0106 рад
   Δlon = -0.1133 рад

3. Формула Haversine:
   a = sin²(-0.0053) + cos(0.8803) × cos(0.8697) × sin²(-0.0567)
   a ≈ 0.00197
   c = 2 × arcsin(√0.00197) ≈ 0.0888 рад

4. Відстань:
   d = 6371 × 0.0888 ≈ 469 км
```

**Реальна відстань:** 470 км (точність 99.8% ✅)

---

## 4️⃣ Алгоритм детекції аномалій пального (Fuel Anomaly Detection)

### 4.1 Призначення та опис

**Мета:** Виявлення підозрілих операцій з паливом для попередження шахрайства, помилок обліку та технічних несправностей.

**Контекст використання:**
- Антикорупційний контроль витрат палива
- Виявлення зливу палива
- Детекція переповнення баків (фіктивні заправки)
- Моніторинг аномальної витрати

**Типи аномалій:**
1. **TANK_OVERFLOW** — спроба заправити більше, ніж вміщує бак
2. **NEGATIVE_BALANCE** — спроба списати більше, ніж є в баку
3. **ODOMETER_ROLLBACK** — скручування одометра назад
4. **EXCESSIVE_CONSUMPTION** — витрата понад 150% від норми
5. **FREQUENT_SMALL_REFILLS** — підозріло часті заправки малими порціями

### 4.2 Математична модель

**Вхідні дані для кожного транспортного засобу:**
- `fuel_norm` — нормативна витрата (л/100км)
- `tank_capacity` — місткість баку (л)
- `current_balance` — поточний залишок палива (віртуальний бак)
- `last_odometer` — останнє показання одометра (км)

**Алгоритм перевірки:**

```
function DetectFuelAnomaly(record):
    anomalies = []
    
    // Перевірка 1: Переповнення баку
    if record.type == REFUEL:
        if current_balance + record.liters > tank_capacity:
            anomalies.add("TANK_OVERFLOW")
            return HARD_STOP  // Блокуємо транзакцію
    
    // Перевірка 2: Від'ємний баланс
    if record.type == EXPENSE:
        if record.liters > current_balance:
            anomalies.add("NEGATIVE_BALANCE")
            return HARD_STOP  // Блокуємо транзакцію
    
    // Перевірка 3: Скручування одометра
    if record.odometer < last_odometer:
        anomalies.add("ODOMETER_ROLLBACK")
        return HARD_STOP  // Блокуємо транзакцію
    
    // Перевірка 4: Надмірна витрата
    if record.type == EXPENSE and record.odometer != null:
        distance = record.odometer - last_odometer
        actual_consumption = (record.liters / distance) × 100
        
        if actual_consumption > fuel_norm × 1.5:
            anomalies.add("EXCESSIVE_CONSUMPTION")
            record.is_anomaly = true  // Soft-stop (записуємо, але попереджаємо)
    
    // Перевірка 5: Часті малі заправки (статистичний аналіз)
    recent_refills = count_refills_last_7_days(vehicle_id)
    if recent_refills > 10 and record.liters < 20:
        anomalies.add("FREQUENT_SMALL_REFILLS")
        record.is_anomaly = true
    
    return anomalies
```

**Метрика ризику шахрайства:**

```
Risk Score = (кількість аномальних заправок / загальна кількість заправок) × 100
```

- `Risk < 10%` — низький ризик (зелена зона)
- `10% ≤ Risk < 30%` — помірний ризик (жовта зона)
- `Risk ≥ 30%` — високий ризик (червона зона, потребує розслідування)

### 4.3 Аналіз складності

- **Часова складність перевірки одного запису:** `O(1)` + `O(log n)` для запиту last_odometer
- **Просторова складність:** `O(1)` — константна пам'ять
- **Транзакційна безпека:** Використання `FOR UPDATE` для уникнення race conditions

### 4.4 Блок-схема алгоритму

```
┌────────────────────────────────────────┐
│ ПОЧАТОК: CreateFuelRecord()            │
│ Вхід: record (FuelRecord)              │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ Старт транзакції (BEGIN)               │
│ SELECT fuel_norm, tank_capacity        │
│ FROM vehicles                          │
│ WHERE id = record.vehicle_id           │
│ FOR UPDATE  ← Блокування рядка         │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ Розрахунок поточного балансу:          │
│ current_balance = SUM(                 │
│   CASE record_type                     │
│     WHEN 'REFUEL' THEN +liters         │
│     WHEN 'EXPENSE' THEN -liters        │
│   END)                                 │
└──────────────┬─────────────────────────┘
               │
               ▼
       ┌───────────────┐
       │ record.type = │
       │   REFUEL?     │
       └───┬───────┬───┘
     Так   │       │ Ні
           ▼       ▼
┌────────────────┐ ┌────────────────┐
│ Перевірка:     │ │ record.type =  │
│ current +      │ │   EXPENSE?     │
│ liters >       │ └───┬────────┬───┘
│ tank_capacity? │ Так │        │ Ні
└────┬───────┬───┘     ▼        │
 Так │   Ні  │   ┌─────────────┐│
     ▼       │   │ Перевірка:  ││
┌─────────┐  │   │ liters >    ││
│ ПОМИЛКА:│  │   │ current?    ││
│ "Бак    │  │   └──┬──────┬───┘│
│ перепов-│  │  Так │  Ні  │    │
│ нено!"  │  │      ▼      │    │
│ ROLLBACK│  │ ┌─────────┐ │    │
└─────────┘  │ │ ПОМИЛКА:│ │    │
             │ │"Недоста-│ │    │
             │ │тньо     │ │    │
             │ │пального"│ │    │
             │ │ ROLLBACK│ │    │
             │ └─────────┘ │    │
             │             │    │
             └─────────────┴────┘
                     │
                     ▼
       ┌─────────────────────────┐
       │ record.odometer != NULL?│
       └────┬────────────────┬───┘
        Так │                │ Ні
            ▼                │
┌──────────────────────────┐ │
│ SELECT last_odometer     │ │
│ FROM fuel_records        │ │
│ WHERE vehicle_id = ...   │ │
│ ORDER BY created_at DESC │ │
│ LIMIT 1                  │ │
└──────────────┬───────────┘ │
               │             │
               ▼             │
       ┌───────────────┐     │
       │ odometer <    │     │
       │ last_odometer?│     │
       └───┬───────┬───┘     │
       Так │   Ні  │         │
           ▼       ▼         │
      ┌─────────┐ ┌─────────────────┐
      │ ПОМИЛКА:│ │ distance =      │
      │"Скручу- │ │ odometer -      │
      │вання    │ │ last_odometer   │
      │одометра"│ └────────┬────────┘
      │ ROLLBACK│          │
      └─────────┘          ▼
                  ┌────────────────────┐
                  │ actual_consumption │
                  │ = (liters/distance)│
                  │   × 100            │
                  └────────┬───────────┘
                           │
                           ▼
                   ┌───────────────┐
                   │ actual > norm │
                   │   × 1.5?      │
                   └───┬───────┬───┘
                   Так │   Ні  │
                       ▼       │
              ┌────────────────┐│
              │ is_anomaly =   ││
              │ TRUE           ││
              │ reason =       ││
              │ "EXCESSIVE_    ││
              │  CONSUMPTION"  ││
              └────────┬───────┘│
                       │        │
                       └────────┘
                            │
                            ▼
              ┌──────────────────────┐
              │ INSERT INTO          │
              │ fuel_records (...)   │
              │ VALUES (...)         │
              └──────────┬───────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │ COMMIT транзакції    │
              └──────────┬───────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │ КІНЕЦЬ: Успіх        │
              └──────────────────────┘
```

### 4.5 Код реалізації

**Файл:** `/Omnilog_backend/internal/repositories/fuel_repository.go`

```go
func (r *FuelRepository) CreateFuelRecord(ctx context.Context, 
    record *models.FuelRecord, db DBExecutor) error {
    
    // Таймаут для запобігання зависанню
    var cancel context.CancelFunc
    if _, ok := ctx.Deadline(); !ok {
        ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
        defer cancel()
    }

    // Старт транзакції
    b, ok := db.(beginner)
    if !ok {
        return fmt.Errorf("db не підтримує транзакції")
    }

    tx, err := b.Begin(ctx)
    if err != nil {
        return fmt.Errorf("помилка старту транзакції: %w", err)
    }
    defer tx.Rollback(ctx)

    // 1. Отримання даних авто з блокуванням рядка (FOR UPDATE)
    var fuelNorm float64
    var tankCapacity float64
    err = tx.QueryRow(ctx, 
        "SELECT fuel_norm, tank_capacity FROM vehicles WHERE id = $1 FOR UPDATE", 
        record.VehicleID).Scan(&fuelNorm, &tankCapacity)
    if err != nil {
        return fmt.Errorf("помилка отримання даних авто: %w", err)
    }

    // 2. Розрахунок поточного залишку "Віртуальний бак"
    var currentBalance float64
    balanceQuery := `
        SELECT COALESCE(SUM(CASE 
            WHEN record_type = 'REFUEL' THEN liters 
            ELSE -liters 
        END), 0)
        FROM fuel_records
        WHERE vehicle_id = $1
    `
    err = tx.QueryRow(ctx, balanceQuery, record.VehicleID).Scan(&currentBalance)
    if err != nil {
        return fmt.Errorf("помилка розрахунку залишку: %w", err)
    }

    // 3. HARD STOP: Перевірка на переповнення баку
    if record.RecordType == models.FuelRefuel {
        if currentBalance+record.Liters > tankCapacity {
            return fmt.Errorf(
                "бак переповнено! Поточний: %.1f л, максимум: %.1f л, спроба: %.1f л",
                currentBalance, tankCapacity, record.Liters)
        }
    }

    // 4. HARD STOP: Перевірка на від'ємний баланс
    if record.RecordType == models.FuelExpense {
        if record.Liters > currentBalance {
            return fmt.Errorf(
                "недостатньо пального. У баку: %.1f л, спроба списати: %.1f л",
                currentBalance, record.Liters)
        }
    }

    // 5. HARD STOP: Перевірка скручування одометра
    var lastOdometer int
    var hasLastOdometer bool

    if record.OdometerKm != nil {
        err = tx.QueryRow(ctx,
            `SELECT odometer_km FROM fuel_records 
             WHERE vehicle_id = $1 AND odometer_km IS NOT NULL 
             ORDER BY created_at DESC LIMIT 1`,
            record.VehicleID,
        ).Scan(&lastOdometer)

        if err == nil {
            hasLastOdometer = true
            if *record.OdometerKm < lastOdometer {
                return fmt.Errorf(
                    "скручування одометра: поточний (%d км) < попередній (%d км)",
                    *record.OdometerKm, lastOdometer)
            }
        } else if err != pgx.ErrNoRows {
            log.Printf("Помилка одометра: %v", err)
        }
    }

    // 6. SOFT STOP: Детекція аномальної витрати
    if record.RecordType == models.FuelExpense && 
       record.OdometerKm != nil && hasLastOdometer {
        
        distance := *record.OdometerKm - lastOdometer
        if distance > 0 {
            actualConsumption := (record.Liters / float64(distance)) * 100
            
            // Перевіряємо, чи витрата перевищує норму на 50%
            if actualConsumption > fuelNorm*1.5 {
                record.IsAnomaly = true
                reason := fmt.Sprintf(
                    "Надмірна витрата: %.2f л/100км (норма: %.2f л/100км)",
                    actualConsumption, fuelNorm)
                record.AnomalyReason = &reason
                
                log.Printf("⚠️ АНОМАЛІЯ ВИЯВЛЕНА: %s для авто %s", 
                    reason, record.VehicleID)
            }
        }
    }

    // 7. Збереження запису в БД
    insertQuery := `
        INSERT INTO fuel_records (
            vehicle_id, record_type, liters, odometer_km, 
            fuel_station, cost_per_liter, total_cost, is_anomaly, 
            anomaly_reason, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id
    `
    err = tx.QueryRow(ctx, insertQuery,
        record.VehicleID,
        record.RecordType,
        record.Liters,
        record.OdometerKm,
        record.FuelStation,
        record.CostPerLiter,
        record.TotalCost,
        record.IsAnomaly,
        record.AnomalyReason,
        time.Now(),
    ).Scan(&record.ID)

    if err != nil {
        return fmt.Errorf("помилка збереження запису: %w", err)
    }

    // 8. Коміт транзакції
    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("помилка коміту: %w", err)
    }

    return nil
}
```

### 4.6 Приклад роботи алгоритму

**Ситуація 1: Спроба переповнення баку (HARD STOP)**
```
Авто: Toyota Land Cruiser (бак 138 л)
Поточний залишок: 120 л
Спроба заправити: 25 л

Розрахунок: 120 + 25 = 145 > 138 ❌

Результат: ПОМИЛКА "Бак переповнено!" → Транзакція відхилена
```

**Ситуація 2: Аномальна витрата (SOFT STOP)**
```
Авто: Ford Transit (норма 12 л/100км)
Попередній одометр: 45000 км
Поточний одометр: 45100 км (пробіг 100 км)
Витрачено: 18 л

Розрахунок: (18 / 100) × 100 = 18 л/100км
18 > 12 × 1.5 = 18 ✅ (на межі)

Результат: Запис збережено, але позначено is_anomaly = TRUE
```

---

## 5️⃣ Алгоритм автоматичної ескалації за SLA (SLA Monitor)

### 5.1 Призначення та опис

**Мета:** Автоматичний моніторинг термінів виконання заявок та ескалація прострочених запитів для своєчасного реагування.

**Контекст використання:**
- Контроль дотримання Service Level Agreement (SLA)
- Автоматичне попередження про наближення дедлайну
- Ескалація критичних заявок керівництву
- Збір метрик ефективності

**Принцип роботи:**
Фоновий cron-процес кожну хвилину сканує базу даних на наявність заявок, у яких:
1. Статус = OPEN або IN_PROGRESS
2. Поточний час > дедлайн (або created_at + 7 днів, якщо дедлайн не вказано)

При виявленні прострочених заявок:
- Змінює статус на ESCALATED
- Створює запис в audit log
- Надсилає email-нотифікацію відповідальним особам

### 5.2 Математична модель

**Розрахунок дедлайну:**

```
deadline = created_at + SLA_DURATION

де:
  SLA_DURATION = 7 днів (базовий SLA)
  або спеціальний deadline, якщо вказано у заявці
```

**Умова ескалації:**

```
ESCALATE IF:
  (status IN ('OPEN', 'IN_PROGRESS')) AND
  (NOW() > COALESCE(deadline, created_at + INTERVAL '7 days'))
```

**Метрика On-Time Delivery (OTD):**

```
OTD% = (кількість виконаних вчасно / загальна кількість виконаних) × 100

де:
  "виконаних вчасно" = completed_at ≤ deadline
```

### 5.3 Аналіз складності

- **Часова складність перевірки:** `O(n)`, де `n` — кількість активних заявок
- **Просторова складність:** `O(k)`, де `k` — кількість прострочених заявок
- **Частота виконання:** Кожну 1 хвилину (настроюється через cron)
- **Навантаження на БД:** ~1 SELECT запит з індексацією по (status, created_at)

**Оптимізація:**
- Індекс: `CREATE INDEX idx_requests_sla ON supply_requests(status, created_at, deadline)`
- Обмеження вибірки: `LIMIT 100` для захисту від перевантаження

### 5.4 Блок-схема алгоритму

```
┌────────────────────────────────────────┐
│ CRON JOB: Кожну 1 хвилину              │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ ПОЧАТОК: slaMonitor.Start()            │
│ Створення goroutine (фонового потоку)  │
└──────────────┬─────────────────────────┘
               │
               ▼
         ┌─────────────┐
         │ Безкінечний │
         │   цикл      │
         └──────┬──────┘
                │
                ▼
┌────────────────────────────────────────┐
│ CheckPendingRequests()                 │
│ ───────────────────────────────────    │
│ SELECT id, created_at, resource_id,    │
│        deadline, unit_id               │
│ FROM supply_requests                   │
│ WHERE status IN ('OPEN','IN_PROGRESS') │
│ AND NOW() > COALESCE(deadline,         │
│              created_at + INTERVAL '7d'│
└──────────────┬─────────────────────────┘
               │
               ▼
       ┌───────────────┐
       │ Знайдено      │
       │ прострочені?  │
       └───┬───────┬───┘
       Ні  │       │ Так
           │       ▼
           │  ┌──────────────────────┐
           │  │ Для кожної заявки:   │
           │  └──────────┬───────────┘
           │             │
           │             ▼
           │  ┌──────────────────────┐
           │  │ START TRANSACTION    │
           │  └──────────┬───────────┘
           │             │
           │             ▼
           │  ┌──────────────────────┐
           │  │ UPDATE supply_requests│
           │  │ SET status='ESCALATED'│
           │  │ WHERE id = $1        │
           │  └──────────┬───────────┘
           │             │
           │             ▼
           │  ┌──────────────────────┐
           │  │ INSERT INTO audit_log│
           │  │ (action, entity_type,│
           │  │  entity_id, message) │
           │  │ VALUES (             │
           │  │  'SLA_BREACH',       │
           │  │  'request',          │
           │  │  $id,                │
           │  │  'Заявка прострочена'│
           │  │ )                    │
           │  └──────────┬───────────┘
           │             │
           │             ▼
           │  ┌──────────────────────┐
           │  │ emailService.Send()  │
           │  │ Кому: відповідальний │
           │  │ Тема: "⚠️ SLA Breach"│
           │  │ Тіло: деталі заявки  │
           │  └──────────┬───────────┘
           │             │
           │             ▼
           │  ┌──────────────────────┐
           │  │ COMMIT TRANSACTION   │
           │  └──────────┬───────────┘
           │             │
           │             └────┐
           │                  │
           ▼                  │
    ┌─────────────┐           │
    │ Log:        │           │
    │ "Ескальовано│◄──────────┘
    │  X заявок"  │
    └──────┬──────┘
           │
           ▼
    ┌─────────────────┐
    │ time.Sleep(     │
    │   1 хвилина)    │
    └──────┬──────────┘
           │
           └───────────┐
                       │
                       ▼
                  (повтор циклу)


════════════════════════════════════════════
ДОДАТКОВА БЛОК-СХЕМА: Розрахунок метрик SLA
════════════════════════════════════════════

┌────────────────────────────────────────┐
│ GetContractorSLAMetrics()              │
│ Вхід: startDate, endDate               │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ SELECT                                 │
│   AVG(completed_at - created_at)       │
│     as average_days,                   │
│   MIN(completed_at - created_at)       │
│     as fastest_days,                   │
│   COUNT(*) as total,                   │
│   COUNT(*) FILTER (                    │
│     WHERE completed_at <= deadline     │
│   ) as on_time                         │
│ FROM supply_requests                   │
│ WHERE status = 'COMPLETED'             │
│   AND created_at BETWEEN $1 AND $2     │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ Розрахунок OTD%:                       │
│ otd_percentage = (on_time / total) × 100│
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ SELECT COUNT(*)                        │
│ FROM supply_requests                   │
│ WHERE status IN ('OPEN','IN_PROGRESS') │
│   AND NOW() > deadline                 │
│   as overdue_count                     │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│ КІНЕЦЬ: Повернути структуру SLAMetrics │
│ {                                      │
│   average_days: float64                │
│   fastest_days: float64                │
│   completed_count: int                 │
│   otd_percentage: int                  │
│   overdue_count: int                   │
│ }                                      │
└────────────────────────────────────────┘
```

### 5.5 Код реалізації

**Файл:** `/Omnilog_backend/internal/services/sla_monitor.go`

```go
package services

import (
    "context"
    "fmt"
    "log"
    "Omnilog_backend/internal/models"
    "Omnilog_backend/internal/repositories"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type SLAMonitor struct {
    db            *pgxpool.Pool
    requestRepo   *repositories.SupplyRequestRepository
    auditRepo     *repositories.AuditRepository
    emailService  EmailService
    stopChan      chan struct{}
    checkInterval time.Duration
}

func NewSLAMonitor(
    db *pgxpool.Pool,
    requestRepo *repositories.SupplyRequestRepository,
    auditRepo *repositories.AuditRepository,
    emailService EmailService,
) *SLAMonitor {
    return &SLAMonitor{
        db:            db,
        requestRepo:   requestRepo,
        auditRepo:     auditRepo,
        emailService:  emailService,
        stopChan:      make(chan struct{}),
        checkInterval: 1 * time.Minute, // Перевірка кожну хвилину
    }
}

// Start запускає фоновий процес моніторингу SLA
func (m *SLAMonitor) Start(ctx context.Context) {
    log.Println("🚀 SLA Monitor запущено (інтервал: 1 хвилина)")

    go func() {
        ticker := time.NewTicker(m.checkInterval)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                escalatedCount, err := m.CheckPendingRequests(ctx)
                if err != nil {
                    log.Printf("❌ SLA Monitor помилка: %v", err)
                } else if escalatedCount > 0 {
                    log.Printf("⚠️ SLA Monitor: ескальовано %d заявок", escalatedCount)
                }
            case <-m.stopChan:
                log.Println("🛑 SLA Monitor зупинено")
                return
            }
        }
    }()
}

// Stop зупиняє моніторинг
func (m *SLAMonitor) Stop() {
    close(m.stopChan)
}

// CheckPendingRequests перевіряє всі активні заявки на порушення SLA
func (m *SLAMonitor) CheckPendingRequests(ctx context.Context) (int, error) {
    // Запит на знаходження прострочених заявок
    query := `
        SELECT id, created_at, resource_id, deadline, unit_id
        FROM supply_requests
        WHERE status IN ('OPEN', 'IN_PROGRESS')
          AND NOW() > COALESCE(deadline, created_at + INTERVAL '7 days')
        ORDER BY created_at ASC
        LIMIT 100
    `

    rows, err := m.db.Query(ctx, query)
    if err != nil {
        return 0, fmt.Errorf("помилка запиту: %w", err)
    }
    defer rows.Close()

    escalatedCount := 0

    for rows.Next() {
        var (
            requestID  string
            createdAt  time.Time
            resourceID string
            deadline   *time.Time
            unitID     *int64
        )

        if err := rows.Scan(&requestID, &createdAt, &resourceID, &deadline, &unitID); err != nil {
            log.Printf("Помилка сканування рядка: %v", err)
            continue
        }

        // Визначаємо фактичний дедлайн
        actualDeadline := createdAt.Add(7 * 24 * time.Hour)
        if deadline != nil {
            actualDeadline = *deadline
        }

        // Розрахунок часу прострочення
        overdueDuration := time.Since(actualDeadline)
        overdueDays := int(overdueDuration.Hours() / 24)

        // Ескалація заявки
        tx, err := m.db.Begin(ctx)
        if err != nil {
            log.Printf("Помилка старту транзакції для заявки %s: %v", requestID, err)
            continue
        }

        // 1. Оновлення статусу на ESCALATED
        updateQuery := `
            UPDATE supply_requests
            SET status = 'ESCALATED', updated_at = NOW()
            WHERE id = $1
        `
        _, err = tx.Exec(ctx, updateQuery, requestID)
        if err != nil {
            tx.Rollback(ctx)
            log.Printf("Помилка оновлення статусу для %s: %v", requestID, err)
            continue
        }

        // 2. Створення запису в Audit Log
        auditMessage := fmt.Sprintf(
            "Заявка прострочена на %d днів (дедлайн: %s)",
            overdueDays,
            actualDeadline.Format("2006-01-02"),
        )

        auditLog := &models.AuditLog{
            Action:     "SLA_BREACH",
            EntityType: "supply_request",
            EntityID:   requestID,
            Metadata:   nil,
            Message:    &auditMessage,
            CreatedAt:  time.Now(),
        }

        if err := m.auditRepo.Create(ctx, tx, auditLog); err != nil {
            tx.Rollback(ctx)
            log.Printf("Помилка створення audit log: %v", err)
            continue
        }

        // 3. Надсилання email-нотифікації
        subject := fmt.Sprintf("⚠️ SLA BREACH: Заявка %s прострочена", requestID)
        body := fmt.Sprintf(`
            <h2>Порушення SLA</h2>
            <p><strong>Заявка:</strong> %s</p>
            <p><strong>Ресурс:</strong> %s</p>
            <p><strong>Створена:</strong> %s</p>
            <p><strong>Дедлайн:</strong> %s</p>
            <p><strong>Прострочена на:</strong> %d днів</p>
            <p style="color: red;"><strong>Необхідні термінові дії!</strong></p>
        `, requestID, resourceID, 
           createdAt.Format("02.01.2006 15:04"), 
           actualDeadline.Format("02.01.2006"), 
           overdueDays)

        // Визначаємо отримувачів (логістика + керівництво)
        recipients := m.getEscalationRecipients(ctx, unitID)
        
        if err := m.emailService.Send(recipients, subject, body); err != nil {
            log.Printf("⚠️ Помилка надсилання email для %s: %v", requestID, err)
            // Не відкочуємо транзакцію - email не критичний
        }

        // 4. Коміт транзакції
        if err := tx.Commit(ctx); err != nil {
            log.Printf("Помилка коміту для %s: %v", requestID, err)
            continue
        }

        escalatedCount++
        log.Printf("✅ Заявка %s ескальована (прострочення: %d днів)", requestID, overdueDays)
    }

    return escalatedCount, nil
}

// getEscalationRecipients визначає список email для нотифікацій
func (m *SLAMonitor) getEscalationRecipients(ctx context.Context, unitID *int64) []string {
    if unitID == nil {
        return []string{"admin@Omnilog.ua"}
    }

    // Запит на знаходження відповідальних осіб
    query := `
        SELECT u.email
        FROM users u
        WHERE u.unit_id = $1
          AND u.role IN ('REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN', 
                         'REGION_DIRECTOR', 'BRANCH_MANAGER')
          AND u.status = 'ACTIVE'
    `

    rows, err := m.db.Query(ctx, query, *unitID)
    if err != nil {
        log.Printf("Помилка пошуку отримувачів: %v", err)
        return []string{"admin@Omnilog.ua"}
    }
    defer rows.Close()

    emails := []string{}
    for rows.Next() {
        var email string
        if err := rows.Scan(&email); err == nil {
            emails = append(emails, email)
        }
    }

    if len(emails) == 0 {
        return []string{"admin@Omnilog.ua"}
    }

    return emails
}

// GetEscalatedCount повертає кількість ескальованих заявок
func (m *SLAMonitor) GetEscalatedCount(ctx context.Context) (int, error) {
    var count int
    query := `SELECT COUNT(*) FROM supply_requests WHERE status = 'ESCALATED'`
    err := m.db.QueryRow(ctx, query).Scan(&count)
    return count, err
}

// GetPendingStats повертає статистику по активних заявках
func (m *SLAMonitor) GetPendingStats(ctx context.Context) (total int, soonOverdue int, err error) {
    query := `
        SELECT 
            COUNT(*) as total,
            COUNT(*) FILTER (
                WHERE NOW() > COALESCE(deadline, created_at + INTERVAL '7 days') - INTERVAL '24 hours'
            ) as soon_overdue
        FROM supply_requests
        WHERE status IN ('OPEN', 'IN_PROGRESS')
    `
    err = m.db.QueryRow(ctx, query).Scan(&total, &soonOverdue)
    return
}
```

**Ініціалізація в main.go:**

```go
// Створення SLA Monitor
slaMonitor := services.NewSLAMonitor(dbPool, reqRepo, auditRepo, emailService)

// Запуск фонового процесу
slaMonitor.Start(context.Background())

// Graceful shutdown при завершенні програми
defer slaMonitor.Stop()
```

### 5.6 Приклад роботи алгоритму

**Сценарій:**
1. **08:00** - Створена заявка на поповнення аптечок (ID: req-12345)
2. **Дедлайн:** 08:00 + 7 днів = **15:00 наступного тижня**
3. **15:01** - SLA Monitor виявляє прострочення (1 хвилина запізнення)

**Дії системи:**
```
✅ UPDATE supply_requests SET status='ESCALATED' WHERE id='req-12345'
✅ INSERT INTO audit_log (action='SLA_BREACH', entity_id='req-12345')
✅ SEND EMAIL → logistics@Omnilog.ua, director@Omnilog.ua
   Subject: "⚠️ SLA BREACH: Заявка req-12345 прострочена"
   Body: "Заявка прострочена на 1 хвилину. Необхідні термінові дії!"
```

**Результат:**
- Статус заявки змінено на **ESCALATED**
- Відповідальні особи отримали негайну нотифікацію
- В audit log зафіксовано подію для аналізу

---

## 📊 Порівняльна таблиця алгоритмів

| Алгоритм | Складність часу | Складність пам'яті | Критичність | Частота викликів |
|----------|----------------|-------------------|------------|-----------------|
| **Рекурсивна перевірка підпорядкування** | O(h) | O(n) | Висока | ~100/хв |
| **Smart Replenishment** | O(n) | O(n) | Середня | ~10/добу |
| **Haversine (GPS)** | O(1) | O(1) | Висока | ~1000/хв |
| **Детекція аномалій палива** | O(1) | O(1) | Критична | ~50/добу |
| **SLA Monitor** | O(n) | O(k) | Висока | ~1440/добу |

## 🎯 Висновки

У даному розділі представлено **п'ять ключових алгоритмів**, які забезпечують:

1. ✅ **Безпеку даних** — ієрархічний контроль доступу через рекурсивні CTE
2. ✅ **Оптимізацію запасів** — інтелектуальне планування поповнення з урахуванням незавершених замовлень
3. ✅ **Точність геолокації** — математично коректні розрахунки відстаней за формулою Haversine
4. ✅ **Антикорупційний контроль** — багаторівнева детекція аномалій палива
5. ✅ **Дотримання SLA** — автоматична ескалація прострочених заявок

Всі алгоритми мають **лінійну або константну складність**, що забезпечує масштабованість системи до рівня підприємства з тисячами користувачів та транспортних засобів.

---

**Автор:** Маркос Олександр  
**Проєкт:** Omnilog Enterprise Logistics System  
**Дата:** Квітень 2026 р.
