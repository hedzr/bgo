# hedzr/bgo

[![Go](https://github.com/hedzr/bgo/actions/workflows/go.yml/badge.svg)](https://github.com/hedzr/bgo/actions/workflows/go.yml)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/bgo.svg?label=release)](https://github.com/hedzr/bgo/releases)
[![go.dev](https://img.shields.io/badge/go-dev-green)](https://pkg.go.dev/github.com/hedzr/bgo)
[![Docker Pulls](https://img.shields.io/docker/pulls/hedzr/bgo)](https://hub.docker.com/r/hedzr/bgo)
![Docker Image Version (latest semver)](https://img.shields.io/docker/v/hedzr/bgo)
![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/hedzr/bgo)

`bgo` provides a super easy way to build your go apps.
All things you need to do is entering the golang project directory and entering `bgo` and pressing <kbd>Enter</kbd>.

> powered by [cmdr](https://github.com/hedzr/cmdr).

![tip-for-bgo](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20220202111546956.png)

## Features

`bgo` makes your golang building life easier, it's an efficient and extensible build tool.

- `bgo` or `bgo build`: Run go building with or without a config file `.bgo.yml`
- `bgo run -- ...`: Forward command-line to `go run` as is
- `bgo test -- ...`: Forward command-line to `go test` as is
- `bgo init`: Scan the directory to grab all main packages and initial `.bgo.yml`
- `bgo list`: list managed projects in all of `.bgo.yaml`
- `bgo sbom`: dump [#SBOM](#sbom) information itself or specified executables
- More extensible commands
  - pre-builtins (*bundled by [`80.aliases.yml`](https://github.com/hedzr/bgo/blob/master/ci/etc/bgo/conf.d/80.aliases.yml)*)
    - `bgo chk|check-code-qualities`
    - `bgo cov|coverage`

    see also [#Extensible Commands](https://github.com/hedzr/bgo#extensible-commands)
- While

  - you have lots of small cli apps in many sub-directories
  - one project wants go1.18beta1 and others stays in go1.17
  - too much `-tags`, `-ldflags`, `-asmflags`, ...
  - multiple targets, cross compiling
  - etc.

    have a try with `bgo`.

## History

- v0.5.25
  - security patch: bump golang.org/x/crypto from 0.15.0 to 0.17.0, and more

- v0.5.21/23
  - removed sending to coverall since a unknown/untagged error when gh action triggering the release 
  - security patch: bump golang.org/x/net from 0.15.0 to 0.17.0, and more
  - fixed a bug from errors.v3
  - upgrade deps

- v0.5.19
  - upgrade deps
  - fixed: `StripOrderPrefix`, ...

- v0.5.18
  - upgrade deps
  - fixed: change `SaveAs`(`--output`) field from string slice to string

- v0.5.16
  - upgrade deps
  - fixed a little

- More in [CHANGELOG](https://github.com/hedzr/bgo/blob/master/CHANGELOG)

## Getting Started

### Install

Download the prebuilt binaries from Release page.

Or Built from source code:

```bash
go install github.com/hedzr/bgo@latest
bgo --help         # -h
bgo --version      # -V
bgo --build-info   # -#
```

You may clone and compile `bgo` from source code:

```bash
git clone https://github.com/hedzr/bgo.git
cd bgo
go run . -s
```

`go run . -s` will run bgo from source code and install itself to ~/go/bin/.

![image-20220130101919648](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20220130101919648.png)

You could run bgo via docker way:

```bash
docker pull hedzr/bgo
# or
docker pull ghcr.io/hedzr/bgo:latest
```

And run it:

```bash
docker run -it --rm -v $PWD:/app -v /tmp:/tmp -v /tmp/go-pkg:/go/pkg hedzr/bgo
```

For macOS/Linux, there is a brew formula:

```bash
brew install hedzr/brew/bgo
```

### Run directly

To work without `.bgo.yml`, simply go into a folder and run bgo, the cli apps under the folder will be found out and built.

```bash
cd my-projects
bgo  # start `bgo` in auto mode, see Scopes chapter below
```

![image-20220128104835837](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20220128104835837.png)

Filter the target systems by `-for OS/ARCH`, `-os OS` and `-arch ARCH`:

```bash
bgo --for linux/386 -for linux/amd64,darwin/arm64
bgo -os linux -arch 386 -arch amd64 -arch arm64
```

Both long and short options are available for `for`, `os` and `arch`.
Their values are stacked by specifying multiple times, so `-arch 386 -arch amd64 -arch arm64` is same with `--arch 386,amd64,arm64`.

### Run with `.bgo.yml`

#### Create `.bgo.yml` at first

```bash
cd my-projects
bgo init  # create bgo.yml by scanning
mv bgo.yml .bgo.yml # rename it
```

#### Tune `.bgo.yml`

See sample of `.bgo.yml`

[.bgo.yml](https://github.com/hedzr/bgo/blob/master/.bgo.yaml)

#### Run

```bash
bgo
```

bgo will load projects from `.bgo.yml` and build them

### Scopes

1. `bgo -s`: short mode - this will build first project with current GOOS/GOARCH.
2. `bgo`|`bgo -a`: auto mode - build projects in `.bgo.yml`
3. `bgo -f`: full mode - build by scanning current directory

### Using shell auto-completion

Run `bgo gen sh --zsh` to install auto-completion script to the proper location and enable the feature:

![image-20220130092618399](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20220130092618399.png)

For macOS, you may install zsh completion script manually (not recommended, the better way is by HomeBrew):

```bash
bgo gen sh --zsh --output "$(brew --prefix)/share/zsh/site-functions/_bgo"
```

Run `bgo gen sh --bash -o=bgo.bash` to get bash completions script and put it to the right location. Generally it should be:

```bash
bgo generate shell --bash -o=bgo.bash
mv bgo.bash /etc/bash-completion.d/bgo
```

Nothing needs to do if installed via brew (since v0.3.3+).

### Others Subcommands

#### `run`

`bgo run` can forward arguments to `go run`. For the command line

```bash
bgo run -- ./...
```

The underlying real command will be run:

```bash
go run ./...
```

#### `test`

`bgo test` can forward arguments to `go test`. For the command line

```bash
bgo test -- ./...
```

The underlying real command will be run:

```bash
go test ./...
```

### `SBOM`

Checking executable's [SBOM (`Software Bill Of Materials`)](ttps://www.argon.io/blog/guide-to-sbom/) without golang installed:

```bash
bgo sbom   # show sbom of bgo itself

# Show these executables
bgo sbom ~/go/bin/gopls ~/go/bin/golangci-lint
```

### Extensible Commands

The commands system of `bgo` can be extensible by editing `.bgo.yaml`.

#### Config Files of `bgo`

`bgo` loads all `bgo.yaml` as its configurations from these locations:

1. _Main_ Config Files:
   - `/etc/bgo/bgo.yml` and `conf.d/*.yml` in it
   - `/usr/local/etc/bgo/bgo.yml` and `conf.d/*.yml` in it
   Any of them will be loaded as default config set.
2. _Secondary_ Config Files:
   - `$HOME/.bgo/bgo.yml` and `conf.d/*.yml` in it
   - `$HOME/.config/bgo/bgo.yml` and `conf.d/*.yml` in it
   Any of them will be loaded as extensible config set.
3. _Alternative_ Config File:
   - `./.bgo.yaml`
   This config file is used to current building project by `bgo`.
   The scanning results of golang CLI apps from current directory will be written into this file.

#### Customize the extensible commands in Main/Secondary Config Files

The extensible commands from config file is a feature of `hedzr/cmdr`. It generally looks like:

<details>
  <summary> Expand to source codes of `80.aliases.yml` </summary>

```yaml
# ~/.config/bgo/conf.d/80.aliases.yml
app:
    aliases:
        # group:                                  # group-name (optional). such as: "别名".
        commands:
            # - title: list
            #   short-name: ls
            #   # aliases: []
            #   # name: ""
            #   invoke-sh: ls -la -G                # for macOS, -G = --color; for linux: -G = --no-group
            #   # invoke: "another cmdr command"
            #   # invoke-proc: "..." # same with invoke-sh
            #   desc: list the current directory

            # - title: echo
            #   invoke-sh: |
            #     # pwd
            #     echo "{{$flg := index .Cmd.Flags 0}}{{$path :=$flg.GetDottedNamePath}} {{$fullpath := .Store.Wrap $path}} {{$fullpath}} | {{.Store.GetString $fullpath}}"
            #   desc: print the name
            #   flags:
            #     - title: name
            #       default:              # default value
            #       type: string          # bool, string, duration, int, uint, ...
            #       group:
            #       toggle-group:
            #       desc: specify the name to be printed

            - title: check-code-qualities
              short-name: chk
                # aliases: [check]
                # name: ""
              invoke-sh: |
                  which golint || go install golang.org/x/lint/golint
                  which gocyclo || go install github.com/fzipp/gocyclo
                  echo
                  echo
                  echo "Command hit: {{.Cmd.GetDottedNamePath}}"
                  echo "fmt {{.ArgsString}} ..."
                  gofmt -l -s -w {{range .Args}}{{.}}{{end}}
                  echo "lint {{.ArgsString}} ..."
                  golint {{.ArgsString}}
                  echo "cyclo ."
                  gocyclo -top 13 .
                # invoke: "another cmdr command"
                # invoke-proc: "..." # same with invoke-sh
              shell: /usr/bin/env bash # optional, default is /bin/bash
              desc: pre-options before releasing. typically fmt,lint,cyclo,...

            - title: coverage
              short-name: cov
              invoke-sh: |
                  # pwd
                  # echo "{{$flg := index .Cmd.Flags 0}}{{$path := $flg.GetDottedNamePath}} {{$fullpath := .Store.Wrap $path}} {{$fullpath}} | {{.Store.GetString $fullpath}}"
                  # echo "{{$flg = index .Cmd.Flags 1}}{{$path = $flg.GetDottedNamePath}} {{$fullpath2 := .Store.Wrap $path}} {{$fullpath2}} | {{.Store.GetString $fullpath2}}"
                  # echo "{{$flg = index .Cmd.Flags 2}}{{$path = $flg.GetDottedNamePath}} {{$fullpath3 := .Store.Wrap $path}} {{$fullpath3}} | {{.Store.GetString $fullpath3}}"
                  go test -v -race \
                    -coverprofile={{.Store.GetString $fullpath2}} \
                    -covermode=atomic -timeout=20m \
                    {{if not .Args}}.{{else}}{{.ArgsString}}{{end}} \
                    | tee {{.Store.GetString $fullpath3}}
                  go tool cover \
                    -html={{.Store.GetString $fullpath2}} \
                    -o={{.Store.GetString $fullpath}}
              desc: run coverage, produce coverage.txt and cover.html
              examples: bgo cov ./... or bgo cov
              flags:
                  - title: html
                    default: cover.html # default value
                    default-placeholder: FILE
                    type: string # bool, string, duration, int, uint, ...
                    group:
                    toggle-group:
                    desc: specify the html filename to be output
                  - title: text
                    default: coverage.txt # default value
                    default-placeholder: FILE
                    type: string # bool, string, duration, int, uint, ...
                    group:
                    toggle-group:
                    desc: specify the coverprofile filename
                  - title: log
                    default: coverage.log # default value
                    default-placeholder: FILE
                    type: string # bool, string, duration, int, uint, ...
                    group:
                    toggle-group:
                    desc: specify the name to be logged
```

</details>

The config adds two top-level subcommands: `chk`(`check-code-qualities`) and `cov`(`coverage`)

Another sample is the bundled `91.more.aliases.yml`:

<details>
  <summary> Expand to source codes of `91.more.aliases.yml`</summary>

```yaml
app:
  aliases:
    # group:                                  # group-name (optional). such as: "别名".
    commands:
      # - title: list
      #   short-name: ls
      #   # aliases: []
      #   # name: ""
      #   invoke-sh: ls -la -G                # for macOS, -G = --color; for linux: -G = --no-group
      #   # invoke: "another cmdr command"
      #   # invoke-proc: "..." # same with invoke-sh
      #   desc: list the current directory

      - title: yolo
        short-name: y
        invoke-sh: |
          which yolo || { echo "installing yolo..." && go install -v github.com/azer/yolo; }
          yolo -i *.go -c 'go build' -a localhost:8080 {{.ArgsString}}
        desc: print the name

      - title: echo
        invoke-sh: |
          # pwd
          echo "{{$flg := index .Cmd.Flags 0}}{{$path :=$flg.GetDottedNamePath}} {{$fullpath := .Store.Wrap $path}} {{$fullpath}} | {{.Store.GetString $fullpath}}"
        desc: print the name
        flags:
          - title: name
            default: # default value
            type: string # bool, string, duration, int, uint, ...
            group:
            toggle-group:
            desc: specify the name to be printed
```

</details>

## Inspired By

`bgo` building subcommand was inspired by many tools:

- <https://github.com/mitchellh/gox>
- <https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5>
- <https://github.com/davecheney/golang-crosscompile>
- <https://github.com/laher/goxc>
- Makefile

## LICENSE

Apache 2.0
