###  Guideline summary
- Fork it
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request

### Install Go
Please setup Golang environment, follow [Getting Started \- The Go Programming Language](https://golang.org/doc/instal)

___(required version 1.5 or later)___, and recommend version 1.7 or later)

### Fork
* Fork [lsec2](https://github.com/goldeneggg/lsec2/fork)
* Download forked repositry in your `$GOPATH`

```sh
$ mkdir -p $GOPATH/src/github.com/goldeneggg
$ cd $GOPATH/src/github.com/goldeneggg
$ git clone <forked repositry URL>
$ cd lsec2
$ git remote add upstream https://github.com/goldeneggg/lsec2
```

### Install godep
* Install [tools/godep: dependency tool for go](https://github.com/tools/godep)

```sh
$ go get -u github.com/tools/godep

# (required version v7x or later, please check `godep version`)
$ godep version
```

### Add or remove dependencies

* Run `go get hoge/huga`
* Edit your code to [not] `import hoge/huga`
* Run `make dep-save`

### Run tests

* Run `make test`

### Run vet

* Run `make vet`

### Build and install

* Run `make`

### Run golint

* Run `make lint`
