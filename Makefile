export BUILDKIT_PROGRESS=plain

watch:
	air

fmt:
	go fmt main.go

lint:
	golangci-lint run

sql:
	. ./.env && sqlite3 ${SQLITE_DB}

test:
	set -a && . ./.env.test && go test ./...

build:
	docker build -t annie .

run:
	docker run --env-file=.env -e SQLITE_DB=:memory: annie
