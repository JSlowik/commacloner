PROJ=commacloner
ORG_PATH=github.com/jslowik
REPO_PATH=$(ORG_PATH)/$(PROJ)
export PATH := $(PWD)/bin:$(PATH)

VERSION ?= $(shell ./scripts/git-version)

$( shell mkdir -p bin )

user=$(shell id -u -n)
group=$(shell id -g -n)

export GOBIN=$(PWD)

LD_FLAGS="-s -w -X $(REPO_PATH)/version.Version=$(VERSION)"


build: bin/commacloner

bin/commacloner:
	@mkdir -p bin/
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/commacloner

.PHONY: release-binary
release-binary:
	@go build -o /go/bin/commacloner -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/dex

action-binaries:
	@go build -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/commacloner

makecover:
	@go test ./... -cover -coverprofile=coverage.out

cover: makecover
	@go tool cover -func=coverage.out | grep total

cover_html: makecover
	@go tool cover -html=coverage.out

bench:
	@go test -bench=. -run=^a ./... | grep Benchmark

test:
	@go test -v ./...

testrace:
	@go test -v --race ./...

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- $(shell ./scripts/get-latest-golangci-lint)
	@mv bin/golangci-lint $@

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run -v

.PHONY: fix
fix: bin/golangci-lint ## Fix lint violations
	bin/golangci-lint run --fix

clean:
	@rm -rf bin/

testall: testrace

FORCE:

.PHONY: test testrace testall
