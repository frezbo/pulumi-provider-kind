SHELL := bash

PROJECT_NAME := Pulumi KIND Resource Provider

PACK             := kind
PACKDIR          := sdk
PROJECT          := github.com/frezbo/pulumi-provider-kind
NODE_MODULE_NAME := @frezbo/kind

PROVIDER        := pulumi-resource-${PACK}
CODEGEN         := pulumi-gen-${PACK}
VERSION         := 0.0.1
PROVIDER_PATH   := provider/v3
VERSION_PATH     := ${PROVIDER_PATH}/pkg/version.Version

SCHEMA_FILE     := provider/cmd/${PROVIDER}/schema.json
GOPATH			:= $(shell go env GOPATH)

WORKING_DIR     := $(shell pwd)
TESTPARALLELISM := 4

kindgen:
	(cd provider && go build -a -o $(WORKING_DIR)/bin/${CODEGEN} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" ${PROJECT}/${PROVIDER_PATH}/cmd/$(CODEGEN))

go_sdk:: kindgen
	$(WORKING_DIR)/bin/${CODEGEN} $(SCHEMA_FILE) $(PACKDIR)

kindprovider:
	(cd provider && VERSION=${VERSION} go generate cmd/${PROVIDER}/main.go)
	(cd provider && go build -a -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

.PHONY: build
build:: kindgen kindprovider go_sdk

lint::
	for DIR in "provider" "sdk" "tests" ; do \
		pushd $$DIR && golangci-lint run -c ../.golangci.yml --timeout 10m && popd ; \
	done

install::
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin
