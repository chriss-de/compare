package compare

import "reflect"

func (c *Comparer) cmpFloat(path []string, a, b reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if a.Float() != b.Float() {
		if a.CanInterface() {
			c.changes.add(CHANGE, path, getAsAny(a), getAsAny(b), parent)
		} else {
			c.changes.add(CHANGE, path, a.Float(), b.Float(), parent)
		}
	}

	return nil
}
