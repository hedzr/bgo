package logic

import (
	"debug/buildinfo"
	"fmt"
	"github.com/hedzr/cmdr"
	"gopkg.in/hedzr/errors.v3"
)

func cmdrSubCmdSBOM(root cmdr.OptCmd) {
	sbom := cmdr.NewSubCmd().Titles("sbom", "s").
		Description("show SBOM of executable", `
			show SBOM of executable.
			
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
	var ec = errors.New("processing executables")
	for _, file := range args {
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
