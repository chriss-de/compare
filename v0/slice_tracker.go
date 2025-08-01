package compare

import "reflect"

// keeps track of elements that have already been matched, to stop duplicate matches from occurring
type sliceTracker []bool

func (st *sliceTracker) has(s, v reflect.Value, c *Comparer) bool {
	if len(*st) != s.Len() {
		*st = make([]bool, s.Len())
	}

	for i := 0; i < s.Len(); i++ {
		// skip already matched elements
		if (*st)[i] {
			continue
		}

		x := s.Index(i)

		var nc Comparer

		err := nc.compare([]string{}, x, v, nil)
		if err != nil {
			continue
		}

		if len(nc.changes) == 0 {
			(*st)[i] = true
			return true
		}
	}

	return false
}
