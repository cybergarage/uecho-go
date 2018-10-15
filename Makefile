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

GITHUB_ROOT=github.com/cybergarage/uecho-go
PACKAGE_ROOT=${GITHUB_ROOT}/net/echonet

PACKAGE_ID=${PACKAGE_ROOT}
PACKAGES=\
	${PACKAGE_ID} \
	${PACKAGE_ID}/encoding \
	${PACKAGE_ID}/protocol \
	${PACKAGE_ID}/transport

SOURCE_ROOT_DIR=src/${PACKAGE_ROOT}

EXSAMPLE_ROOT=${GITHUB_ROOT}/examples

BINARIES=\
	${EXSAMPLE_ROOT}/uechopost \
	${EXSAMPLE_ROOT}/uechosearch \
	${EXSAMPLE_ROOT}/uecholight

.PHONY: version clean

all: test

VERSION_GO=${SOURCE_ROOT_DIR}/version.go

${VERSION_GO}: ${SOURCE_ROOT_DIR}/version.gen
	$< > $@

version: ${VERSION_GO}

format:
	gofmt -w src/${PACKAGE_ROOT}

vet: format
	go vet ${PACKAGES}

build: vet
	go build -v ${PACKAGES}

test: vet
	go test -v -cover -timeout 30s ${PACKAGES}

install: vet
	go install ${BINARIES}

clean:
	-rm ${PREFIX}/bin/*
	rm -rf _obj
	go clean -i ${PACKAGES}
