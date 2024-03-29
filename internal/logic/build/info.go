package build

//nolint:goimports //so what
import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hedzr/cmdr"
	"github.com/hedzr/log/exec"

	"github.com/hedzr/bgo/internal/logx"
)

type (
	DynBuildInfo struct {
		ProjectName         string
		AppName             string
		Version             string // copied from Project.Version
		BgoGroupKey         string // project-group key in .bgo.yml
		BgoGroupLeadingText string // same above,
		HasGoMod            bool   //
		GoModFile           string //
		GOROOT              string // force using a special GOROOT
		Dir                 string
	}

	Info struct {
		GoVersion      string // the result from 'go version'
		GoVersionMajor int    // 1
		GoVersionMinor int    // 17
		GitVersion     string // the result from `git describe --tags --abbrev=0`
		GitRevision    string // revision, git hash code, from `git rev-parse --short HEAD`
		GitSummary     string // `git describe --tags --dirty --always`
		GitDesc        string // `git log --online -1`
		BuilderComment string //
		BuildTime      string // format is like '2023-01-22T09:26:07+08:00'
		GOOS           string // a copy from runtime.GOOS
		GOARCH         string // a copy from runtime.GOARCH
		GOVERSION      string // a copy from runtime.Version()

		RandomString string
		RandomInt    int
		Serial       int64
	}
)

func NewBuildInfo() *Info {
	prepareBuildInfo()
	return gBuildInfo
}

func newDynBuildInfo() *DynBuildInfo {
	return &DynBuildInfo{
		AppName:             "",
		Version:             "",
		ProjectName:         "",
		BgoGroupKey:         "",
		BgoGroupLeadingText: "",
	}
}

//nolint:gochecknoglobals //no
var (
	gBuildInfo    *Info
	onceBuildInfo sync.Once
)

func (c *Info) VersionIsGreaterThan(major, minor int) bool {
	if c.GoVersionMajor > major {
		return true
	} else if major == c.GoVersionMajor && c.GoVersionMinor > minor {
		return true
	}
	return false
}

func prepareBuildInfo() {
	onceBuildInfo.Do(func() {
		if gBuildInfo != nil {
			return
		}

		gBuildInfo = new(Info)
		gBuildInfo.GOOS = runtime.GOOS
		gBuildInfo.GOARCH = runtime.GOARCH
		gBuildInfo.GOVERSION = runtime.Version()

		var err error

		err = exec.CallQuiet("go version", func(retCode int, stdoutText string) {
			gBuildInfo.GoVersion = strings.ReplaceAll(
				strings.TrimSuffix(strings.TrimPrefix(stdoutText, "go version "), "\n"),
				" ", "_")
			logx.Colored(logx.Green, "go.version: %v", gBuildInfo.GoVersion)
			if strings.HasPrefix(gBuildInfo.GoVersion, "go") {
				a := strings.Split(gBuildInfo.GoVersion[2:], "_")
				b := strings.Split(a[0], ".")
				gBuildInfo.GoVersionMajor, _ = strconv.Atoi(b[0])
				gBuildInfo.GoVersionMinor, _ = strconv.Atoi(b[1])
			}
		})
		if err != nil {
			logx.Warn("No suitable Golang executable 'go' found, use runtime.Version() instead.")
			// os.Exit(1)
			gBuildInfo.GoVersion = runtime.Version()
		} else if logx.CountOfVerbose() > 1 {
			logx.Warn("is warn ok")
			logx.Trace("is trace ok")
			logx.Log("is ok")
			logx.Colored(logx.Green, "is green ok")
		}

		err = exec.CallQuiet("git describe --tags --abbrev=0", func(retCode int, stdoutText string) {
			gBuildInfo.GitVersion = strings.TrimSuffix(stdoutText, "\n") // such as v0.3.13
			logx.Colored(logx.Green, "git.version: %v", gBuildInfo.GitVersion)
		})
		if err != nil {
			logx.Warn("No suitable 'git' executable found or cannot git describe.")
			logx.Log("Error: %v", err)
			logx.Log("Env:   %v", os.Environ())
			// os.Exit(1)
			gBuildInfo.GitVersion = "-unknown-"
		}

		err = exec.CallQuiet("git rev-parse --short HEAD", func(retCode int, stdoutText string) {
			gBuildInfo.GitRevision = strings.TrimSuffix(stdoutText, "\n") // such as 3e6fd96
			logx.Colored(logx.Green, "git.revision: %v", gBuildInfo.GitRevision)
		})
		if err != nil {
			logx.Warn("Cannot get revision from git commits. No suitable 'git' executable found or not a git repo?")
			logx.Log("%v", err)
			// os.Exit(1)
			gBuildInfo.GitRevision = "-unknown-"
		}

		err = exec.CallQuiet("git describe --tags --dirty --always", func(retCode int, stdoutText string) {
			gBuildInfo.GitSummary = strings.TrimSuffix(stdoutText, "\n")
			logx.Colored(logx.Green, "git.summary: %v", gBuildInfo.GitSummary)
		})
		if err != nil {
			logx.Warn("Cannot git describe. No suitable 'git' executable found or not a git repo?")
			logx.Log("%v", err)
			// os.Exit(1)
			gBuildInfo.GitSummary = ""
		}

		err = exec.CallQuiet("git log --oneline -1", func(retCode int, stdoutText string) {
			gBuildInfo.GitDesc = strings.TrimSuffix(stdoutText, "\n")
			logx.Colored(logx.Green, "git.description: %v", gBuildInfo.GitDesc)
		})
		if err != nil {
			logx.Warn("Cannot retrieve last commit detail. No suitable 'git' executable found or not a git repo?")
			logx.Log("%v", err)
			// os.Exit(1)
			gBuildInfo.GitDesc = ""
		}

		gBuildInfo.BuildTime = time.Now().Format(time.RFC3339)
		logx.Colored(logx.Green, "build.time: %v", gBuildInfo.BuildTime)

		prepareWithGoEnv()
		retrieveGoToolDists()
	})
}

