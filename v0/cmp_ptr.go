package compare

import "reflect"

func (c *Comparer) cmpPtr(path []string, left, right reflect.Value, parent any) error {
	if left.Kind() != right.Kind() {
		if left.Kind() == reflect.Invalid {
			if !right.IsNil() {
				return c.compare(path, reflect.ValueOf(nil), reflect.Indirect(right), parent)
			}

			c.changes.add(ADD, path, nil, getAsAny(right), parent)
			return nil
		}

		if right.Kind() == reflect.Invalid {
			if !left.IsNil() {
				return c.compare(path, reflect.Indirect(left), reflect.ValueOf(nil), parent)
			}

			c.changes.add(REMOVE, path, getAsAny(left), nil, parent)
			return nil
		}

		return ErrTypeMismatch
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

	return c.compare(path, reflect.Indirect(left), reflect.Indirect(right), parent)
}
