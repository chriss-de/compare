package compare

import "reflect"

func (c *Comparer) cmpFloat(path []string, left, right reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, left, right); err != nil || changed {
		return err
	}

	if left.Float() != right.Float() {
		if left.CanInterface() {
			c.changes.add(CHANGE, path, getAsAny(left), getAsAny(right), parent)
		} else {
			c.changes.add(CHANGE, path, left.Float(), right.Float(), parent)
		}
	}

	return nil
}
