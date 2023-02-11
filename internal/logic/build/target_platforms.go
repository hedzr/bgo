package build

//nolint:goimports //so what
import (
	"path"
	"runtime"
	"strings"

	"github.com/hedzr/cmdr"

	"github.com/hedzr/bgo/internal/logic/logx"
)

type TargetPlatforms struct {
	// OsArchMap is a map by key OS and CPUArch
	// key1: os
	// key2: arch
	OsArchMap map[string]map[string]bool `yaml:"os-arch-map"`
	// Sources is a map by key PackageName
	// key: packageName
	Sources map[string]*DynBuildInfo `yaml:"sources"`
	// Locked identify only os,arch pairs in OsArchMap should
	// be applied for building.
	Locked bool
}

func NewTargetPlatforms() *TargetPlatforms {
	tp := &TargetPlatforms{
		OsArchMap: map[string]map[string]bool{},
		Sources:   make(map[string]*DynBuildInfo),
	}
	if cmdr.GetUsedAlterConfigFile() == "" {
		tp.SetOsArch(runtime.GOOS, runtime.GOARCH)
	}
	return tp
}

func (ss *TargetPlatforms) Init() (err error) {
	// Prepare()
	// for 'bgo.dists' key in cmdr Option Store
	prepareBuildInfo()

	// 'bgo.dists' was prepared by extracting from
	// 'go tool dist list', see also
	// prepareBuildInfo()
	err = cmdr.GetSectionFrom("bgo.dists", &ss)

	// cmdr.DebugOutputTildeInfo(true)

	return
}

func (ss *TargetPlatforms) setSource(pkgName string) { //nolint:unused //no
	ss.Sources[pkgName] = &DynBuildInfo{
		AppName:     path.Base(pkgName),
		Version:     "",
		ProjectName: "",
	}
}

func (ss *TargetPlatforms) SetOsArch(os, arch string) {
	if _, ok := ss.OsArchMap[os]; !ok {
		ss.OsArchMap[os] = make(map[string]bool)
	}
	ss.OsArchMap[os][arch] = true
}

func (ss *TargetPlatforms) isValidPair(os, arch string) (ok bool) {
	if ss.OsArchMap != nil {
		var m map[string]bool
		if m, ok = ss.OsArchMap[os]; ok {
			_, ok = m[arch]
		}
	}
	return
}

func (ss *TargetPlatforms) filterByFor(oamNew *TargetPlatforms, forSlice []string) {
	for _, s := range forSlice {
		a := strings.Split(s, "/")
		aos, aarch := a[0], a[1]
		// if scope == "full" {
		//	for oss, osv := range ss.OsArchMap {
		//		for arch, _ := range osv {
		//			if oss == aos && arch == aarch {
		//				oamNew.SetOsArch(aos, aarch)
		//			}
		//		}
		//	}
		// } else {
		//	oamNew.SetOsArch(aos, aarch)
		// }
		if ss.isValidPair(aos, aarch) {
			oamNew.SetOsArch(aos, aarch)
		}
	}
}

func (ss *TargetPlatforms) filterByOs(oamNew *TargetPlatforms, osSlice []string) {
	for _, aos := range osSlice {
		if m, ok := ss.OsArchMap[aos]; ok {
			for aarch := range m {
				oamNew.SetOsArch(aos, aarch)
			}
		}
		// for oss, osv := range ss.OsArchMap {
		//	for arch, _ := range osv {
		//		if oss == aos {
		//			oamNew.SetOsArch(aos, arch)
		//		}
		//	}
		// }
	}
}

func (ss *TargetPlatforms) filterByArch(oamNew *TargetPlatforms, archSlice []string) {
	for _, aarch := range archSlice {
		for oss, osv := range ss.OsArchMap {
			if _, ok := osv[aarch]; ok {
				oamNew.SetOsArch(oss, aarch)
			}
			// for arch, _ := range osv {
			//	if arch == aarch {
			//		oamNew.SetOsArch(oss, aarch)
			//	}
			// }
		}
	}
}

func (ss *TargetPlatforms) filterByOsArchBoth(oamNew *TargetPlatforms, osSlice, archSlice []string) {
	for _, aos := range osSlice {
		for _, aarch := range archSlice {
			// // ?? why
			// if scope == "full" {
			//	for oss, osv := range ss.OsArchMap {
			//		for arch, _ := range osv {
			//			if oss == aos && arch == aarch {
			//				oamNew.SetOsArch(aos, aarch)
			//			}
			//		}
			//	}
			// } else {
			//	oamNew.SetOsArch(aos, aarch)
			// }
			if ss.isValidPair(aos, aarch) {
				oamNew.SetOsArch(aos, aarch)
			}
		}
	}
}

func (ss *TargetPlatforms) FilterBy(scope string, forSlice, osSlice, archSlice []string) {
	if scope == "short" {
		return
	}

	oamNew := NewTargetPlatforms()

	ss.filterByFor(oamNew, forSlice)

	switch {
	case len(osSlice) > 0 && len(archSlice) == 0:
		ss.filterByOs(oamNew, osSlice)
	case len(osSlice) == 0 && len(archSlice) > 0:
		ss.filterByArch(oamNew, archSlice)
	default:
		ss.filterByOsArchBoth(oamNew, osSlice, archSlice)
	}

	ss.OsArchMap = oamNew.OsArchMap
	logx.DimV("After FilterBy(), the TargetPlatforms are: %v", ss.OsArchMap)
}
