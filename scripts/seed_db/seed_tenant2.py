#!/usr/bin/env python3
"""
Seed-скрипт: створює другий тенант «ТОВ Карго-Логістика» та заповнює БД
тестовими даними, що відповідають реальній схемі (docker/postgres/init.sql).

Схема враховує tenant_id у всіх таблицях:
  tenants, users, units, resource_categories, resources,
  vehicles, fuel_records, supply_requests, contractor_requests.

Опціональні таблиці (якщо вже є через міграції бекенду):
  warehouses, maintenance_records, gps_locations, geofences,
  notifications, audit_logs, shipments, resource_assignments.

Запуск:
  cd scripts/seed_db
  python3 -m venv .venv && source .venv/bin/activate
  pip install -r requirements.txt
  python seed_tenant2.py
  python seed_tenant2.py --reset          # видалити дані тільки цього тенанта, потім залити
  python seed_tenant2.py --dsn 'postgres://postgres:postgres@localhost:5432/omnilog'

Пароль усіх користувачів: password123
"""

from __future__ import annotations

import argparse
import os
import random
import secrets
import string
import sys
from dataclasses import dataclass, field
from datetime import datetime, timedelta, timezone
from typing import Optional

import bcrypt
import psycopg
from psycopg.rows import dict_row

# ---------------------------------------------------------------------------
# Конфіг
# ---------------------------------------------------------------------------

DEFAULT_DSN = os.environ.get(
    "SEED_DSN",
    "postgres://postgres:postgres@localhost:5432/omnilog?sslmode=disable",
)
DEFAULT_PASSWORD = "password123"

TENANT_NAME = "ТОВ «Карго-Логістика»"
TENANT_SLUG = "cargo-logistics"
TENANT_TIER = "PRO"           # BASIC | PRO | ENTERPRISE
TENANT_OWNER_EMAIL = "owner@cargo-logistics.local"

# ---------------------------------------------------------------------------
# Утиліти
# ---------------------------------------------------------------------------

def hash_password(pw: str) -> str:
    return bcrypt.hashpw(pw.encode(), bcrypt.gensalt()).decode()


PASSWORD_HASH = hash_password(DEFAULT_PASSWORD)


def rnd_plate() -> str:
    """Генерує псевдо-номерний знак типу AA1234BB."""
    letters = "ABEKMHOPCTX"
    return (
        random.choice(letters)
        + random.choice(letters)
        + "".join(random.choices(string.digits, k=4))
        + random.choice(letters)
        + random.choice(letters)
    )


def rnd_serial(prefix: str = "SN") -> str:
    return f"{prefix}-{secrets.token_hex(4).upper()}"


def tables_in_db(cur) -> set[str]:
    return {
        r["table_name"]
        for r in cur.execute(
            "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
        ).fetchall()
    }


# ---------------------------------------------------------------------------
# Ensure-функції (upsert або insert за наявності)
# ---------------------------------------------------------------------------

def ensure_tenant(cur) -> str:
    row = cur.execute("SELECT id FROM tenants WHERE slug = %s", (TENANT_SLUG,)).fetchone()
    if row:
        print(f"   ℹ Тенант '{TENANT_SLUG}' вже існує — оновлюю.")
        cur.execute(
            """UPDATE tenants
               SET name=%s, subscription_tier=%s, owner_email=%s, is_active=TRUE, updated_at=NOW()
               WHERE slug=%s""",
            (TENANT_NAME, TENANT_TIER, TENANT_OWNER_EMAIL, TENANT_SLUG),
        )
        return row["id"]
    row = cur.execute(
        """INSERT INTO tenants (name, slug, subscription_tier, owner_email, is_active)
           VALUES (%s, %s, %s, %s, TRUE) RETURNING id""",
        (TENANT_NAME, TENANT_SLUG, TENANT_TIER, TENANT_OWNER_EMAIL),
    ).fetchone()
    return row["id"]


def delete_tenant_data(cur, tenant_id: str) -> None:
    """Видаляє всі дані тенанта в безпечному порядку (CASCADE через FK)."""
    # Достатньо видалити самого тенанта — CASCADE прибере все.
    cur.execute("DELETE FROM tenants WHERE id = %s", (tenant_id,))
    print(f"   • Дані тенанта {tenant_id} видалено (CASCADE).")


def ensure_unit(
    cur,
    tenant_id: str,
    parent_id: Optional[int],
    name: str,
    unit_type: str,
) -> int:
    """
    units(tenant_id, parent_id, name, unit_type)
    Немає колонки subscription_tier — вона на рівні тенанта.
    """
    row = cur.execute(
        "SELECT id FROM units WHERE tenant_id = %s AND name = %s AND unit_type = %s",
        (tenant_id, name, unit_type),
    ).fetchone()
    if row:
        return row["id"]
    row = cur.execute(
        """INSERT INTO units (tenant_id, parent_id, name, unit_type)
           VALUES (%s, %s, %s, %s) RETURNING id""",
        (tenant_id, parent_id, name, unit_type),
    ).fetchone()
    return row["id"]


