package go_compare

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type SimpleStructNoTag struct {
	Name  string
	Value int
}

type SimpleStructWithTag struct {
	Name  string `cmp:"name,identifier"`
	Value int    `cmp:"value"`
}

type NoIdentifierStruct struct {
	Value int `cmp:"value"`
}

type EmbeddedStruct struct {
	SimpleStructNoTag
	Baz bool `cmp:"baz"`
}

type ComplexStruct struct {
	ID              string                `cmp:"id"`
	Name            string                `cmp:"name"`
	Value           int                   `cmp:"value"`
	Bool            bool                  `cmp:"bool"`
	Values          []string              `cmp:"values"`
	Map             map[string]string     `cmp:"map"`
	Time            time.Time             `cmp:"time"`
	Pointer         *string               `cmp:"pointer"`
	Ignored         bool                  `cmp:"-"`
	Identifiables   []SimpleStructWithTag `cmp:"identifiables"`
	Unidentifiables []NoIdentifierStruct  `cmp:"unidentifiables"`
	private         int                   `cmp:"private"`
}

type RealWorldSubStruct struct {
	Name string `cmp:"name"`
	ID   int64  `cmp:"-,identifier"`
}
type RealWorldStruct struct {
	Name   string                `cmp:"name,identifier"`
	Value  int                   `cmp:"value"`
	Addons []*RealWorldSubStruct `cmp:"addons"`
}

var testTimeA, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
var testTimeB, _ = time.Parse(time.RFC3339, "2007-01-02T15:04:05Z")

func getStringPointer(s string) *string {
	return &s
}

