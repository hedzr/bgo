module github.com/hedzr/bgo

go 1.18

// replace gopkg.in/hedzr/errors.v3 => ../../go-cmdr/05.errors

// replace github.com/hedzr/log => ../../go-cmdr/10.log

//replace github.com/hedzr/logex => ../../cmdr-series/libs/logex

//replace github.com/hedzr/cmdr => ../../cmdr-series/cmdr

// replace github.com/hedzr/cmdr-addons => ../53.cmdr-addons

require (
	github.com/hedzr/cmdr v1.11.8
	github.com/hedzr/evendeep v0.3.1
	github.com/hedzr/log v1.6.0
	github.com/hedzr/logex v1.6.0
	golang.org/x/crypto v0.6.0
	golang.org/x/mod v0.4.2
	gopkg.in/hedzr/errors.v3 v3.1.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/hedzr/cmdr-base v1.0.0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/kr/text v0.1.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/net v0.6.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/term v0.5.0 // indirect
)
