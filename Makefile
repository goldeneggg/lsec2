NAME := lsec2
SRCS := $(shell find . -type f -name '*.go' | \grep -v 'vendor')
PACKAGES := $(shell ./scripts/_packages.sh)
PROF_DIR := ./.profile
GOVERSION := $(shell go version | awk '{print $$3;}')

.DEFAULT_GOAL := bin/$(NAME)

.PHONY: version
version:
	@echo $(shell ./scripts/_version.sh)

mod-dl:
	@GO111MODULE=on go mod download

mod-tidy:
	@GO111MODULE=on go mod tidy

bin/$(NAME): $(SRCS)
	@./scripts/build.sh bin/$(NAME)

.PHONY: install
install:
	@./scripts/install.sh

.PHONY: test
test:
	@go test -race -cover -v $(PACKAGES)

.PHONY: prof
prof:
	@[ ! -d $(PROF_DIR) ] && mkdir $(PROF_DIR); go test -bench . -benchmem -cover -coverprofile $(PROF_DIR)/cover.out $(PACKAGES)

prof-full-awsec2:
	@[ ! -d $(PROF_DIR) ] && mkdir $(PROF_DIR); go test -bench . -benchmem -blockprofile $(PROF_DIR)/block.out -cover -coverprofile $(PROF_DIR)/cover.out -cpuprofile $(PROF_DIR)/cpu.out -memprofile $(PROF_DIR)/mem.out ./awsec2

prof-full-printer:
	@[ ! -d $(PROF_DIR) ] && mkdir $(PROF_DIR); go test -bench . -benchmem -blockprofile $(PROF_DIR)/block.out -cover -coverprofile $(PROF_DIR)/cover.out -cpuprofile $(PROF_DIR)/cpu.out -memprofile $(PROF_DIR)/mem.out ./printer

prof-full-cmd-cli:
	@[ ! -d $(PROF_DIR) ] && mkdir $(PROF_DIR); go test -bench . -benchmem -blockprofile $(PROF_DIR)/block.out -cover -coverprofile $(PROF_DIR)/cover.out -cpuprofile $(PROF_DIR)/cpu.out -memprofile $(PROF_DIR)/mem.out ./cmd/lsec2/cli

.PHONY: vet
vet:
	@go vet $(PACKAGES)

.PHONY: lint
lint:
	@golint -set_exit_status $(PACKAGES)

.PHONY: validate
validate: vet lint

ci-test:
	@./scripts/ci-test.sh

.PHONY: ci
ci: ci-test validate

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor

vendor-build:
	@./scripts/build.sh bin/$(NAME) "-mod vendor"

lint-travis:
	@travis lint --org --debug .travis.yml

test-goreleaser:
	@GOVERSION=$(GOVERSION) goreleaser release --snapshot --skip-publish --rm-dist

ci-goreleaser:
	@export GOVERSION=$(GOVERSION) && curl -sL http://git.io/goreleaser | bash -s -- release --snapshot --skip-publish --rm-dist

.PHONY: clean
clean:
	@go clean $(PACKAGES)
	@rm -f bin/$(NAME)
	@rm -fr dist pkg
	@find . -name '*.test' -delete
	@rm -fr $(PROF_DIR)
