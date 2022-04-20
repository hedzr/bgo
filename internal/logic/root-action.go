package logic

import ( //nolint:goimports //so what
	"fmt"
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/cmdr"
)

func rootAction(cmd *cmdr.Command, args []string) (err error) { // nolint
	if cmdr.GetTraceMode() {
		cmdr.DebugOutputTildeInfo(true)

		if s := cmdr.GetUsedConfigSubDir(); s != "" {
			fmt.Printf("config.dir: %v\n", s)
		}
		if s := cmdr.GetUsedConfigFile(); s != "" {
			fmt.Printf("config.file: %v\n", s)
		}
		fmt.Printf("build/scope.full: %v\n", cmdr.GetBoolR("full"))
		fmt.Printf("build/scope.short: %v\n", cmdr.GetBoolR("short"))
		fmt.Printf("build/scope: %v\n", cmdr.GetStringR(cmd.GetDottedNamePath()+".Scope")) // show 'app.bgo.Scope'
		fmt.Printf("build/osarch.list: %v\n", cmdr.GetStringSliceR("build.osarch"))
		fmt.Printf("build/os.list: %v\n", cmdr.GetStringSliceR("build.os"))
	}

	if cmdr.GetDebugMode() {
		logx.Dim("Debug Mode On")
		// log.SetLevel(log.DebugLevel)
	}

	if cmdr.GetTraceMode() {
		logx.DimV("Trace Mode On")
		// log.SetLevel(log.TraceLevel)
	}

	err = cmdr.InvokeCommand("build")

	// buildScope := buildScopeFromCmdr(cmd)
	// err = actionGoBuild(buildScope)
	return
}
