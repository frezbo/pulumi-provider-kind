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

ensure::
	cd provider && go mod tidy
	cd sdk && go mod tidy
	cd tests && go mod tidy

kindgen::
	(cd provider && go build -a -o $(WORKING_DIR)/bin/${CODEGEN} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" ${PROJECT}/${PROVIDER_PATH}/cmd/$(CODEGEN))

schema:: kindgen
	@echo "Generating Pulumi schema..."
	$(WORKING_DIR)/bin/${CODEGEN} schema "" $(CURDIR)
	@echo "Finished generating schema."

kindprovider::
	(cd provider && VERSION=${VERSION} go generate cmd/${PROVIDER}/main.go)
	(cd provider && go build -a -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

kindprovider_debug::
	(cd provider && go build -a -o $(WORKING_DIR)/bin/${PROVIDER} -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

test_provider::
	cd provider/pkg && go test -short -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM} ./...

go_sdk::
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION} go $(SCHEMA_FILE) $(CURDIR)

nodejs_sdk:: VERSION := $(shell pulumictl get version --language javascript)
nodejs_sdk::
	$(WORKING_DIR)/bin/$(CODEGEN) -version=${VERSION} nodejs $(SCHEMA_FILE) $(CURDIR)
	cd ${PACKDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc
	cp README.md LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock ${PACKDIR}/nodejs/bin/
	sed -i.bak 's/$${VERSION}/$(VERSION)/g' ${PACKDIR}/nodejs/bin/package.json

.PHONY: build
build:: kindgen schema kindprovider go_sdk nodejs_sdk

# Required for the codegen action that runs in pulumi/pulumi
only_build:: build

lint::
	for DIR in "provider" "sdk" "tests" ; do \
		pushd $$DIR && golangci-lint run -c ../.golangci.yml --timeout 10m && popd ; \
	done

install:: install_go_sdk install_nodejs_sdk
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin

generate_schema:: schema

install_go_sdk::
	#target intentionally blank

install_nodejs_sdk::
	-yarn unlink --cwd $(WORKING_DIR)/sdk/nodejs/bin
	yarn link --cwd $(WORKING_DIR)/sdk/nodejs/bin
