package compare

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
	changes                 Differences
}

// NewComparer creates a new configurable diffing object
func NewComparer(opts ...CompareOptsFunc) (*Comparer, error) {
	d := Comparer{tagName: "cmp", summarizeMissingStructs: false, sliceOrdering: false, structMapKeys: false, embeddedStructFields: true}

	for _, opt := range opts {
		err := opt(&d)
		if err != nil {
			return nil, err
		}
	}

	return &d, nil
}

func (c *Comparer) clone() *Comparer {
	nc := &Comparer{
		tagName:                 c.tagName,
		summarizeMissingStructs: c.summarizeMissingStructs,
		sliceOrdering:           c.sliceOrdering,
		structMapKeys:           c.structMapKeys,
		embeddedStructFields:    c.embeddedStructFields,
	}
	return nc
}

func (c *Comparer) getCompareFunc(left, right reflect.Value) (Type, CompareFunc) {
	switch {
	case areType(left, right, reflect.TypeOf(time.Time{})):
		return TIME, c.cmpTime
	case areKind(left, right, reflect.Struct, reflect.Invalid):
		return STRUCT, c.cmpStruct
	case areKind(left, right, reflect.Slice, reflect.Invalid):
		return SLICE, c.cmpSlice
	case areKind(left, right, reflect.Array, reflect.Invalid):
		return ARRAY, c.cmpSlice
	case areKind(left, right, reflect.String, reflect.Invalid):
		return STRING, c.cmpString
	case areKind(left, right, reflect.Bool, reflect.Invalid):
		return BOOL, c.cmpBool
	case areKind(left, right, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Invalid):
		return INT, c.cmpInt
	case areKind(left, right, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Invalid):
		return UINT, c.cmpUint
	case areKind(left, right, reflect.Float32, reflect.Float64, reflect.Invalid):
		return FLOAT, c.cmpFloat
	case areKind(left, right, reflect.Map, reflect.Invalid):
		return MAP, c.cmpMap
	case areKind(left, right, reflect.Ptr, reflect.Invalid):
		return PTR, c.cmpPtr
	case areKind(left, right, reflect.Interface, reflect.Invalid):
		return INTERFACE, c.cmpInterface
	default:
		return UNSUPPORTED, nil
	}
}

// Compare returns a changelog of all mutated values from both
func (c *Comparer) Compare(left, right any) (Differences, error) {
	// reset the state of the compare
	c.changes = Differences{}

	return c.changes, c.compare([]string{}, reflect.ValueOf(left), reflect.ValueOf(right), nil)
}

func (c *Comparer) compare(path []string, left, right reflect.Value, parent any) error {
	// check if types match or areKind
	if isInvalid(left, right) {
		//if c.AllowTypeMismatch {
		//	c.changes.Add(CHANGE, path, left.Interface(), right.Interface())
		//	return nil
		//}
		return ErrTypeMismatch
	}

	// get the compare type and the corresponding built-int compare function to handle this type
	cmpType, compareFunc := c.getCompareFunc(left, right)

	// first go through custom compare functions
	//if len(c.customValueDiffers) > 0 {
	//	for _, vd := range c.customValueDiffers {
	//		if vd.Match(left, right) {
	//			err := vd.Compare(cmpType, compareFunc, &c.changes, path, left, right, parent)
	//			if err != nil {
	//				return err
	//			}
	//			return nil
	//		}
	//	}
	//}

	// then
	if cmpType == UNSUPPORTED {
		return errors.New("unsupported type: " + left.Kind().String())
	}

	return compareFunc(path, left, right, parent)
}

// cmpDefault does basic compare operations
func (c *Comparer) cmpDefault(path []string, left, right reflect.Value) (changed bool, err error) {
	if left.Kind() == reflect.Invalid {
		c.changes.add(ADD, path, nil, getAsAny(right))
		return true, nil
	}

	if right.Kind() == reflect.Invalid {
		c.changes.add(REMOVE, path, getAsAny(left), nil)
		return true, nil
	}

	if left.Kind() != right.Kind() {
		return false, ErrTypeMismatch
	}

	return false, nil
}
