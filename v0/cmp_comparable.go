package compare

import "reflect"

func (c *Comparer) processComparableList(path []string, cmpList *ComparableList, parent any) error {
	for _, k := range cmpList.keys {
		id := idString(k)
		if c.structMapKeys {
			id = idComplex(k)
		}

		nv := reflect.ValueOf(nil)

		if cmpList.m[k].LEFT == nil {
			cmpList.m[k].LEFT = &nv
		}

		if cmpList.m[k].RIGHT == nil {
			cmpList.m[k].RIGHT = &nv
		}

		fpath := copyAppend(path, id)

		err := c.compare(fpath, *cmpList.m[k].LEFT, *cmpList.m[k].RIGHT, parent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Comparer) isComparable(left, right reflect.Value) bool {
	if left.Len() > 0 {
		leftElem := left.Index(0)
		leftVal := getFinalValue(leftElem)

		if leftVal.Kind() == reflect.Struct {
			if hasIdentifier(c.tagName, leftVal) != nil {
				return true
			}
		}
	}

	if right.Len() > 0 {
		rightElem := right.Index(0)
		rightVal := getFinalValue(rightElem)

		if rightVal.Kind() == reflect.Struct {
			if hasIdentifier(c.tagName, rightVal) != nil {
				return true
			}
		}
	}

	return false
}
