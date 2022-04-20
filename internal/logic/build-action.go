package logic

import (
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log"
	"github.com/hedzr/log/isdelve"
)

func listAction(cmd *cmdr.Command, args []string) (err error) {
	return buildAction(cmd, args)
}

func buildAction(cmd *cmdr.Command, args []string) (err error) {
	if cmdr.GetDebugMode() {
		logx.Dim("Debug Mode On, level = %v", log.GetLevel())
	}

	if isdelve.Enabled {
		logx.Dim("isdelve Mode On")
	}

	buildScope := buildScopeFromCmdr(cmd)
	err = actionGoBuild(buildScope, cmd, args)
	return
}

func actionGoBuild(buildScope string, cmd *cmdr.Command, args []string) (err error) {
	logx.Verbose("Build Scope: %v, Using main config: %v\n", buildScope, cmdr.GetUsedConfigFile())

	// var buildSettings = new(BgoSettings)
	// err = cmdr.GetSectionFrom("bgo.build", &buildSettings)
	// logDim("build.settings: %+v", buildSettings)
	// logHiLight("Starting...")

	switch buildScope {
	case "short", "current": //nolint:goconst //i like it
		err = buildCurr(buildScope, cmd, args)
	// case "full":
	// 	fallthrough
	default:
		err = buildAuto(buildScope, cmd, args)
	}
	return
}
