package compare

import "reflect"

func (d *Comparer) cmpBool(path []string, a, b reflect.Value, parent any) error {
	if changed, err := d.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if a.Bool() != b.Bool() {
		d.changes.add(CHANGE, path, getAsAny(a), getAsAny(b), parent)
	}

	return nil
}
