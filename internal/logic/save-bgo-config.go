package logic

import (
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log/dir"
	"os"
	"path"
)

func saveNewBgoYamlFile(bs *BgoSettings) (err error) {
	_ = cmdr.SaveCheckpoint()
	defer cmdr.RestoreCheckpoint()

	cmdr.ResetOptions()

	var bsCopy = new(BgoSettings)

	if err = cmdr.CloneViaGob(bsCopy, bs); err != nil {
		return
	}

	cleanupBs(bsCopy)

	err = cmdr.MergeWith(map[string]interface{}{
		"app": map[string]interface{}{
			"bgo": map[string]interface{}{
				"build": bsCopy,
			},
		},
	})

	//cmdr.DebugOutputTildeInfo(false)
	if err != nil {
		logx.Fatal("Error: %v", err)
	}

	return saveBgoConfigAs(bs)
}

func saveBgoConfigAs(bs *BgoSettings) (err error) {
	logx.Log("%q saved\n", path.Join(dir.GetCurrentDir(), bs.SavedAs))

	switch ext := path.Ext(bs.SavedAs); ext {
	case ".yml", ".yaml":
		err = cmdr.SaveAsYaml(bs.SavedAs)
	case ".toml":
		err = cmdr.SaveAsToml(bs.SavedAs)
	case ".json":
		err = cmdr.SaveAsJSON(bs.SavedAs)
	}

	if err == nil {
		err = appendComments(bs.SavedAs)
	}
	return
}

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
