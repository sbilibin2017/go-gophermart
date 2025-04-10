.PHONY: test test-cov migrate lint mockgen docker-up

test:
	go test ./...

test-cov:
	go test -coverprofile=coverage.out ./...
	

lint:
	staticcheck ./...

mockgen:
	@echo "Generating mock for: $(file)"
	@mockgen -source=$(file) \
		-destination=$(dir $(file))$(notdir $(basename $(file)))_mock.go \
		-package=$(shell basename $(dir $(file)))

migrate:
	goose postgres "user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) dbname=$(POSTGRES_DB) sslmode=disable" up

docker-up:
	docker run --name pg-gophermart -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_DB) -p $(POSTGRES_PORT):5432 -d postgres:latest
	
docker-stop:
	docker stop pg-gophermart
	docker rm pg-gophermart