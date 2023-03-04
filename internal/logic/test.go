package logic

import (
	"fmt"
	"strconv"

	"github.com/hedzr/log/exec"
	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr"

	"github.com/hedzr/bgo/internal/logic/logx"
)

func cmdrSubCmdTest(root cmdr.OptCmd) {
	cmdr.NewSubCmd().Titles("test", "t").
		Description("synonym to `go test`", `
			synonym to `+"`go test`"+`.`).
		Action(goTestAction).
		TailPlaceholder("", "-- [go-test-arguments...]").
		AttachTo(root)
}

func goTestAction(cmd *cmdr.Command, args []string) (err error) {
	var ec = errors.New("processing executables")
	defer ec.Defer(&err)

	logx.Dim("INVOKING: %v", args)

	c := exec.New().
		WithCommand(asArray2M("go", "test", args...)...).
		WithOnOK(runOkHandler(ec)).
		WithOnError(func(err error, retCode int, stdoutText, stderrText string) {
			logx.Error("ERROR TEXT:\n%v\nError:\n%v\nRetCode: %v\nCommands: %v\n",
				logx.ToColor(logx.Red, leftPad(stderrText, 4)),  //nolint:gomnd //so what
				logx.ToColor(logx.Red, leftPad(err.Error(), 4)), //nolint:gomnd //so what
				logx.ToDim(strconv.Itoa(retCode)),
				logx.ToDim(fmt.Sprintf("%v", cmd)))
		})
	ec.Attach(c.RunAndCheckError())
	if !ec.IsEmpty() {
		// caller will discard goBuild error, so we print it to notify end-user.
		logx.Error("Error occurs: %v", ec)
	}
	return
}
