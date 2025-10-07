###################################################################
#
# uecho-go
#
# Copyright (C) The uecho-go Authors 2017
#
# This is licensed under BSD-style license, see file COPYING.
#
###################################################################

SHELL := bash

GOBIN := $(shell go env GOPATH)/bin
PATH := $(GOBIN):$(PATH)

export CGO_ENABLED=0

PKG_NAME=net/echonet

MODULE_ROOT=github.com/cybergarage/uecho-go
PKG_SOURCE_ROOT=${PKG_NAME}
PKG_ROOT=${MODULE_ROOT}/${PKG_NAME}

PKG_COVER=uecho-cover

PKG_SOURCES=\
	${PKG_SOURCE_ROOT} \
	${PKG_SOURCE_ROOT}/encoding \
	${PKG_SOURCE_ROOT}/protocol \
	${PKG_SOURCE_ROOT}/transport

PKG_ID=${PKG_ROOT}
PKGES=\
	${PKG_ID} \
	${PKG_ID}/encoding \
	${PKG_ID}/protocol \
	${PKG_ID}/transport

EXAMPLE_PKG_SOURCE_ROOT=examples
EXAMPLE_ROOT=${MODULE_ROOT}/${EXAMPLE_PKG_SOURCE_ROOT}

CMD_PKG_SOURCE_ROOT=cmd
CMD_ROOT=${MODULE_ROOT}/${CMD_PKG_SOURCE_ROOT}

BINARIES=\
	${CMD_ROOT}/uechoctl \
	${EXAMPLE_ROOT}/uechopost \
	${EXAMPLE_ROOT}/uechosearch \
	${EXAMPLE_ROOT}/uecholight \
	${EXAMPLE_ROOT}/uechobench

CMD_DOC_ROOT=doc/cmd

.PHONY: version clean
.IGNORE: lint

all: test

version:
	@pushd ${PKG_SOURCE_ROOT} && ./version.gen > version.go && popd
	-git commit ${PKG_SOURCE_ROOT}/version.go -m "Update version"

format: version
	gofmt -w ${PKG_SOURCE_ROOT} ${CMD_PKG_SOURCE_ROOT} ${EXAMPLE_PKG_SOURCE_ROOT}

vet: format
	go vet ${PKG_ROOT}

lint: vet
	golangci-lint run ${PKG_SOURCES}

build: lint
	go build -v ${PKGES}

test: lint
	go test -v -p 1 -timeout 10m -cover -coverpkg=${PKG_ROOT}/... -coverprofile=${PKG_COVER}.out ${PKG_ROOT}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

cover: test
	open ${PKG_COVER}.html || xdg-open ${PKG_COVER}.html || gnome-open ${PKG_COVER}.html

godoc:
	go install golang.org/x/tools/cmd/godoc@latest
	open http://localhost:6060/pkg/${PKG_ID}/ || xdg-open http://localhost:6060/pkg/${PKG_ID}/ || gnome-open http://localhost:6060/pkg/${PKG_ID}/
	godoc -http=:6060 -play

gendoc:
	@mkdir -p ${CMD_DOC_ROOT}
	@rm -f ${CMD_DOC_ROOT}/*.md
	go run ./scripts/gendoc.go && git add ${CMD_DOC_ROOT} && git commit ${CMD_DOC_ROOT} -m "docs: update CLI documentation"

install: gendoc
	go install ${BINARIES}

clean:
	go clean -i ${PKGES}
