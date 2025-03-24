FROM golang:1.23

WORKDIR /app

# Устанавливаем зависимости для OpenSSL
RUN apt-get update && apt-get install -y pkg-config libssl-dev

COPY go.mod go.sum ./
RUN go mod tidy

COPY init.lua /opt/tarantool/init.lua
COPY . .
RUN go build -o main ./cmd/kv-server

CMD ["/app/main"]