##  Guideline summary
- Fork it
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request

## Install Go
Please setup Golang environment, follow [Getting Started \- The Go Programming Language](https://golang.org/doc/instal)

__(required version 1.11 or later)__

## Setup
* Fork [lsec2](https://github.com/goldeneggg/lsec2/fork)
* Clone forked repositry in your machine.
* Run `make`

```sh
$ git clone <forked repositry URL>
$ cd lsec2
$ git remote add upstream https://github.com/goldeneggg/lsec2

# 1st go build
# At first, it will find and download dependency libraries on [Module-aware mode](https://golang.org/cmd/go/#hdr-Module_aware_go_get) before building.
$ make
go: finding github.com/mattn/go-colorable v0.1.1
go: finding github.com/aws/aws-sdk-go v1.17.14
:
:
```

## Build

* Run `make`

## Run tests

* Run `make test`

## install

* Run `make install`

## Validate codes

* Run `make validate`

## Add some new dependency libraries
___On Module-aware mode___

1. Edit your code and add import library.
1. Run `make rebuild` or `make test`
