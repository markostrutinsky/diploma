# Міграція з Caddy на nginx + mkcert ✅

## Що було зроблено

### ❌ Було (Caddy)
- Самопідписані сертифікати
- Браузери показували попередження "небезпечний сайт"
- Потрібно було кожного разу підтверджувати виключення

### ✅ Стало (nginx + mkcert)
- Довірені локальні сертифікати
- Браузери автоматично довіряють
- Жодних попереджень безпеки

## Зміни в проекті

### 1. Додано файли
- `setup-certs.sh` - скрипт для генерації сертифікатів
- `docker/nginx/nginx.conf` - конфігурація nginx
- `certs/` - директорія для сертифікатів
- `.gitignore` - ігнорування сертифікатів
- `HTTPS_SETUP.md` - детальна документація

### 2. Змінено файли
- `docker-compose.yml` - замінено Caddy на nginx

### 3. Видалено
- `Caddyfile` - більше не потрібен
- Caddy volumes з docker-compose.yml

## Як запустити проект

```bash
# 1. Згенерувати сертифікати (один раз)
./setup-certs.sh

# 2. Запустити проект
docker-compose up --build -d

# 3. Відкрити в браузері
# https://localhost
```

## Перевірка

Ваш сайт тепер доступний за адресою **https://localhost** без жодних попереджень! 🎉

### Тестування
```bash
# Перевірити статус
docker-compose ps

# Перевірити HTTPS
curl -I https://localhost

# Переглянути логи nginx
docker-compose logs nginx
```

## Важливо

🔒 Сертифікати мають термін дії **2 роки** (до 26 липня 2028)  
🔄 Після закінчення просто виконайте `./setup-certs.sh` знову  
🚫 Сертифікати в `certs/` НЕ комітяться в git  
⚠️ Після встановлення CA **перезапустіть браузер**

## Налаштування для команди

Кожен розробник повинен:

1. Встановити mkcert (один раз):
```bash
# Ubuntu/Debian
sudo apt install libnss3-tools
wget https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-linux-amd64
sudo mv mkcert-v1.4.4-linux-amd64 /usr/local/bin/mkcert
sudo chmod +x /usr/local/bin/mkcert
```

2. Згенерувати сертифікати (в директорії проекту):
```bash
./setup-certs.sh
```

3. Перезапустити браузер

4. Запустити проект:
```bash
docker-compose up --build
```

## Troubleshooting

### Браузер все ще показує попередження?
- Переконайтесь, що запустили `./setup-certs.sh`
- Перезапустіть браузер
- Очистіть кеш (Ctrl+Shift+Delete)
- Переконайтесь, що встановлено `libnss3-tools`

### Порт 443 зайнятий?
```bash
sudo lsof -i :443
# Зупиніть інший сервіс або змініть порт в docker-compose.yml
```

### nginx не стартує?
```bash
docker-compose logs nginx
# Перевірте, чи існують файли сертифікатів:
ls -la certs/
```

---

Дякую за міграцію! Тепер локальна розробка стала зручнішою! 🚀
