DB_NAME ?= beaconvalidators
DB_USER ?= your_default_database_user
DB_PASSWORD ?= your_default_database_password
DB_HOST ?= localhost
DB_PORT ?= 5432

build:
	@go build -o bin/beaconactor  -v

run: build
	@./bin/beaconactor

test:
	@go test -v ./... -count=1

create-db:
	psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME)"

drop-db:
	psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME)"

.PHONY: create-db drop-db build run test
