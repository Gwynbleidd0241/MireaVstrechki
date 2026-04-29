.PHONY: up down run migrate seed test test-integration fuzz lint front-start db-up docs docs-serve
DSN ?= postgres://postgres:postgres@localhost:5432/meeting_service?sslmode=disable

up:
	docker compose up --build

down:
	docker compose down

run:
	cd backend && go run ./cmd/api

migrate:
	cd backend && POSTGRES_DSN="$(DSN)" go run ./cmd/api migrate

seed:
	cd backend && POSTGRES_DSN="$(DSN)" go run ./cmd/api seed

test:
	cd backend && go test -v ./...

test-integration:
	cd backend && INTEGRATION_DSN="$(DSN)" go test -tags=integration -v ./internal/repository/postgres

fuzz:
	cd backend && go test -run=^$$ -fuzz=FuzzIsValidEmail -fuzztime=10s ./internal/validation
	cd backend && go test -run=^$$ -fuzz=FuzzParseToken -fuzztime=10s ./internal/auth
	cd backend && go test -run=^$$ -fuzz=FuzzIsValidTaskStatus -fuzztime=10s ./internal/service
	cd backend && go test -run=^$$ -fuzz=FuzzEventIDFromPath -fuzztime=10s ./internal/http/handlers

lint:
	cd backend && go vet ./...
	@cd backend && unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt: следующие файлы не отформатированы:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

front-start:
	cd frontend && npm start

db-up:
	docker compose up -d postgres

docs:
	cd backend && swag init -g cmd/api/main.go -o ../docs/swagger

docs-serve: docs
	docker run --rm -p 4000:4000 \
		-e SWAGGER_JSON=/swagger.json \
		-v $(CURDIR)/docs/swagger/swagger.json:/swagger.json \
		swaggerapi/swagger-ui
