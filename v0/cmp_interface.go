package go_compare

import "reflect"

func (d *Comparer) cmpInterface(path []string, a, b reflect.Value, parent any) error {
	if changed, err := d.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	if a.IsNil() && b.IsNil() {
		return nil
	}

	if a.IsNil() {
		d.changes.add(CHANGE, path, nil, getAsAny(b), parent)
		return nil
	}

	if b.IsNil() {
		d.changes.add(CHANGE, path, getAsAny(a), nil, parent)
		return nil
	}

	return d.compare(path, a.Elem(), b.Elem(), parent)
}
