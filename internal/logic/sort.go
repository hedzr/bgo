package logic

import "sort"

type KeyIndex struct {
	// index of map[string]pkgInfo
	Index    string
	grp, prj string
}

func getSortedV(m map[string]*pkgInfo) (ki []KeyIndex) {
	for k, v := range m {
		ki = append(ki, KeyIndex{k, v.groupKey, v.projectName})
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
