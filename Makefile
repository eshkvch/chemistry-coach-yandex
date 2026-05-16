.PHONY: run build swagger docker-up docker-down

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

swagger:
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down
