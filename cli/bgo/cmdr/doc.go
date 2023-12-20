package cmdr

//goland:noinspection GoNameStartsWithPackageName,GoUnusedConst
const (
	appName   = "bgo"
	version   = "0.5.25"
	copyright = "bgo - A super easy way to build your go apps - cmdr series"
	desc      = "bgo provides a super easy way to build your go apps"
	longDesc  = `bgo provides a super easy way to build your go apps.

To get help for bgo building options, run 
'bgo build --help', or 'bgo b -h'.
`
	examples = `
$ {{.AppName}} init
  generate bgo.yml. Please rename to .bgo.yml so that {{.AppName}} can autoload it
$ {{.AppName}}
  run 'go build' in accord with .bgo.yml
$ {{.AppName}} -s
  run 'go build' for current GOOS/GOARCH
$ {{.AppName}} --for linux/riscv64
  run 'go build' for specified os/arch
$ {{.AppName}} --for linux/riscv64 --for linux/386
  run 'go build' for specified os/arch list
$ {{.AppName}} --help
  show help screen.`
	//nolint:unused,varcheck //like it
	examplesLong = `
$ {{.AppName}} gen shell [--bash|--zsh|--auto]
  generate bash/shell completion scripts
$ {{.AppName}} gen man
  generate linux man page 1
$ {{.AppName}} --help
  show help screen.
`
	//nolint:unused,varcheck //like it
	overview = ``

	zero = 0

	//nolint:varcheck //like it
	defaultTraceEnabled  = true
	defaultDebugEnabled  = false
	defaultLoggerLevel   = "info"
	defaultLoggerBackend = "logrus"
)
