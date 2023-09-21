IMAGE_NAME=linkding-telegram

VERSION=latest

DOCKERFILE_PATH=.

lint:
	golangci-lint run
	
build: clean
	docker build -t $(IMAGE_NAME):$(VERSION) $(DOCKERFILE_PATH)

clean: stop
	docker rmi $(IMAGE_NAME):$(VERSION) || true

run: build
	docker compose up -d

stop:
	docker compose stop
	docker compose rm

unit-test:
	go test ./...

.PHONY: build clean run stop unit-test
