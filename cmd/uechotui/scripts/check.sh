#!/bin/sh
set -eu
cd "$(dirname "$0")/../../.."
export GOWORK=off GOPROXY=off GOTOOLCHAIN=local
unformatted=$(gofmt -l internal/tui cmd/uechoctl cmd/uechotui)
if [ -n "$unformatted" ]; then printf '%s\n' "$unformatted"; exit 1; fi
go vet ./internal/tui/... ./cmd/uechoctl ./cmd/uechotui
go test -count=1 ./internal/tui/... ./cmd/uechoctl ./cmd/uechotui
# Root Makefile exports CGO_ENABLED=0 for portable builds; Linux race
# instrumentation needs cgo, so opt in only for this check.
CGO_ENABLED=1 go test -race -count=1 ./internal/tui/... ./cmd/uechoctl ./cmd/uechotui
mkdir -p cmd/uechotui/bin
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o cmd/uechotui/bin/uechoctl-darwin-arm64 ./cmd/uechoctl
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o cmd/uechotui/bin/uechoctl-linux-arm64 ./cmd/uechoctl
