.PHONY: migrate lint mockgen docker-run

export ACCRUAL_RUN_ADDRESS=:8081
export GOPHERMART_RUN_ADDRESS=:8080

export POSTGRES_USER=user
export POSTGRES_PASSWORD=password
export POSTGRES_DB=db
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432

run-accrual:	
	cd cmd/accrual && go build -o accrual && ./accrual \
		-a ${ACCRUAL_RUN_ADDRESS} \
		-d postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)

run-gophermart:	
	cd cmd/gophermart && go build -o gophermart && ./gophermart \
		-a ${GOPHERMART_RUN_ADDRESS} \
		-d postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB) \
		-r ${ACCRUAL_RUN_ADDRESS}

migrate:
	goose -dir ./migrations \
		postgres "user=$(POSTGRES_USER) dbname=$(POSTGRES_DB) password=$(POSTGRES_PASSWORD) host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) sslmode=disable" up
		

lint:
	staticcheck ./...

mockgen:	
	mockgen -source=$(file) \
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

test:
	go test ./internal/... -v

test-cov:
	go test ./internal/... -coverprofile=coverage.out > /dev/null
	go tool cover -func=coverage.out | grep total: | awk '{print "Total coverage:", $$3}'
	rm coverage.out
