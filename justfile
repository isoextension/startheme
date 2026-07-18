default:
    @just --list

version := "2.0.0"
git_hash := `git rev-parse --short HEAD`
ldflags  := "-X main.Version=" + version + " -X main.GitHash=" + git_hash

build-linux: build-linux-amd64 build-linux-arm64
build-windows: build-windows-x86-64 build-windows-arm64
build-all: build-linux build-windows

build:
    mkdir -p bin
    go build -ldflags '{{ldflags}}' -o startheme ./main

build-linux-amd64:
    mkdir -p bin
    GOOS=linux GOARCH=amd64 go build -ldflags '{{ldflags}}' -o bin/startheme-linux-amd64 ./main

build-linux-arm64:
    mkdir -p bin
    GOOS=linux GOARCH=arm64 go build -ldflags '{{ldflags}}' -o bin/startheme-linux-arm64 ./main

build-windows-x86-64:
    mkdir -p bin
    GOOS=windows GOARCH=amd64 go build -ldflags '{{ldflags}}' -o bin/startheme-windows-x86_64.exe ./main

build-windows-arm64:
    mkdir -p bin
    GOOS=windows GOARCH=arm64 go build -ldflags '{{ldflags}}' -o bin/startheme-windows-arm64.exe ./main

clean:
    rm -rf bin startheme*