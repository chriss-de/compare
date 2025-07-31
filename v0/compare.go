package go_compare

import (
	"errors"
	"reflect"
	"time"
)

// CompareType represents an enum with all the supported compare types
type Type uint8

const (
	UNSUPPORTED Type = iota
	TIME
	STRUCT
	SLICE
	ARRAY
	STRING
	BOOL
	INT
	UINT
	FLOAT
	MAP
	PTR
	INTERFACE
)

// CompareFunc represents the built-in compare functions
type CompareFunc func([]string, reflect.Value, reflect.Value, any) error

// Comparer a configurable compare instance
type Comparer struct {
	tagName                 string
	summarizeMissingStructs bool
	sliceOrdering           bool
	structMapKeys           bool
	embeddedStructFields    bool
	//useStructIdInSlices     bool
	//customValueDiffers []ValueDiffer
	changes Changes
}

// NewComparer creates a new configurable diffing object
func NewComparer(opts ...func(d *Comparer) error) (*Comparer, error) {
	d := Comparer{tagName: "cmp", summarizeMissingStructs: false, sliceOrdering: false, structMapKeys: false, embeddedStructFields: true}

	for _, opt := range opts {
		err := opt(&d)
		if err != nil {
			return nil, err
		}
	}

	return &d, nil
}

func (d *Comparer) getCompareFunc(a, b reflect.Value) (Type, CompareFunc) {
	switch {
	case areType(a, b, reflect.TypeOf(time.Time{})):
		return TIME, d.cmpTime
	case areKind(a, b, reflect.Struct, reflect.Invalid):
		return STRUCT, d.cmpStruct
	case areKind(a, b, reflect.Slice, reflect.Invalid):
		return SLICE, d.cmpSlice
	case areKind(a, b, reflect.Array, reflect.Invalid):
		return ARRAY, d.cmpSlice
	case areKind(a, b, reflect.String, reflect.Invalid):
		return STRING, d.cmpString
	case areKind(a, b, reflect.Bool, reflect.Invalid):
		return BOOL, d.cmpBool
	case areKind(a, b, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Invalid):
		return INT, d.cmpInt
	case areKind(a, b, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Invalid):
		return UINT, d.cmpUint
	case areKind(a, b, reflect.Float32, reflect.Float64, reflect.Invalid):
		return FLOAT, d.cmpFloat
	case areKind(a, b, reflect.Map, reflect.Invalid):
		return MAP, d.cmpMap
	case areKind(a, b, reflect.Ptr, reflect.Invalid):
		return PTR, d.cmpPtr
	case areKind(a, b, reflect.Interface, reflect.Invalid):
		return INTERFACE, d.cmpInterface
	default:
		return UNSUPPORTED, nil
	}
}

// Compare returns a changelog of all mutated values from both
func (d *Comparer) Compare(a, b any) (Changes, error) {
	// reset the state of the compare
	d.changes = Changes{}

	return d.changes, d.compare([]string{}, reflect.ValueOf(a), reflect.ValueOf(b), nil)
}

func (d *Comparer) compare(path []string, a, b reflect.Value, parent any) error {
	// check if types match or areKind
	if isInvalid(a, b) {
		//if d.AllowTypeMismatch {
		//	d.changes.Add(CHANGE, path, a.Interface(), b.Interface())
		//	return nil
		//}
		return ErrTypeMismatch
	}

	// get the compare type and the corresponding built-int compare function to handle this type
	cmpType, compareFunc := d.getCompareFunc(a, b)

	// first go through custom compare functions
	//if len(d.customValueDiffers) > 0 {
	//	for _, vd := range d.customValueDiffers {
	//		if vd.Match(a, b) {
	//			err := vd.Compare(cmpType, compareFunc, &d.changes, path, a, b, parent)
	//			if err != nil {
	//				return err
	//			}
	//			return nil
	//		}
	//	}
	//}

	// then
	if cmpType == UNSUPPORTED {
		return errors.New("unsupported type: " + a.Kind().String())
	}

	return compareFunc(path, a, b, parent)
}

// cmpDefault does basic compare operations
func (d *Comparer) cmpDefault(path []string, a, b reflect.Value) (changed bool, err error) {
	if a.Kind() == reflect.Invalid {
		d.changes.add(ADD, path, nil, getAsAny(b))
		return true, nil
	}

	if b.Kind() == reflect.Invalid {
		d.changes.add(REMOVE, path, getAsAny(a), nil)
		return true, nil
	}

	if a.Kind() != b.Kind() {
		return false, ErrTypeMismatch
	}

	return false, nil
}
