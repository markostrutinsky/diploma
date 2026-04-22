#!/usr/bin/env python3
"""
GPS Simulator для MilLog / OmniLog.

Імітує передачу GPS-координат від транспортних засобів, які ЗАРАЗ перебувають
у рейсі (таблиця `shipments`). На карті «GPS Трекінг» з'являються тільки машини
з активним рейсом; коли рейс позначено DELIVERED — машина їде назад, а коли
повертається у вихідний склад — пропадає з карти.

Фази життя машини в симуляції:
    OUTBOUND  — рухається від from_warehouse до to_warehouse (shipment=DISPATCHED)
    WAITING   — доїхала до to_warehouse, стоїть (швидкість 0) поки приймаючий
                склад не натисне «прийняв». У БД shipment все ще DISPATCHED.
    RETURN    — shipment перейшов у DELIVERED → машина повертається назад.
    DONE      — повернулась у вихідний склад → рядки GPS видаляються і
                машина зникає з /gps.

Запуск:
  cd scripts/gps_simulator
  python3 -m venv .venv && source .venv/bin/activate
  pip install -r requirements.txt
  python simulate.py                          # нескінченно, tick кожні 5 сек
  python simulate.py --interval 2             # частіше (для демо)
  python simulate.py --iterations 100         # N тіків і вихід
  python simulate.py --mode free --vehicles 5 # старий хаотичний режим

Вимоги до БД:
  * міграції Go-бекенду застосовані (gps_locations, shipments, warehouses);
  * у shipments є запис зі статусом DISPATCHED і координатами в from/to
    складах (створюється через UI «Створити рейс»).
"""

from __future__ import annotations

import argparse
import math
import os
import random
import signal
import sys
import time
from dataclasses import dataclass
from datetime import datetime, timezone

import psycopg
from psycopg.rows import dict_row

DEFAULT_DSN = os.environ.get(
    "SEED_DSN",
    "postgres://postgres:postgres@localhost:5432/omnilog?sslmode=disable",
)

# Використовується тільки в --mode free
START_POINTS = [
    (50.4501, 30.5234),
    (50.4021, 30.6525),
    (50.5010, 30.4720),
    (49.5883, 34.5514),
    (49.8397, 24.0297),
]

# ---------------------------------------------------------------------------
# Гео-утиліти
# ---------------------------------------------------------------------------

EARTH_KM = 6371.0


def haversine_km(a_lat: float, a_lon: float, b_lat: float, b_lon: float) -> float:
    lat1, lat2 = math.radians(a_lat), math.radians(b_lat)
    dlat = math.radians(b_lat - a_lat)
    dlon = math.radians(b_lon - a_lon)
    h = math.sin(dlat / 2) ** 2 + math.cos(lat1) * math.cos(lat2) * math.sin(dlon / 2) ** 2
    return 2 * EARTH_KM * math.asin(math.sqrt(h))


def bearing_deg(a_lat: float, a_lon: float, b_lat: float, b_lon: float) -> float:
    lat1, lat2 = math.radians(a_lat), math.radians(b_lat)
    dlon = math.radians(b_lon - a_lon)
    y = math.sin(dlon) * math.cos(lat2)
    x = math.cos(lat1) * math.sin(lat2) - math.sin(lat1) * math.cos(lat2) * math.cos(dlon)
    return (math.degrees(math.atan2(y, x)) + 360) % 360


def move_towards(
    lat: float, lon: float, target_lat: float, target_lon: float, km: float
) -> tuple[float, float, float]:
    """Пересуває точку у напрямку цілі на km. Якщо залишок <= km — паркує в цілі.
    Повертає (new_lat, new_lon, distance_left_km)."""
    dist = haversine_km(lat, lon, target_lat, target_lon)
    if dist <= km or dist < 1e-6:
        return target_lat, target_lon, 0.0
    brg = math.radians(bearing_deg(lat, lon, target_lat, target_lon))
    lat_delta = (km / 111.0) * math.cos(brg)
    lon_delta = (km / (111.0 * math.cos(math.radians(lat)))) * math.sin(brg)
    return lat + lat_delta, lon + lon_delta, dist - km


