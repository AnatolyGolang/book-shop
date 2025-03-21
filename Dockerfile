FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod +x cmd/wait-for-it.sh
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o app ./cmd/main.go

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache bash dos2unix

RUN adduser -D -g '' appuser

COPY --from=builder /app/app .
COPY --from=builder /app/cmd/wait-for-it.sh /app/wait-for-it.sh
COPY --from=builder /app/internal/app/migrations /app/internal/app/migrations
COPY --from=builder /app/config/local.env /app/config/local.env

RUN dos2unix /app/config/local.env && chmod 644 /app/config/local.env
RUN dos2unix /app/wait-for-it.sh && chmod +x /app/wait-for-it.sh

USER appuser

EXPOSE 8080
CMD ["bash", "/app/wait-for-it.sh", "postgres:5432", "--timeout=60", "--", "./app"]