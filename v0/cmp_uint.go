package compare

import "reflect"

func (d *Comparer) cmpUint(path []string, a, b reflect.Value, parent any) error {
	if changed, err := d.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if a.Uint() != b.Uint() {
		if a.CanInterface() {
			d.changes.add(CHANGE, path, getAsAny(a), getAsAny(b), parent)
		} else {
			d.changes.add(CHANGE, path, a.Uint(), b.Uint(), parent)
		}
	}

	return nil
}
