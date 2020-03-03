.PHONY: build install test fmt start_backend stop_backend run

build:
	rm -rf .bin
	mkdir .bin
	cd .bin && go build ../src/github.com/ktraff/load_balancer/main.go

install:
	go install src/github.com/ktraff/load_balancer/main.go

test:
	go test ./...

fmt:
	go fmt ./...
	-git add .
	-git commit -m "formatting"

start_backend:
	cd backend && docker-compose up --detach --remove-orphans

stop_backend:
	cd backend && docker-compose down

run: build
	BACKEND_1=http://localhost:8000 \
	BACKEND_2=http://localhost:8001 \
	BACKEND_3=http://localhost:8002 \
	./.bin/main $(WORKERS) $(REQUESTS)

push: fmt
	git push