def ensure_user(
    cur,
    tenant_id: str,
    email: str,
    username: str,
    full_name: str,
    role: str,
    status: str,
    unit_id: Optional[int],
    phone: Optional[str] = None,
) -> str:
    row = cur.execute(
        "SELECT id FROM users WHERE tenant_id = %s AND email = %s",
        (tenant_id, email),
    ).fetchone()
    if row:
        cur.execute(
            """UPDATE users
               SET full_name=%s, role=%s, status=%s, unit_id=%s,
                   password_hash=%s, username=%s, updated_at=NOW()
               WHERE id=%s""",
            (full_name, role, status, unit_id, PASSWORD_HASH, username, row["id"]),
        )
        return row["id"]
    row = cur.execute(
        """INSERT INTO users
               (tenant_id, username, email, full_name, password_hash, role, status, unit_id, phone)
           VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s) RETURNING id""",
        (tenant_id, username, email, full_name, PASSWORD_HASH, role, status, unit_id, phone),
    ).fetchone()
    return row["id"]


def ensure_category(cur, tenant_id: str, name: str, description: str) -> str:
    row = cur.execute(
        "SELECT id FROM resource_categories WHERE tenant_id = %s AND name = %s",
        (tenant_id, name),
    ).fetchone()
    if row:
        return row["id"]
    row = cur.execute(
        """INSERT INTO resource_categories (tenant_id, name, description)
           VALUES (%s, %s, %s) RETURNING id""",
        (tenant_id, name, description),
    ).fetchone()
    return row["id"]


def ensure_resource(
    cur,
    tenant_id: str,
    category_id: str,
    unit_id: int,
    name: str,
    description: str,
    quantity: int,
    condition: str,
    min_quantity: int,
    serial_number: str,
    location: str,
) -> str:
    # Дедуплікуємо по (tenant_id, name, unit_id)
    row = cur.execute(
        "SELECT id FROM resources WHERE tenant_id=%s AND name=%s AND unit_id=%s",
        (tenant_id, name, unit_id),
    ).fetchone()
    if row:
        cur.execute(
            "UPDATE resources SET quantity=%s, updated_at=NOW() WHERE id=%s",
            (quantity, row["id"]),
        )
        return row["id"]
    row = cur.execute(
        """INSERT INTO resources
               (tenant_id, category_id, unit_id, name, description, quantity,
                serial_number, location, condition, min_quantity)
           VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s) RETURNING id""",
        (
            tenant_id, category_id, unit_id, name, description, quantity,
            serial_number, location, condition, min_quantity,
        ),
    ).fetchone()
    return row["id"]


def ensure_vehicle(
    cur,
    tenant_id: str,
    plate: str,
    brand: str,
    model: str,
    driver_id: Optional[str],
) -> str:
    row = cur.execute(
        "SELECT id FROM vehicles WHERE tenant_id=%s AND plate_number=%s",
        (tenant_id, plate),
    ).fetchone()
    if row:
        return row["id"]
    row = cur.execute(
        """INSERT INTO vehicles (tenant_id, brand, model, plate_number, status, driver_id)
           VALUES (%s, %s, %s, %s, 'ACTIVE', %s) RETURNING id""",
        (tenant_id, brand, model, plate, driver_id),
    ).fetchone()
    return row["id"]


# ---------------------------------------------------------------------------
# Тестові дані
# ---------------------------------------------------------------------------

CATEGORIES = [
    ("Канцелярія",       "Ручки, папір, офісне приладдя"),
    ("Електроніка",      "Ноутбуки, телефони, зарядки"),
    ("Інструмент",       "Ручний та електричний інструмент"),
    ("Медикаменти",      "Аптечки, ліки, перев'язувальні матеріали"),
    ("Продукти",         "Консерви та продукти довготривалого зберігання"),
    ("Одяг",             "Форма, спецодяг, ЗІЗ"),
    ("Паливно-мастильні","Оливи, фільтри, технічні рідини"),
    ("Запчастини",       "Авто-запчастини та комплектуючі"),
]

RESOURCE_NAMES = [
    "Ноутбук Lenovo ThinkPad E15",
    "Планшет Samsung Galaxy Tab",
    "Комплект аптечок IFAK",
    "Набір викруток 40 шт.",
    "Генератор бензиновий 5 кВт",
    "Кабель-подовжувач 50 м",
    "Захисна каска EN397",
    "Рукавиці робочі (пара)",
    "Ліхтарик тактичний Fenix",
    "Радіостанція Motorola DP2400",
    "Термобілизна комплект",
    "Дрон DJI Mini 3",
    "Акумулятор LiFePO4 100 Ah",
    "Рюкзак тактичний 60 л",
    "Спальний мішок -10°C",
    "Фільтр для води LifeStraw",
    "Тент армований 6×8 м",
    "Мультитул Leatherman Wave",
    "Олива моторна 5W-40 (4 л)",
    "Фільтр повітряний Mann-Filter",
    "Аптечка автомобільна",
    "Вогнегасник ВП-5",
    "Насос ручний для шин",
    "Трос буксирувальний 5 т",
    "Лопата складна армійська",
]

