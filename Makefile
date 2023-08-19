include .env

start:
	foreman start

fmt:
	go fmt main.go

lint:
	golangci-lint run

sql:
	sqlite3 ${SQLITE_DB}

test:
	SQLITE_DB=:memory: go test ./...
