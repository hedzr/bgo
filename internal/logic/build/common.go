package build

//nolint:goimports //so what
import (
	"runtime"

	"github.com/hedzr/evendeep"
	"github.com/hedzr/evendeep/flags/cms"

	"github.com/hedzr/bgo/internal/logic/logx"
)

type (
	//nolint:lll //no
	CommonBase struct {
		OS   string `yaml:"-" json:"-" toml:"-"` // just for string template expansion
		ARCH string `yaml:"-" json:"-" toml:"-"` // just for string template expansion

		Args       []string `yaml:"args,omitempty,flow" json:"args,omitempty" toml:"args,omitempty"`
		Ldflags    []string `yaml:"ldflags,omitempty,flow" json:"ldflags,omitempty" toml:"ldflags,omitempty"`          // default ldflags is to get the smaller build for releasing
		Asmflags   []string `yaml:"asmflags,omitempty,flow" json:"asmflags,omitempty" toml:"asmflags,omitempty"`       //
		Gcflags    []string `yaml:"gcflags,omitempty,flow" json:"gcflags,omitempty" toml:"gcflags,omitempty"`          //
		Gccgoflags []string `yaml:"gccgoflags,omitempty,flow" json:"gccgoflags,omitempty" toml:"gccgoflags,omitempty"` //
		Tags       []string `yaml:"tags,omitempty,flow" json:"tags,omitempty" toml:"tags,omitempty"`                   //
		// Cgo option
		Cgo bool `yaml:",omitempty" json:"cgo,omitempty" toml:"cgo,omitempty"` //
		// Race option enables data race detection.
		//		Supported only on linux/amd64, freebsd/amd64, darwin/amd64, windows/amd64,
		//		linux/ppc64le and linux/arm64 (only for 48-bit VMA).
		Race bool `yaml:",omitempty" json:"race,omitempty" toml:"race,omitempty"` //
		// Msan option enables interoperation with memory sanitizer.
		//		Supported only on linux/amd64, linux/arm64
		//		and only with Clang/LLVM as the host C compiler.
		//		On linux/arm64, pie build mode will be used.
		Msan bool `yaml:",omitempty" json:"msan,omitempty" toml:"msan,omitempty"` //
		// Asan enables interoperation with address sanitizer.
		//		Supported only on linux/arm64, linux/amd64.
		//		Supported only on linux/amd64 or linux/arm64 and only with GCC 7 and higher
		//		or Clang/LLVM 9 and higher.
		Asan bool `yaml:",omitempty" json:"asan,omitempty" toml:"asan,omitempty"` //
		//
		Mod           string `yaml:",omitempty" json:"mod,omitempty" toml:"mod,omitempty"`                                     // -mod reaonly,vendor,mod
		Amd64         string `yaml:",omitempty" json:"goamd64,omitempty" toml:"goamd64,omitempty"`                             // GOAMD64 v1,v2,v3,v4
		Gocmd         string `yaml:",omitempty" json:"gocmd,omitempty" toml:"gocmd,omitempty"`                                 // -gocmd go
		Gen           bool   `yaml:",omitempty" json:"gen,omitempty" toml:"gen,omitempty"`                                     // go generate at first?
		Install       bool   `yaml:",omitempty" json:"install,omitempty" toml:"install,omitempty"`                             // install binary to $GOPATH/bin like 'go install' ?
		Debug         bool   `yaml:",omitempty" json:"debug,omitempty" toml:"debug,omitempty"`                                 // true to produce a larger build with debug info
		DisableResult bool   `yaml:"disable-result,omitempty" json:"disable_result,omitempty" toml:"disable_result,omitempty"` // no ll (Shell list) building result

		// -X for -ldflags,
		// -X importpath.name=value
		//    Set the value of the string variable in importpath named name to value.
		//    Note that before Go 1.5 this option took two separate arguments.
		//    Now it takes one argument split on the first = sign.
		Extends      []PackageNameValues `yaml:"extends,omitempty" json:"extends,omitempty" toml:"extends,omitempty"` //
		CmdrSpecials bool                `yaml:"cmdr,omitempty" json:"cmdr,omitempty" toml:"cmdr,omitempty"`
		CmdrVersion  string              `yaml:"cmdr-version,omitempty" json:"cmdr-version,omitempty" toml:"cmdr-version,omitempty"`
	}

	// Upx params package.
	Upx struct {
		Enable bool     `yaml:",omitempty" json:"enable,omitempty" toml:"enable,omitempty"`
		Params []string `yaml:",omitempty" json:"params,omitempty" toml:"params,omitempty"`
	}

	PackageNameValues struct {
		Package string            `yaml:"pkg,omitempty" json:"package,omitempty" toml:"package,omitempty"`
		Values  map[string]string `yaml:"values,omitempty" json:"values,omitempty" toml:"values,omitempty"`
	}

	//nolint:lll //no
	Common struct {
		CommonBase `yaml:",omitempty,inline,flow" json:",omitempty" toml:""`

		Disabled       bool     `yaml:"disabled,omitempty" json:"disabled,omitempty" toml:"disabled,omitempty"`
		KeepWorkdir    bool     `yaml:"keep-workdir,omitempty" json:"keep_workdir,omitempty" toml:"keep_workdir,omitempty"`
		UseWorkDir     string   `yaml:"use-workdir,omitempty" json:"use_work_dir,omitempty" toml:"use_work_dir,omitempty"`
		For            []string `yaml:"for,omitempty,flow" json:"for,omitempty" toml:"for,omitempty"`
		Os             []string `yaml:"os,omitempty,flow" json:"os,omitempty" toml:"os,omitempty"`
		Arch           []string `yaml:"arch,omitempty,flow" json:"arch,omitempty" toml:"arch,omitempty"`
		Goroot         string   `yaml:"goroot,omitempty,flow" json:"goroot,omitempty" toml:"goroot,omitempty"`
		PreAction      string   `yaml:"pre-action,omitempty" json:"pre_action,omitempty" toml:"pre_action,omitempty"`                   // bash script
		PostAction     string   `yaml:"post-action,omitempty" json:"post_action,omitempty" toml:"post_action,omitempty"`                // bash script
		PreActionFile  string   `yaml:"pre-action-file,omitempty" json:"pre_action_file,omitempty" toml:"pre_action_file,omitempty"`    // bash script
		PostActionFile string   `yaml:"post-action-file,omitempty" json:"post_action_file,omitempty" toml:"post_action_file,omitempty"` // bash script

		DeepReduce bool `yaml:"reduce,omitempty" json:"reduce,omitempty" toml:"reduce,omitempty"`
		Upx        Upx  `yaml:",omitempty" json:"upx,omitempty" toml:"upx,omitempty"`
	}
)

