package compare

import "slices"

func WhereOr(filterFunc ...DiffFilterFunc) DiffFilterFunc {
	return func(d Difference) bool {
		for _, ff := range filterFunc {
			if ff(d) {
				return true
			}
		}
		return false
	}
}

func WherePath(p string) DiffFilterFunc {
	return func(d Difference) bool {
		return slices.Contains(d.Path, p)
	}
}

func WherePathAt(p string, idx int) DiffFilterFunc {
	return func(d Difference) bool {
		return len(d.Path) > idx && d.Path[idx] == p
	}
}

func WherePathDepth(l int) DiffFilterFunc {
	return func(d Difference) bool {
		return len(d.Path) == l
	}
}

func WherePathDepthGt(l int) DiffFilterFunc {
	return func(d Difference) bool {
		return len(d.Path) > l
	}
}

func WherePathDepthLt(l int) DiffFilterFunc {
	return func(d Difference) bool {
		return len(d.Path) < l
	}
}

func WhereDiffType(dt DiffType) DiffFilterFunc {
	return func(d Difference) bool {
		return d.Type == dt
	}
}