VEHICLE_TYPES = [
    ("Renault",    "Master 2.3 dCi",   "VAN"),
    ("Ford",       "Transit 2.0 TDCi", "VAN"),
    ("Mercedes",   "Sprinter 319 CDI", "TRUCK"),
    ("Toyota",     "Hilux 2.8 D-4D",   "PICKUP"),
    ("Volkswagen", "Crafter 35",       "VAN"),
    ("MAN",        "TGE 5.180",        "TRUCK"),
    ("Iveco",      "Daily 70C18",      "TRUCK"),
    ("Peugeot",    "Boxer 2.2 HDi",    "VAN"),
]

SUPPLY_STATUSES = ["PENDING", "PENDING", "APPROVED", "REJECTED", "COMPLETED"]
CONTRACTOR_STATUSES = ["OPEN", "OPEN", "IN_PROGRESS", "DELIVERED", "COMPLETED"]

CONTRACTOR_TITLES = [
    "Потрібен генератор 5 кВт",
    "Потрібні аптечки IFAK (10 комплектів)",
    "Доставка продуктів харчування",
    "Ноутбуки — 2 шт. для офісу",
    "Комплект рацій Motorola (5 шт.)",
    "Акумулятори LiFePO4 (2 шт.)",
    "Запчастини для Ford Transit",
    "Спальні мішки для команди (6 шт.)",
]

# ---------------------------------------------------------------------------
# Ієрархія підрозділів тенанта
# ---------------------------------------------------------------------------
# Регіони → Філії → Відділи → Команди

UNIT_HIERARCHY = [
    {
        "key":    "north",
        "name":   "Регіон «Північ»",
        "type":   "REGION",
        "parent": None,
        "children": [
            {
                "key":    "north_chernihiv",
                "name":   "Філія Чернігів",
                "type":   "BRANCH",
                "parent": "north",
                "children": [
                    {
                        "key":    "north_ch_dept",
                        "name":   "Відділ логістики Чернігів",
                        "type":   "DEPARTMENT",
                        "parent": "north_chernihiv",
                        "children": [
                            {
                                "key":    "north_ch_team",
                                "name":   "Команда доставки",
                                "type":   "TEAM",
                                "parent": "north_ch_dept",
                                "children": [],
                            }
                        ],
                    }
                ],
            },
            {
                "key":    "north_sumy",
                "name":   "Філія Суми",
                "type":   "BRANCH",
                "parent": "north",
                "children": [],
            },
        ],
    },
    {
        "key":    "south",
        "name":   "Регіон «Південь»",
        "type":   "REGION",
        "parent": None,
        "children": [
            {
                "key":    "south_odesa",
                "name":   "Філія Одеса",
                "type":   "BRANCH",
                "parent": "south",
                "children": [
                    {
                        "key":    "south_od_dept",
                        "name":   "Відділ митної логістики",
                        "type":   "DEPARTMENT",
                        "parent": "south_odesa",
                        "children": [],
                    }
                ],
            },
            {
                "key":    "south_mykolaiv",
                "name":   "Філія Миколаїв",
                "type":   "BRANCH",
                "parent": "south",
                "children": [],
            },
        ],
    },
    {
        "key":    "central",
        "name":   "Регіон «Центр»",
        "type":   "REGION",
        "parent": None,
        "children": [
            {
                "key":    "central_kyiv",
                "name":   "Філія Київ",
                "type":   "BRANCH",
                "parent": "central",
                "children": [
                    {
                        "key":    "central_kyiv_dept",
                        "name":   "Відділ складської логістики",
                        "type":   "DEPARTMENT",
                        "parent": "central_kyiv",
                        "children": [
                            {
                                "key":    "central_kyiv_team_a",
                                "name":   "Команда А (нічна зміна)",
                                "type":   "TEAM",
                                "parent": "central_kyiv_dept",
                                "children": [],
                            },
                            {
                                "key":    "central_kyiv_team_b",
                                "name":   "Команда Б (денна зміна)",
                                "type":   "TEAM",
                                "parent": "central_kyiv_dept",
                                "children": [],
                            },
                        ],
                    }
                ],
            },
        ],
    },
]


@dataclass
class UserSpec:
    email: str
    full_name: str
    role: str
    unit_key: Optional[str]
    status: str = "ACTIVE"
    phone: Optional[str] = None


