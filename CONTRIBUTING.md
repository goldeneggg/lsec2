###  Guideline summary
- Fork it
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request

### Install Go
Please setup Golang environment, follow [Getting Started \- The Go Programming Language](https://golang.org/doc/instal)

__(required version 1.11 or later)__

### Fork
* Fork [lsec2](https://github.com/goldeneggg/lsec2/fork)
* Download forked repositry in your machine.

```sh
$ git clone <forked repositry URL>
$ cd lsec2
$ git remote add upstream https://github.com/goldeneggg/lsec2
```

### Download dependency libraries

```sh
$ make mod-dl
```

### Run tests

* Run `make test`

### Build

* Run `make`

### install

* Run `make install`

### Validate codes

* Run `make validate`
