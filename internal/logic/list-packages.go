package logic

import (
	"github.com/hedzr/bgo/internal/logic/build"
	"github.com/hedzr/bgo/internal/logic/logx"
	"github.com/hedzr/bgo/internal/logic/tool"
	"github.com/hedzr/log/dir"
)

func listPackages(tpBase *build.TargetPlatforms, bc *build.Context, bs *BgoSettings, packages map[string]*pkgInfo) (err error) {
	// ensureProject(pi, bc, bs)
	//
	// logx.Colored(logx.Green, "> building for package %v (dir: %q)...", pi.p.Package, pi.dirname)

	err = loopAllProjects(tpBase, bc, bs, listProject)
	return
}

func listProject(bc *build.Context, bs *BgoSettings) (err error) {
	//logx.Log("  >> %v/%v", bc.Os, bc.Arch)
	logx.Log("      >> %v/%v, Need Install: %v\n", bc.OS, bc.ARCH, bc.Install)
	//logx.Dim("     project.Common: %+v\n", *p.Common)

	return
}

func loopAllProjects(
	tpBase *build.TargetPlatforms,
	bc *build.Context, bs *BgoSettings,
	cb func(bc *build.Context, bs *BgoSettings) (err error),
) (err error) {
	var seq int

	kiSlice := getSortedProjectGroupKeys(bs.Projects)
	for _, ki := range kiSlice {
		gn, grp := ki.grp, bs.Projects[ki.grp]
		if grp.Common != nil && grp.Common.Disabled {
			continue
		}
		logx.Log("-> Group %v: %v ...\n", logx.ToColor(logx.Yellow, tool.StripOrderPrefix(gn)), logx.ToColor(logx.Green, grp.LeadingText))

		kiSliceP := getSortedProjectKeys(gn, &grp)
		for _, kiP := range kiSliceP {
			//for pn, it := range grp.Items {
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
			for oss, osv := range it.tp.OsArchMap {
				for arch, _ := range osv {
					prepareBuildContextForEachProjectTarget(bc, oss, arch, it,
						pn, gn, grp.LeadingText)
					if err = cb(bc, bs); err != nil {
						break STOP
					}
				}
			}
		}
	}

	return
}
