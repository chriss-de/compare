package compare

import "reflect"

func (c *Comparer) cmpSlice(path []string, a, b reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if c.isComparable(a, b) {
		return c.cmpSliceComparable(path, a, b)
	}

	return c.cmpSliceGeneric(path, a, b)
}

func (c *Comparer) cmpSliceGeneric(path []string, a, b reflect.Value) error {
	missing := NewComparableList()

	sliceA := sliceTracker{}
	for i := 0; i < a.Len(); i++ {
		ae := a.Index(i)

		if (c.sliceOrdering && !hasAtSameIndex(b, ae, i)) || (!c.sliceOrdering && !sliceA.has(b, ae, c)) {
			missing.addA(i, &ae)
		}
	}

	sliceB := sliceTracker{}
	for i := 0; i < b.Len(); i++ {
		be := b.Index(i)

		if (c.sliceOrdering && !hasAtSameIndex(a, be, i)) || (!c.sliceOrdering && !sliceB.has(a, be, c)) {
			missing.addB(i, &be)
		}
	}

	// fallback to comparing based on order in slice if item is missing
	if len(missing.keys) == 0 {
		return nil
	}

	return c.processComparableList(path, missing, getAsAny(a))
}

func (c *Comparer) cmpSliceComparable(path []string, a, b reflect.Value) error {
	c := NewComparableList()

	for i := 0; i < a.Len(); i++ {
		aElem := a.Index(i)
		aVal := getFinalValue(aElem)

		id := hasIdentifier(c.tagName, aVal)
		if id != nil {
			c.addA(id, &aElem)
		}
	}

	for i := 0; i < b.Len(); i++ {
		bElem := b.Index(i)
		bVal := getFinalValue(bElem)

		id := hasIdentifier(c.tagName, bVal)
		if id != nil {
			c.addB(id, &bElem)
		}
	}

	return c.processComparableList(path, c, getAsAny(a))
}
