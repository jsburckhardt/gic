lint:
	golangci-lint run ./...

fmt:
	gofmt -l -w -s 
