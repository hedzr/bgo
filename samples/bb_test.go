package samples

//nolint:goimports
import (
	errors2 "errors"
	"gopkg.in/hedzr/errors.v3"
	"strconv"
	"testing"
)

// TestErrorsIs _
func TestErrorsIs(t *testing.T) {
	_, err := strconv.ParseFloat("hello", 64)
	t.Logf("err = %+v", err)
	e1 := errors2.Unwrap(err)
	t.Logf("e1 = %+v", e1)

	t.Logf("errors.Is(err, strconv.ErrSyntax): %v", errors.Is(err, strconv.ErrSyntax))
	t.Logf("errors.Is(err, &strconv.NumError{}): %v", errors.Is(err, &strconv.NumError{}))

	var e2 *strconv.NumError
	if errors.As(err, &e2) {
		t.Logf("As() ok, e2 = %v", e2)
	} else {
		t.Logf("As() not ok")
	}
}
