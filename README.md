# Room Booking Service

Сервис бронирования переговорок на Go с JWT-авторизацией, PostgreSQL и `goose`-миграциями.

Реализованы обязательные сценарии из `TASK.md` и `api.yaml`:

- `dummyLogin` для ролей `admin` и `user`
- создание и просмотр переговорок
- создание immutable-расписания для комнаты
- получение доступных слотов на дату
- создание и отмена бронирования
- список всех броней для `admin`
- список будущих броней пользователя

Дополнительные сценарии:

- `register/login` по email/паролю с JWT
- `createConferenceLink` через mock `Conference Service`
- `make seed` для наполнения БД тестовыми данными
- `swagger-gen` для генерации Swagger-артефактов из аннотаций в коде

## Стек

- Go `1.25`
- PostgreSQL `16`
- `pgx/v5` + `pgxpool`
- `goose` для миграций
- `golangci-lint`
- `docker compose` для локального запуска
- отдельный smoke-контур на реальном PostgreSQL

## Архитектура

Сервис собран как `modular monolith` с feature-first структурой:

```text
internal/
  auth/
  room/
  schedule/
  slot/
  booking/
  platform/
  worker/
```

Внутри каждого модуля используются папки:

- `handler` — HTTP-слой, валидация запроса, mapping ошибок в API response
- `service` — use-case логика, транзакции, orchestration
- `repository` — SQL и доступ к PostgreSQL
- `model` — доменные сущности, enum'ы и ошибки

Общие инфраструктурные зависимости вынесены в `internal/platform`:

- `postgres` — DB pool, тонкий transaction manager, запуск миграций
- `jwt` — выпуск и валидация JWT
- `middleware` — auth/logging/recover
- `httpcommon` — JSON bind/render + error responses
- `password` — bcrypt для optional `register/login`

## Основные доменные сущности

- `User` — пользователь с ролью `admin` или `user`
- `Room` — переговорка
- `Schedule` — immutable расписание комнаты по дням недели и окну времени
- `Slot` — 30-минутный слот, заранее сгенерированный для конкретной даты
- `Booking` — бронь слота

Связи:

- `Room 1:1 Schedule`
- `Room 1:N Slot`
- `User 1:N Booking`
- `Slot 1:N Booking`, но активная бронь на слот может быть только одна

## Стратегия работы со слотами

Используется гибридный подход:

- при создании расписания слоты генерируются сразу на окно `SLOT_GENERATION_WINDOW_DAYS` вперёд
- при чтении `/rooms/{roomId}/slots/list` сервис лениво догенерирует слоты, если запрошенная дата вышла за текущий горизонт
- фоновый worker периодически подтягивает `generated_until`, чтобы горячий endpoint оставался fast-path

Почему так:

- `slotId` стабилен и хранится в БД
- endpoint слотов работает как дешёвое indexed read
- не нужна тяжёлая on-the-fly генерация на каждый запрос
- нет отдельного кеша и проблем с invalidation

## Конкурентность и консистентность

- бронирование защищено partial unique index:
  `bookings(slot_id) WHERE status = 'active'`
- создание расписания защищено `UNIQUE(room_id)` на `schedules`
- отмена брони выполняется в транзакции с `SELECT ... FOR UPDATE`
- генерация слотов идемпотентна через
  `ON CONFLICT (room_id, start_at) DO NOTHING`
- `generated_until` обновляется через `GREATEST(...)`

Это позволяет переживать параллельные запросы без `SERIALIZABLE` и без отдельного lock-service.

## Схема БД

Миграции лежат в [`migrations/`](./migrations):

- `00001_init.sql` — таблицы `users`, `rooms`, `schedules`, `slots`, `bookings`
- `00002_indexes.sql` — индексы под горячие запросы и ограничения консистентности
- `00003_dummy_users.sql` — тестовые пользователи под `dummyLogin`

Бизнес-правила по `role` и `status` валидируются на уровне приложения. В базе для них не используются enum-like `CHECK`-ограничения, чтобы не размазывать прикладную логику между слоями.

Ключевые индексы:

- `slots (room_id, slot_date, start_at)`
- `UNIQUE slots (room_id, start_at)`
- `UNIQUE bookings (slot_id) WHERE status = 'active'`
- `bookings (user_id, created_at DESC)`
- `bookings (created_at DESC, id DESC)`

## Авторизация

JWT кладётся в `Authorization: Bearer <token>`.

Роли:

- `admin` — создание комнат и расписаний, просмотр всех броней
- `user` — просмотр комнат/слотов, создание и отмена своих броней, просмотр своих будущих броней

Для обязательной части можно использовать `POST /dummyLogin`.

Дополнительно реализованы:

- `POST /register`
- `POST /login`

## Conference Link (Mock)

Для `POST /bookings/create` поле `createConferenceLink=true` вызывает mock-провайдер конференций:

- реализация: `internal/conference/mock`
- интерфейс интеграции: `internal/booking/service.ConferenceService`
- настройка через env:
  - `CONFERENCE_BASE_URL` (по умолчанию `https://conference.local/booking/`)
  - `CONFERENCE_MOCK_MODE` (`ok`, `create_unavailable`, `cancel_unavailable`)

Принятые решения по сбоям:

- если `Conference Service` недоступен при создании ссылки (`create_unavailable`), бронь всё равно создаётся, но без `conferenceLink`;
- если ссылка успешно создана, но запись брони в БД не удалась, выполняется best-effort компенсация: `CancelBookingLink`;
- если компенсация недоступна (`cancel_unavailable`), запрос на бронь возвращает исходную ошибку БД, а риск orphan-link фиксируется как технический долг (для production-версии рекомендуется outbox/async compensator).

