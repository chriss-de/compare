package compare

import (
	"reflect"
)

// Comparable ...
type Comparable struct {
	LEFT, RIGHT *reflect.Value
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

func (cmpList *ComparableList) addLeft(key any, val *reflect.Value) {
	if (*cmpList).m[key] == nil {
		(*cmpList).m[key] = &Comparable{}
		(*cmpList).keys = append((*cmpList).keys, key)
	}
	(*cmpList).m[key].LEFT = val
}

func (cmpList *ComparableList) addRight(key any, val *reflect.Value) {
	if (*cmpList).m[key] == nil {
		(*cmpList).m[key] = &Comparable{}
		(*cmpList).keys = append((*cmpList).keys, key)
	}
	(*cmpList).m[key].RIGHT = val
}
