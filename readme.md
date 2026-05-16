# Alice Chemistry Coach — Backend (Go)

MVP API для AI-тренажёра романтической коммуникации. Стек: **Go**, **Gin**, **GORM**, **PostgreSQL**, **Yandex AI Studio**.

## Быстрый старт

```bash
cp .env.example .env
# при необходимости укажите YANDEX_AI_API_KEY и YANDEX_FOLDER_ID

docker compose up --build
```

API: `http://localhost:8080/api/v1`  
Swagger: `http://localhost:8080/swagger/index.html`  
Health: `http://localhost:8080/health`

Без ключей Yandex сервис стартует с **mock LLM** (детерминированные ответы для разработки фронтенда).

## Локальная разработка

```bash
docker compose up -d db
go run ./cmd/server
```

## Swagger

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o docs
```

## Аутентификация (MVP)

Заголовок `X-User-Id` (получается из `POST /api/v1/auth/start`).

## Эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/v1/auth/start` | Онбординг |
| GET | `/api/v1/profile` | Профиль |
| GET | `/api/v1/goals` | Цели |
| GET | `/api/v1/personas` | Персоны |
| POST | `/api/v1/sessions` | Новая сессия |
| POST | `/api/v1/sessions/:id/messages` | Сообщение в чат |
| POST | `/api/v1/sessions/:id/suggest` | Подсказка |
| POST | `/api/v1/sessions/:id/finish` | Разбор |
| GET | `/api/v1/sessions` | История |
| GET | `/api/v1/sessions/:id` | Разбор сессии |

## Архитектура (clean)

- `domain` — сущности и интерфейсы
- `usecase` — бизнес-логика
- `infrastructure` — PostgreSQL, Yandex AI
- `delivery/http` — Gin handlers, middleware
