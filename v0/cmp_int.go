package compare

import "reflect"

func (c *Comparer) cmpInt(path []string, a, b reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if a.Int() != b.Int() {
		if a.CanInterface() {
			c.changes.add(CHANGE, path, getAsAny(a), getAsAny(b), parent)
		} else {
			c.changes.add(CHANGE, path, a.Int(), b.Int(), parent)
		}
	}

	return nil
}
