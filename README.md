# Alice Chemistry Coach

Тренажёр романтической коммуникации с AI-персонажем.  
Стек: **Go + Gin + PostgreSQL** (бэкенд) · **Next.js 16 + React 19 + Tailwind 4** (фронтенд).

---

## Быстрый старт (Docker Compose)

### Требования
- Docker ≥ 24 и Docker Compose v2 (`docker compose version`)
- Свободные порты: **3032** (фронтенд), **8080** (API), **5432** (PostgreSQL)

### 1. Клонировать / распаковать репозиторий

```bash
git clone <repo-url> chemistry-coach
cd chemistry-coach
```

### 2. Создать `.env` (опционально — для Yandex AI)

```bash
cp .env.example .env
# Заполните YANDEX_AI_API_KEY и YANDEX_FOLDER_ID если есть ключи.
# Без ключей приложение работает на mock-LLM (встроенные ответы).
```

### 3. Запустить все сервисы

```bash
docker compose up --build
```

Первый запуск занимает 2–5 минут (сборка образов).  
После старта:

| Сервис | URL |
|--------|-----|
| Фронтенд | http://localhost:3032 |
| API | http://localhost:8080 |
| Swagger UI | http://localhost:8080/swagger/index.html |

### 4. Остановить

```bash
docker compose down          # остановить контейнеры
docker compose down -v       # + удалить данные PostgreSQL
```

---

## Локальный запуск (без Docker)

### Бэкенд

```bash
# Требования: Go ≥ 1.23, PostgreSQL 15

# 1. Создать БД
createdb acc_db

# 2. Настроить окружение
cp .env.example .env
# Установить DATABASE_URL=postgresql://postgres:password@localhost:5432/acc_db?sslmode=disable

# 3. Запустить
go run ./cmd/server
# Сервер слушает :8080
```

### Фронтенд

```bash
cd frontend

# 1. Установить зависимости
npm install

# 2. Указать URL API (если бэкенд не на :8080)
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local

# 3. Dev-режим
npm run dev
# http://localhost:3032

# 4. Или production-сборка
npm run build && npm start
```

---

## Переменные окружения

### Бэкенд (`docker-compose.yml` / `.env`)

| Переменная | По умолчанию | Описание |
|------------|-------------|----------|
| `DATABASE_URL` | `postgresql://postgres:password@db:5432/acc_db?sslmode=disable` | Строка подключения к PostgreSQL |
| `PORT` | `8080` | Порт HTTP-сервера |
| `NODE_ENV` | `development` | `development` включает Swagger и verbose-логи |
| `YANDEX_AI_API_KEY` | _(пусто)_ | API-ключ Yandex AI Studio. Без него — mock-режим |
| `YANDEX_AI_MODEL` | `yandexgpt` | Модель Yandex GPT |
| `YANDEX_FOLDER_ID` | _(пусто)_ | Folder ID Yandex Cloud |
| `CORS_ALLOWED_ORIGINS` | `*` | Разрешённые CORS-origins (через запятую) |

### Фронтенд (build arg / `.env.local`)

| Переменная | По умолчанию | Описание |
|------------|-------------|----------|
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080` | Базовый URL бэкенда. **Важно:** значение запекается в JS-бандл при сборке — при изменении нужен `docker compose build frontend` |

---

## Архитектура

```
chemistry-coach-yandex/
├── cmd/server/          # точка входа Go
├── internal/
│   ├── catalog/         # справочник целей и персонажей
│   ├── delivery/http/   # Gin-роутер, хэндлеры, middleware (CORS, auth)
│   ├── domain/          # сущности и интерфейсы репозиториев
│   ├── infrastructure/  # PostgreSQL-репозитории + Yandex AI / mock-LLM
│   └── usecase/         # бизнес-логика (auth, session, profile, catalog)
├── migrations/          # SQL-миграции (применяются автоматически при старте)
├── frontend/            # Next.js приложение
│   ├── app/             # страницы (onboarding, goal, persona, chat, debrief, history)
│   ├── components/ui/   # дизайн-система (shadcn-based)
│   ├── lib/api.ts       # HTTP-клиент для всех эндпоинтов бэкенда
│   └── Dockerfile       # multi-stage сборка (deps → builder → runner)
├── Dockerfile           # multi-stage сборка бэкенда
└── docker-compose.yml   # db + api + frontend
```

### Пользовательский сценарий

```
/ (онбординг) → /goal → /persona → /chat → /debrief → /history
```

1. **Онбординг** — год рождения + 3 согласия → `POST /api/v1/auth/start` → `userId` в `localStorage`
2. **Цель** — выбор сценария → `GET /api/v1/goals`
3. **Персонаж** — выбор AI-собеседника → `GET /api/v1/personas`
4. **Чат** — переписка → `POST /sessions`, `POST /sessions/:id/messages`, `POST /sessions/:id/suggest`
5. **Разбор** — итоги сессии → `POST /sessions/:id/finish` → `GET /sessions/:id`
6. **История** — список прошлых сессий → `GET /sessions`

---

## API

Swagger UI доступен по адресу `http://localhost:8080/swagger/index.html` при `NODE_ENV=development`.

Все защищённые эндпоинты требуют заголовок `X-User-Id: <userId>`.

---

## Mock-режим LLM

Если `YANDEX_AI_API_KEY` не задан, бэкенд автоматически использует встроенный mock:
- Персонаж отвечает заготовленными фразами
- Анализ сообщений и итоговый разбор генерируются детерминированно
- Подходит для демонстрации без ключей Yandex AI
