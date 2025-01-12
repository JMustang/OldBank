postgres:
	docker run --name postgres17 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root OldBank
	
dropdb:
	docker exec -it postgres17 dropdb OldBank

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5433/OldBank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5433/OldBank?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown