.PHONY: build-run build run-exe run-go migrate-up migrate_down run_tests

build-run:
	go build -o shop cmd/app/main.go
	./shop

run-tests:
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/merch_shop_test?sslmode=disable" -path ./migrations up
	go test -v ./internal/service > ./tests/service_tests.log 2>&1
	go test -v ./tests/e2e > ./tests/e2e_tests.log 2>&1

run-exe:
	./shop

run-go:
	go run cmd/app/main.go

migrate-up:
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/merch_shop?sslmode=disable" -path ./migrations up

migrate-down:
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/merch_shop?sslmode=disable" -path ./migrations down
