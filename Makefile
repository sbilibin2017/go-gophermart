.PHONY: test test-cov migrate lint mockgen docker-up

test:
	go test ./...

test-cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tee coverage.txt
	grep -v '_mock' coverage.txt | grep -v 'main.go' | grep -v '^total:' > coverage_filtered.txt
	rm coverage.txt coverage.out
	mv coverage_filtered.txt coverage.txt

lint:
	staticcheck ./...

mockgen:
	@echo "Generating mock for: $(file)"
	@mockgen -source=$(file) \
		-destination=$(dir $(file))$(notdir $(basename $(file)))_mock.go \
		-package=$(shell basename $(dir $(file)))

migrate:
	cd ./migrations && goose postgres "user=user password=password host=localhost port=5432 dbname=db sslmode=disable" up

docker-up:
	docker run --name pg-gophermart -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=db -p 5432:5432 -d postgres:latest
	
docker-stop:
	docker stop pg-gophermart
	docker rm pg-gophermart