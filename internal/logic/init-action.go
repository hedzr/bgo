package logic

import "github.com/hedzr/cmdr"

func initAction(cmd *cmdr.Command, args []string) (err error) {
	setSaveMode(true)

	buildScope := "full"
	setBuildScope(buildScope)

	on := cmdr.GetStringRP(cmd.GetDottedNamePath(), "output", "bgo.yml")
	cmdr.Set("settings-filename", on)

	err = buildFull(buildScope)
	return
}
