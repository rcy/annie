export BUILDKIT_PROGRESS=plain

watch:
	. ./.env && air

fmt:
	go fmt main.go

lint:
	golangci-lint run

sql:
	set -a && . ./.env && sqlite3 ${SQLITE_DB}

test:
	set -a && . ./.env.test && go test ./...

build:
	docker build -t annie .

run:
	docker run --env-file=.env -e SQLITE_DB=:memory: annie
