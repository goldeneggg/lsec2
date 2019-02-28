NAME := lsec2
SRCS := $(shell find . -type f -name '*.go' | \grep -v 'vendor')
VERSION := $(shell ./scripts/_version.sh)
LDFLAGS := $(shell ./scripts/_ldflags.sh)
PACKAGES := $(shell ./scripts/_packages.sh)
PROF_DIR := ./.profile
PROF_TARGET := ./awsec2

.DEFAULT_GOAL := bin/$(NAME)

.PHONY: version
version:
	@echo $(VERSION)

mod-dl:
	@GO111MODULE=on go mod download

bin/$(NAME): $(SRCS)
	@./scripts/build.sh bin/$(NAME)

.PHONY: install
install:
	@go install -v -ldflags='$(LDFLAGS)'

.PHONY: test
test:
	@go test -race -cover -v $(PACKAGES)

ci-test:
	@./scripts/ci-test.sh

.PHONY: prof
prof:
	@[ ! -d $(PROF_DIR) ] && mkdir $(PROF_DIR); go test -bench . -benchmem -blockprofile $(PROF_DIR)/block.out -cover -coverprofile $(PROF_DIR)/cover.out -cpuprofile $(PROF_DIR)/cpu.out -memprofile $(PROF_DIR)/mem.out $(PROF_TARGET)

.PHONY: vet
vet:
	@go vet -n -x $(PACKAGES)

.PHONY: lint
lint:
	@${GOBIN}/golint -set_exit_status $(PACKAGES)

.PHONY: validate
validate: vet lint

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor

vendor-build:
	@./scripts/build.sh bin/$(NAME) "-mod vendor"

lint-travis:
	@travis lint --org --debug .travis.yml

test-goreleaser:
	@goreleaser release --snapshot --skip-publish --rm-dist
