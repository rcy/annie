-include ../.env

gen: schema.sql query.sql sqlc.yaml
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.25.0 generate

schema.sql:
	echo .schema | sqlite3 ${SQLITE_DB} > $@

clean:
	rm -f schema.sql
