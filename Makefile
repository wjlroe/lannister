.PHONY: build doc fmt lint run test vet

default: build

docs:
	@pandoc -s -w man -o lannister.1 README.md
	@godoc -html > docs/lannister.html

build: vet
	go build -v -o bin/lannister ./src/*.go

vet:
	go vet ./src/...

lint:
	golint ./src

fmt:
	go fmt ./src/...

test:
	go test ./src/...

run: build
	./bin/lannister
