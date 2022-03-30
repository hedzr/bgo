package logic

//nolint:goimports
import (
	"bufio"
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/log"
	"github.com/hedzr/log/dir"
	"github.com/hedzr/log/exec"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

type pkgInfo struct {
	dirname          string
	packageName      string
	groupKey         string
	groupLeadingText string
	projectName      string
	appName          string
	p                *ProjectWrap
}

//nolint:nakedret
func findMainPackages(bs *BgoSettings) (packages map[string]*pkgInfo, err error) {
	packages = make(map[string]*pkgInfo)

	keys := make([]string, 0, len(bs.Projects))
	{
		for k := range bs.Projects {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	for _, groupKey := range keys {
		group := bs.Projects[groupKey]
		kps := make([]string, 0, len(group.Items))
		{
			for k := range group.Items {
				kps = append(kps, k)
			}
			sort.Strings(kps)
		}

		logx.Verbose("curr-dir: %v", dir.GetCurrentDir())
		for _, projectKey := range kps {
			p := group.Items[projectKey]
			if p.Disabled || !dir.FileExists(p.Dir) {
				logx.Verbose("lookup group %v, project %v, disabled: %v, dir %q NOT EXISTS, ignored.", groupKey, projectKey, p.Disabled, p.Dir)
				continue
			}
			logx.Verbose("lookup group %v, project %v, disabled: %v, dir %q", groupKey, projectKey, p.Disabled, p.Dir)

			if err = dissectTheMainDir(p.Dir, packages, bs); err == nil {
				k := path.Clean(p.Dir)
				if _, ok := packages[k]; !ok {
					packages[k] = &pkgInfo{
						dirname:          k,
						packageName:      "",
						groupKey:         groupKey,
						groupLeadingText: group.LeadingText,
						projectName:      projectKey,
						appName:          p.Name,
						p:                p,
					}
				}
				packages[k].p = p
				packages[k].groupKey = groupKey
				packages[k].groupLeadingText = group.LeadingText
				packages[k].projectName = projectKey
				packages[k].appName = p.Name
				// err = fmp(bs, p.Dir, packages)
				logx.Trace("groupKey: %v, projKey: %v", groupKey, projectKey)
				if bs.Scope == "short" {
					return
				}
			}
		}
	}
	return
}

//nolint:nakedret
func scanWorkDir(workdir, scope string, packages map[string]*pkgInfo, bs *BgoSettings) (err error) {
	if len(packages) > 0 && scope == "short" {
		return
	}

	defer dir.PushDir(workdir)()

	logx.Log("> scanning %v...\n", dir.AbsPath("."))
	// find all packages no matter whether it's a main package

	pd := make(map[string]string) // dir -> package
	err = dir.ForDir(".", func(depth int, dirname string, fi os.FileInfo) (stop bool, err error) {
		dirName := dirname // path.Join(dirname, fi.Name())
		logx.Verbose("  >> %v/%v [%v], %v", dirname, fi.Name(), dirName, dir.GetCurrentDir())

		var mainFound bool

		err = dir.ForFileMax(dirName, 0, 1, func(depth int, cwd string, fi os.FileInfo) (stop bool, err error) {
			logx.Trace("     . %v", fi.Name())
			if dir.IsWildMatch(fi.Name(), "*.go") {
				fileName, pn := path.Join(cwd, fi.Name()), ""
				if pn, err = extractPackageName(fileName); err == nil {
					if pn == "main" { //nolint:goconst
						logx.Verbose("       > ADD: %v -> %v", pn, dirName)
						k := path.Clean(dirName)
						pd[k] = pn
						stop = true      // stop populating after first .go found
						mainFound = true // told up-level at least one 'main' pkg found
					}
				}
			}
			return
		}, ".*")

		if mainFound && scope == "short" {
			stop = true
		}
		return
	}, ".git", ".svn", ".hg", ".github", ".vscode", ".*")

	if err != nil {
		logx.Error("Error: %v", err)
		return
	}

	logx.Verbose("pd: %q", pd)
	for folder, pn := range pd {
		if pn == "main" {
			if err = dissectTheMainDir(folder, packages, bs); err == nil {
				logx.Colored(logx.Yellow, "  scanWorkDir: pkg main - %v", folder)
			}
		}
	}
	return
}

func dissectTheMainDir(dirname string, packages map[string]*pkgInfo, bs *BgoSettings) (err error) {
	if dirname == "" {
		return
	}

	if dirname != "." {
		defer dir.PushDir(dirname)()
	}

	_ = dissectTheMainDirImpl(dirname, packages, bs)
	return
}

//nolint:nakedret
func dissectTheMainDirImpl(dirname string, packages map[string]*pkgInfo, bs *BgoSettings) (err error) {
	cmd := []string{"go", "list", "-f", "{{.Name}},{{.ImportPath}}", "."}
	err = exec.CallSliceQuiet(cmd, func(retCode int, stdoutText string) {
		scanner := bufio.NewScanner(bufio.NewReader(strings.NewReader(stdoutText)))
		for scanner.Scan() {
			a := strings.Split(scanner.Text(), ",")
			if a[0] == "main" {
				k := path.Clean(dirname)
				if _, ok := packages[k]; ok {
					packages[k].packageName = a[1]
				} else {
					packages[k] = &pkgInfo{
						dirname:     k,
						packageName: a[1],
						p:           nil,
					}
				}
				if packages[k].p == nil {
					// associate with project entry in bgo-settings
					var found bool
					for _, g := range bs.Projects {
						for _, p := range g.Items {
							if path.Clean(p.Dir) == k {
								packages[k].p = p
								found = true
								break
							}
						}
					}
					if !found {
						// add this entry into bgo-settings too
						p := newProjectClean(packages[k].packageName, k)
						pk := path.Base(packages[k].packageName)
						for gk := range bs.Projects {
							bs.Projects[gk].Items[pk] = p
							packages[k].p = p
							break
						}
					}
				}
			}
		}
		if err = scanner.Err(); err != nil {
			logx.Error("Error: %v", err)
		}
	})
	if err != nil {
		// This error can be ignored safety, we will build the
		// target without go modules later:
		//
		// go: cannot find main module, but found .git/config in ...
		//
		if logx.IsVerboseMode() {
			log.Warnf("dissectTheMainDirImpl(%q) failed: %v", dirname, err)
		}
	}
	return
}

var rePackageName = regexp.MustCompile(`^package[ \t]+(\w+)`)

func extractPackageName(fileName string) (pn string, err error) {
	var f *os.File
	if f, err = os.Open(fileName); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(bufio.NewReader(f))
		for scanner.Scan() {
			r := rePackageName.FindAllStringSubmatch(scanner.Text(), -1)
			for _, r1 := range r {
				pn = r1[1]
				// pd[pn] = fullName
			}
		}
		if err = scanner.Err(); err != nil {
			logx.Error("Error: %v", err)
		}
	}
	return
}
