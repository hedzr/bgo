package logic

import (
	"fmt"
	"github.com/hedzr/bgo/internal/logic/build"
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/bgo/internal/logic/tool"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log/dir"
	"github.com/hedzr/log/exec"
	"gopkg.in/hedzr/errors.v2"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func buildCurr(buildScope string, cmd *cmdr.Command, args []string) (err error) {
	tp := build.NewTargetPlatforms()
	tp.SetOsArch(runtime.GOOS, runtime.GOARCH)
	err = buildFor(buildScope, tp, nil, nil, cmd, args)
	return
}

func buildAuto(buildScope string, cmd *cmdr.Command, args []string) (err error) {
	var tp = build.NewTargetPlatforms()

	err = tp.Init()
	if err != nil {
		return
	}

	err = buildFor(buildScope, tp, nil, nil, cmd, args)
	return
}

func buildFor(buildScope string, tp *build.TargetPlatforms,
	bc *build.Context, bs *BgoSettings,
	cmd *cmdr.Command, args []string) (err error) {

	if bc == nil {
		// initial build-context, preparing the build information
		bc = build.NewContext()
		//logx.Log("context: %v", bc)
		//logx.Log("context.build-info: %v", bc.BuildInfo)
	}

	if bs == nil {
		bs = newBgoSettings(buildScope)
		if isSaveMode() {
			bs.Scope = "full"
		}
	}

	if bs != nil {
		logx.Trace("tp = %v\n", tp)
		logx.Trace("bs = %+v\n", bs)
		// reduce the top-level targets to a minimal set
		tp.FilterBy(bs.Scope, bs.For, bs.Os, bs.Arch)
	} else {
		logx.Fatal("BAD, bgoSettings == nil!!")
	}

	err = buildProjects(tp, bc, bs, cmd, args)
	return
}

func buildProjects(tp *build.TargetPlatforms, bc *build.Context, bs *BgoSettings,
	cmd *cmdr.Command, args []string) (err error) {
	var packages map[string]*pkgInfo
	if packages, err = findMainPackages(bs); err != nil {
		return
	}

	projectName, singleProjectOrPackage := checkSingleProjectNameSpecified(packages, bs)

	if bs.Scope != "auto" || cmdr.GetUsedAlterConfigFile() == "" {
		if (singleProjectOrPackage == nil && projectName != "") || projectName == "" {
			err = scanWorkDir(bc.WorkDir, bs.Scope, packages, bs)

			for _, pi := range packages {
				ensureProject(pi, bc, bs)
			}
		}
	}

	if !isDryRunMode() && isSaveMode() {
		if err = saveBackToBs(packages, bs); err == nil {
			err = saveNewBgoYamlFile(bs)
		}
		return
	}

	return buildProjectsImpl(packages, tp, bc, bs, cmd, args)
}

func buildProjectsImpl(
	packages map[string]*pkgInfo,
	tp *build.TargetPlatforms, bc *build.Context, bs *BgoSettings,
	cmd *cmdr.Command, args []string,
) (err error) {

	if cmd.GetTitleName() == "list" {
		err = listPackages(tp, bc, bs, packages)
		return
	}

	// start building now | loop for all modules/projects now
	ki := getSortedV(packages)
	for _, k := range ki {
		if err = buildPackages(tp, bc, bs, packages[k.Index]); err != nil {
			break
		}
	}
	return
}

func checkSingleProjectNameSpecified(packages map[string]*pkgInfo, bs *BgoSettings) (projectName string, singleProjectOrPackage *pkgInfo) {
	projectName = cmdr.GetStringR("build.project-name")
	if projectName != "" {
		for _, p := range packages {
			if tool.StripOrderPrefix(p.projectName) == projectName ||
				tool.StripOrderPrefix(p.appName) == projectName {
				singleProjectOrPackage = p
				break
			}
		}
		if singleProjectOrPackage != nil {
			packages = make(map[string]*pkgInfo)
			packages[path.Clean(singleProjectOrPackage.dirname)] = singleProjectOrPackage
			for _, g := range bs.Projects {
				for pn, p := range g.Items {
					if tool.StripOrderPrefix(pn) == projectName ||
						tool.StripOrderPrefix(p.Name) == projectName {
						continue
					}
					p.Disabled = true
				}
			}
		}
	}
	return
}

func ensureProject(pi *pkgInfo, bc *build.Context, bs *BgoSettings) {
	p := pi.p
	if p == nil {
		p = newProject(pi.packageName, bc, bs)
	} else {
		p.apply(pi.packageName, bs)
	}
	p.applyPI(pi)
	pi.p = p
}

func saveBackToBs(packages map[string]*pkgInfo, bs *BgoSettings) (err error) {
	for _, pi := range packages {
		var found bool
		for _, g := range bs.Projects {
			for _, p := range g.Items {
				if p.Dir == pi.dirname {
					found = true
					cmdr.Clone(pi.p, p)
					break
				}
			}
		}

		if !found {
			pn := pi.appName
			if pn == "" {
				pn = pi.projectName
				if pn == "" {
					pn = path.Base(pi.dirname)
				}
			}
			if len(bs.Projects) == 0 {
				bs.Projects = make(map[string]ProjectGroup)
				bs.Projects["000-default-group"] = ProjectGroup{Items: make(map[string]*ProjectWrap)}
			}
			//if bs.Projects["000-default-group"].Items == nil {
			//	bs.Projects["000-default-group"].Items = make(map[string]*ProjectWrap)
			//}
			bs.Projects["000-default-group"].Items[pn] = pi.p
		}
	}
	return
}

func buildPackages(tpBase *build.TargetPlatforms, bc *build.Context, bs *BgoSettings, pi *pkgInfo) (err error) {
	ensureProject(pi, bc, bs)

	//logx.Colored(logx.Green, "> building for package %v (dir: %q)...", pi.p.Package, pi.dirname)

	if cmdr.GetTraceMode() {
		//logx.Dim("  - BS:")
		//logx.Dim("%v", leftPad(yamlText(bs.Projects), 4))

		logx.Dim("  - Project:")
		logx.Dim("%v", leftPad(yamlText(pi.p), 4))
		//if isDryRunMode() {
		//	os.Exit(0)
		//}
	}

	err = loopAllProjects(tpBase, bc, bs, buildProject)

	//STOP:
	//	for oss, osv := range tp.OsArchMap {
	//		for arch, _ := range osv {
	//			if cmdr.GetTraceMode() {
	//				logx.Dim("%v\n", leftPad(yamlText(pi.p.Common), 5))
	//			}
	//
	//			prepareBuildContextForEachProjectTarget(bc, oss, arch, pi.p,
	//				pi.projectName, pi.groupKey, pi.groupLeadingText)
	//
	//			err = buildProject(bc, bs)
	//			if err != nil {
	//				break STOP
	//			}
	//		}
	//	}
	return
}

func prepareBuildContextForEachProjectTarget(bc *build.Context, os, arch string,
	p *ProjectWrap, knownProjectName, groupKey, groupLeadingText string) {
	// dynBuildInfo
	bc.BgoGroupKey = groupKey
	bc.BgoGroupLeadingText = groupLeadingText
	bc.HasGoMod = false
	bc.GOROOT, bc.Dir = "", p.Dir

	// build info
	bc.Serial++
	bc.RandomInt = tool.NextIn(100)
	bc.RandomString = tool.NextString(24)

	// Bug: Clone/CloneFrom can't process unhashable/uncomparable
	// field such as the element of 'Extends', i.e.,
	// PackageNameValues can't be cloned properly since its
	// 'Values' field is a map (map is not a comparable type in go).
	//
	// So, we have to clear this field before we'll CloneFrom somewhere.
	bc.Common = build.NewCommon()

	// bc.DynBuildInfo
	bc.Common.CloneFrom(p.Common)

	// bc.Gen, bc.Cgo, bc.Race, bc.Gocmd, bc.Install = p.Gen, p.Cgo, p.Race, p.Gocmd, p.Install
	bc.OS, bc.ARCH, bc.Version = os, arch, bc.GitVersion
	if p.Version != "" {
		bc.Version = p.Version
	}

	bc.FindAppName(p.Name, knownProjectName, p.Package)
}

func buildProject(bc *build.Context, bs *BgoSettings) (err error) {
	//logx.Log("  >> %v/%v", bc.Os, bc.Arch)
	logx.Log("      >> %v/%v, Need Install: %v\n", bc.OS, bc.ARCH, bc.Install)
	//logx.Dim("     project.Common: %+v\n", *p.Common)

	var cmd []interface{}
	cmd, err = prepareCommandLine(bc, bs)

	var outBinary string
	outBinary, err = getBuildTargetBinaryPath(bc, bs)
	if err != nil {
		return
	}

	if bc.KeepWorkdir {
		//wd := bc.WorkDir
		if !dir.FileExists(bc.Dir) {
			return
		}

		relOut := outBinary
		bc.Output.Set(relOut)
		cmd = append(cmd, "-o", relOut)

		if !path.IsAbs(bc.Dir) && !strings.HasPrefix(bc.Dir, "./") {
			if bc.Dir != "." {
				bc.Dir = "./" + bc.Dir
			}
		}
		cmd = append(cmd, bc.Dir)
		bc.PackageDir = bc.Dir
	} else {
		wd := bc.Dir
		if bc.UseWorkDir != "" {
			wd = bc.UseWorkDir
		}

		dir.PushDir(wd)()
		//if c, e := dir.PushDirEx(wd); e != nil {
		//	logx.Warn("%v\n", logx.ToColor(logx.LightRed, "         The project ignored since not exists."))
		//	return nil
		//} else {
		//	defer c()
		//}
		logx.Verbose("         entering dir: %v\n", wd)

		var relOut string
		relOut, err = filepath.Rel(wd, outBinary)
		if err != nil {
			return
		}
		bc.Output.Set(relOut)
		cmd = append(cmd, "-o", relOut)

		var relpkg string
		relpkg, err = filepath.Rel(wd, bc.Dir)
		if !path.IsAbs(relpkg) && !strings.HasPrefix(relpkg, "./") {
			if relpkg != "." {
				relpkg = "./" + relpkg
			}
		}
		cmd = append(cmd, relpkg)
		bc.PackageDir = relpkg
	}

	logx.Verbose("         PWD: %v\n", logx.ToDim(dir.GetCurrentDir()))
	logx.Verbose("         OUT: %v\n", logx.ToDim(bc.Output.Path))
	logx.Verbose("         Go.mod: %v\n", logx.ToColor(logx.Cyan, bc.GoModFile))
	logx.Verbose("         CommandLine: %v\n", logx.ToDim("%v", cmd))

	if isDryRunMode() {
		logx.Colored(logx.Yellow, "         STOP since dry-run mode specified,\n")
		return
	}

	err = goBuild(bc, bs, cmd...)
	err = nil
	return
}

func prepareCommandLine(bc *build.Context, bs *BgoSettings) (cmd []interface{}, err error) {
	cmd = []interface{}{"go", "build"}

	st := path.Join(bc.Dir, "go.mod")
	if dir.FileExists(st) {
		bc.HasGoMod, bc.GoModFile = true, st
	} else {
		st = path.Join(bc.WorkDir, "go.mod")
		if dir.FileExists(st) {
			bc.HasGoMod, bc.GoModFile = true, st
		}
	}

	if bc.Gocmd != "" {
		gocmd := os.ExpandEnv(bc.Gocmd)
		if x, e := exec.LookPath(gocmd); e == nil {
			y := dir.FollowSymLink(x)
			cmd[0] = y
			yup1 := path.Dir(y)
			if bin := path.Base(yup1); bin == "bin" {
				bc.GOROOT = path.Dir(yup1)
			}
		}
	}

	if cmdr.GetBoolR("build.verbose") {
		cmd = append(cmd, "-v")
	}

	if bc.Debug == false {
		if cmdr.GetBoolR("build.no-trimpath") == false {
			cmd = append(cmd, "-trimpath")
		}

		bc.Ldflags = uniappend(bc.Ldflags, "-s")
		bc.Ldflags = uniappend(bc.Ldflags, "-w")
	}

	if cmdr.GetBoolR("build.rebuild") {
		cmd = append(cmd, "-a")
	}

	bc.Gen = cmdr.GetBoolR("build.generate", bc.Gen)

	x := pclMore1(bc, bs)
	cmd = append(cmd, x...)
	return
}

func pclMore1(bc *build.Context, bs *BgoSettings) (cmd []interface{}) {
	if str := cmdr.GetStringR("build.mod"); str != "" {
		cmd = append(cmd, "-mod", str)
	}

	if bc.Race {
		cmd = append(cmd, "-race")
	}

	if bc.Msan {
		cmd = append(cmd, "-msan")
	}

	if len(bc.Asmflags) > 0 {
		cmd = append(cmd, strings.Join([]string{"-asmflags",
			strings.Join(bc.Asmflags, " ")}, "="))
	}
	if len(bc.Gcflags) > 0 {
		cmd = append(cmd, strings.Join([]string{"-gcflags",
			strings.Join(bc.Gcflags, " ")}, "="))
	}
	if len(bc.Gccgoflags) > 0 {
		cmd = append(cmd, strings.Join([]string{"-gccgoflags",
			strings.Join(bc.Gccgoflags, " ")}, "="))
	}
	if len(bc.Tags) > 0 {
		cmd = append(cmd, strings.Join([]string{"-tags",
			strings.Join(bc.Tags, " ")}, "="))
	}

	ifLdflags(bc)

	if len(bc.Ldflags) > 0 {
		cmd = append(cmd, strings.Join([]string{"-ldflags",
			strings.Join(bc.Ldflags, " ")}, "="))
	}

	return
}

func goBuild(bc *build.Context, bs *BgoSettings, cmd ...interface{}) (err error) {
	if err = goBuildPreChecks(bc, bs); err != nil {
		return
	}

	var opts []exec.Opt
	var cgo bool
	if opts, cgo, err = goBuildPrepareOpts(bc, bs); err != nil {
		return
	}

	ec := errors.NewContainer("holder")
	c := exec.New(opts...).
		WithCommand(cmd...).
		WithEnv("GOOS", bc.OS).
		WithEnv("GOARCH", bc.ARCH).
		WithEnv("CGO_ENABLED", boolToString(cgo)).
		//WithStdoutCaught(). // can be removed
		//WithStderrCaught(). // can be removed
		WithOnOK(okHandler(ec, bc, bs)).
		WithOnError(func(err error, retCode int, stdoutText, stderrText string) {
			logx.Error("ERROR TEXT:\n%v\nError:\n%v\nRetCode: %v\nCommands: %v\n",
				logx.ToColor(logx.Red, leftPad(stderrText, 4)),
				logx.ToColor(logx.Red, leftPad(fmt.Sprintf("%v", err), 4)),
				logx.ToDim(strconv.Itoa(retCode)),
				logx.ToDim(fmt.Sprintf("%v", cmd)))
		})
	if err = c.RunAndCheckError(); err != nil {
		ec.Attach(err)
	}

	if err = ec.Error(); err != nil {
		logx.Error("Error occurs: %v", err)
	}
	return
}

func goBuildPreChecks(bc *build.Context, bs *BgoSettings) (err error) {
	if bc.PreAction != "" {
		if err = iaRunScript(bc.PreAction, bc, "pre-action"); err != nil {
			return
		}
	}
	if bc.PreActionFile != "" && dir.FileExists(bc.PreActionFile) {
		if err = iaRunScriptFile(bc.PreActionFile, bc, "pre-action-file"); err != nil {
			return
		}
	}

	if bc.Gen {
		if err = iaGenerate(bc, bs); err != nil {
			return
		}
	}
	return
}

func goBuildPrepareOpts(bc *build.Context, bs *BgoSettings) (opts []exec.Opt, cgo bool, err error) {
	if bc.HasGoMod {
		opts = append(opts, exec.WithEnv("GO111MODULE", "on"))
		logx.DimV("           GO111MODULE: ON\n")
	} else {
		opts = append(opts, exec.WithEnv("GO111MODULE", "off"))
		logx.DimV("           GO111MODULE: off\n")
	}
	if bs.Goproxy != "" {
		opts = append(opts, exec.WithEnv("GOPROXY", bs.Goproxy))
	}
	if bc.GOROOT != "" {
		opts = append(opts, exec.WithEnv("GOROOT", bc.GOROOT))
		//WithEnv("GOPATH", os.ExpandEnv("$HOME/go")).
	}

	cgo = bc.Cgo
	if cgo && runtime.GOOS != bc.OS || runtime.GOARCH != bc.ARCH {
		cgo = false
	}

	return
}

func okHandler(ec *errors.WithCauses, bc *build.Context, bs *BgoSettings) (onOK func(retCode int, stdoutText string)) {
	return func(retCode int, stdoutText string) {
		var err error
		if bc.Install {
			if err = iaInstall(bc.Output.Path, bc, bs); err != nil {
				ec.Attach(err)
				return
			}
		}
		if bc.PostAction != "" {
			if err = iaRunScript(bc.PostAction, bc, "post-action"); err != nil {
				ec.Attach(err)
				return
			}
		}
		if bc.PostActionFile != "" && dir.FileExists(bc.PostActionFile) {
			if err = iaRunScriptFile(bc.PostActionFile, bc, "post-action-file"); err != nil {
				ec.Attach(err)
				return
			}
		}
		if !bc.DisableResult {
			if err = iaLL(bc.Output.Path, bc); err != nil {
				ec.Attach(err)
				return
			}
		}
		if len(stdoutText) > 0 {
			logx.Dim("OUTPUT:\n%v\n", stdoutText)
		}

		// exec.New().WithCommandString("bash -c 'echo hello world!'", '\'').WithContext(context.Background()).Run()

		return
	}
}

//goland:noinspection ALL
func getBuildTargetBinaryPath(bc *build.Context, bs *BgoSettings) (outBinary string, err error) {
	//var outBinary string
	if outBin := cmdr.GetStringR("build.output", bs.Output.NamedAs); outBin != "" {
		if bs.Scope == "short" {
			outBin = "{{.AppName}}"
		}
		if runtime.GOOS == "windows" {
			outBin += ".exe"
		}

		tpl := path.Join(bs.Output.Dir, bs.Output.SplitTo, outBin)
		if outBinary, err = tplExpand(tpl, "output-binary-name", bc); err != nil {
			return
		}
	}
	return
}

func iaGenerate(bc *build.Context, bs *BgoSettings) (err error) {
	logx.Log("         > Run 'go generate' at %q...\n", bc.PackageDir)
	return exec.New().
		WithCommand("go", "generate", bc.PackageDir).
		RunAndCheckError()
}

func iaInstall(outBinary string, bc *build.Context, bs *BgoSettings) (err error) {
	if bc.OS == bc.GOOS && bc.ARCH == bc.GOARCH {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = os.ExpandEnv("$HOME/go")
		}
		goBin := path.Join(gopath, "bin")

		logx.Log("         > Installing to %v...\n", goBin)
		err = exec.New().WithCommand("cp", outBinary, goBin).RunAndCheckError()
	}
	return
}

