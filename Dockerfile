FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext nano musl-dev
# dependencies
COPY ["app/go.mod", "app/go.sum", "./"]

RUN go mod download

# build

COPY app ./
RUN go build -o ./bin/app cmd/app/main.go

FROM alpine AS runner

COPY --from=builder usr/local/src/bin/app  ./app
RUN chmod +x ./app
EXPOSE 8000
COPY --from=builder /usr/local/src/internal/service/database ./internal/service/database
CMD ["./app"]