.PHONY: lint fmt test

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

test:
	go test ./...
