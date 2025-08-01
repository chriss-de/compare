package compare

import "reflect"

func (c *Comparer) cmpInt(path []string, left, right reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, left, right); err != nil || changed {
		return err
	}

	if left.Int() != right.Int() {
		if left.CanInterface() {
			c.changes.add(CHANGE, path, getAsAny(left), getAsAny(right), parent)
		} else {
			c.changes.add(CHANGE, path, left.Int(), right.Int(), parent)
		}
	}

	return nil
}
