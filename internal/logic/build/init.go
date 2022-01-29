package build

import (
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/log/closers"
	"io/ioutil"
	"os"
)

var TempDir string

func init() {
	var err error
	TempDir, err = ioutil.TempDir("", "bgo")
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
