package compare

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"reflect"
)

func (c *Comparer) cmpMap(path []string, left, right reflect.Value, parent any) error {
	if left.Kind() == reflect.Invalid {
		return c.cmpMapValuesForInvalid(ADD, path, right)
	}

	if right.Kind() == reflect.Invalid {
		return c.cmpMapValuesForInvalid(REMOVE, path, left)
	}

	cmpList := NewComparableList()

	for _, k := range left.MapKeys() {
		leftElem := left.MapIndex(k)
		cmpList.addLeft(getAsAny(k), &leftElem)
	}

	for _, k := range right.MapKeys() {
		rightElem := right.MapIndex(k)
		cmpList.addRight(getAsAny(k), &rightElem)
	}

	return c.processComparableList(path, cmpList, getAsAny(left))
}

func (c *Comparer) cmpMapValuesForInvalid(t DiffType, path []string, val reflect.Value) error {
	if t != ADD && t != REMOVE {
		return ErrInvalidChangeType
	}

	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	if val.Kind() != reflect.Map {
		return ErrTypeMismatch
	}

	x := reflect.New(val.Type()).Elem()

	for _, k := range val.MapKeys() {
		ae := val.MapIndex(k)
		xe := x.MapIndex(k)

		var err error
		if c.config.structMapKeys {
			var bWriter = new(bytes.Buffer)
			if err = gob.NewEncoder(bWriter).Encode(k.Interface()); err == nil {
				key := base64.RawStdEncoding.EncodeToString(bWriter.Bytes())
				err = c.compare(append(path, key), xe, ae, val.Interface())
			}

		} else {
			err = c.compare(append(path, fmt.Sprint(k.Interface())), xe, ae, val.Interface())
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
