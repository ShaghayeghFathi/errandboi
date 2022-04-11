all: up lint run

build:
	go build github.com/ShaghayeghFathi/errandboi

docker-build:
	go docker build -t errandboi

run:
	go run github.com/ShaghayeghFathi/errandboi serve

up:
	docker-compose up -d

down:
	docker-compose down

lint:
	golangci-lint run
.PHONY: lint
