NAME := lsec2
VERSION := $(shell ./scripts/_version.sh)
LD_FLAGS := $(shell ./scripts/_ldflags.sh)
SRCS := $(shell find . -type f -name '*.go')
BASE_PACKAGE := 'github.com/goldeneggg/lsec2'
SAVE_TARGET := ./...
PROF_DIR := ./.profile
PROF_TARGET := ./awsec2

.DEFAULT_GOAL := bin/$(NAME)

version:
	@echo $(VERSION)

ldflags:
	@echo $(LD_FLAGS)

dep:
	@dep ensure -v

dep-status:
	@dep status -v

bin/$(NAME): $(SRCS)
	@go build -v -a -tags netgo -installsuffix netgo -ldflags='$(LD_FLAGS)' -o bin/$(NAME)

install:
	@go install -v -ldflags='$(LD_FLAGS)'

test:
	@go test -race -cover -v $(BASE_PACKAGE)
	@go test -race -cover -v $(BASE_PACKAGE)/awsec2...

ci-test:
	@./scripts/ci-test.sh

prof:
	@[ ! -d $(PROF_DIR) ] && mkdir $(PROF_DIR); go test -bench . -benchmem -blockprofile $(PROF_DIR)/block.out -cover -coverprofile $(PROF_DIR)/cover.out -cpuprofile $(PROF_DIR)/cpu.out -memprofile $(PROF_DIR)/mem.out $(PROF_TARGET)

vet:
	@go tool vet -v --all -shadow ./*.go
	@go tool vet -v -all -shadow ./awsec2

lint:
	@${GOBIN}/golint -set_exit_status $(BASE_PACKAGE)
	@${GOBIN}/golint -set_exit_status $(BASE_PACKAGE)/awsec2

validate: vet lint
  
release:
	@echo "Releasing"
	@./scripts/release.sh

release-darwin-386:
	@echo "Releasing darwin-386"
	@./scripts/release.sh darwin 386

release-darwin-amd64:
	@echo "Releasing darwin-amd64"
	@./scripts/release.sh darwin amd64

release-linux-386:
	@echo "Releasing linux-386"
	@./scripts/release.sh linux 386

release-linux-amd64:
	@echo "Releasing linux-amd64"
	@./scripts/release.sh linux amd64

release-linux-arm:
	@echo "Releasing linux-arm"
	@./scripts/release.sh linux arm

release-windows-386:
	@echo "Releasing windows-386"
	@./scripts/release.sh windows 386

release-windows-amd64:
	@echo "Releasing windows-amd64"
	@./scripts/release.sh windows amd64

release-freebsd-386:
	@echo "Releasing freebsd-386"
	@./scripts/release.sh freebsd 386

release-freebsd-amd64:
	@echo "Releasing freebsd-amd64"
	@./scripts/release.sh freebsd amd64

upload:
	@echo "Uploading releases to github"
	@./scripts/upload.sh

formula:
	@echo "Generating formula"
	@./scripts/upload.sh formula-only

publish: release upload
