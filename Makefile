start:
	foreman start

fmt:
	go fmt main.go

deploy:
	flyctl deploy

lint:
	golangci-lint run
