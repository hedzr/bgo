package logic

import "github.com/hedzr/bgo/internal/logic/build"

//const BgoYamlFilename = ".bgo.yml"

type (
	ProjectWrap struct {
		Project `yaml:",omitempty,inline,flow" json:"project" toml:"project,omitempty"`
	}

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

func newProjectClean(pkg string, dir string) *ProjectWrap {
	p := &ProjectWrap{Project: Project{
		Name:    "",
		Dir:     dir,
		Package: pkg,
		MainGo:  "",
		Common:  build.NewCommon(),
	}}
	return p
}

func (p *Project) apply(packageName string, bs *BgoSettings) {
	if p.Common == nil {
		p.Common = build.NewCommon()
	}
	p.Common.CloneFrom(bs.Common)
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
