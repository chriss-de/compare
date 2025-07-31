package go_compare

import "reflect"

func (d *Comparer) cmpInt(path []string, a, b reflect.Value, parent any) error {
	if changed, err := d.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if a.Int() != b.Int() {
		if a.CanInterface() {
			d.changes.add(CHANGE, path, getAsAny(a), getAsAny(b), parent)
		} else {
			d.changes.add(CHANGE, path, a.Int(), b.Int(), parent)
		}
	}

	return nil
}
