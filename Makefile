BINNAME := lsec2
PGM_PATH := 'github.com/goldeneggg/lsec2'
SAVE_TARGET := ./...

all: build

build:
	@echo "Building ${GOBIN}/$(BINNAME)"
	@GO15VENDOREXPERIMENT=1 godep go build -o ${GOBIN}/$(BINNAME) $(PGM_PATH)

test:
	@echo "Testing"
	@GO15VENDOREXPERIMENT=1 godep go test -race -v $(PGM_PATH)
	@GO15VENDOREXPERIMENT=1 godep go test -race -v $(PGM_PATH)/awsec2...

vet:
	@echo "Vetting"
	@GO15VENDOREXPERIMENT=1 godep go tool vet --all -shadow ./*.go
	@GO15VENDOREXPERIMENT=1 godep go tool vet -all -shadow ./awsec2

dep-save:
	@echo "Run godep save"
	@GO15VENDOREXPERIMENT=1 godep save -v $(SAVE_TARGET)

dep-saved-build: dep-save build

lint:
	@echo "Linting"
	@GO15VENDOREXPERIMENT=1 ${GOBIN}/golint $(PGM_PATH)
	@GO15VENDOREXPERIMENT=1 ${GOBIN}/golint $(PGM_PATH)/awsec2

release:
	@echo "Releasing"
	@./scripts/release.sh

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
	@echo "Releasing linux-amd64"
	@./scripts/release.sh linux amd64

release-freebsd-amd64:
	@echo "Releasing freebsd-amd64"
	@./scripts/release.sh freebsd amd64

publish: release
	@echo "Publishing releases to github"
	@./scripts/publish.sh

formula:
	@echo "Generating formula"
	@./scripts/publish.sh formula-only
