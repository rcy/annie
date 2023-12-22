include .env

watch:
	air

fmt:
	go fmt main.go

lint:
	golangci-lint run

sql:
	sqlite3 ${SQLITE_DB}

test:
	SQLITE_DB=:memory: go test ./...

build:
	docker build -t annie .

run:
	docker run --env-file=.env -e SQLITE_DB=:memory: annie

