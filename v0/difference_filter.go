package compare

func WherePath(p string) DiffFilterFunc {
	return func(d Difference) bool {
		for _, dp := range d.Path {
			if dp == p {
				return true
			}
		}
		return false
	}
}

func WherePathAt(p string, idx int) DiffFilterFunc {
	return func(d Difference) bool {
		if len(d.Path) >= idx && d.Path[idx] == p {
			return true
		}
		return false
	}
}

func WhereDiffType(dt DiffType) DiffFilterFunc {
	return func(d Difference) bool {
		if d.Type == dt {
			return true
		}
		return false
	}
}

func WherePathDepth(i int) DiffFilterFunc {
	return func(d Difference) bool {
		if i == len(d.Path) {
			return true
		}
		return false
	}
}
