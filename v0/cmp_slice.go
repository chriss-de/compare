package compare

import "reflect"

func (d *Comparer) cmpSlice(path []string, a, b reflect.Value, parent any) error {
	if changed, err := d.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if d.isComparable(a, b) {
		return d.cmpSliceComparable(path, a, b)
	}

	return d.cmpSliceGeneric(path, a, b)
}

func (d *Comparer) cmpSliceGeneric(path []string, a, b reflect.Value) error {
	missing := NewComparableList()

	sliceA := sliceTracker{}
	for i := 0; i < a.Len(); i++ {
		ae := a.Index(i)

		if (d.sliceOrdering && !hasAtSameIndex(b, ae, i)) || (!d.sliceOrdering && !sliceA.has(b, ae, d)) {
			missing.addA(i, &ae)
		}
	}

	sliceB := sliceTracker{}
	for i := 0; i < b.Len(); i++ {
		be := b.Index(i)

		if (d.sliceOrdering && !hasAtSameIndex(a, be, i)) || (!d.sliceOrdering && !sliceB.has(a, be, d)) {
			missing.addB(i, &be)
		}
	}

	// fallback to comparing based on order in slice if item is missing
	if len(missing.keys) == 0 {
		return nil
	}

	return d.processComparableList(path, missing, getAsAny(a))
}

func (d *Comparer) cmpSliceComparable(path []string, a, b reflect.Value) error {
	c := NewComparableList()

	for i := 0; i < a.Len(); i++ {
		aElem := a.Index(i)
		aVal := getFinalValue(aElem)

		id := hasIdentifier(d.tagName, aVal)
		if id != nil {
			c.addA(id, &aElem)
		}
	}

	for i := 0; i < b.Len(); i++ {
		bElem := b.Index(i)
		bVal := getFinalValue(bElem)

		id := hasIdentifier(d.tagName, bVal)
		if id != nil {
			c.addB(id, &bElem)
		}
	}

	return d.processComparableList(path, c, getAsAny(a))
}
