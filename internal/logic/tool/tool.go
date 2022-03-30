package tool

//nolint:goimports
import (
	"github.com/hedzr/cmdr/tool/randomizer"
	"regexp"
)

func StripOrderPrefix(s string) string {
	if len(s) == 0 {
		return s
	}
	a := xre.FindStringSubmatch(s)
	return a[2]
	// if xre.MatchString(s) {
	//	s = s[strings.Index(s, ".")+1:]
	// }
	// return s
}

// HasOrderPrefix tests whether an order prefix is present or not.
// An order prefix is a dotted string with multiple alphabet and digit. Such as:
// "zzzz.", "0001.", "700.", "A1." ...
func HasOrderPrefix(s string) bool {
	return xre.MatchString(s)
}

// NextIn returns a random number with upper bound 'max'
func NextIn(max int) int {
	return rr.NextIn(max)
}

// NextString returns a random string with max length 'length'
func NextString(length int) string {
	return rr.AsStrings().NextString(length)
}

var (
	xre = regexp.MustCompile(`^([0-9A-Za-z]+[.-])?(.+)$`)
	rr  = randomizer.New()
)
