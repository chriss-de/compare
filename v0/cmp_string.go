package compare

import "reflect"

func (d *Comparer) cmpString(path []string, a, b reflect.Value, parent any) error {
	if changed, err := d.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if a.String() != b.String() {
		if a.CanInterface() {
			// If a and/or b is of a type that is an alias for String, store that type in changelog
			d.changes.add(CHANGE, path, getAsAny(a), getAsAny(b), parent)
		} else {
			d.changes.add(CHANGE, path, a.String(), b.String(), parent)
		}
	}

	return nil
}
