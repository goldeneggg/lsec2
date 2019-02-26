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

# Duplicated tasks as follows

.PHONY: release
release:
	@./scripts/release.sh

release-darwin-386:
	@./scripts/release.sh darwin 386

release-darwin-amd64:
	@./scripts/release.sh darwin amd64

release-linux-386:
	@./scripts/release.sh linux 386

release-linux-amd64:
	@./scripts/release.sh linux amd64

release-linux-arm:
	@./scripts/release.sh linux arm

release-windows-386:
	@./scripts/release.sh windows 386

release-windows-amd64:
	@./scripts/release.sh windows amd64

release-freebsd-386:
	@./scripts/release.sh freebsd 386

release-freebsd-amd64:
	@./scripts/release.sh freebsd amd64

.PHONY: upload
upload:
	@./scripts/upload.sh

.PHONY: formula
formula:
	@./scripts/upload.sh formula-only

.PHONY: publish
publish: release upload