USER_SPECS: list[UserSpec] = [
    # Tenant admin
    UserSpec("admin@cargo-logistics.local",        "Адміністратор Тенанта",        "TENANT_ADMIN",         None),
    UserSpec("owner@cargo-logistics.local",        "Власник Карго-Логістика",      "TENANT_ADMIN",         None),

    # Регіон Північ (PRO)
    UserSpec("dir.north@cargo-logistics.local",    "Олексій Дирєкторенко",         "REGION_DIRECTOR",      "north"),
    UserSpec("log.north@cargo-logistics.local",    "Ірина Логістенко",             "REGION_LOGISTICIAN",   "north"),
    UserSpec("store.north@cargo-logistics.local",  "Дмитро Комірниченко",          "REGION_STOREKEEPER",   "north"),

    UserSpec("mgr.chernihiv@cargo-logistics.local","Василь Менеджеренко",          "BRANCH_MANAGER",       "north_chernihiv"),
    UserSpec("log.chernihiv@cargo-logistics.local","Наталія Логіст Чернігів",      "BRANCH_LOGISTICIAN",   "north_chernihiv"),
    UserSpec("store.ch@cargo-logistics.local",     "Роман Комірник Чернігів",      "BRANCH_STOREKEEPER",   "north_chernihiv"),

    UserSpec("dept.ch@cargo-logistics.local",      "Андрій Начвідділу Чернігів",   "DEPT_MANAGER",         "north_ch_dept"),
    UserSpec("sup.ch@cargo-logistics.local",       "Юлія Супервайзер Чернігів",    "DEPT_SUPERVISOR",      "north_ch_dept"),
    UserSpec("lead.ch@cargo-logistics.local",      "Сергій Тімлід Чернігів",       "TEAM_LEAD",            "north_ch_team"),
    UserSpec("emp1.ch@cargo-logistics.local",      "Микола Працівник 1",           "EMPLOYEE",             "north_ch_team",
             phone="+380671234567"),
    UserSpec("emp2.ch@cargo-logistics.local",      "Оксана Працівниця 2",          "EMPLOYEE",             "north_ch_team",
             phone="+380671234568"),

    UserSpec("mgr.sumy@cargo-logistics.local",     "Катерина Менеджер Суми",       "BRANCH_MANAGER",       "north_sumy"),
    UserSpec("log.sumy@cargo-logistics.local",     "Євген Логіст Суми",            "BRANCH_LOGISTICIAN",   "north_sumy"),

    # Регіон Південь
    UserSpec("dir.south@cargo-logistics.local",    "Тетяна Директор Південь",      "REGION_DIRECTOR",      "south"),
    UserSpec("log.south@cargo-logistics.local",    "Богдан Логіст Південь",        "REGION_LOGISTICIAN",   "south"),
    UserSpec("store.south@cargo-logistics.local",  "Людмила Комірник Південь",     "REGION_STOREKEEPER",   "south"),

    UserSpec("mgr.odesa@cargo-logistics.local",    "Вікторія Менеджер Одеса",      "BRANCH_MANAGER",       "south_odesa"),
    UserSpec("log.odesa@cargo-logistics.local",    "Ярослав Логіст Одеса",         "BRANCH_LOGISTICIAN",   "south_odesa"),
    UserSpec("dept.od@cargo-logistics.local",      "Ганна Начвідділу Митниця",     "DEPT_MANAGER",         "south_od_dept"),

    UserSpec("mgr.mykolaiv@cargo-logistics.local", "Петро Менеджер Миколаїв",      "BRANCH_MANAGER",       "south_mykolaiv"),

    # Регіон Центр
    UserSpec("dir.central@cargo-logistics.local",  "Ростислав Директор Центр",     "REGION_DIRECTOR",      "central"),
    UserSpec("log.central@cargo-logistics.local",  "Марина Логіст Центр",          "REGION_LOGISTICIAN",   "central"),
    UserSpec("store.central@cargo-logistics.local","Іван Комірник Центр",          "REGION_STOREKEEPER",   "central"),

    UserSpec("mgr.kyiv@cargo-logistics.local",     "Олена Менеджер Київ",          "BRANCH_MANAGER",       "central_kyiv"),
    UserSpec("log.kyiv@cargo-logistics.local",     "Павло Логіст Київ",            "BRANCH_LOGISTICIAN",   "central_kyiv"),
    UserSpec("store.kyiv@cargo-logistics.local",   "Аліна Комірник Київ",          "BRANCH_STOREKEEPER",   "central_kyiv"),
    UserSpec("dept.kyiv@cargo-logistics.local",    "Максим Начвідділу Київ",       "DEPT_MANAGER",         "central_kyiv_dept"),
    UserSpec("sup.kyiv@cargo-logistics.local",     "Світлана Супервайзер Київ",    "DEPT_SUPERVISOR",      "central_kyiv_dept"),
    UserSpec("lead.kyiv.a@cargo-logistics.local",  "Тарас Тімлід A",               "TEAM_LEAD",            "central_kyiv_team_a"),
    UserSpec("lead.kyiv.b@cargo-logistics.local",  "Лариса Тімлід Б",              "TEAM_LEAD",            "central_kyiv_team_b"),
    UserSpec("emp3.kyiv@cargo-logistics.local",    "Денис Працівник 3",            "EMPLOYEE",             "central_kyiv_team_a"),
    UserSpec("emp4.kyiv@cargo-logistics.local",    "Вероніка Працівниця 4",        "EMPLOYEE",             "central_kyiv_team_b"),

    # Контрактори (без unit_id)
    UserSpec("contractor1@cargo-logistics.local",  "Волонтер Антон",               "CONTRACTOR",           None),
    UserSpec("contractor2@cargo-logistics.local",  "Волонтер Марія",               "CONTRACTOR",           None),

    # Тест-статуси
    UserSpec("blocked@cargo-logistics.local",      "Заблокований Тест",            "EMPLOYEE",             "north_ch_team",
             status="BLOCKED"),
    UserSpec("pending@cargo-logistics.local",      "Новий Тест (PENDING)",         "EMPLOYEE",             "central_kyiv_dept",
             status="PENDING"),
]


