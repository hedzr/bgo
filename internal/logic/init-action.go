package logic

import "github.com/hedzr/cmdr"

func initAction(cmd *cmdr.Command, args []string) (err error) {
	setSaveMode(true)

	buildScope := "full"
	setBuildScope(buildScope)

	on := cmdr.GetStringSliceRP(cmd.GetDottedNamePath(), "output", "bgo.yml")
	cmdr.Set("settings-filename", on)

	err = buildAuto(buildScope)
	return
}
