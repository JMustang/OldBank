postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
createdb:
	docker exec -it postgres createdb --username=root --owner=root OldBank
dropdb:
	docker exec -it postgres dropdb OldBank

migrateup:
	migrate -path db/migration -database "postgresql://docker:docker@localhost:5433/OldBank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://docker:docker@localhost:5433/OldBank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc