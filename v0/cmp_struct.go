package compare

import "reflect"

func (c *Comparer) cmpStruct(path []string, a, b reflect.Value, parent any) error {
	if a.Kind() == reflect.Invalid {
		if c.summarizeMissingStructs {
			c.changes.add(ADD, path, nil, getAsAny(b))
			return nil
		} else {
			return c.cmpStructValuesForInvalid(ADD, path, b)
		}
	}

	if b.Kind() == reflect.Invalid {
		if c.summarizeMissingStructs {
			c.changes.add(REMOVE, path, getAsAny(a), nil)
			return nil
		} else {
			return c.cmpStructValuesForInvalid(REMOVE, path, a)
		}
	}

	for i := 0; i < a.NumField(); i++ {
		field := a.Type().Field(i)
		tName := getTagName(c.tagName, field)

		if tName == "-" || hasTagOption(c.tagName, field, "immutable") {
			continue
		}

		if tName == "" {
			tName = field.Name
		}

		af := a.Field(i)
		bf := b.FieldByName(field.Name)

		fpath := path
		if !(c.embeddedStructFields && field.Anonymous) {
			fpath = copyAppend(fpath, tName)
		}

		//if c.Filter != nil && !c.Filter(fpath, a.Type(), field) {
		//	continue
		//}

		// skip private fields
		if !a.CanInterface() {
			continue
		}

		if err := c.compare(fpath, af, bf, getAsAny(a)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Comparer) cmpStructValuesForInvalid(t DiffType, path []string, a reflect.Value) error {
	var nd Comparer
	//nd.Filter = c.Filter
	//nd.customValueDiffers = c.customValueDiffers

	if t != ADD && t != REMOVE {
		return ErrInvalidChangeType
	}

	if a.Kind() == reflect.Ptr {
		a = reflect.Indirect(a)
	}

	if a.Kind() != reflect.Struct {
		return ErrTypeMismatch
	}

	x := reflect.New(a.Type()).Elem()

	for i := 0; i < a.NumField(); i++ {

		field := a.Type().Field(i)
		tName := getTagName(c.tagName, field)

		if tName == "-" {
			continue
		}

		if tName == "" {
			tName = field.Name
		}

		af := a.Field(i)
		xf := x.FieldByName(field.Name)

		fpath := copyAppend(path, tName)

		//if nd.Filter != nil && !nd.Filter(fpath, a.Type(), field) {
		//	continue
		//}

		err := nd.compare(fpath, xf, af, getAsAny(a))
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(nd.changes); i++ {
		c.changes = append(c.changes, patchChange(t, nd.changes[i]))
	}

	return nil
}
