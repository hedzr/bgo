package logic

import (
	"fmt"
	"strconv"

	"github.com/hedzr/log/exec"
	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr"

	"github.com/hedzr/bgo/internal/logic/logx"
)

func cmdrSubCmdRun(root cmdr.OptCmd) {
	cmdr.NewSubCmd().Titles("run", "r").
		Description("synonym to `go run`", `
			synonym to `+"`go run`"+`.`).
		Action(goRunAction).
		TailPlaceholder("", "-- [go-run-arguments...]").
		AttachTo(root)
}

func goRunAction(cmd *cmdr.Command, args []string) (err error) {
	var ec = errors.New("processing executables")
	defer ec.Defer(&err)

	logx.Dim("INVOKING: %v", args)

	c := exec.New().
		WithCommand(asArray2M("go", "run", args...)...).
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

func asArray2M(a, b string, args ...string) (ret []interface{}) {
	ret = append(ret, a, b)
	for _, x := range args {
		ret = append(ret, x)
	}
	return
}

func runOkHandler(ec errors.Error) (onOK func(retCode int, stdoutText string)) {
	return func(retCode int, stdoutText string) {
		if len(stdoutText) > 0 {
			logx.Dim("OUTPUT (Code=%v):\n", retCode)
			logx.Text("%v\n", stdoutText)
		}

		// exec.New().WithCommandString("bash -c 'echo hello world!'", '\'').WithContext(context.Background()).Run()

		// return
	}
}
