package cmdr

//goland:noinspection GoNameStartsWithPackageName
const (
	appName   = "bgo"
	version   = "0.2.5"
	copyright = "bgo - A devops tool - cmdr series"
	desc      = "bgo is an effective devops tool. It make an demo application for 'cmdr'"
	longDesc  = `bgo is an effective devops tool. It make an demo application for 'cmdr'.

To get help for bgo building options, run 
'bgo --help', or 'bgo -h'.
`
	examples = `
$ {{.AppName}} gen shell [--bash|--zsh|--auto]
  generate bash/shell completion scripts
$ {{.AppName}} gen man
  generate linux man page 1
$ {{.AppName}} --help
  show help screen.
`
	overview = ``

	zero = 0

	defaultTraceEnabled  = true
	defaultDebugEnabled  = false
	defaultLoggerLevel   = "info"
	defaultLoggerBackend = "logrus"
)