# ---------------------------------------------------------------------------
# Shipment-aware режим
# ---------------------------------------------------------------------------

PHASE_OUTBOUND = "OUTBOUND"
PHASE_WAITING = "WAITING"
PHASE_RETURN = "RETURN"


@dataclass
class Trip:
    shipment_id: str
    vehicle_id: str
    unit_id: int
    from_lat: float
    from_lon: float
    to_lat: float
    to_lon: float
    lat: float
    lon: float
    heading: float
    speed: float
    cruise_kmh: float
    phase: str = PHASE_OUTBOUND


def fetch_dispatched_trips(conn: psycopg.Connection) -> list[dict]:
    """Усі активні рейси (DISPATCHED + DELIVERED) з координатами складів.
    DELIVERED потрібні, щоб докатати машину назад і потім видалити з карти."""
    query = """
        SELECT
            s.id::text        AS shipment_id,
            s.vehicle_id::text AS vehicle_id,
            s.status          AS status,
            wf.unit_id        AS unit_id,
            wf.latitude       AS from_lat,
            wf.longitude      AS from_lon,
            wt.latitude       AS to_lat,
            wt.longitude      AS to_lon
        FROM shipments s
        JOIN warehouses wf ON wf.id = s.from_warehouse_id
        JOIN warehouses wt ON wt.id = s.to_warehouse_id
        WHERE s.status IN ('DISPATCHED', 'DELIVERED')
          AND wf.latitude IS NOT NULL AND wf.longitude IS NOT NULL
          AND wt.latitude IS NOT NULL AND wt.longitude IS NOT NULL
    """
    with conn.cursor(row_factory=dict_row) as cur:
        cur.execute(query)
        return cur.fetchall()


def insert_point(
    conn: psycopg.Connection,
    vehicle_id: str,
    unit_id: int,
    lat: float,
    lon: float,
    speed: float,
    heading: float,
) -> None:
    with conn.cursor() as cur:
        cur.execute(
            """
            INSERT INTO gps_locations
                (vehicle_id, unit_id, latitude, longitude, altitude, speed, heading, accuracy, timestamp)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
            """,
            (
                vehicle_id,
                unit_id,
                round(lat, 8),
                round(lon, 8),
                round(random.uniform(120, 220), 2),
                round(speed, 2),
                round(heading, 2),
                round(random.uniform(3, 8), 2),
                datetime.now(tz=timezone.utc),
            ),
        )
    conn.commit()


def delete_vehicle_points(conn: psycopg.Connection, vehicle_id: str) -> None:
    """Прибирає всі точки машини — після цього вона зникає з fleet-map."""
    with conn.cursor() as cur:
        cur.execute("DELETE FROM gps_locations WHERE vehicle_id = %s", (vehicle_id,))
    conn.commit()


