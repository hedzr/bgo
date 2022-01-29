package build

import (
	"bufio"
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log/exec"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
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
		GoVersion   string // the result from 'go version'
		GitVersion  string // the result from 'git describe --tags --abbrev=0'
		GitRevision string // revision, git hash code, from 'git rev-parse --short HEAD'
		BuildTime   string //
		GOOS        string // a copy from runtime.GOOS
		GOARCH      string // a copy from runtime.GOARCH
		GOVERSION   string // a copy from runtime.Version()

		RandomString string
		RandomInt    int
		Serial       int
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

func Prepare() {
	prepareBuildInfo()
}

var (
	gBuildInfo    *Info
	onceBuildInfo sync.Once
)

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
		})
		if err != nil {
			logx.Warn("No suitable Golang executable 'go' found, use runtime.Version() instead.")
			//os.Exit(1)
			gBuildInfo.GoVersion = runtime.Version()
		}

		err = exec.CallQuiet("git describe --tags --abbrev=0", func(retCode int, stdoutText string) {
			gBuildInfo.GitVersion = strings.TrimSuffix(stdoutText, "\n")
			logx.Colored(logx.Green, "git.version: %v", gBuildInfo.GitVersion)
		})
		if err != nil {
			logx.Warn("No suitable 'git' executable found or cannot git describe.")
			logx.Log("%v", err)
			//os.Exit(1)
			gBuildInfo.GitVersion = "-unknown-"
		}

		err = exec.CallQuiet("git rev-parse --short HEAD", func(retCode int, stdoutText string) {
			gBuildInfo.GitRevision = strings.TrimSuffix(stdoutText, "\n")
			logx.Colored(logx.Green, "git.revision: %v", gBuildInfo.GitRevision)
		})
		if err != nil {
			logx.Warn("No suitable 'git' executable found or not a git repo.")
			logx.Log("%v", err)
			//os.Exit(1)
			gBuildInfo.GitRevision = "-unknown-"
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

		if err := scanner.Err(); err != nil {
			logx.Fatal("Error: %v", err)
		}

		//logx.Colored(logx.Green, "git.revision: %v", gBuildInfo.GitRevision)
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

		if err := scanner.Err(); err != nil {
			logx.Fatal("Error: %v", err)
		}

		//cmdr.Set("bgo.dists", TargetPlatforms)
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
		//logColored(green, "os.arch.map: %v", TargetPlatforms)
	})
}
