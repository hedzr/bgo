package cmdr

//goland:noinspection GoNameStartsWithPackageName
const (
	appName   = "your-starter"
	version   = "0.2.5"
	copyright = "your-starter - A devops tool - cmdr series"
	desc      = "your-starter is an effective devops tool. It make an demo application for 'cmdr'"
	longDesc  = `your-starter is an effective devops tool. It make an demo application for 'cmdr'.

To get help for your-starter building options, run 
'your-starter --help', or 'your-starter -h'.
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
