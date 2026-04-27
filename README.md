# Meeting Service

Сервис для управления рабочими мероприятиями. Каждая встреча — это не строчка в календаре, а контейнер для совместной работы: внутри неё ведутся чек-листы, задачи и список участников. Идея — собрать в одном месте то, для чего обычно используют отдельные сервисы.

## Стек

**Backend:** Go 1.25,PostgreSQL
**Frontend:** React + TypeScript

## Структура проекта

```
meeting-service/
├── backend/
│   ├── cmd/api/                  # точка входа сервиса
│   ├── internal/
│   │   ├── auth/                 # генерация и парсинг JWT
│   │   ├── config/               # загрузка конфигурации из env
│   │   ├── http/
│   │   │   ├── handlers/         # HTTP-обработчики
│   │   │   ├── middleware/       # auth, CORS
│   │   │   └── server.go         # роутинг и сборка зависимостей
│   │   ├── logger/               # zap-логгер
│   │   ├── model/                # модели
│   │   ├── repository/postgres/  # работа с PostgreSQL
│   │   └── service/              # бизнес-логика
│   ├── migrations/               # SQL-миграции для goose
│   └── go.mod
├── frontend/
│   └── src/
│       ├── api/                  # клиент REST API
│       ├── components/           # переиспользуемые компоненты
│       └── pages/                # страницы приложения
├── docker-compose.yml
└── README.md
```

Бэкенд организован по слоям: `handlers → service → repository → postgres`. Слой `service` содержит бизнес-правила, `handlers` отвечают только за транспорт, `repository` инкапсулирует SQL.

## Запуск для разработки

### 1. Запустить PostgreSQL

```bash
docker compose up -d
```

База поднимается на `localhost:5432`, пользователь и пароль — `postgres`, база — `meeting_service`.

### 2. Накатить миграции

```bash
cd backend
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/meeting_service?sslmode=disable" up
```

### 3. Запустить backend

```bash
cd backend
go run ./cmd/api
```

Сервис слушает `:8080`.

### 4. Запустить frontend

```bash
cd frontend
npm install
npm start
```

Дев-сервер откроется на `http://localhost:3000` и будет ходить в API на `http://localhost:8080`.

## Конфигурация

Все параметры передаются через переменные окружения.

| Переменная       | Назначение                       | Дефолт                                                                         |
| ---------------- | -------------------------------- | ------------------------------------------------------------------------------ |
| `SERVER_PORT`    | Порт HTTP-сервера                | `8080`                                                                         |
| `POSTGRES_DSN`   | DSN для PostgreSQL               | `postgres://postgres:postgres@localhost:5432/meeting_service?sslmode=disable`  |
| `JWT_SECRET`     | Секрет для подписи JWT           | `dev-secret` (в проде обязательно переопределить)                              |