func NewCommon() *Common {
	return &Common{
		Disabled: false,
		For:      nil,
		// Os:            nil,
		// Arch:          nil,
		PreAction:  "",
		PostAction: "",
		// Ldflags:       "",
		// Asmflags:      nil,
		// Gcflags:       nil,
		// Tags:          nil,
		// Cgo:           false,

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
			Asmflags:      nil, // []string{"-trimpath=$GOPATH"},
			Gcflags:       nil, // []string{"-trimpath=$GOPATH"},
			Ldflags:       []string{},
		},
	}
}

func (c *CommonBase) CloneFrom(from *CommonBase) {
	cc := evendeep.New(
		evendeep.WithCopyStrategyOpt,
		evendeep.WithStrategies(cms.OmitIfEmpty),
	)
	// cc := cmdr.StandardCopier
	// cc.KeepIfFromIsNil = true
	// cc.KeepIfFromIsZero = true
	// cc.EachFieldAlways = true
	if err := cc.CopyTo(from, &c); err != nil {
		logx.Error("CommonBase.CloneFrom failed: %v", err)
	}
}

func (c *Common) CloneFrom(from *Common) {
	cc := evendeep.New(
		evendeep.WithCopyStrategyOpt,
		evendeep.WithStrategies(cms.OmitIfEmpty),
	)
	// cc := cmdr.StandardCopier
	// cc.KeepIfFromIsNil = true
	// cc.KeepIfFromIsZero = true
	// cc.EachFieldAlways = true
	if err := cc.CopyTo(from, &c); err != nil {
		logx.Error("Common.CloneFrom failed: %v", err)
	}
}

func (c *Common) MergeFrom(from *Common) {
	cc := evendeep.New(
		evendeep.WithMergeStrategyOpt,
		evendeep.WithStrategies(cms.OmitIfEmpty),
	)
	// cc := cmdr.StandardCopier
	// cc.KeepIfFromIsNil = true
	// cc.KeepIfFromIsZero = true
	// cc.EachFieldAlways = true
	if err := cc.CopyTo(from, &c); err != nil {
		logx.Error("Common.CloneFrom failed: %v", err)
	}
}
