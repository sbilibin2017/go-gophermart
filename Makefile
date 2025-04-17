.PHONY: migrate lint mockgen docker-run

export ACCRUAL_RUN_ADDRESS=:8080
export GOPHERMART_RUN_ADDRESS=:8081

export POSTGRES_USER=user
export POSTGRES_PASSWORD=password
export POSTGRES_DB=db
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432

run-accrual:
	go run cmd/accrual/main.go \
		-a ${ACCRUAL_RUN_ADDRESS} \
		-d postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)

run-gophermart:
	go run cmd/gophermart/main.go \
		-a ${GOPHERMART_RUN_ADDRESS} \
		-d postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB) \
		-r ${ACCRUAL_RUN_ADDRESS}

migrate:
	goose -dir ./migrations \
		postgres "user=$(POSTGRES_USER) dbname=$(POSTGRES_DB) password=$(POSTGRES_PASSWORD) host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) sslmode=disable" up
		

lint:
	staticcheck ./...

mockgen:
	@echo "Generating mock for: $(file)"
	@mockgen -source=$(file) \
		-destination=$(dir $(file))$(notdir $(basename $(file)))_mock.go \
		-package=$(shell basename $(dir $(file)))

docker-run:
	docker run -d \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p $(POSTGRES_PORT):5432 \
		--name postgres-gopohermart \
		postgres:16