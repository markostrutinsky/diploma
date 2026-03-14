# Millog Frontend

Фронтенд системи військової логістики Millog.

## Запуск

```bash
npm install
npm run dev
```

Додаток буде доступний на http://localhost:5173

## Збірка

```bash
npm run build
```

## Сторінки

- **/** — головна
- **/admin/users** — додавання користувачів (адмін)
- **/setup-password** — встановлення пароля за invite-посиланням (`?token=...`)

## API

Фронт працює через Vite proxy: запити на `/api/*` перенаправляються на бекенд (localhost:8080).
