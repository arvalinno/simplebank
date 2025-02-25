DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres16.3 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16.3-alpine3.20

createdb:
	docker exec -it postgres16.3 createdb --username=root --owner=root simple_bank 

execpsql:
	docker exec -it postgres16.3 psql -U root -d simple_bank

dropdb:
	docker exec -it postgres16.3 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc: 
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/arvalinno/simplebank/db/sqlc Store

dockerbuild:
	docker build -t simplebank:latest .

createnewdbmigration:
	migrate create -ext sql -dir db/migration -seq add_sessions

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge,merge_file_name=simple_bank\
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: postgres createdb execpsql dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1 dockerbuild db_docs db_schema proto evans