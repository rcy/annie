start:
	PORT=8080 go run main.go

fmt:
	go fmt main.go

deploy:
	flyctl deploy
