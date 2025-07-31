package compare

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"reflect"
)

func (c *Comparer) cmpMap(path []string, a, b reflect.Value, parent any) error {
	if a.Kind() == reflect.Invalid {
		return c.cmpMapValuesForInvalid(ADD, path, b)
	}

	if b.Kind() == reflect.Invalid {
		return c.cmpMapValuesForInvalid(REMOVE, path, a)
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

	return c.processComparableList(path, c, getAsAny(a))
}

func (c *Comparer) cmpMapValuesForInvalid(t DiffType, path []string, a reflect.Value) error {
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
		if c.structMapKeys {
			var bWriter = new(bytes.Buffer)
			if err = gob.NewEncoder(bWriter).Encode(k.Interface()); err == nil {
				key := base64.RawStdEncoding.EncodeToString(bWriter.Bytes())
				err = c.compare(append(path, key), xe, ae, a.Interface())
			}

		} else {
			err = c.compare(append(path, fmt.Sprint(k.Interface())), xe, ae, a.Interface())
		}
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(c.changes); i++ {
		// only swap changes on the relevant map
		if pathMatches(path, c.changes[i].Path) {
			c.changes[i] = patchChange(t, c.changes[i])
		}
	}

	return nil
}
