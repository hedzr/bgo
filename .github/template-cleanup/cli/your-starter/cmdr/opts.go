package cmdr

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/pprof"
	"github.com/hedzr/log"
	"github.com/hedzr/log/isdelve"
	"github.com/hedzr/logex/build"
)

var options []cmdr.ExecOption

func init() {
	options = append(options, cmdr.WithUnhandledErrorHandler(onUnhandledErrorHandler))
	//cmdr.WithUnhandledErrorHandler(onUnhandledErrorHandler)

	options = append(options, cmdr.WithLogx(build.New(build.NewLoggerConfigWith(defaultDebugEnabled, defaultLoggerBackend, defaultLoggerLevel, log.WithTimestamp(true, "")))))

	options = append(options, cmdr.WithHelpTailLine(`
# Type '-h'/'-?' or '--help' to get command help screen.
# Star me if it's helpful: https://github.com/%REPOSITORY%
`))

	if isDebugBuild() {
		options = append(options, pprof.GetCmdrProfilingOptions())
	}

	//dex.WithDaemon(
	//	svr.NewDaemon(svr.WithRouterImpl(sth.NewGinMux())),
	//	dex.WithCommandsModifier(modifier),
	//	dex.WithLoggerForward(true),
	//),
	// server.WithCmdrDaemonSupport(),
	// server.WithCmdrHook(),

	// enable '--trace' command line option to toggle a internal trace mode (can be retrieved by cmdr.GetTraceMode())
	// import "github.com/hedzr/cmdr-addons/pkg/plugins/trace"
	// trace.WithTraceEnable(defaultTraceEnabled)
	optAddTraceOption := cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
		cmdr.NewBool(false).
			Titles("trace", "tr").
			Description("enable trace mode for tcp/mqtt send/recv data dump", "").
			//Action(func(cmd *cmdr.Command, args []string) (err error) {
			//	println("trace mode on")
			//	cmdr.SetTraceMode(true)
			//	return
			//}).
			Group(cmdr.SysMgmtGroup).
			AttachToRoot(root)
	}, nil)
	options = append(options, optAddTraceOption)
	//options = append(options, optAddServerExtOpt«ion)

	//// allow and search '.bgo.yml' at first
	//locations := []string{".$APPNAME.yml"}
	//locations = append(locations, cmdr.GetPredefinedLocations()...)
	//options = append(options, cmdr.WithPredefinedLocations(locations...))
}

func isDebugBuild() bool { return isdelve.Enabled }

//var optAddTraceOption cmdr.ExecOption
//var optAddServerExtOption cmdr.ExecOption

//func init() {
//	//// attaches `--trace` to root command
//	//optAddTraceOption = cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
//	//	cmdr.NewBool(false).
//	//		Titles("trace", "tr").
//	//		Description("enable trace mode for tcp/mqtt send/recv data dump", "").
//	//		AttachToRoot(root)
//	//}, nil)
//
//	//// the following statements show you how to attach an option to a sub-command
//	//optAddServerExtOption = cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
//	//	serverCmd := cmdr.FindSubCommandRecursive("server", nil)
//	//	serverStartCmd := cmdr.FindSubCommand("start", serverCmd)
//	//	cmdr.NewInt(5100).
//	//		Titles("vnc-server", "vnc").
//	//		Description("start as a vnc server (just a faked demo)", "").
//	//		Placeholder("PORT").
//	//		AttachTo(cmdr.NewCmdFrom(serverStartCmd))
//	//}, nil)
//}
