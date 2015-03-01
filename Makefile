.PHONY: build doc fmt lint run test vet create generate serve

default: build

docs:
	@pandoc -s -w man -o lannister.1 README.md
	@godoc -html > docs/lannister.html

build: vet
	go build -v

vet:
	go vet

lint:
	golint

fmt:
	go fmt

test:
	go test

clean:
	rm -rf testsite

create: clean build
	./lannister createsite testsite

generate: create
	./lannister generate testsite

serve:
	livereloadx -s -p 4366 testsite/site
