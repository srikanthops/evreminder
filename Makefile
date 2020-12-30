VERSION ?= dev

pkgs := $(shell go list ./... | grep -v tck)
files := $(shell find . -path ./... -prune -o -name '*.go' -print)

.PHONY: all
all: format test build-binaries

VERSION:=$(shell ./scripts/version.sh)
BUILDTIME:=$(shell date +%FT%T%z)
binary=cstester

.PHONY: install
install:
	go install -ldflags "-X main.version=${VERSION}" cmd/

.PHONY: setup_devtools
setup_devtools:
	go get -u github.com/cloudstateio/go-support/cloudstate
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/onsi/ginkgo
	go get -u github.com/onsi/gomega

.PHONY: format
format:
	goimports -w $(files)

.PHONY: test
test: checkformat vet lint gotest

.PHONY: checkformat
checkformat:
	@if [ -n "$(shell goimports -l $(files))" ]; then \
		echo "unformatted files: $(shell goimports -l $(files))"; \
		echo "run make format"; \
		exit 1; \
	fi

.PHONY: vet
vet:
	go vet $(pkgs)

.PHONY: lint
lint:
	@for pkg in $(pkgs); do \
		golint -set_exit_status $$pkg || exit 1; \
	done;

.PHONY: gotest
gotest:
	go test -race $(pkgs)

.PHONY: gotestnocache
gotestnocache:
	go clean -testcache
	go test -race $(pkgs)


.PHONY: protoc
protoc:
	scripts/protogen.sh

.PHONY: run
run:
	go run ./cmd/cloudrun/ -v -c ./cmd/cloudrun/config.yaml

.PHONY: build-binaries
build-binaries:
	go build -ldflags "-X main.version=${VERSION}" -a -o "${binary}" ./cmd/cloudrun/
