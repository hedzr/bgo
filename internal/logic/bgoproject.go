package logic

import (
	"github.com/hedzr/bgo/internal/logic/build"
	"github.com/hedzr/bgo/internal/tool"
)

// const BgoYamlFilename = ".bgo.yml"

type (
	ProjectWrap struct {
		Project `yaml:",omitempty,inline,flow" json:"project" toml:"project,omitempty"`
	}
	//nolint:lll //no
	Project struct {
		Name    string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`          // appName if specified
		Dir     string `yaml:"dir" json:"dir" toml:"dir"`                                           // root dir of this module/cli-app
		Package string `yaml:"package,omitempty" json:"package,omitempty" toml:"package,omitempty"` // pkgName or mainGo
		MainGo  string `yaml:"main-go,omitempty" json:"main_go,omitempty" toml:"main_go,omitempty"` // default: 'main.go', // or: "./cli", "./cli/main.go", ...
		Version string `yaml:"version,omitempty" json:"version,omitempty" toml:"version,omitempty"` //

		// In current and earlier releases, MainGo is reserved and never used.

		*build.Common `yaml:",omitempty,inline,flow" json:",omitempty" toml:",omitempty"`

		// Cmdr building opts

		// More Posts:
		//   packaging methods
		//   coverage
		//   uploading
		//   docker
		//   github release
		// Before
		//   fmt, lint, test, cyclo, ...

		keyName string
		tp      *build.TargetPlatforms
	}
)

func newProject(pkg string, bc *build.Context, bs *BgoSettings) *ProjectWrap {
	p := &ProjectWrap{Project: Project{
		Name:    "",
		Dir:     bc.WorkDir,
		Package: pkg,
		MainGo:  "",
		// Common:  *bs.Common,
	},
	}
	p.apply("", bs)
	p.Common.KeepWorkdir = true
	return p
}

func newProjectClean(pkg, dir string) *ProjectWrap {
	p := &ProjectWrap{Project: Project{
		Name:    "",
		Dir:     dir,
		Package: pkg,
		MainGo:  "",
		Common:  build.NewCommon(),
	}}
	return p
}

func (p *Project) GetTitleName() string {
	if p.Name != "" {
		return p.Name
	}
	if p.keyName != "" {
		return tool.StripOrderPrefix(p.keyName)
	}
	return "Name?"
}

func (p *Project) apply(packageName string, bs *BgoSettings) {
	if p.Common == nil {
		p.Common = build.NewCommon()
	}
	// logx.Log(`Project.apply(bs.Common)`)
	p.Common.MergeFrom(bs.Common)
	if p.Package == "" && packageName != "" {
		p.Package = packageName
	}
}

func (p *Project) applyPI(pi *pkgInfo) {
	if p.Package == "" && pi.packageName != "" {
		p.Package = pi.packageName
	}
	if p.Name == "" {
		if pi.appName != "" {
			p.Name = pi.appName
		} else if pi.projectName != "" {
			p.Name = pi.projectName
		}
	}
	if p.Dir != pi.dirname {
		p.Dir = pi.dirname
	}
}

func (p *Project) inIntSlice(val int, slice []int) (yes bool) { //nolint:unused //no
	for _, v := range slice {
		if yes = val == v; yes {
			break
		}
	}
	return
}

func (p *Project) inStringSlice(val string, slice []string) (yes bool) {
	for _, v := range slice {
		if yes = val == v; yes {
			break
		}
	}
	return
}

func (p *Project) overspreadByTP(scope string, tpBase *build.TargetPlatforms) {
	for os, oss := range tpBase.OsArchMap {
		for arch := range oss {
			if !p.inStringSlice(os, p.Os) {
				p.Os = append(p.Os, os)
			}
			if !p.inStringSlice(arch, p.Arch) {
				p.Arch = append(p.Arch, arch)
			}
		}
	}

	p.tp = build.NewTargetPlatforms()
	if err := p.tp.Init(); err != nil {
		return
	}
	p.tp.FilterBy(scope, p.For, p.Os, p.Arch)
}
