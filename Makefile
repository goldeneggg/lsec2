BINNAME := lsec2
GO ?= go
GOLINT ?= golint
PGMPKGPATH := .
TESTTARGET := ./...
SAVETARGET := ./...
LINTTARGET := ./...

GODEP ?= godep

all: build

build:
	$(GODEP) $(GO) build -ldflags="-w" -o $(GOBIN)/$(BINNAME) $(PGMPKGPATH)

test: vet
	$(GODEP) $(GO) test -race -v $(TESTTARGET)

vet:
	$(GODEP) $(GO) vet -n $(TESTTARGET)

depsave:
	$(GODEP) save $(SAVETARGET)

depbuild: depsave build

deprestore:
	$(GODEP) restore

lint:
	$(GOLINT) $(LINTTARGET)