func prepareWithGoEnv() {
	var err error
	err = exec.Call("go env", func(retCode int, stdoutText string) {
		scanner := bufio.NewScanner(bufio.NewReader(strings.NewReader(stdoutText)))
		for scanner.Scan() {
			a := strings.Split(scanner.Text(), "=")
			switch a[0] {
			case "GOOS", "GOARCH":
				v, _ := strconv.Unquote(a[1])
				if err = os.Setenv(a[0], v); err != nil {
					logx.Fatal("Error: %v", err)
				}
			default:
				if strings.HasSuffix(a[0], "FLAGS") {
					v, _ := strconv.Unquote(a[1])
					if err = os.Setenv(a[0], v); err != nil {
						logx.Fatal("Error: %v", err)
					}
				}
			}
		}

		if err = scanner.Err(); err != nil {
			logx.Fatal("Error: %v", err)
		}

		// logx.Colored(logx.Green, "git.revision: %v", gBuildInfo.GitRevision)
	})
}

func retrieveGoToolDists() {
	var err error
	err = exec.Call("go tool dist list", func(retCode int, stdoutText string) {
		var osArchMap = make(map[string]map[string]bool)
		scanner := bufio.NewScanner(bufio.NewReader(strings.NewReader(stdoutText)))
		for scanner.Scan() {
			a := strings.Split(scanner.Text(), "/")
			if _, ok := osArchMap[a[0]]; !ok {
				osArchMap[a[0]] = make(map[string]bool)
			}
			osArchMap[a[0]][a[1]] = true
		}

		if err = scanner.Err(); err != nil {
			logx.Fatal("Error: %v", err)
		}

		// merge all distros into cmdr Option Store, key path: app.bgo.dists.os-arch-map

		// cmdr.Set("bgo.dists", TargetPlatforms)
		err = cmdr.MergeWith(map[string]interface{}{
			"app": map[string]interface{}{
				"bgo": map[string]interface{}{
					"dists": map[string]interface{}{
						"os-arch-map": osArchMap,
					},
				},
			},
		})
		if err != nil {
			logx.Fatal("Error: %v", err)
		}
		// logColored(green, "os.arch.map: %v", TargetPlatforms)
	})
}
