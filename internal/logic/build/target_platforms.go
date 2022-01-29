package build

import (
	"github.com/hedzr/cmdr"
	"path"
	"runtime"
	"strings"
)

type TargetPlatforms struct {
	// OsArchMap is a map by key OS and CPUArch
	OsArchMap map[string]map[string]bool `yaml:"os-arch-map"`
	// Sources is a map by key PackageName
	Sources map[string]*DynBuildInfo `yaml:"sources"`
}

func NewTargetPlatforms() *TargetPlatforms {
	tp := &TargetPlatforms{
		OsArchMap: map[string]map[string]bool{},
		Sources:   make(map[string]*DynBuildInfo),
	}
	if cmdr.GetUsedConfigFile() == "" {
		tp.SetOsArch(runtime.GOOS, runtime.GOARCH)
	}
	return tp
}

func (ss *TargetPlatforms) setSource(pkgName string) {
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

func (ss *TargetPlatforms) FilterBy(scope string, forSlice, osSlice, archSlice []string) {
	if scope == "short" {
		return
	}

	oamNew := NewTargetPlatforms()

	for _, s := range forSlice {
		a := strings.Split(s, "/")
		aos, aarch := a[0], a[1]
		if scope == "full" {
			for oss, osv := range ss.OsArchMap {
				for arch, _ := range osv {
					if oss == aos && arch == aarch {
						oamNew.SetOsArch(aos, aarch)
					}
				}
			}
		} else {
			oamNew.SetOsArch(aos, aarch)
		}
	}

	if len(osSlice) > 0 && len(archSlice) == 0 {
		for _, aos := range osSlice {
			for oss, osv := range ss.OsArchMap {
				for arch, _ := range osv {
					if oss == aos {
						oamNew.SetOsArch(aos, arch)
					}
				}
			}
		}

	} else if len(osSlice) == 0 && len(archSlice) > 0 {
		for _, aarch := range archSlice {
			for oss, osv := range ss.OsArchMap {
				for arch, _ := range osv {
					if arch == aarch {
						oamNew.SetOsArch(oss, aarch)
					}
				}
			}
		}

	} else {
		for _, aos := range osSlice {
			for _, aarch := range archSlice {
				if scope == "full" {
					for oss, osv := range ss.OsArchMap {
						for arch, _ := range osv {
							if oss == aos && arch == aarch {
								oamNew.SetOsArch(aos, aarch)
							}
						}
					}
				} else {
					oamNew.SetOsArch(aos, aarch)
				}
			}
		}
	}
	ss.OsArchMap = oamNew.OsArchMap
}
