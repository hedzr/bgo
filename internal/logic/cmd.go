package logic

import (
	"github.com/hedzr/cmdr"
)

func AttachToCmdr(root *cmdr.RootCmdOpt) {
	// root.Action(rootAction)

	// make root command as an alias of subcommand 'build'
	root.RunAsSubCommand("build")

	cmdrOptsScopes(root)

	cmdr.NewBool().
		Titles("dry-run", "dr").
		Description("Just listing without Go build", "").
		// ToggleGroup("").
		Group("Control").
		EnvKeys("DRY_RUN").
		AttachTo(root)

	// cmdrOptsBuildCommons(root)
	// cmdrOptsBuildControls(root)

	cmdrSubCmdInit(root)
	cmdrSubCmdBuild(root)
	cmdrSubCmdList(root)

	cmdrSubCmdSBOM(root)

	cmdrSubCmdRun(root)  // forwarding to `go run`
	cmdrSubCmdTest(root) // forwarding to `go test`
}

const gnBuild = "Build Target"
const gnGoBuild = "Go Build Params"
const gnControl = "Control"

func cmdrOptsScopes(cmd cmdr.OptCmd) {
	cmdr.NewBool().
		Titles("short", "s").
		Description("Build for current CPU and OS Arch ONLY", "").
		ToggleGroup("Scope").
		// Group("Build").
		EnvKeys("SHORT").
		AttachTo(cmd)
	cmdr.NewBool().
		Titles("full", "f").
		Description("Build all CLIs under work directory recursively", "").
		ToggleGroup("Scope").
		// Group("Build").
		EnvKeys("FULL").
		AttachTo(cmd)
	cmdr.NewBool(true).
		Titles("auto", "a").
		Description("Build all modules defined in .bgo.yaml recursively (auto mode)", "").
		ToggleGroup("Scope").
		// Group("Build").
		EnvKeys("AUTO").
		AttachTo(cmd)
}

func cmdrOptsBuildCommons(cmd cmdr.OptCmd) { //nolint:funlen //so what
	cmdr.NewStringSlice().
		Titles("osarch", "arch", "arch").
		//nolint:lll //no
		Description("Specify the CPU-Arch list (comma separated or multi times)", `Specify the CPU-Arch list (comma separated or multi times)
		Such as '--arch arm64 --arch amd64,386'.
		The available list can be dumped by 'go tool dist list'`).
		ToggleGroup("").
		Group(gnBuild).
		Placeholder("ARCH").
		EnvKeys("GOARCH").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("os", "os").
		//nolint:lll //no
		Description("Specify the OS list (comma separated or multi times)", `Specify the OS list (comma separated or multi times)
		Such as '--os linux --os plan9 --os darwin,windows'.
		The available list can be dumped by 'go tool dist list'`).
		ToggleGroup("").
		Group(gnBuild).
		Placeholder("OS").
		EnvKeys("GOOS").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("for", "for").
		//nolint:lll //no
		Description("Specify the OS/arch list (comma separated or multi times)", `Specify the OS list (comma separated or multi times)
		Such as '--for linux/arm64 --os plan9/amd64'.
		The available list can be dumped by 'go tool dist list'`).
		ToggleGroup("").
		Group(gnBuild).
		Placeholder("OS/ARCH").
		EnvKeys("FOR").
		AttachTo(cmd)

	cmdr.NewStringSlice().
		Titles("tags", "tags").
		//nolint:lll //no
		Description("Additional '-tags' value to pass to go build (comma separated or multi times)", `Additional '-tags' value to pass to go build (comma separated or multi times)
		Such as '--tags isdelve --tags private --os aws-build,lan-build'.`).
		ToggleGroup("").
		Group(gnGoBuild).
		Placeholder("tag,list").
		EnvKeys("").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("ldflags", "ldflags").
		Description("Additional '-ldflags' value to pass to go build (comma separated or multi times)", ``).
		ToggleGroup("").
		Group(gnGoBuild).
		Placeholder("[pattern=]arg list").
		EnvKeys("LDFLAGS").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("gcflags", "gcflags").
		Description("Additional '-gcflags' value to pass to go build (comma separated or multi times)", ``).
		ToggleGroup("").
		Group(gnGoBuild).
		Placeholder("[pattern=]arg list").
		EnvKeys("GCFLAGS").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("gccgoflags", "gccgoflags").
		Description("Additional '-gccgoflags' value to pass to go build (comma separated or multi times)", ``).
		ToggleGroup("").
		Group(gnGoBuild).
		Placeholder("[pattern=]arg list").
		EnvKeys("GCCGOFLAGS").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("asmflags", "asmflags").
		Description("Additional '-asmflags' value to pass to go build (comma separated or multi times)", ``).
		ToggleGroup("").
		Group(gnGoBuild).
		Placeholder("[pattern=]arg list").
		EnvKeys("ASMFLAGS").
		AttachTo(cmd)

	cmdr.NewString("").
		Titles("goamd64", "goamd64", "amd64").
		Description("Additional 'GOAMD64d' envvar value", ``).
		ToggleGroup("").
		Group(gnGoBuild).
		Placeholder("GOAMD64").
		ValidArgs("v1", "v2", "v3", "v4").
		EnvKeys("GOAMD64").
		AttachTo(cmd)

	cmdr.NewString("").
		Titles("mod", "mod").
		Description("Additional '-mod' value to pass to go build", ``).
		ToggleGroup("").
		Group(gnGoBuild).
		Placeholder("MOD").
		ValidArgs("readonly", "vendor", "mod").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("cgo", "cgo").
		Description("enable CGO mode", "").
		// ToggleGroup("").
		Group(gnGoBuild).
		EnvKeys("CGO_ENABLED").
		AttachTo(cmd)
	cmdr.NewBool().
		Titles("race", "race").
		Description("enable --race building", "").
		// ToggleGroup("").
		Group(gnGoBuild).
		EnvKeys("RACE").
		AttachTo(cmd)
	cmdr.NewBool().
		Titles("msan", "msan").
		//nolint:lll //no
		Description("enable --msan building (enable interoperation with memory sanitizer)", `enable --msan building (enable interoperation with memory sanitizer).
		Supported only on linux/amd64, linux/arm64
		and only with Clang/LLVM as the host C compiler.
		On linux/arm64, pie build mode will be used.`).
		// ToggleGroup("").
		Group(gnGoBuild).
		EnvKeys("MSAN").
		AttachTo(cmd)

	cmdr.NewString("{{.AppName}}-{{.OS}}-{{.ARCH}}").
		Titles("output", "o").
		Description("Specify the binary filename pattern", ``).
		ToggleGroup("").
		Group(gnGoBuild).
		Placeholder("PATTERN").
		AttachTo(cmd)
}

