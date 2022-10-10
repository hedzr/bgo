module github.com/hedzr/bgo

go 1.18

// replace gopkg.in/hedzr/errors.v3 => ../../go-cmdr/05.errors

// replace github.com/hedzr/log => ../../go-cmdr/10.log

// replace github.com/hedzr/logex => ../../go-cmdr/15.logex

// replace github.com/hedzr/cmdr => ../../go-cmdr/50.cmdr

// replace github.com/hedzr/cmdr-addons => ../53.cmdr-addons

require (
	github.com/hedzr/cmdr v1.10.50
	github.com/hedzr/log v1.5.56
	github.com/hedzr/logex v1.5.56
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	gopkg.in/hedzr/errors.v3 v3.0.23
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/hedzr/cmdr-base v0.1.3 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
	golang.org/x/term v0.0.0-20220526004731-065cf7ba2467 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)
