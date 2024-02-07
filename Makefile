include .env

export POSTGRES_USER
export POSTGRES_PASSWORD
export POSTGRES_DB
export POSTGRES_HOST
export POSTGRES_PORT

DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

initdb:
	migrate -path migration -database "$(DB_URL)" -verbose up

dropdb:
	migrate -path migration -database "$(DB_URL)" -verbose down 1

test:
	go test ./...

bot:
	go run .

mocks:
	go generate ./...

.PHONY: initdb dropdb test bot mocks