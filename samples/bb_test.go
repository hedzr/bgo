package samples

//nolint:goimports //so what
import (
	errors2 "errors"
	"strconv"
	"testing"

	"gopkg.in/hedzr/errors.v3"
)

// TestErrorsIs _
func TestErrorsIs(t *testing.T) {
	_, err := strconv.ParseFloat("hello", 64)

	t.Logf("err = %+v", err)

	e1 := errors2.Unwrap(err)
	t.Logf("e1 = %+v", e1)

	t.Logf("err = %+v", err)
	t.Logf("errors.Is(err, strconv.ErrSyntax): %v", errors.Is(err, strconv.ErrSyntax))
	t.Logf("err = %+v", err)
	t.Logf("errors.Is(err, &strconv.NumError{Err: strconv.ErrRange}): %v", errors.Is(err, &strconv.NumError{Err: strconv.ErrRange}))
	t.Logf("err = %+v", err)

	var e2 *strconv.NumError
	if errors.As(err, &e2) {
		t.Logf("As() ok, e2 = %v", e2)
	} else {
		t.Logf("As() not ok")
	}
}
