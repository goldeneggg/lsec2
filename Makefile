BINNAME := lsec2
PGM_PATH := 'github.com/goldeneggg/lsec2'
TEST_TARGET := ./...
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
