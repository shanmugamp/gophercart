.PHONY: test test-race fmt vet build clean

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

build:
	go build ./...

clean:
	go clean