# ---------------------------------------------------------------------------
# Допоміжні рекурсивні функції для hierarchy
# ---------------------------------------------------------------------------

def flatten_units(nodes: list[dict]) -> list[dict]:
    """Повертає плоский список вузлів для ітерації."""
    result = []
    for node in nodes:
        result.append(node)
        result.extend(flatten_units(node["children"]))
    return result


def build_units(cur, tenant_id: str, nodes: list[dict], parent_map: dict[str, int]) -> dict[str, int]:
    unit_ids: dict[str, int] = {}
    for node in nodes:
        parent_id = parent_map.get(node["parent"]) if node["parent"] else None
        uid = ensure_unit(cur, tenant_id, parent_id, node["name"], node["type"])
        unit_ids[node["key"]] = uid
        # Рекурсивно діти
        child_ids = build_units(cur, tenant_id, node["children"], {**parent_map, **unit_ids})
        unit_ids.update(child_ids)
    return unit_ids


# ---------------------------------------------------------------------------
# Основний seed
# ---------------------------------------------------------------------------

def seed(conn: psycopg.Connection, reset: bool) -> None:
    random.seed(2025)
    with conn.cursor(row_factory=dict_row) as cur:

        existing = tables_in_db(cur)

        # ── Тенант ──────────────────────────────────────────────────────────
        if reset:
            print("→ Шукаю існуючий тенант для скидання…")
            row = cur.execute("SELECT id FROM tenants WHERE slug=%s", (TENANT_SLUG,)).fetchone()
            if row:
                delete_tenant_data(cur, row["id"])
                conn.commit()

        print("→ Створюю тенант…")
        tenant_id = ensure_tenant(cur)
        print(f"   tenant_id = {tenant_id}")

        # ── Підрозділи ───────────────────────────────────────────────────────
        print("→ Будую ієрархію підрозділів (units)…")
        unit_ids = build_units(cur, tenant_id, UNIT_HIERARCHY, {})
        print(f"   створено {len(unit_ids)} unit(s): {list(unit_ids.keys())}")

        # ── Користувачі ──────────────────────────────────────────────────────
        print("→ Створюю користувачів…")
        user_ids: dict[str, str] = {}
        for spec in USER_SPECS:
            uid_for_unit = unit_ids.get(spec.unit_key) if spec.unit_key else None
            username = spec.email.split("@")[0]
            user_ids[spec.email] = ensure_user(
                cur,
                tenant_id=tenant_id,
                email=spec.email,
                username=username,
                full_name=spec.full_name,
                role=spec.role,
                status=spec.status,
                unit_id=uid_for_unit,
                phone=spec.phone,
            )
        print(f"   створено/оновлено {len(user_ids)} користувач(ів)")

        # ── Категорії ────────────────────────────────────────────────────────
        print("→ Створюю категорії ресурсів…")
        cat_ids: dict[str, str] = {
            name: ensure_category(cur, tenant_id, name, desc)
            for name, desc in CATEGORIES
        }

        # ── Ресурси ──────────────────────────────────────────────────────────
        print("→ Заповнюю ресурси…")
        # Розподіляємо ресурси по регіонам
        region_keys_for_resources = [
            "north", "north_chernihiv", "north_sumy",
            "south", "south_odesa", "south_mykolaiv",
            "central", "central_kyiv",
        ]
        resource_ids_by_unit: dict[str, list[str]] = {k: [] for k in region_keys_for_resources}
        counts_per_unit = {
            "north": 15, "north_chernihiv": 20, "north_sumy": 12,
            "south": 10, "south_odesa": 18, "south_mykolaiv": 10,
            "central": 8,  "central_kyiv": 25,
        }
        locations_north   = ["Склад Чернігів, стелаж A", "Склад Чернігів, стелаж B", "Склад Суми, зона 1"]
        locations_south   = ["Склад Одеса, ангар 1", "Склад Одеса, ангар 2", "Склад Миколаїв, зона 3"]
        locations_central = ["Склад Київ, поверх 1", "Склад Київ, поверх 2", "Склад Київ, антресоль"]
        location_map = {
            "north": locations_north, "north_chernihiv": locations_north,
            "north_sumy": locations_north, "south": locations_south,
            "south_odesa": locations_south, "south_mykolaiv": locations_south,
            "central": locations_central, "central_kyiv": locations_central,
        }
        total_resources = 0
        for unit_key, count in counts_per_unit.items():
            uid = unit_ids[unit_key]
            for i in range(count):
                base_name = random.choice(RESOURCE_NAMES)
                name = f"{base_name} #{random.randint(100, 999)}"
                cat_id = random.choice(list(cat_ids.values()))
                qty = random.randint(1, 200)
                min_q = random.choice([2, 5, 10, 20])
                condition = random.choice(["NEW", "NEW", "GOOD", "GOOD", "FAIR", "DAMAGED"])
                location = random.choice(location_map[unit_key])
                rid = ensure_resource(
                    cur,
                    tenant_id=tenant_id,
                    category_id=cat_id,
                    unit_id=uid,
                    name=name,
                    description=f"Тестовий ресурс #{i + 1} для підрозділу {unit_key}",
                    quantity=qty,
                    condition=condition,
                    min_quantity=min_q,
                    serial_number=rnd_serial(),
                    location=location,
                )
                resource_ids_by_unit[unit_key].append(rid)
                total_resources += 1
        print(f"   ресурсів: {total_resources}")

        # ── Транспорт ────────────────────────────────────────────────────────
        print("→ Створюю транспорт (vehicles)…")
        # Водії — EMPLOYEE/TEAM_LEAD/BRANCH_LOGISTICIAN
        driver_pool_per_region = {
            "north":   ["emp1.ch@cargo-logistics.local", "emp2.ch@cargo-logistics.local",
                        "lead.ch@cargo-logistics.local"],
            "south":   ["log.odesa@cargo-logistics.local", "log.sumy@cargo-logistics.local"],
            "central": ["emp3.kyiv@cargo-logistics.local", "emp4.kyiv@cargo-logistics.local",
                        "lead.kyiv.a@cargo-logistics.local"],
        }
        vehicles_per_region = {"north": 6, "south": 5, "central": 7}
        vehicle_ids_by_region: dict[str, list[str]] = {}

        for reg_key, count in vehicles_per_region.items():
            ids = []
            drivers = driver_pool_per_region[reg_key]
            for i in range(count):
                brand, model_name, _vtype = VEHICLE_TYPES[i % len(VEHICLE_TYPES)]
                year = random.randint(2019, 2024)
                plate = rnd_plate()
                driver_id = user_ids.get(random.choice(drivers))
                vid = ensure_vehicle(
                    cur,
                    tenant_id=tenant_id,
                    plate=plate,
                    brand=brand,
                    model=f"{model_name} ({year})",
                    driver_id=driver_id,
                )
                ids.append(vid)
            vehicle_ids_by_region[reg_key] = ids
        total_vehicles = sum(len(v) for v in vehicle_ids_by_region.values())
        print(f"   транспортних засобів: {total_vehicles}")

        # ── Записи палива ────────────────────────────────────────────────────
        if "fuel_records" in existing:
            print("→ Заповнюю fuel_records…")
            fuel_cnt = 0
            for reg_key, veh_ids in vehicle_ids_by_region.items():
                for vid in veh_ids:
                    odometer = random.randint(25_000, 150_000)
                    for day in range(0, 45, random.choice([2, 3, 5])):
                        liters = round(random.uniform(20, 80), 2)
                        odometer += random.randint(100, 600)
                        ts = datetime.now(timezone.utc) - timedelta(days=45 - day)
                        cur.execute(
                            """INSERT INTO fuel_records
                                   (tenant_id, vehicle_id, liters, odometer_km, record_type, created_at)
                               VALUES (%s, %s, %s, %s, 'REFUEL', %s)""",
                            (tenant_id, vid, liters, odometer, ts),
                        )
                        fuel_cnt += 1
            print(f"   fuel_records: {fuel_cnt}")
        else:
            print("   ⚠ Таблиця fuel_records не знайдена — пропускаю.")

        # ── Заявки на постачання ─────────────────────────────────────────────
        print("→ Створюю supply_requests…")
        request_creators = [
            ("dir.north@cargo-logistics.local",   "north"),
            ("log.north@cargo-logistics.local",   "north"),
            ("mgr.chernihiv@cargo-logistics.local","north_chernihiv"),
            ("dir.south@cargo-logistics.local",   "south"),
            ("log.odesa@cargo-logistics.local",   "south_odesa"),
            ("mgr.kyiv@cargo-logistics.local",    "central_kyiv"),
            ("dept.kyiv@cargo-logistics.local",   "central_kyiv_dept"),
        ]
        supply_cnt = 0
        for creator_email, unit_key in request_creators:
            all_resources = resource_ids_by_unit.get(unit_key, [])
            if not all_resources:
                continue
            for _ in range(random.randint(3, 6)):
                res_id = random.choice(all_resources)
                status = random.choice(SUPPLY_STATUSES)
                approved_by = None
                approved_at = None
                # Знаходимо відповідного регіон-директора для затвердження
                region_key = unit_key.split("_")[0]  # "north", "south", "central"
                approver_email = f"dir.{region_key}@cargo-logistics.local"
                if status in ("APPROVED", "REJECTED", "COMPLETED"):
                    approved_by = user_ids.get(approver_email)
                    approved_at = datetime.now(timezone.utc) - timedelta(
                        days=random.randint(1, 20)
                    )
                cur.execute(
                    """INSERT INTO supply_requests
                           (tenant_id, created_by, resource_id, quantity, status,
                            approved_by, approved_at, comment, created_at, updated_at)
                       VALUES (%s, %s, %s, %s, %s, %s, %s, %s,
                               NOW() - INTERVAL '%s days',
                               NOW() - INTERVAL '%s days')""",
                    (
                        tenant_id,
                        user_ids[creator_email],
                        res_id,
                        random.randint(1, 30),
                        status,
                        approved_by,
                        approved_at,
                        random.choice(["Терміново", "Планова поставка", "Поповнення резерву", ""]),
                        random.randint(1, 30),
                        random.randint(0, 5),
                    ),
                )
                supply_cnt += 1
        print(f"   supply_requests: {supply_cnt}")

        # ── Волонтерські заявки ──────────────────────────────────────────────
        print("→ Створюю contractor_requests…")
        vol_creators = [
            ("dir.north@cargo-logistics.local",   "north"),
            ("mgr.odesa@cargo-logistics.local",   "south_odesa"),
            ("dir.central@cargo-logistics.local", "central"),
            ("dept.kyiv@cargo-logistics.local",   "central_kyiv_dept"),
        ]
        contractor_users = [
            user_ids["contractor1@cargo-logistics.local"],
            user_ids["contractor2@cargo-logistics.local"],
        ]
        cr_cnt = 0
        for creator_email, unit_key in vol_creators:
            for _ in range(random.randint(2, 4)):
                status = random.choice(CONTRACTOR_STATUSES)
                taken_by = None
                taken_at = None
                completed_at = None
                if status in ("IN_PROGRESS", "DELIVERED", "COMPLETED"):
                    taken_by = random.choice(contractor_users)
                    taken_at = datetime.now(timezone.utc) - timedelta(
                        days=random.randint(1, 10)
                    )
                if status == "COMPLETED":
                    completed_at = taken_at + timedelta(days=random.randint(1, 4))
                cur.execute(
                    """INSERT INTO contractor_requests
                           (tenant_id, created_by, unit_id, title, description,
                            status, taken_by, taken_at, completed_at, created_at)
                       VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s,
                               NOW() - INTERVAL '%s days')""",
                    (
                        tenant_id,
                        user_ids[creator_email],
                        unit_ids[unit_key],
                        random.choice(CONTRACTOR_TITLES),
                        "Заявка створена автоматично seed-скриптом для тестування.",
                        status,
                        taken_by,
                        taken_at,
                        completed_at,
                        random.randint(1, 25),
                    ),
                )
                cr_cnt += 1
        print(f"   contractor_requests: {cr_cnt}")

        # ── Опціональні таблиці ──────────────────────────────────────────────

        # warehouses (можуть бути від розширеної міграції)
        if "warehouses" in existing:
            print("→ Створюю warehouses…")
            wh_data = [
                (unit_ids["north_chernihiv"], "Головний склад Чернігів", 51.4982, 31.2893, "STATIONARY"),
                (unit_ids["north_sumy"],      "Склад Суми",              50.9077, 34.7981, "STATIONARY"),
                (unit_ids["south_odesa"],     "Склад Одеса-Порт",        46.4774, 30.7326, "STATIONARY"),
                (unit_ids["south_odesa"],     "Склад Одеса-Центр",       46.4825, 30.7233, "STATIONARY"),
                (unit_ids["south_mykolaiv"],  "Склад Миколаїв",          46.9750, 31.9946, "STATIONARY"),
                (unit_ids["central_kyiv"],    "Склад Київ-Північ",        50.5260, 30.5100, "STATIONARY"),
                (unit_ids["central_kyiv"],    "Склад Київ-Південь",       50.3930, 30.5450, "STATIONARY"),
                (unit_ids["central_kyiv"],    "Мобільний хаб Київ",      50.4501, 30.5234, "MOBILE"),
            ]
            wh_cnt = 0
            for (uid, wh_name, lat, lon, loc_type) in wh_data:
                row = cur.execute(
                    "SELECT id FROM warehouses WHERE name=%s AND unit_id=%s",
                    (wh_name, uid),
                ).fetchone()
                if not row:
                    cur.execute(
                        """INSERT INTO warehouses (tenant_id, unit_id, name, location_type, latitude, longitude)
                           VALUES (%s, %s, %s, %s, %s, %s)""",
                        (tenant_id, uid, wh_name, loc_type, lat, lon),
                    )
                    wh_cnt += 1
            print(f"   warehouses: {wh_cnt}")

        # gps_locations + geofences
        has_gps = {"geofences", "gps_locations"}.issubset(existing)
        if has_gps:
            print("→ Додаю GPS-трекінг та геозони…")
            geofences = [
                (unit_ids["north_chernihiv"], "Безпечна зона штабу Чернігів", 51.4982, 31.2893, 400, "SAFE"),
                (unit_ids["south_odesa"],     "Портова зона Одеса",           46.4774, 30.7326, 600, "SAFE"),
                (unit_ids["central_kyiv"],    "Логістичний центр Київ",       50.4501, 30.5234, 500, "SAFE"),
                (unit_ids["north"],           "Заборонена зона",              51.6000, 31.5000, 800, "FORBIDDEN"),
            ]
            for (uid, gname, lat, lon, radius, gtype) in geofences:
                cur.execute(
                    """INSERT INTO geofences (tenant_id, unit_id, name, latitude, longitude, radius, type, active)
                       VALUES (%s, %s, %s, %s, %s, %s, %s, TRUE)
                       ON CONFLICT DO NOTHING""",
                    (tenant_id, uid, gname, lat, lon, radius, gtype),
                )
            gps_cnt = 0
            for reg_key, veh_ids in vehicle_ids_by_region.items():
                base_coords = {
                    "north":   (51.4982, 31.2893),
                    "south":   (46.4774, 30.7326),
                    "central": (50.4501, 30.5234),
                }
                lat0, lon0 = base_coords.get(reg_key, (49.0, 32.0))
                for vid in veh_ids:
                    for i in range(15):
                        lat = lat0 + random.uniform(-0.04, 0.04)
                        lon = lon0 + random.uniform(-0.04, 0.04)
                        ts = datetime.now(timezone.utc) - timedelta(minutes=15 - i)
                        try:
                            cur.execute(
                                """INSERT INTO gps_locations
                                       (tenant_id, vehicle_id, unit_id, latitude, longitude,
                                        speed, heading, accuracy, timestamp)
                                   VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)""",
                                (
                                    tenant_id, vid, unit_ids[reg_key],
                                    lat, lon,
                                    round(random.uniform(0, 90), 2),
                                    round(random.uniform(0, 360), 2),
                                    round(random.uniform(2, 15), 2),
                                    ts,
                                ),
                            )
                            gps_cnt += 1
                        except psycopg.errors.Error:
                            conn.rollback()
                            break
            print(f"   geofences: {len(geofences)}, gps_locations: {gps_cnt}")
        else:
            print("   ⚠ Таблиці gps_locations/geofences не знайдені — пропускаю (запусти міграції).")

        # notifications
        if "notifications" in existing:
            print("→ Додаю тестові сповіщення (notifications)…")
            notif_users = [
                user_ids["dir.north@cargo-logistics.local"],
                user_ids["dir.south@cargo-logistics.local"],
                user_ids["dir.central@cargo-logistics.local"],
                user_ids["admin@cargo-logistics.local"],
            ]
            notif_types = ["LOW_STOCK", "SUPPLY_APPROVED", "SUPPLY_REJECTED", "VEHICLE_MAINTENANCE", "NEW_USER"]
            notif_cnt = 0
            for uid in notif_users:
                for _ in range(random.randint(3, 7)):
                    cur.execute(
                        """INSERT INTO notifications
                               (tenant_id, user_id, type, title, message, is_read, created_at)
                           VALUES (%s, %s, %s, %s, %s, %s, NOW() - INTERVAL '%s hours')""",
                        (
                            tenant_id, uid,
                            random.choice(notif_types),
                            "Тестове сповіщення",
                            "Це автоматично згенероване сповіщення від seed-скрипта.",
                            random.choice([True, False]),
                            random.randint(1, 72),
                        ),
                    )
                    notif_cnt += 1
            print(f"   notifications: {notif_cnt}")

        conn.commit()
        print("\n✅  Seed завершено успішно!")
        print_summary(tenant_id, unit_ids, user_ids)


