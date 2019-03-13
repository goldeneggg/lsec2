NAME := lsec2
SRCS := $(shell find . -type f -name '*.go' | \grep -v 'vendor')
PACKAGES := $(shell ./scripts/_packages.sh)
PROF_DIR := ./.profile

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

ci-test:
	@./scripts/ci-test.sh

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
