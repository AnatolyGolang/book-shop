dc:
	docker-compose up --remove-orphans ---build

test:
	go test -race ./...

run_local:
	go run -race cmd/main.go -env=config/local.env

# Предполагается что у вас уже установлен golangci
lint:
	golangci-lint run