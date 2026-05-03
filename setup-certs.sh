#!/bin/bash

# Скрипт для налаштування локальних HTTPS сертифікатів за допомогою mkcert

set -e

echo "🔐 Налаштування локальних HTTPS сертифікатів..."

# Перевірка чи встановлено mkcert
if ! command -v mkcert &> /dev/null; then
    echo "❌ mkcert не встановлено!"
    echo ""
    echo "Встановіть mkcert:"
    echo "  Ubuntu/Debian: sudo apt install libnss3-tools && wget https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-linux-amd64 && sudo mv mkcert-v1.4.4-linux-amd64 /usr/local/bin/mkcert && sudo chmod +x /usr/local/bin/mkcert"
    echo "  Arch Linux: sudo pacman -S mkcert"
    echo "  Fedora: sudo dnf install mkcert"
    echo ""
    echo "Або використайте офіційну інструкцію: https://github.com/FiloSottile/mkcert#installation"
    exit 1
fi

# Створення директорії для сертифікатів
mkdir -p certs

# Встановлення локального CA (Certificate Authority)
echo "📝 Встановлюємо локальний Certificate Authority..."
mkcert -install

# Генерація сертифікатів для localhost
echo "🔑 Генеруємо сертифікати для localhost..."
cd certs
mkcert localhost 127.0.0.1 ::1

# Перейменування файлів для зручності
mv localhost+2.pem localhost.pem
mv localhost+2-key.pem localhost-key.pem

echo ""
echo "✅ Сертифікати успішно створено!"
echo "📁 Файли знаходяться в директорії ./certs/"
echo ""
echo "Тепер можете запустити проект:"
echo "  docker-compose up --build"
echo ""
echo "Сайт буде доступний за адресою: https://localhost"
