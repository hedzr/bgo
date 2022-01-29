package logx

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log"
	"os"
	"strings"
)

func Error(format string, args ...interface{}) {
	_internalLogTo(func(sb strings.Builder, ln bool) {
		if ln {
			_, _ = fmt.Fprint(os.Stderr, sb.String())
		} else {
			//_, _ = fmt.Fprintln(os.Stderr, sb.String())
			log.Skip(extrasLogSkip).Errorf("%v", sb.String())
		}
	}, format, args...)
}

func Fatal(format string, args ...interface{}) {
	_internalLogTo(func(sb strings.Builder, ln bool) {
		if ln {
			log.Skip(extrasLogSkip).Fatalf("%v", sb.String())
		} else {
			log.Skip(extrasLogSkip).Fatalf("%v", sb.String())
		}
	}, format, args...)
}

const extrasLogSkip = 4

// Log will print the formatted message to stdout.
//
// While the message ends with '\n', it will be printed by
// print(), so you'll see it always.
// But if not, the message will be printed by hedzr/log. In
// this case, its outputting depends on
// hedzr/log.GetLogLevel() >= log.DebugLevel.
func Log(format string, args ...interface{}) {
	_internalLogTo(func(sb strings.Builder, ln bool) {
		if ln {
			print(sb.String())
		} else {
			//println(sb.String())
			log.Skip(extrasLogSkip).Debugf("%v", sb.String())
		}
	}, format, args...)
}

func Warn(format string, args ...interface{}) {
	_internalLogTo(func(sb strings.Builder, ln bool) {
		if ln {
			print(ToColor(Yellow, sb.String()))
		} else {
			//println(sb.String())
			log.Skip(extrasLogSkip).Warnf("%v", ToColor(Yellow, sb.String()))
		}
	}, format, args...)
}

func Verbose(format string, args ...interface{}) {
	if cmdr.GetVerboseMode() {
		_internalLogTo(func(sb strings.Builder, ln bool) {
			if ln {
				print(sb.String())
			} else {
				//println(sb.String())
				log.Skip(extrasLogSkip).Debugf("%v", sb.String())
			}
		}, format, args...)
	}
}

func Trace(format string, args ...interface{}) {
	_internalLogTo(func(sb strings.Builder, ln bool) {
		if ln {
			if cmdr.GetTraceMode() {
				//log.Skip(extrasLogSkip).Tracef("%v", sb.String())
				Colored(LightGray, "%v", sb.String())
			}
		} else {
			log.Skip(extrasLogSkip).Tracef("%v", sb.String())
		}
	}, format, args...)
}

func Hilight(format string, args ...interface{}) {
	_internalLogTo(func(sb strings.Builder, ln bool) {
		_, _ = fmt.Fprintf(os.Stdout, "\x1b[0;1m%v\x1b[0m", sb.String())
		if !ln {
			println()
		}
	}, format, args...)
}

func DimV(format string, args ...interface{}) {
	if cmdr.GetVerboseMode() {
		_internalLogTo(func(sb strings.Builder, ln bool) {
			_, _ = fmt.Fprintf(os.Stdout, "\x1b[2m\x1b[37m%v\x1b[0m", sb.String())
			if !ln {
				println()
			}
		}, format, args...)
	}
}

func Dim(format string, args ...interface{}) {
	_internalLogTo(func(sb strings.Builder, ln bool) {
		_, _ = fmt.Fprintf(os.Stdout, "\x1b[2m\x1b[37m%v\x1b[0m", sb.String())
		if !ln {
			println()
		}
	}, format, args...)
}

func ToDim(format string, args ...interface{}) string {
	str := fmt.Sprintf(format, args...)
	str = fmt.Sprintf("\x1b[2m\x1b[37m%v\x1b[0m", str)
	return str
}

func ToColor(clr Color, format string, args ...interface{}) string {
	str := fmt.Sprintf(format, args...)
	str = fmt.Sprintf("\u001B[%dm%v\x1b[0m", clr, str)
	return str
}

func ColoredV(clr Color, format string, args ...interface{}) {
	if cmdr.GetVerboseMode() {
		_internalLogTo(func(sb strings.Builder, ln bool) {
			color(clr)
			_, _ = fmt.Fprintf(os.Stdout, "%v\x1b[0m", sb.String())
			if !ln {
				println()
			}
		}, format, args...)
	}
}

func Colored(clr Color, format string, args ...interface{}) {
	_internalLogTo(func(sb strings.Builder, ln bool) {
		color(clr)
		_, _ = fmt.Fprintf(os.Stdout, "%v\x1b[0m", sb.String())
		if !ln {
			println()
		}
	}, format, args...)
}

func color(c Color) {
	_, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
}

func resetColor(c Color) {
	_, _ = fmt.Fprint(os.Stdout, "\x1b[0m")
}

func _internalLogTo(tofn func(sb strings.Builder, ln bool), format string, args ...interface{}) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(format, args...))
	tofn(sb, strings.HasSuffix(sb.String(), "\n"))
}

type Color int

const (
	Black        Color = 30
	Red          Color = 31
	Green        Color = 32
	Yellow       Color = 33
	Blue         Color = 34
	Magenta      Color = 35
	Cyan         Color = 36
	LightGray    Color = 37
	DarkGray     Color = 90
	LightRed     Color = 91
	LightGreen   Color = 92
	LightYellow  Color = 93
	LightBlue    Color = 94
	LightMagenta Color = 95
	LightCyan    Color = 96
	White        Color = 97

	BgNormal       Color = 0
	BgBoldOrBright Color = 1
	BgDim          Color = 2
	BgUnderline    Color = 4
	BgUlink        Color = 5
	BgHidden       Color = 8

	DarkColor = LightGray
)
