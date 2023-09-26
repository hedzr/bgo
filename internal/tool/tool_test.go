package tool

import (
	"regexp"
	"testing"
)

type testCase struct {
	in, expect string
}
type testCaseR struct {
	expect, in string
}

func TestHasOrderPrefix(t *testing.T) {
	xre = regexp.MustCompile(`^([0-9A-Za-z]+[.-]){1}(.+)$`)
	if xre.MatchString("LogIt") {
		t.Error("BAD!")
	}

	for i, c := range []struct {
		in  string
		has bool
		out string
	}{
		{"zzzz.LogIt", true, "LogIt"},
		{"001.LogIt", true, "LogIt"},
		{"A1.LogIt", true, "LogIt"},
		{"LogIt", false, "LogIt"},
	} {
		r := HasOrderPrefix(c.in)
		if r != c.has {
			t.Errorf("%5d. HasOrderPrefix(%q) should return %v, but got %v", i, c.in, c.has, r)
		}
		ret := StripOrderPrefix(c.in)
		if ret != c.out {
			t.Errorf("%5d. StripOrderPrefix(%q) should return %q, but got %q", i, c.in, c.out, ret)
		}
	}
}

func TestNextIn(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(NextIn(8))
	}
}

func TestNextString(t *testing.T) {
	// var ret string

	// for ix, c := range []testCase{
	// 	{"firstLine(874)", "FIRST_LINE(874)"},
	// 	{"secondLine(15),", "SECOND_LINE(15),"},
	// 	{"thirdLineOfText87(0x0001);", "THIRD_LINE_OF_TEXT87(0x0001);"},
	// 	// {"", ""},
	// 	// {"", ""},
	// 	// {"", ""},
	// } {
	// 	ret := ToUpperSnakeCase(c.in)
	// 	if ret != c.expect {
	// 		t.Errorf("%5d. Failed for KebabCase2SnakeCase(%q), want %q but got %q", ix, c.in, c.expect, ret)
	// 	}
	// }

	// for ix, c := range []testCase{
	// 	{"firstLine(874)", "first_line(874)"},
	// 	{"secondLine(15),", "second_line(15),"},
	// 	{"thirdLineOfText87(0x0001);", "third_line_of_text87(0x0001);"},
	// 	// {"", ""},
	// 	// {"", ""},
	// 	// {"", ""},
	// } {
	// 	ret := ToSnakeCase(c.in)
	// 	if ret != c.expect {
	// 		t.Errorf("%5d. Failed for KebabCase2SnakeCase(%q), want %q but got %q", ix, c.in, c.expect, ret)
	// 	}
	// }

	t.Log(NextString(8))
}
