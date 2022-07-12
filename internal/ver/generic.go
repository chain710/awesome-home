package ver

import (
	"regexp"
	"sort"
	"strconv"
)

func NewGeneric(s string) (*GenericVersion, error) {
	var gv GenericVersion
	gv.raw = s
	r := regexp.MustCompile(`\d+`)
	segments := r.FindAllString(s, -1)
	for _, segment := range segments {
		if v, err := strconv.Atoi(segment); err != nil {
			return nil, err
		} else {
			gv.segments = append(gv.segments, v)
		}
	}
	return &gv, nil
}

type GenericVersion struct {
	raw      string
	segments []int
}

func (v *GenericVersion) Less(w *GenericVersion) bool {
	lenv := len(v.segments)
	lenw := len(w.segments)
	end := lenv
	if lenw < end {
		end = lenw
	}

	for i := 0; i < end; i++ {
		if v.segments[i] == w.segments[i] {
			continue
		} else {
			return v.segments[i] < w.segments[i]
		}
	}

	return lenv < lenw
}

func (v *GenericVersion) String() string {
	return v.raw
}

type GenericVersions []GenericVersion

func (g GenericVersions) Len() int {
	return len(g)
}

func (g GenericVersions) Less(i, j int) bool {
	return g[i].Less(&g[j])

}

func (g GenericVersions) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

// UpperBound return element that is greater than or equal to every element of versions, must call after sort
func (g GenericVersions) UpperBound(v *GenericVersion) *GenericVersion {
	n := len(g)
	if n == 0 {
		panic("try to find ver in empty slice")
	}
	find := sort.Search(n, func(i int) bool {
		return v.Less(&g[i])
	})
	if find == n {
		return &g[n-1]
	} else {
		return &g[find]
	}
}
