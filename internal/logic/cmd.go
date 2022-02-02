package logic

import (
	"github.com/hedzr/cmdr"
)

func AttachToCmdr(root cmdr.OptCmd) {
	root.Action(rootAction)

	cmdrOptsScopes(root)
	//cmdrOptsBuildCommons(root)
	//cmdrOptsBuildControls(root)

	cmdrSubCmdInit(root)
	cmdrSubCmdBuild(root)
}

func cmdrOptsScopes(cmd cmdr.OptCmd) {
	cmdr.NewBool().
		Titles("short", "s").
		Description("Build for current CPU and OS Arch ONLY", "").
		ToggleGroup("Scope").
		//Group("Build").
		EnvKeys("SHORT").
		AttachTo(cmd)
	cmdr.NewBool().
		Titles("full", "f").
		Description("Build all CLIs under work directory recursively", "").
		ToggleGroup("Scope").
		//Group("Build").
		EnvKeys("FULL").
		AttachTo(cmd)
	cmdr.NewBool(true).
		Titles("auto", "a").
		Description("Build all modules defined in .bgo.yaml recursively (auto mode)", "").
		ToggleGroup("Scope").
		//Group("Build").
		EnvKeys("AUTO").
		AttachTo(cmd)
}

func cmdrOptsBuildCommons(cmd cmdr.OptCmd) {
	cmdr.NewStringSlice().
		Titles("osarch", "arch", "arch").
		Description("Specify the CPU-Arch list (comma separated or multi times)", `Specify the CPU-Arch list (comma separated or multi times)
		Such as '--arch arm64 --arch amd64,386'.
		The available list can be dumped by 'go tool dist list'`).
		ToggleGroup("").
		Group("Build").
		Placeholder("ARCH").
		EnvKeys("GOARCH").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("os", "os").
		Description("Specify the OS list (comma separated or multi times)", `Specify the OS list (comma separated or multi times)
		Such as '--os linux --os plan9 --os darwin,windows'.
		The available list can be dumped by 'go tool dist list'`).
		ToggleGroup("").
		Group("Build").
		Placeholder("OS").
		EnvKeys("GOOS").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("for", "for").
		Description("Specify the OS/arch list (comma separated or multi times)", `Specify the OS list (comma separated or multi times)
		Such as '--for linux/arm64 --os plan9/amd64'.
		The available list can be dumped by 'go tool dist list'`).
		ToggleGroup("").
		Group("Build").
		Placeholder("OS/ARCH").
		EnvKeys("FOR").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("tags", "tags").
		Description("Additional '-tags' value to pass to go build (comma separated or multi times)", `Additional '-tags' value to pass to go build (comma separated or multi times)
		Such as '--tags isdelve --tags private --os aws-build,lan-build'.`).
		ToggleGroup("").
		Group("Go Build").
		Placeholder("tag,list").
		EnvKeys("").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("ldflags", "ldflags").
		Description("Additional '-ldflags' value to pass to go build (comma separated or multi times)", ``).
		ToggleGroup("").
		Group("Go Build").
		Placeholder("[pattern=]arg list").
		EnvKeys("LDFLAGS").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("gcflags", "gcflags").
		Description("Additional '-gcflags' value to pass to go build (comma separated or multi times)", ``).
		ToggleGroup("").
		Group("Go Build").
		Placeholder("[pattern=]arg list").
		EnvKeys("GCFLAGS").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("gccgoflags", "gccgoflags").
		Description("Additional '-gccgoflags' value to pass to go build (comma separated or multi times)", ``).
		ToggleGroup("").
		Group("Go Build").
		Placeholder("[pattern=]arg list").
		EnvKeys("GCCGOFLAGS").
		AttachTo(cmd)
	cmdr.NewStringSlice().
		Titles("asmflags", "asmflags").
		Description("Additional '-asmflags' value to pass to go build (comma separated or multi times)", ``).
		ToggleGroup("").
		Group("Go Build").
		Placeholder("[pattern=]arg list").
		EnvKeys("ASMFLAGS").
		AttachTo(cmd)

	cmdr.NewString("").
		Titles("mod", "mod").
		Description("Additional '-mod' value to pass to go build", ``).
		ToggleGroup("").
		Group("Go Build").
		Placeholder("MOD").
		ValidArgs("readonly", "vendor", "mod").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("cgo", "cgo").
		Description("enable CGO mode", "").
		//ToggleGroup("").
		Group("Go Build").
		EnvKeys("CGO_ENABLED").
		AttachTo(cmd)
	cmdr.NewBool().
		Titles("race", "race").
		Description("enable --race building", "").
		//ToggleGroup("").
		Group("Go Build").
		EnvKeys("RACE").
		AttachTo(cmd)
	cmdr.NewBool().
		Titles("msan", "msan").
		Description("enable --msan building (enable interoperation with memory sanitizer)", `enable --msan building (enable interoperation with memory sanitizer).
		Supported only on linux/amd64, linux/arm64
		and only with Clang/LLVM as the host C compiler.
		On linux/arm64, pie build mode will be used.`).
		//ToggleGroup("").
		Group("Go Build").
		EnvKeys("MSAN").
		AttachTo(cmd)

	cmdr.NewString("{{.AppName}}-{{.OS}}-{{.ARCH}}").
		Titles("output", "o").
		Description("Specify the binary filename pattern", ``).
		ToggleGroup("").
		Group("Go Build").
		Placeholder("PATTERN").
		AttachTo(cmd)
}