func cmdrOptsBuildControls(cmd cmdr.OptCmd) {
	cmdr.NewString("").
		Titles("gocmd", "gocmd").
		Description("Additional '-gocmd' value to pass to go build", ``).
		ToggleGroup("").
		Group(gnControl).
		Placeholder("GOCMD").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("rebuild", "r").
		Description("Force rebuilding of package that were up to date", "").
		// ToggleGroup("").
		Group(gnControl).
		EnvKeys("REBUILD").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("no-trimpath", "ntp").
		Description("Don't use -trimapth", "").
		// ToggleGroup("").
		Group(gnControl).
		EnvKeys("NO_TRIMPATH").
		AttachTo(cmd)

	// cmdr.NewString("").
	//	Titles("for", "for").
	//	Description("Build for os/arch, such as 'linux/riscv64'", "").
	//	//ToggleGroup("").
	//	Group(gnControl).
	//	EnvKeys("FOR").
	//	AttachTo(cmd)
	//
	// cmdr.NewString("").
	//	Titles("os", "os").
	//	Description("Build for os, such as 'linux'", "").
	//	//ToggleGroup("").
	//	Group(gnControl).
	//	EnvKeys("OS").
	//	AttachTo(cmd)
	//
	// cmdr.NewString("").
	//	Titles("arch", "arch").
	//	Description("Build for arch, such as 'riscv64'", "").
	//	//ToggleGroup("").
	//	Group(gnControl).
	//	EnvKeys("ARCH").
	//	AttachTo(cmd)

	cmdr.NewString().
		Titles("project-name", "pn", "name").
		Description("Build one project with its name", "").
		// ToggleGroup("").
		Group(gnControl).
		EnvKeys("PROJECT_NAME").
		AttachTo(cmd)

	cmdr.NewInt().
		Titles("parallel", "j").
		Description("TODO: Use parallel building with CPU Core Count, 0 to Auto", "").
		// ToggleGroup("").
		Group(gnControl).
		Placeholder("#").
		Hidden(true). // parallel building not in planning
		EnvKeys("PARALLEL").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("save", "save").
		Description("Save this session as a .bgo.yml", "").
		// ToggleGroup("").
		Group(gnControl).
		EnvKeys("SAVE").
		AttachTo(cmd)

	// cmdr.NewBool().
	//	Titles("yes", "y").
	//	Description("Assume 'yes' while sth has been asking", "").
	//	//ToggleGroup("").
	//	Group(gnControl).
	//	EnvKeys("").
	//	AttachTo(cmd)
}

func cmdrSubCmdInit(root cmdr.OptCmd) {
	initCmd := cmdr.NewSubCmd().Titles("init", "i").
		Description(`scan folder and save <i>result</i> to <code>bgo.yml</code>, as <mark>project settings</mark>`).
		//nolint:lll //no
		// Description(`<del>scan</del> <u><font color="yellow">folder</font></u> and save <i>result</i> to <code>bgo.yml</code>, as <mark>project settings</mark>`).
		Action(initAction).
		AttachTo(root)

	cmdr.NewStringSlice("bgo.yml").
		Titles("output", "o").
		Description("filename to be saved setting as", "").
		// ToggleGroup("").
		Group("").
		AttachTo(initCmd)
}

func cmdrSubCmdBuild(root cmdr.OptCmd) {
	buildCmd := cmdr.NewSubCmd().Titles("build", "b").
		Description("do 'go build' with <code>.bgo.yml</code>", `
			do 'go build' with <code>.bgo.yml</code>.
			
			This command is a synonym to root ('bgo'), that means, type
				'bgo -os linux'
			is equal to 'bgo build -or linux'.`).
		Action(buildAction).
		AttachTo(root)

	// cmdr.NewString("bgo.yml").
	//	Titles("output", "o").
	//	Description("filename to be saved setting as", "").
	//	//ToggleGroup("").
	//	Group("").
	//	AttachTo(buildCmd)

	// cmdrOptsScopes(buildCmd)

	cmdrOptsBuildCommons(buildCmd)
	cmdrOptsBuildControls(buildCmd)
}

func cmdrSubCmdList(root cmdr.OptCmd) {
	listCmd := cmdr.NewSubCmd().Titles("list", "ls").
		Description("list projects in .bgo.yml", `
			list projects in .bgo.yml.
			
			prints them in brief mode.`).
		Action(listAction).
		AttachTo(root)

	cmdr.NewBool().
		Titles("more", "m").
		Description("more details", "").
		// ToggleGroup("").
		Group("").
		VendorHidden(true).
		AttachTo(listCmd)
}
