# MireaVstrechki
Сервис для управления рабочими встречами. Внутри каждой встречи — участники, задачи, повестка и email-уведомления.

**Стек:** Go 1.25, PostgreSQL, React + TypeScript

## Запуск

```bash
make db-up       # поднять postgres
make migrate     # накатить миграции
make seed        # залить тестовые данные
make run         # запустить backend 
make front-start # запустить frontend
```

Или всё сразу через Docker:

```bash
make up
```

## Конфигурация

| Переменная       | Дефолт                                                        |
| ---------------- | ------------------------------------------------------------- |
| `SERVER_PORT`    | `8080`                                                        |
| `POSTGRES_DSN`   | `postgres://postgres:postgres@localhost:5432/meeting_service?sslmode=disable` |
| `JWT_SECRET`     | `tut vasha reklama`                                           |
| `SMTP_HOST`      | —                                                             |
| `SMTP_PORT`      | `587`                                                         |
| `SMTP_USERNAME`  | —                                                             |
| `SMTP_PASSWORD`  | —                                                             |
| `SMTP_FROM`      | —                                                             |

Если SMTP не задан — email-уведомления просто не отправляются.

`YANDEX_GEO_KEY` — ключ Яндекс Геокодера. Если не задан — адреса показываются как текст без ссылки на карту.

## Структура

```
backend/
  cmd/api/           — точка входа
  internal/
    auth/            — JWT
    config/          — env-конфиг
    http/            — handlers, middleware, router
    model/           — модели
    notification/    — email: приглашения и напоминания
    repository/      — SQL
    service/         — бизнес-логика
  migrations/        — goose-миграции

frontend/
  src/
    api/             — REST-клиент
    pages/           — страницы
    components/      — переиспользуемые компоненты
```

## API

Swagger-доки: `make docs-serve` → `http://localhost:4000`

Авторизация — Bearer-токен, который возвращает `POST /login`.

## Деплой

Пуш тега `v*` запускает GitHub Actions: собирает бинари через GoReleaser, пушит Docker-образы в GHCR и деплоит на VPS по SSH.

Нужные секреты в репо: `VPS_HOST`, `VPS_USER`, `VPS_SSH_KEY`, `POSTGRES_DSN`, `JWT_SECRET`, `YANDEX_GEO_KEY`, SMTP-переменные.

## Роли

| Роль          | Что может                                              |
| ------------- | ------------------------------------------------------ |
| `admin`       | Всё                                                    |
| `organizer`   | Создавать встречи, управлять участниками и задачами    |
| `employee`    | Просматривать встречи, подтверждать/отклонять участие  |

## Тесты

```bash
make test   # unit-тесты
make fuzz   # фаззинг (10 секунд на каждый)
```

```bash
make load-test                                       # дефолт: admin@demo.local / password
make load-test LOGIN=sergey.lazarenko.0241@mail.ru PASSWORD=sereginrparolik   # свой пользователь
```
