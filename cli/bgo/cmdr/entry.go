package cmdr

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log"
	"gopkg.in/hedzr/errors.v3"
)

func Entry() {
	root := buildRootCmd()
	if err := cmdr.Exec(root, options...); err != nil {
		log.Fatalf("error occurs in app running: %+v\n", err)
	}
}

func onUnhandledErrorHandler(err interface{}) {
	if cmdr.GetBoolR("enable-ueh") {
		dumpStacks()
		// return
	}

	// internal.App().Close()

	panic(err)
}

func dumpStacks() {
	//nolint:lll //no
	log.Printf("\n\n=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===\n\n", errors.DumpStacksAsString(true))
}
