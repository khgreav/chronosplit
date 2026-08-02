export DEBFULLNAME ?= Karel Hanák
export DEBEMAIL ?= khgreav@gmail.com

.PHONY: build clean update-deps run changelog deb-package

# dev
build:
	mkdir -p output
	go build -o output/app

clean:
	rm -rf output

update-deps:
	go get -u
	go mod tidy

run:
	go run .

# packaging
changelog:
	gbp dch -aS

deb-package:
	debuild -us -uc -b -tc
