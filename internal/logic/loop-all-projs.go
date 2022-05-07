package logic

import (
	"github.com/hedzr/bgo/internal/logic/build"
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/bgo/internal/logic/tool"
	"github.com/hedzr/log/dir"
	"runtime"
)

//nolint:gocognit //needs split
func loopAllProjects(
	tpBase *build.TargetPlatforms,
	bc *build.Context, bs *BgoSettings,
	cb func(bc *build.Context, bs *BgoSettings) (err error),
) (err error) {
	seq := 0
	shortMode := bs.Scope == "short"

	kiSlice := getSortedProjectGroupKeys(bs.Projects)
RETURN:
	for _, ki := range kiSlice {
		gn, grp := ki.grp, bs.Projects[ki.grp]
		if grp.Common != nil && grp.Common.Disabled {
			continue
		}
		logx.Log("-> Group %v: %v ...\n", logx.ToColor(logx.Yellow, tool.StripOrderPrefix(gn)), logx.ToColor(logx.Green, grp.LeadingText))

		kiSliceP := getSortedProjectKeys(gn, &grp)
		for _, kiP := range kiSliceP {
			// for pn, it := range grp.Items {
			pn, it := kiP.prj, grp.Items[kiP.Index]
			if it.Disabled {
				continue
			}

			seq++
			it.keyName = pn
			logx.Colored(logx.Green, "   %d. Project %v ...\n", seq, it.GetTitleName())

			if !dir.FileExists(it.Dir) {
				logx.Dim("      %v\n", "Ignored since it's not exists")
				continue
			}

			it.overspreadByTP(bs.Scope, tpBase)

		STOP:
			for osName, osv := range it.tp.OsArchMap {
				for archName := range osv {
					if shortMode {
						if osName != runtime.GOOS || archName != runtime.GOARCH {
							continue
						}
					}

					prepareBuildContextForEachProjectTarget(
						bc, osName, archName,
						it, pn, gn, grp.LeadingText,
					)
					if err = cb(bc, bs); err != nil {
						break STOP
					}

					if shortMode {
						break RETURN
					}
				}
			}
		}
	}

	return
}
