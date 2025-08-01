package compare

import "reflect"

func (c *Comparer) cmpStruct(path []string, left, right reflect.Value, parent any) error {
	if left.Kind() == reflect.Invalid {
		if c.summarizeMissingStructs {
			c.changes.add(ADD, path, nil, getAsAny(right))
			return nil
		} else {
			return c.cmpStructValuesForInvalid(ADD, path, right)
		}
	}

	if right.Kind() == reflect.Invalid {
		if c.summarizeMissingStructs {
			c.changes.add(REMOVE, path, getAsAny(left), nil)
			return nil
		} else {
			return c.cmpStructValuesForInvalid(REMOVE, path, left)
		}
	}

	for i := 0; i < left.NumField(); i++ {
		field := left.Type().Field(i)
		tName := getTagName(c.tagName, field)

		if tName == "-" || hasTagOption(c.tagName, field, "immutable") {
			continue
		}

		if tName == "" {
			tName = field.Name
		}

		af := left.Field(i)
		bf := right.FieldByName(field.Name)

		fpath := path
		if !(c.embeddedStructFields && field.Anonymous) {
			fpath = copyAppend(fpath, tName)
		}

		// skip private fields
		if !left.CanInterface() {
			continue
		}

		if err := c.compare(fpath, af, bf, getAsAny(left)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Comparer) cmpStructValuesForInvalid(t DiffType, path []string, val reflect.Value) error {
	var nc *Comparer = c.clone()

	if t != ADD && t != REMOVE {
		return ErrInvalidChangeType
	}

	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	if val.Kind() != reflect.Struct {
		return ErrTypeMismatch
	}

	x := reflect.New(val.Type()).Elem()

	for i := 0; i < val.NumField(); i++ {

		field := val.Type().Field(i)
		tName := getTagName(c.tagName, field)

		if tName == "-" {
			continue
		}

		if tName == "" {
			tName = field.Name
		}

		af := val.Field(i)
		xf := x.FieldByName(field.Name)

		fpath := copyAppend(path, tName)

		err := nc.compare(fpath, xf, af, getAsAny(val))
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(nc.changes); i++ {
		c.changes = append(c.changes, patchChange(t, nc.changes[i]))
	}

	return nil
}
