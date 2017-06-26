SHELL          = /bin/bash

BASE_IMAGE     = golang:1.7-alpine

GROUP_NAME     = answer1991
TARGET         = discovery
IMAGE_NAME     = $(GROUP_NAME)/$(TARGET)
MAJOR_VERSION = $(shell cat VERSION)
DATE = $(shell date +%Y%m%d)

MAJOR_VERSION = $(shell cat VERSION)
GIT_VERSION   = $(shell git log -1 --pretty=format:%h)
GIT_NOTES     = $(shell git log -1 --oneline)

default: image

compile:
	docker run --rm -v $(shell pwd):/go/src/github.com/answer1991/discovery -w /go/src/github.com/answer1991/discovery ${BASE_IMAGE} go build -v

image: compile
	docker build --rm -t ${IMAGE_NAME}:${MAJOR_VERSION}-${GIT_VERSION}-${DATE}  .
	@echo Image is : ${IMAGE_NAME}:${MAJOR_VERSION}-${GIT_VERSION}-${DATE}