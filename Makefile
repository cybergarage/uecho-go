###################################################################
#
# uecho-go
#
# Copyright (C) Satoshi Konno 2017
#
# This is licensed under BSD-style license, see file COPYING.
#
###################################################################

SHELL := bash

PREFIX?=$(shell pwd)
GOPATH:=$(shell pwd)
export GOPATH

PACKAGE_NAME=uecho
BIN_DUMP_NAME=uechodump

GITHUB_ROOT=github.com/cybergarage/uecho-go
GITHUB=${GITHUB_ROOT}/net/${PACKAGE_NAME}

SOURCE_ROOT_DIR=src/${GITHUB}
PACKAGE_ID=${SOURCE_ROOT_DIR}
PACKAGES=\
	${PACKAGE_ID}

BINARY_DUMP=${GITHUB}/${BIN_DUMP_NAME}
BINARIES=${BINARY_DUMP}

.PHONY: version clean

all: test

VERSION_GO=${SOURCE_ROOT_DIR}/version.go

${VERSION_GO}: ${SOURCE_ROOT_DIR}/version.gen
	$< > $@

version: ${VERSION_GO}

format:
	gofmt -w src/${GITHUB}

vet: format
	go vet ${PACKAGES}

build: antlr vet
	go build -v ${PACKAGES}

test: antlr vet
	go test -v -cover ${PACKAGES}

install: antlr vet
	go install ${BINARIES}

clean:
	-rm ${PREFIX}/bin/*
	rm -rf _obj
	go clean -i ${PACKAGES}
