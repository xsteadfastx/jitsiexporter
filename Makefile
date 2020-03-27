.PHONY: build clean test lint dep-update

VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
            echo v0)

build:
	export GOFLAGS=-mod=vendor
	go generate ./...
	CGO_ENABLED=0 gox -osarch="linux/amd64" -mod vendor -ldflags '-extldflags "-static" -X "main.version=${VERSION}"' github.com/xsteadfastx/jitsiexporter/cmd/jitsiexporter

clean:
	rm -f jitsiexporter

test:
	export GOFLAGS=-mod=vendor
	go test ./...

lint:
	golangci-lint run --enable-all--timeout 5m

dep-update:
	go get -u ./...
	go test ./...
	go mod tidy
	go mod vendor
