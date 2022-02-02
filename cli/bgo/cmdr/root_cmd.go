package cmdr

import (
	"github.com/hedzr/bgo/internal/logic"
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/conf"
)

func buildRootCmd() (rootCmd *cmdr.RootCommand) {
	ver := version
	if conf.ServerID != "" {
		ver = conf.Version
	}
	root := cmdr.Root(appName, ver).
		AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
			// fmt.Printf("# global pre-action 1, curr-dir: %v\n", cmdr.GetCurrentDir())
			// cmdr.Set("enable-ueh", true)
			//err = internal.App().Init(cmd, args) // App() will be auto-closed
			logx.Verbose("VERBOSE: %v, DEBUG: %v\n", logx.ToDim("%v", cmdr.GetVerboseMode()), logx.ToDim("%v", cmdr.GetDebugMode()))
			logx.Verbose("Config File: %v\n", logx.ToDim("%v", cmdr.GetUsedConfigFile()))
			logx.Verbose("Config File (2ndry): %v\n", logx.ToDim("%v", cmdr.GetUsedSecondaryConfigFile()))
			logx.Verbose("Config File (alter): %v\n", logx.ToDim("%v", cmdr.GetUsedAlterConfigFile()))
			return
		}).
		//AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
		//	//fmt.Printf("# global pre-action 2, exe-path: %v\n", cmdr.GetExecutablePath())
		//	return
		//}).
		//AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
		//	//fmt.Println("# global post-action 1")
		//}).
		//AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
		//	//fmt.Println("# global post-action 2")
		//}).
		Copyright(copyright, "bgo Authors").
		Description(desc, longDesc).
		Examples(examples)
	rootCmd = root.RootCommand()

	// core.AttachToCmdr(root.RootCmdOpt())

	logic.AttachToCmdr(root.RootCmdOpt())

	//cmdr.NewBool(false).
	//	Titles("enable-ueh", "ueh").
	//	Description("Enables the unhandled exception handler?").
	//	AttachTo(root)
	//soundex(root)
	//panicTest(root)

	//pprof.AttachToCmdr(root.RootCmdOpt())

	return
}
