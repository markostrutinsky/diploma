# 📋 TODO LIST - Omnilog System

**Дата створення:** 29 квітня 2026 р.  
**Статус:** Active Development  
**Джерело:** Аудит кодової бази та документації SCENARIO_ORGANIZATION.md

---

## 🚨 КРИТИЧНІ БАГИ (Priority 0 - Зламана функціональність)

### ❌ Bug #1: EMPLOYEE не може створювати заявки
**Проблема:** Роль EMPLOYEE виключена з `SupplyRequestCreatorRoles`, що робить роль нефункціональною.

**Файл:** `Omnilog_backend/internal/models/auth.go`

**Поточний код (лінія ~75-82):**
```go
var SupplyRequestCreatorRoles = []UserRole{
	RoleRegionDirector,
	RoleRegionLogistician,
	RoleBranchManager,
	RoleBranchLogistician,
	RoleDeptManager,
	RoleTeamLead,
	// RoleEmployee відсутній!
}
```

**Виправлення:**
```go
var SupplyRequestCreatorRoles = []UserRole{
	RoleRegionDirector,
	RoleRegionLogistician,
	RoleBranchManager,
	RoleBranchLogistician,
	RoleDeptManager,
	RoleTeamLead,
	RoleEmployee,  // ← ДОДАТИ
}
```

**Вплив:** 
- Співробітники не можуть запитувати ресурси для роботи
- Порушує основний бізнес-процес
- Workaround: TEAM_LEAD створює заявки за них

**Estimate:** 5 хвилин  
**Тести:** Перевірити, що EMPLOYEE може POST /api/requests

---

### ❌ Bug #2: Логісти не можуть керувати інвентарем
**Проблема:** REGION_LOGISTICIAN та BRANCH_LOGISTICIAN відсутні в `InventoryManagerRoles`.

**Файл:** `Omnilog_backend/internal/models/auth.go`

**Поточний код (лінія ~93-97):**
```go
var InventoryManagerRoles = []UserRole{
	RoleRegionStorekeeper,
	RoleBranchStorekeeper,
	RoleDeptSupervisor,
	// Логісти відсутні!
}
```

**Виправлення:**
```go
var InventoryManagerRoles = []UserRole{
	RoleRegionLogistician,   // ← ДОДАТИ
	RoleBranchLogistician,   // ← ДОДАТИ
	RoleRegionStorekeeper,
	RoleBranchStorekeeper,
	RoleDeptSupervisor,
}
```

**Вплив:**
- Логісти не можуть створювати/редагувати ресурси
- Не можуть керувати категоріями
- Порушує відповідальність ролі

**Estimate:** 5 хвилин  
**Тести:** Перевірити POST/PATCH /api/inventory/resources для логістів

---

### ⚠️ Bug #3: Комірники не можуть затверджувати заявки (Спірно)
**Проблема:** REGION_STOREKEEPER та BRANCH_STOREKEEPER відсутні в `SupplyRequestApproverRoles`.

**Файл:** `Omnilog_backend/internal/models/auth.go`

**Поточний код (лінія ~85-90):**
```go
var SupplyRequestApproverRoles = []UserRole{
	RoleRegionDirector,
	RoleRegionLogistician,
	RoleBranchManager,
	RoleBranchLogistician,
	RoleDeptManager,
	// Комірники відсутні
}
```

**Варіанти рішення:**

**Опція A (рекомендовано):** Додати комірників з обмеженням
```go
var SupplyRequestApproverRoles = []UserRole{
	RoleRegionDirector,
	RoleRegionLogistician,
	RoleBranchManager,
	RoleBranchLogistician,
	RoleDeptManager,
	RoleRegionStorekeeper,   // ← Можуть self-approve невеликі заявки
	RoleBranchStorekeeper,   // ← Можуть self-approve невеликі заявки
}

// + Додати логіку обмеження в handler:
if user.Role == RoleRegionStorekeeper || user.Role == RoleBranchStorekeeper {
	if request.TotalValue > 10000 { // Поріг
		return errors.New("потрібне затвердження логіста")
	}
}
```

**Опція B:** Залишити як є, але додати роль DEPT_SUPERVISOR до approvers

**Пріоритет:** MEDIUM (залежить від бізнес-процесу)  
**Estimate:** 1 година (з логікою обмеження)

