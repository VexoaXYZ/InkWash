.PHONY: build build-all install clean test run

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o inkwash

build-all:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/inkwash-windows-amd64.exe
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/inkwash-linux-amd64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/inkwash-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/inkwash-darwin-arm64

install:
	go build $(LDFLAGS) -o $(GOPATH)/bin/inkwash

clean:
	rm -rf bin/ inkwash inkwash.exe

test:
	go test -v ./...

run:
	go run main.go

# Compress binaries (optional, requires UPX)
compress:
	upx --best --lzma bin/inkwash-*
