package compare

import "reflect"

func (c *Comparer) cmpUint(path []string, left, right reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, left, right); err != nil || changed {
		return err
	}

	if left.Uint() != right.Uint() {
		if left.CanInterface() {
			c.changes.add(CHANGE, path, getAsAny(left), getAsAny(right), parent)
		} else {
			c.changes.add(CHANGE, path, left.Uint(), right.Uint(), parent)
		}
	}

	return nil
}