---

## 🔧 АРХІТЕКТУРНІ ОБМЕЖЕННЯ (Priority 1 - Потребує рефакторингу)

### 📌 Limitation #1: Відсутність гранулярних прав
**Проблема:** Система використовує тільки групи ролей, без окремих permissions.

**Наслідок:**
- Не можна дати часткові права (наприклад, "читати аналітику, але не експортувати")
- Якщо роль входить у групу, вона отримує ВСІ права цієї групи

**Приклад:**
```go
// Зараз: або всі права на інвентар, або жодних
middleware.RequireAnyRole(models.InventoryManagerRoles)

// Хотілося б:
middleware.RequirePermission("inventory.read")
middleware.RequirePermission("inventory.write")
middleware.RequirePermission("inventory.delete")
```

**Рішення (довгострокове):**
1. Створити таблицю `permissions`
2. Створити таблицю `role_permissions` (many-to-many)
3. Створити middleware `RequirePermission(permName string)`
4. Міграція існуючих ролей на нову систему

**Estimate:** 2-3 дні  
**Пріоритет:** LOW (працює, але негнучко)

---

### 📌 Limitation #2: Audit logs доступні занадто багатьом ролям
**Проблема:** Логи доступні через `SupplyRequestApproverRoles`, що дає доступ менеджерам.

**Файл:** `Omnilog_backend/main.go` (лінія ~283)
```go
api.GET("/audit", middleware.RequireAnyRole(models.SupplyRequestApproverRoles), auditHandler.GetLogs)
```

**Рекомендація:**
```go
// Створити окрему групу
var AuditViewerRoles = []UserRole{
	RoleSystemAdmin,
	RoleTenantAdmin,
	RoleRegionDirector,  // Тільки top-level ролі
}

api.GET("/audit", middleware.RequireAnyRole(models.AuditViewerRoles), auditHandler.GetLogs)
```

**Estimate:** 15 хвилин  
**Пріоритет:** MEDIUM (безпека)

---

### 📌 Limitation #3: Smart Distribution endpoint відсутній
**Проблема:** Задокументовано в API, але не реалізовано.

**Очікуваний endpoint:** `POST /api/requests/smart-distribute`

**Що робити:**
1. Видалити з документації (якщо не планується)
2. Або реалізувати (якщо критично)

**Estimate:** Невідомо (залежить від алгоритму)  
**Пріоритет:** LOW (nice-to-have)

---

### 📌 Limitation #4: Cascade deletion для Units
**Проблема:** При видаленні Company → Region → Branch каскад може не працювати.

**Рекомендація:**
```sql
-- Додати в міграцію
ALTER TABLE regions 
ADD CONSTRAINT fk_company 
FOREIGN KEY (company_id) 
REFERENCES companies(id) 
ON DELETE CASCADE;

ALTER TABLE branches 
ADD CONSTRAINT fk_region 
FOREIGN KEY (region_id) 
REFERENCES regions(id) 
ON DELETE CASCADE;
```

**Estimate:** 30 хвилин  
**Пріоритет:** MEDIUM (data integrity)

---

## 🚀 FEATURE ENHANCEMENTS (Priority 2 - Покращення UX)

### 🎯 EMPLOYEE Role Enhancement

#### Priority 1 - Критично необхідні

- [ ] **GET /api/inventory/my-equipment** - Перегляд призначеного обладнання
  ```go
  // Handler
  func (h *InventoryHandler) GetMyEquipment(c *gin.Context) {
      userID := c.GetString("user_id")
      resources, err := h.repo.GetAssignedToUser(ctx, userID)
      // ...
  }
  ```
  **Estimate:** 2 години  
  **Тести:** EMPLOYEE бачить тільки свої ресурси

- [ ] **POST /api/requests/:id/confirm-receipt** - Підтвердження отримання
  ```go
  func (h *RequestHandler) ConfirmReceipt(c *gin.Context) {
      // Статус: COMPLETED → RECEIVED
      // Audit log
      // Notification до логіста
  }
  ```
  **Estimate:** 3 години  
  **Тести:** Тільки призначений EMPLOYEE може підтвердити

#### Priority 2 - Важливі для UX

