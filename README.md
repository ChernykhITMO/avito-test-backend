# Room Booking Service

Реализация тестового задания на стажировку в Авито.

Условие: [TASK.md](TASK.md)  
API: [api.yaml](api.yaml)

---

## Запуск

    cp .env.example .env
    docker compose up --build

Сервис доступен:  
http://localhost:8080

Swagger UI:  
http://localhost:8080/swagger/index.html

---

## Стек

- Go 1.25
- PostgreSQL 16
- pgx/v5
- goose
- docker compose

---

## Реализовано

- JWT авторизация (`dummyLogin`, `register`, `login`)
- CRUD переговорок (admin)
- immutable расписание
- генерация слотов
- бронирование и отмена
- список броней (admin / user)
- mock conference service
- swagger docs
- seed данных

---

## Тесты

    make test
    make smoke ENV_FILE=.env

Покрытие: ~43%

---

## Основные команды

    make up ENV_FILE=.env
    make down
    make seed
    make migrate-up
    make migrate-down
    make lint
    make swagger-gen

---

## Замечания

- слоты генерируются заранее + лениво догенерируются
- конкурентность обеспечивается через уникальные индексы и транзакции
- проверочные тесты прошли на 5/5
