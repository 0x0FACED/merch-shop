.PHONY: build-run build run run-go migrate-up

build-run:
	go build -o shop cmd/app/main.go
	./shop

run:
	./shop

run-go:
	go run cmd/app/main.go

migrate-up:
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/merch_shop?sslmode=disable" -path ./migrations up
