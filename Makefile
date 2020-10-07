.PHONY: generate build release clean test lint dep-update

generate:
	GOFLAGS=-mod=vendor go generate ./...

build:
	goreleaser build --rm-dist --snapshot

release:
	goreleaser release --rm-dist --snapshot --skip-publish

clean:
	rm -f dist/

test:
	export GOFLAGS=-mod=vendor ; \
	go test -v -race -cover ./...

lint:
	golangci-lint run --enable-all --disable=godox

dep-update:
	go get -u ./...
	go test ./...
	go mod tidy
	go mod vendor
