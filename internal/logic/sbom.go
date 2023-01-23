package logic

import (
	"debug/buildinfo"
	"fmt"
	"github.com/hedzr/bgo/internal/logic/logx"
	"os"
	"path/filepath"

	"github.com/hedzr/cmdr"
	"gopkg.in/hedzr/errors.v3"
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
	for _, file := range args {
		logx.Colored(logx.Green, "checking %v ...", file)
		ec.Attach(sbomOne(file))
		caught = true
	}
	if !caught {
		p, _ := os.Executable()
		p, _ = filepath.Abs(p)
		file := p
		//file := exec.GetExecutablePath()
		logx.Colored(logx.Green, "SBOM on %v", file)
		ec.Attach(sbomOne(file))
	}
	return
}

func sbomOne(file string) (err error) {
	var inf *buildinfo.BuildInfo
	if inf, err = buildinfo.ReadFile(file); err != nil {
		return
	}

	fmt.Printf(`SBOM:
  executable: %q
  go-version: %v
  path: %v
  module-path: %v
  module-version: %v
  module-sum: %v
  module-replace: <ignored>
  settings:
`,
		file, inf.GoVersion, inf.Path,
		inf.Main.Path, inf.Main.Version, inf.Main.Sum,
	)

	for _, d := range inf.Settings {
		fmt.Printf("    - %q: %v\n", d.Key, d.Value)
	}
	fmt.Println("  depends:")
	for _, d := range inf.Deps {
		// str := fmt.Sprintf("%#v", *d)
		fmt.Printf("    - debug-module: { path: %q, version: %q, sum: %q, replace: %#v } \n", d.Path, d.Version, d.Sum, d.Replace)
	}
	return
}
