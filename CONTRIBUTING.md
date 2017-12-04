###  Guideline summary
- Fork it
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request

### Install Go
Please setup Golang environment, follow [Getting Started \- The Go Programming Language](https://golang.org/doc/instal)

__(required version 1.6 or later)__ , and recommend version 1.9 or later

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

### Install dep
* [Setup golang/dep: Go dependency management tool](https://github.com/golang/dep#setup)

```sh
$ go get -u github.com/golang/dep/cmd/dep
```

### Add dependencies on local

* Run `make dep`

### Run tests

* Run `make test`

### Build

* Run `make`

### install

* Run `make install`

### Validate codes

* Run `make validate`
