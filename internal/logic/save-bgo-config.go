package logic

import (
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log/dir"
	"os"
	"path"
)

func saveNewBgoYamlFile(bs *BgoSettings) (err error) {
	makeCopyOf := func(bs *BgoSettings) *BgoSettings {
		var bsCopy = new(BgoSettings)

		if err = cmdr.CloneViaGob(bsCopy, bs); err != nil {
			return nil
		}

		cleanupBs(bsCopy)
		return bsCopy
	}

	_ = cmdr.SaveCheckpoint()
	defer func() { err = cmdr.RestoreCheckpoint() }()
	cmdr.ResetOptions()

	bsCopy := makeCopyOf(bs)
	return saveBgoConfigAs(bsCopy, bs.SavedAs)
}

func saveBgoConfigAs(bs *BgoSettings, savedAs []string) (err error) {

	bs.SavedAs = nil
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
	//cmdr.DebugOutputTildeInfo(false)

	for _, fn := range savedAs {
		if fn != "" {
			switch ext := path.Ext(fn); ext {
			case ".toml":
				err = cmdr.SaveAsToml(fn)
			case ".json":
				err = cmdr.SaveAsJSONExt(fn, true)
			case ".yml", ".yaml":
				fallthrough
			default:
				err = cmdr.SaveAsYaml(fn)
			}

			if err == nil {
				logx.Log("%q saved\n", path.Join(dir.GetCurrentDir(), fn))
				err = appendComments(fn)
			}
		}
	}
	return
}

// appendComments todo append some sample yaml comments into .bgo.yml when `bgo init`
func appendComments(file string) (err error) {
	var f *os.File
	f, err = os.OpenFile(file, os.O_APPEND|os.O_RDWR, 0644)
	if err == nil {
		defer func() {
			err = f.Close()
		}()

		_, err = f.WriteString(`

`)
	}
	return
}
