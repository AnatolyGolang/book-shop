.PHONY: dc run test lint

GOEXE :=
ifeq ($(OS),Windows_NT)
    GOEXE := .exe
endif

dc:
	@docker-compose up  --remove-orphans --build

build:
	@go build -race -o app$(GOEXE) cmd/main.go

run:
	@go build -race -o app cmd/main.go
ifeq ($(OS),Windows_NT)
	@cmd /c app.exe
else
	@./app
endif

test:
	@go test -race ./...

install-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2

lint:
	@golangci-lint run ./...
