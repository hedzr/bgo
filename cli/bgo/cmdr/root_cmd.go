package cmdr

//nolint:goimports //so what
import (
	"os"
	"strings"

	"github.com/hedzr/bgo/internal/logic"
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/conf"
)

//nolint:nakedret //so what
func buildRootCmd() (rootCmd *cmdr.RootCommand) {
	ver := version

	if //goland:noinspection ALL
	conf.ServerID != "" {
		ver = conf.Version
	}

	root := cmdr.Root(appName, ver).
		AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
			// fmt.Printf("# global pre-action 1, curr-dir: %v\n", cmdr.GetCurrentDir())
			// cmdr.Set("enable-ueh", true)
			// err = internal.App().Init(cmd, args) // App() will be auto-closed

			// our local logx initialize itself here. logx is an enhancement and wrapper for hedzr/log(ex)
			logx.LazyInit()

			logx.Verbose("VERBOSE: %v, DEBUG: %v", logx.ToDim("%v", logx.IsVerboseMode()), logx.ToDim("%v", cmdr.GetDebugMode()))
			logx.Verbose("Config File: %v\n", logx.ToDim("%v", cmdr.GetUsedConfigFile()))
			logx.Verbose("Config File (2ndry): %v\n", logx.ToDim("%v", cmdr.GetUsedSecondaryConfigFile()))
			logx.Verbose("Config File (alter): %v\n", logx.ToDim("%v", cmdr.GetUsedAlterConfigFile()))
			logx.Verbose("Args: %v\n", strings.Join(os.Args, ", "))
			logx.Verbose("cmd: %v, args: %v\n", cmd.GetTitleName(), args)
			if logx.CountOfVerbose() > 1 {
				logx.Error("an error ok, 1")
				logx.Warn("a warn ok, 1")
				logx.Log("a log ok, 1")
				logx.Verbose("a verbose ok, 1")
			}
			return
		}).
		// AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
		//	//fmt.Printf("# global pre-action 2, exe-path: %v\n", cmdr.GetExecutablePath())
		//	return
		// }).
		// AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
		//	//fmt.Println("# global post-action 1")
		// }).
		// AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
		//	//fmt.Println("# global post-action 2")
		// }).
		Copyright(copyright, "bgo Authors").
		Description(desc, longDesc).
		Examples(examples)
	rootCmd = root.RootCommand()

	// core.AttachToCmdr(root.RootCmdOpt())

	logic.AttachToCmdr(root.RootCmdOpt())
	cmdrMoreCommandsForTest(root.RootCmdOpt())

	// pprof.AttachToCmdr(root.RootCmdOpt())

	return
}
