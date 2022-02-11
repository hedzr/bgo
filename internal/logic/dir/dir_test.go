package dir

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestShortDir(t *testing.T) {
	x := "./study/a1"
	y := path.Clean(x)
	t.Log(y)
}

func TestForDirMax(t *testing.T) {
	err := ForDirMax("..", 0, 3, func(depth int, dirname string, fi os.FileInfo) (stop bool, err error) {
		println(depth, dirname, fi.Name(), fi.IsDir(), fi.Size(), fi.ModTime().Format(time.RFC3339Nano))
		return
	}, ".git", ".vscode", ".svn", "test*")
	if err != nil {
		t.Errorf("%v", err)
	}
}
