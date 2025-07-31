package compare

import "reflect"

func (d *Comparer) processComparableList(path []string, c *ComparableList, parent any) error {
	for _, k := range c.keys {
		id := idString(k)
		if d.structMapKeys {
			id = idComplex(k)
		}

		nv := reflect.ValueOf(nil)

		if c.m[k].A == nil {
			c.m[k].A = &nv
		}

		if c.m[k].B == nil {
			c.m[k].B = &nv
		}

		fpath := copyAppend(path, id)

		err := d.compare(fpath, *c.m[k].A, *c.m[k].B, parent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Comparer) isComparable(a, b reflect.Value) bool {
	if a.Len() > 0 {
		aElem := a.Index(0)
		aVal := getFinalValue(aElem)

		if aVal.Kind() == reflect.Struct {
			if hasIdentifier(d.tagName, aVal) != nil {
				return true
			}
		}
	}

	if b.Len() > 0 {
		bElem := b.Index(0)
		bVal := getFinalValue(bElem)

		if bVal.Kind() == reflect.Struct {
			if hasIdentifier(d.tagName, bVal) != nil {
				return true
			}
		}
	}

	return false
}
