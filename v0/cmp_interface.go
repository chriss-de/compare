package compare

import "reflect"

func (c *Comparer) cmpInterface(path []string, left, right reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, left, right); err != nil || changed {
		return err
	}

	if left.IsNil() && right.IsNil() {
		return nil
	}

	if left.IsNil() {
		c.changes.add(CHANGE, path, nil, getAsAny(right), parent)
		return nil
	}

	if right.IsNil() {
		c.changes.add(CHANGE, path, getAsAny(left), nil, parent)
		return nil
	}

	return c.compare(path, left.Elem(), right.Elem(), parent)
}
