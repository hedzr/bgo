# hedzr/bgo

[![Go](https://github.com/hedzr/bgo/actions/workflows/go.yml/badge.svg)](https://github.com/hedzr/bgo/actions/workflows/go.yml)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/bgo.svg?label=release)](https://github.com/hedzr/bgo/releases)
[![](https://img.shields.io/badge/go-dev-green)](https://pkg.go.dev/github.com/hedzr/bgo)
[![Docker Pulls](https://img.shields.io/docker/pulls/hedzr/bgo)](https://hub.docker.com/r/hedzr/bgo)
![Docker Image Version (latest semver)](https://img.shields.io/docker/v/hedzr/bgo)
![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/hedzr/bgo)


`bgo` provides a super easy way to build your go apps.
All things you need to do is entering the golang project directory and entering `bgo` and pressing <kbd>Enter</kbd>.  

> powered by [cmdr](https://github.com/hedzr/cmdr).

## Features

- Run go building with or without a config file `.bgo.yml`
- Scan the directory to grab all main packages and initial `.bgo.yml`
- 

## Getting Started

To run the CLI app:

```bash
# go install -v github.com/swaggo/swag/cmd/swag
go generate ./...          # run it once at least, for gen the swagger-doc files from skeletons
go run ./cli/app/cli/app   # build the mainly main.go
```

### Use Makefile for building and CI

You may use `make` simply:

```bash
make help    # list all available make targets, such as info, build, ...
make info    # print and review the golang build env

make build
```



## Inspired By:

- https://github.com/mitchellh/gox
- https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5
- https://github.com/davecheney/golang-crosscompile
- https://github.com/laher/goxc
- Makefile

## LICENSE

Apache 2.0


