package logx

//nolint:goimports //i like it
import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/hedzr/log"

	"github.com/hedzr/cmdr"
)

type lxS struct {
	log.Logger
}

var onceLxInitializer sync.Once //nolint:gochecknoglobals //no
var lx *lxS                     //nolint:gochecknoglobals //no

const extrasLogSkip = 4

// LazyInit initials local lx instance properly.
// While you are storing a log.Logger copy locally, DO NOT put these
// codes into func init() since it's too early to get a unconfigurec
// log.Logger.
func LazyInit() { lazyInit() }
func lazyInit() *lxS {
	onceLxInitializer.Do(func() {
		lx = &lxS{
			log.Skip(extrasLogSkip),
		}
	})
	return lx
}

func IsVerboseMode() bool { return cmdr.GetVerboseMode() }
func CountOfVerbose() int { return cmdr.GetVerboseModeHitCount() }

// Error outputs formatted message to stderr.
func Error(format string, args ...interface{}) { //nolint:goprintffuncname //so what
	_internalLogTo(func(sb strings.Builder, ln bool) {
		if ln {
			_, _ = fmt.Fprint(os.Stderr, sb.String())
		} else {
			// _, _ = fmt.Fprintln(os.Stderr, sb.String())
			lx.Errorf("%v", sb.String())
		}
	}, format, args...)
}

// Fatal outputs formatted message to stderr.
func Fatal(format string, args ...interface{}) { //nolint:goprintffuncname //so what
	_internalLogTo(func(sb strings.Builder, ln bool) {
		if ln {
			lx.Fatalf("%v", sb.String())
		} else {
			lx.Panicf("%v", sb.String())
		}
	}, format, args...)
}

// Warn outputs formatted message to stderr while logger level
// less than log.WarnLevel.
// For log.SetLevel(log.ErrorLevel), the text will be discarded.
func Warn(format string, args ...interface{}) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(log.ErrorLevel)
	if log.GetLevel() < log.WarnLevel {
		return
	}

	_internalLogTo(func(sb strings.Builder, ln bool) {
		if ln {
			print(ToColor(Yellow, sb.String())) //nolint:forbidigo //no
		} else {
			// println(sb.String())
			lx.Warnf("%v", ToColor(Yellow, sb.String()))
		}
	}, format, args...)
}

// Log will print the formatted message to stdout.
//
// While the message ends with '\n', it will be printed by
// print(), so you'll see it always.
// But if not, the message will be printed by hedzr/log. In
// this case, its outputting depends on
// hedzr/log.GetLogLevel() >= log.DebugLevel.
//
// Log outputs formatted message to stdout while logger level
// less than log.WarnLevel.
// For log.SetLevel(log.ErrorLevel), the text will be discarded.
func Log(format string, args ...interface{}) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(log.ErrorLevel)
	if log.GetLevel() < log.WarnLevel {
		return
	}

	_internalLogTo(func(sb strings.Builder, ln bool) {
		if ln {
			print(sb.String()) //nolint:forbidigo //no
		} else {
			// println(sb.String())
			lx.Printf("%v", sb.String())
		}
	}, format, args...)
}

// Verbose outputs formatted message to stdout while cmdr is in
// VERBOSE mode.
// For log.SetLevel(log.ErrorLevel), the text will be discarded.
func Verbose(format string, args ...interface{}) { //nolint:goprintffuncname //so what
	if IsVerboseMode() {
		_internalLogTo(func(sb strings.Builder, ln bool) {
			if ln {
				print(sb.String()) //nolint:forbidigo //no
			} else {
				// println(sb.String())
				lx.Printf("%v", sb.String())
			}
		}, format, args...)
	}
}

// Trace outputs formatted message to stdout while logger level
// is log.TraceLevel, or cmdr is in TRACE mode or trace module
// is enabled.
func Trace(format string, args ...interface{}) { //nolint:goprintffuncname //so what
	if log.GetLevel() == log.TraceLevel || !cmdr.GetTraceMode() {
		return
	}

	_internalLogTo(func(sb strings.Builder, ln bool) {
		if ln {
			// log.Skip(extrasLogSkip).Tracef("%v", sb.String())
			Colored(LightGray, "%v", sb.String())
		} else {
			lx.Tracef("%v", sb.String())
		}
	}, format, args...)
}

