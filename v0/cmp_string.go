package compare

import "reflect"

func (c *Comparer) cmpString(path []string, left, right reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, left, right); err != nil || changed {
		return err
	}

	if left.String() != right.String() {
		if left.CanInterface() {
			// If left and/or right is of left type that is an alias for String, store that type in changelog
			c.changes.add(CHANGE, path, getAsAny(left), getAsAny(right), parent)
		} else {
			c.changes.add(CHANGE, path, left.String(), right.String(), parent)
		}
	}

	return nil
}
