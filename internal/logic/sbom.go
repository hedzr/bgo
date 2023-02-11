package logic

import (
	"debug/buildinfo"
	"os"
	"path/filepath"

	"github.com/hedzr/bgo/internal/logic/logx"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr"
)

func cmdrSubCmdSBOM(root cmdr.OptCmd) {
	sbom := cmdr.NewSubCmd().Titles("sbom", "s").
		Description("Print SBOM information of this or specified executable", `
			Print SBOM information of this or specified executable(s).
			
			The outputs is YAML compliant.
			
			Just like 'go version -m' but no need to install Go Runtime.`).
		Action(sbomAction).
		TailPlaceholder("<executables...>").
		AttachTo(root)

	cmdr.NewBool().
		Titles("more", "m").
		Description("more details", "").
		// ToggleGroup("").
		Group("").
		VendorHidden(true).
		AttachTo(sbom)
}

func sbomAction(cmd *cmdr.Command, args []string) (err error) {
	var caught bool
	var ec = errors.New("processing executables")
	defer ec.Defer(&err)

	logx.Log("SBOM-LIST:\n")
	for _, file := range args {
		logx.Colored(logx.Green, "  # Checking %v ...", file)
		ec.Attach(sbomOne(file))
		caught = true
	}
	if !caught {
		p, _ := os.Executable()
		p, _ = filepath.Abs(p)
		file := p
		// file := exec.GetExecutablePath()
		logx.Colored(logx.Green, `  # EXECUTABLE: "%v"`, file)
		ec.Attach(sbomOne(file))
	}
	return
}

func sbomOne(file string) (err error) {
	var inf *buildinfo.BuildInfo
	if inf, err = buildinfo.ReadFile(file); err != nil {
		return
	}

	var defaultForeColor = logx.LightGray
	var valColor = func(v any) string {
		return logx.ToColor(logx.White, "%v", v)
	}
	var valSharpColor = func(v any) string {
		return logx.ToColor(logx.White, "%#v", v)
	}
	logx.Colored(defaultForeColor, `  - SBOM:
    executable: "%v"
    go-version: %v
    path: %v
    module-path: %v
    module-version: %v
    module-sum: %v
    module-replace: <ignored>
    settings:
`,
		valColor(file), valColor(inf.GoVersion), valColor(inf.Path),
		valColor(inf.Main.Path), valColor(inf.Main.Version), valColor(inf.Main.Sum),
	)

	for _, d := range inf.Settings {
		logx.Colored(defaultForeColor, "      - %q: %v\n", d.Key, valColor(d.Value))
	}
	logx.Colored(defaultForeColor, "    depends:")
	for _, d := range inf.Deps {
		// str := fmt.Sprintf("%#v", *d)
		logx.Colored(defaultForeColor, `      - debug-module: { path: "%v", version: "%v", sum: "%v", replace: "%v" }
`,
			valColor(d.Path), valColor(d.Version), valColor(d.Sum), valSharpColor(d.Replace))
	}
	return
}
