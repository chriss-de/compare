package compare

import "reflect"

func (c *Comparer) cmpPtr(path []string, a, b reflect.Value, parent any) error {
	if a.Kind() != b.Kind() {
		if a.Kind() == reflect.Invalid {
			if !b.IsNil() {
				return c.compare(path, reflect.ValueOf(nil), reflect.Indirect(b), parent)
			}

			c.changes.add(ADD, path, nil, getAsAny(b), parent)
			return nil
		}

		if b.Kind() == reflect.Invalid {
			if !a.IsNil() {
				return c.compare(path, reflect.Indirect(a), reflect.ValueOf(nil), parent)
			}

			c.changes.add(REMOVE, path, getAsAny(a), nil, parent)
			return nil
		}

		return ErrTypeMismatch
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

	return c.compare(path, reflect.Indirect(a), reflect.Indirect(b), parent)
}
