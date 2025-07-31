package compare

import "reflect"

func (c *Comparer) cmpBool(path []string, a, b reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if a.Bool() != b.Bool() {
		c.changes.add(CHANGE, path, getAsAny(a), getAsAny(b), parent)
	}

	return nil
}
