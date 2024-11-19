lint:
	golangci-lint run ./...

fmt:
	gofmt -l -s .

gic:
	go run main.go