- [ ] **POST /api/inventory/report-damage** - Звіт про пошкодження
  ```go
  type DamageReport struct {
      ResourceID string
      Description string
      Photo string // base64
      Severity string // "minor", "major", "critical"
  }
  ```
  **Estimate:** 4 години (з завантаженням фото)

- [ ] **GET /api/notifications** - Перегляд сповіщень
  ```go
  // Повернути: заявки схвалено, обладнання призначено, нагадування про повернення
  ```
  **Estimate:** 5 годин (з реалтайм через WebSocket)

- [ ] **GET /api/requests/history** - Історія своїх заявок
  ```go
  // Фільтр: my requests + where I'm mentioned
  ```
  **Estimate:** 1 година

#### Priority 3 - Nice-to-have

- [ ] **GET /api/documents/personal** - Документи (інструкції, форми)
- [ ] **POST /api/feedback** - Зворотній зв'язок про систему
- [ ] **GET /api/training** - Навчальні матеріали
- [ ] **GET /api/schedule** - Графік роботи зі складом
- [ ] **POST /api/inventory/request-maintenance** - Заявка на техобслуговування

---

## 📱 KIOSK TERMINAL (Priority 3 - Додаткові можливості)

### ❌ НЕ РЕКОМЕНДУЄТЬСЯ: Режим приймання товару в термінал

**Причини:**
1. Комірники не входять у `InventoryManagerRoles` → не можуть створювати ресурси
2. Приймання = складний процес (категорія, постачальник, фото, документи)
3. Термінал призначений для швидкої видачі, не приймання
4. Є повноцінний інтерфейс для додавання ресурсів

### ✅ АЛЬТЕРНАТИВА 1: Окрема сторінка "Приймання товару"
- [ ] Створити `/receive-goods` сторінку
- [ ] Сканер штрих-кодів (використати Html5QrcodeScanner)
- [ ] Всі необхідні поля: категорія, постачальник, місце зберігання
- [ ] Завантаження фото
- [ ] Генерація QR-кодів для нових товарів
- [ ] Друк етикеток

