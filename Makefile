# include variables from .envrc file
include .envrc
# ================================ #
# DEVELOPMENT
# ================================ #
.PHONY: run
run:
	go run ./cmd/api

# =============================== #
# MIGRATIONS
# =============================== #
.PHONY: migrate/create
migrate/create:
	migrate create -seq -ext=.sql -dir=./migrations create_urls_table

.PHONY: migrate/db 
migrate/db:
	psql --host=localhost --dbname=urlshortening --username=urlshortening 

.PHONY: migrate/up
migrate/up:
	migrate -path ./migrations -database ${URLSHORTENING_DB_DSN} up

.PHONY: migrate/down
migrate/down:
	migrate -path ./migrations -database ${URLSHORTENING_DB_DSN} down

.PHONY: migrate/version
migrate/version:
	migrate -path ./migrations -database ${URLSHORTENING_DB_DSN} version

.PHONY: migrate/force
migrate/force:
	migrate -path ./migrations -database ${URLSHORTENING_DB_DSN} force 1


# ========================= #
# BUILD
# ========================= #
.PHONY: build 
build:
	go build -o shortner ./cmd/api