// Hilight outputs formatted message to stdout while logger level
// less than log.WarnLevel.
// For log.SetLevel(log.ErrorLevel), the text will be discarded.
func Hilight(format string, args ...interface{}) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(log.ErrorLevel)
	if log.GetLevel() < log.WarnLevel {
		return
	}

	_internalLogTo(func(sb strings.Builder, ln bool) {
		if cmdr.GetNoColorMode() {
			_, _ = fmt.Fprintf(os.Stdout, "%v", sb.String())
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "\x1b[0;1m%v\x1b[0m", sb.String())
		}
		if !ln {
			println() //nolint:forbidigo //no
		}
	}, format, args...)
}

// DimV outputs formatted message to stdout while logger level less
// than log.WarnLevel and cmdr is in verbose mode.
//
// For example, after log.SetLevel(log.ErrorLevel), the text via DimV will be discarded.
//
// While env-var VERBOSE=1, the text via DimV will be shown.
func DimV(format string, args ...interface{}) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(log.ErrorLevel)
	if log.GetLevel() < log.WarnLevel {
		return
	}

	if IsVerboseMode() {
		_internalLogTo(func(sb strings.Builder, ln bool) {
			if cmdr.GetNoColorMode() {
				_, _ = fmt.Fprintf(os.Stdout, "%v", sb.String())
			} else {
				_, _ = fmt.Fprintf(os.Stdout, "\x1b[2m\x1b[37m%v\x1b[0m", sb.String())
			}
			if !ln {
				println() //nolint:forbidigo //no
			}
		}, format, args...)
	}
}

// Text prints formatted message without any predefined ansi escaping.
func Text(format string, args ...interface{}) { //nolint:goprintffuncname //so what
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

// Dim outputs formatted message to stdout while logger level
// less than log.WarnLevel.
//
// For example, after log.SetLevel(log.ErrorLevel), the text via Dim will be discarded.
func Dim(format string, args ...interface{}) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(log.ErrorLevel)
	if log.GetLevel() < log.WarnLevel {
		return
	}

	_internalLogTo(func(sb strings.Builder, ln bool) {
		if cmdr.GetNoColorMode() {
			_, _ = fmt.Fprintf(os.Stdout, "%v", sb.String())
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "\x1b[2m\x1b[37m%v\x1b[0m", sb.String())
		}
		if !ln {
			println() //nolint:forbidigo //no
		}
	}, format, args...)
}

func ToDim(format string, args ...interface{}) (str string) {
	str = fmt.Sprintf(format, args...)
	if cmdr.GetNoColorMode() {
		return
	}
	str = fmt.Sprintf("\x1b[2m\x1b[37m%v\x1b[0m", str)
	return
}

func ToColor(clr Color, format string, args ...interface{}) (str string) {
	str = fmt.Sprintf(format, args...)
	if cmdr.GetNoColorMode() {
		return
	}
	str = fmt.Sprintf("\u001B[%dm%v\x1b[0m", clr, str)
	return
}

// ColoredV outputs formatted message to stdout while logger level
// less than log.WarnLevel and cmdr is in VERBOSE mode.
func ColoredV(clr Color, format string, args ...interface{}) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(log.ErrorLevel)
	if log.GetLevel() < log.WarnLevel {
		return
	}

	if IsVerboseMode() {
		_internalLogTo(func(sb strings.Builder, ln bool) {
			if cmdr.GetNoColorMode() {
				_, _ = fmt.Fprintf(os.Stdout, "%v", sb.String())
			} else {
				color(clr)
				_, _ = fmt.Fprintf(os.Stdout, "%v\x1b[0m", sb.String())
			}
			if !ln {
				println() //nolint:forbidigo //no
			}
		}, format, args...)
	}
}

// Colored outputs formatted message to stdout while logger level
// less than log.WarnLevel.
// For log.SetLevel(log.ErrorLevel), the text will be discarded.
func Colored(clr Color, format string, args ...interface{}) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(log.ErrorLevel)
	if log.GetLevel() < log.WarnLevel {
		return
	}

	_internalLogTo(func(sb strings.Builder, ln bool) {
		if cmdr.GetNoColorMode() {
			_, _ = fmt.Fprintf(os.Stdout, "%v", sb.String())
		} else {
			color(clr)
			_, _ = fmt.Fprintf(os.Stdout, "%v\x1b[0m", sb.String())
		}
		if !ln {
			println() //nolint:forbidigo //no
		}
	}, format, args...)
}

func color(c Color) {
	_, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
}

func resetColor(c Color) { //nolint:unused //no
	_, _ = fmt.Fprint(os.Stdout, "\x1b[0m")
}

//nolint:lll //no
func _internalLogTo(tofn func(sb strings.Builder, ln bool), format string, args ...interface{}) { //nolint:goprintffuncname //so what
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
