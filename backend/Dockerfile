FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /out/api ./cmd/api
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.27.1

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /out/api /app/api
COPY --from=builder /go/bin/goose /app/goose
COPY migrations /app/migrations

CMD ["/app/api"]
