package compare

import "reflect"

func (c *Comparer) cmpSlice(path []string, left, right reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, left, right); err != nil || changed {
		return err
	}

	if c.isComparable(left, right) {
		return c.cmpSliceComparable(path, left, right)
	}

	return c.cmpSliceGeneric(path, left, right)
}

func (c *Comparer) cmpSliceGeneric(path []string, left, right reflect.Value) error {
	missing := NewComparableList()

	sliceLeft := sliceTracker{}
	for i := 0; i < left.Len(); i++ {
		leftElem := left.Index(i)

		if (c.config.sliceOrdering && !hasAtSameIndex(right, leftElem, i)) || (!c.config.sliceOrdering && !sliceLeft.has(right, leftElem, c)) {
			missing.addLeft(i, &leftElem)
		}
	}

	sliceRight := sliceTracker{}
	for i := 0; i < right.Len(); i++ {
		rightElem := right.Index(i)

		if (c.config.sliceOrdering && !hasAtSameIndex(left, rightElem, i)) || (!c.config.sliceOrdering && !sliceRight.has(left, rightElem, c)) {
			missing.addRight(i, &rightElem)
		}
	}

	// fallback to comparing based on order in slice if item is missing
	if len(missing.keys) == 0 {
		return nil
	}

	return c.processComparableList(path, missing, getAsAny(left))
}

func (c *Comparer) cmpSliceComparable(path []string, left, right reflect.Value) error {
	cmpList := NewComparableList()

	for i := 0; i < left.Len(); i++ {
		leftElem := left.Index(i)
		leftVal := getFinalValue(leftElem)

		leftID := getIdentifier(c.config.tagName, leftVal, string(c.config.combinedIdentifierJoinSep))
		if leftID != nil {
			cmpList.addLeft(leftID, &leftElem)
		}
	}

	for i := 0; i < right.Len(); i++ {
		rightElem := right.Index(i)
		rightVal := getFinalValue(rightElem)

		leftID := getIdentifier(c.config.tagName, rightVal, string(c.config.combinedIdentifierJoinSep))
		if leftID != nil {
			cmpList.addRight(leftID, &rightElem)
		}
	}

	return c.processComparableList(path, cmpList, getAsAny(left))
}
