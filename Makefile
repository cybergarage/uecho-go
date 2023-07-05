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

#PREFIX?=$(shell pwd)
#GOPATH:=$(shell pwd)
#export GOPATH
export CGO_ENABLED=0

PACKAGE_NAME=net/echonet

MODULE_ROOT=github.com/cybergarage/uecho-go
SOURCE_ROOT=${PACKAGE_NAME}
PACKAGE_ROOT=${MODULE_ROOT}/${PACKAGE_NAME}

SOURCES=\
	${SOURCE_ROOT} \
	${SOURCE_ROOT}/encoding \
	${SOURCE_ROOT}/protocol \
	${SOURCE_ROOT}/transport

PACKAGE_ID=${PACKAGE_ROOT}
PACKAGES=\
	${PACKAGE_ID} \
	${PACKAGE_ID}/encoding \
	${PACKAGE_ID}/protocol \
	${PACKAGE_ID}/transport

BINARY_ROOT=${MODULE_ROOT}/examples

BINARIES=\
	${BINARY_ROOT}/uechopost \
	${BINARY_ROOT}/uechosearch \
	${BINARY_ROOT}/uecholight \
	${BINARY_ROOT}/uechobench

.PHONY: version clean

all: test

version:
	@pushd ${SOURCE_ROOT} && ./version.gen > version.go && popd

format:
	gofmt -w ${SOURCE_ROOT}

vet: format
	go vet ${PACKAGE_ROOT}

lint: vet
	golangci-lint run ${SOURCES}

build: lint
	go build -v ${PACKAGES}

test: lint
	go test -v -cover -timeout 300s ${PACKAGES}

test_only:
	go test -v -cover -timeout 300s ${PACKAGES}

install:
	go install ${BINARIES}

clean:
	go clean -i ${PACKAGES}
