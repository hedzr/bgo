package logic

import ( //nolint:goimports //so what

	"log"

	"github.com/hedzr/cmdr"

	"github.com/hedzr/bgo/internal/logx"
)

func rootAction(cmd *cmdr.Command, args []string) (err error) { //nolint:unused //no
	if cmdr.GetTraceMode() {
		cmdr.DebugOutputTildeInfo(true)

		if s := cmdr.GetUsedConfigSubDir(); s != "" {
			log.Printf("config.dir: %v\n", s)
		}
		if s := cmdr.GetUsedConfigFile(); s != "" {
			log.Printf("config.file: %v\n", s)
		}
		log.Printf("build/scope.full: %v\n", cmdr.GetBoolR("full"))
		log.Printf("build/scope.short: %v\n", cmdr.GetBoolR("short"))
		log.Printf("build/scope: %v\n", cmdr.GetStringR(cmd.GetDottedNamePath()+".Scope")) // show 'app.bgo.Scope'
		log.Printf("build/osarch.list: %v\n", cmdr.GetStringSliceR("build.osarch"))
		log.Printf("build/os.list: %v\n", cmdr.GetStringSliceR("build.os"))
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