def print_summary(tenant_id: str, unit_ids: dict[str, int], user_ids: dict[str, str]) -> None:
    print("\n============================================================")
    print(f"  Тенант : {TENANT_NAME}")
    print(f"  Slug   : {TENANT_SLUG}")
    print(f"  Tier   : {TENANT_TIER}")
    print(f"  ID     : {tenant_id}")
    print(f"  Units  : {len(unit_ids)}")
    print(f"  Users  : {len(user_ids)}")
    print()
    print("  Акаунти для входу (пароль: password123)")
    print(f"  {'TENANT_ADMIN':<20} admin@cargo-logistics.local")
    print(f"  {'REGION_DIRECTOR (N)':<20} dir.north@cargo-logistics.local")
    print(f"  {'REGION_DIRECTOR (S)':<20} dir.south@cargo-logistics.local")
    print(f"  {'REGION_DIRECTOR (C)':<20} dir.central@cargo-logistics.local")
    print(f"  {'BRANCH_MANAGER':<20} mgr.kyiv@cargo-logistics.local")
    print(f"  {'TEAM_LEAD':<20} lead.ch@cargo-logistics.local")
    print(f"  {'EMPLOYEE':<20} emp1.ch@cargo-logistics.local")
    print(f"  {'CONTRACTOR':<20} contractor1@cargo-logistics.local")
    print(f"  {'BLOCKED':<20} blocked@cargo-logistics.local")
    print("============================================================\n")


# ---------------------------------------------------------------------------
# Точка входу
# ---------------------------------------------------------------------------

def main() -> int:
    parser = argparse.ArgumentParser(
        description="Створити тенант 'Карго-Логістика' та заповнити БД тестовими даними."
    )
    parser.add_argument("--dsn", default=DEFAULT_DSN, help="PostgreSQL DSN")
    parser.add_argument(
        "--reset",
        action="store_true",
        help="Видалити всі дані тенанта та залити заново (CASCADE через tenant_id)",
    )
    args = parser.parse_args()

    print(f"→ Підключаюсь до {args.dsn}")
    try:
        with psycopg.connect(args.dsn, autocommit=False) as conn:
            seed(conn, reset=args.reset)
    except psycopg.OperationalError as exc:
        print(f"❌ Не вдалось під'єднатись: {exc}", file=sys.stderr)
        print(
            "   Перевір, що PostgreSQL запущений:  docker compose up -d postgres",
            file=sys.stderr,
        )
        return 1
    except Exception as exc:
        print(f"❌ Помилка: {exc}", file=sys.stderr)
        raise
    return 0


if __name__ == "__main__":
    sys.exit(main())
