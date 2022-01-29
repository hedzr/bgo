package build

import (
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/cmdr"
	"runtime"
)

type (
	CommonBase struct {
		OS   string `yaml:"-"` // just for string template expansion
		ARCH string `yaml:"-"` // just for string template expansion

		Ldflags    []string `yaml:"ldflags,omitempty,flow"`    // default ldflags is to get the smaller build for releasing
		Asmflags   []string `yaml:"asmflags,omitempty,flow"`   //
		Gcflags    []string `yaml:"gcflags,omitempty,flow"`    //
		Gccgoflags []string `yaml:"gccgoflags,omitempty,flow"` //
		Tags       []string `yaml:"tags,omitempty,flow"`       //
		// Cgo option
		Cgo bool `yaml:",omitempty"` //
		// Race option enables data race detection.
		//		Supported only on linux/amd64, freebsd/amd64, darwin/amd64, windows/amd64,
		//		linux/ppc64le and linux/arm64 (only for 48-bit VMA).
		Race bool `yaml:",omitempty"` //
		// Msan option enables interoperation with memory sanitizer.
		//		Supported only on linux/amd64, linux/arm64
		//		and only with Clang/LLVM as the host C compiler.
		//		On linux/arm64, pie build mode will be used.
		Msan          bool   `yaml:",omitempty"`               //
		Gocmd         string `yaml:",omitempty"`               // -gocmd go
		Gen           bool   `yaml:",omitempty"`               // go generate at first?
		Install       bool   `yaml:",omitempty"`               // install binary to $GOPATH/bin like 'go install' ?
		Debug         bool   `yaml:",omitempty"`               // true to produce a larger build with debug info
		DisableResult bool   `yaml:"disable-result,omitempty"` // no ll (Shell list) building result

		// -X for -ldflags,
		// -X importpath.name=value
		//    Set the value of the string variable in importpath named name to value.
		//    Note that before Go 1.5 this option took two separate arguments.
		//    Now it takes one argument split on the first = sign.
		Extends      []PackageNameValues `yaml:"extends,omitempty"` //
		CmdrSpecials bool                `yaml:"cmdr,omitempty"`
	}

	PackageNameValues struct {
		Package string            `yaml:"pkg,omitempty"`
		Values  map[string]string `yaml:"values,omitempty"`
	}

	Common struct {
		CommonBase     `yaml:"base,omitempty,inline,flow"`
		Disabled       bool     `yaml:"disabled,omitempty"`
		KeepWorkdir    bool     `yaml:"keep-workdir,omitempty"`
		UseWorkDir     string   `yaml:"use-workdir,omitempty"`
		For            []string `yaml:"for,omitempty,flow"`
		Os             []string `yaml:"os,omitempty,flow"`
		Arch           []string `yaml:"arch,omitempty,flow"`
		Goroot         string   `yaml:"goroot,omitempty,flow"`
		PreAction      string   `yaml:"pre-action,omitempty"`       // bash script
		PostAction     string   `yaml:"post-action,omitempty"`      // bash script
		PreActionFile  string   `yaml:"pre-action-file,omitempty"`  // bash script
		PostActionFile string   `yaml:"post-action-file,omitempty"` // bash script
	}
)

func NewCommon() *Common {
	return &Common{
		Disabled: false,
		For:      nil,
		//Os:            nil,
		//Arch:          nil,
		PreAction:  "",
		PostAction: "",
		//Ldflags:       "",
		//Asmflags:      nil,
		//Gcflags:       nil,
		//Tags:          nil,
		//Cgo:           false,

		Os:   []string{runtime.GOOS},
		Arch: []string{runtime.GOARCH},
		CommonBase: CommonBase{
			Gen:           false,
			Install:       false,
			Debug:         false,
			DisableResult: false,
			Cgo:           false,
			Race:          false,
			Tags:          nil,
			Asmflags:      nil, //[]string{"-trimpath=$GOPATH"},
			Gcflags:       nil, //[]string{"-trimpath=$GOPATH"},
			Ldflags:       []string{},
		},
	}
}

func (c *CommonBase) CloneFrom(from *CommonBase) {
	cc := cmdr.StandardCopier
	cc.KeepIfFromIsNil = true
	cc.KeepIfFromIsZero = true
	cc.EachFieldAlways = true
	if err := cc.Copy(c, from); err != nil {
		logx.Error("CommonBase.CloneFrom failed: %v", err)
	}
}

func (c *Common) CloneFrom(from *Common) {
	cc := cmdr.StandardCopier
	cc.KeepIfFromIsNil = true
	cc.KeepIfFromIsZero = true
	cc.EachFieldAlways = true
	if err := cc.Copy(c, from); err != nil {
		logx.Error("Common.CloneFrom failed: %v", err)
	}
}
