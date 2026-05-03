#!/usr/bin/env python3
"""
Omnilog / OmniLog — скрипт заповнення тестової БД.

Створює багаторегіональну структуру (single-tenant, але з ізольованими
піддеревами одиниць), щоб можна було перевірити:
  - PRO-функціонал у регіоні з тарифом PRO;
  - paywall (HTTP 402) у регіонах з тарифом BASIC;
  - ліміти BASIC (склади, юзери, ресурси, авто) у близькому до граничного стані;
  - адміна-супер'юзера (bypass premium-middleware, але ліміти піддерев лишаються);
  - ENTERPRISE — повний необмежений тариф;
  - заявки на постачання, волонтерські заявки, пальне, GPS-трекінг.

Запуск:
  cd scripts/seed_db
  python3 -m venv .venv && source .venv/bin/activate
  pip install -r requirements.txt
  python seed.py                 # заповнити (не чистить існуюче, але дедуплікує по email/plate/name)
  python seed.py --reset         # TRUNCATE усіх таблиць, потім заповнити наново
  python seed.py --dsn 'postgres://postgres:postgres@localhost:5432/omnilog'

За замовчуванням DSN збігається з docker-compose.yml (порт 5432 проброшений на хост).

Паролі всіх користувачів: password123
"""

from __future__ import annotations

import argparse
import os
import random
import secrets
import string
import sys
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from typing import Iterable

import bcrypt
import psycopg
from psycopg.rows import dict_row

DEFAULT_DSN = os.environ.get(
    "SEED_DSN",
    "postgres://postgres:postgres@localhost:5432/omnilog?sslmode=disable",
)

DEFAULT_PASSWORD = "password123"

# ---------------------------------------------------------------------------
# Утиліти
# ---------------------------------------------------------------------------

def hash_password(pw: str) -> str:
    return bcrypt.hashpw(pw.encode(), bcrypt.gensalt()).decode()


PASSWORD_HASH = hash_password(DEFAULT_PASSWORD)


def rnd_plate() -> str:
    letters = "ABEKMHOPCTX"
    return (
        "AA"
        + "".join(random.choices(string.digits, k=4))
        + random.choice(letters)
        + random.choice(letters)
    )


def rnd_serial(prefix: str = "SN") -> str:
    return f"{prefix}-{secrets.token_hex(4).upper()}"


# ---------------------------------------------------------------------------
# Таблиці для TRUNCATE (правильний порядок — spread через CASCADE однаково буде)
# ---------------------------------------------------------------------------

TRUNCATE_TABLES = [
    "gps_locations",
    "geofence_alerts",
    "geofences",
    "fuel_records",
    "maintenance_records",
    "vehicle_driver_history",
    "inventory_check_items",
    "inventory_checks",
    "shipment_items",
    "shipments",
    "resource_assignments",
    "supply_requests",
    "contractor_requests",
    "audit_logs",
    "resources",
    "resource_categories",
    "warehouses",
    "refresh_tokens",
    "invite_tokens",
    "users",
    "units",
]


def reset_database(cur) -> None:
    existing = {
        row["table_name"]
        for row in cur.execute(
            "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
        ).fetchall()
    }
    tables = [t for t in TRUNCATE_TABLES if t in existing]
    if not tables:
        return
    cur.execute(f"TRUNCATE TABLE {', '.join(tables)} RESTART IDENTITY CASCADE")
    print(f"  • TRUNCATE {len(tables)} tables")


# ---------------------------------------------------------------------------
# Домен
# ---------------------------------------------------------------------------

@dataclass
class UserSpec:
    email: str
    full_name: str
    role: str
    unit_key: str | None  # ключ у словнику unit_ids; None для адмінів / контракторів
    status: str = "ACTIVE"


