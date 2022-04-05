all: up run

build:
	go build github.com/ShaghayeghFathi/errandboi

run:
	go run github.com/ShaghayeghFathi/errandboi serve

up:
	docker-compose up -d

down:
	docker-compose down

lint:
	golangci-lint run
.PHONY: lint
