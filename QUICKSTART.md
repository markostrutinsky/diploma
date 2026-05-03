# 🚀 Швидкий старт - Omnilog

Цей гайд допоможе вам запустити Omnilog за 5 хвилин!

## 📋 Вимоги

- **Docker** 20.10+ ([Встановити](https://docs.docker.com/get-docker/))
- **Docker Compose** 2.0+ (зазвичай йде з Docker Desktop)
- **Git** ([Встановити](https://git-scm.com/downloads))

## ⚡ Швидкий старт (5 хвилин)

### Крок 1: Клонування репозиторію

```bash
git clone https://github.com/markostrutinsky/diploma.git
cd diploma
```

### Крок 2: Налаштування змінних середовища

```bash
# Скопіюйте приклад конфігурації
cp .env.example .env

# Відредагуйте .env файл (опціонально для локального запуску)
nano .env  # або vim, або ваш улюблений редактор
```

**Мінімальна конфігурація для локального запуску:**
```env
DATABASE_URL=postgresql://Omnilog:Omnilog123@postgres:5432/Omnilog?sslmode=disable
JWT_SECRET=dev-secret-for-local-testing
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your@gmail.com
SMTP_PASSWORD=your-app-password
FRONTEND_URL=http://localhost
PORT=8080
ALLOW_BOOTSTRAP_OVERRIDE=true
```

> 💡 **Порада**: Для email нотифікацій використовуйте [Gmail App Password](https://support.google.com/accounts/answer/185833)

### Крок 3: Запуск

```bash
docker compose up -d --build
```

Це займе 2-3 хвилини. Docker завантажить та зібере всі необхідні образи.

### Крок 4: Перевірка статусу

```bash
docker compose ps
```

Ви повинні побачити 3 запущені контейнери:
- ✅ `diploma-frontend`
- ✅ `diploma-backend`
- ✅ `diploma-postgres`

### Крок 5: Відкрийте додаток

🌐 **Перейдіть за адресою:** http://localhost

---

## 🎯 Перші кроки після запуску

### Варіант 1: Platform Admin (управління всією системою)

1. Відкрийте http://localhost/platform
2. Введіть credentials:
   ```
   Email:    platform@omnilog.system
   Password: AdminSystem2024!
   ```
3. Ви побачите дашборд з усіма організаціями

**Що можна робити:**
- ✅ Переглядати всі організації
- ✅ Змінювати тарифні плани
- ✅ Блокувати/активувати організації
- ✅ Переглядати глобальну статистику

---

### Варіант 2: Створити свою організацію

1. Відкрийте http://localhost/signup
2. Заповніть форму:
   ```
   Назва організації: 26-та Окрема Бригада
   Email:            commander@brigade26.mil
   Пароль:           SecurePassword123!
   ```
3. Автоматично отримаєте:
   - 🎟️ Тариф: **FREE**
   - 👤 Роль: **TENANT_ADMIN**
   - 🏢 Власний ізольований простір

**Що далі:**
- ✅ Створіть підрозділи (Units)
- ✅ Додайте користувачів (Users)
- ✅ Створіть склади (Warehouses)
- ✅ Додайте категорії та ресурси (Inventory)
- ✅ Створіть першу заявку (Requests)

---

### Варіант 3: Волонтер/Підрядник

1. Відкрийте http://localhost/register
2. Зареєструйтесь:
   ```
   Ім'я:         Іван Петренко
   Email:        volunteer@ngo.org
   Пароль:       SecurePass123!
   Організація:  NGO Допомога (опціонально)
   ```
3. Автоматично отримаєте роль **CONTRACTOR**

**Що можна робити:**
- ✅ Переглядати відкриті заявки від військових
- ✅ Брати заявки в роботу
- ✅ Позначати доставку виконаною
- ✅ Комунікувати з військовими підрозділами

---

## 📚 Наступні кроки

### 1. Ознайомтеся з документацією

- 📖 [README.md](README.md) - Загальний огляд проєкту
- 📖 [API_DOCUMENTATION.md](API_DOCUMENTATION.md) - Детальна API документація (2841 рядків!)
- 📖 [ARCHITECTURE.md](ARCHITECTURE.md) - Архітектура системи

### 2. Експериментуйте з API

```bash
# Логін
curl -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "commander@brigade26.mil",
    "password": "SecurePassword123!"
  }'

# Збережіть токен
export TOKEN="your-jwt-token-here"

# Отримати список ресурсів
curl -X GET http://localhost/api/inventory/resources \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Спробуйте Premium Features

Для тестування PRO features:

1. Увійдіть як **Platform Admin**
2. Оновіть тариф вашої організації до **PRO**
3. Тепер доступні:
   - 📊 Advanced KPI Dashboard
   - 🔮 Demand Forecasting (AI)
   - 🔧 Predictive Maintenance
   - 🛡️ Fuel Anti-Fraud Detection
   - 🌍 Real-Time GPS Tracking

---

## 🔧 Корисні команди

### Перезапустити сервіси

```bash
docker compose restart
```

### Переглянути логи

```bash
# Всі сервіси
docker compose logs -f

# Тільки backend
docker compose logs -f backend

# Тільки frontend
docker compose logs -f frontend
```

### Зупинити систему

```bash
docker compose down
```

### Повністю очистити (включно з даними БД)

```bash
docker compose down -v
```

### Перезібрати образи

```bash
docker compose up -d --build
```

### Підключитись до БД

```bash
docker compose exec postgres psql -U Omnilog -d Omnilog
```

Приклади SQL запитів:
```sql
-- Кількість користувачів
SELECT COUNT(*) FROM users;

-- Список організацій
SELECT id, name, subscription_tier FROM tenants;

-- Критичні ресурси
SELECT name, quantity, min_quantity 
FROM resources 
WHERE quantity <= min_quantity;
```

---

## 🐛 Troubleshooting

### Проблема: Порт вже зайнятий

**Помилка:**
```
Error: bind: address already in use
```

**Рішення:**
```bash
# Знайти процес на порті 80
sudo lsof -i :80

# Або зупинити всі Docker контейнери
docker stop $(docker ps -aq)
```

---

### Проблема: Backend не стартує

**Помилка в логах:**
```
Failed to connect to database
```

**Рішення:**
```bash
# Перевірте, чи запущений PostgreSQL
docker compose ps postgres

# Перезапустіть БД
docker compose restart postgres

# Почекайте 10 секунд та перезапустіть backend
docker compose restart backend
```

---

### Проблема: Email не відправляються

**Рішення:**

1. Перевірте `.env` файл:
   ```env
   SMTP_EMAIL=your@gmail.com
   SMTP_PASSWORD=your-16-char-app-password  # НЕ звичайний пароль!
   ```

2. Для Gmail створіть App Password:
   - Перейдіть на https://myaccount.google.com/security
   - Увімкніть 2FA (якщо не увімкнено)
   - App passwords → Generate
   - Скопіюйте 16-значний код

3. Перезапустіть backend:
   ```bash
   docker compose restart backend
   ```

---

### Проблема: Frontend показує білий екран

**Рішення:**
```bash
# Перевірте логи frontend
docker compose logs frontend

# Перезберіть frontend
docker compose up -d --build frontend
```

---

### Проблема: Міграції не застосувались

**Симптоми:** Помилки про відсутні таблиці

**Рішення:**
```bash
# Повністю очистіть БД та перезапустіть
docker compose down -v
docker compose up -d --build

# Або вручну запустіть міграції
docker compose exec backend /app/Omnilog_backend migrate
```

---

## 📊 Тестові дані

Хочете швидко заповнити систему тестовими даними?

```bash
# GPS симулятор (створює траєкторії транспорту)
cd scripts/gps_simulator
pip install -r requirements.txt
python simulate.py --vehicles 5 --duration 3600

# Database seeder (створює тестові ресурси, заявки)
cd scripts/seed_db
pip install -r requirements.txt
python seed.py
```

---

## 🎓 Навчальні матеріали

### Відео туторіали (Coming soon)
- [ ] Створення організації
- [ ] Управління користувачами
- [ ] Складський облік
- [ ] GPS трекінг
- [ ] Аналітика та звіти

### Приклади використання
- [Scenario 1: Військова бригада](docs/scenarios/military-brigade.md)
- [Scenario 2: Волонтерська організація](docs/scenarios/volunteer-ngo.md)
- [Scenario 3: Комерційна логістика](docs/scenarios/commercial-logistics.md)

---

## 💬 Підтримка

**Потрібна допомога?**

- 📧 Email: support@omnilog.system
- 💬 GitHub Discussions: [Задати питання](https://github.com/markostrutinsky/diploma/discussions)
- 🐛 Bug Reports: [Створити Issue](https://github.com/markostrutinsky/diploma/issues)
- 📖 Документація: [Повний список docs](README.md#-зміст)

---

## ✅ Checklist для першого запуску

- [ ] Docker встановлено та запущено
- [ ] Репозиторій клоновано
- [ ] `.env` файл створено
- [ ] SMTP credentials налаштовано (для email)
- [ ] `docker compose up -d --build` виконано успішно
- [ ] Всі 3 контейнери запущені (`docker compose ps`)
- [ ] http://localhost відкривається в браузері
- [ ] Увійшли в систему (Platform Admin або створили організацію)
- [ ] Створили перший ресурс або заявку

**Все працює? Чудово! 🎉**

---

## 🚀 Готові до production?

Перед розгортанням на production:

1. **Безпека:**
   - [ ] Змініть `JWT_SECRET` на сильний випадковий ключ
   - [ ] Встановіть `ALLOW_BOOTSTRAP_OVERRIDE=false`
   - [ ] Використовуйте сильний пароль для БД
   - [ ] Увімкніть HTTPS (Caddy робить це автоматично)

2. **Backup:**
   - [ ] Налаштуйте автоматичні backup БД
   - [ ] Налаштуйте моніторинг (Prometheus, Grafana)
   - [ ] Налаштуйте алерти (Sentry, PagerDuty)

3. **Продуктивність:**
   - [ ] Налаштуйте connection pooling
   - [ ] Увімкніть Redis для кешування
   - [ ] Налаштуйте CDN для статичних файлів

4. **Документація:**
   - [ ] Прочитайте [Production Deployment Guide](docs/DEPLOYMENT.md)
   - [ ] Ознайомтесь з [Security Best Practices](docs/SECURITY.md)

---

**Вдалого використання! 🎖️**

Made with ❤️ in Ukraine 🇺🇦