## Запуск локально

Секреты и локальные параметры вынесены в [`.env.example`](./.env.example).

Перед запуском:

```bash
cp .env.example .env
```

### Вариант 1. Через Docker

```bash
docker compose up --build
```

Либо через `Makefile`:

```bash
make up ENV_FILE=.env
```

Сервис будет доступен на `http://localhost:8080`.
Swagger UI: `http://localhost:8080/swagger/index.html`.

### Вариант 2. Локально через Go

1. Поднять PostgreSQL.
2. Экспортировать переменные из `.env`.
3. Запустить сервис:

```bash
set -a && source .env && set +a
make run
```

Миграции запускаются автоматически на старте приложения.

## Полезные команды

```bash
make test
make lint
make swagger-gen
make smoke ENV_FILE=.env
make migrate-status
make migrate-up
make migrate-down
DATABASE_URL=postgres://postgres:postgres@localhost:5432/room_booking?sslmode=disable make seed
```

Для `goose` теперь нужно явно передать `DATABASE_URL` из окружения:

```bash
set -a && source .env && set +a
make migrate-up
```

## Основные переменные окружения

```bash
APP_PORT=8080
APP_LOG_LEVEL=info
DATABASE_URL=postgres://postgres:postgres@localhost:5432/room_booking?sslmode=disable
DB_MAX_CONNS=10
DB_MIN_CONNS=2
JWT_SECRET=dev-secret
JWT_TTL_SECONDS=86400
SLOT_GENERATION_WINDOW_DAYS=30
SLOT_REFILL_INTERVAL_SECONDS=300
CONFERENCE_BASE_URL=https://conference.local/booking/
CONFERENCE_MOCK_MODE=ok
POSTGRES_DB=room_booking
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_PORT=5432
```

## Swagger Codegen

Генерация Swagger-артефактов из аннотаций:

```bash
make swagger-gen
```

Результат:

- `docs/swagger/swagger.json`
- `docs/swagger/swagger.yaml`
- `docs/swagger/docs.go`

## Быстрый smoke flow

### 1. Получить токен администратора

```bash
curl -s localhost:8080/dummyLogin \
  -H 'Content-Type: application/json' \
  -d '{"role":"admin"}'
```

### 2. Создать переговорку

```bash
curl -s localhost:8080/rooms/create \
  -H "Authorization: Bearer <admin-token>" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Blue","capacity":8}'
```

### 3. Создать расписание

```bash
curl -s localhost:8080/rooms/<room-id>/schedule/create \
  -H "Authorization: Bearer <admin-token>" \
  -H 'Content-Type: application/json' \
  -d '{"daysOfWeek":[1,2,3,4,5],"startTime":"09:00","endTime":"18:00"}'
```

### 4. Получить токен пользователя

```bash
curl -s localhost:8080/dummyLogin \
  -H 'Content-Type: application/json' \
  -d '{"role":"user"}'
```

### 5. Получить свободные слоты

```bash
curl -s "localhost:8080/rooms/<room-id>/slots/list?date=2026-03-23" \
  -H "Authorization: Bearer <user-token>"
```

### 6. Создать бронь

```bash
curl -s localhost:8080/bookings/create \
  -H "Authorization: Bearer <user-token>" \
  -H 'Content-Type: application/json' \
  -d '{"slotId":"<slot-id>","createConferenceLink":true}'
```

## Тесты и проверочный контур

В проекте есть:

- unit tests для `service`, `handler`, `repository`
- e2e-сценарии в [`tests/e2e/e2e_test.go`](./tests/e2e/e2e_test.go)
- smoke-сценарий в [`tests/smoke/smoke_test.go`](./tests/smoke/smoke_test.go), который запускается через [`docker-compose.e2e.yaml`](./docker-compose.e2e.yaml) и ходит в реальный HTTP-сервис поверх настоящего PostgreSQL

Проверенные команды:

```bash
GOCACHE=/tmp/gocache go test ./...
GOCACHE=/tmp/gocache go test ./... -coverprofile=coverage.out
make lint
```

Проверочный контур:

```bash
make smoke ENV_FILE=.env
make smoke-down ENV_FILE=.env
```

Итоговое покрытие по `coverage.out`: `43.0%`.

## Нагрузочное тестирование

Нагрузочный прогон выполнен **23 марта 2026** для горячего endpoint:

- endpoint: `GET /rooms/{roomId}/slots/list?date=...`
- инструмент: `hey`
- длительность: `30s`
- целевая нагрузка: `~100 RPS` (`-q 5 -c 20`)
- окружение: локальный Docker Compose (`app + postgres`)

Ключевые метрики:

- фактический `RPS`: `99.95`
- `p95`: `14.6 ms`
- `p99`: `17.6 ms`
- `error rate`: `0%` (`3000` ответов с `200`)

Вывод: на текущем локальном стенде endpoint укладывается в ориентир SLI `200 ms`.

## Конфигурация линтера

Конфигурация `golangci-lint` находится в корне проекта:

- `.golangci.yaml`

Включены линтеры: `errcheck`, `govet`, `ineffassign`, `misspell`, `revive`, `staticcheck`, `unused`.

## Что можно улучшить дальше

- перейти с offset-pagination на keyset для `GET /bookings/list`
- вынести conference-link integration в async job/outbox
- добавить OpenAPI codegen для DTO и server stubs
- добавить интеграционные тесты с реальным PostgreSQL в CI
- добавить refresh tokens / key rotation для JWT
