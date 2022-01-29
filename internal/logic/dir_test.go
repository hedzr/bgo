package logic

import (
	"path"
	"testing"
)

func TestShortDir(t *testing.T) {
	x := "./study/a1"
	y := path.Clean(x)
	t.Log(y)
}
