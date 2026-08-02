# Copyright (c) 2026 Karel Hanák
# SPDX-License-Identifier: MIT
.PHONY: build clean update-deps run check changelog deb-package

# dev
build: update-deps
	mkdir -p bin && CGO_ENABLED=0 go build -trimpath -o bin/chronosplit .

clean:
	rm -rf bin

update-deps:
	go get -u
	go mod tidy

run:
	go run .

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

deb-package:
	debuild -us -uc -b -tc
