
PASS = root
USER = root
HOST = localhost
PORT = 5432
DB = simplebank

DB_CONTAINER = postgres12

db:
	docker start $(DB_CONTAINER)

postgres:
	docker run --name $(DB_CONTAINER) -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it postgres12 dropdb simplebank

migrateup:
	migrate -path db/migration -database "postgres://$(USER):$(PASS)@$(HOST):$(PORT)/$(DB)?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://$(USER):$(PASS)@$(HOST):$(PORT)/$(DB)?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc db test
