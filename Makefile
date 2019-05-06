NAME := lsec2
PROF_DIR := ./.profile
AWS_SDK_GO_PKG := github.com/aws/aws-sdk-go

# Note: NOT use lazy initializer because make is unstable.
#SRCS = $(eval SRCS := $(shell find . -type f -name '*.go' | \grep -v 'vendor'))$(SRCS)
#PACKAGES = $(eval PACKAGES := $(shell ./scripts/_packages.sh))$(PACKAGES)
#GOVERSION = $(eval GOVERSION := $(shell go version | awk '{print $$3;}'))$(GOVERSION)
SRCS = $(shell find . -type f -name '*.go' | \grep -v 'vendor')
PACKAGES = $(shell ./scripts/_packages.sh)
GOVERSION = $(shell go version | awk '{print $$3;}')

.DEFAULT_GOAL := bin/$(NAME)

.PHONY: version
version:
	@echo $(shell ./scripts/_version.sh)

mod-dl:
	@GO111MODULE=on go mod download

mod-tidy:
	@GO111MODULE=on go mod tidy

list-versions-aws-sdk-go:
	@go list -u -m -versions $(AWS_SDK_GO_PKG) | tr ' ' '\n'

update-aws-sdk-go:
	@read -p 'Input Module Query(e.g. "<v1.20")?: ' query; echo query=$$query; GO111MODULE=on go get $(AWS_SDK_GO_PKG)@''$$query''

bin/$(NAME): $(SRCS)
	@./scripts/build.sh bin/$(NAME)

.PHONY: rm
rm:
	@rm bin/$(NAME)

.PHONY: rebuild
rebuild: rm bin/$(NAME)

.PHONY: install
install:
	@./scripts/install.sh

.PHONY: test
test:
	@go test -race -cover -v $(PACKAGES)

.PHONY: prof
prof:
	@[ ! -d $(PROF_DIR) ] && mkdir $(PROF_DIR); go test -bench . -benchmem -cover -coverprofile $(PROF_DIR)/cover.out $(PACKAGES)

func_prof_full = [ ! -d $(PROF_DIR) ] && mkdir $(PROF_DIR); go test -bench . -benchmem -blockprofile $(PROF_DIR)/block.out -cover -coverprofile $(PROF_DIR)/cover.out -cpuprofile $(PROF_DIR)/cpu.out -memprofile $(PROF_DIR)/mem.out ./$1

prof-full-awsec2:
	@$(call func_prof_full,awsec2)

prof-full-printer:
	@$(call func_prof_full,printer)

prof-full-cmd-cli:
	@$(call func_prof_full,cmd/lsec2/cli)

.PHONY: vet
vet:
	@go vet $(PACKAGES)

.PHONY: lint
lint:
	@golint -set_exit_status $(PACKAGES)

.PHONY: validate
validate: vet lint

chk_latest = go list -u -m $1
chk_versions = go list -u -m -versions $1 | tr ' ' '\n'

chk-latest-all:
	@$(call chk_latest,all)

chk-latest-aws-sdk-go:
	@$(call chk_latest,github.com/aws/aws-sdk-go)

chk-versions-aws-sdk-go:
	@$(call chk_versions,github.com/aws/aws-sdk-go)

chk-latest-urfave-cli:
	@$(call chk_latest,github.com/urfave/cli)

chk-versions-urfave-cli:
	@$(call chk_versions,github.com/urfave/cli)

mod-dl:
	@GO111MODULE=on go mod download

mod-tidy:
	@GO111MODULE=on go mod tidy

ci-test:
	@./scripts/ci-test.sh

.PHONY: ci
ci: ci-test vet

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor

vendor-build:
	@./scripts/build.sh bin/$(NAME) "-mod vendor"

lint-travis:
	@travis lint --org --debug .travis.yml

# Note: require "brew install rpmbuild" on OS X
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

mod-clean:
	@go clean -modcache
