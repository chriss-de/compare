package go_compare

import (
	"reflect"
)

// Comparable ...
type Comparable struct {
	A, B *reflect.Value
}

// ComparableList : stores indexed isComparable
type ComparableList struct {
	m    map[any]*Comparable
	keys []any
}

// NewComparableList : returns a new isComparable list
func NewComparableList() *ComparableList {
	return &ComparableList{
		m:    make(map[any]*Comparable),
		keys: make([]any, 0),
	}
}

func (cl *ComparableList) addA(k any, v *reflect.Value) {
	if (*cl).m[k] == nil {
		(*cl).m[k] = &Comparable{}
		(*cl).keys = append((*cl).keys, k)
	}
	(*cl).m[k].A = v
}

func (cl *ComparableList) addB(k any, v *reflect.Value) {
	if (*cl).m[k] == nil {
		(*cl).m[k] = &Comparable{}
		(*cl).keys = append((*cl).keys, k)
	}
	(*cl).m[k].B = v
}
