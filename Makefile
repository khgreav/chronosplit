# Copyright (c) 2026 Karel Hanák
# SPDX-License-Identifier: MIT
.PHONY: build clean update-deps run check changelog deb-package

VERSION=dev
BUILD_DATE=$(shell date -u +'%Y-%m-%d')
MODULE_PATH=github.com/khgreav/chronosplit

LDFLAGS := -ldflags="-X $(MODULE_PATH)/version.Version=$(VERSION) -X $(MODULE_PATH)/version.BuildDate=$(BUILD_DATE)"

# dev
build: update-deps
	mkdir -p bin && CGO_ENABLED=0 go build $(LDFLAGS) -o bin/chronosplit .

clean:
	rm -rf bin chronosplit.db

update-deps:
	go get -u && go mod tidy
	go mod vendor

run:
	CHRONOSPLIT_DB="./chronosplit.db" go run $(LDFLAGS) .

# analysis
vet:
	go vet ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	golangci-lint run ./...

vulncheck:
	@command -v govulncheck >/dev/null 2>&1 || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

check: vet lint vulncheck

# packaging
changelog:
	gbp dch -aS

deb-package: update-deps
	debuild -us -uc -b -tc
