FROM golang:1.23.11 AS builder

WORKDIR /app

copy ../go.mod ./
copy ../go.sum ./
run go mod downloader

copy .. .

RUN go build -o app ./app/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

COPY internal/service/database/database.db ./internal/service/database/database.db

EXPOSE 8000

CMD ["./app"]