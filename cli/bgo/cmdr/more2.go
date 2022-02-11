//go:build !delve && !test
// +build !delve,!test

package cmdr

import (
	"github.com/hedzr/cmdr"
)

func cmdrMoreCommandsForTest(root cmdr.OptCmd) {
	// normal build, no more subcommands for testing purpose
}
