include .env

postgres:
	docker run --name postgres_chatApp -p ${DATABASE_PORT}:${DATABASE_PORT} -e POSTGRES_USER=${DATABASE_USER} -e POSTGRES_PASSWORD=${DATABASE_PASSWORD} -d postgres:16-alpine3.19

createdb:
	docker exec -it postgres_chatApp createdb --username=${DATABASE_USER} --owner=${DATABASE_USER} ${DATABASE_NAME}

dropdb:
	docker exec -it postgres_chatApp dropdb ${DATABASE_NAME}

migrateup:
	migrate -path migration -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose up

migratedown:
	migrate -path migration -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose down

server:
	go run cmd/main.go

.PHONY : server