func cmdrOptsBuildControls(cmd cmdr.OptCmd) {
	cmdr.NewBool().
		Titles("rebuild", "r").
		Description("Force rebuilding of package that were up to date", "").
		//ToggleGroup("").
		Group("Control").
		EnvKeys("REBUILD").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("no-trimpath", "ntp").
		Description("Don't use -trimapth", "").
		//ToggleGroup("").
		Group("Control").
		EnvKeys("NO_TRIMPATH").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("dry-run", "dr").
		Description("Just listing without Go build", "").
		//ToggleGroup("").
		Group("Control").
		EnvKeys("DRY_RUN").
		AttachTo(cmd)

	//cmdr.NewString("").
	//	Titles("for", "for").
	//	Description("Build for os/arch, such as 'linux/riscv64'", "").
	//	//ToggleGroup("").
	//	Group("Control").
	//	EnvKeys("FOR").
	//	AttachTo(cmd)
	//
	//cmdr.NewString("").
	//	Titles("os", "os").
	//	Description("Build for os, such as 'linux'", "").
	//	//ToggleGroup("").
	//	Group("Control").
	//	EnvKeys("OS").
	//	AttachTo(cmd)
	//
	//cmdr.NewString("").
	//	Titles("arch", "arch").
	//	Description("Build for arch, such as 'riscv64'", "").
	//	//ToggleGroup("").
	//	Group("Control").
	//	EnvKeys("ARCH").
	//	AttachTo(cmd)

	cmdr.NewString().
		Titles("project-name", "pn", "name").
		Description("Build one project with its name", "").
		//ToggleGroup("").
		Group("Control").
		EnvKeys("PROJECT_NAME").
		AttachTo(cmd)

	cmdr.NewInt().
		Titles("parallel", "j").
		Description("TODO: Use parallel building with CPU Core Count, 0 to Auto", "").
		//ToggleGroup("").
		Group("Control").
		Placeholder("#").
		Hidden(true). // parallel building not in planning
		EnvKeys("PARALLEL").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("save", "save").
		Description("Save this session as a .bgo.yml", "").
		//ToggleGroup("").
		Group("Control").
		EnvKeys("SAVE").
		AttachTo(cmd)

	//cmdr.NewBool().
	//	Titles("yes", "y").
	//	Description("Assume 'yes' while sth has been asking", "").
	//	//ToggleGroup("").
	//	Group("Control").
	//	EnvKeys("").
	//	AttachTo(cmd)
}

func cmdrSubCmdInit(root cmdr.OptCmd) {
	initCmd := root.NewSubCommand("init", "i").
		Description("scan folder and save config file as bgo.yml").
		Action(initAction)

	cmdr.NewString("bgo.yml").
		Titles("output", "o").
		Description("filename to be saved setting as", "").
		//ToggleGroup("").
		Group("").
		AttachTo(initCmd)
}

func cmdrSubCmdBuild(root cmdr.OptCmd) {
	buildCmd := root.NewSubCommand("build", "b").
		Description("do 'go build' with .bgo.yml <default>").
		Action(buildAction)

	//cmdr.NewString("bgo.yml").
	//	Titles("output", "o").
	//	Description("filename to be saved setting as", "").
	//	//ToggleGroup("").
	//	Group("").
	//	AttachTo(buildCmd)

	//cmdrOptsScopes(buildCmd)

	cmdrOptsBuildCommons(buildCmd)
	cmdrOptsBuildControls(buildCmd)
}