def ensure_unit(cur, parent_id: int | None, name: str, unit_type: str, tier: str) -> int:
    row = cur.execute(
        "SELECT id FROM units WHERE name = %s AND unit_type = %s", (name, unit_type)
    ).fetchone()
    if row:
        cur.execute(
            "UPDATE units SET subscription_tier = %s, parent_id = %s WHERE id = %s",
            (tier, parent_id, row["id"]),
        )
        return row["id"]
    row = cur.execute(
        """INSERT INTO units (parent_id, name, unit_type, subscription_tier)
           VALUES (%s, %s, %s, %s) RETURNING id""",
        (parent_id, name, unit_type, tier),
    ).fetchone()
    return row["id"]


def ensure_user(cur, spec: UserSpec, unit_ids: dict[str, int]) -> str:
    username = spec.email.split("@")[0]
    unit_id = unit_ids.get(spec.unit_key) if spec.unit_key else None
    row = cur.execute("SELECT id FROM users WHERE email = %s", (spec.email,)).fetchone()
    if row:
        cur.execute(
            """UPDATE users SET full_name=%s, role=%s, status=%s, unit_id=%s,
                                password_hash=%s, username=%s, updated_at=NOW()
               WHERE id=%s""",
            (spec.full_name, spec.role, spec.status, unit_id, PASSWORD_HASH, username, row["id"]),
        )
        return row["id"]
    row = cur.execute(
        """INSERT INTO users (username, email, full_name, password_hash, role, status, unit_id)
           VALUES (%s, %s, %s, %s, %s, %s, %s) RETURNING id""",
        (username, spec.email, spec.full_name, PASSWORD_HASH, spec.role, spec.status, unit_id),
    ).fetchone()
    return row["id"]


def ensure_category(cur, name: str, description: str) -> str:
    row = cur.execute(
        "SELECT id FROM resource_categories WHERE name = %s", (name,)
    ).fetchone()
    if row:
        return row["id"]
    row = cur.execute(
        "INSERT INTO resource_categories (name, description) VALUES (%s, %s) RETURNING id",
        (name, description),
    ).fetchone()
    return row["id"]


def ensure_warehouse(
    cur, unit_id: int, name: str, lat: float, lon: float, location_type: str = "STATIONARY"
) -> str:
    row = cur.execute(
        "SELECT id FROM warehouses WHERE name = %s AND unit_id = %s", (name, unit_id)
    ).fetchone()
    if row:
        return row["id"]
    row = cur.execute(
        """INSERT INTO warehouses (unit_id, name, location_type, latitude, longitude)
           VALUES (%s, %s, %s, %s, %s) RETURNING id""",
        (unit_id, name, location_type, lat, lon),
    ).fetchone()
    return row["id"]


def insert_resource(
    cur,
    *,
    category_id: str,
    unit_id: int,
    warehouse_id: str,
    name: str,
    quantity: int,
    min_quantity: int = 5,
    unit_type: str = "PCS",
    weight_kg: float = 1.0,
    condition: str = "NEW",
    description: str = "",
) -> str:
    # Перевірка по (name + warehouse_id) аби не дублювати
    row = cur.execute(
        "SELECT id FROM resources WHERE name = %s AND warehouse_id = %s",
        (name, warehouse_id),
    ).fetchone()
    if row:
        cur.execute(
            "UPDATE resources SET quantity=%s, updated_at=NOW() WHERE id=%s",
            (quantity, row["id"]),
        )
        return row["id"]
    row = cur.execute(
        """INSERT INTO resources
               (category_id, unit_id, warehouse_id, name, description, quantity,
                unit_type, serial_number, condition, min_quantity, weight_kg)
           VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) RETURNING id""",
        (
            category_id, unit_id, warehouse_id, name, description, quantity,
            unit_type, rnd_serial(), condition, min_quantity, weight_kg,
        ),
    ).fetchone()
    return row["id"]


def ensure_vehicle(
    cur,
    *,
    plate: str,
    brand: str,
    model: str,
    vtype: str,
    capacity_kg: float,
    tank: float,
    norm: float,
    driver_id: str | None,
) -> str:
    row = cur.execute("SELECT id FROM vehicles WHERE plate_number = %s", (plate,)).fetchone()
    if row:
        return row["id"]
    row = cur.execute(
        """INSERT INTO vehicles
               (brand, model, plate_number, type, capacity_kg, status, driver_id,
                tank_capacity, fuel_norm)
           VALUES (%s, %s, %s, %s, %s, 'ACTIVE', %s, %s, %s) RETURNING id""",
        (brand, model, plate, vtype, capacity_kg, driver_id, tank, norm),
    ).fetchone()
    return row["id"]


