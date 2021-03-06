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

![tip-for-bgo](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20220202111546956.png)

## Features

`bgo` makes your golang building life easier, it's an efficent and extensible build tool.

-   Run go building with or without a config file `.bgo.yml`
-   Scan the directory to grab all main packages and initial `.bgo.yml`
-   While

    -   you have lots of small cli apps in many sub-directories
    -   one project wants go1.18beta1 and others stays in go1.17
    -   too much `-tags`, `-ldflags`, `-asmflags`, ...
    -   multiple targets, cross compiling
    -   etc.

    have a try with `bgo`.

## History

- v0.3.23

    - add support to GOAMD64

- v0.3.21

    -   add SBOM support: `bgo sbom <executable>`
    -   what is SBOM: [here](https://www.argon.io/blog/guide-to-sbom/), and [here](https://blog.sonatype.com/sbom-from-the-idea-of-transparency-to-the-reality-of-code).

- v0.3.18

    -   code reviewd

- v0.3.17

    -   upgrade log,errors

- v0.3.15

    -   fix: version command and help screen not work

- v0.3.13

    -   fix: script-file not work (pre/post-action-file)
    -   fix: --auto/--short get wrong target platforms matrix sometimes
    -   fix; cross-platform ls
    -   imp: `bgo` has been adapted onto **Windows 11** (should work for winx too)
    -   imp: more shell completions supported (bash, zsh, fish-shell, etc)
    -   fix/imp: .bgo.yml for myself has wrong params (such as .randomString in extends section, or wrong githash, version)

- v0.3.12

    -   temporary build
    -   fea: `bgo init -o bgo.{yml,yaml,json,toml}`
        -   imp: optimized json and toml outputting
        -   fea: support bgo init multiple outputs once: `bgo init -o=a.{yml,toml}`, an inessential feature
        -   imp: better json, toml outputting
    -   fix: the wrong template expansion in post/preAction was covered silently; and fixed the typo in postAction
    -   .bgo.yml: bgo - linux+darwin in auto mode.
    -   args:
        -   move `--dry-run` up to root level
        -   fix: buildtags might not work
        -   **fixed**: `-os` `-arch` and more `build` options cannot work in root command level  
            **TEMP WORKAROUND for older versions**  
            uses full path command `bgo build -os linux` instead of `bgo -os linux` till our new release arrive.
    -   imp: build - rewrite loopAllPackages, enable leadingText field
    -   imp: logx - LazyInit() and better performance
    -   imp: review codes

- v0.3.3

    -   fea: **Aliases** definitions in primary config directory can be merged into `bgo` command system now
        -   fea: `check-code-qualities` alias command added and play `gofmt`, `golint` and `golint` at once.
        -   fea: Extend `bgo` command system with Aliases definitions.
    -   fea: `bgo init -o bgo.{yml,yaml,json,toml}` writes different config file formats with giving suffix
    -   fix: TargetPlatforms.FilterBy not very ok
    -   imp: added cmdr global pre-action: verbose info for debugging
    -   CHANGE: `.bgo.yml` is loaded as an alternative config file now
    -   CHANGE: `$HOME/.bgo` and `conf.d` subdirectory is primary config files now
    -   CHANGE: primary config files will be distributed with binary executable

- v0.3.0

    -   fix: `bgo init` not work
    -   fix: zsh completion file not fully written

- v0.2.17 and older
    -   pre-releases

## Getting Started

### Install

Download the prebuilt binaries from Release page.

Or Built from source code:

```bash
git clone http://github.com/hedzr/bgo.git
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
bgo
```

![image-20220128104835837](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20220128104835837.png)

Filter the target systems by `-for OS/ARCH`, `-os OS` and `-arch ARCH`:

```bash
bgo --for linux/386 -for linux/amd64,darwin/arm64
bgo -os linux -arch 386 -arch amd64 -arch arm64
```

Both long and short options are available for `for`, `os` and `arch`.

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

Run `bgo gen sh --bash -o=bgo.bash` to get bash completions script and put it to the right location. Generally it should be:

```bash
bgo generate shell --bash -o=bgo.bash
mv bgo.bash /etc/bash-completion.d/bgo
```

Nothing needs to do if installed via brew (since v0.3.3+).

## Inspired By:

-   https://github.com/mitchellh/gox
-   https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5
-   https://github.com/davecheney/golang-crosscompile
-   https://github.com/laher/goxc
-   Makefile

## LICENSE

Apache 2.0
