package logic

import (
	"github.com/hedzr/bgo/internal/logic/build"
	"github.com/hedzr/cmdr"
)

type (
	// BgoSettings for go build.
	// Generally it's loaded from '.bgo.yml'.
	BgoSettings struct {
		*build.Common `yaml:",omitempty,inline"`

		Scope    string                  `yaml:"scope,omitempty"`
		Projects map[string]ProjectGroup `yaml:"projects"`

		Output   Output   `yaml:",omitempty"`
		Excludes []string `yaml:",omitempty,flow"`
		Goproxy  string   `yaml:",omitempty"`

		SavedAs string `yaml:"saved-as,omitempty"`
	}

	Output struct {
		Dir         string `yaml:",omitempty"`              // Base Dir
		SplitTo     string `yaml:"split-to,omitempty"`      // see build.Context, Group, Project, ... | see also BuildContext struct
		NamedAs     string `yaml:"named-as,omitempty"`      // see build.Context,
		SuffixAs    string `yaml:"suffix-as,omitempty"`     // NEVER USED.
		ZipSuffixAs string `yaml:"zip-suffix-as,omitempty"` // NEVER USED. supported: gz/tgz, bz2, xz/txz, 7z
	}

	ProjectGroup struct {
		LeadingText   string                  `yaml:"leading-text,omitempty"`
		Items         map[string]*ProjectWrap `yaml:"items"`
		*build.Common `yaml:",omitempty,inline"`
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

	if bs.SavedAs == "" {
		bs.SavedAs = "bgo.yml"
	}
	bs.SavedAs = cmdr.GetStringR("settings-filename", bs.SavedAs)

	if str := cmdr.GetStringR("build.gocmd"); str != "" {
		bs.Gocmd = str
	}

	if cmdr.GetBoolR("build.race") {
		bs.Race = true
	}

	if cmdr.GetBoolR("build.cgo") {
		bs.Cgo = true
	}

	if cmdr.GetBoolR("build.msan") {
		bs.Msan = true
	}

	if slice := cmdr.GetStringSliceR("build.ldflags"); len(slice) > 0 {
		bs.Ldflags = append(bs.Ldflags, slice...)
	}
	if slice := cmdr.GetStringSliceR("build.gcflags"); len(slice) > 0 {
		bs.Gcflags = append(bs.Gcflags, slice...)
	}
	if slice := cmdr.GetStringSliceR("build.asmflags"); len(slice) > 0 {
		bs.Asmflags = append(bs.Asmflags, slice...)
	}
	if slice := cmdr.GetStringSliceR("build.gccgoflags"); len(slice) > 0 {
		bs.Gccgoflags = append(bs.Gccgoflags, slice...)
	}
	if slice := cmdr.GetStringSliceR("build.tags"); len(slice) > 0 {
		bs.Tags = append(bs.Tags, slice...)
	}

	var manually, os, arch bool
	if slice := cmdr.GetStringSliceR("build.os"); len(slice) > 0 {
		//bs.Os = append(bs.Os, slice...)
		manually, os, bs.Os = true, true, slice
	} else if bs.Scope != "full" && len(bs.Os) == 0 {
		bs.Os = []string{"linux", "darwin", "windows"}
	}
	if slice := cmdr.GetStringSliceR("build.osarch"); len(slice) > 0 {
		//bs.Arch = append(bs.Arch, slice...)
		manually, arch, bs.Arch = true, true, slice
	} else if bs.Scope != "full" && len(bs.Arch) == 0 {
		bs.Arch = []string{"amd64"}
	}

	if manually {
		if os && !arch {
			bs.Arch = nil
		} else if arch && !os {
			bs.Os = nil
		}
	}

	if slice := cmdr.GetStringSliceR("build.for"); len(slice) > 0 {
		bs.For = slice
		if manually == false {
			bs.Os = nil
			bs.Arch = nil
		}
	}

	if str := cmdr.GetStringR("build.output"); str != "" {
		bs.Output.NamedAs = str
	}

	if err == nil {
		bs.PullDownCommonSettings()
	}

	//logHiLight("Starting...")

	//logx.Dim("- BS:")
	//logx.Dim("%v", leftPad(yamlText(bs.Projects), 2))
	// os.Exit(0)

	//cmdr.DebugOutputTildeInfo(true)

	return bs
}

func (s *BgoSettings) PullDownCommonSettings() {
	c := cmdr.StandardCopier // make a copy
	c.KeepIfFromIsNil = true
	c.KeepIfFromIsZero = true
	c.EachFieldAlways = true
	c.IgnoreIfNotEqual = true

	// Copying the common settings from the up level if necessary.
	// It means:
	//   the top level (app.bgo.build.xxx) -> project.common
	//   the project-group level -> project.common
	for gn, grp := range s.Projects {
		for pn, prj := range grp.Items {
			if prj.Common == nil {
				prj.Common = build.NewCommon()
			}
			_ = c.Copy(prj.Common, s.Common)
			_ = c.Copy(prj.Common, grp.Common)
			s.Projects[gn].Items[pn] = prj
		}
	}
}
