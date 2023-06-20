include .env

start:
	foreman start

fmt:
	go fmt main.go

deploy:
	flyctl deploy --build-arg commit=$(shell git rev-parse HEAD)

lint:
	golangci-lint run

sql:
	sqlite3 ${SQLITE_DB}