# ---------------------------------------------------------------------------
# Дані для сіду
# ---------------------------------------------------------------------------

REGIONS: list[dict] = [
    {
        "key": "west",
        "name": "Регіон «Захід»",
        "tier": "PRO",
        "branches": [
            {"key": "west_lviv", "name": "Філія Львів"},
            {"key": "west_ivf", "name": "Філія Івано-Франківськ"},
        ],
        "dept": {"key": "west_lviv_dept", "name": "Відділ логістики Львів", "parent_key": "west_lviv"},
        "team": {"key": "west_lviv_team", "name": "Команда кур'єрів Львів", "parent_key": "west_lviv_dept"},
        "warehouses": [
            ("west", "Центральний склад Захід", 49.8397, 24.0297),
            ("west_lviv", "Склад Львів-1", 49.8420, 24.0315),
            ("west_lviv", "Склад Львів-2 Зимна Вода", 49.8100, 23.9460),
            ("west_ivf", "Склад ІФ", 48.9215, 24.7106),
        ],
        "vehicles": 6,          # добряче — щоб протестувати PRO GPS/maintenance/fuel
        "resources_per_wh": 18,
    },
    {
        "key": "center",
        "name": "Регіон «Центр»",
        "tier": "BASIC",
        "branches": [
            {"key": "center_kyiv", "name": "Філія Київ"},
        ],
        "dept": {"key": "center_kyiv_dept", "name": "Відділ Київ", "parent_key": "center_kyiv"},
        "team": None,
        "warehouses": [
            ("center", "Центральний склад Центр", 50.4501, 30.5234),
            ("center_kyiv", "Склад Київ-Поділ", 50.4660, 30.5200),
        ],
        "vehicles": 3,
        "resources_per_wh": 10,
    },
    {
        "key": "east",
        "name": "Регіон «Схід»",
        "tier": "BASIC",
        "branches": [
            {"key": "east_kh", "name": "Філія Харків"},
        ],
        "dept": {"key": "east_kh_dept", "name": "Відділ Харків", "parent_key": "east_kh"},
        "team": None,
        # BASIC: MaxWarehouses=10 — створимо 9, щоби показати наближення до ліміту
        "warehouses": [
            ("east", "Склад Схід-1", 49.9935, 36.2304),
            ("east", "Склад Схід-2", 49.9800, 36.2500),
            ("east", "Склад Схід-3", 49.9700, 36.2600),
            ("east_kh", "Склад Харків-А", 49.9500, 36.3000),
            ("east_kh", "Склад Харків-Б", 49.9400, 36.3100),
            ("east_kh", "Склад Харків-В", 49.9300, 36.3200),
            ("east_kh", "Склад Харків-Г", 49.9200, 36.3300),
            ("east_kh", "Склад Харків-Д", 49.9100, 36.3400),
            ("east_kh", "Склад Харків-Е", 49.9000, 36.3500),
        ],
        "vehicles": 4,
        "resources_per_wh": 8,
    },
    {
        "key": "test",
        "name": "Регіон «Тест-ENTERPRISE»",
        "tier": "ENTERPRISE",
        "branches": [],
        "dept": None,
        "team": None,
        "warehouses": [
            ("test", "Демо-склад Enterprise", 48.4647, 35.0462),
        ],
        "vehicles": 2,
        "resources_per_wh": 5,
    },
]


CATEGORIES = [
    ("Канцелярія", "Ручки, папір, тощо"),
    ("Електроніка", "Ноутбуки, телефони"),
    ("Інструмент", "Ручний/електричний інструмент"),
    ("Медикаменти", "Аптечки та ліки"),
    ("Продукти", "Довготривале зберігання"),
    ("Одяг", "Форма та спецодяг"),
    ("Паливо-мастильні", "Оливи, фільтри, присадки"),
]

