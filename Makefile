.PHONY: proto proto-lint proto-generate test test-race fmt vet

proto-lint:
	buf lint

proto-generate:
	buf generate

proto: proto-lint proto-generate

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	gofmt -w .

vet:
	go vet ./...