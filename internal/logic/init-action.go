package logic

import (
	"github.com/hedzr/cmdr"

	"github.com/hedzr/bgo/internal/logic/logx"
)

func initAction(cmd *cmdr.Command, args []string) (err error) {
	setSaveMode(true)

	buildScope := "full"
	setBuildScope(buildScope)

	on := cmdr.GetStringSliceRP(cmd.GetDottedNamePath(), "output", "bgo.yml")
	cmdr.Set("settings-filename", on)

	logx.Log("Init Action, scope = %v, cfg = %v", buildScope, cmdr.GetUsedConfigFile())

	err = buildAuto(buildScope, cmd, args)
	return
}
