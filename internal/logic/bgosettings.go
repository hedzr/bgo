package logic

import (
	"github.com/hedzr/bgo/internal/logic/build"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/evendeep"
)

type (
	// BgoSettings for go build.
	// Generally it's loaded from '.bgo.yml'.
	BgoSettings struct {
		*build.Common `yaml:",omitempty,inline" json:",omitempty" toml:",omitempty"`

		Scope    string                  `yaml:"scope,omitempty" json:"scope,omitempty" toml:"scope,omitempty"`
		Projects map[string]ProjectGroup `yaml:"projects" json:"projects" toml:"projects"`

		Output   Output   `yaml:",omitempty" json:"output" toml:"output,omitempty"`
		Excludes []string `yaml:",omitempty,flow" json:"excludes,omitempty" toml:"excludes,omitempty"`
		Goproxy  string   `yaml:",omitempty" json:"goproxy,omitempty" toml:"goproxy,omitempty"`

		SavedAs []string `yaml:"saved-as,omitempty" json:"saved_as,omitempty" toml:"saved_as,omitempty"`
	}

	Output struct {
		Dir         string `yaml:",omitempty" json:"dir,omitempty" toml:"dir,omitempty"`                   // Base Dir
		SplitTo     string `yaml:"split-to,omitempty" json:"split_to,omitempty" toml:"split_to,omitempty"` // see build.Context, Group, Project, ... | see also BuildContext struct
		NamedAs     string `yaml:"named-as,omitempty" json:"named_as,omitempty" toml:"named_as,omitempty"` // see build.Context,
		SuffixAs    string `yaml:"suffix-as,omitempty" json:"-" toml:"-"`                                  // NEVER USED.
		ZipSuffixAs string `yaml:"zip-suffix-as,omitempty" json:"-" toml:"-"`                              // NEVER USED. supported: gz/tgz, bz2, xz/txz, 7z
	}

	ProjectGroup struct {
		LeadingText   string                  `yaml:"leading-text,omitempty" json:"leading_text,omitempty" toml:"leading_text,omitempty"`
		Items         map[string]*ProjectWrap `yaml:"items" json:"items,omitempty" toml:"items,omitempty"`
		*build.Common `yaml:",omitempty,inline" json:",omitempty" toml:",omitempty"`
	}
)

func newBgoSettings(buildScope string) *BgoSettings {
	bs := new(BgoSettings)
	bs.Common = build.NewCommon()
	bs.Scope = buildScope
	bs.Excludes = []string{"study*", "test*"}
	bs.Output.Dir = "./bin"
	bs.Output.NamedAs = "{{.AppName}}-{{.OS}}-{{.ARCH}}"

	// get build section in .bgo.yml
	// .bgo.yml will be loaded automatically as a cmdr feature
	err := cmdr.GetSectionFrom("bgo.build", &bs)

	bs.initFromCmdr()

	if err == nil {
		bs.PullDownCommonSettings()
	}

	// logHiLight("Starting...")

	// logx.Dim("- BS:")
	// logx.Dim("%v", leftPad(yamlText(bs.Projects), 2))
	// os.Exit(0)

	// cmdr.DebugOutputTildeInfo(true)

	return bs
}

func (s *BgoSettings) initFromCmdr() {
	if s.SavedAs == nil {
		s.SavedAs = []string{"bgo.yml"}
	}
	s.SavedAs = cmdr.GetStringSliceR("settings-filename", s.SavedAs...)

	if str := cmdr.GetStringR("build.output"); str != "" {
		s.Output.NamedAs = str
	}

	if str := cmdr.GetStringR("build.gocmd"); str != "" {
		s.Gocmd = str
	}

	s._init1()
	s._init2()
}

func (s *BgoSettings) _init1() {
	if cmdr.GetBoolR("build.race") {
		s.Race = true
	}

	if cmdr.GetBoolR("build.cgo") {
		s.Cgo = true
	}

	if cmdr.GetBoolR("build.msan") {
		s.Msan = true
	}

	if slice := cmdr.GetStringSliceR("build.ldflags"); len(slice) > 0 {
		s.Ldflags = append(s.Ldflags, slice...)
	}
	if slice := cmdr.GetStringSliceR("build.gcflags"); len(slice) > 0 {
		s.Gcflags = append(s.Gcflags, slice...)
	}
	if slice := cmdr.GetStringSliceR("build.asmflags"); len(slice) > 0 {
		s.Asmflags = append(s.Asmflags, slice...)
	}
	if slice := cmdr.GetStringSliceR("build.gccgoflags"); len(slice) > 0 {
		s.Gccgoflags = append(s.Gccgoflags, slice...)
	}
	if slice := cmdr.GetStringSliceR("build.tags"); len(slice) > 0 {
		s.Tags = append(s.Tags, slice...)
	}
}

func (s *BgoSettings) _init2() {
	var manually, os, arch bool
	if slice := cmdr.GetStringSliceR("build.os"); len(slice) > 0 {
		// s.Os = append(s.Os, slice...)
		manually, os, s.Os = true, true, slice
	} else if s.Scope != "full" && len(s.Os) == 0 { //nolint:goconst //i like it
		s.Os = []string{"linux", "darwin", "windows"}
	}
	if slice := cmdr.GetStringSliceR("build.osarch"); len(slice) > 0 {
		// s.Arch = append(s.Arch, slice...)
		manually, arch, s.Arch = true, true, slice
	} else if s.Scope != "full" && len(s.Arch) == 0 {
		s.Arch = []string{"amd64"}
	}

	if manually {
		if os && !arch {
			s.Arch = nil
		} else if arch && !os {
			s.Os = nil
		}
	}

	if slice := cmdr.GetStringSliceR("build.for"); len(slice) > 0 {
		s.For = slice
		if !manually {
			s.Os = nil
			s.Arch = nil
		}
	}
}

func (s *BgoSettings) PullDownCommonSettings() {
	c := evendeep.New(evendeep.WithCopyStrategyOpt, evendeep.WithWipeTargetSliceFirstOpt)
	// c := cmdr.StandardCopier // make a copy
	// c.KeepIfFromIsNil = true
	// c.KeepIfFromIsZero = true
	// c.EachFieldAlways = true
	// c.IgnoreIfNotEqual = true

	// Copying the common settings from the up level if necessary.
	// It means:
	//   the top level (app.bgo.build.xxx) -> project.common
	//   the project-group level -> project.common
	for gn, grp := range s.Projects {
		for pn, prj := range grp.Items {
			if prj.Common == nil {
				prj.Common = build.NewCommon()
			}
			_ = c.CopyTo(s.Common, prj.Common)
			_ = c.CopyTo(grp.Common, prj.Common)
			s.Projects[gn].Items[pn] = prj
		}
	}
}
