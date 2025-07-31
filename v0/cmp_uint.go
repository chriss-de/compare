package compare

import "reflect"

func (c *Comparer) cmpUint(path []string, a, b reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if a.Uint() != b.Uint() {
		if a.CanInterface() {
			c.changes.add(CHANGE, path, getAsAny(a), getAsAny(b), parent)
		} else {
			c.changes.add(CHANGE, path, a.Uint(), b.Uint(), parent)
		}
	}

	return nil
}
