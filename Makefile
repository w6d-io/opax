SHELL=/bin/bash -o pipefail

export GO111MODULE        := on
export PATH               := bin:${PATH}
export PWD                := $(shell pwd)
export BUILD_DATE         := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
export VCS_REF            := $(shell git rev-parse HEAD)
export NEXT_TAG           ?=

ifeq (,$(shell go env GOOS))
GOOS       = $(shell echo $OS)
else
GOOS       = $(shell go env GOOS)
endif


GO_DEPENDENCIES = golang.org/x/tools/cmd/goimports@latest

define make-go-dependency
  # go install is responsible for not re-building when the code hasn't changed
  bin/$(notdir $1): go.mod go.sum Makefile
	GOBIN=$(PWD)/bin/ go install $1
endef
$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))
$(call make-lint-dependency)

# Formats the code
.PHONY: format
GOBIN = $(shell pwd)/bin
format:
	$(GOBIN)/goimports -w -local $(PWD) .

.PHONY: changelog
changelog:
	git-chglog -o CHANGELOG.md --next-tag $(NEXT_TAG)

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: fmt vet
	go test -v -coverpkg=./... -coverprofile=cover.out ./...
	@go tool cover -func cover.out | grep total

.PHONY: bin/goreadme
bin/goreadme:
	GOBIN=$(PWD)/bin \
	go install github.com/posener/goreadme/cmd/goreadme@latest

.PHONY: readme
readme: bin/goreadme
	./build/create_readme.sh

.PHONY: opa
OPA_BINARY = $(shell pwd)/bin/opa
SCRIPTBASH = $(shell pwd)/makefile.sh
GOBIN = $(shell pwd)/bin
ifeq (darwin,$(GOOS))
OPA_BINARY_URL=https://openpolicyagent.org/downloads/v0.43.0/opa_darwin_amd64
else
OPA_BINARY_URL=https://openpolicyagent.org/downloads/v0.43.0/opa_linux_amd64_static
endif
opa: ##init opa
ifeq (,$(wildcard $(OPA_BINARY)))
	mkdir -p $(GOBIN)
	wget $(OPA_BINARY_URL) -O $(OPA_BINARY)
	chmod +x $(OPA_BINARY)
else
	$(info ************ BINARY ALREADY EXIST **********)
endif

start:
	sh $(shell pwd)/makefile.sh config
	sh $(shell pwd)/makefile.sh run

stop: clean
ifeq (darwin,$(GOOS))
	lsof -i -P | grep 8181 | sed -e 's/.*opa     *//' -e 's#/.*##' | sed 's/ .*//' | xargs kill
else
	netstat -lnp | grep 8181 | sed -e 's/.*LISTEN *//' -e 's#/.*##' | xargs kill
endif

clean:
	rm -rf bin
	$(info ************ BIN FOLDER IS DELETED **********)