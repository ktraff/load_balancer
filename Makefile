.PHONY: build test run

build:
	rm -rf .bin
	mkdir .bin
	cd .bin && go build ../src/github.com/ktraff/load_balancer/main.go

install:
	go install src/github.com/ktraff/load_balancer/main.go

test:
	go test ./...

run: build
	BACKEND_1=http://localhost:8000 \
	BACKEND_2=http://localhost:8001 \
	BACKEND_3=http://localhost:8002 \
	./.bin/main