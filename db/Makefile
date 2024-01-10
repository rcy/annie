-include ../.env

gen: schema.sql query.sql sqlc.yaml
	sqlc generate

schema.sql:
	echo .schema | sqlite3 ${SQLITE_DB} > $@

clean:
	rm -f schema.sql
