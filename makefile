lint:
	golangci-lint run ./...

fmt:
	gofmt -l -s .
