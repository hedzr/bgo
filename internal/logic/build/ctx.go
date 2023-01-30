package build

//nolint:goimports //so what
import (
	"path"

	"github.com/hedzr/bgo/internal/logic/tool"
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

func (bc *Context) CalcVersion() (ver string) {
	ver = bc.Version
	if bc.GitRevision[0] == 'v' { //nolint:gocritic //no
		ver = bc.GitRevision
	} else if bc.GitVersion[0] == 'v' {
		ver = bc.GitVersion
	} else if bc.GitSummary[0] == 'v' {
		ver = bc.GitSummary
	}
	return
}
