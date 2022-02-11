package logic

import "sort"

type keyIndex struct {
	// index of map[string]pkgInfo
	Index    string
	grp, prj string
}

func getSortedV(m map[string]*pkgInfo) (ki []keyIndex) {
	for k, v := range m {
		ki = append(ki, keyIndex{k, v.groupKey, v.projectName})
	}

	sort.Slice(ki, func(i, j int) bool {
		if ki[i].grp < ki[j].grp {
			return true
		}
		if ki[i].grp > ki[j].grp {
			return false
		}
		return ki[i].prj < ki[j].prj
	})
	return
}

func getSortedProjectGroupKeys(m map[string]ProjectGroup) (ki []keyIndex) {
	for k := range m {
		ki = append(ki, keyIndex{k, k, ""})
	}

	sort.Slice(ki, func(i, j int) bool {
		if ki[i].grp < ki[j].grp {
			return true
		}
		if ki[i].grp > ki[j].grp {
			return false
		}
		return ki[i].prj < ki[j].prj
	})
	return
}

func getSortedProjectKeys(gn string, g *ProjectGroup) (ki []keyIndex) {
	for k := range g.Items {
		ki = append(ki, keyIndex{k, gn, k})
	}

	sort.Slice(ki, func(i, j int) bool {
		if ki[i].grp < ki[j].grp {
			return true
		}
		if ki[i].grp > ki[j].grp {
			return false
		}
		return ki[i].prj < ki[j].prj
	})
	return
}