func TestCompare(t *testing.T) {
	cases := []struct {
		Name    string
		A, B    any
		Changes Changes
		Error   error
	}{
		{
			"uint-equal", uint(1), uint(1),
			Changes{},
			nil,
		},
		{
			"uint-not-equal", uint(1), uint(2),
			Changes{
				Change{Type: CHANGE, Path: []string{}, From: uint(1), To: uint(2)},
			},
			nil,
		},
		{
			"int-equal", int(1), int(1),
			Changes{},
			nil,
		},
		{
			"int-not-equal", int(1), int(2),
			Changes{
				Change{Type: CHANGE, Path: []string{}, From: int(1), To: int(2)},
			},
			nil,
		},
		{
			"float-equal", float64(1.1), float64(1.1),
			Changes{},
			nil,
		},
		{
			"float-not-equal", float64(1.1), float64(2.2),
			Changes{
				Change{Type: CHANGE, Path: []string{}, From: float64(1.1), To: float64(2.2)},
			},
			nil,
		},
		{
			"string-equal", "hello", "hello",
			Changes{},
			nil,
		},
		{
			"string-not-equal", "hello", "world",
			Changes{
				Change{Type: CHANGE, Path: []string{}, From: "hello", To: "world"},
			},
			nil,
		},
		{
			"time-equal", testTimeA, testTimeA,
			Changes{},
			nil,
		},
		{
			"time-not-equal", testTimeA, testTimeB,
			Changes{
				Change{Type: CHANGE, Path: []string{}, From: testTimeA, To: testTimeB},
			},
			nil,
		},
		{
			"SimpleStructNoTag-equal", SimpleStructNoTag{Name: "test A", Value: 123}, SimpleStructNoTag{Name: "test A", Value: 123},
			Changes{},
			nil,
		},
		{
			"SimpleStructNoTag-not-equal-name", SimpleStructNoTag{Name: "test A", Value: 123}, SimpleStructNoTag{Name: "test B", Value: 123},
			Changes{
				Change{Type: CHANGE, Path: []string{"Name"}, From: "test A", To: "test B"},
			},
			nil,
		},
		{
			"SimpleStructNoTag-not-equal-name-and-value", SimpleStructNoTag{Name: "test A", Value: 123}, SimpleStructNoTag{Name: "test B", Value: 456},
			Changes{
				Change{Type: CHANGE, Path: []string{"Name"}, From: "test A", To: "test B"},
				Change{Type: CHANGE, Path: []string{"Value"}, From: 123, To: 456},
			},
			nil,
		},
		{
			"SimpleStructWithTag-equal", SimpleStructWithTag{Name: "test A", Value: 123}, SimpleStructWithTag{Name: "test A", Value: 123},
			Changes{},
			nil,
		},
		{
			"SimpleStructWithTag-not-equal", SimpleStructWithTag{Name: "test A", Value: 123}, SimpleStructWithTag{Name: "test B", Value: 456},
			Changes{
				Change{Type: CHANGE, Path: []string{"name"}, From: "test A", To: "test B"},
				Change{Type: CHANGE, Path: []string{"value"}, From: 123, To: 456},
			},
			nil,
		},
		{
			"SimpleStructWithTag-equal", SimpleStructWithTag{Name: "test A", Value: 123}, SimpleStructWithTag{Name: "test A", Value: 123},
			Changes{},
			nil,
		},
		{
			"SimpleStructWithTag-not-equal", SimpleStructWithTag{Name: "test A", Value: 123}, SimpleStructWithTag{Name: "test B", Value: 123},
			Changes{
				Change{Type: CHANGE, Path: []string{"name"}, From: "test A", To: "test B"},
			},
			nil,
		},
		{
			"SimpleStructWithTag-not-equal", SimpleStructWithTag{Name: "test A", Value: 123}, SimpleStructWithTag{Name: "test B", Value: 456},
			Changes{
				Change{Type: CHANGE, Path: []string{"name"}, From: "test A", To: "test B"},
				Change{Type: CHANGE, Path: []string{"value"}, From: 123, To: 456},
			},
			nil,
		},
		{
			"SimpleStructWithTag-not-equal", SimpleStructWithTag{Name: "test A", Value: 123}, SimpleStructWithTag{Name: "test B", Value: 456},
			Changes{
				Change{Type: CHANGE, Path: []string{"name"}, From: "test A", To: "test B"},
				Change{Type: CHANGE, Path: []string{"value"}, From: 123, To: 456},
			},
			nil,
		},
		{
			"different-structs", SimpleStructWithTag{Name: "test A", Value: 123}, SimpleStructNoTag{Name: "test A", Value: 123},
			Changes{},
			nil,
		},
		{
			"different-structs", SimpleStructNoTag{Name: "test A", Value: 123}, SimpleStructWithTag{Name: "test A", Value: 123},
			Changes{},
			nil,
		},
		{
			"int-slice-insert", []int{1, 2, 3}, []int{1, 2, 3, 4},
			Changes{
				Change{Type: ADD, Path: []string{"3"}, To: 4},
			},
			nil,
		},
		{
			"int-array-insert", [3]int{1, 2, 3}, [4]int{1, 2, 3, 4},
			Changes{
				Change{Type: ADD, Path: []string{"3"}, To: 4},
			},
			nil,
		},
		{
			"int-slice-delete", []int{1, 2, 3}, []int{1, 3},
			Changes{
				Change{Type: REMOVE, Path: []string{"1"}, From: 2},
			},
			nil,
		},
		{
			"int-array-delete", [3]int{1, 2, 3}, [2]int{1, 3},
			Changes{
				Change{Type: REMOVE, Path: []string{"1"}, From: 2},
			},
			nil,
		},
		{
			"uint-slice-insert", []uint{1, 2, 3}, []uint{1, 2, 3, 4},
			Changes{
				Change{Type: ADD, Path: []string{"3"}, To: uint(4)},
			},
			nil,
		},
		{
			"uint-array-insert", [3]uint{1, 2, 3}, [4]uint{1, 2, 3, 4},
			Changes{
				Change{Type: ADD, Path: []string{"3"}, To: uint(4)},
			},
			nil,
		},
		{
			"uint-slice-delete", []uint{1, 2, 3}, []uint{1, 3},
			Changes{
				Change{Type: REMOVE, Path: []string{"1"}, From: uint(2)},
			},
			nil,
		},
		{
			"uint-array-delete", [3]uint{1, 2, 3}, [2]uint{1, 3},
			Changes{
				Change{Type: REMOVE, Path: []string{"1"}, From: uint(2)},
			},
			nil,
		},
		{
			"string-slice-insert", []string{"1", "2", "3"}, []string{"1", "2", "3", "4"},
			Changes{
				Change{Type: ADD, Path: []string{"3"}, To: "4"},
			},
			nil,
		},
		{
			"string-array-insert", [3]string{"1", "2", "3"}, [4]string{"1", "2", "3", "4"},
			Changes{
				Change{Type: ADD, Path: []string{"3"}, To: "4"},
			},
			nil,
		},
		{
			"string-slice-delete", []string{"1", "2", "3"}, []string{"1", "3"},
			Changes{
				Change{Type: REMOVE, Path: []string{"1"}, From: "2"},
			},
			nil,
		},
		{
			"string-slice-delete", [3]string{"1", "2", "3"}, [2]string{"1", "3"},
			Changes{
				Change{Type: REMOVE, Path: []string{"1"}, From: "2"},
			},
			nil,
		},
		{
			"string-slice-insert-delete", []string{"1", "2", "3"}, []string{"1", "3", "4"},
			Changes{
				Change{Type: REMOVE, Path: []string{"1"}, From: "2"},
				Change{Type: ADD, Path: []string{"2"}, To: "4"},
			},
			nil,
		},
		{
			"string-array-insert-delete", [3]string{"1", "2", "3"}, [3]string{"1", "3", "4"},
			Changes{
				Change{Type: REMOVE, Path: []string{"1"}, From: "2"},
				Change{Type: ADD, Path: []string{"2"}, To: "4"},
			},
			nil,
		},
		{
			"isComparable-slice-insert", []SimpleStructWithTag{{"one", 1}}, []SimpleStructWithTag{{"one", 1}, {"two", 2}},
			Changes{
				Change{Type: ADD, Path: []string{"two", "name"}, To: "two"},
				Change{Type: ADD, Path: []string{"two", "value"}, To: 2},
			},
			nil,
		},
		{
			"isComparable-array-insert", [1]SimpleStructWithTag{{"one", 1}}, [2]SimpleStructWithTag{{"one", 1}, {"two", 2}},
			Changes{
				Change{Type: ADD, Path: []string{"two", "name"}, To: "two"},
				Change{Type: ADD, Path: []string{"two", "value"}, To: 2},
			},
			nil,
		},
		{
			"isComparable-slice-delete", []SimpleStructWithTag{{"one", 1}, {"two", 2}}, []SimpleStructWithTag{{"one", 1}},
			Changes{
				Change{Type: REMOVE, Path: []string{"two", "name"}, From: "two"},
				Change{Type: REMOVE, Path: []string{"two", "value"}, From: 2},
			},
			nil,
		},
		{
			"isComparable-array-delete", [2]SimpleStructWithTag{{"one", 1}, {"two", 2}}, [1]SimpleStructWithTag{{"one", 1}},
			Changes{
				Change{Type: REMOVE, Path: []string{"two", "name"}, From: "two"},
				Change{Type: REMOVE, Path: []string{"two", "value"}, From: 2},
			},
			nil,
		},
		{
			"isComparable-slice-update", []SimpleStructWithTag{{"one", 1}}, []SimpleStructWithTag{{"one", 50}},
			Changes{
				Change{Type: CHANGE, Path: []string{"one", "value"}, From: 1, To: 50},
			},
			nil,
		},
		{
			"isComparable-array-update", [1]SimpleStructWithTag{{"one", 1}}, [1]SimpleStructWithTag{{"one", 50}},
			Changes{
				Change{Type: CHANGE, Path: []string{"one", "value"}, From: 1, To: 50},
			},
			nil,
		},
		{
			"map-slice-insert", []map[string]string{{"test": "123"}}, []map[string]string{{"test": "123", "tset": "456"}},
			Changes{
				Change{Type: ADD, Path: []string{"0", "tset"}, To: "456"},
			},
			nil,
		},
		{
			"map-array-insert", [1]map[string]string{{"test": "123"}}, [1]map[string]string{{"test": "123", "tset": "456"}},
			Changes{
				Change{Type: ADD, Path: []string{"0", "tset"}, To: "456"},
			},
			nil,
		},
		{
			"map-slice-update", []map[string]string{{"test": "123"}}, []map[string]string{{"test": "456"}},
			Changes{
				Change{Type: CHANGE, Path: []string{"0", "test"}, From: "123", To: "456"},
			},
			nil,
		},
		{
			"map-array-update", [1]map[string]string{{"test": "123"}}, [1]map[string]string{{"test": "456"}},
			Changes{
				Change{Type: CHANGE, Path: []string{"0", "test"}, From: "123", To: "456"},
			},
			nil,
		},
		{
			"map-slice-delete", []map[string]string{{"test": "123", "tset": "456"}}, []map[string]string{{"test": "123"}},
			Changes{
				Change{Type: REMOVE, Path: []string{"0", "tset"}, From: "456"},
			},
			nil,
		},
		{
			"map-array-delete", [1]map[string]string{{"test": "123", "tset": "456"}}, [1]map[string]string{{"test": "123"}},
			Changes{
				Change{Type: REMOVE, Path: []string{"0", "tset"}, From: "456"},
			},
			nil,
		},
		{
			"map-interface-slice-update", []map[string]interface{}{{"test": nil}}, []map[string]interface{}{{"test": "456"}},
			Changes{
				Change{Type: CHANGE, Path: []string{"0", "test"}, From: nil, To: "456"},
			},
			nil,
		},
		{
			"map-interface-array-update", [1]map[string]interface{}{{"test": nil}}, [1]map[string]interface{}{{"test": "456"}},
			Changes{
				Change{Type: CHANGE, Path: []string{"0", "test"}, From: nil, To: "456"},
			},
			nil,
		},
		{
			"map-nil", map[string]string{"one": "test"}, nil,
			Changes{
				Change{Type: REMOVE, Path: []string{"one"}, From: "test", To: nil},
			},
			nil,
		},
		{
			"nil-map", nil, map[string]string{"one": "test"},
			Changes{
				Change{Type: ADD, Path: []string{"one"}, From: nil, To: "test"},
			},
			nil,
		},
		{
			"nested-map-insert", map[string]map[string]string{"a": {"test": "123"}}, map[string]map[string]string{"a": {"test": "123", "tset": "456"}},
			Changes{
				Change{Type: ADD, Path: []string{"a", "tset"}, To: "456"},
			},
			nil,
		},
		{
			"nested-map-interface-insert", map[string]map[string]interface{}{"a": {"test": "123"}}, map[string]map[string]interface{}{"a": {"test": "123", "tset": "456"}},
			Changes{
				Change{Type: ADD, Path: []string{"a", "tset"}, To: "456"},
			},
			nil,
		},
		{
			"nested-map-update", map[string]map[string]string{"a": {"test": "123"}}, map[string]map[string]string{"a": {"test": "456"}},
			Changes{
				Change{Type: CHANGE, Path: []string{"a", "test"}, From: "123", To: "456"},
			},
			nil,
		},
		{
			"nested-map-delete", map[string]map[string]string{"a": {"test": "123"}}, map[string]map[string]string{"a": {}},
			Changes{
				Change{Type: REMOVE, Path: []string{"a", "test"}, From: "123", To: nil},
			},
			nil,
		},
		{
			"nested-slice-insert", map[string][]int{"a": {1, 2, 3}}, map[string][]int{"a": {1, 2, 3, 4}},
			Changes{
				Change{Type: ADD, Path: []string{"a", "3"}, To: 4},
			},
			nil,
		},
		{
			"nested-array-insert", map[string][3]int{"a": {1, 2, 3}}, map[string][4]int{"a": {1, 2, 3, 4}},
			Changes{
				Change{Type: ADD, Path: []string{"a", "3"}, To: 4},
			},
			nil,
		},
		{
			"nested-slice-update", map[string][]int{"a": {1, 2, 3}}, map[string][]int{"a": {1, 4, 3}},
			Changes{
				Change{Type: CHANGE, Path: []string{"a", "1"}, From: 2, To: 4},
			},
			nil,
		},
		{
			"nested-array-update", map[string][3]int{"a": {1, 2, 3}}, map[string][3]int{"a": {1, 4, 3}},
			Changes{
				Change{Type: CHANGE, Path: []string{"a", "1"}, From: 2, To: 4},
			},
			nil,
		},
		{
			"nested-slice-delete", map[string][]int{"a": {1, 2, 3}}, map[string][]int{"a": {1, 3}},
			Changes{
				Change{Type: REMOVE, Path: []string{"a", "1"}, From: 2, To: nil},
			},
			nil,
		},
		{
			"nested-array-delete", map[string][3]int{"a": {1, 2, 3}}, map[string][2]int{"a": {1, 3}},
			Changes{
				Change{Type: REMOVE, Path: []string{"a", "1"}, From: 2, To: nil},
			},
			nil,
		},

		{
			"struct-string-update", ComplexStruct{Name: "one"}, ComplexStruct{Name: "two"},
			Changes{
				Change{Type: CHANGE, Path: []string{"name"}, From: "one", To: "two"},
			},
			nil,
		},
		{
			"struct-int-update", ComplexStruct{Value: 1}, ComplexStruct{Value: 50},
			Changes{
				Change{Type: CHANGE, Path: []string{"value"}, From: 1, To: 50},
			},
			nil,
		},
		{
			"struct-bool-update", ComplexStruct{Bool: true}, ComplexStruct{Bool: false},
			Changes{
				Change{Type: CHANGE, Path: []string{"bool"}, From: true, To: false},
			},
			nil,
		},
		{
			"struct-time-update", ComplexStruct{}, ComplexStruct{Time: testTimeA},
			Changes{
				Change{Type: CHANGE, Path: []string{"time"}, From: time.Time{}, To: testTimeA},
			},
			nil,
		},
		{
			"struct-map-update", ComplexStruct{Map: map[string]string{"test": "123"}}, ComplexStruct{Map: map[string]string{"test": "456"}},
			Changes{
				Change{Type: CHANGE, Path: []string{"map", "test"}, From: "123", To: "456"},
			},
			nil,
		},
		{
			"struct-string-pointer-update", ComplexStruct{Pointer: getStringPointer("test")}, ComplexStruct{Pointer: getStringPointer("test2")},
			Changes{
				Change{Type: CHANGE, Path: []string{"pointer"}, From: "test", To: "test2"},
			},
			nil,
		},
		{
			"struct-nil-string-pointer-update", ComplexStruct{Pointer: nil}, ComplexStruct{Pointer: getStringPointer("test")},
			Changes{
				Change{Type: CHANGE, Path: []string{"pointer"}, From: nil, To: getStringPointer("test")},
			},
			nil,
		},
		{
			"struct-generic-slice-insert", ComplexStruct{Values: []string{"one"}}, ComplexStruct{Values: []string{"one", "two"}},
			Changes{
				Change{Type: ADD, Path: []string{"values", "1"}, From: nil, To: "two"},
			},
			nil,
		},
		//{
		//	"struct-identifiable-slice-insert", ComplexStruct{Identifiables: []tistruct{{"one", 1}}}, ComplexStruct{Identifiables: []tistruct{{"one", 1}, {"two", 2}}},
		//	Changes{
		//		Change{Type: ADD, Path: []string{"identifiables", "two", "name"}, From: nil, To: "two"},
		//		Change{Type: ADD, Path: []string{"identifiables", "two", "value"}, From: nil, To: 2},
		//	},
		//	nil,
		//},
		{
			"struct-generic-slice-delete", ComplexStruct{Values: []string{"one", "two"}}, ComplexStruct{Values: []string{"one"}},
			Changes{
				Change{Type: REMOVE, Path: []string{"values", "1"}, From: "two", To: nil},
			},
			nil,
		},
		//{
		//	"struct-identifiable-slice-delete", ComplexStruct{Identifiables: []tistruct{{"one", 1}, {"two", 2}}}, ComplexStruct{Identifiables: []tistruct{{"one", 1}}},
		//	Changes{
		//		Change{Type: REMOVE, Path: []string{"identifiables", "two", "name"}, From: "two", To: nil},
		//		Change{Type: REMOVE, Path: []string{"identifiables", "two", "value"}, From: 2, To: nil},
		//	},
		//	nil,
		//},
		//{
		//	"struct-unidentifiable-slice-insert-delete", ComplexStruct{Unidentifiables: []tuistruct{{1}, {2}, {3}}}, ComplexStruct{Unidentifiables: []tuistruct{{5}, {2}, {3}, {4}}},
		//	Changes{
		//		Change{Type: CHANGE, Path: []string{"unidentifiables", "0", "value"}, From: 1, To: 5},
		//		Change{Type: ADD, Path: []string{"unidentifiables", "3", "value"}, From: nil, To: 4},
		//	},
		//	nil,
		//},
		//{
		//	"struct-with-private-value", privateValueStruct{Public: "one", Private: new(sync.RWMutex)}, privateValueStruct{Public: "two", Private: new(sync.RWMutex)},
		//	Changes{
		//		Change{Type: CHANGE, Path: []string{"Public"}, From: "one", To: "two"},
		//	},
		//	nil,
		//},
		//{
		//	"mismatched-values-struct-map", map[string]string{"test": "one"}, &ComplexStruct{Identifiables: []tistruct{{"one", 1}}},
		//	Changes{},
		//	ErrTypeMismatch,
		//},
		{
			"omittable", ComplexStruct{Ignored: false}, ComplexStruct{Ignored: true},
			Changes{},
			nil,
		},
		//{
		//	"slice", &ComplexStruct{}, &ComplexStruct{Nested: tnstruct{Slice: []tmstruct{{"one", 1}, {"two", 2}}}},
		//	Changes{
		//		Change{Type: ADD, Path: []string{"nested", "slice", "0", "foo"}, From: nil, To: "one"},
		//		Change{Type: ADD, Path: []string{"nested", "slice", "0", "bar"}, From: nil, To: 1},
		//		Change{Type: ADD, Path: []string{"nested", "slice", "1", "foo"}, From: nil, To: "two"},
		//		Change{Type: ADD, Path: []string{"nested", "slice", "1", "bar"}, From: nil, To: 2},
		//	},
		//	nil,
		//},
		{
			"slice-duplicate-items", []int{1}, []int{1, 1},
			Changes{
				Change{Type: ADD, Path: []string{"1"}, From: nil, To: 1},
			},
			nil,
		},
		{
			"mixed-slice-map", []map[string]interface{}{{"name": "name1", "type": []string{"null", "string"}}}, []map[string]interface{}{{"name": "name1", "type": []string{"null", "int"}}, {"name": "name2", "type": []string{"null", "string"}}},
			Changes{
				Change{Type: CHANGE, Path: []string{"0", "type", "1"}, From: "string", To: "int"},
				Change{Type: ADD, Path: []string{"1", "name"}, From: nil, To: "name2"},
				Change{Type: ADD, Path: []string{"1", "type"}, From: nil, To: []string{"null", "string"}},
			},
			nil,
		},
		//{
		//	"map-string-pointer-create",
		//	map[string]*tmstruct{"one": &struct1},
		//	map[string]*tmstruct{"one": &struct1, "two": &struct2},
		//	Changes{
		//		Change{Type: ADD, Path: []string{"two", "foo"}, From: nil, To: "two"},
		//		Change{Type: ADD, Path: []string{"two", "bar"}, From: nil, To: 2},
		//	},
		//	nil,
		//},
		//{
		//	"map-string-pointer-delete",
		//	map[string]*tmstruct{"one": &struct1, "two": &struct2},
		//	map[string]*tmstruct{"one": &struct1},
		//	Changes{
		//		Change{Type: REMOVE, Path: []string{"two", "foo"}, From: "two", To: nil},
		//		Change{Type: REMOVE, Path: []string{"two", "bar"}, From: 2, To: nil},
		//	},
		//	nil,
		//},
		{
			"private-struct-field",
			ComplexStruct{private: 1},
			ComplexStruct{private: 4},
			Changes{
				Change{Type: CHANGE, Path: []string{"private"}, From: int64(1), To: int64(4)},
			},
			nil,
		},
		{
			"embedded-struct-field",
			EmbeddedStruct{SimpleStructNoTag{Name: "a", Value: 2}, true},
			EmbeddedStruct{SimpleStructNoTag{Name: "b", Value: 3}, false},
			Changes{
				Change{Type: CHANGE, Path: []string{"Name"}, From: "a", To: "b"},
				Change{Type: CHANGE, Path: []string{"Value"}, From: 2, To: 3},
				Change{Type: CHANGE, Path: []string{"baz"}, From: true, To: false},
			},
			nil,
		},
		{
			"embedded-struct-field-as-extra-field",
			EmbeddedStruct{SimpleStructNoTag{Name: "a", Value: 2}, true},
			EmbeddedStruct{SimpleStructNoTag{Name: "b", Value: 3}, false},
			Changes{
				Change{Type: CHANGE, Path: []string{"SimpleStructNoTag", "Name"}, From: "a", To: "b"},
				Change{Type: CHANGE, Path: []string{"SimpleStructNoTag", "Value"}, From: 2, To: 3},
				Change{Type: CHANGE, Path: []string{"baz"}, From: true, To: false},
			},
			nil,
		},
		{
			"real-world-struct",
			RealWorldStruct{Name: "TestA", Value: 1, Addons: []*RealWorldSubStruct{{Name: "Sub1", ID: 10}, {Name: "Sub2", ID: 20}}},
			RealWorldStruct{Name: "TestB", Value: 1, Addons: []*RealWorldSubStruct{{Name: "Sub1", ID: 10}, {Name: "Sub3", ID: 30}}},
			Changes{
				Change{Type: CHANGE, Path: []string{"name"}, From: "TestA", To: "TestB"},
				Change{Type: REMOVE, Path: []string{"addons", "20", "name"}, From: "Sub2"},
				Change{Type: ADD, Path: []string{"addons", "30", "name"}, To: "Sub3"},
			},
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var options []func(d *Comparer) error
			switch tc.Name {
			//case "mixed-slice-map", "nil-map", "map-nil":
			//	options = append(options, WithStructMapKeys())
			case "embedded-struct-field-as-extra-field":
				options = append(options, WithEmbeddedStructsAsField())
			case "custom-tags":
				options = append(options, WithTagName("json"))
			}
			cl, err := Compare(tc.A, tc.B, options...)

			assert.Equal(t, tc.Error, err)
			assert.Equal(t, len(tc.Changes), len(cl))

			for i, c := range cl {
				assert.Equal(t, tc.Changes[i].Type, c.Type)
				assert.Equal(t, tc.Changes[i].Path, c.Path)
				assert.Equal(t, tc.Changes[i].From, c.From)
				assert.Equal(t, tc.Changes[i].To, c.To)
			}
		})
	}
}
