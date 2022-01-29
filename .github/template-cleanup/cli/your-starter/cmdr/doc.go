package cmdr

//goland:noinspection GoNameStartsWithPackageName
const (
	appName   = "%NAME%"
	version   = "0.2.5"
	copyright = "%NAME% - A devops tool - cmdr series"
	desc      = "%NAME% is an effective devops tool. It make an demo application for 'cmdr'"
	longDesc  = `%NAME% is an effective devops tool. It make an demo application for 'cmdr'.

To get help for %NAME% building options, run 
'%NAME% --help', or '%NAME% -h'.
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