def run_shipments_mode(conn: psycopg.Connection, args: argparse.Namespace, stop: dict) -> None:
    trips: dict[str, Trip] = {}
    iteration = 0
    while not stop["flag"]:
        iteration += 1
        try:
            rows = fetch_dispatched_trips(conn)
        except Exception as e:  # noqa: BLE001
            print(f"⚠️  DB помилка: {e}", file=sys.stderr)
            conn.rollback()
            rows = []

        active_ids = set()
        for row in rows:
            sid = row["shipment_id"]
            active_ids.add(sid)
            trip = trips.get(sid)
            if trip is None:
                trip = Trip(
                    shipment_id=sid,
                    vehicle_id=row["vehicle_id"],
                    unit_id=int(row["unit_id"] or 0),
                    from_lat=float(row["from_lat"]),
                    from_lon=float(row["from_lon"]),
                    to_lat=float(row["to_lat"]),
                    to_lon=float(row["to_lon"]),
                    lat=float(row["from_lat"]),
                    lon=float(row["from_lon"]),
                    heading=bearing_deg(
                        float(row["from_lat"]), float(row["from_lon"]),
                        float(row["to_lat"]), float(row["to_lon"]),
                    ),
                    speed=random.uniform(45, 65),
                    cruise_kmh=random.uniform(55, 75),
                )
                trips[sid] = trip
                print(
                    f"➕ новий рейс {sid[:8]}…  vehicle={trip.vehicle_id[:8]}…  "
                    f"{trip.from_lat:.3f},{trip.from_lon:.3f} → {trip.to_lat:.3f},{trip.to_lon:.3f}"
                )

            # Статус у БД змінився — переводимо в RETURN
            if row["status"] == "DELIVERED" and trip.phase in (PHASE_OUTBOUND, PHASE_WAITING):
                trip.phase = PHASE_RETURN
                trip.speed = max(trip.speed, 40.0)
                trip.heading = bearing_deg(trip.lat, trip.lon, trip.from_lat, trip.from_lon)
                print(f"↩️  рейс {sid[:8]}… DELIVERED — машина їде назад")

        # Рух
        for sid, trip in list(trips.items()):
            if sid not in active_ids:
                # рейс видалено / скасовано в БД
                delete_vehicle_points(conn, trip.vehicle_id)
                trips.pop(sid, None)
                continue

            km_per_step = trip.speed * (args.interval / 3600.0)

            if trip.phase == PHASE_OUTBOUND:
                trip.lat, trip.lon, left = move_towards(
                    trip.lat, trip.lon, trip.to_lat, trip.to_lon, km_per_step
                )
                trip.heading = bearing_deg(trip.lat, trip.lon, trip.to_lat, trip.to_lon)
                trip.speed = trip.cruise_kmh + random.uniform(-8, 8)
                if left == 0.0:
                    trip.phase = PHASE_WAITING
                    trip.speed = 0.0
                    print(
                        f"🅿️  машина {trip.vehicle_id[:8]}… на місці призначення "
                        f"(рейс {sid[:8]}…) — чекає приймання"
                    )
            elif trip.phase == PHASE_WAITING:
                trip.speed = 0.0
            elif trip.phase == PHASE_RETURN:
                trip.lat, trip.lon, left = move_towards(
                    trip.lat, trip.lon, trip.from_lat, trip.from_lon, km_per_step
                )
                trip.heading = bearing_deg(trip.lat, trip.lon, trip.from_lat, trip.from_lon)
                trip.speed = trip.cruise_kmh + random.uniform(-8, 8)
                if left == 0.0:
                    print(
                        f"🏁 машина {trip.vehicle_id[:8]}… повернулась на базу "
                        f"(рейс {sid[:8]}…) — зникає з карти"
                    )
                    delete_vehicle_points(conn, trip.vehicle_id)
                    trips.pop(sid, None)
                    continue

            try:
                insert_point(
                    conn, trip.vehicle_id, trip.unit_id,
                    trip.lat, trip.lon, trip.speed, trip.heading,
                )
            except Exception as e:  # noqa: BLE001
                print(f"⚠️  insert failed vehicle={trip.vehicle_id[:8]}…: {e}", file=sys.stderr)
                conn.rollback()

        print(
            f"✓ tick #{iteration}  active_trips={len(trips)}  "
            f"({datetime.now().strftime('%H:%M:%S')})"
        )

        if args.iterations and iteration >= args.iterations:
            print(f"✅ Досягнуто ліміту {args.iterations} ітерацій.")
            return

        slept = 0.0
        while slept < args.interval and not stop["flag"]:
            time.sleep(min(0.5, args.interval - slept))
            slept += 0.5


# ---------------------------------------------------------------------------
# Старий вільний режим (залишено як fallback)
# ---------------------------------------------------------------------------


@dataclass
class FreeState:
    vehicle_id: str
    unit_id: int
    lat: float
    lon: float
    heading: float
    speed: float


