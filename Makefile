.PHONY: build update-deps

update-deps:
	go get -u
	go mod tidy

build:
	mkdir -p output
	go build -o output/app

clean:
	rm -rf output

run:
	go run .
