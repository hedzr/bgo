package logic

import (
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log"
	"github.com/hedzr/log/isdelve"
)

func buildAction(cmd *cmdr.Command, args []string) (err error) {
	if cmdr.GetDebugMode() {
		logx.Dim("Debug Mode On, level = %v", log.GetLevel())
	}

	if isdelve.Enabled {
		logx.Dim("isdelve Mode On")
	}

	buildScope := buildScopeFromCmdr(cmd)
	err = actionGoBuild(buildScope)
	return
}

func actionGoBuild(buildScope string) (err error) {
	logx.Verbose("Build Scope: %v, Using main config: %v\n", buildScope, cmdr.GetUsedConfigFile())

	//var buildSettings = new(BgoSettings)
	//err = cmdr.GetSectionFrom("bgo.build", &buildSettings)
	//logDim("build.settings: %+v", buildSettings)
	//logHiLight("Starting...")

	switch buildScope {
	case "short", "current":
		err = buildCurr(buildScope)
	case "full":
		fallthrough
	default:
		err = buildAuto(buildScope)
	}
	return
}
