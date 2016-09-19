BINNAME := lsec2
GO ?= go
GOLINT ?= golint
PGMPKGPATH := .
TESTTARGET := ./...
SAVETARGET := ./...
LINTTARGET := ./...

GODEP ?= godep

all: depbuild

depbuild: depsave
	$(GODEP) $(GO) build -o $(GOBIN)/$(BINNAME) $(PGMPKGPATH)

deptest: depvet
	$(GODEP) $(GO) test -race -v $(TESTTARGET)

depvet: depsave
	$(GODEP) $(GO) vet -n $(TESTTARGET)

depsave:
	$(GODEP) save $(SAVETARGET)

depget:
	$(GODEP) get

lint:
	$(GOLINT) $(LINTTARGET)
