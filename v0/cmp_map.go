package compare

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"reflect"
)

func (d *Comparer) cmpMap(path []string, a, b reflect.Value, parent any) error {
	if a.Kind() == reflect.Invalid {
		return d.cmpMapValuesForInvalid(ADD, path, b)
	}

	if b.Kind() == reflect.Invalid {
		return d.cmpMapValuesForInvalid(REMOVE, path, a)
	}

	c := NewComparableList()

	for _, k := range a.MapKeys() {
		ae := a.MapIndex(k)
		c.addA(getAsAny(k), &ae)
	}

	for _, k := range b.MapKeys() {
		be := b.MapIndex(k)
		c.addB(getAsAny(k), &be)
	}

	return d.processComparableList(path, c, getAsAny(a))
}

func (d *Comparer) cmpMapValuesForInvalid(t ChangeType, path []string, a reflect.Value) error {
	if t != ADD && t != REMOVE {
		return ErrInvalidChangeType
	}

	if a.Kind() == reflect.Ptr {
		a = reflect.Indirect(a)
	}

	if a.Kind() != reflect.Map {
		return ErrTypeMismatch
	}

	x := reflect.New(a.Type()).Elem()

	for _, k := range a.MapKeys() {
		ae := a.MapIndex(k)
		xe := x.MapIndex(k)

		var err error
		if d.structMapKeys {
			var bWriter = new(bytes.Buffer)
			if err = gob.NewEncoder(bWriter).Encode(k.Interface()); err == nil {
				key := base64.RawStdEncoding.EncodeToString(bWriter.Bytes())
				err = d.compare(append(path, key), xe, ae, a.Interface())
			}

		} else {
			err = d.compare(append(path, fmt.Sprint(k.Interface())), xe, ae, a.Interface())
		}
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(d.changes); i++ {
		// only swap changes on the relevant map
		if pathMatches(path, d.changes[i].Path) {
			d.changes[i] = patchChange(t, d.changes[i])
		}
	}

	return nil
}
