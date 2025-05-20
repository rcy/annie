-include ../.env

schema.sql:
	echo .schema | sqlite3 ${SQLITE_DB} > $@

clean:
	rm -f schema.sql
