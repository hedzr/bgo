package core

import "github.com/hedzr/cmdr"

func AttachToCmdr(root *cmdr.RootCmdOpt) {
	var cmd1 = cmdr.NewCmd().
		Titles("cmd1", "c1", "command1").
		Description("command 1 here", "").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})
	cmd1.AttachTo(root)
	cmdr.NewBool().
		Titles("bool1", "b1", "bool-flag-1").
		Description("bool flag 1 here", "").
		AttachTo(cmd1)
}
