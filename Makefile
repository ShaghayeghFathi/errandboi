all: up run

build:
	go build errandboi

run:
	go run errandboi serve

up:
	docker-compose up -d

down:
	docker-compose down

lint:
	golangci-lint run
.PHONY: lint
