SHORT_NAME := object-storage-cli
REPO_PATH := github.com/deis/${SHORT_NAME}

HOST_OS := $(shell uname)
ifeq ($(HOST_OS),Darwin)
	GOOS=darwin
else
	GOOS=linux
endif

DEV_ENV_IMAGE := quay.io/deis/go-dev:0.11.1
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_PREFIX := docker run --rm -e GO15VENDOREXPERIMENT=1 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}
DIST_DIR := _dist
BINARY_NAME := objstorage


GO_LDFLAGS = -ldflags "-s -X ${REPO_PATH}/version.BuildVersion=${VERSION}"

VERSION ?= $(shell git rev-parse --short HEAD)

bootstrap:
	${DEV_ENV_CMD} glide install

build: binary-build

build-all:
	${DEV_ENV_CMD} gox -verbose \
	${GO_LDFLAGS} \
	-os="linux darwin " \
	-arch="amd64 386" \
	-output="$(DIST_DIR)/${BINARY_NAME}-latest-{{.OS}}-{{.Arch}}" .
	${DEV_ENV_CMD} gox -verbose ${GO_LDFLAGS} -os="linux darwin " -arch="amd64 386" -output="$(DIST_DIR)/deis-latest-{{.OS}}-{{.Arch}}" .
ifdef TRAVIS_TAG
	${DEV_ENV_CMD} gox -verbose ${GO_LDFLAGS} -os="linux darwin" -arch="amd64 386" -output="$(DIST_DIR)/${TRAVIS_TAG}/deis-${TRAVIS_TAG}-{{.OS}}-{{.Arch}}" .
else
	${DEV_ENV_CMD} gox -verbose ${GO_LDFLAGS} -os="linux darwin" -arch="amd64 386" -output="$(DIST_DIR)/${VERSION}/deis-${VERSION}-{{.OS}}-{{.Arch}}" .
endif


binary-build:
	${DEV_ENV_PREFIX} -e GOOS=${GOOS} ${DEV_ENV_IMAGE} go build -a ${GO_LDFLAGS} -o ${BINARY_NAME} .

dist: build-all

test:
	${DEV_ENV_CMD} go test $$(glide nv)