**Estimate:** 2 дні  
**Доступ:** InventoryManagerRoles + розширити для комірників (див. Bug #2)

### ✅ АЛЬТЕРНАТИВА 2: "Швидке додавання" в Inventory
- [ ] Додати компонент `<QuickAddResource />` у основний інтерфейс
- [ ] Мінімум полів: назва, штрих-код, кількість, категорія
- [ ] Автозаповнення з існуючих ресурсів
- [ ] Сканер як опція

**Estimate:** 1 день

---

## 🧪 TESTING (Priority 1)

### Матриця тестування ролей

| Роль | Створення заявок | Затвердження | Керування інвентарем | Створення складів | Транспорт |
|------|------------------|--------------|----------------------|-------------------|-----------|
| SYSTEM_ADMIN | ✅ PASS | ✅ PASS | ✅ PASS | ✅ PASS | ✅ PASS |
| TENANT_ADMIN | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO |
| REGION_DIRECTOR | ⏳ TODO | ⏳ TODO | ❌ Expected | ⏳ TODO | ⏳ TODO |
| REGION_LOGISTICIAN | ⏳ TODO | ⏳ TODO | ❌ BUG #2 | ⏳ TODO | ⏳ TODO |
| REGION_STOREKEEPER | ⏳ TODO | ❌ Expected | ⏳ TODO | ❌ Expected | ❌ Expected |
| BRANCH_MANAGER | ⏳ TODO | ⏳ TODO | ❌ Expected | ⏳ TODO | ⏳ TODO |
| BRANCH_LOGISTICIAN | ⏳ TODO | ⏳ TODO | ❌ BUG #2 | ⏳ TODO | ⏳ TODO |
| BRANCH_STOREKEEPER | ⏳ TODO | ❌ Expected | ⏳ TODO | ❌ Expected | ❌ Expected |
| DEPT_MANAGER | ⏳ TODO | ⏳ TODO | ❌ Expected | ⏳ TODO | ⏳ TODO |
| DEPT_SUPERVISOR | ⏳ TODO | ❌ Expected | ⏳ TODO | ❌ Expected | ⏳ TODO |
| TEAM_LEAD | ⏳ TODO | ❌ Expected | ❌ Expected | ❌ Expected | ❌ Expected |
| EMPLOYEE | ❌ BUG #1 | ❌ Expected | ❌ Expected | ❌ Expected | ❌ Expected |
| CONTRACTOR | ❌ Expected | ❌ Expected | ❌ Expected | ❌ Expected | ❌ Expected |

**Тестові сценарії:**
```bash
# Bug #1: EMPLOYEE requests
curl -X POST /api/requests \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" \
  -d '{"resource_id": "...", "quantity": 5}'
# Очікується: 200 OK (після фіксу)

# Bug #2: LOGISTICIAN inventory
curl -X POST /api/inventory/resources \
  -H "Authorization: Bearer $LOGISTICIAN_TOKEN" \
  -d '{"name": "Test", "quantity": 10}'
# Очікується: 200 OK (після фіксу)

# Bug #3: STOREKEEPER approve (optional)
curl -X PATCH /api/requests/:id/approve \
  -H "Authorization: Bearer $STOREKEEPER_TOKEN"
# Очікується: 200 OK або 403 (business decision)
```

---

## 📊 DOCUMENTATION UPDATES

- [x] SCENARIO_ORGANIZATION.md - оновлено з реальними обмеженнями
- [x] Додано секцію "КРИТИЧНІ БАГИ"
- [x] Додано матриці доступу з ✅/❌/⚠️
- [x] Додано EMPLOYEE enhancement рекомендації
- [x] Додано Kiosk Terminal аналіз
- [ ] API_DOCUMENTATION.md - перевірити відповідність реальним endpoint'ам
- [ ] README.md - додати секцію "Known Issues"
- [ ] ARCHITECTURE.md - документувати систему RBAC

---

## 🎯 PRIORITY ROADMAP

### Week 1: Критичні баги
- [ ] Day 1: Fix Bug #1 (EMPLOYEE requests) + тести
- [ ] Day 2: Fix Bug #2 (LOGISTICIAN inventory) + тести
- [ ] Day 3: Decide on Bug #3 (STOREKEEPER approve) + implement
- [ ] Day 4-5: Full role matrix testing

### Week 2: EMPLOYEE enhancements
- [ ] Day 1-2: My Equipment endpoint + UI
- [ ] Day 3: Confirm Receipt endpoint + UI
- [ ] Day 4-5: Report Damage endpoint + photo upload

### Week 3: Security & Architecture
- [ ] Day 1: Audit logs access restriction
- [ ] Day 2: Cascade deletion for units
- [ ] Day 3-4: Consider permission system design
- [ ] Day 5: Documentation updates

### Week 4: Optional features
- [ ] "Приймання товару" сторінка (alternative to kiosk mode)
- [ ] Notifications system for EMPLOYEE
- [ ] Request history for all roles

---

## 📈 SUCCESS METRICS

**Після виправлення Bug #1-3:**
- ✅ EMPLOYEE може створювати заявки
- ✅ Логісти можуть керувати інвентарем
- ✅ Всі 14 ролей мають функціональні права

**Після EMPLOYEE enhancements:**
- ✅ Self-service для співробітників
- ✅ Зменшення навантаження на TEAM_LEAD
- ✅ Покращення прозорості процесів

**Після архітектурних покращень:**
- ✅ Безпечніші audit logs
- ✅ Надійніше видалення даних
- ✅ Документація відповідає реальності

---

## 🔗 RELATED FILES

**Backend:**
- `Omnilog_backend/internal/models/auth.go` - визначення ролей та груп
- `Omnilog_backend/main.go` - роути з middleware
- `Omnilog_backend/internal/middleware/auth.go` - перевірка прав
- `Omnilog_backend/internal/handlers/*.go` - бізнес-логіка

**Frontend:**
- `Omnilog_frontend/src/constants/roles.ts` - frontend role groups
- `Omnilog_frontend/src/pages/KioskTerminal.tsx` - термінал
- `Omnilog_frontend/src/contexts/AuthContext.tsx` - автентифікація

**Documentation:**
- `SCENARIO_ORGANIZATION.md` - головна документація ролей
- `API_DOCUMENTATION.md` - документація API
- `ARCHITECTURE.md` - архітектура системи

---

**Last Updated:** 29 квітня 2026 р.  
**Maintainer:** Development Team  
**Status:** 🔴 3 критичні баги, 🟡 10+ покращень, 🟢 документація актуальна