def pick_free_vehicles(conn: psycopg.Connection, limit: int) -> list[FreeState]:
    with conn.cursor(row_factory=dict_row) as cur:
        cur.execute("SELECT id FROM units ORDER BY id LIMIT 1")
        unit_row = cur.fetchone()
        default_unit_id = int(unit_row["id"]) if unit_row else 1

        cur.execute(
            """
            SELECT v.id::text AS id
            FROM vehicles v
            WHERE v.status IN ('ACTIVE', 'ON_MISSION')
              AND v.type IN ('VAN', 'TRUCK', 'PICKUP')
            ORDER BY random()
            LIMIT %s
            """,
            (limit,),
        )
        rows = cur.fetchall()

    if not rows:
        print("⚠️  Не знайдено жодної машини ACTIVE/ON_MISSION.", file=sys.stderr)
        sys.exit(1)

    states: list[FreeState] = []
    for i, row in enumerate(rows):
        s = START_POINTS[i % len(START_POINTS)]
        states.append(FreeState(
            vehicle_id=row["id"],
            unit_id=default_unit_id,
            lat=s[0] + random.uniform(-0.005, 0.005),
            lon=s[1] + random.uniform(-0.005, 0.005),
            heading=random.uniform(0, 360),
            speed=random.uniform(25, 45),
        ))
    return states


def run_free_mode(conn: psycopg.Connection, args: argparse.Namespace, stop: dict) -> None:
    states = pick_free_vehicles(conn, args.vehicles)
    print(f"🚚 [free] Симулюю {len(states)} машин хаотично.")
    iteration = 0
    while not stop["flag"]:
        iteration += 1
        for s in states:
            s.heading = (s.heading + random.uniform(-15, 15)) % 360
            s.speed = max(5.0, min(70.0, s.speed + random.uniform(-5, 5)))
            km = s.speed * (args.interval / 3600.0)
            brg = math.radians(s.heading)
            s.lat += (km / 111.0) * math.cos(brg)
            s.lon += (km / (111.0 * math.cos(math.radians(s.lat)))) * math.sin(brg)
            insert_point(conn, s.vehicle_id, s.unit_id, s.lat, s.lon, s.speed, s.heading)
        print(f"✓ tick #{iteration}  ({datetime.now().strftime('%H:%M:%S')})")
        if args.iterations and iteration >= args.iterations:
            return
        slept = 0.0
        while slept < args.interval and not stop["flag"]:
            time.sleep(min(0.5, args.interval - slept))
            slept += 0.5


# ---------------------------------------------------------------------------
# entrypoint
# ---------------------------------------------------------------------------


def main() -> None:
    parser = argparse.ArgumentParser(description="GPS simulator for MilLog vehicles")
    parser.add_argument("--dsn", default=DEFAULT_DSN, help="PostgreSQL DSN (env SEED_DSN)")
    parser.add_argument(
        "--mode",
        choices=["shipments", "free"],
        default="shipments",
        help="shipments — тільки машини в рейсі; free — хаотичне блукання N машин",
    )
    parser.add_argument("--vehicles", type=int, default=5, help="[free] скільки машин")
    parser.add_argument("--interval", type=float, default=5.0, help="Інтервал між тіками, сек")
    parser.add_argument("--iterations", type=int, default=0, help="К-сть тіків (0 = нескінченно)")
    args = parser.parse_args()

    print(f"🛰️  GPS Simulator  mode={args.mode}  interval={args.interval}s  DSN={args.dsn}")
    conn = psycopg.connect(args.dsn)

    # Очищаємо застарілі точки (старші 5 хв) — щоб на карті не висіли фантоми
    # від попередніх запусків чи вже доставлених рейсів.
    with conn.cursor() as cur:
        cur.execute(
            "DELETE FROM gps_locations WHERE timestamp < NOW() - INTERVAL '5 minutes'"
        )
        deleted = cur.rowcount
    conn.commit()
    if deleted:
        print(f"🧹 Видалено {deleted} застарілих GPS-точок (старші 5 хв).")

    stop = {"flag": False}

    def _sigint(_sig, _frm):
        stop["flag"] = True
        print("\n🛑 Зупинка після поточної ітерації…")

    signal.signal(signal.SIGINT, _sigint)
    signal.signal(signal.SIGTERM, _sigint)

    try:
        if args.mode == "shipments":
            run_shipments_mode(conn, args, stop)
        else:
            run_free_mode(conn, args, stop)
    finally:
        conn.close()
        print("👋 GPS Simulator завершено.")


if __name__ == "__main__":
    main()
