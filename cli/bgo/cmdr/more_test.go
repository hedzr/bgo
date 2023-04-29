package cmdr

import (
	"testing"

	"github.com/hedzr/log/detects"
)

func TestAny(t *testing.T) {
	detects.InDockerEnvSimple()
}
