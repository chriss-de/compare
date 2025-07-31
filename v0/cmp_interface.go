package compare

import "reflect"

func (c *Comparer) cmpInterface(path []string, a, b reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if a.IsNil() && b.IsNil() {
		return nil
	}

	if a.IsNil() {
		c.changes.add(CHANGE, path, nil, getAsAny(b), parent)
		return nil
	}

	if b.IsNil() {
		c.changes.add(CHANGE, path, getAsAny(a), nil, parent)
		return nil
	}

	return c.compare(path, a.Elem(), b.Elem(), parent)
}
