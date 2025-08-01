package compare

import "reflect"

func (c *Comparer) cmpBool(path []string, left, right reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, left, right); err != nil || changed {
		return err
	}

	if left.Bool() != right.Bool() {
		c.changes.add(CHANGE, path, getAsAny(left), getAsAny(right), parent)
	}

	return nil
}
