.PHONY: build test run

GOPATH=`pwd`

build:
	GOPATH=${GOPATH} go build ./...

install:
	GOPATH=${GOPATH} go install ./...

test:
	GOPATH=${GOPATH} go test ./...

run: build
	./bin/load_balancer