# Документация API

API описан через [swaggo/swag](https://github.com/swaggo/swag) — аннотации
лежат прямо над handler-ами в `backend/internal/http/handlers/*.go`.

## Сгенерировать документацию

```bash
make docs
```

Эта команда:

1. Установит `swag` CLI, если его ещё нет (`go install github.com/swaggo/swag/cmd/swag@latest`).
2. Проанализирует все аннотации и создаст три файла в `docs/swagger/`:
   - `swagger.json`
   - `swagger.yaml`
   - `docs.go` (Go-пакет, который можно встроить в бинарь для serving UI).

Файлы генерируются автоматически — **не редактируйте их руками**, правьте
только аннотации над хендлерами.

## Просмотр локально

### Самый быстрый способ — `make docs-serve`

```bash
make docs-serve
```

Команда сначала пере-генерирует документацию (`make docs`), затем поднимает
контейнер `swaggerapi/swagger-ui` на `http://localhost:8888`. Откройте
браузер — там будет полноценная интерактивная Swagger UI с возможностью
«Try it out» прямо по ручкам бэкенда.

Чтобы остановить — Ctrl+C, контейнер удалится автоматически (`--rm`).

### Альтернативы

**Swagger Editor онлайн.** Открыть https://editor.swagger.io/ →
File → Import file → выбрать `docs/swagger/swagger.yaml`. Подойдёт, если
не хочется ничего поднимать локально.

**VS Code.** Расширение `42Crunch.vscode-openapi` рендерит preview прямо
в редакторе при открытии `swagger.yaml`.

## Как добавлять новые ручки

В файле хендлера над методом написать комментарий вида:

```go
// MethodName godoc
// @Summary  Краткое описание
// @Tags     events
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id    path     int           true "Event ID"
// @Param    input body     createRequest true "Тело запроса"
// @Success  201   {object} response
// @Failure  400   {string} string "validation error"
// @Router   /events [post]
func (h *EventHandler) MethodName(w http.ResponseWriter, r *http.Request) { ... }
```

После этого пере-генерировать через `make docs`.

Полный справочник тегов: https://github.com/swaggo/swag#declarative-comments-format
