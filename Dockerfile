# Стадия сборки: образ golang собирает статический бинарник под Linux.
FROM golang:1.26.2 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/scheduler .

# Финальный образ: только Alpine, исполняемый файл и каталог web.
FROM alpine:latest

WORKDIR /app

RUN mkdir -p /data

COPY --from=builder /out/scheduler /app/scheduler
COPY web /app/web

ENV TODO_PORT=7540 \
    TODO_DBFILE=/data/scheduler.db \
    TODO_PASSWORD="1234567!"

EXPOSE ${TODO_PORT}

VOLUME ["/data"]

CMD ["/app/scheduler"]
