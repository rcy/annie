
start:
	PORT=8080 SQLITE_DB=/tmp/annie.db IRC_NICK=PvxqcjQxgd4C9cmbv IRC_CHANNEL=#embx IRC_SERVER=irc.libera.chat:6697 go run main.go

fmt:
	go fmt main.go

deploy:
	flyctl deploy
