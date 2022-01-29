package logic

import (
	"github.com/hedzr/log"
	"github.com/hedzr/log/dir"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func relativeDir(dirname, base string) string {
	b := dir.NormalizeDir(base)
	c := strings.Split(b, string(filepath.Separator))
	var pre []string
	for _, _ = range c {
		pre = append(pre, "..")
	}
	pre = append(pre, dirname)
	return path.Join(pre...)
}

// ForDir walks on `root` directory and its children
// fix for log/dir.ForDir
func ForDir(
	root string,
	cb func(depth int, dirname string, fi os.FileInfo) (stop bool, err error),
	excludes ...string,
) (err error) {
	err = ForDirMax(root, 0, -1, cb, excludes...)
	return
}

// ForDirMax walks on `root` directory and its children with nested levels up to `maxLength`.
//
// Example - discover folder just one level
//
//      _ = ForDirMax(dir, 0, 1, func(depth int, dirname string, fi os.FileInfo) (stop bool, err error) {
//			if fi.IsDir() {
//				return
//			}
//          // ... doing something for a file,
//			return
//		})
//
// maxDepth = -1: no limit.
// initialDepth: 0 if no idea.
func ForDirMax(
	root string,
	initialDepth int,
	maxDepth int,
	cb func(depth int, dirname string, fi os.FileInfo) (stop bool, err error),
	excludes ...string,
) (err error) {
	if maxDepth > 0 && initialDepth >= maxDepth {
		return
	}

	var dirs []os.FileInfo
	rootDir := os.ExpandEnv(root)
	dirs, err = ioutil.ReadDir(rootDir)
	if err != nil {
		// Logger.Fatalf("error in ForDirMax(): %v", err)
		return
	}

	var stop bool

	//files, err :=os.ReadDir(rootDir)
	var fi os.FileInfo
	fi, err = os.Stat(rootDir)
	if err != nil {
		return
	}
	if stop, err = cb(initialDepth, rootDir, fi); stop {
		return
	}

	for _, f := range dirs {
		//Logger.Printf("  - %v", f.Name())
		if err != nil {
			log.NewStdLogger().Errorf("error in ForDirMax().cb: %v", err)
			continue
		}

		if f.IsDir() && (maxDepth <= 0 || (maxDepth > 0 && initialDepth+1 < maxDepth)) {
			e := false
			d := path.Join(root, f.Name())
			for _, ptn := range excludes {
				if dir.IsWildMatch(d, ptn) {
					e = true
					break
				}
			}
			if e {
				continue
			}

			if stop, err = cb(initialDepth, d, f); stop {
				return
			}
			if err = ForDirMax(d, initialDepth+1, maxDepth, cb); err != nil {
				log.NewStdLogger().Errorf("error in ForDirMax(): %v", err)
			}
		}
	}

	return
}
