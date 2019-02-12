all: test

test:
	go test -v ./...
	go generate ./example
