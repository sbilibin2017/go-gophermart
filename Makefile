.PHONY: drpg dspg m du

POSTGRES_USER ?= user
POSTGRES_PASSWORD ?= password
POSTGRES_DB ?= db
POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432

export POSTGRES_USER
export POSTGRES_PASSWORD
export POSTGRES_DB
export POSTGRES_HOST
export POSTGRES_PORT

DSN := user=$(POSTGRES_USER) dbname=$(POSTGRES_DB) password=$(POSTGRES_PASSWORD) host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) sslmode=disable

dr:
	$(MAKE) dspg
	$(MAKE) drpg	
	sleep 5	
	$(MAKE) m

m:
	goose -dir ./migrations postgres "$(DSN)" up

drpg:
	docker run -d \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p $(POSTGRES_PORT):5432 \
		--name postgres-gopohermart \
		postgres:16

dspg:
	-@docker stop postgres-gopohermart 2>/dev/null || true
	-@docker container prune -f
