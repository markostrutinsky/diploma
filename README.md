# MilLog — Система управління військовою логістикою

## Запуск (Docker)

```bash
docker compose up -d --build
```

- **Фронт:** http://localhost
- **API:** http://localhost/api

### Перший запуск

1. Відкрийте http://localhost/bootstrap
2. Створіть обліковий запис адміністратора
3. Увійдіть на http://localhost/login

### Змінні середовища (.env)

```env
JWT_SECRET=your-secret-key
SMTP_EMAIL=your@gmail.com
SMTP_PASSWORD=app-password
FRONTEND_URL=http://localhost
# Дозволити bootstrap коли вже є користувачі (для dev/reset)
ALLOW_BOOTSTRAP_OVERRIDE=true
```

## Функціонал

- **Auth:** JWT, логін, setup-password (invite), bootstrap (перший адмін)
- **Користувачі:** Адмін створює користувачів (invite по email)
- **Ресурси:** Категорії, CRUD ресурсів, критичні залишки (min_quantity)
- **Заявки:** Волонтер створює, командир/адмін затверджує

## Ролі

| Роль | Можливості |
|------|------------|
| ADMIN | Користувачі, категорії, ресурси, заявки |
| WAREHOUSE | Категорії, ресурси |
| COMMANDER | Перегляд, затвердження заявок |
| CONTRACTOR | Заявки на постачання |

## Розробка

```bash
# Backend
cd millog_backend && go run .

# Frontend
cd millog_frontend && npm install && npm run dev
```

**Примітка:** Якщо БД вже існує, нові таблиці (inventory, requests, vehicles) можуть не з’явитися. Для чистого старту або якщо міграції не пройшли:
```bash
docker compose down -v
docker compose up -d --build
```
