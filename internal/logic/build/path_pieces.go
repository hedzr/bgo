package build

//nolint:goimports //so what
import (
	"path"
	"path/filepath"

	"github.com/hedzr/log/dir"
)

type PathPieces struct {
	Path    string
	Dir     string
	Base    string
	Ext     string
	AbsPath string
}

func (s *PathPieces) Set(file string) {
	s.Path = file
	s.Dir = path.Dir(file)
	s.Base = path.Base(file)
	s.Ext = filepath.Ext(file)
	s.AbsPath = dir.AbsPath(file)
}
