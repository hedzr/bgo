package build

//nolint:goimports //so what
import (
	"os"

	"github.com/hedzr/log/closers"

	"github.com/hedzr/bgo/internal/logx"
)

var TempDir string //nolint:gochecknoglobals //no

//nolint:gochecknoinits //so what
func init() {
	var err error
	TempDir, err = os.MkdirTemp("", "bgo")
	if err != nil {
		logx.Error("%v", err)
		return
	}

	// hedzr/cmdr provides a mechanism to close all closers automatically.
	closers.RegisterCloseFns(func() {
		_ = os.RemoveAll(TempDir)
	})

	logx.DimV("Using temp dir: %q", TempDir)
}
