package cmdr

import (
	"github.com/hedzr/bgo/internal/logic"
	"github.com/hedzr/cmdr"
)

func buildRootCmd() (rootCmd *cmdr.RootCommand) {
	root := cmdr.Root(appName, version).
		//AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
		//	// fmt.Printf("# global pre-action 1, curr-dir: %v\n", cmdr.GetCurrentDir())
		//	// cmdr.Set("enable-ueh", true)
		//	err = internal.App().Init(cmd, args) // App() will be auto-closed
		//	return
		//}).
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
