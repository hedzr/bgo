package logic

import ( //nolint:goimports //i like it
	"github.com/hedzr/bgo/internal/logic/build"
	"github.com/hedzr/bgo/internal/logic/logx"
	//nolint:goimports //i like it
)

func listPackages(tpBase *build.TargetPlatforms, bc *build.Context, bs *BgoSettings, packages map[string]*pkgInfo) (err error) { // nolint
	// ensureProject(pi, bc, bs)
	//
	// logx.Colored(logx.Green, "> building for package %v (dir: %q)...", pi.p.Package, pi.dirname)

	err = loopAllProjects(tpBase, bc, bs, listProject)
	return
}

func listProject(bc *build.Context, bs *BgoSettings) (err error) {
	// logx.Log("  >> %v/%v", bc.Os, bc.Arch)
	logx.Log("      >> %v/%v, Need Install: %v\n", bc.OS, bc.ARCH, bc.Install)
	// logx.Dim("     project.Common: %+v\n", *p.Common)

	return
}
