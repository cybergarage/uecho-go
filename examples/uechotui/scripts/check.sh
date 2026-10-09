#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
export GOWORK=off GOPROXY=off GOTOOLCHAIN=local
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then printf '%s\n' "$unformatted"; exit 1; fi
go vet ./...
go test -count=1 ./...
go test -race -count=1 ./...
mkdir -p bin
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/uechotui-darwin-arm64 .
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/uechotui-linux-arm64 .
