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

PREFIX?=$(shell pwd)
GOPATH:=$(shell pwd)
export GOPATH

PACKAGE_NAME=uecho
BIN_DUMP_NAME=uechodump
BIN_SEARCH_NAME=uechosearch
BIN_LIGHT_NAME=uecholight

GITHUB_ROOT=github.com/cybergarage/uecho-go
GITHUB=${GITHUB_ROOT}/net/${PACKAGE_NAME}

PACKAGE_ID=${GITHUB}
PACKAGES=\
	${PACKAGE_ID} \
	${PACKAGE_ID}/encoding \
	${PACKAGE_ID}/protocol \
	${PACKAGE_ID}/transport

SOURCE_ROOT_DIR=src/${PACKAGE_ROOT}

BINARIES=\
	${GITHUB}/${BIN_DUMP_NAME} \
	${GITHUB}/${BIN_SEARCH_NAME} \
	${GITHUB}/${BIN_LIGHT_NAME}

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

build: vet
	go build -v ${PACKAGES}

test: vet
	go test -v -cover ${PACKAGES}

install: vet
	go install ${BINARIES}

clean:
	-rm ${PREFIX}/bin/*
	rm -rf _obj
	go clean -i ${PACKAGES}
