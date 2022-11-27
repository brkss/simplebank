
PASS = root
USER = root
HOST = localhost
PORT = 5432
DB = simplebank

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it postgres12 dropdb simplebank

migrateup:
	migrate -path db/migration -database "postgres://$(USER):$(PASS)@$(HOST):$(PORT)/$(DB)?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://$(USER):$(PASS)@$(HOST):$(PORT)/$(DB)?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown
