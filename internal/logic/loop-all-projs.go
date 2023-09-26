package logic

import (
	"runtime"

	"github.com/hedzr/log/dir"

	"github.com/hedzr/bgo/internal/logic/build"
	"github.com/hedzr/bgo/internal/logx"
	"github.com/hedzr/bgo/internal/tool"
)

func loopTargetPlatforms(
	tp *build.TargetPlatforms,
	cb func(os, arch string) (stop bool, err error),
) (err error) {
	var stop bool
	for oss, osv := range tp.OsArchMap {
		for arch := range osv {
			if stop, err = cb(oss, arch); stop {
				return
			}
		}
	}
	return
}

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
		//nolint:lll //no
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

			var terminateNow bool
			if terminateNow, err = loopChildren(bc, bs, cb, it.tp.OsArchMap, shortMode, pn, gn, grp, it); terminateNow {
				break RETURN
			}
			// STOP:
			// 	for osName, osv := range it.tp.OsArchMap {
			// 		for archName := range osv {
			// 			if shortMode {
			// 				if osName != runtime.GOOS || archName != runtime.GOARCH {
			// 					continue
			// 				}
			// 			}
			//
			// 			prepareBuildContextForEachProjectTarget(
			// 				bc, osName, archName,
			// 				it, pn, gn, grp.LeadingText,
			// 			)
			// 			if err = cb(bc, bs); err != nil {
			// 				break STOP
			// 			}
			//
			// 			if shortMode {
			// 				break RETURN
			// 			}
			// 		}
			// 	}
		}
	}

	return
}

func loopChildren(
	bc *build.Context, bs *BgoSettings,
	cb func(bc *build.Context, bs *BgoSettings) (err error),
	osArchMap map[string]map[string]bool, shortMode bool,
	pn, gn string, grp ProjectGroup, it *ProjectWrap,
) (terminateNow bool, err error) {
	for osName, osv := range osArchMap {
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
				return
			}

			if shortMode {
				terminateNow = true
				return
			}
		}
	}
	return
}