RESOURCE_NAMES = [
    "Ноутбук Lenovo ThinkPad", "Папка архівна", "Комплект аптечок IFAK",
    "Набір викруток 40 шт.", "Генератор 5 кВт", "Кабель-подовжувач 50м",
    "Шолом захисний", "Рукавиці робочі", "Ліхтарик тактичний",
    "Радіостанція Motorola", "Термобілизна комплект", "Дрон DJI Mavic",
    "Батарея LiFePO4 100Ah", "Рюкзак 60л", "Спальний мішок -10°C",
    "Фільтр для води", "Тент 6x8м", "Мультитул Leatherman",
    "Оливa моторна 5W40 4л", "Фільтр повітряний",
]


# ---------------------------------------------------------------------------
# Основний сідер
# ---------------------------------------------------------------------------

def seed(conn: psycopg.Connection, reset: bool) -> None:
    random.seed(42)
    with conn.cursor(row_factory=dict_row) as cur:
        if reset:
            print("→ Чищу БД…")
            reset_database(cur)

        print("→ Створюю одиниці (units)…")
        unit_ids: dict[str, int] = {}
        for region in REGIONS:
            rid = ensure_unit(cur, None, region["name"], "REGION", region["tier"])
            unit_ids[region["key"]] = rid
            for br in region["branches"]:
                bid = ensure_unit(cur, rid, br["name"], "BRANCH", region["tier"])
                unit_ids[br["key"]] = bid
            if region.get("dept"):
                d = region["dept"]
                did = ensure_unit(cur, unit_ids[d["parent_key"]], d["name"], "DEPARTMENT", region["tier"])
                unit_ids[d["key"]] = did
            if region.get("team"):
                t = region["team"]
                tid = ensure_unit(cur, unit_ids[t["parent_key"]], t["name"], "TEAM", region["tier"])
                unit_ids[t["key"]] = tid
        print(f"   створено/оновлено {len(unit_ids)} unit(s)")

        print("→ Створюю користувачів…")
        user_specs: list[UserSpec] = [
            # Super-admin без unit_id
            UserSpec("admin@Omnilog.local", "Головний Адмін", "ADMIN", None),
            UserSpec("admin2@Omnilog.local", "Резервний Адмін", "ADMIN", None),

            # --- Захід (PRO) ---
            UserSpec("director.west@Omnilog.local", "Оксана Директорова (PRO)", "REGION_DIRECTOR", "west"),
            UserSpec("logist.west@Omnilog.local", "Петро Логіст (PRO)", "REGION_LOGISTICIAN", "west"),
            UserSpec("storekeeper.west@Omnilog.local", "Марта Комірникова (PRO)", "REGION_STOREKEEPER", "west"),
            UserSpec("manager.lviv@Omnilog.local", "Богдан Менеджерів", "BRANCH_MANAGER", "west_lviv"),
            UserSpec("logist.lviv@Omnilog.local", "Ірина Логіст Львів", "BRANCH_LOGISTICIAN", "west_lviv"),
            UserSpec("storekeeper.lviv@Omnilog.local", "Олег Комірник Львів", "BRANCH_STOREKEEPER", "west_lviv"),
            UserSpec("manager.ivf@Omnilog.local", "Ярослав Менеджер ІФ", "BRANCH_MANAGER", "west_ivf"),
            UserSpec("dept.lviv@Omnilog.local", "Роман Начвідділу Львів", "DEPT_MANAGER", "west_lviv_dept"),
            UserSpec("supervisor.lviv@Omnilog.local", "Анна Супервайзер Львів", "DEPT_SUPERVISOR", "west_lviv_dept"),
            UserSpec("team.lviv@Omnilog.local", "Василь Тімлід Львів", "TEAM_LEAD", "west_lviv_team"),
            UserSpec("employee1.lviv@Omnilog.local", "Михайло Працівник", "EMPLOYEE", "west_lviv_team"),
            UserSpec("employee2.lviv@Omnilog.local", "Тарас Працівник", "EMPLOYEE", "west_lviv_team"),

            # --- Центр (BASIC) ---
            UserSpec("director.center@Omnilog.local", "Сергій Директор (BASIC)", "REGION_DIRECTOR", "center"),
            UserSpec("logist.center@Omnilog.local", "Юлія Логіст Центр", "REGION_LOGISTICIAN", "center"),
            UserSpec("storekeeper.center@Omnilog.local", "Дмитро Комірник Центр", "REGION_STOREKEEPER", "center"),
            UserSpec("manager.kyiv@Omnilog.local", "Наталя Менеджер Київ", "BRANCH_MANAGER", "center_kyiv"),
            UserSpec("logist.kyiv@Omnilog.local", "Олексій Логіст Київ", "BRANCH_LOGISTICIAN", "center_kyiv"),
            UserSpec("storekeeper.kyiv@Omnilog.local", "Світлана Комірник Київ", "BRANCH_STOREKEEPER", "center_kyiv"),
            UserSpec("dept.kyiv@Omnilog.local", "Віктор Начвідділу Київ", "DEPT_MANAGER", "center_kyiv_dept"),

            # --- Схід (BASIC, наближення до ліміту) ---
            UserSpec("director.east@Omnilog.local", "Ніна Директор (BASIC)", "REGION_DIRECTOR", "east"),
            UserSpec("logist.east@Omnilog.local", "Андрій Логіст Схід", "REGION_LOGISTICIAN", "east"),
            UserSpec("storekeeper.east@Omnilog.local", "Євген Комірник Схід", "REGION_STOREKEEPER", "east"),
            UserSpec("manager.kh@Omnilog.local", "Катерина Менеджер Харків", "BRANCH_MANAGER", "east_kh"),

            # --- Test ENTERPRISE ---
            UserSpec("director.test@Omnilog.local", "Enterprise Director", "REGION_DIRECTOR", "test"),
            UserSpec("logist.test@Omnilog.local", "Enterprise Logist", "REGION_LOGISTICIAN", "test"),

            # --- Контрактори (волонтери), без unit_id ---
            UserSpec("contractor1@Omnilog.local", "Волонтер Богдан", "CONTRACTOR", None),
            UserSpec("contractor2@Omnilog.local", "Волонтер Іванна", "CONTRACTOR", None),
            UserSpec("contractor3@Omnilog.local", "Волонтер Руслан", "CONTRACTOR", None),

            # Один BLOCKED для тестів
            UserSpec("blocked@Omnilog.local", "Заблокований", "EMPLOYEE", "west_lviv_team", status="BLOCKED"),
            # Один PENDING
            UserSpec("pending@Omnilog.local", "Новий Співробітник", "EMPLOYEE", "center_kyiv_dept", status="PENDING"),
        ]
        user_ids: dict[str, str] = {}
        for spec in user_specs:
            user_ids[spec.email] = ensure_user(cur, spec, unit_ids)
        print(f"   створено/оновлено {len(user_ids)} користувач(ів)")

        print("→ Створюю категорії ресурсів…")
        cat_ids: dict[str, str] = {name: ensure_category(cur, name, desc) for name, desc in CATEGORIES}

        print("→ Створюю склади та ресурси…")
        warehouses_created: dict[str, str] = {}   # "region_key|name" -> wh_id
        resources_by_region: dict[str, list[str]] = {r["key"]: [] for r in REGIONS}
        for region in REGIONS:
            reg_key = region["key"]
            for (unit_key, wh_name, lat, lon) in region["warehouses"]:
                wh_id = ensure_warehouse(cur, unit_ids[unit_key], wh_name, lat, lon)
                warehouses_created[f"{reg_key}|{wh_name}"] = wh_id
                # ресурси
                for i in range(region["resources_per_wh"]):
                    name = random.choice(RESOURCE_NAMES) + f" #{random.randint(100, 999)}"
                    cat_id = random.choice(list(cat_ids.values()))
                    qty = random.randint(1, 250)
                    min_q = random.choice([2, 5, 10, 20])
                    weight = round(random.uniform(0.1, 20.0), 2)
                    rid = insert_resource(
                        cur,
                        category_id=cat_id,
                        unit_id=unit_ids[unit_key],
                        warehouse_id=wh_id,
                        name=name,
                        quantity=qty,
                        min_quantity=min_q,
                        weight_kg=weight,
                    )
                    resources_by_region[reg_key].append(rid)
        print(f"   складів: {len(warehouses_created)}, ресурсів: {sum(len(v) for v in resources_by_region.values())}")

        print("→ Створюю авто (vehicles) та історію палива…")
        # Водіями роблю EMPLOYEE/TEAM_LEAD/BRANCH_*
        drivers_by_region = {
            "west": [user_ids["employee1.lviv@Omnilog.local"], user_ids["employee2.lviv@Omnilog.local"],
                    user_ids["team.lviv@Omnilog.local"], user_ids["logist.lviv@Omnilog.local"]],
            "center": [user_ids["logist.kyiv@Omnilog.local"], user_ids["manager.kyiv@Omnilog.local"]],
            "east": [user_ids["logist.east@Omnilog.local"], user_ids["manager.kh@Omnilog.local"]],
            "test": [user_ids["logist.test@Omnilog.local"]],
        }
        vehicle_ids_by_region: dict[str, list[str]] = {}
        veh_types = [("Renault", "Master", "VAN", 1500, 80, 9.0),
                     ("Ford", "Transit", "VAN", 1800, 75, 10.0),
                     ("Mercedes", "Sprinter", "TRUCK", 3500, 90, 11.5),
                     ("Toyota", "Hilux", "PICKUP", 1000, 80, 9.5),
                     ("Volkswagen", "Crafter", "VAN", 2500, 75, 10.5),
                     ("MAN", "TGE", "TRUCK", 3500, 100, 12.0)]
        for region in REGIONS:
            ids = []
            for i in range(region["vehicles"]):
                brand, model, vtype, cap, tank, norm = veh_types[i % len(veh_types)]
                driver_pool = drivers_by_region.get(region["key"], [])
                driver_id = random.choice(driver_pool) if driver_pool else None
                vid = ensure_vehicle(
                    cur, plate=rnd_plate(), brand=brand, model=f"{model} {random.randint(2018, 2024)}",
                    vtype=vtype, capacity_kg=cap, tank=tank, norm=norm, driver_id=driver_id,
                )
                ids.append(vid)
            vehicle_ids_by_region[region["key"]] = ids

        # Записи про пальне (fuel_records) + одна аномалія
        fuel_inserted = 0
        for reg_key, veh_ids in vehicle_ids_by_region.items():
            for vid in veh_ids:
                odometer = random.randint(30_000, 120_000)
                for day in range(0, 30, random.choice([2, 3, 4])):
                    liters = round(random.uniform(25, 70), 2)
                    odometer += random.randint(150, 400)
                    created_at = datetime.now(timezone.utc) - timedelta(days=30 - day)
                    cur.execute(
                        """INSERT INTO fuel_records (vehicle_id, liters, odometer_km, record_type, created_at)
                           VALUES (%s, %s, %s, 'REFUEL', %s)""",
                        (vid, liters, odometer, created_at),
                    )
                    fuel_inserted += 1
                # штучна аномалія
                cur.execute(
                    """INSERT INTO fuel_records (vehicle_id, liters, odometer_km, record_type,
                                                 created_at, is_anomaly, anomaly_reason)
                       VALUES (%s, %s, %s, 'REFUEL', NOW() - INTERVAL '1 day', TRUE,
                               'Аномальна витрата: у 2.3x перевищує норму')""",
                    (vid, 180.5, odometer + 500),
                )
                fuel_inserted += 1
        print(f"   авто: {sum(len(v) for v in vehicle_ids_by_region.values())}, fuel_records: {fuel_inserted}")

        print("→ Додаю ТО (maintenance_records)…")
        for veh_ids in vehicle_ids_by_region.values():
            for vid in veh_ids:
                if random.random() < 0.5:
                    cur.execute(
                        """INSERT INTO maintenance_records
                               (vehicle_id, odometer_km, description, performed_by, cost_amount, created_at)
                           VALUES (%s, %s, %s, %s, %s, NOW() - INTERVAL '%s days')""",
                        (vid, random.randint(40000, 110000),
                         random.choice(["Заміна оливи", "Ремонт підвіски", "Заміна колодок", "Діагностика"]),
                         "СТО Центральне", round(random.uniform(1500, 12000), 2),
                         random.randint(5, 90)),
                    )

        print("→ Створюю заявки на постачання (supply_requests)…")
        # Беремо ресурси з Заходу/Центру/Сходу і робимо заявки в різних статусах
        statuses = ["PENDING", "PENDING", "APPROVED", "REJECTED", "COMPLETED"]
        creators = [
            ("director.west@Omnilog.local", "west"),
            ("logist.lviv@Omnilog.local", "west"),
            ("manager.kyiv@Omnilog.local", "center"),
            ("logist.east@Omnilog.local", "east"),
            ("dept.lviv@Omnilog.local", "west"),
        ]
        for creator_email, reg_key in creators:
            for _ in range(4):
                if not resources_by_region[reg_key]:
                    continue
                res_id = random.choice(resources_by_region[reg_key])
                status = random.choice(statuses)
                approved_by = None
                approved_at = None
                if status in ("APPROVED", "REJECTED", "COMPLETED"):
                    approved_by = user_ids[f"director.{reg_key}@Omnilog.local"]
                    approved_at = datetime.now(timezone.utc) - timedelta(days=random.randint(1, 15))
                cur.execute(
                    """INSERT INTO supply_requests
                           (created_by, resource_id, quantity, status, approved_by, approved_at,
                            comment, created_at)
                       VALUES (%s, %s, %s, %s, %s, %s, %s,
                               NOW() - INTERVAL '%s days')""",
                    (user_ids[creator_email], res_id, random.randint(1, 20), status,
                     approved_by, approved_at,
                     random.choice(["Терміново", "Планова поставка", "Резерв", ""]),
                     random.randint(1, 30)),
                )

        print("→ Створюю волонтерські заявки (contractor_requests)…")
        vol_creators = [("director.west@Omnilog.local", "west"),
                        ("manager.kyiv@Omnilog.local", "center"),
                        ("logist.east@Omnilog.local", "east")]
        contractor_users = [user_ids["contractor1@Omnilog.local"],
                            user_ids["contractor2@Omnilog.local"],
                            user_ids["contractor3@Omnilog.local"]]
        for creator_email, reg_key in vol_creators:
            for _ in range(3):
                status = random.choice(["OPEN", "OPEN", "IN_PROGRESS", "DELIVERED", "COMPLETED"])
                taken_by = None
                taken_at = None
                completed_at = None
                if status in ("IN_PROGRESS", "DELIVERED", "COMPLETED"):
                    taken_by = random.choice(contractor_users)
                    taken_at = datetime.now(timezone.utc) - timedelta(days=random.randint(1, 10))
                if status == "COMPLETED":
                    completed_at = taken_at + timedelta(days=random.randint(1, 3))
                cur.execute(
                    """INSERT INTO contractor_requests
                           (created_by, unit_id, title, description, status, taken_by, taken_at,
                            completed_at, created_at)
                       VALUES (%s, %s, %s, %s, %s, %s, %s, %s, NOW() - INTERVAL '%s days')""",
                    (user_ids[creator_email], unit_ids[reg_key],
                     random.choice(["Потрібен генератор 5 кВт", "Потрібні аптечки IFAK",
                                    "Потрібна доставка продуктів", "Ноутбуки (2 шт)"]),
                     "Тестова заявка, згенерована сідером.", status,
                     taken_by, taken_at, completed_at, random.randint(1, 20)),
                )

        print("→ Додаю GPS-трекінг та геозони (PRO фіча — тільки для Заходу і Тест-ENT)…")
        # Перевіряємо, чи таблиці GPS взагалі існують (міграції backend-у могли ще не запускатись).
        existing_tables = {
            r["table_name"]
            for r in cur.execute(
                "SELECT table_name FROM information_schema.tables WHERE table_schema='public'"
            ).fetchall()
        }
        has_gps = {"geofences", "gps_locations"}.issubset(existing_tables)
        if not has_gps:
            print("   ⚠ Таблиці gps_locations/geofences ще не створені — пропускаю.")
            print("     (Підніми бекенд `docker compose up backend` — він зробить міграції — і перезапусти сідер.)")
        # geofences
        if has_gps:
            for reg_key in ("west", "test"):
                for (name, lat, lon, radius, gtype) in [
                    ("Безпечна зона — штаб", 49.8397, 24.0297, 500, "SAFE"),
                    ("Забороненa зона — полігон", 49.9000, 24.1000, 1000, "FORBIDDEN"),
                ]:
                    cur.execute(
                        """INSERT INTO geofences (unit_id, name, latitude, longitude, radius, type, active)
                           VALUES (%s, %s, %s, %s, %s, %s, TRUE)""",
                        (unit_ids[reg_key], name, lat, lon, radius, gtype),
                    )
        # gps_locations для авто Заходу — треки по 20 точок.
        # ПІСЛЯ ФІКСУ МІГРАЦІЇ: vehicle_id — UUID, тож передаємо UUID напряму.
        gps_ok = 0
        if has_gps:
            for reg_key in ("west", "test"):
                for vid in vehicle_ids_by_region[reg_key]:
                    lat0, lon0 = 49.84, 24.03
                    for i in range(20):
                        lat = lat0 + random.uniform(-0.05, 0.05)
                        lon = lon0 + random.uniform(-0.05, 0.05)
                        ts = datetime.now(timezone.utc) - timedelta(minutes=20 - i)
                        try:
                            cur.execute(
                                """INSERT INTO gps_locations
                                       (vehicle_id, unit_id, latitude, longitude, speed, heading,
                                        accuracy, timestamp)
                                   VALUES (%s, %s, %s, %s, %s, %s, %s, %s)""",
                                (vid, unit_ids[reg_key], lat, lon,
                                 round(random.uniform(0, 90), 2), round(random.uniform(0, 360), 2),
                                 round(random.uniform(2, 15), 2), ts),
                            )
                            gps_ok += 1
                        except psycopg.errors.Error:
                            conn.rollback()
                            break
        print(f"   gps_locations: {gps_ok}")

        conn.commit()
        print("✅ Готово!")
        print_summary(cur, unit_ids)


