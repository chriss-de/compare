package go_compare

import "reflect"

func (d *Comparer) cmpFloat(path []string, a, b reflect.Value, parent any) error {
	if changed, err := d.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if a.Float() != b.Float() {
		if a.CanInterface() {
			d.changes.add(CHANGE, path, getAsAny(a), getAsAny(b), parent)
		} else {
			d.changes.add(CHANGE, path, a.Float(), b.Float(), parent)
		}
	}

	return nil
}
