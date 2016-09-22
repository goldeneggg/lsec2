lsec2 [![Build Status](http://drone.io/github.com/goldeneggg/lsec2/status.png)](https://drone.io/github.com/goldeneggg/lsec2/latest) [![Go Report Card](https://goreportcard.com/badge/github.com/goldeneggg/lsec2)](https://goreportcard.com/report/github.com/goldeneggg/lsec2) [![MIT License](http://img.shields.io/badge/license-MIT-lightgrey.svg)](https://github.com/goldeneggg/lsec2/blob/master/LICENSE)
==========
List view of aws ec2 instances

## Getting Started

### for Mac using homebrew

```bash
$ brew tap goldeneggg/lsec2
$ brew install lsec2
```

### go get

```bash
$ go get -u github.com/goldeneggg/lsec2
```

### set environment variables

2 variables are required

```bash
$ export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY_ID
$ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_ACCESS_KEY
```

## Usage

Please type `lsec2`

```bash
# show all instances info by list view
$ lsec2 show
```

You can get informations of instances as follows

* instance id
* private ip address
* public ip address
* instance type
* instance state
* tags

### Filter by tag

```bash
# filterd by a tag
$ lsec2 show TagName1=tagvalue1

# filterd by some tags
$ lsec2 show TagName1=tagvalue1 TagNameN=tagvalueN

# filterd by a tag multiple values separated comma
$ lsec2 show TagName1=tagvalue11,tagvalue12,tagvalue1N
```

### Options

```bash
# with header
$ lsec2 show -H

# show only private IP address
$ lsec2 show -p

# assign region using "--region" global option (default region is "ap-northeast-1")
$ lsec2 --region YOUR_REGION show
```

## Tips

### With peco
[peco](https://github.com/peco/peco) is a very useful interactive filtering tool.

* Example: show instances => select instance => SSH to selected instance

```bash
# add function to your .bashrc or .bash_profile or other shell dotfile
$ vi YOUR_DOTFILE

lssh () {
  IP=$(lsec2 show $@ | peco | awk -F "\t" '{print $2}')
  if [ $? -eq 0 -a "${IP}" != "" ]
  then
      eval "ssh ${IP}"
  fi
}


# load dotfile
$ source YOUR_DOTFILE

# shortcut "lsec2 OPTIONS TAG_FILTERS" => "ssh PRIVATE_IP"
$ lssh TagName1=tagvalue1
```

## Development

Install Go, and setup using `go get`

```bash
$ go get -u github.com/goldeneggg/lsec2
```

Initialize vendoring libraries

```bash
$ make depget
```

Build

```bash
$ make
```

## Contact

* Bugs: [issues](https://github.com/goldeneggg/lsec2/issues)


## ChangeLog
[CHANGELOG](CHANGELOG.md) file for details.


## License

[LICENSE](LICENSE) file for details.

## Special Thanks
[@sugitak](https://github.com/sugitak)

## TODO
* [ ] Add tests
* [ ] Add some sub commands
    * `lsec2 ssh` - execute ssh
