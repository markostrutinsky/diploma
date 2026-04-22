# GPS Simulator

Імітує IoT-трекери вантажівок на сторінці **GPS Трекінг** (PRO-фіча).

За замовчуванням працює в **shipment-aware** режимі: симулює тільки ті машини,
які зараз в рейсі (`shipments.status = 'DISPATCHED'`), везе їх між складами,
а після `DELIVERED` — повертає на вихідний склад і прибирає з карти.

## Фази рейсу

1. **OUTBOUND** — машина з'являється у вихідному складі (`from_warehouse_id`)
   і їде в склад призначення (`to_warehouse_id`).
2. **WAITING** — доїхала до призначення, стоїть там з нульовою швидкістю, поки
   приймаючий склад не натисне «прийняв» у UI → `shipments.status` стає `DELIVERED`.
3. **RETURN** — симулятор помічає `DELIVERED` і розвертає машину назад до
   `from_warehouse_id`.
4. **DONE** — доїхала назад → `DELETE FROM gps_locations WHERE vehicle_id=...`,
   машина зникає з `/gps`.

## Запуск

```bash
cd scripts/gps_simulator
python3 -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt

# Нескінченно, тільки машини з активних рейсів, tick кожні 5 сек:
python simulate.py

# Прискорена демонстрація (tick кожні 2 сек):
python simulate.py --interval 2

# 100 тіків і вихід:
python simulate.py --iterations 100

# Старий хаотичний режим (N машин блукають випадково, без рейсів):
python simulate.py --mode free --vehicles 5

# Інший DSN:
python simulate.py --dsn 'postgres://postgres:postgres@localhost:5432/omnilog'
```

## Передумови

- Бекенд з мігрованою БД (таблиці `gps_locations`, `shipments`, `warehouses`).
- Щонайменше один рейс у БД зі статусом `DISPATCHED`. Створити можна через UI
  («Склади → Створити рейс» / «Запити → Відправити»).
- Склади-учасники рейсу повинні мати заповнені `latitude`/`longitude`.

Якщо активних рейсів нема — симулятор працює, просто нічого не пише і друкує
`active_trips=0` на кожному тіку. Щойно з'явиться новий DISPATCHED — машина
автоматично з'явиться на карті.

Зупинка: `Ctrl+C`.
