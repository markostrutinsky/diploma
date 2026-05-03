# Seed DB — сідер для Omnilog / OmniLog

Скрипт заповнює БД синтетичними даними так, щоб одразу можна було протестувати
**і безкоштовний (BASIC), і платний (PRO/ENTERPRISE)** функціонал.

## Структура, яку створює скрипт

| Регіон (unit type = `REGION`) | Тариф | Для чого            |
|-------------------------------|-------|---------------------|
| Регіон «Захід»                | `PRO` | Повний преміум      |
| Регіон «Центр»                | `BASIC` | Має бачити paywall 402 на PRO-ендпоінтах |
| Регіон «Схід»                 | `BASIC` | 9 складів із 10 (близько до ліміту) |
| Регіон «Тест-ENTERPRISE»      | `ENTERPRISE` | Необмежений тариф |

У кожному регіоні — гілки `BRANCH → DEPARTMENT → TEAM` (там, де сенс), юзери
з усіма ролями (`REGION_DIRECTOR`, `BRANCH_MANAGER`, `DEPT_MANAGER`,
`TEAM_LEAD`, `_LOGISTICIAN`, `_STOREKEEPER`, `DEPT_SUPERVISOR`, `EMPLOYEE`),
а також 3 незалежних `CONTRACTOR`-и (волонтери) без `unit_id`.

## Паролі

Для всіх створених юзерів — `password123` (bcrypt, сумісно з Go-сервером).

Ключові акаунти для демо:

| Email | Роль | Тариф (через unit) |
|-------|------|--------------------|
| `admin@Omnilog.local`              | ADMIN              | bypass (див. нижче) |
| `director.west@Omnilog.local`      | REGION_DIRECTOR    | PRO |
| `logist.west@Omnilog.local`        | REGION_LOGISTICIAN | PRO |
| `director.center@Omnilog.local`    | REGION_DIRECTOR    | BASIC |
| `director.east@Omnilog.local`      | REGION_DIRECTOR    | BASIC (біля ліміту) |
| `director.test@Omnilog.local`      | REGION_DIRECTOR    | ENTERPRISE |
| `contractor1@Omnilog.local`        | CONTRACTOR         | — |

## Як себе поводить адмін (важливий нюанс)

- У middleware `RequireSubscriptionTier` (файл
  `Omnilog_backend/internal/middleware/subscription.go`) для ролі `ADMIN` є
  явний **bypass** — адмін завжди отримує доступ до PRO/ENTERPRISE фіч,
  навіть якщо `unit_id = NULL`. Це задумано для підтримки/демо.
- `LimitationService` (файл `internal/services/limitation_service.go`) bypass-у
  для адміна **не має**: коли адмін створює склад у BASIC-регіоні, діє ліміт
  саме того регіону (10 складів у BASIC, 100 у PRO).

Тобто сценарій, який тебе хвилював — «три регіони, один з PRO, інші на BASIC»
— працює **без конфлікту**: тариф визначається за `unit_id` користувача і
лазить вгору по дереву до кореня (`WITH RECURSIVE`). Кожен з трьох регіонів —
окремий корінь (`parent_id = NULL`), тож юзери одного регіону не успадковують
тариф іншого.

## Запуск

```bash
cd scripts/seed_db

# перший раз — створити venv і залежності
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

# PostgreSQL з docker-compose має бути піднятий:
#   docker compose up -d postgres

# заповнити (ідемпотентно, не чистить):
python seed.py

# повне очищення і перезаповнення:
python seed.py --reset

# інший DSN:
python seed.py --dsn 'postgres://postgres:postgres@localhost:5432/omnilog'
```

Після сіду можна логінитись у фронті будь-яким з акаунтів у таблиці вище з
паролем `password123`.

## Що саме сідиться

- `units` — 4 регіони + їхні філії/відділи/команди.
- `users` — ~30 користувачів з усіх ролей, плюс один `BLOCKED` та один `PENDING`.
- `resource_categories` — 7 категорій.
- `warehouses` — по 1–9 на регіон (на Сході навмисно 9/10 для демо ліміту BASIC).
- `resources` — десятки ресурсів на кожному складі, різні категорії/кількості.
- `vehicles` — по 2–6 авто на регіон, з водіями.
- `fuel_records` — історія заправок за 30 днів, плюс одна штучна **аномалія**
  (PRO-аналітика має її показати).
- `maintenance_records` — випадкові записи ТО.
- `supply_requests` — заявки у різних статусах (PENDING / APPROVED / REJECTED / COMPLETED).
- `contractor_requests` — волонтерські заявки (OPEN / IN_PROGRESS / DELIVERED / COMPLETED).
- `geofences` — 2 геозони на PRO- та ENTERPRISE-регіонах.
- `gps_locations` — треки авто Заходу/Тест-регіону (для PRO GPS-фічі).

## Сценарії для тестування

1. **BASIC paywall.** Залогінься як `director.center@Omnilog.local` і спробуй
   `GET /api/analytics/dashboard` — повинен повернутись **HTTP 402** з текстом
   «Функція доступна тільки на тарифі: PRO».
2. **PRO функціонал.** Той самий запит від `director.west@Omnilog.local` має
   повернути дані аналітики.
3. **Адмін як супер'юзер.** `admin@Omnilog.local` отримує **усі** PRO-фічі,
   незважаючи на `unit_id = NULL`.
4. **Ліміти BASIC.** `director.east@Omnilog.local` спробує створити 2-й склад
   (у регіоні вже 9/10) — пройде. Третій спробувати — отримає помилку ліміту.
5. **ENTERPRISE без меж.** `director.test@Omnilog.local` — жодних обмежень.
6. **Аналітика аномалій пального.** PRO-юзер запитує `GET /api/analytics/fuel-anomalies`
   і бачить штучно проставлений `is_anomaly = TRUE`.
7. **GPS-трекінг (PRO).** PRO-юзер запитує
   `GET /api/gps/fleet-map` — бачить 20 останніх точок на кожне авто Заходу/Тест-регіону,
   а `GET /api/gps/trajectory?vehicle_id=<UUID>&start_time=...&end_time=...` повертає трек.
