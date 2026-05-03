# Налаштування HTTPS для локальної розробки

Цей проект використовує **mkcert** для створення довірених локальних HTTPS сертифікатів замість самопідписаних сертифікатів.

## Переваги mkcert

✅ Браузери довіряють сертифікатам автоматично  
✅ Немає попереджень "Ваше з'єднання не є приватним"  
✅ Просте налаштування одною командою  
✅ Працює з localhost, 127.0.0.1 та іншими доменами

## Швидкий старт

### 1. Встановлення mkcert

#### Ubuntu/Debian
```bash
sudo apt install libnss3-tools
wget https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-linux-amd64
sudo mv mkcert-v1.4.4-linux-amd64 /usr/local/bin/mkcert
sudo chmod +x /usr/local/bin/mkcert
```

#### Arch Linux
```bash
sudo pacman -S mkcert
```

#### Fedora
```bash
sudo dnf install mkcert
```

#### macOS
```bash
brew install mkcert
```

### 2. Генерація сертифікатів

Виконайте скрипт налаштування:

```bash
./setup-certs.sh
```

Цей скрипт:
- Встановить локальний Certificate Authority
- Згенерує сертифікати для localhost
- Збереже їх у директорію `./certs/`

### 3. Запуск проекту

```bash
docker-compose up --build
```

Сайт буде доступний за адресою: **https://localhost**

## Структура файлів

```
diploma/
├── certs/                      # Сертифікати (не комітьте в git!)
│   ├── localhost.pem          # SSL сертифікат
│   └── localhost-key.pem      # Приватний ключ
├── docker/
│   └── nginx/
│       └── nginx.conf         # Конфігурація nginx з HTTPS
├── docker-compose.yml         # Docker compose з nginx замість caddy
└── setup-certs.sh            # Скрипт для генерації сертифікатів
```

## Усунення проблем

### Браузер все ще показує попередження

1. Переконайтеся, що ви запустили `./setup-certs.sh`
2. Перезапустіть браузер після встановлення сертифікатів
3. Очистіть кеш браузера (Ctrl+Shift+Delete)

### Помилка "connection refused"

Переконайтеся, що контейнери запущені:
```bash
docker-compose ps
```

### Порти вже зайняті

Якщо порти 80 або 443 вже використовуються, зупиніть інші сервіси:
```bash
sudo systemctl stop apache2  # або nginx, якщо встановлено системно
sudo lsof -i :80
sudo lsof -i :443
```

## Додаткові домени

Якщо потрібно додати інші локальні домени (наприклад, `Omnilog.local`):

1. Додайте домен до `/etc/hosts`:
```bash
echo "127.0.0.1 Omnilog.local" | sudo tee -a /etc/hosts
```

2. Згенеруйте сертифікат з цим доменом:
```bash
cd certs
mkcert localhost Omnilog.local 127.0.0.1 ::1
```

3. Оновіть `nginx.conf`, додавши домен до `server_name`:
```nginx
server_name localhost Omnilog.local;
```

## Видалення

Щоб видалити локальний CA:
```bash
mkcert -uninstall
```

## Примітки безпеки

⚠️ **Важливо:** Файли сертифікатів (`certs/`) додані до `.gitignore` і НЕ повинні комітитись в репозиторій!

⚠️ Ці сертифікати призначені **тільки для локальної розробки**. Для production використовуйте Let's Encrypt або інші довірені CA.

## Посилання

- [mkcert на GitHub](https://github.com/FiloSottile/mkcert)
- [Документація nginx](https://nginx.org/en/docs/)