def print_summary(cur, unit_ids: dict[str, int]) -> None:
    print("\n================== ПІДСУМОК ==================")
    for reg_key, uid in unit_ids.items():
        if "_" in reg_key:
            continue
        row = cur.execute("SELECT name, subscription_tier FROM units WHERE id = %s", (uid,)).fetchone()
        print(f"  {row['name']:30s}  tier={row['subscription_tier']:10s}  id={uid}")
    print("\n-- Акаунти для логіну (пароль: password123) --")
    print("  ADMIN:              admin@Omnilog.local")
    print("  PRO director:       director.west@Omnilog.local")
    print("  PRO logist:         logist.west@Omnilog.local")
    print("  BASIC director:     director.center@Omnilog.local  (буде отримувати 402 Payment Required)")
    print("  BASIC (near limit): director.east@Omnilog.local    (9/10 складів на ліміті BASIC)")
    print("  ENTERPRISE:         director.test@Omnilog.local")
    print("  CONTRACTOR:         contractor1@Omnilog.local")
    print("==============================================\n")


def main() -> int:
    parser = argparse.ArgumentParser(description="Seed Omnilog database")
    parser.add_argument("--dsn", default=DEFAULT_DSN, help="PostgreSQL DSN")
    parser.add_argument("--reset", action="store_true",
                        help="TRUNCATE усі таблиці перед сідом")
    args = parser.parse_args()

    print(f"→ Підключаюсь до {args.dsn}")
    try:
        with psycopg.connect(args.dsn, autocommit=False) as conn:
            seed(conn, reset=args.reset)
    except psycopg.OperationalError as exc:
        print(f"❌ Не вдалось під'єднатись: {exc}", file=sys.stderr)
        print("   Переконайся, що PostgreSQL запущений:  docker compose up -d postgres", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
