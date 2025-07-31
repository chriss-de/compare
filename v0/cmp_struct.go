package compare

import "reflect"

func (d *Comparer) cmpStruct(path []string, a, b reflect.Value, parent any) error {
	if a.Kind() == reflect.Invalid {
		if d.summarizeMissingStructs {
			d.changes.add(ADD, path, nil, getAsAny(b))
			return nil
		} else {
			return d.cmpStructValuesForInvalid(ADD, path, b)
		}
	}

	if b.Kind() == reflect.Invalid {
		if d.summarizeMissingStructs {
			d.changes.add(REMOVE, path, getAsAny(a), nil)
			return nil
		} else {
			return d.cmpStructValuesForInvalid(REMOVE, path, a)
		}
	}

	for i := 0; i < a.NumField(); i++ {
		field := a.Type().Field(i)
		tName := getTagName(d.tagName, field)

		if tName == "-" || hasTagOption(d.tagName, field, "immutable") {
			continue
		}

		if tName == "" {
			tName = field.Name
		}

		af := a.Field(i)
		bf := b.FieldByName(field.Name)

		fpath := path
		if !(d.embeddedStructFields && field.Anonymous) {
			fpath = copyAppend(fpath, tName)
		}

		//if d.Filter != nil && !d.Filter(fpath, a.Type(), field) {
		//	continue
		//}

		// skip private fields
		if !a.CanInterface() {
			continue
		}

		if err := d.compare(fpath, af, bf, getAsAny(a)); err != nil {
			return err
		}
	}

	return nil
}

func (d *Comparer) cmpStructValuesForInvalid(t ChangeType, path []string, a reflect.Value) error {
	var nd Comparer
	//nd.Filter = d.Filter
	//nd.customValueDiffers = d.customValueDiffers

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
		tName := getTagName(d.tagName, field)

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
		d.changes = append(d.changes, patchChange(t, nd.changes[i]))
	}

	return nil
}
