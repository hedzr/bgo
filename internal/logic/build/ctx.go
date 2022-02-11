package build

import (
	"github.com/hedzr/bgo/internal/logic/tool"
	"github.com/hedzr/log/dir"
	"path"
)

type (
	Context struct {
		WorkDir    string
		TempDir    string
		PackageDir string

		// Output collects the target binary executable path pieces in building
		Output PathPieces

		*Common
		*Info
		*DynBuildInfo
	}
)

func NewContext() *Context {
	return &Context{
		WorkDir: dir.GetCurrentDir(),
		TempDir: TempDir,

		Common:       NewCommon(),
		Info:         NewBuildInfo(),
		DynBuildInfo: newDynBuildInfo(),
	}
}

func (bc *Context) FindAppName(knownName, knownProjectName, knownPackageName string) {
	bc.AppName, bc.ProjectName = knownName, knownProjectName
	if bc.ProjectName == "" {
		bc.ProjectName = path.Base(knownPackageName)
	}
	if bc.AppName == "" {
		bc.AppName = bc.ProjectName
	}
	bc.ProjectName = tool.StripOrderPrefix(bc.ProjectName)
	bc.AppName = tool.StripOrderPrefix(bc.AppName)
}
