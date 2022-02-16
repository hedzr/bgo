module github.com/hedzr/bgo

go 1.17

//replace github.com/hedzr/log => ../../go-cmdr/10.log

//replace github.com/hedzr/logex => ../../go-cmdr/15.logex

//replace github.com/hedzr/cmdr => ../../go-cmdr/50.cmdr

//replace github.com/hedzr/cmdr-addons => ../53.cmdr-addons

require (
	github.com/hedzr/cmdr v1.10.24
	github.com/hedzr/log v1.5.17
	github.com/hedzr/logex v1.5.17
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	gopkg.in/hedzr/errors.v2 v2.2.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/hedzr/cmdr-base v0.1.3 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)
