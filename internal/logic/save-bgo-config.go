package logic

//nolint:goimports // i like it
import (
	"os"
	"path"

	"github.com/hedzr/cmdr"

	"github.com/hedzr/bgo/internal/logx"
	"github.com/hedzr/evendeep"
	"github.com/hedzr/evendeep/dbglog"

	"github.com/hedzr/log/dir"
)

func saveNewBgoYamlFile(bs *BgoSettings) (err error) {
	makeCopyOf := func(bs *BgoSettings) *BgoSettings {
		var bsCopy = new(BgoSettings)

		if err = evendeep.DefaultCopyController.CopyTo(bs, bsCopy); err != nil {
			logx.Error("error in makeCopyOf(bs): %v", err)
			return nil
		}
		// if err = cmdr.CloneViaGob(bsCopy, bs); err != nil {
		// 	return nil
		// }

		cleanupBs(bsCopy)
		return bsCopy
	}

	saveCmdrOptions := func() func() {
		logx.Log(`saving cmdr checkpoint`)
		defer dbglog.DisableLog()()
		err = cmdr.SaveCheckpoint()
		if err != nil {
			logx.Error("CANNOT save cmdr options checkpoint. err: %+v", err)
		}
		logx.Log(`reset cmdr options, prepare a clean, new options store`)
		cmdr.ResetOptions()
		return func() {
			logx.Log(`restoring cmdr checkpoint`)
			err = cmdr.RestoreCheckpoint()
		}
	}

	defer saveCmdrOptions()()

	logx.Log(`make a copy of BgoSettings as bsCopy | because saveBgoConfigAs will modify it`)
	bsCopy := makeCopyOf(bs)
	logx.Log(`and saving bsCopy as yaml file. bs = %v`, bs)
	return saveBgoConfigAs(bsCopy, bs.SavedAs)
}

func saveBgoConfigAs(bs *BgoSettings, savedAs string) (err error) {
	bs.SavedAs = ""
	err = cmdr.MergeWith(map[string]interface{}{
		"app": map[string]interface{}{
			"bgo": map[string]interface{}{
				"build": bs,
			},
		},
	})
	if err != nil {
		logx.Fatal("Error: %v", err)
	}
	// cmdr.DebugOutputTildeInfo(false)

	fn := savedAs
	if fn != "" {
		switch ext := path.Ext(fn); ext {
		case ".toml":
			err = cmdr.SaveAsToml(fn)
		case ".json":
			err = cmdr.SaveAsJSONExt(fn, true)
		// case ".yml", ".yaml":
		// 	fallthrough
		default:
			err = cmdr.SaveAsYaml(fn)
		}

		if err == nil {
			logx.Log("%q saved\n", path.Join(dir.GetCurrentDir(), fn))
			err = appendComments(fn)
		}
	}

	return
}

// appendComments todo append some sample yaml comments into .bgo.yml when `bgo init`
func appendComments(file string) (err error) {
	var f *os.File
	f, err = os.OpenFile(file, os.O_APPEND|os.O_RDWR, 0o644)
	if err == nil {
		defer func() {
			err = f.Close()
		}()

		_, err = f.WriteString(`

`)
	}
	return
}
