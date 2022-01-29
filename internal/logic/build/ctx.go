package build

import (
	"github.com/hedzr/log/dir"
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
