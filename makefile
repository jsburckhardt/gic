lint:
	golangci-lint run ./...

fmt:
	gofmt -l -s -w .
	goimports -w .

fix: fmt
	golangci-lint run --fix ./...

gic:
	go run main.go
