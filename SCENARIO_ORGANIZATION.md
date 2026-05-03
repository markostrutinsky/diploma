# Сценарна організація застосунку

## Зміст

1. [Вступ](#вступ)
2. [Основні принципи організації](#основні-принципи-організації)
   - [Рольова сегментація](#1-рольова-сегментація)
   - [Multi-tenant архітектура](#2-multi-tenant-архітектура)
   - [Підписочна модель](#3-підписочна-модель)
3. [Детальний опис ролей та їх функцій](#детальний-опис-ролей-та-їх-функцій)
   - [SYSTEM_ADMIN](#1-system_admin--системний-адміністратор-платформи)
   - [TENANT_ADMIN / ADMIN](#2-tenant_admin--admin--адміністратор-організації)
   - [REGION_DIRECTOR](#3-region_director--директор-регіону)
   - [REGION_LOGISTICIAN](#4-region_logistician--логіст-регіону)
   - [REGION_STOREKEEPER](#5-region_storekeeper--комірник-регіону)
   - [BRANCH_MANAGER](#6-branch_manager--менеджер-філії)
   - [BRANCH_LOGISTICIAN](#7-branch_logistician--логіст-філії)
   - [BRANCH_STOREKEEPER](#8-branch_storekeeper--комірник-філії)
   - [DEPT_MANAGER](#9-dept_manager--менеджер-локального-підрозділу)
   - [DEPT_SUPERVISOR](#10-dept_supervisor--супервайзер-відділу)
   - [TEAM_LEAD](#11-team_lead--керівник-команди)
   - [EMPLOYEE](#12-employee--співробітник)
   - [CONTRACTOR](#13-contractor--підрядникволонтер)
4. [Матриці доступу](#матриця-доступу-до-розділів-системи)
   - [Доступ до розділів системи](#матриця-доступу-до-розділів-системи)
   - [Матриця дій (Permissions)](#матриця-дій-permissions-matrix)
   - [Платні функції (PRO Features)](#платні-функції-pro-features)
   - [Матриця створення ролей](#матриця-створення-ролей)
   - [Матриця затвердження заявок](#матриця-затвердження-заявок-approval-matrix)
5. [Сценарії першого входу](#сценарії-першого-входу)
6. [Сценарії адміністрування](#сценарії-адміністрування)
7. [Сценарії управління інвентарем](#сценарії-управління-інвентарем)
8. [Сценарії управління заявками](#сценарії-управління-заявками)
9. [Сценарії управління транспортом](#сценарії-управління-транспортом)
10. [Сценарії палива та детекції шахрайства](#сценарії-палива-та-детекції-шахрайства)
11. [Сценарії аналітики та звітності](#сценарії-аналітики-та-звітності)
12. [Сценарії роботи з контрагентами](#сценарії-роботи-з-контрагентами)
13. [Сценарії адміністрування платформи](#сценарії-адміністрування-платформи)
14. [Висновки](#висновки)

---

## Вступ

Система Omnilog (Omnilog) є комплексною платформою для управління логістикою організацій, що реалізує багаторівневу модель взаємодії користувачів з різними типами доступу та можливостей. Сценарна організація застосунку визначає типові шляхи взаємодії користувачів з системою залежно від їхньої ролі, контексту та бізнес-завдань.

### Ключові показники системи

- **14 ролей користувачів** з чіткою ієрархією та розподілом відповідальності
- **28 детальних сценаріїв** використання від реєстрації до аналітики
- **4 рівні організаційної ієрархії** (Компанія → Регіон → Філія → Відділ)
- **37+ дій (permissions)** з матрицею доступу
- **8 PRO-функцій** з AI та advanced аналітикою
- **Multi-tenant архітектура** з повною ізоляцією даних
- **JWT + RBAC** безпека з повним audit trail

### Цільова аудиторія документа

- **Розробники** — для розуміння бізнес-логіки та імплементації
- **Архітектори** — для проєктування системи доступу
- **Тестувальники** — для складання тест-кейсів
- **Бізнес-аналітики** — для валідації вимог
- **Керівники проєкту** — для оцінки обсягу функціональності
- **Технічні письменники** — для створення користувацької документації

## Основні принципи організації

### 1. Рольова сегментація

Система передбачає 14 ролей користувачів, об'єднаних у 4 основні категорії:

- **Платформні адміністратори** (SYSTEM_ADMIN) — крос-тенантне управління
- **Адміністратори організацій** (TENANT_ADMIN, ADMIN) — управління в межах tenant
- **Операційний персонал** (керівники, логісти, комірники) — виконання логістичних операцій
- **Зовнішні контрагенти** (CONTRACTOR) — обмежена участь у процесах постачання

#### Повний перелік ролей та рівнів ієрархії:

| Код ролі | Українська назва | Рівень ієрархії | Призначення |
|----------|------------------|-----------------|-------------|
| **SYSTEM_ADMIN** | Системний адміністратор | Платформа | Власник платформи, крос-тенантний доступ |
| **TENANT_ADMIN** | Адміністратор організації | Компанія | Повне управління в межах tenant |
| **ADMIN** | Адміністратор (застаріле) | Компанія | Синонім TENANT_ADMIN |
| **REGION_DIRECTOR** | Директор регіону | Регіональний відділ | Управління регіональним відділом |
| **REGION_LOGISTICIAN** | Логіст регіону | Регіональний відділ | Координація логістики регіону |
| **REGION_STOREKEEPER** | Комірник регіону | Регіональний відділ | Облік регіональних складів |
| **BRANCH_MANAGER** | Менеджер філії | Філія | Управління філією |
| **BRANCH_LOGISTICIAN** | Логіст філії | Філія | Координація постачань філії |
| **BRANCH_STOREKEEPER** | Комірник філії | Філія | Облік складів філії |
| **DEPT_MANAGER** | Менеджер відділу | Локальний підрозділ | Управління локальним відділом |
| **DEPT_SUPERVISOR** | Супервайзер відділу | Локальний підрозділ | Контроль операцій відділу |
| **TEAM_LEAD** | Керівник команди | Локальний підрозділ | Координація команди |
| **EMPLOYEE** | Співробітник | Будь-який | Базовий виконавець |
| **CONTRACTOR** | Підрядник/Волонтер | Поза організацією | Виконання зовнішніх заявок |

---

## Детальний опис ролей та їх функцій

### 1. SYSTEM_ADMIN — Системний адміністратор платформи

**Сфера відповідальності:** Глобальне управління платформою

**Функціональні можливості:**

#### Управління організаціями (Tenants)
- ✅ Перегляд всіх організацій у системі
- ✅ Зміна тарифних планів організацій (FREE → PRO)
- ✅ Блокування/розблокування організацій
- ✅ Доступ до dashboard всіх tenants

#### Користувачі
- ✅ Перегляд користувачів усіх організацій
- ✅ Створення користувачів будь-якої ролі в будь-якій організації
- ✅ Зміна ролей без обмежень
- ✅ Блокування користувачів крос-тенантно

#### Аудит та моніторинг
- ✅ Повний доступ до audit logs всіх організацій
- ✅ Перегляд метрик платформи
- ✅ Моніторинг підписок та оплат

#### Технічне управління
- ✅ Доступ до `/platform-admin` (платформна панель)
- ✅ Bootstrap ініціалізації
- ✅ Налаштування глобальних параметрів системи

**Обмеження:**
- ❌ Не створює звичайних користувачів через UI організацій (це роль TENANT_ADMIN)

**Технічна реалізація:**
- Middleware: `IsPlatformAdmin()` перевіряє `role === "SYSTEM_ADMIN"`
- Доступ: крос-тенантні запити без фільтрації `tenant_id`

---

### 2. TENANT_ADMIN / ADMIN — Адміністратор організації

**Сфера відповідальності:** Повне управління організацією в межах свого tenant

**Функціональні можливості:**

#### Управління організаційною структурою
- ✅ Створення/редагування/видалення підрозділів всіх рівнів
- ✅ Перегляд повної ієрархії організації
- ✅ Зміна батьківських підрозділів

#### Управління користувачами
- ✅ Створення користувачів будь-яких ролей (крім SYSTEM_ADMIN)
- ✅ Зміна ролей у межах матриці прав
- ✅ Блокування/розблокування користувачів
- ✅ Генерація invite-токенів
- ✅ Скидання паролів

#### Управління ресурсами та інвентарем
- ✅ Створення/редагування категорій ресурсів
- ✅ Створення/редагування складів
- ✅ Перегляд всіх ресурсів організації
- ✅ Списання ресурсів
- ✅ Трансфер між складами

#### Заявки
- ✅ Створення заявок на постачання
- ✅ Затвердження/відхилення будь-яких заявок
- ✅ Створення заявок для підрядників
- ✅ Комплектація та відправка

#### Транспорт
- ✅ Додавання/редагування транспортних засобів
- ✅ Реєстрація заправок
- ✅ Призначення авто на маршрути
- ✅ Доступ до GPS-трекінгу (якщо PRO)
- ✅ Перегляд графіку ТО

#### Аналітика
- ✅ Перегляд базового dashboard
- ✅ Перегляд Advanced KPI (якщо PRO)
- ✅ Доступ до Demand Forecasting (якщо PRO)
- ✅ Fuel Anti-Fraud детекція (якщо PRO)

#### Адміністративні функції
- ✅ Управління підпискою та біллінгом
- ✅ Перегляд audit logs організації
- ✅ Налаштування SLA політик
- ✅ Експорт звітів

**Обмеження:**
- ❌ Не бачить інші організації (multi-tenant ізоляція)
- ❌ Не може змінювати системні налаштування платформи
- ❌ PRO-фічі недоступні на FREE плані

**Права створення ролей:**
Може створювати: REGION_DIRECTOR, BRANCH_MANAGER, DEPT_MANAGER, TEAM_LEAD, всі LOGISTICIAN, всі STOREKEEPER, DEPT_SUPERVISOR, EMPLOYEE, CONTRACTOR

**Технічна реалізація:**
- Middleware: `IsTenantOwner()` перевіряє `role IN ("TENANT_ADMIN", "ADMIN")`
- Всі запити фільтруються за `tenant_id` з JWT claims

---

### 3. REGION_DIRECTOR — Директор регіону

**Сфера відповідальності:** Управління регіональним відділом та підпорядкованими філіями

**Функціональні можливості:**

#### Управління структурою
- ✅ Створення філій у своєму регіоні
- ✅ Редагування підрозділів регіону
- ✅ Перегляд ієрархії регіону

#### Персонал
- ✅ Створення користувачів: BRANCH_MANAGER, DEPT_MANAGER, TEAM_LEAD, логістів, комірників, EMPLOYEE, CONTRACTOR
- ✅ Зміна ролей підлеглих (згідно матриці)
- ✅ Блокування користувачів регіону

#### Ресурси та склади
- ✅ Перегляд всіх ресурсів регіону
- ✅ Створення складів на рівні філій
- ✅ Управління категоріями ресурсів

#### Заявки
- ✅ Створення заявок на постачання
- ✅ Затвердження заявок від BRANCH_MANAGER, DEPT_MANAGER
- ✅ Створення заявок для підрядників
- ✅ Перегляд всіх заявок регіону

#### Транспорт
- ✅ Додавання транспорту для регіону
- ✅ Призначення авто на маршрути
- ✅ Перегляд статистики автопарку

#### Аналітика
- ✅ Базовий dashboard регіону
- ✅ Звіти по виконанню заявок
- ✅ KPI регіону (якщо PRO)

**Обмеження:**
- ❌ Не бачить інші регіони (якщо їх кілька)
- ❌ Не може видаляти вищі ролі (TENANT_ADMIN)

**Матриця погодження:**
Затверджує заявки від: BRANCH_MANAGER, BRANCH_LOGISTICIAN, BRANCH_STOREKEEPER, DEPT_MANAGER, TEAM_LEAD

---

### 4. REGION_LOGISTICIAN — Логіст регіону

**Сфера відповідальності:** Координація логістичних операцій на рівні регіону

**Функціональні можливості:**

#### Інвентар
- ✅ Перегляд всіх ресурсів регіону (GET `/api/inventory`)
- ❌ **НЕ МОЖЕ** створювати/редагувати категорії (обмеження: `InventoryManagerRoles` не включає логістів)
- ❌ **НЕ МОЖЕ** додавати/редагувати ресурси (обмеження: `InventoryManagerRoles`)
- ❌ **НЕ МОЖЕ** списувати ресурси (обмеження: `InventoryManagerRoles`)

#### Склади
- ✅ Перегляд складів регіону (GET `/api/warehouses`)
- ✅ Створення складів у філіях (POST `/api/warehouses` - має `WarehouseManagerRoles`)
- ✅ Редагування параметрів складів (PATCH `/api/warehouses/:id`)
- ✅ Видалення складів (DELETE `/api/warehouses/:id`)

#### Заявки
- ✅ Створення заявок (POST `/api/requests` - має `SupplyRequestCreatorRoles`)
- ✅ Затвердження заявок (POST `/api/requests/:id/approve` - має `SupplyRequestApproverRoles`)
- ✅ Відхилення заявок (POST `/api/requests/:id/reject`)
- ✅ Перегляд всіх заявок регіону (GET `/api/requests`)

#### Транспорт
- ✅ Додавання транспорту (POST `/api/vehicles` - має `FuelRecordCreatorRoles`)
- ✅ Перегляд автопарку (GET `/api/vehicles`)
- ✅ Реєстрація заправок (POST `/api/vehicles/:id/fuel`)
- ✅ Призначення водіїв (PATCH `/api/vehicles/:id/driver`)
- ✅ Реєстрація технічного обслуговування (POST `/api/vehicles/:id/maintenance`)

#### Персонал
- ✅ Створення користувачів (POST `/api/admin/users` - має `UserCreatorRoles`)
- ✅ Може створювати: REGION_STOREKEEPER, BRANCH_LOGISTICIAN, BRANCH_STOREKEEPER (згідно `RoleCreationMap`)

**КРИТИЧНІ ОБМЕЖЕННЯ (баги/обмеження реалізації):**
- ❌ **НЕ МОЖЕ** керувати інвентарем напряму (додавати/списувати ресурси) - тільки комірники
- ❌ **НЕ МОЖЕ** створювати категорії ресурсів - тільки `InventoryManagerRoles`
- ⚠️ **ОБМЕЖЕНИЙ** доступ до управління ресурсами - в основному перегляд

**Реальна роль:** Координатор логістики з правами створення заявок, управління складами та транспортом, але без прямого управління інвентарем

**Матриця погодження:**
Затверджує заявки від: BRANCH_LOGISTICIAN, BRANCH_STOREKEEPER, REGION_STOREKEEPER, DEPT_MANAGER, TEAM_LEAD

---

### 5. REGION_STOREKEEPER — Комірник регіону

**Сфера відповідальності:** Облік ресурсів на регіональних складах

**Функціональні можливості:**

#### Інвентар
- ✅ Додавання нових ресурсів (POST `/api/inventory/resources` - має `InventoryManagerRoles`)
- ✅ Оновлення кількості (PATCH `/api/inventory/resources/:id`)
- ✅ Списання ресурсів (POST `/api/inventory/resources/:id/write-off`)
- ✅ Призначення ресурсів (POST `/api/inventory/resources/:id/assign`)
- ✅ Видалення ресурсів (DELETE `/api/inventory/resources/:id`)
- ✅ Створення/редагування категорій (POST/PATCH `/api/inventory/categories`)
- ✅ Перегляд ресурсів своїх складів (GET `/api/inventory`)

#### Заявки
- ✅ Створення заявок (POST `/api/requests` - має `SupplyRequestCreatorRoles`)
- ❌ **НЕ МОЖЕ** затверджувати заявки (немає в `SupplyRequestApproverRoles`)
- ✅ Перегляд заявок (GET `/api/requests`)
- ⚠️ **Комплектація** - немає окремого ендпоінту в поточній реалізації

#### Склади
- ✅ Перегляд складів (GET `/api/warehouses`)
- ❌ **НЕ МОЖЕ** створювати склади (немає в `WarehouseManagerRoles`)
- ⚠️ **Аудит складів** - функціонал POST `/api/inventory/audit` існує, але перевірка прав не чітка

#### Заявки для підрядників
- ✅ Створення заявок (POST `/api/volunteer-requests` - має `ContractorRequestCreatorRoles`)
- ✅ Прийняття/відхилення доставок (POST `/api/volunteer-requests/:id/accept|reject`)

**КРИТИЧНІ ОБМЕЖЕННЯ:**
- ❌ **НЕ МОЖЕ** затверджувати заявки (тільки створювати)
- ❌ **НЕ МОЖЕ** створювати склади (тільки використовувати існуючі)
- ✅ **МОЖЕ** повністю керувати ресурсами (додавання, редагування, списання)

**Реальна роль:** Відповідальний за облік та управління ресурсами на складах, створює заявки на поповнення, але не має прав затвердження

**Матриця погодження:**
Його заявки затверджують: REGION_DIRECTOR, REGION_LOGISTICIAN, TENANT_ADMIN

---

### 6. BRANCH_MANAGER — Менеджер філії

**Сфера відповідальності:** Управління філією та локальними підрозділами

**Функціональні можливості:**

#### Структура
- ✅ Створення локальних підрозділів (DEPT)
- ✅ Редагування філії

#### Персонал
- ✅ Створення: DEPT_MANAGER, TEAM_LEAD, BRANCH_LOGISTICIAN, BRANCH_STOREKEEPER, DEPT_SUPERVISOR, EMPLOYEE
- ✅ Блокування користувачів філії

#### Ресурси
- ✅ Перегляд ресурсів філії
- ✅ Створення складів локальних підрозділів

#### Заявки
- ✅ Створення заявок
- ✅ Затвердження заявок від DEPT_MANAGER, TEAM_LEAD, BRANCH_LOGISTICIAN, BRANCH_STOREKEEPER
- ✅ Створення заявок для підрядників

#### Аналітика
- ✅ Dashboard філії
- ✅ Звіти по підрозділах

**Матриця погодження:**
Затверджує заявки від: DEPT_MANAGER, DEPT_SUPERVISOR, TEAM_LEAD, BRANCH_LOGISTICIAN, BRANCH_STOREKEEPER

---

### 7. BRANCH_LOGISTICIAN — Логіст філії

**Сфера відповідальності:** Логістика на рівні філії

**Функціональні можливості:**

#### Інвентар
- ✅ Перегляд ресурсів філії (GET `/api/inventory`)
- ❌ **НЕ МОЖЕ** створювати категорії (обмеження: `InventoryManagerRoles`)
- ❌ **НЕ МОЖЕ** додавати/редагувати/списувати ресурси (обмеження: `InventoryManagerRoles`)

#### Склади
- ✅ Створення складів (POST `/api/warehouses` - має `WarehouseManagerRoles`)
- ✅ Редагування складів філії (PATCH `/api/warehouses/:id`)
- ✅ Видалення складів (DELETE `/api/warehouses/:id`)
- ✅ Перегляд складів (GET `/api/warehouses`)

#### Заявки
- ✅ Створення заявок (POST `/api/requests` - має `SupplyRequestCreatorRoles`)
- ✅ Затвердження заявок (POST `/api/requests/:id/approve` - має `SupplyRequestApproverRoles`)
- ✅ Відхилення заявок (POST `/api/requests/:id/reject`)
- ✅ Перегляд заявок філії (GET `/api/requests`)

#### Заявки для підрядників
- ✅ Створення заявок для підрядників (POST `/api/volunteer-requests` - має `ContractorRequestCreatorRoles`)
- ✅ Прийняття/відхилення доставок від підрядників (POST `/api/volunteer-requests/:id/accept`)

#### Транспорт
- ✅ Додавання транспорту (POST `/api/vehicles` - має `FuelRecordCreatorRoles`)
- ✅ Реєстрація заправок (POST `/api/vehicles/:id/fuel`)
- ✅ Призначення водіїв (PATCH `/api/vehicles/:id/driver`)
- ✅ Реєстрація ТО (POST `/api/vehicles/:id/maintenance`)
- ✅ Оновлення статусу транспорту (PATCH `/api/vehicles/:id/status`)

#### Персонал
- ✅ Створення користувачів (POST `/api/admin/users` - має `UserCreatorRoles`)
- ✅ Може створювати: BRANCH_STOREKEEPER (згідно `RoleCreationMap`)

**КРИТИЧНІ ОБМЕЖЕННЯ:**
- ❌ **НЕ МОЖЕ** напряму керувати інвентарем (додавання/списання ресурсів) - це роль комірників
- ❌ **НЕ МОЖЕ** створювати категорії ресурсів
- ⚠️ **Комплектація заявок** - технічно немає окремого ендпоінту, це робить комірник

**Реальна роль:** Координатор логістики філії з правами на створення/затвердження заявок, управління складами та транспортом

**Матриця погодження:**
Затверджує заявки від: BRANCH_STOREKEEPER, DEPT_SUPERVISOR, DEPT_MANAGER (якщо в межах філії)

---

### 8. BRANCH_STOREKEEPER — Комірник філії

**Сфера відповідальності:** Облік ресурсів складів філії

**Функціональні можливості:**

#### Інвентар
- ✅ Додавання нових ресурсів (POST `/api/inventory/resources` - має `InventoryManagerRoles`)
- ✅ Оновлення кількості (PATCH `/api/inventory/resources/:id`)
- ✅ Списання ресурсів (POST `/api/inventory/resources/:id/write-off`)
- ✅ Призначення ресурсів (POST `/api/inventory/resources/:id/assign`)
- ✅ Видалення ресурсів (DELETE `/api/inventory/resources/:id`)
- ✅ Створення/редагування категорій (POST/PATCH `/api/inventory/categories`)
- ✅ Перегляд ресурсів складів філії (GET `/api/inventory`)

#### Заявки
- ✅ Створення заявок (POST `/api/requests` - має `SupplyRequestCreatorRoles`)
- ❌ **НЕ МОЖЕ** затверджувати заявки (немає в `SupplyRequestApproverRoles`)
- ✅ Перегляд заявок (GET `/api/requests`)

#### Заявки для підрядників
- ✅ Створення заявок (POST `/api/volunteer-requests` - має `ContractorRequestCreatorRoles`)
- ✅ Прийняття доставок від підрядників (POST `/api/volunteer-requests/:id/accept`)

#### Склади
- ✅ Перегляд складів (GET `/api/warehouses`)
- ❌ **НЕ МОЖЕ** створювати склади (немає в `WarehouseManagerRoles`)

**КРИТИЧНІ ОБМЕЖЕННЯ:**
- ❌ **НЕ МОЖЕ** затверджувати заявки (тільки створювати)
- ❌ **НЕ МОЖЕ** створювати склади
- ✅ **МОЖЕ** повністю керувати ресурсами на своїх складах
- ⚠️ **Обмежений** видимістю складів філії (tenant scoping через middleware)

**Реальна роль:** Відповідальний за облік ресурсів на складах філії, створює заявки на поповнення, управляє наявністю

**Матриця погодження:**
Його заявки затверджують: BRANCH_MANAGER, BRANCH_LOGISTICIAN, REGION_DIRECTOR, REGION_LOGISTICIAN, TENANT_ADMIN

---

### 9. DEPT_MANAGER — Менеджер локального підрозділу

**Сфера відповідальності:** Управління локальним відділом

**Функціональні можливості:**

#### Персонал
- ✅ Створення: TEAM_LEAD, DEPT_SUPERVISOR, EMPLOYEE

#### Структура
- ✅ Редагування свого підрозділу
- ✅ Створення дочірніх підрозділів (якщо дозволено глибиною)

#### Ресурси
- ✅ Перегляд ресурсів підрозділу

#### Заявки
- ✅ Створення заявок
- ✅ Затвердження заявок від TEAM_LEAD, DEPT_SUPERVISOR, EMPLOYEE
- ✅ Створення заявок для підрядників

**Матриця погодження:**
Затверджує заявки від: TEAM_LEAD, DEPT_SUPERVISOR, EMPLOYEE

---

### 10. DEPT_SUPERVISOR — Супервайзер відділу

**Сфера відповідальності:** Операційний контроль роботи відділу

**Функціональні можливості:**

#### Інвентар
- ✅ Додавання/оновлення ресурсів (POST/PATCH `/api/inventory/resources` - має `InventoryManagerRoles`)
- ✅ Списання ресурсів (POST `/api/inventory/resources/:id/write-off`)
- ✅ Призначення ресурсів (POST `/api/inventory/resources/:id/assign`)
- ✅ Створення/редагування категорій (POST/PATCH `/api/inventory/categories`)
- ✅ Перегляд ресурсів (GET `/api/inventory`)

#### Заявки
- ✅ Створення заявок (POST `/api/requests` - має `SupplyRequestCreatorRoles`)
- ❌ **НЕ МОЖЕ** затверджувати заявки (немає в `SupplyRequestApproverRoles`)
- ✅ Перегляд заявок (GET `/api/requests`)

#### Заявки для підрядників
- ✅ Створення заявок (POST `/api/volunteer-requests` - має `ContractorRequestCreatorRoles`)
- ✅ Прийняття доставок (POST `/api/volunteer-requests/:id/accept`)

#### Транспорт
- ✅ Додавання транспорту (POST `/api/vehicles` - має `FuelRecordCreatorRoles`)
- ✅ Перегляд транспорту підрозділу (GET `/api/vehicles`)
- ✅ Реєстрація заправок (POST `/api/vehicles/:id/fuel`)
- ✅ Реєстрація ТО (POST `/api/vehicles/:id/maintenance`)
- ✅ Оновлення статусу (PATCH `/api/vehicles/:id/status`)

#### Персонал
- ✅ Створення користувачів (POST `/api/admin/users` - має `UserCreatorRoles`)
- ✅ Може створювати: EMPLOYEE (згідно `RoleCreationMap` - немає явно, але має `UserCreatorRoles`)

**КРИТИЧНІ ОБМЕЖЕННЯ:**
- ❌ **НЕ МОЖЕ** затверджувати заявки (тільки створювати)
- ❌ **НЕ МОЖЕ** створювати склади (немає в `WarehouseManagerRoles`)
- ✅ **МОЖЕ** повністю керувати інвентарем (як комірник)
- ✅ **МОЖЕ** керувати транспортом відділу

**Реальна роль:** Операційний менеджер з правами комірника та доступом до транспорту, але без прав затвердження заявок

**Матриця погодження:**
Його заявки затверджують: DEPT_MANAGER, BRANCH_MANAGER, BRANCH_LOGISTICIAN, REGION_DIRECTOR, REGION_LOGISTICIAN, TENANT_ADMIN

---

### 11. TEAM_LEAD — Керівник команди

**Сфера відповідальності:** Координація невеликої команди виконавців

**Функціональні можливості:**

#### Персонал
- ✅ Створення користувачів (POST `/api/admin/users` - має `UserCreatorRoles`)
- ✅ Може створювати: EMPLOYEE (згідно `RoleCreationMap`)
- ✅ Перегляд користувачів свого підрозділу (GET `/api/users/visible`)

#### Заявки
- ✅ Створення заявок (POST `/api/requests` - має `SupplyRequestCreatorRoles`)
- ❌ **НЕ МОЖЕ** затверджувати заявки (немає в `SupplyRequestApproverRoles`)
- ✅ Перегляд заявок (GET `/api/requests`)

#### Ресурси
- ✅ Перегляд доступних ресурсів (GET `/api/inventory`)
- ❌ **НЕ МОЖЕ** додавати/редагувати ресурси (немає в `InventoryManagerRoles`)

#### Склади
- ✅ Перегляд складів (GET `/api/warehouses`)
- ❌ **НЕ МОЖЕ** створювати склади (немає в `WarehouseManagerRoles`)

#### Підрозділи
- ❌ **НЕ МОЖЕ** створювати підрозділи (немає в `UnitManagerRoles`)

**КРИТИЧНІ ОБМЕЖЕННЯ:**
- ❌ **НЕ МОЖЕ** затверджувати заявки (навіть свої власні)
- ❌ **НЕ МОЖЕ** керувати інвентарем (тільки перегляд)
- ❌ **НЕ МОЖЕ** створювати склади або підрозділи
- ✅ **МОЖЕ** створювати співробітників (EMPLOYEE)
- ⚠️ **Обмежений доступ** до аналітики (тільки базовий dashboard)

**Реальна роль:** Координатор команди з мінімальними правами - створює заявки та співробітників, але більшість операцій вимагає затвердження вище

**Матриця погодження:**
Його заявки затверджують: DEPT_MANAGER, DEPT_SUPERVISOR, BRANCH_MANAGER, BRANCH_LOGISTICIAN, REGION_DIRECTOR, REGION_LOGISTICIAN, TENANT_ADMIN

---

### 12. EMPLOYEE — Співробітник

**Сфера відповідальності:** Базовий виконавець, рядовий співробітник організації

**Поточний стан (має баги ❌):**

#### Заявки
- ❌ **НЕ МОЖЕ** створювати заявки (немає в `SupplyRequestCreatorRoles` - **це баг!**)
- ✅ Перегляд заявок (GET `/api/requests` - але обмежений своїми через tenant scoping)

#### Ресурси
- ✅ Перегляд доступних ресурсів (GET `/api/inventory`)
- ❌ **НЕ МОЖЕ** додавати/редагувати ресурси (немає в `InventoryManagerRoles` - **це нормально**)

#### Профіль
- ✅ Перегляд свого профілю (GET `/api/auth/me`)
- ✅ Оновлення свого профілю (PATCH `/api/users/me`)

#### Базовий доступ
- ✅ Доступ до головної сторінки Dashboard (GET `/`)
- ✅ Перегляд підрозділів організації (GET `/api/units`)

---

**🔧 РЕКОМЕНДОВАНІ ДОПОВНЕННЯ ДЛЯ РОЛІ EMPLOYEE:**

### Пріоритет 1 - Критично необхідні (без них роль нефункціональна):

#### 1. ✅ Створення заявок на ресурси
```go
var SupplyRequestCreatorRoles = []UserRole{
    // ... існуючі ролі
    RoleEmployee, // ← ДОДАТИ
}
```
**Обґрунтування:** Співробітник повинен мати можливість запитувати ресурси для своєї роботи (спорядження, інструменти, матеріали). Це основний бізнес-процес.

**Приклад використання:**
- Військовий запитує термобілизну
- Водій запитує паливо для авто
- Медик запитує медикаменти

**Матриця погодження:** Заявки від EMPLOYEE затверджують TEAM_LEAD, DEPT_MANAGER, DEPT_SUPERVISOR або вище.

---

### Пріоритет 2 - Бажані (покращують функціональність):

#### 2. ✅ Перегляд свого особистого обладнання
**Новий ендпоінт:** `GET /api/inventory/my-equipment`
```go
// В InventoryHandler
func (h *InventoryHandler) GetMyEquipment(c *gin.Context) {
    userID := c.GetString("user_id")
    // Повертає ресурси, призначені на цього користувача
}
```
**Обґрунтування:** Співробітник повинен бачити список обладнання, яке йому видано (зброя, спорядження, інструменти).

---

#### 3. ✅ Підтвердження отримання ресурсів
**Новий ендпоінт:** `POST /api/requests/:id/confirm-receipt`
```go
// Співробітник підтверджує, що отримав ресурс
```
**Обґрунтування:** Підвищує прозорість - не тільки комірник позначає "видано", але й співробітник підтверджує "отримав".

---

#### 4. ✅ Повідомлення про пошкодження/втрату обладнання
**Новий ендпоінт:** `POST /api/inventory/report-damage`
```json
{
  "resource_id": 123,
  "damage_type": "DAMAGED|LOST|DESTROYED",
  "description": "Пошкоджено в бою",
  "photo_url": "..."
}
```
**Обґрунтування:** Співробітник повинен мати можливість повідомити про пошкодження без створення повноцінної заявки.

---

#### 5. ⚠️ Перегляд статусу своїх заявок
**Покращення існуючого:** `GET /api/requests?created_by=me`
**Обґрунтування:** Співробітник бачить тільки свої заявки та їх статус (PENDING, APPROVED, REJECTED, COMPLETED).

---

#### 6. ⚠️ Скасування своїх заявок (до затвердження)
**Новий ендпоінт:** `POST /api/requests/:id/cancel`
```go
// Дозволити скасувати тільки якщо status = PENDING і created_by = current_user
```
**Обґрунтування:** Співробітник помилково створив заявку - повинен мати право скасувати до затвердження.

---

### Пріоритет 3 - Додаткові (nice-to-have):

#### 7. 📱 Перегляд нотифікацій
**Новий ендпоінт:** `GET /api/notifications`
```json
[
  {
    "type": "REQUEST_APPROVED",
    "message": "Вашу заявку #123 затверджено",
    "created_at": "2026-04-29T10:30:00Z",
    "read": false
  }
]
```
**Обґрунтування:** Співробітник отримує сповіщення про статус своїх заявок.

---

#### 8. 📝 Коментарі до заявок
**Новий ендпоінт:** `POST /api/requests/:id/comments`
```json
{
  "text": "Дуже терміново, виїжджаю на позиції завтра"
}
```
**Обґрунтування:** Співробітник може додати контекст до своєї заявки або відповісти на запитання керівника.

---

#### 9. 📊 Базова статистика (свій профіль)
**Новий ендпоінт:** `GET /api/users/me/stats`
```json
{
  "requests_created": 15,
  "requests_approved": 12,
  "requests_rejected": 1,
  "equipment_assigned": 5,
  "avg_approval_time_hours": 3.5
}
```
**Обґрунтування:** Співробітник бачить свою активність (геймифікація, прозорість).

---

#### 10. 🎓 Перегляд документації/інструкцій
**Новий модуль:** Knowledge Base
```
GET /api/knowledge-base/articles
GET /api/knowledge-base/categories
```
**Обґрунтування:** Співробітники можуть читати інструкції по використанню обладнання, процедури безпеки тощо.

---

### ❌ ЩО НЕ ПОТРІБНО ДОДАВАТИ ДЛЯ EMPLOYEE:

1. ❌ Затвердження заявок - це роль керівників
2. ❌ Управління інвентарем - це роль комірників
3. ❌ Створення складів - це роль логістів
4. ❌ Управління персоналом - це роль менеджерів
5. ❌ Управління транспортом - це роль логістів/менеджерів
6. ❌ Доступ до фінансової інформації
7. ❌ Доступ до audit logs всієї організації
8. ❌ Зміна структури організації

---

### 📋 ПІДСУМОК РЕКОМЕНДАЦІЙ:

**Мінімально необхідні зміни (щоб роль працювала):**
```go
// В Omnilog_backend/internal/models/auth.go
var SupplyRequestCreatorRoles = []UserRole{
    RoleSystemAdmin, RoleTenantAdmin, RoleAdmin,
    RoleRegionDirector, RoleBranchManager, RoleDeptManager, RoleTeamLead,
    RoleRegionLogistician, RoleBranchLogistician, RoleDeptSupervisor,
    RoleEmployee, // ← КРИТИЧНО: додати цей рядок
}
```

**Додаткові ендпоінти (покращують UX):**
```go
// В main.go
api.Group("/api/inventory")
{
    inv.GET("/my-equipment", invHandler.GetMyEquipment) // Тільки для EMPLOYEE
}

api.Group("/api/requests")
{
    requests.POST("/:id/confirm-receipt", reqHandler.ConfirmReceipt) // Для EMPLOYEE
    requests.POST("/:id/cancel", reqHandler.CancelOwn) // Для EMPLOYEE (тільки свої PENDING)
}
```

**Реальна роль після виправлень:** Базовий виконавець, який може запитувати ресурси, бачити своє обладнання та статус заявок.

**Матриця погодження:**
Заявки від EMPLOYEE затверджують: TEAM_LEAD, DEPT_MANAGER, DEPT_SUPERVISOR, BRANCH_MANAGER, BRANCH_LOGISTICIAN, REGION_DIRECTOR, REGION_LOGISTICIAN, TENANT_ADMIN

---

### 13. CONTRACTOR — Підрядник/Волонтер

**Сфера відповідальності:** Виконання зовнішніх заявок на постачання

**Функціональні можливості:**

#### Заявки підрядників
- ✅ Перегляд відкритих заявок організацій
- ✅ Взяття заявки в роботу
- ✅ Позначення заявки як "Доставлено"
- ✅ Комунікація з організацією через коментарі

**Обмеження:**
- ❌ Не має доступу до внутрішніх даних організацій
- ❌ Не бачить інвентар, склади, транспорт
- ❌ Не може створювати внутрішні заявки
- ❌ Не бачить користувачів організацій
- ❌ Не має доступу до аналітики

**Реєстрація:**
- Публічна (без схвалення адміністратора)
- Не прив'язаний до tenant

**Технічна особливість:**
- `tenant_id = NULL` в таблиці users
- Окремий маршрут реєстрації: `POST /api/auth/register`

---

## Матриця доступу до розділів системи

**ВАЖЛИВО:** Ця матриця базується на реальних обмеженнях middleware в коді (`RequireAnyRole`), а не на теоретичних можливостях.

| Розділ (Route) | Дозволені ролі | Реальне обмеження |
|----------------|----------------|-------------------|
| **Dashboard** (`/`) | Всі авторизовані | ✅ Базовий доступ для всіх |
| **Inventory** (`/inventory`) | Всі авторизовані | ✅ Перегляд для всіх, але редагування тільки для `InventoryManagerRoles` |
| **Warehouses** (`/warehouses`) | `WarehouseManagerRoles` | ⚠️ ADMIN, REGION_DIRECTOR, REGION_LOGISTICIAN, BRANCH_MANAGER, BRANCH_LOGISTICIAN, DEPT_MANAGER, TEAM_LEAD |
| **Requests** (`/requests`) | Всі авторизовані | ✅ Перегляд для всіх, створення - `SupplyRequestCreatorRoles`, затвердження - `SupplyRequestApproverRoles` |
| **Volunteer Requests** (`/volunteer-requests`) | Всі авторизовані | ✅ Перегляд для всіх, створення - `ContractorRequestCreatorRoles`, виконання - CONTRACTOR |
| **Units** (`/units`) | Всі авторизовані | ✅ Перегляд для всіх, створення/редагування - `UnitManagerRoles` |
| **Vehicles** (`/vehicles`) | `FuelRecordCreatorRoles` | ⚠️ Включає логістів, директорів, менеджерів, супервайзерів |
| **Admin Users** (`/admin/users`) | `UserCreatorRoles` | ⚠️ Виключає SYSTEM_ADMIN, EMPLOYEE, CONTRACTOR |
| **Analytics** (`/analytics`) | Всі авторизовані | ✅ Базовий dashboard для всіх, PRO функції потребують підписки |
| **KPI Dashboard** (`/analytics/kpi`) | Всі + PRO план | 🔒 Потребує `RequireSubscriptionTier("PRO")` |
| **GPS Tracking** (`/gps`) | Всі + PRO план | 🔒 Потребує `RequireSubscriptionTier("PRO")` |
| **Fuel Anomalies** (`/analytics/fuel-anomalies`) | Всі + PRO план | 🔒 Потребує `RequireSubscriptionTier("PRO")` |
| **Maintenance** (`/analytics/maintenance`) | Всі + PRO план | 🔒 Потребує `RequireSubscriptionTier("PRO")` |
| **Audit Logs** (`/admin/audit-logs`) | `UserCreatorRoles` | ⚠️ Не тільки ADMIN, але всі хто може створювати користувачів |
| **Billing** (`/billing`) | Frontend обмеження | ⚠️ Немає middleware обмеження на backend |
| **Platform Admin** (`/platform-admin`) | Тільки SYSTEM_ADMIN | ❌ Немає в поточній реалізації |

**Легенда:**
- ✅ - Працює як задумано
- ⚠️ - Працює, але з відхиленнями від документації
- ❌ - Не реалізовано
- 🔒 - Платна функція (PRO)

---

## Матриця дій (Permissions Matrix)

**ВАЖЛИВО:** Поточна реалізація використовує рольові групи (`RequireAnyRole`), а не детальні permissions. Нижче наведено реальну матрицю прав на основі коду.

### Ресурси (Inventory)

| Дія | Група ролей | Реальні ролі |
|-----|-------------|--------------|
| `resource_view` | Всі авторизовані | ✅ Всі можуть переглядати |
| `resource_create` | `InventoryManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_STOREKEEPER, BRANCH_STOREKEEPER, DEPT_SUPERVISOR |
| `resource_update` | `InventoryManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_STOREKEEPER, BRANCH_STOREKEEPER, DEPT_SUPERVISOR |
| `resource_delete` | `InventoryManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_STOREKEEPER, BRANCH_STOREKEEPER, DEPT_SUPERVISOR |
| `resource_writeoff` | `InventoryManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_STOREKEEPER, BRANCH_STOREKEEPER, DEPT_SUPERVISOR |
| `resource_assign` | `InventoryManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_STOREKEEPER, BRANCH_STOREKEEPER, DEPT_SUPERVISOR |
| `category_create` | `InventoryManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_STOREKEEPER, BRANCH_STOREKEEPER, DEPT_SUPERVISOR |
| `category_update` | `InventoryManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_STOREKEEPER, BRANCH_STOREKEEPER, DEPT_SUPERVISOR |

**⚠️ ПОМИЛКА:** Логісти (REGION_LOGISTICIAN, BRANCH_LOGISTICIAN) НЕ включені в `InventoryManagerRoles`, хоча за логікою повинні мати права на управління інвентарем.

### Заявки на постачання

| Дія | Група ролей | Реальні ролі |
|-----|-------------|--------------|
| `request_view` | Всі авторизовані | ✅ Всі в межах свого tenant |
| `request_create` | `SupplyRequestCreatorRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, BRANCH_MANAGER, DEPT_MANAGER, TEAM_LEAD, REGION_LOGISTICIAN, BRANCH_LOGISTICIAN, DEPT_SUPERVISOR |
| `request_approve` | `SupplyRequestApproverRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, BRANCH_MANAGER, DEPT_MANAGER, REGION_LOGISTICIAN, BRANCH_LOGISTICIAN |
| `request_reject` | `SupplyRequestApproverRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, BRANCH_MANAGER, DEPT_MANAGER, REGION_LOGISTICIAN, BRANCH_LOGISTICIAN |

**❌ ОБМЕЖЕННЯ:** EMPLOYEE не може створювати заявки (не включений в `SupplyRequestCreatorRoles`).

**❌ ОБМЕЖЕННЯ:** Комірники (REGION_STOREKEEPER, BRANCH_STOREKEEPER) не можуть затверджувати заявки (не включені в `SupplyRequestApproverRoles`).

### Склади

| Дія | Група ролей | Реальні ролі |
|-----|-------------|--------------|
| `warehouse_view` | Всі авторизовані | ✅ Всі можуть переглядати |
| `warehouse_create` | `WarehouseManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, REGION_LOGISTICIAN, BRANCH_MANAGER, BRANCH_LOGISTICIAN, DEPT_MANAGER, TEAM_LEAD |
| `warehouse_update` | `WarehouseManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, REGION_LOGISTICIAN, BRANCH_MANAGER, BRANCH_LOGISTICIAN, DEPT_MANAGER, TEAM_LEAD |
| `warehouse_delete` | `WarehouseManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, REGION_LOGISTICIAN, BRANCH_MANAGER, BRANCH_LOGISTICIAN, DEPT_MANAGER, TEAM_LEAD |

**✅ ПРАВИЛЬНО:** Логісти включені, комірники виключені (вони працюють зі складами, але не створюють їх).

### Автопарк

| Дія | Група ролей | Реальні ролі |
|-----|-------------|--------------|
| `vehicle_view` | `FuelRecordCreatorRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, REGION_LOGISTICIAN, BRANCH_MANAGER, BRANCH_LOGISTICIAN, DEPT_MANAGER, DEPT_SUPERVISOR |
| `vehicle_create` | `FuelRecordCreatorRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, REGION_LOGISTICIAN, BRANCH_MANAGER, BRANCH_LOGISTICIAN, DEPT_MANAGER, DEPT_SUPERVISOR |
| `vehicle_update` | `FuelRecordCreatorRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, REGION_LOGISTICIAN, BRANCH_MANAGER, BRANCH_LOGISTICIAN, DEPT_MANAGER, DEPT_SUPERVISOR |
| `vehicle_fuel_log` | `FuelRecordCreatorRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, REGION_LOGISTICIAN, BRANCH_MANAGER, BRANCH_LOGISTICIAN, DEPT_MANAGER, DEPT_SUPERVISOR |
| `vehicle_maintenance` | `FuelRecordCreatorRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, REGION_LOGISTICIAN, BRANCH_MANAGER, BRANCH_LOGISTICIAN, DEPT_MANAGER, DEPT_SUPERVISOR |

### Підрозділи

| Дія | Група ролей | Реальні ролі |
|-----|-------------|--------------|
| `unit_view` | Всі авторизовані | ✅ Всі можуть переглядати структуру |
| `unit_create` | `UnitManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, BRANCH_MANAGER, DEPT_MANAGER, REGION_LOGISTICIAN |
| `unit_update` | `UnitManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, BRANCH_MANAGER, DEPT_MANAGER, REGION_LOGISTICIAN |
| `unit_delete` | `UnitManagerRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, BRANCH_MANAGER, DEPT_MANAGER, REGION_LOGISTICIAN |

### Користувачі

| Дія | Група ролей | Реальні ролі |
|-----|-------------|--------------|
| `user_view` | Всі авторизовані | ✅ GET `/api/users/visible` |
| `user_create` | `UserCreatorRoles` | TENANT_ADMIN, ADMIN, REGION_DIRECTOR, BRANCH_MANAGER, DEPT_MANAGER, TEAM_LEAD, REGION_LOGISTICIAN, BRANCH_LOGISTICIAN, REGION_STOREKEEPER, BRANCH_STOREKEEPER, DEPT_SUPERVISOR |
| `user_update_role` | Backend перевірка | ⚠️ Використовує `RoleCreationMap` в handler |
| `user_block` | Backend перевірка | ⚠️ Тільки TENANT_ADMIN і вище |

**⚠️ ЗВЕРНІТЬ УВАГУ:** SYSTEM_ADMIN виключений з `UserCreatorRoles` - він працює крос-тенантно через інші механізми.

### Заявки підрядникам

| Дія | Група ролей | Реальні ролі |
|-----|-------------|--------------|
| `contractor_request_view` | Всі авторизовані | ✅ Всі можуть переглядати публічні заявки |
| `contractor_request_create` | `ContractorRequestCreatorRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, BRANCH_MANAGER, DEPT_MANAGER, REGION_LOGISTICIAN, BRANCH_LOGISTICIAN, REGION_STOREKEEPER, BRANCH_STOREKEEPER, DEPT_SUPERVISOR |
| `contractor_request_take` | Тільки CONTRACTOR | ❌ Тільки роль CONTRACTOR |
| `contractor_request_deliver` | Тільки CONTRACTOR | ❌ Тільки роль CONTRACTOR |
| `contractor_request_accept` | `ContractorRequestCreatorRoles` | SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR, BRANCH_MANAGER, DEPT_MANAGER, REGION_LOGISTICIAN, BRANCH_LOGISTICIAN, REGION_STOREKEEPER, BRANCH_STOREKEEPER, DEPT_SUPERVISOR |

### Адміністрування

| Дія | Група ролей | Реальні ролі |
|-----|-------------|--------------|
| `audit_view` | `UserCreatorRoles` | ⚠️ Не тільки ADMIN - всі хто може створювати користувачів |
| `analytics_basic` | Всі авторизовані | ✅ Базовий dashboard |
| `analytics_pro` | PRO підписка | 🔒 `RequireSubscriptionTier("PRO")` |
| `gps_tracking` | PRO підписка | 🔒 `RequireSubscriptionTier("PRO")` |
| `excel_import` | `InventoryManagerRoles` + PRO | 🔒 Комірники з PRO підпискою |
| `kiosk_terminal` | Frontend `ROLE_GROUPS.kiosk` | ✅ **БЕЗКОШТОВНО** - ADMIN, REGION_STOREKEEPER, BRANCH_STOREKEEPER, DEPT_SUPERVISOR |

**⚠️ ЗВЕРНІТЬ УВАГУ:** Kiosk Terminal - це **безкоштовна** функція, доступна на всіх тарифах. Це інструмент для швидкої видачі ресурсів через сканування штрих-кодів/QR-кодів.

---

## Платні функції (PRO Features)

**ВАЖЛИВО:** Всі PRO функції перевіряються через `RequireSubscriptionTier("PRO", dbPool)` middleware. SYSTEM_ADMIN має доступ незалежно від тарифу.

| Функція | Endpoint | Мін. тариф | Доступ для ролей | Статус реалізації |
|---------|----------|------------|------------------|-------------------|
| **Smart Replenish** | POST `/api/analytics/auto-replenish` | PRO | Всі авторизовані | ✅ Реалізовано |
| **Advanced KPI** | GET `/api/analytics/kpi` | PRO | Всі авторизовані | ✅ Реалізовано |
| **Demand Forecast** | GET `/api/analytics/forecast` | PRO | Всі авторизовані | ✅ Реалізовано |
| **Predictive Maintenance** | GET `/api/analytics/maintenance` | PRO | Всі авторизовані | ✅ Реалізовано |
| **Fuel Anomaly Detection** | GET `/api/analytics/fuel-anomalies` | PRO | Всі авторизовані | ✅ Реалізовано |
| **GPS Tracking** | POST/GET `/api/gps/*` | PRO | Всі авторизовані | ✅ Реалізовано |
| **GPS Fleet Map** | GET `/api/gps/fleet-map` | PRO | Всі авторизовані | ✅ Реалізовано |
| **Geofencing** | POST/GET `/api/gps/geofences` | PRO | Всі авторизовані | ✅ Реалізовано |
| **Excel Import** | POST `/api/inventory/resources/import` | PRO | `InventoryManagerRoles` | ✅ Реалізовано |

**Механізм перевірки:**
```go
middleware.RequireSubscriptionTier("PRO", dbPool)
```

**Логіка:**
1. Витягує `tenant_id` з JWT claims
2. Запитує `tenants.subscription_tier` з бази
3. Порівнює з вимогою (`PRO` > `BASIC`)
4. SYSTEM_ADMIN завжди пропускається
5. При недостатньому тарифі повертає `402 Payment Required` з деталями

**Відповідь при блокуванні:**
```json
{
  "error": "Функція доступна тільки на тарифі: PRO",
  "current_tier": "BASIC",
  "required_tier": "PRO",
  "upgrade_url": "/billing",
  "message": "Оновіть підписку для доступу до цієї функції"
}
```

**Аудит спроб:**
Всі спроби доступу до PRO функцій без підписки логуються в `audit_logs` з типом `UNAUTHORIZED_PREMIUM_ACCESS`.

---

## КРИТИЧНІ БАГИ ТА ОБМЕЖЕННЯ РЕАЛІЗАЦІЇ

### 🐛 БАГ 1: EMPLOYEE не може створювати заявки

**Проблема:**
```go
var SupplyRequestCreatorRoles = []UserRole{
    RoleSystemAdmin, RoleTenantAdmin, RoleAdmin,
    RoleRegionDirector, RoleBranchManager, RoleDeptManager, RoleTeamLead,
    RoleRegionLogistician, RoleBranchLogistician, RoleDeptSupervisor,
}
// ❌ EMPLOYEE відсутній!
```

**Наслідок:** Базові співробітники не можуть запитувати ресурси для роботи - це руйнує основний бізнес-процес.

**Очікувана поведінка:** EMPLOYEE повинен мати право створювати заявки, які затверджує TEAM_LEAD або вище.

**Виправлення:** Додати `RoleEmployee` до `SupplyRequestCreatorRoles`.

---

### 🐛 БАГ 2: Логісти не можуть керувати інвентарем

**Проблема:**
```go
var InventoryManagerRoles = []UserRole{
    RoleSystemAdmin, RoleTenantAdmin, RoleAdmin,
    RoleRegionStorekeeper, RoleBranchStorekeeper, RoleDeptSupervisor,
}
// ❌ REGION_LOGISTICIAN і BRANCH_LOGISTICIAN відсутні!
```

**Наслідок:** Логісти можуть тільки переглядати інвентар, але не можуть:
- Додавати нові ресурси
- Редагувати кількість
- Списувати ресурси
- Створювати категорії

**Очікувана поведінка:** Логісти повинні мати повні права на управління інвентарем (це їхня основна роль).

**Виправлення:** Додати `RoleRegionLogistician, RoleBranchLogistician` до `InventoryManagerRoles`.

---

### 🐛 БАГ 3: Комірники не можуть затверджувати заявки

**Проблема:**
```go
var SupplyRequestApproverRoles = []UserRole{
    RoleSystemAdmin, RoleTenantAdmin, RoleAdmin,
    RoleRegionDirector, RoleBranchManager, RoleDeptManager,
    RoleRegionLogistician, RoleBranchLogistician,
}
// ❌ REGION_STOREKEEPER і BRANCH_STOREKEEPER відсутні!
```

**Наслідок:** Комірники можуть тільки створювати заявки, але не можуть їх затверджувати. Це створює зайвий бюрократичний бар'єр.

**Дискусійно:** Це може бути навмисне рішення (розподіл обов'язків), але в реальності комірники часто мають право самостійно розподіляти ресурси.

**Рекомендація:** Додати опціонально `RoleRegionStorekeeper, RoleBranchStorekeeper` до `SupplyRequestApproverRoles` або створити окремий permission `can_self_approve_small_requests`.

---

### ⚠️ ОБМЕЖЕННЯ 4: Немає детальних permissions

**Проблема:** Поточна реалізація використовує тільки рольові групи (`RequireAnyRole`), без детальних permissions типу:
- `inventory:create`
- `inventory:read`
- `inventory:update`
- `inventory:delete`

**Наслідок:** Неможливо дати комусь право "читати але не редагувати" в межах одного модуля.

**Приклад:** TEAM_LEAD може переглядати інвентар, але якби хотілося дати йому право додавати ресурси (але не видаляти) - це неможливо без зміни всієї групи.

**Рекомендація:** Впровадити RBAC з детальними permissions (можна використати бібліотеку `casbin`).

---

### ⚠️ ОБМЕЖЕННЯ 5: Аудит доступний не тільки ADMIN

**Проблема:**
```go
admin.GET("/audit-logs", auditHandler.GetLogs)
// Використовує UserCreatorRoles, а не окремий AdminRoles
```

**Наслідок:** Audit logs можуть переглядати не тільки TENANT_ADMIN, але й REGION_DIRECTOR, BRANCH_MANAGER, DEPT_MANAGER, TEAM_LEAD та навіть комірники/логісти.

**Дискусійно:** Це може бути фіча (прозорість), але зазвичай audit logs - це привілейована інформація тільки для адміністраторів.

**Рекомендація:** Створити окремий middleware `RequireAdminRoles` для `/admin/audit-logs`.

---

### ⚠️ ОБМЕЖЕННЯ 6: Немає каскадного видалення

**Проблема:** При видаленні підрозділу (unit) немає автоматичного:
- Перенесення користувачів в інший підрозділ
- Перенесення складів
- Перенесення транспорту

**Наслідок:** Можливі orphaned records (записи без батьківського підрозділу).

**Рекомендація:** Додати foreign key constraints з `ON DELETE CASCADE` або `ON DELETE SET NULL`, або перевірку в handler перед видаленням.

---

### ⚠️ ОБМЕЖЕННЯ 7: Smart Distribution не реалізовано повністю

**Проблема:** В документації описано `POST /api/requests/smart-distribute`, але в коді `main.go` немає такого ендпоінту.

**Наслідок:** AI-розподіл заявок по транспорту недоступний.

**Рекомендація:** Або додати ендпоінт, або видалити з документації.

---

### ✅ ЩО ПРАЦЮЄ ПРАВИЛЬНО

1. **Multi-tenant ізоляція:** Middleware автоматично фільтрує дані за `tenant_id` ✅
2. **JWT токени:** Refresh tokens працюють, автоматично оновлюються ✅
3. **Підписки:** PRO функції коректно блокуються для FREE тарифу ✅
4. **GPS трекінг:** Всі ендпоінти реалізовані (`/api/gps/*`) ✅
5. **Fuel anomalies:** Детекція працює (`/api/analytics/fuel-anomalies`) ✅
6. **Аналітика:** KPI, forecast, maintenance - всі ендпоінти є ✅
7. **Заявки підрядникам:** Повний цикл (create → take → deliver → accept) ✅

---

## Матриця створення ролей

Показує, які ролі може створювати кожна роль:

| Роль-створювач | Може створити |
|----------------|---------------|
| **SYSTEM_ADMIN** | Будь-яку роль (включно з TENANT_ADMIN) |
| **TENANT_ADMIN/ADMIN** | Всі, крім SYSTEM_ADMIN |
| **REGION_DIRECTOR** | BRANCH_MANAGER, DEPT_MANAGER, TEAM_LEAD, логісти, комірники, EMPLOYEE, CONTRACTOR |
| **BRANCH_MANAGER** | DEPT_MANAGER, TEAM_LEAD, BRANCH_LOGISTICIAN, BRANCH_STOREKEEPER, DEPT_SUPERVISOR, EMPLOYEE |
| **DEPT_MANAGER** | TEAM_LEAD, DEPT_SUPERVISOR, EMPLOYEE |
| **TEAM_LEAD** | EMPLOYEE |
| **REGION_LOGISTICIAN** | REGION_STOREKEEPER, BRANCH_LOGISTICIAN, BRANCH_STOREKEEPER |
| **BRANCH_LOGISTICIAN** | BRANCH_STOREKEEPER |

**Інші ролі** не можуть створювати користувачів.

---

## Матриця затвердження заявок (Approval Matrix)

Показує, хто може затверджувати заявки від кого:

| Автор заявки | Хто затверджує |
|--------------|----------------|
| **EMPLOYEE** | TEAM_LEAD, DEPT_MANAGER, DEPT_SUPERVISOR, BRANCH_MANAGER, BRANCH_LOGISTICIAN, REGION_DIRECTOR, REGION_LOGISTICIAN, ADMIN |
| **TEAM_LEAD** | DEPT_MANAGER, DEPT_SUPERVISOR, BRANCH_MANAGER, BRANCH_LOGISTICIAN, REGION_DIRECTOR, REGION_LOGISTICIAN, ADMIN |
| **DEPT_SUPERVISOR** | DEPT_MANAGER, BRANCH_MANAGER, BRANCH_LOGISTICIAN, REGION_DIRECTOR, REGION_LOGISTICIAN, ADMIN |
| **DEPT_MANAGER** | BRANCH_MANAGER, BRANCH_LOGISTICIAN, REGION_DIRECTOR, REGION_LOGISTICIAN, ADMIN |
| **BRANCH_MANAGER** | REGION_DIRECTOR, REGION_LOGISTICIAN, ADMIN |
| **BRANCH_LOGISTICIAN** | BRANCH_MANAGER, REGION_DIRECTOR, REGION_LOGISTICIAN, ADMIN |
| **BRANCH_STOREKEEPER** | BRANCH_MANAGER, BRANCH_LOGISTICIAN, REGION_DIRECTOR, REGION_LOGISTICIAN, ADMIN |
| **REGION_STOREKEEPER** | REGION_DIRECTOR, REGION_LOGISTICIAN, ADMIN |
| **REGION_DIRECTOR** | Тільки ADMIN |
| **REGION_LOGISTICIAN** | REGION_DIRECTOR, ADMIN |

**Правило:** Заявку затверджує безпосередній керівник або вище за ієрархією.

---

## Візуалізація ієрархії ролей

### Організаційна структура з розподілом ролей

```
┌─────────────────────────────────────────────────────────────────┐
│                      ПЛАТФОРМА (Platform)                        │
│                   SYSTEM_ADMIN (крос-тенант)                    │
└─────────────────────────────────────────────────────────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │    TENANT (Організація)  │
                    │   TENANT_ADMIN / ADMIN   │
                    └────────────┬────────────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        │                        │                        │
┌───────▼────────┐      ┌────────▼────────┐     ┌────────▼────────┐
│ REGIONAL OFFICE │      │ REGIONAL OFFICE │     │ REGIONAL OFFICE │
│ (Регіональний   │      │ (Регіональний   │     │ (Регіональний   │
│  відділ)        │      │  відділ)        │     │  відділ)        │
│                 │      │                 │     │                 │
│ • REGION_       │      │ • REGION_       │     │ • REGION_       │
│   DIRECTOR      │      │   DIRECTOR      │     │   DIRECTOR      │
│ • REGION_       │      │ • REGION_       │     │ • REGION_       │
│   LOGISTICIAN   │      │   LOGISTICIAN   │     │   LOGISTICIAN   │
│ • REGION_       │      │ • REGION_       │     │ • REGION_       │
│   STOREKEEPER   │      │   STOREKEEPER   │     │   STOREKEEPER   │
└────────┬────────┘      └────────┬────────┘     └────────┬────────┘
         │                        │                        │
    ┌────┴────┐              ┌────┴────┐              ┌────┴────┐
    │         │              │         │              │         │
┌───▼───┐ ┌───▼───┐      ┌───▼───┐ ┌───▼───┐      ┌───▼───┐ ┌───▼───┐
│BRANCH │ │BRANCH │      │BRANCH │ │BRANCH │      │BRANCH │ │BRANCH │
│(Філія)│ │(Філія)│      │(Філія)│ │(Філія)│      │(Філія)│ │(Філія)│
│       │ │       │      │       │ │       │      │       │ │       │
│ • BRA-│ │ • BRA-│      │ • BRA-│ │ • BRA-│      │ • BRA-│ │ • BRA-│
│   NCH_│ │   NCH_│      │   NCH_│ │   NCH_│      │   NCH_│ │   NCH_│
│   MANA│ │   MANA│      │   MANA│ │   MANA│      │   MANA│ │   MANA│
│   GER │ │   GER │      │   GER │ │   GER │      │   GER │ │   GER │
│ • BRA-│ │ • BRA-│      │ • BRA-│ │ • BRA-│      │ • BRA-│ │ • BRA-│
│   NCH_│ │   NCH_│      │   NCH_│ │   NCH_│      │   NCH_│ │   NCH_│
│   LOGI│ │   LOGI│      │   LOGI│ │   LOGI│      │   LOGI│ │   LOGI│
│   STIC│ │   STIC│      │   STIC│ │   STIC│      │   STIC│ │   STIC│
│   IAN │ │   IAN │      │   IAN │ │   IAN │      │   IAN │ │   IAN │
│ • BRA-│ │ • BRA-│      │ • BRA-│ │ • BRA-│      │ • BRA-│ │ • BRA-│
│   NCH_│ │   NCH_│      │   NCH_│ │   NCH_│      │   NCH_│ │   NCH_│
│   STOR│ │   STOR│      │   STOR│ │   STOR│      │   STOR│ │   STOR│
│   EKEE│ │   EKEE│      │   EKEE│ │   EKEE│      │   EKEE│ │   EKEE│
│   PER │ │   PER │      │   PER │ │   PER │      │   PER │ │   PER │
└───┬───┘ └───┬───┘      └───┬───┘ └───┬───┘      └───┬───┘ └───┬───┘
    │         │              │         │              │         │
┌───▼───┐ ┌───▼───┐      ┌───▼───┐ ┌───▼───┐      ┌───▼───┐ ┌───▼───┐
│ DEPT  │ │ DEPT  │      │ DEPT  │ │ DEPT  │      │ DEPT  │ │ DEPT  │
│(Локал)│ │(Локал)│      │(Локал)│ │(Локал)│      │(Локал)│ │(Локал)│
│       │ │       │      │       │ │       │      │       │ │       │
│ • DEP-│ │ • DEP-│      │ • DEP-│ │ • DEP-│      │ • DEP-│ │ • DEP-│
│   T_MA│ │   T_MA│      │   T_MA│ │   T_MA│      │   T_MA│ │   T_MA│
│   NAGE│ │   NAGE│      │   NAGE│ │   NAGE│      │   NAGE│ │   NAGE│
│   R   │ │   R   │      │   R   │ │   R   │      │   R   │ │   R   │
│ • DEP-│ │ • DEP-│      │ • DEP-│ │ • DEP-│      │ • DEP-│ │ • DEP-│
│   T_SU│ │   T_SU│      │   T_SU│ │   T_SU│      │   T_SU│ │   T_SU│
│   PERV│ │   PERV│      │   PERV│ │   PERV│      │   PERV│ │   PERV│
│   ISOR│ │   ISOR│      │   ISOR│ │   ISOR│      │   ISOR│ │   ISOR│
│ • TEAM│ │ • TEAM│      │ • TEAM│ │ • TEAM│      │ • TEAM│ │ • TEAM│
│   _LEA│ │   _LEA│      │   _LEA│ │   _LEA│      │   _LEA│ │   _LEA│
│   D   │ │   D   │      │   D   │ │   D   │      │   D   │ │   D   │
│ • EMPL│ │ • EMPL│      │ • EMPL│ │ • EMPL│      │ • EMPL│ │ • EMPL│
│   OYEE│ │   OYEE│      │   OYEE│ │   OYEE│      │   OYEE│ │   OYEE│
└───────┘ └───────┘      └───────┘ └───────┘      └───────┘ └───────┘

                    ┌─────────────────────────┐
                    │  ЗОВНІШНІ КОНТРАГЕНТИ   │
                    │      CONTRACTOR         │
                    │  (поза організацією)    │
                    └─────────────────────────┘
```

### Функціональні групи ролей

```
┌────────────────────────────────────────────────────────────────┐
│                   СТРАТЕГІЧНЕ УПРАВЛІННЯ                       │
│  SYSTEM_ADMIN, TENANT_ADMIN, ADMIN, REGION_DIRECTOR           │
│  • Управління організацією                                     │
│  • Стратегічні рішення                                         │
│  • Фінанси та підписки                                         │
│  • Повна аналітика                                             │
└────────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
┌───────▼───────┐  ┌────────▼────────┐  ┌──────▼──────┐
│   ЛОГІСТИКА   │  │   СКЛАДСЬКИЙ    │  │ ОПЕРАТИВНЕ  │
│               │  │     ОБЛІК       │  │  УПРАВЛІННЯ │
│ • REGION_LOGI-│  │                 │  │             │
│   STICIAN     │  │ • REGION_STORE- │  │ • BRANCH_   │
│ • BRANCH_LOGI-│  │   KEEPER        │  │   MANAGER   │
│   STICIAN     │  │ • BRANCH_STORE- │  │ • DEPT_     │
│               │  │   KEEPER        │  │   MANAGER   │
│ Функції:      │  │ • DEPT_SUPER-   │  │ • TEAM_LEAD │
│ • Планування  │  │   VISOR         │  │             │
│   поставок    │  │                 │  │ Функції:    │
│ • Розподіл    │  │ Функції:        │  │ • Управління│
│   транспорту  │  │ • Облік ресурсів│  │   персоналом│
│ • Оптимізація │  │ • Інвентаризація│  │ • Координація│
│   маршрутів   │  │ • Видача ресурсів│ │   роботи    │
│ • Контроль    │  │ • Контроль залиш│  │ • Створення │
│   витрат      │  │   ків           │  │   заявок    │
└───────────────┘  └─────────────────┘  └─────────────┘
        │                   │                   │
        └───────────────────┼───────────────────┘
                            │
                ┌───────────▼───────────┐
                │  БАЗОВІ ВИКОНАВЦІ     │
                │                       │
                │ • EMPLOYEE            │
                │                       │
                │ Функції:              │
                │ • Створення заявок    │
                │ • Отримання ресурсів  │
                │ • Базові операції     │
                └───────────────────────┘
```

### Потік погодження заявки (Approval Flow)

```
EMPLOYEE створює заявку
         │
         ▼
    ┌────────────┐
    │  PENDING   │ ◄─── Очікує затвердження
    └─────┬──────┘
          │
          ├──► TEAM_LEAD (може затвердити)
          │          │
          │          ▼
          ├──► DEPT_MANAGER (може затвердити)
          │          │
          │          ▼
          ├──► BRANCH_MANAGER (може затвердити)
          │          │
          │          ▼
          └──► REGION_DIRECTOR (може затвердити)
                     │
                     ▼
              ┌─────────────┐
              │  APPROVED   │
              └──────┬──────┘
                     │
                     ▼
         Комірник комплектує (DISPATCH)
                     │
                     ▼
              ┌─────────────┐
              │ IN_TRANSIT  │
              └──────┬──────┘
                     │
                     ▼
           Одержувач приймає
                     │
                     ▼
              ┌─────────────┐
              │  COMPLETED  │
              └─────────────┘
```



### 2. Multi-tenant архітектура

Кожна організація працює у власному ізольованому просторі даних з чотирирівневою ієрархією підрозділів:

```
Компанія (Tenant)
  └── Регіональний відділ
      └── Філія
          └── Локальний підрозділ
```

### 3. Підписочна модель

Доступ до функціональності контролюється тарифним планом:
- **FREE** — базові можливості обліку
- **PRO** — розширена аналітика, GPS-трекінг, AI-прогнози

---

## Сценарії першого входу

### Сценарій 1: Ініціалізація платформи (Bootstrap)

**Актор:** Власник платформи  
**Передумови:** Система розгорнута вперше, БД порожня  
**Тригер:** Перший запуск системи

**Потік подій:**

1. Користувач відкриває веб-інтерфейс за адресою `/bootstrap`
2. Система перевіряє, чи існує хоча б один SYSTEM_ADMIN
3. Якщо ні — відображається форма створення першого адміністратора
4. Користувач вводить:
   - Email
   - Пароль (мін. 8 символів)
   - Повне ім'я
5. Система створює запис SYSTEM_ADMIN у таблиці `users`
6. Логується подія аудиту "Ініціалізація головного адміністратора"
7. Користувач перенаправляється на `/login`

**Постумова:** Платформа готова до реєстрації організацій

**Технічна реалізація:**
- Handler: `AuthHandler.BootstrapAdmin()`
- Endpoint: `POST /api/auth/bootstrap`
- Валідація: перевірка відсутності існуючих адміністраторів

---

### Сценарій 2: Реєстрація нової організації (Tenant Signup)

**Актор:** Представник організації  
**Передумови:** Платформа ініціалізована  
**Тригер:** Бажання організації почати використовувати систему

**Потік подій:**

1. Користувач відкриває `/signup`
2. Заповнює форму реєстрації:
   - Назва організації
   - Email адміністратора
   - Пароль
3. Система валідує унікальність email
4. Створюється:
   - Запис у таблиці `tenants` (організація)
   - Головний підрозділ з `hierarchy_level = 'COMPANY'`
   - Користувач з роллю `TENANT_ADMIN`
5. Генерується пара JWT токенів (access + refresh)
6. Логується створення tenant в audit log
7. Користувач автоматично авторизується

**Постумова:** Організація створена з безкоштовним FREE планом

**Альтернативні потоки:**
- Email зайнятий → повернення 409 Conflict
- Невалідні дані → повернення 400 Bad Request з описом помилок

**Технічна реалізація:**
- Handler: `AuthHandler.SignupTenant()`
- Endpoint: `POST /api/auth/tenants/signup`
- Service: `AuthService.CreateTenant()`

---

### Сценарій 3: Авторизація користувача

**Актор:** Зареєстрований користувач  
**Передумови:** Користувач має активний акаунт  
**Тригер:** Відкриття сторінки `/login`

**Потік подій:**

1. Користувач вводить email і пароль
2. Система знаходить користувача за email
3. Перевіряється:
   - Статус блокування (`blocked = false`)
   - Відповідність паролю (bcrypt hash)
4. Генерується пара JWT токенів
5. Refresh token зберігається в БД
6. Повертається відповідь з токенами та даними користувача:
   ```json
   {
     "access_token": "eyJhbGc...",
     "refresh_token": "eyJhbGc...",
     "user": {
       "id": "uuid",
       "email": "user@example.com",
       "role": "BRANCH_MANAGER",
       "tenant_id": "tenant-uuid",
       "unit_id": 42
     }
   }
   ```
7. Логується успішний вхід в audit log
8. Frontend зберігає токени в localStorage
9. Перенаправлення на головну сторінку `/`

**Постумова:** Користувач авторизований, має доступ до захищених ресурсів

**Альтернативні потоки:**
- Невірний пароль → 401 Unauthorized
- Користувач заблокований → 403 Forbidden
- Неіснуючий email → 401 Unauthorized

**Технічна реалізація:**
- Handler: `AuthHandler.Login()`
- Endpoint: `POST /api/auth/login`
- Middleware: токени додаються до HTTP-заголовків `Authorization: Bearer <token>`

---

## Сценарії адміністрування

### Сценарій 4: Створення користувача адміністратором

**Актор:** TENANT_ADMIN або вищі ролі  
**Передумови:** Адміністратор авторизований  
**Тригер:** Натискання кнопки "Додати користувача"

**Потік подій:**

1. Адміністратор відкриває сторінку `/admin/users`
2. Натискає "Створити користувача"
3. Заповнює форму:
   - Email (обов'язково)
   - Повне ім'я
   - Роль (випадаючий список доступних ролей)
   - Підрозділ (якщо роль не CONTRACTOR)
4. Система валідує:
   - Унікальність email в межах tenant
   - Права створювача на призначення вибраної ролі
   - Відповідність підрозділу і ролі
5. Генерується унікальний invite token
6. Створюється користувач з `blocked = true`
7. Відправляється email з посиланням `/setup-password?token=<token>`
8. Логується створення користувача в audit log
9. Користувач з'являється у списку зі статусом "Очікує активації"

**Постумова:** Користувач зареєстрований, має 24 години на встановлення пароля

**Альтернативні потоки:**
- Email зайнятий → помилка "Email вже використовується"
- Недостатньо прав → 403 Forbidden
- Невалідна роль для підрозділу → помилка валідації

**Технічна реалізація:**
- Handler: `AuthHandler.RegisterUser()`
- Endpoint: `POST /api/auth/users`
- Матриця прав: `models.AuthRoleMatrix`

---

### Сценарій 5: Активація облікового запису за токеном

**Актор:** Новий користувач  
**Передумови:** Отримано email з invite посиланням  
**Тригер:** Перехід за посиланням з email

**Потік подій:**

1. Користувач відкриває `/setup-password?token=<token>`
2. Система валідує токен:
   - Існування в таблиці `invite_tokens`
   - Час дії (не більше 24 годин)
   - Статус `used = false`
3. Відображається форма встановлення пароля
4. Користувач вводить пароль (мін. 8 символів)
5. Система:
   - Хешує пароль (bcrypt)
   - Оновлює запис користувача
   - Встановлює `blocked = false`
   - Позначає токен як використаний
6. Логується активація в audit log
7. Відображається повідомлення "Пароль встановлено, можете увійти"
8. Перенаправлення на `/login`

**Постумова:** Користувач може авторизуватися з новим паролем

**Альтернативні потоки:**
- Токен протермінований → помилка "Посилання недійсне"
- Токен вже використаний → помилка "Посилання вже використане"

**Технічна реалізація:**
- Handler: `AuthHandler.SetupPassword()`
- Endpoint: `POST /api/auth/setup-password`
- Repository: `InviteTokenRepository`

---

### Сценарій 6: Управління структурою підрозділів

**Актор:** TENANT_ADMIN, REGION_DIRECTOR, BRANCH_MANAGER  
**Передумови:** Користувач має права управління організаційною структурою  
**Тригер:** Необхідність створити новий підрозділ

**Потік подій:**

1. Керівник відкриває сторінку `/units`
2. Вибирає батьківський підрозділ у дереві
3. Натискає "Створити дочірній підрозділ"
4. Заповнює форму:
   - Назва підрозділу
   - Тип (автоматично визначається за батьківським):
     - Компанія → Регіональний відділ
     - Регіональний відділ → Філія
     - Філія → Локальний підрозділ
5. Система валідує:
   - Максимальну глибину ієрархії (4 рівні)
   - Права на створення в цьому підрозділі
6. Створюється запис у `units` з `parent_id` та `hierarchy_level`
7. Логується в audit log
8. Дерево підрозділів оновлюється

**Постумова:** Новий підрозділ доступний для призначення користувачів і складів

**Візуалізація:**
```
26-та Окрема Бригада [COMPANY]
├── Центральний регіон [REGIONAL_OFFICE]
│   ├── Київська філія [BRANCH]
│   │   ├── Склад Святошино [DEPT - warehouse]
│   │   └── Склад Дарниця [DEPT - warehouse]
│   └── Одеська філія [BRANCH]
└── Східний регіон [REGIONAL_OFFICE]
```

**Технічна реалізація:**
- Handler: `UnitsHandler.Create()`
- Endpoint: `POST /api/units`
- Frontend: Рекурсивний компонент `UnitTreeNode`

---

## Сценарії управління інвентарем

### Сценарій 7: Створення категорії ресурсів

**Актор:** Комірник (STOREKEEPER) або логіст  
**Передумови:** Користувач має доступ до управління інвентарем  
**Тригер:** Необхідність класифікувати нові типи ресурсів

**Потік подій:**

1. Комірник відкриває `/inventory`
2. Переходить на вкладку "Категорії"
3. Натискає "Створити категорію"
4. Вводить назву (наприклад, "Медикаменти", "Боєприпаси", "Продовольство")
5. Система створює запис у `categories`
6. Категорія стає доступною для класифікації ресурсів

**Постумова:** Категорія доступна в усіх складах tenant

---

### Сценарій 8: Додавання ресурсу до складу

**Актор:** Комірник складу  
**Передумови:** Існує склад і категорія  
**Тригер:** Надходження ресурсу на склад

**Потік подій:**

1. Комірник обирає свій склад на сторінці `/inventory`
2. Натискає "Додати ресурс"
3. Заповнює форму:
   - Назва ресурсу (наприклад, "Джгут кровоспинний")
   - Категорія (вибір зі списку)
   - Початкова кількість
   - Одиниця виміру (шт, кг, л)
   - Мінімальний залишок (для сповіщень)
4. Система створює запис у `resources` з прив'язкою до `warehouse_id`
5. Логується створення в audit log
6. Ресурс з'являється у таблиці з поточним балансом

**Постумова:** Ресурс доступний для включення в заявки на постачання

**Технічна реалізація:**
- Handler: `InventoryHandler.Create()`
- Endpoint: `POST /api/inventory`
- Валідація: перевірка прав комірника на цей склад

---

### Сценарій 9: Моніторинг критичних залишків

**Актор:** Логіст або комірник  
**Передумови:** Ресурси мають встановлений `min_quantity`  
**Тригер:** Автоматична перевірка або відкриття `/inventory`

**Потік подій:**

1. Система щодня (або при змінах) перевіряє умову:
   ```sql
   SELECT * FROM resources 
   WHERE quantity <= min_quantity 
     AND min_quantity > 0
   ```
2. Ресурси з дефіцитом позначаються червоним індикатором
3. На dashboard з'являється сповіщення "Критичні залишки: 5 ресурсів"
4. Логіст відкриває список і бачить:
   - Джгут кровоспинний: 12 шт (мін. 50)
   - Бинт стерильний: 8 уп (мін. 20)
5. Ініціює створення заявки на постачання

**Постумова:** Відповідальні особи поінформовані про необхідність поповнення

**Технічна реалізація:**
- Repository: `ResourceRepository.GetLowStockResources()`
- UI: компонент `LowStockAlert`

---

## Сценарії управління заявками

### Сценарій 10: Створення внутрішньої заявки на постачання

**Актор:** Керівник підрозділу (DEPT_MANAGER, TEAM_LEAD)  
**Передумови:** Підрозділ потребує ресурс зі складу  
**Тригер:** Виробнича необхідність

**Потік подій:**

1. Керівник відкриває `/requests`
2. Натискає "Створити заявку"
3. Обирає:
   - Ресурс зі списку доступних у складах вище за ієрархією
   - Кількість
   - Цільовий склад (свого підрозділу або нижче)
   - Пріоритет (НИЗЬКИЙ, НОРМАЛЬНИЙ, ВИСОКИЙ, КРИТИЧНИЙ)
4. Система валідує:
   - Достатність ресурсу на складі-джерелі
   - Права на створення заявки
5. Створюється запис у `supply_requests` зі статусом `PENDING`
6. Відправляється нотифікація керівнику, що має право затверджувати (визначається матрицею `APPROVAL_MATRIX`)
7. Логується створення в audit log

**Постумова:** Заявка очікує розгляду затверджуючим

**Матриця погодження:**
```typescript
const APPROVAL_MATRIX = {
  'TEAM_LEAD': ['DEPT_MANAGER', 'BRANCH_MANAGER', 'REGION_DIRECTOR'],
  'DEPT_MANAGER': ['BRANCH_MANAGER', 'REGION_DIRECTOR'],
  'BRANCH_MANAGER': ['REGION_DIRECTOR', 'ADMIN']
}
```

**Технічна реалізація:**
- Handler: `RequestsHandler.Create()`
- Endpoint: `POST /api/requests`
- Валідація ресурсів: `ResourceRepository.GetAvailableForUnit()`

---

### Сценарій 11: Затвердження заявки керівником

**Актор:** Затверджувач (APPROVER, BRANCH_MANAGER, REGION_DIRECTOR)  
**Передумови:** Існує заявка зі статусом PENDING  
**Тригер:** Отримання нотифікації про нову заявку

**Потік подій:**

1. Затверджувач відкриває `/requests`
2. Бачить заявку зі статусом "Очікує"
3. Переглядає деталі:
   - Автор заявки
   - Запитуваний ресурс і кількість
   - Обґрунтування
4. Приймає рішення:
   - **Затвердити** → натискає "Затвердити"
   - **Відхилити** → натискає "Відхилити", вводить коментар
5. При затвердженні:
   - Статус міняється на `APPROVED`
   - `approved_by` = ID затверджувача
   - `approved_at` = поточний час
6. При відхиленні:
   - Статус міняється на `REJECTED`
   - Зберігається коментар відхилення
7. Логується дія в audit log
8. Автор заявки отримує нотифікацію

**Постумова затвердження:** Заявка готова до виконання комірником

**Технічна реалізація:**
- Handler: `RequestsHandler.Approve()`, `RequestsHandler.Reject()`
- Endpoints: 
  - `PUT /api/requests/:id/approve`
  - `PUT /api/requests/:id/reject`

---

### Сценарій 12: Виконання заявки та комплектація

**Актор:** Комірник складу  
**Передумови:** Заявка затверджена, ресурс є на складі  
**Тригер:** Заявка з'явилась у списку "Затверджені"

**Потік подій:**

1. Комірник відкриває заявку зі статусом `APPROVED`
2. Натискає "Укомплектувати і відправити"
3. Система показує модальне вікно:
   - Склад-відправник (автоматично обраний)
   - Склад-одержувач (з заявки)
   - Транспортний засіб (список доступних між цими підрозділами)
   - Пріоритет доставки
4. Комірник обирає автомобіль
5. Натискає "Підтвердити відправку"
6. Система:
   - Віднімає кількість зі складу-джерела
   - Створює запис у `vehicle_cargo` (вантаж у транспорті)
   - Змінює статус заявки на `IN_TRANSIT`
   - Відправляє нотифікацію одержувачу
7. Логується відправка в audit log

**Постумова:** Ресурс у дорозі, баланс складу оновлено

**Альтернативні потоки:**
- Недостатньо ресурсу → комірник коригує кількість або скасовує
- Немає доступних авто → можливість створити заявку на авто

**Технічна реалізація:**
- Handler: `RequestsHandler.Dispatch()`
- Endpoint: `POST /api/requests/dispatch`
- Транзакція: зміна балансу + створення cargo + оновлення статусу

---

### Сценарій 13: Прийомка ресурсу одержувачем

**Актор:** Комірник складу-одержувача  
**Передумови:** Вантаж доставлено  
**Тригер:** Фізичне прибуття транспорту

**Потік подій:**

1. Комірник бачить заявку зі статусом `IN_TRANSIT`
2. Натискає "Прийняти вантаж"
3. Система:
   - Додає кількість до балансу складу-одержувача
   - Видаляє запис з `vehicle_cargo`
   - Змінює статус заявки на `COMPLETED`
   - Оновлює `completed_at`
4. Автор заявки отримує нотифікацію про виконання
5. Логується прийомка в audit log

**Постумова:** Ресурс на балансі одержувача, заявка закрита

**Технічна реалізація:**
- Handler: `RequestsHandler.Complete()`
- Endpoint: `PUT /api/requests/:id/complete`

---

### Сценарій 14: Smart AI розподіл заявок по транспорту

**Актор:** Логіст або комірник (PRO план)  
**Передумови:** 
- Активовано PRO підписку
- Є кілька затверджених заявок
- Є транспортні засоби з різною вантажністю  
**Тригер:** Натискання "Smart Розподіл" на сторінці заявок

**Потік подій:**

1. Логіст відкриває `/requests`
2. Вибирає чекбоксами 5-10 затверджених заявок
3. Натискає кнопку "Smart Розподіл (AI)"
4. Система викликає AI-алгоритм `SmartDistribution`:
   ```go
   routes := SmartDistribution(items, vehicles)
   ```
5. Алгоритм виконує:
   - Групування заявок за маршрутами (from → to warehouse)
   - Bin Packing для розміщення вантажів у доступні авто
   - Оптимізацію за пріоритетом і вантажністю
6. Показується попередній перегляд:
   ```
   Маршрут 1: Київ → Одеса
   └── ГАЗель (1500 кг): Заявка #45 (500 кг), #47 (800 кг)
   
   Маршрут 2: Київ → Львів
   └── Камаз (10 т): Заявка #48 (7 т)
   
   Нерозподілені: Заявка #50 (немає підходящого авто)
   ```
7. Логіст підтверджує або редагує розподіл
8. Система:
   - Комплектує всі заявки одночасно
   - Призначає вантажі на транспорт
   - Змінює статуси на `IN_TRANSIT`
9. Логується масова відправка

**Постумова:** Оптимально розподілені ресурси по транспорту

**Переваги:**
- Економія часу (замість 10 окремих відправок — 1 клік)
- Ефективне використання вантажопідйомності
- Пріоритизація критичних заявок

**Технічна реалізація:**
- Handler: `RequestsHandler.SmartDistribute()`
- Endpoint: `POST /api/requests/smart-distribute`
- Algorithm: First Fit Decreasing Bin Packing

---

## Сценарії управління транспортом

### Сценарій 15: Додавання транспортного засобу

**Актор:** Менеджер автопарку (BRANCH_MANAGER, REGION_DIRECTOR)  
**Передумови:** Організація має транспорт  
**Тригер:** Постановка нового авто на баланс

**Потік подій:**

1. Менеджер відкриває `/vehicles`
2. Натискає "Додати транспорт"
3. Заповнює форму:
   - Назва/модель (наприклад, "ГАЗель #12")
   - Реєстраційний номер (АА1234ВВ)
   - Підрозділ-власник
   - Вантажопідйомність (кг)
   - Пробіг поточний (км)
   - Інтервал ТО (наприклад, кожні 10000 км)
4. Опціонально: прикріплює GPS-трекер (для PRO)
5. Система створює запис у `vehicles`
6. Транспорт стає доступним для призначення на маршрути

**Постумова:** Авто доступне для логістичних операцій

**Технічна реалізація:**
- Handler: `VehicleHandler.Create()`
- Endpoint: `POST /api/vehicles`

---

### Сценарій 16: Відстеження GPS в реальному часі (PRO)

**Актор:** Логіст або диспетчер  
**Передумови:** PRO підписка, GPS-трекери встановлені  
**Тригер:** Необхідність контролювати флот

**Потік подій:**

1. Диспетчер відкриває `/gps`
2. Система завантажує актуальні координати всіх авто:
   ```sql
   SELECT * FROM gps_tracking 
   WHERE timestamp > NOW() - INTERVAL '15 minutes'
   ORDER BY timestamp DESC
   ```
3. Відображається карта з маркерами транспорту
4. Кольори маркерів:
   - 🟢 Зелений: рухається (швидкість > 5 км/год)
   - 🟡 Жовтий: стоїть (< 5 км/год, двигун працює)
   - 🔴 Червоний: заглушено
5. Клік на маркер показує:
   - Назва авто
   - Поточна швидкість
   - Останнє оновлення
   - Вантаж (якщо є)
6. Фільтри:
   - За підрозділом
   - За статусом (у дорозі / вільні)
   - За геозоною

**Постумова:** Диспетчер бачить актуальне розташування флоту

**Технічна реалізація:**
- Handler: `GPSHandler.GetLatestLocations()`
- Endpoint: `GET /api/gps/latest`
- Frontend: Leaflet.js для відображення карти

---

### Сценарій 17: Прогнозне технічне обслуговування (PRO)

**Актор:** Механік або менеджер автопарку  
**Передумови:** PRO підписка, ведеться облік пробігу  
**Тригер:** Наближення планового ТО

**Потік подій:**

1. Система щоденно перевіряє умову:
   ```go
   nextTOat := vehicle.LastServiceKm + vehicle.ServiceInterval
   remaining := nextTOat - vehicle.CurrentKm
   if remaining <= 500 {
     alert()
   }
   ```
2. За 500 км до ТО створюється сповіщення
3. Механік відкриває `/maintenance`
4. Бачить список транспорту з метриками:
   - Поточний пробіг: 48700 км
   - Наступне ТО: 50000 км
   - Залишилось: 1300 км ≈ 3 дні
   - Статус: 🟡 Скоро потрібне ТО
5. Система автоматично пропонує дату ТО за прогнозом використання
6. Механік створює заявку на обслуговування
7. Після виконання:
   - Оновлюється `last_service_km`
   - Обнуляється лічильник до наступного ТО

**Постумова:** Транспорт обслуговується вчасно, зменшується ризик поломок

**Переваги:**
- Запобігання аваріям через несправності
- Оптимізація графіку роботи автопарку
- Зменшення простоїв

**Технічна реалізація:**
- Handler: `MaintenanceHandler.GetSchedule()`
- Endpoint: `GET /api/maintenance/schedule`
- Algorithm: прогноз на основі середнього денного пробігу

---

## Сценарії палива та детекції шахрайства

### Сценарій 18: Реєстрація заправки палива

**Актор:** Водій або диспетчер  
**Передумови:** Транспорт заправлено  
**Тригер:** Фізична заправка

**Потік подій:**

1. Водій відкриває `/vehicles/:id`
2. Натискає "Зареєструвати заправку"
3. Вводить:
   - Кількість літрів
   - Вартість за літр
   - Дата і час
   - Опціонально: фото чеку
4. Система зберігає запис у `fuel_logs`
5. Оновлюється загальна статистика витрат
6. Якщо активовано PRO — запускається аналіз аномалій

**Постумова:** Заправка зафіксована в системі

---

### Сценарій 19: AI-детекція шахрайства з паливом (PRO)

**Актор:** Контролер або система (автоматично)  
**Передумови:** PRO підписка, історія заправок і GPS-треків  
**Тригер:** Нова заправка або щоденний аналіз

**Потік подій:**

1. Система збирає дані за останні 30 днів:
   - Всі заправки авто
   - Пройдений шлях (за GPS)
   - Нормативна витрата (л/100км)
2. Алгоритм `DetectFuelAnomalies()` виконує перевірки:

   **Аномалія 1: Перевитрата**
   ```go
   expectedFuel := distance * fuelRate / 100
   actualFuel := sum(refuels)
   if actualFuel > expectedFuel * 1.3 {
     flag = "EXCESSIVE_CONSUMPTION"
   }
   ```

   **Аномалія 2: Заправка без руху**
   ```sql
   SELECT * FROM fuel_logs f
   WHERE NOT EXISTS (
     SELECT 1 FROM gps_tracking g
     WHERE g.vehicle_id = f.vehicle_id
       AND g.timestamp BETWEEN f.date - INTERVAL '2 hours' 
                           AND f.date + INTERVAL '2 hours'
   )
   ```

   **Аномалія 3: Нереалістично часта заправка**
   ```go
   if timeBetweenRefuels < 6 hours && volume > 50L {
     flag = "SUSPICIOUS_FREQUENCY"
   }
   ```

3. Виявлені аномалії зберігаються в `fuel_anomalies`
4. Контролер відкриває `/fuel-anomalies` і бачить:
   ```
   ⚠️ ГАЗель #12: Заправка 80 л, але авто стояло на місці
   📍 Локація: 50.4501, 30.5234 (база)
   🕐 Час: 2026-04-15 14:30
   ```
5. Контролер:
   - Вимагає пояснення від водія
   - Або позначає як false positive (помилкове спрацювання)
6. Логується перевірка аномалії

**Постумова:** Виявлено потенційне шахрайство, ініційовано розслідування

**Економічний ефект:**
- За статистикою, зловживання паливом сягають 10-15% витрат
- Система окупається за 2-3 місяці

**Технічна реалізація:**
- Handler: `FuelHandler.GetAnomalies()`
- Endpoint: `GET /api/fuel/anomalies`
- Algorithm: статистичний аналіз з ML-елементами

---

## Сценарії аналітики та звітності

### Сценарій 20: Перегляд базового dashboard

**Актор:** Будь-який авторизований користувач  
**Передумови:** FREE або PRO підписка  
**Тригер:** Відкриття головної сторінки `/`

**Потік подій:**

1. Користувач авторизується і потрапляє на `/`
2. Система завантажує основні метрики:
   - Кількість заявок за статусами
   - Кількість ресурсів
   - Кількість транспорту
   - Загальний пробіг за місяць
3. Відображаються віджети:
   - "Активні заявки: 12"
   - "Критичні залишки: 3 ресурси"
   - "Транспорт в дорозі: 4 авто"
4. Швидкі дії:
   - Створити заявку
   - Переглянути інвентар
   - Відкрити GPS-карту (якщо PRO)

**Постумова:** Користувач отримав огляд стану системи

---

### Сценарій 21: Advanced KPI Dashboard (PRO)

**Актор:** Керівник або аналітик  
**Передумови:** PRO підписка  
**Тригер:** Відкриття `/kpi`

**Потік подій:**

1. Система обчислює KPI за останні 30 днів:

   **SLA (Service Level Agreement)**
   ```go
   onTime := countRequestsCompletedInTime()
   total := countAllRequests()
   sla := (onTime / total) * 100
   ```
   Відображається:
   - SLA: 87.3% 🟢 (норма > 85%)
   - Вчасно: 124 заявки
   - Прострочено: 18 заявок
   - Середня затримка: 2.4 години

   **TCO (Total Cost of Ownership)**
   ```go
   fuelCosts := sum(fuel_logs.amount)
   unitsShipped := sum(supply_requests.quantity)
   costPerUnit := fuelCosts / unitsShipped
   ```
   Відображається:
   - Загальні витрати на паливо: 45 780 грн
   - Доставлено одиниць: 2 340
   - Вартість доставки одиниці: 19.56 грн

   **Risk Score (Індекс ризику)**
   ```go
   criticalResources := countResourcesBelow(minQuantity)
   totalResources := countAllResources()
   riskPercent := (criticalResources / totalResources) * 100
   ```
   Відображається:
   - Ресурсів під ризиком: 8 з 125 (6.4%) 🟡

   **Depletion Forecast (Прогноз вичерпання)**
   ```go
   for resource in resources {
     avgDailyConsumption := calculateAvgConsumption(30days)
     daysLeft := resource.quantity / avgDailyConsumption
     if daysLeft < 7 {
       alert(resource, daysLeft)
     }
   }
   ```
   Відображається:
   - Вичерпаються за 7 днів: Бинт (4 дні), Джгут (6 днів)
   - Вичерпаються за 14 днів: Шини 215/75R17.5 (11 днів)

2. Графіки:
   - SLA за тиждень (лінійний)
   - TCO динаміка (стовпчиковий)
   - Розподіл ризиків за категоріями (кругова діаграма)

**Постумова:** Керівник отримав стратегічні метрики для прийняття рішень

**Переваги:**
- Швидке виявлення проблемних зон
- Обґрунтування інвестицій даними
- Порівняння ефективності між періодами

**Технічна реалізація:**
- Handler: `AnalyticsHandler.GetAdvancedKPIs()`
- Endpoint: `GET /api/analytics/kpi`
- Обчислення: агрегуючі SQL-запити з window functions

---

### Сценарій 22: Demand Forecasting — AI-прогноз попиту (PRO)

**Актор:** Керівник відділу постачання  
**Передумови:** PRO підписка, історія заявок за 90+ днів  
**Тригер:** Планування закупівель на наступний місяць

**Потік подій:**

1. Керівник відкриває `/analytics`
2. Обирає "Прогноз попиту"
3. Вказує ресурс (наприклад, "Дизпаливо")
4. Система аналізує:
   - Споживання за останні 90 днів
   - Сезонність (якщо є)
   - Тренди (зростання/спад)
5. Алгоритм `ForecastDemand()` використовує:
   - Moving Average для згладжування
   - Експоненціальне згладжування для тренду
   - Додавання стандартного відхилення для сценаріїв
6. Повертається прогноз на 30 днів:
   ```json
   {
     "resource": "Дизпаливо",
     "scenarios": {
       "low": 1200,      // оптимістичний
       "medium": 1680,   // реалістичний
       "high": 2100      // песимістичний
     },
     "confidence": 0.78,
     "recommendation": "Закупити 1700 л (+5% запас)"
   }
   ```
7. Відображається графік:
   - Історичне споживання (синя лінія)
   - Прогноз Medium (зелена лінія)
   - Коридор Low-High (заштрихована область)

**Постумова:** Планувальник знає, скільки замовляти, уникаючи дефіциту і надлишку

**Економічний ефект:**
- Зменшення заморожених коштів у надлишках на 20-30%
- Запобігання дефіциту (критичні ситуації -50%)

**Технічна реалізація:**
- Handler: `AnalyticsHandler.ForecastDemand()`
- Endpoint: `POST /api/analytics/forecast`
- Algorithm: Exponential Smoothing (Holt-Winters)

---

## Сценарії роботи з контрагентами

### Сценарій 23: Публікація зовнішньої заявки для волонтерів

**Актор:** Логіст або керівник  
**Передумови:** Потреба в ресурсі, якого немає у внутрішніх складах  
**Тригер:** Неможливість задовольнити потребу внутрішніми силами

**Потік подій:**

1. Логіст відкриває `/volunteer-requests`
2. Натискає "Створити публічну заявку"
3. Заповнює:
   - Назва потреби (наприклад, "100 шт Термобілизна розмір L")
   - Опис (деталі, вимоги якості)
   - Пріоритет
   - Підрозділ-одержувач
4. Система створює запис у `volunteer_requests` зі статусом `OPEN`
5. Заявка стає видимою для всіх користувачів з роллю `CONTRACTOR`
6. Надсилається email-нотифікація зареєстрованим волонтерам (якщо налаштовано)

**Постумова:** Заявка опублікована, волонтери можуть її взяти

**Технічна реалізація:**
- Handler: `VolunteerRequestHandler.Create()`
- Endpoint: `POST /api/volunteer-requests`

---

### Сценарій 24: Взяття заявки волонтером

**Актор:** CONTRACTOR (волонтер, підрядник)  
**Передумови:** Користувач зареєстрований як CONTRACTOR  
**Тригер:** Бачить заявку, яку може виконати

**Потік подій:**

1. Волонтер авторизується і відкриває `/volunteer-requests`
2. Бачить список відкритих заявок:
   ```
   🟢 100 шт Термобілизна L | Пріоритет: ВИСОКИЙ
   📍 Підрозділ: Київська філія
   📅 Створено: 2026-04-20
   ```
3. Натискає "Взяти в роботу"
4. Система:
   - Змінює статус на `TAKEN`
   - Встановлює `taken_by = user_id`
   - `taken_at = NOW()`
5. Заявка більше не відображається іншим волонтерам
6. Організація отримує нотифікацію "Заявку взяв волонтер [ім'я]"
7. Логується в audit log

**Постумова:** Волонтер зобов'язаний виконати заявку

---

### Сценарій 25: Доставка та прийомка від волонтера

**Актор:** Волонтер + Комірник організації  
**Передумови:** Волонтер доставив ресурс фізично  
**Тригер:** Прибуття волонтера на склад

**Потік подій:**

1. Волонтер позначає в системі "Доставлено"
   - Статус міняється на `DELIVERED`
2. Комірник організації отримує нотифікацію
3. Відкриває заявку і бачить кнопки:
   - ✅ Прийняти на баланс
   - ❌ Відхилити (якщо не відповідає вимогам)
4. При прийнятті:
   - Комірник вибирає:
     - **Новий ресурс**: вводить категорію, назву, одиниці
     - **Існуючий ресурс**: обирає зі списку
   - Вказує кількість
5. Система:
   - Створює/оновлює ресурс у `resources`
   - Додає кількість до балансу
   - Змінює статус заявки на `ACCEPTED`
   - `accepted_at = NOW()`
   - `accepted_by = комірник_id`
6. Волонтер отримує нотифікацію про прийняття
7. Логується прийомка

**Постумова:** Ресурс на балансі організації, заявка закрита

**Альтернативний потік — Відхилення:**
- Комірник натискає "Відхилити"
- Вводить причину (наприклад, "Неправильний розмір")
- Статус → `REJECTED`
- Волонтер отримує нотифікацію з причиною

**Технічна реалізація:**
- Handler: `VolunteerRequestHandler.Accept()`
- Endpoint: `POST /api/volunteer-requests/:id/accept`
- Транзакція: створення ресурсу + оновлення заявки

---

## Сценарії адміністрування платформи

### Сценарій 26: Управління підписками (Billing)

**Актор:** TENANT_ADMIN організації  
**Передумови:** Організація на FREE плані, хоче PRO  
**Тригер:** Необхідність доступу до PRO-фіч

**Потік подій:**

1. Адміністратор відкриває `/billing`
2. Бачить поточний план:
   ```
   Поточний план: FREE
   Активно до: Безстроково
   Користувачі: 12 / ∞
   Склади: 5 / ∞
   ```
3. Розділ "Доступні плани":
   - **FREE**: Базові можливості
   - **PRO**: $49/міс або $490/рік (-15%)
4. Натискає "Оновити до PRO"
5. Система показує різницю функцій:
   - ✅ Advanced KPI Dashboard
   - ✅ GPS Tracking
   - ✅ Fuel Anti-Fraud Detection
   - ✅ Demand Forecasting
   - ✅ Predictive Maintenance
6. Обирає період: Місяць / Рік
7. Вводить платіжні дані (реальна інтеграція з Stripe/Fondy)
8. Після оплати:
   - `tenants.subscription_tier = 'PRO'`
   - `tenants.subscription_expires_at = NOW() + 30 days`
9. Всі PRO-фічі стають доступними

**Постумова:** Організація має повний доступ до функціоналу

**Технічна реалізація:**
- Handler: `BillingHandler.UpgradePlan()`
- Endpoint: `POST /api/billing/upgrade`
- Middleware: `subscription.RequireTier("PRO")`

---

### Сценарій 27: Блокування користувача

**Актор:** TENANT_ADMIN  
**Передумови:** Користувач порушив політику  
**Тригер:** Рішення адміністратора

**Потік подій:**

1. Адміністратор відкриває `/admin/users`
2. Знаходить користувача
3. Натискає "Заблокувати"
4. Система:
   - Встановлює `blocked = true`
   - Інвалідує всі refresh токени користувача
5. При спробі логіну користувач отримує 403 Forbidden
6. Логується блокування в audit log

**Постумова:** Користувач не може авторизуватися

**Розблокування:** Аналогічна кнопка "Розблокувати"

---

### Сценарій 28: Audit Logs — журнал подій

**Актор:** SYSTEM_ADMIN або TENANT_ADMIN  
**Передумови:** Система логує всі важливі дії  
**Тригер:** Розслідування інциденту або аудит

**Потік подій:**

1. Адміністратор відкриває `/audit`
2. Бачить таблицю з записами:
   ```
   [2026-04-28 14:35:12] user@example.com | CREATE | SUPPLY_REQUEST | #1234
     → Створено заявку на 50 шт Бинти
   
   [2026-04-28 14:40:55] admin@org.mil | APPROVE | SUPPLY_REQUEST | #1234
     → Затверджено заявку
   
   [2026-04-28 15:12:03] warehouse@org.mil | UPDATE | RESOURCE | #567
     → Змінено кількість: 100 → 50
   ```
3. Фільтри:
   - За користувачем
   - За типом дії (CREATE, UPDATE, DELETE, APPROVE)
   - За типом сутності (USER, REQUEST, RESOURCE, VEHICLE)
   - За датою
4. Експорт в CSV (для звітів)

**Постумова:** Адміністратор має повну історію дій

**Технічна реалізація:**
- Repository: `AuditRepository.List()`
- Асинхронне логування: `go auditService.LogAction(...)`

---

## Висновки

### Реальні функціональні можливості системи (станом на 29.04.2026)

#### ✅ Повністю реалізовано:

1. **Аутентифікація та авторизація:**
   - JWT токени (access + refresh) ✅
   - Рольова модель (14 ролей) ✅
   - Multi-tenant ізоляція ✅
   - Invite tokens для реєстрації ✅
   - Password reset ✅

2. **Управління організаційною структурою:**
   - Створення підрозділів (4 рівні ієрархії) ✅
   - CRUD операції над units ✅
   - Tenant scoping ✅

3. **Управління користувачами:**
   - Створення користувачів з invite-токенами ✅
   - Матриця створення ролей (`RoleCreationMap`) ✅
   - Блокування/розблокування ✅
   - Перегляд видимих користувачів ✅

4. **Інвентар та ресурси:**
   - Створення категорій ✅
   - CRUD ресурсів ✅
   - Списання (write-off) ✅
   - Призначення (assign) ✅
   - Перегляд з фільтрами ✅
   - Excel імпорт (PRO) ✅

5. **Склади:**
   - CRUD складів ✅
   - Прив'язка до підрозділів ✅
   - Перегляд ресурсів по складах ✅

6. **Заявки на постачання:**
   - Створення заявок ✅
   - Затвердження/відхилення ✅
   - Статуси (PENDING → APPROVED → IN_TRANSIT → COMPLETED) ✅
   - Матриця погодження (`ApprovalMatrix`) ✅

7. **Заявки для підрядників:**
   - Створення публічних заявок ✅
   - Взяття в роботу (CONTRACTOR) ✅
   - Позначення доставлено ✅
   - Прийняття/відхилення організацією ✅
   - Скасування ✅

8. **Транспорт:**
   - CRUD транспортних засобів ✅
   - Реєстрація заправок ✅
   - Призначення водіїв ✅
   - Технічне обслуговування ✅
   - Історія ТО та водіїв ✅

9. **GPS трекінг (PRO):**
   - Запис координат ✅
   - Fleet map (карта флоту) ✅
   - Траєкторії руху ✅
   - Geofencing (геозони) ✅
   - Алерти виходу з геозони ✅
   - Статус флоту ✅

10. **Аналітика:**
    - Базовий dashboard ✅
    - Advanced KPI (PRO) ✅
    - Demand forecasting (PRO) ✅
    - Predictive maintenance (PRO) ✅
    - Fuel anomaly detection (PRO) ✅
    - Експорт інвентарю/палива ✅

11. **Підписки:**
    - Перевірка тарифів через middleware ✅
    - Блокування PRO функцій на FREE тарифі ✅
    - Аудит спроб несанкціонованого доступу ✅

12. **Аудит:**
    - Логування всіх критичних дій ✅
    - Перегляд audit logs ✅
    - Фільтрація по користувачам/діям ✅

---

#### ⚠️ Реалізовано з обмеженнями:

1. **Рольові права:**
   - Використовуються групи ролей, а не детальні permissions ⚠️
   - Неможливо дати часткові права в межах модуля ⚠️
   - Немає RBAC з permissions matrix ⚠️

2. **EMPLOYEE роль:**
   - НЕ МОЖЕ створювати заявки ❌ (це баг!)
   - Фактично роль "тільки перегляд" ⚠️

3. **Логісти:**
   - НЕ МОЖУТЬ керувати інвентарем напряму ❌
   - НЕ МОЖУТЬ створювати категорії ❌
   - Працюють тільки із заявками та складами ⚠️

4. **Комірники:**
   - НЕ МОЖУТЬ затверджувати заявки ❌
   - НЕ МОЖУТЬ створювати склади ❌
   - Тільки управління ресурсами ⚠️

5. **Audit logs:**
   - Доступні не тільки ADMIN ⚠️
   - Можуть переглядати всі з `UserCreatorRoles` ⚠️

---

#### ❌ НЕ реалізовано:

1. **Smart Distribution:**
   - Немає ендпоінту `/api/requests/smart-distribute` ❌
   - AI-розподіл заявок по транспорту відсутній ❌

2. **Платформне адміністрування:**
   - Немає `/platform-admin` панелі ❌
   - SYSTEM_ADMIN не має спеціальних ендпоінтів ⚠️

3. **Каскадне видалення:**
   - При видаленні підрозділу немає обробки залежностей ❌

4. **SLA моніторинг:**
   - Є `SLAMonitor` service, але UI та нотифікації не підтверджені ⚠️

5. **Email нотифікації:**
   - Сервіс `EmailService` існує, але інтеграція не протестована ⚠️

6. **Kiosk Terminal:**
   - ✅ **РЕАЛІЗОВАНО** на frontend (`/kiosk`)
   - Frontend компонент `KioskTerminal.tsx` працює ✅
   - Використовує існуючі API (`writeOffResource`) ✅
   - Доступ: ADMIN, комірники, DEPT_SUPERVISOR ✅
   - **БЕЗКОШТОВНО** - не потребує PRO підписки ✅
   - Функціонал:
     * Сканування штрих-кодів/QR-кодів (camera + manual input)
     * Пошук ресурсів за ID, barcode або назвою
     * Кошик для множинної видачі
     * Списання ресурсів зі складу
   - ❌ Немає окремих backend ендпоінтів (використовує стандартні `POST /api/inventory/resources/:id/write-off`)

---

### Рекомендації для виправлення багів:

#### 🔧 Пріоритет 1 (критичні):

```go
// 1. Дозволити EMPLOYEE створювати заявки
var SupplyRequestCreatorRoles = []UserRole{
    // ... існуючі ролі
    RoleEmployee, // ← ДОДАТИ
}

// 2. Дозволити логістам керувати інвентарем
var InventoryManagerRoles = []UserRole{
    // ... існуючі ролі
    RoleRegionLogistician,  // ← ДОДАТИ
    RoleBranchLogistician,  // ← ДОДАТИ
}
```

#### 🔧 Пріоритет 2 (важливі):

```go
// 3. (Опціонально) Дозволити комірникам затверджувати заявки
var SupplyRequestApproverRoles = []UserRole{
    // ... існуючі ролі
    RoleRegionStorekeeper,  // ← ДОДАТИ (дискусійно)
    RoleBranchStorekeeper,  // ← ДОДАТИ (дискусійно)
}

// 4. Обмежити audit logs тільки для адміністраторів
admin.GET("/audit-logs", 
    middleware.RequireAnyRole([]models.UserRole{
        models.RoleSystemAdmin,
        models.RoleTenantAdmin,
        models.RoleAdmin,
    }), 
    auditHandler.GetLogs)
```

#### 🔧 Пріоритет 3 (покращення):

1. Впровадити детальні permissions (RBAC)
2. Додати каскадне видалення з перевірками
3. Реалізувати Smart Distribution
4. Додати UI тести для кожної ролі

---

### Матриця тестування ролей:

Для кожної ролі потрібно протестувати:

| Роль | Може створювати заявки | Може затверджувати | Може керувати інвентарем | Може створювати склади | Може керувати транспортом |
|------|------------------------|---------------------|--------------------------|------------------------|---------------------------|
| **SYSTEM_ADMIN** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **TENANT_ADMIN** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **REGION_DIRECTOR** | ✅ | ✅ | ❌ | ✅ | ✅ |
| **REGION_LOGISTICIAN** | ✅ | ✅ | ❌ БАГ | ✅ | ✅ |
| **REGION_STOREKEEPER** | ✅ | ❌ | ✅ | ❌ | ❌ |
| **BRANCH_MANAGER** | ✅ | ✅ | ❌ | ✅ | ✅ |
| **BRANCH_LOGISTICIAN** | ✅ | ✅ | ❌ БАГ | ✅ | ✅ |
| **BRANCH_STOREKEEPER** | ✅ | ❌ | ✅ | ❌ | ❌ |
| **DEPT_MANAGER** | ✅ | ✅ | ❌ | ✅ | ✅ |
| **DEPT_SUPERVISOR** | ✅ | ❌ | ✅ | ❌ | ✅ |
| **TEAM_LEAD** | ✅ | ❌ | ❌ | ❌ | ❌ |
| **EMPLOYEE** | ❌ БАГ | ❌ | ❌ | ❌ | ❌ |
| **CONTRACTOR** | ❌ | ❌ | ❌ | ❌ | ❌ |

**❌ БАГ** - позначено критичні помилки в правах доступу

---

### Сценарії використання (User Journey)

#### Успішні сценарії:

1. **Комірник додає ресурс → Менеджер створює заявку → Логіст затверджує → Комірник комплектує** ✅

2. **Організація створює заявку для волонтера → CONTRACTOR бере → доставляє → Комірник приймає** ✅

3. **Логіст додає транспорт → реєструє заправку → переглядає аномалії палива (PRO)** ✅

4. **Директор відкриває KPI dashboard (PRO) → бачить SLA, TCO, прогноз попиту** ✅

#### Проблемні сценарії:

1. **EMPLOYEE хоче запитати ресурс → НЕ МОЖЕ створити заявку** ❌
   - **Workaround:** TEAM_LEAD створює заявку за нього

2. **REGION_LOGISTICIAN бачить дефіцит ресурсу → НЕ МОЖЕ додати ресурс** ❌
   - **Workaround:** Просить комірника додати

3. **BRANCH_STOREKEEPER створює заявку → чекає затвердження логіста → логіст на лікарняному → заявка зависає** ⚠️
   - **Workaround:** Дати комірникам право self-approve для невеликих заявок

---

### Підсумок для виправлення багів:

**Цей документ тепер відображає РЕАЛЬНИЙ стан системи**, включаючи:
- ✅ Що працює
- ⚠️ Що працює з обмеженнями
- ❌ Що не працює (баги)
- 🔧 Як виправити

Використовуйте цю документацію як:
1. **Керівництво для тестування** - перевірте кожну роль згідно матриці
2. **План виправлення** - виправте баги в пріоритетному порядку
3. **Документація для розробників** - зрозумійте реальну архітектуру
4. **Baseline для покращень** - знайте від чого відштовхуватися

**Найкритичніші виправлення:**

---

## 🆕 РЕКОМЕНДАЦІЇ З ПОКРАЩЕННЯ KIOSK TERMINAL

### Поточний функціонал:
✅ Швидке списання (видача) ресурсів  
✅ Сканування штрих-кодів/QR  
✅ Пошук за ID/barcode/назвою  
✅ Кошик для множинної видачі  

### 🎯 Рекомендована функція: ПРИЙМАННЯ ТОВАРУ

**Проблема:** Зараз Kiosk Terminal дозволяє тільки **видавати** ресурси зі складу, але не **приймати** їх.

**Рішення:** Додати режим "Прихід товару" в термінал.

#### Технічна реалізація:

**1. Frontend (KioskTerminal.tsx):**

```tsx
const [mode, setMode] = useState<'issue' | 'receive'>('issue');

// Новий метод для приймання
const handleReceiveAll = async () => {
  if (cart.length === 0) return;
  const loadingToast = toast.loading('Приймаємо товар...');
  try {
    for (const item of cart) {
      // Знайти ресурс або створити новий
      const resource = await api.inventory.findByBarcode(item.barcode);
      
      if (resource) {
        // Оновити існуючий
        await api.inventory.updateResource(resource.id, {
          quantity: resource.quantity + item.quantity
        });
      } else {
        // Створити новий ресурс
        await api.inventory.createResource({
          name: item.name,
          barcode: item.barcode,
          quantity: item.quantity,
          // ... інші поля
        });
      }
    }
    toast.success('Товар прийнято!', { id: loadingToast });
    setCart([]);
  } catch (error) {
    toast.error('Помилка приймання', { id: loadingToast });
  }
};

// UI перемикач режимів
<div className="mode-switcher">
  <button 
    className={mode === 'issue' ? 'active' : ''}
    onClick={() => setMode('issue')}
  >
    📤 Видача
  </button>
  <button 
    className={mode === 'receive' ? 'active' : ''}
    onClick={() => setMode('receive')}
  >
    📥 Прийом
  </button>
</div>

// Умовна кнопка підтвердження
{mode === 'issue' ? (
  <button onClick={handleIssueAll}>Видати все</button>
) : (
  <button onClick={handleReceiveAll}>Прийняти все</button>
)}
```

**2. Backend (додати новий ендпоінт або використати існуючий):**

```go
// Опція 1: Використати існуючий PATCH /api/inventory/resources/:id
// (додає до поточної кількості)

// Опція 2: Створити спеціальний ендпоінт
inv.POST("/resources/:id/receive", invHandler.ReceiveStock)

func (h *InventoryHandler) ReceiveStock(c *gin.Context) {
    var req struct {
        Quantity int `json:"quantity" binding:"required,min=1"`
        Note     string `json:"note"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    resourceID := c.Param("id")
    
    // Збільшити кількість
    err := h.service.IncreaseQuantity(c.Request.Context(), resourceID, req.Quantity)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // Логування в audit
    h.auditService.LogAction(c.Request.Context(), &models.AuditLog{
        UserID:     c.GetString("user_id"),
        Action:     "RECEIVE_STOCK",
        EntityType: "RESOURCE",
        EntityID:   resourceID,
        Details:    fmt.Sprintf("Прийнято %d одиниць. %s", req.Quantity, req.Note),
    })
    
    c.JSON(200, gin.H{"message": "Товар прийнято"})
}
```

**3. Додаткові покращення:**

#### 3.1. Генерація QR-кодів для нових товарів
```tsx
import QRCode from 'qrcode';

const generateQRCode = async (resourceId: string) => {
  const qrCodeDataUrl = await QRCode.toDataURL(`Omnilog-resource:${resourceId}`);
  return qrCodeDataUrl; // Показати користувачу або роздрукувати
};
```

#### 3.2. Друк наліпок зі штрих-кодами
```tsx
const printLabel = (resource: Resource) => {
  const printWindow = window.open('', '_blank');
  printWindow?.document.write(`
    <html>
      <head><title>Наліпка ${resource.name}</title></head>
      <body>
        <h2>${resource.name}</h2>
        <img src="${resource.qrCodeUrl}" />
        <p>Barcode: ${resource.barcode}</p>
      </body>
    </html>
  `);
  printWindow?.print();
};
```

#### 3.3. Масовий прийом з Excel/CSV
```tsx
const importFromCSV = async (file: File) => {
  const text = await file.text();
  const rows = text.split('\n').map(row => row.split(','));
  
  for (const [barcode, name, quantity] of rows) {
    await api.inventory.receiveStock({
      barcode,
      name,
      quantity: parseInt(quantity)
    });
  }
};
```

#### 3.4. Фотографування товару при прийманні
```tsx
const [photo, setPhoto] = useState<string | null>(null);

const capturePhoto = async () => {
  const stream = await navigator.mediaDevices.getUserMedia({ video: true });
  const video = document.createElement('video');
  video.srcObject = stream;
  // ... захват кадру як base64
  setPhoto(capturedImageBase64);
};
```

---

### 📊 Порівняння функціоналу:

| Функція | Поточний стан | З покращеннями |
|---------|---------------|----------------|
| **Видача товару** | ✅ Працює | ✅ Працює |
| **Прийом товару** | ❌ Відсутній | ✅ Додати |
| **QR-коди** | ✅ Сканування | ✅ + Генерація та друк |
| **Пошук** | ✅ ID/barcode/назва | ✅ Працює |
| **Множинна операція** | ✅ Кошик | ✅ Працює |
| **Фото товару** | ❌ Відсутній | ✅ Додати |
| **Імпорт CSV/Excel** | ❌ Відсутній | ✅ Додати |
| **Історія операцій** | ⚠️ Через audit logs | ✅ Спеціальна сторінка в терміналі |

---

### 🎯 Пріоритет впровадження:

**Високий пріоритет:**
1. ✅ Режим "Прийом товару" (receive mode)
2. ✅ Оновлення кількості існуючих ресурсів

**Середній пріоритет:**
3. ⚠️ Генерація QR-кодів для нових товарів
4. ⚠️ Друк наліпок

**Низький пріоритет:**
5. 📸 Фотографування товару
6. 📊 Імпорт з CSV/Excel (це вже є в PRO функціях для звичайного інтерфейсу)

---

### ✅ Переваги додавання функції "Прийом":

1. **Швидкість:** Комірник може приймати товар без виходу з термінального режиму
2. **Зручність:** Той самий інтерфейс для видачі та прийому
3. **QR-коди:** Нові товари відразу отримують QR для майбутнього сканування
4. **Аудит:** Всі операції логуються
5. **Мобільність:** Термінал можна використовувати на складі з планшета

---

**Найкритичніші виправлення:**
1. Додати `RoleEmployee` до `SupplyRequestCreatorRoles`
2. Додати `RoleRegionLogistician, RoleBranchLogistician` до `InventoryManagerRoles`
3. Протестувати всі ролі згідно матриці

---

**Документ оновлено:** 29 квітня 2026 року  
**Версія:** 2.0 (Reality Check Edition)  
**Статус:** Відображає реальний стан коду, а не ідеальну модель

---

## 📋 ТАБЛИЦЯ СЦЕНАРІЇВ ВЕБ-ЗАСТОСУНКУ

### Сценарій 1

| Поле | Опис |
|------|------|
| **№** | 1 |
| **Назва та код сценарію** | **Назва:** Авторизація користувача в системі<br>**Код:** AUTH-001<br>**Актори:** Користувач (будь-яка роль)<br>**Пріоритет:** 1<br>**Частота:** Щоденно, кілька разів на день |
| **Опис сценарію** | **Дія:** Користувач входить в систему за допомогою email та пароля<br><br>**Попередні умови:**<br>- Користувач зареєстрований в системі<br>- Користувач має активний акаунт<br>- База даних доступна<br><br>**Подальші умови:**<br>- Користувач успішно авторизований<br>- Створено access та refresh токени<br>- Відкрито головну сторінку відповідно до ролі<br><br>**Нормальна послідовність:**<br>1. Користувач відкриває сторінку входу<br>2. Вводить email та пароль<br>3. Натискає кнопку "Увійти"<br>4. Система перевіряє credentials<br>5. Система генерує JWT токени<br>6. Користувач перенаправляється на головну панель<br><br>**Альтернативний хід:**<br>- Користувач може скористатись функцією "Запам'ятати мене"<br>- Можливе автоматичне оновлення токена при його закінченні<br><br>**Винятки:**<br>- Невірний email або пароль – показати повідомлення про помилку<br>- Акаунт заблокований – показати відповідне повідомлення<br>- Проблеми з мережею – показати помилку з'єднання<br>- База даних недоступна – показати системну помилку<br><br>**Бізнес-правила:**<br>- Токен дійсний 15 хвилин (access) та 7 днів (refresh)<br>- Після 3 невдалих спроб входу – затримка 5 хвилин<br>- Всі спроби входу логуються в audit trail<br>- Паролі зберігаються в хешованому вигляді |

### Сценарій 2

| Поле | Опис |
|------|------|
| **№** | 2 |
| **Назва та код сценарію** | **Назва:** Створення запиту на постачання<br>**Код:** REQ-001<br>**Актори:** Командир підрозділу, Заступник командира, Співробітник<br>**Пріоритет:** 2<br>**Частота:** Щоденно, 5-15 разів на день |
| **Опис сценарію** | **Дія:** Користувач створює запит на необхідні матеріально-технічні засоби<br><br>**Попередні умови:**<br>- Користувач авторизований в системі<br>- Користувач має роль з правами створення запитів<br>- Користувач має активний підрозділ<br><br>**Подальші умови:**<br>- Запит створено в статусі "Pending"<br>- Відправлено повідомлення логісту відповідного рівня<br>- Запис додано до журналу аудиту<br><br>**Нормальна послідовність:**<br>1. Користувач відкриває розділ "Запити"<br>2. Натискає "Створити запит"<br>3. Вибирає категорію ресурсу<br>4. Вибирає конкретний ресурс зі списку<br>5. Вказує необхідну кількість<br>6. Вказує пріоритет (Low/Medium/High/Critical)<br>7. Додає опис та обґрунтування<br>8. Натискає "Відправити запит"<br>9. Система валідує дані<br>10. Створює запит та повідомляє логіста<br><br>**Альтернативний хід:**<br>- Користувач може зберегти запит як чернетку<br>- Можливе прикріплення фото або документів<br>- Можна додати кілька позицій в один запит<br><br>**Винятки:**<br>- Некоректна кількість – показати помилку валідації<br>- Ресурс не знайдено – запропонувати створити новий<br>- Підрозділ неактивний – заборонити створення<br>- Недостатні права – показати помилку доступу<br><br>**Бізнес-правила:**<br>- Запити з пріоритетом Critical потребують додаткового підтвердження<br>- Час обробки залежить від пріоритету: Critical (4 години), High (1 день), Medium (3 дні), Low (тиждень)<br>- Запити можна редагувати тільки в статусі Pending<br>- Всі зміни статусу запиту логуються |

### Сценарій 3

| Поле | Опис |
|------|------|
| **№** | 3 |
| **Назва та код сценарію** | **Назва:** Управління інвентарем на складі<br>**Код:** INV-001<br>**Актори:** Логіст центрального рівня, Логіст регіонального рівня, Логіст галузевого рівня<br>**Пріоритет:** 1<br>**Частота:** Щоденно, 20-50 разів на день |
| **Опис сценарію** | **Дія:** Логіст керує запасами на складі: додає, редагує, списує ресурси<br><br>**Попередні умови:**<br>- Логіст авторизований в системі<br>- Логіст має права управління інвентарем<br>- Склад створений та активний<br><br>**Подальші умови:**<br>- Операція з інвентарем виконана<br>- Оновлено залишки на складі<br>- Створено запис в історії операцій<br>- Якщо залишок критичний – відправлено алерт<br><br>**Нормальна послідовність:**<br>1. Логіст відкриває розділ "Інвентар"<br>2. Вибирає потрібний склад<br>3. Обирає дію: додати/редагувати/списати<br>4. **При додаванні:**<br>   - Вибирає або створює ресурс<br>   - Вказує кількість<br>   - Вводить штрих-код (опціонально)<br>   - Вказує термін придатності<br>   - Додає примітки<br>5. **При списанні:**<br>   - Вибирає ресурс<br>   - Вказує кількість списання<br>   - Вибирає причину (видано, пошкоджено, прострочено)<br>   - Опціонально вказує одержувача<br>6. Підтверджує операцію<br>7. Система оновлює залишки<br><br>**Альтернативний хід:**<br>- Використання сканера штрих-кодів для швидкого пошуку<br>- Масове імпортування через CSV/Excel<br>- Швидке списання через Kiosk Terminal<br><br>**Винятки:**<br>- Недостатньо товару для списання – показати помилку<br>- Некоректний штрих-код – запропонувати ввести вручну<br>- Прострочений термін придатності – показати попередження<br>- Відсутні права на склад – заборонити операцію<br><br>**Бізнес-правила:**<br>- Логіст може керувати лише складами свого рівня доступу<br>- При залишку нижче min_threshold – автоматичне повідомлення<br>- Критичні операції (велике списання) потребують підтвердження<br>- Всі операції логуються з інформацією про користувача<br>- Неможливо видалити ресурс з ненульовим залишком |

### Сценарій 4

| Поле | Опис |
|------|------|
| **№** | 4 |
| **Назва та код сценарію** | **Назва:** Робота з Kiosk Terminal для швидкого списання ресурсів<br>**Код:** KIOSK-001<br>**Актори:** Комірник складу, Логіст регіонального рівня, Логіст галузевого рівня<br>**Пріоритет:** 2<br>**Частота:** Щоденно, 30-100 разів на день |
| **Опис сценарію** | **Дія:** Швидке списання (видача) ресурсів зі складу через спрощений інтерфейс термінала з підтримкою сканування штрих-кодів<br><br>**Попередні умови:**<br>- Користувач авторизований в Kiosk Terminal режимі<br>- Термінал підключений до сканера штрих-кодів<br>- Склад активний та має наявні ресурси<br>- Користувач має права на списання ресурсів<br><br>**Подальші умови:**<br>- Ресурс успішно списаний зі складу<br>- Оновлено залишки в базі даних<br>- Створено audit log запис<br>- Роздруковано чек (опціонально)<br>- Операція відображена в історії<br><br>**Нормальна послідовність:**<br>1. Користувач входить в Kiosk Terminal режим<br>2. Система переключається на спрощений повноекранний інтерфейс<br>3. Користувач сканує штрих-код ресурсу або вводить ID вручну<br>4. Система знаходить ресурс та відображає:<br>   - Назву ресурсу<br>   - Поточні залишки<br>   - Мінімальний поріг<br>   - Розташування на складі<br>5. Користувач вводить кількість для списання<br>6. Система валідує наявність<br>7. Користувач вказує причину списання:<br>   - Видано підрозділу (вибір підрозділу)<br>   - Передано на інший склад<br>   - Пошкоджено/зіпсовано<br>   - Використано<br>8. Підтверджує операцію (Enter або кнопка "Списати")<br>9. Система:<br>   - Оновлює залишки<br>   - Створює запис в історії<br>   - Показує повідомлення про успіх<br>   - Якщо залишок < min_threshold – показує попередження<br>10. Можливість роздрукувати накладну (якщо є принтер)<br>11. Автоматичне повернення до початкового екрану через 3 секунди<br><br>**Альтернативний хід:**<br>- Масове списання: сканування кількох позицій підряд<br>- Використання клавіатурних скорочень для швидкості<br>- Режим "Експрес": пропуск підтвердження для дрібних позицій<br>- Пошук по частковій назві (якщо немає штрих-коду)<br>- Збереження "улюблених" операцій для повторення<br><br>**Винятки:**<br>- Штрих-код не розпізнано – запропонувати ввести вручну<br>- Ресурс не знайдено – запропонувати створити новий<br>- Недостатньо кількості – показати помилку з наявним залишком<br>- Сканер не підключений – автоматичне переключення на ручний ввід<br>- Втрата з'єднання з БД – зберегти в локальну чергу та синхронізувати пізніше<br>- Критично низький залишок – запит додаткового підтвердження<br><br>**Бізнес-правила:**<br>- Режим Kiosk доступний лише на складських локаціях<br>- Всі операції логуються з часовою міткою та користувачем<br>- При залишку < min_threshold – автоматичне створення алерту<br>- Неможливо списати більше ніж є в наявності<br>- Операції групуються в "зміни" (початок/кінець робочого дня)<br>- Підтримка офлайн-режиму з синхронізацією<br>- Автоматичний вихід після 5 хвилин неактивності<br>- Підтримка декількох мов інтерфейсу<br>- Великі кнопки та текст для зручності роботи на планшетах |

### Сценарій 5

| Поле | Опис |
|------|------|
| **№** | 5 |
| **Назва та код сценарію** | **Назва:** Обробка запиту від волонтерів (публічний інтерфейс)<br>**Код:** VOLUNTEER-001<br>**Актори:** Волонтер (неавторизований/зовнішній користувач), Логіст центрального рівня<br>**Пріоритет:** 2<br>**Частота:** 5-20 разів на тиждень |
| **Опис сценарію** | **Дія:** Волонтер подає запит на допомогу через публічну форму, логіст переглядає та обробляє такі запити<br><br>**Попередні умови:**<br>- Публічна форма доступна за URL без авторизації<br>- reCAPTCHA налаштована для захисту від спаму<br>- Email-сервіс налаштований для сповіщень<br>- Логіст має права на перегляд волонтерських запитів<br><br>**Подальші умови:**<br>- Запит створено в системі<br>- Волонтер отримав email з підтвердженням<br>- Логісти отримали сповіщення про новий запит<br>- Запиту присвоєно унікальний tracking номер<br><br>**Нормальна послідовність (Волонтер):**<br>1. Волонтер відкриває публічну сторінку `/volunteer-request`<br>2. Заповнює форму:<br>   - ПІБ або назва організації<br>   - Email для зворотного зв'язку<br>   - Телефон (опціонально)<br>   - Категорія запиту (медицина, продукти, одяг, техніка, інше)<br>   - Опис потреби (текстове поле 500 символів)<br>   - Бажана кількість<br>   - Регіон/місто доставки<br>   - Пріоритет (звичайний/терміновий)<br>3. Додає фото або документи (опціонально, до 5 файлів)<br>4. Проходить перевірку reCAPTCHA<br>5. Натискає "Відправити запит"<br>6. Система валідує дані<br>7. Створює запит зі статусом "New"<br>8. Показує tracking номер (наприклад, VOL-2026-0042)<br>9. Волонтер отримує email з:<br>   - Підтвердженням прийняття запиту<br>   - Tracking номером<br>   - Очікуваним часом обробки<br>   - Посиланням для відстеження статусу<br><br>**Нормальна послідовність (Логіст):**<br>1. Логіст отримує email про новий волонтерський запит<br>2. Заходить в розділ "Волонтерські запити"<br>3. Бачить список запитів з фільтрами:<br>   - Статус (New/In Review/Approved/Rejected/Completed)<br>   - Регіон<br>   - Категорія<br>   - Дата створення<br>4. Відкриває конкретний запит<br>5. Переглядає всю інформацію та прикріплені файли<br>6. Приймає рішення:<br>   - **Схвалити:** Створює внутрішній supply request<br>   - **Відхилити:** Вказує причину<br>   - **Запитати більше інформації:** Відправляє email волонтеру<br>7. Змінює статус запиту<br>8. Додає внутрішні примітки для команди<br>9. Волонтер отримує email з рішенням<br>10. Якщо схвалено – запит прив'язується до внутрішнього supply request<br>11. Волонтер може відстежувати статус за tracking номером<br><br>**Альтернативний хід:**<br>- Волонтер може оновити свій запит до моменту обробки<br>- Можливість завантажити декілька фото (пошкодження, документи)<br>- Автоматична геолокація для визначення регіону<br>- Масова обробка декількох запитів логістом<br>- Експорт списку запитів в Excel<br><br>**Винятки:**<br>- Невалідний email – показати помилку<br>- reCAPTCHA не пройдена – заборонити відправку<br>- Файл занадто великий (>5MB) – показати помилку<br>- Спам-фільтр спрацював – помістити в карантин<br>- Email не доставлено – зберегти в черзі повторних спроб<br>- Tracking номер не знайдено – показати помилку<br><br>**Бізнес-правила:**<br>- Публічна форма працює без авторизації (доступ для всіх)<br>- Обов'язкова перевірка reCAPTCHA v3 (score > 0.5)<br>- Максимум 3 запити з однієї IP за годину (anti-spam)<br>- Термінові запити виділяються червоним в системі<br>- Час обробки: термінові (24 години), звичайні (72 години)<br>- Email сповіщення при кожній зміні статусу<br>- Волонтер може відстежувати статус за публічним посиланням<br>- Історія всіх змін статусу зберігається<br>- Автоматичне закриття неактивних запитів через 30 днів<br>- Статистика по волонтерським запитам в адмін-панелі<br>- GDPR compliance: можливість видалення персональних даних |

### Сценарій 6

| Поле | Опис |
|------|------|
| **№** | 6 |
| **Назва та код сценарію** | **Назва:** Реєстрація нової організації в системі<br>**Код:** TENANT-REG-001<br>**Актори:** Представник організації (новий користувач, не авторизований)<br>**Пріоритет:** 1<br>**Частота:** 2-5 разів на тиждень |
| **Опис сценарію** | **Дія:** Представник організації реєструє нову організацію (tenant) в системі та стає її першим адміністратором<br><br>**Попередні умови:**<br>- Публічна сторінка реєстрації доступна без авторизації<br>- База даних доступна<br>- Email-сервіс налаштований для верифікації<br>- Користувач має унікальний email<br><br>**Подальші умови:**<br>- Створено новий tenant (організацію) в системі<br>- Створено першого користувача з роллю TENANT_ADMIN<br>- Відправлено email з підтвердженням реєстрації<br>- Користувач може увійти та почати роботу<br>- Організація отримала FREE підписку за замовчуванням<br><br>**Нормальна послідовність:**<br>1. Користувач відкриває публічну сторінку реєстрації (`/register`)<br>2. Заповнює форму реєстрації організації:<br>   **Інформація про організацію:**<br>   - Повна назва організації (наприклад, "72-га окрема механізована бригада")<br>   - Скорочена назва (наприклад, "72 ОМБр")<br>   - Код ЄДРПОУ/ID (опціонально)<br>   - Тип організації (Військовий підрозділ / Волонтерська організація / Державна установа / Інше)<br>   **Інформація про адміністратора:**<br>   - ПІБ представника<br>   - Email (буде використовуватись для входу)<br>   - Пароль (мінімум 8 символів, включає цифри та спецсимволи)<br>   - Підтвердження пароля<br>   - Посада в організації<br>   - Телефон (опціонально)<br>   **Додаткова інформація:**<br>   - Регіон розташування<br>   - Місто<br>   - Кількість співробітників (приблизно)<br>3. Погоджується з умовами використання та політикою конфіденційності (checkbox)<br>4. Проходить перевірку reCAPTCHA<br>5. Натискає кнопку "Зареєструватися"<br>6. Система валідує дані:<br>   - Перевіряє унікальність email<br>   - Перевіряє складність пароля<br>   - Перевіряє обов'язкові поля<br>   - Перевіряє reCAPTCHA score<br>7. При успішній валідації система:<br>   - Створює новий tenant в таблиці `tenants`<br>   - Хешує пароль (bcrypt)<br>   - Створює користувача з роллю `TENANT_ADMIN`<br>   - Прив'язує користувача до tenant<br>   - Встановлює `subscription_tier = "FREE"`<br>   - Генерує токен верифікації email<br>8. Відправляє email з посиланням для верифікації:<br>   ```<br>   Тема: Підтвердіть реєстрацію в Omnilog<br>   <br>   Вітаємо!<br>   <br>   Ваша організація "[Назва]" успішно зареєстрована.<br>   <br>   Для завершення реєстрації підтвердіть email:<br>   [Посилання на верифікацію]<br>   <br>   Ваші дані для входу:<br>   Email: [email]<br>   <br>   З повагою, команда Omnilog<br>   ```<br>9. Показує сторінку "Перевірте email"<br>10. Користувач відкриває email та клікає на посилання<br>11. Система верифікує токен та активує акаунт:<br>    - Встановлює `email_verified = true`<br>    - Показує повідомлення "Email підтверджено!"<br>12. Користувач перенаправляється на сторінку входу<br>13. Вводить email та пароль<br>14. Система авторизує та перенаправляє на онбординг<br>15. **Онбординг (Welcome Tour):**<br>    - Крок 1: "Створіть перший підрозділ"<br>    - Крок 2: "Додайте перший склад"<br>    - Крок 3: "Запросіть колег" (генерація invite tokens)<br>    - Крок 4: "Ознайомтесь з можливостями FREE плану"<br>16. Після онбордингу користувач потрапляє на головну панель<br><br>**Альтернативний хід:**<br>- **Реєстрація через Google/Microsoft SSO:**<br>  * Користувач обирає "Увійти через Google"<br>  * Система отримує email та ім'я від провайдера<br>  * Запитує додаткову інформацію про організацію<br>  * Email вже підтверджений автоматично<br>- **Реєстрація за запрошенням:**<br>  * Користувач отримав invite link від SYSTEM_ADMIN<br>  * Токен містить попередньо заповнену назву організації<br>  * Прискорена реєстрація з меншою кількістю полів<br>- **Пропуск онбордингу:**<br>  * Користувач може натиснути "Пропустити" на Welcome Tour<br>  * Онбординг можна пройти пізніше з налаштувань<br><br>**Винятки:**<br>- Email вже зареєстрований – показати помилку "Email вже використовується. Спробуйте увійти або відновити пароль"<br>- Слабкий пароль – показати підказку: "Пароль має містити мінімум 8 символів, цифри та спецсимволи"<br>- reCAPTCHA не пройдена – показати помилку "Не вдалося перевірити, що ви не робот"<br>- Назва організації вже існує – запропонувати додати унікальний ідентифікатор<br>- Email не доставлено – показати кнопку "Відправити повторно"<br>- Токен верифікації застарілий (>24 години) – запропонувати створити новий<br>- Помилка з'єднання з БД – показати "Технічні роботи, спробуйте пізніше"<br>- Поле "Умови використання" не відмічено – заборонити реєстрацію<br><br>**Бізнес-правила:**<br>- Кожна організація отримує унікальний `tenant_id` (UUID)<br>- За замовчуванням всі нові організації отримують FREE план<br>- Перший користувач автоматично стає TENANT_ADMIN<br>- Email має бути унікальним в межах всієї системи<br>- Пароль хешується за допомогою bcrypt (cost factor 10)<br>- Токен верифікації email дійсний 24 години<br>- reCAPTCHA v3 з мінімальним score 0.5<br>- Організація може мати лише одного першого адміністратора<br>- Після реєстрації організація ізольована (multi-tenant isolation)<br>- Неможливо видалити tenant з активними користувачами<br>- Дані організації можуть бачити лише користувачі цієї організації<br>- SYSTEM_ADMIN може бачити всі організації<br>- Назва організації має бути від 3 до 255 символів<br>- Всі дії реєстрації логуються для аудиту<br>- При помилці реєстрації транзакція відкочується (rollback) |

---

**Таблиця сценаріїв створена:** 30 квітня 2026 року  
**Всього сценаріїв:** 6  
**Охоплені функціональні блоки:** Авторизація, Запити на постачання, Управління інвентарем, Kiosk Terminal, Волонтерські запити, Реєстрація організації

**Технічні особливості:**
- 🎯 **Kiosk Terminal** - унікальна функція для швидкого складського обліку з підтримкою штрих-кодів
- 🌐 **Публічний волонтерський інтерфейс** - демонструє роботу з неавторизованими користувачами та інтеграцію зовнішніх запитів
- 🏢 **Multi-tenant архітектура** - реєстрація організацій з повною ізоляцією даних
- 🔒 **Anti-spam захист** - reCAPTCHA та rate limiting
- 📧 **Email notifications** - автоматичні сповіщення для волонтерів та верифікація email
- 📊 **Tracking система** - відстеження статусу публічних запитів
- 💾 **Offline-first** - Kiosk Terminal працює з локальною синхронізацією
- 🚀 **Онбординг процес** - Welcome Tour для нових організацій
- 🔐 **SSO підтримка** - можливість реєстрації через Google/Microsoft
- 🛡️ **Транзакційна безпека** - rollback при помилках реєстрації