func iaRunScript(scriptsSource string, bc *build.Context, title ...string) (err error) {
	var ttl = "invoking-shell-scripts"
	for _, s := range title {
		ttl = s
		break
	}

	var script string
	if script, err = tplExpand(scriptsSource, ttl, bc); err == nil {
		if logx.IsVerboseMode() {
			logx.Log("         > Invoking %v:\n", ttl)
			logx.Dim("%v\n", leftPad(script, 7))
		} else {
			logx.Log("         > Invoking %v...\n", ttl)
		}
		err = exec.New().
			WithPadding(7+4).
			WithCommand("/bin/bash", "-c", script).
			RunAndCheckError()
	}
	return
}

func iaRunScriptFile(scriptsSource string, bc *build.Context, title ...string) (err error) {
	var ttl = "invoking-shell-scripts"
	for _, s := range title {
		ttl = s
		break
	}

	var script string
	if script, err = tplExpand(scriptsSource, ttl, bc); err == nil {
		if logx.IsVerboseMode() {
			logx.Log("         > Invoking %v:\n", ttl)
			logx.Dim("%v\n", leftPad(script, 7))
		} else {
			logx.Log("         > Invoking %v...\n", ttl)
		}
		err = exec.New().
			WithPadding(7+4).
			WithCommand("/bin/bash", "-c", script).
			RunAndCheckError()
	}
	return
}

func iaLL(outBinary string, bc *build.Context) (err error) {
	// ll binary
	c := "--color"
	if runtime.GOOS == "darwin" {
		c = "-G"
	}
	targets := []string{outBinary}
	if bc.Install {
		if bc.OS == bc.GOOS && bc.ARCH == bc.GOARCH {
			gopath := os.Getenv("GOPATH")
			if gopath == "" {
				gopath = os.ExpandEnv("$HOME/go")
			}
			gobin := path.Join(gopath, "bin")

			t := path.Join(gobin, path.Base(outBinary))
			targets = append(targets, t)
		}
	}
	err = exec.New().WithPadding(7+2).WithCommand("ls", "-la", c, targets).RunAndCheckError()
	//err = exec.New().WithPadding(7).WithCommand("gls", "-lh", "--color", targets).RunAndCheckError()
	//err = exec.New().WithCommand("ls", "-la", c, targets).RunAndCheckError()
	return
}
