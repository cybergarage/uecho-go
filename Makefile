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

PACKAGE_NAME=net/echonet

MODULE_ROOT=github.com/cybergarage/uecho-go
SOURCE_ROOT=${PACKAGE_NAME}
PACKAGE_ROOT=${MODULE_ROOT}/${PACKAGE_NAME}

PACKAGE_ID=${PACKAGE_ROOT}
PACKAGES=\
	${PACKAGE_ID} \
	${PACKAGE_ID}/encoding \
	${PACKAGE_ID}/protocol \
	${PACKAGE_ID}/transport


EXSAMPLE_ROOT=${MODULE_ROOT}/examples

BINARIES=\
	${EXSAMPLE_ROOT}/uechopost \
	${EXSAMPLE_ROOT}/uechosearch \
	${EXSAMPLE_ROOT}/uecholight

.PHONY: version clean

all: test

version:
	@pushd ${SOURCE_ROOT} && ./version.gen > version.go && popd

format:
	gofmt -w ${SOURCE_ROOT}

vet: format
	go vet ${PACKAGE_ROOT}

build: vet
	go build -v ${PACKAGES}

test: vet
	go test -v -cover -timeout 60s ${PACKAGES}

install: vet
	go install ${BINARIES}

clean:
	go clean -i ${PACKAGES}
