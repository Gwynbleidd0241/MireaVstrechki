.PHONY: up down run migrate seed test test-integration fuzz front-start db-up docs docs-serve load-test
DSN ?= postgres://postgres:postgres@localhost:5432/meeting_service?sslmode=disable
LOGIN    ?= admin@demo.local
PASSWORD ?= password

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

load-test:
	k6 run -e BASE_URL=http://localhost:8080 -e LOGIN=$(LOGIN) -e PASSWORD=$(PASSWORD) load-test/script.js

front-start:
	cd frontend && npm start

db-up:
	docker compose up -d postgres

docs:
	cd backend && swag init -g cmd/api/main.go -o ../docs/swagger

docs-serve: docs
	@echo "Swagger UI: http://localhost:4000"
	docker run --rm -p 4000:8080 \
		-e SWAGGER_JSON=/swagger.json \
		-v $(CURDIR)/docs/swagger/swagger.json:/swagger.json \
		swaggerapi/swagger-ui
