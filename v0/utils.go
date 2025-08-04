package compare

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unsafe"
)

var isExportFlag uintptr = (1 << 5) | (1 << 6)

// getAsAny returns v's current value as any. It is equivalent to:
// var i any = (v's underlying value)
func getAsAny(v reflect.Value) any {
	// check if we can access the field
	// fake export it if it is unexported
	if !v.CanInterface() {
		flagTmp := (*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + 2*unsafe.Sizeof(uintptr(0))))
		*flagTmp = (*flagTmp) & (^isExportFlag)
	}
	return v.Interface()
}

// areKind checks if a and b areKind of the kinds listed
func areKind(a, b reflect.Value, kinds ...reflect.Kind) bool {
	var aMatch, bMatch bool

	for _, k := range kinds {
		if a.Kind() == k {
			aMatch = true
		}
		if b.Kind() == k {
			bMatch = true
		}
	}

	return aMatch && bMatch
}

func areType(a, b reflect.Value, types ...reflect.Type) bool {
	var aMatch, bMatch bool

	for _, t := range types {
		if a.Kind() != reflect.Invalid {
			if a.Type() == t {
				aMatch = true
			}
		}
		if b.Kind() != reflect.Invalid {
			if b.Type() == t {
				bMatch = true
			}
		}
	}

	return aMatch && bMatch
}

func isInvalid(a, b reflect.Value) bool {
	if a.Kind() == b.Kind() {
		return false
	}

	if a.Kind() == reflect.Invalid {
		return false
	}
	if b.Kind() == reflect.Invalid {
		return false
	}

	return true
}

func copyAppend(src []string, elems ...string) []string {
	dst := make([]string, len(src)+len(elems))
	copy(dst, src)
	for i := len(src); i < len(src)+len(elems); i++ {
		dst[i] = elems[i-len(src)]
	}
	return dst
}

func getTagName(tag string, f reflect.StructField) string {
	t := f.Tag.Get(tag)

	parts := strings.Split(t, ",")
	if len(parts) < 1 {
		return "-"
	}

	return parts[0]
}

func getIdentifier(tag string, v reflect.Value, joinSep string) any {
	if v.Kind() != reflect.Struct {
		return nil
	}

	var combinedIdentifierTemplate string
	var combinedIdentifier map[string]reflect.Value = make(map[string]reflect.Value)

	for i := 0; i < v.NumField(); i++ {
		if hto, toValue := hasTagOption(tag, v.Type().Field(i), "identifier"); hto {
			combinedIdentifier[v.Type().Field(i).Name] = v.Field(i)
			if toValue != "" {
				if combinedIdentifierTemplate == "" {
					combinedIdentifierTemplate = toValue
				} else if combinedIdentifierTemplate != toValue {
					panic("identifier name must be identical")
				}
			}
		}
	}

	switch len(combinedIdentifier) {
	case 0:
		return nil
	case 1:
		for identifier := range maps.Values(combinedIdentifier) {
			return identifier.Interface()
		}
		return nil
	default:
		if combinedIdentifierTemplate == "" {
			var combinedID []string
			notCombinable := []reflect.Kind{
				reflect.Struct, reflect.Slice, reflect.Array,
				reflect.Map, reflect.Ptr, reflect.Interface, reflect.Invalid,
			}

			for _, idVal := range combinedIdentifier {
				res := areKind(idVal, idVal, notCombinable...)
				if res {
					panic(ErrNotCombinableIdentifier)
				}
				combinedID = append(combinedID, strings.Trim(fmt.Sprintf("%#v", idVal), `"`))
			}
			return strings.Join(combinedID, joinSep)
		} else {
			templatedIdentifier := template.Must(template.New("id").Parse(combinedIdentifierTemplate))
			templatedIdentifierOutput := bytes.NewBuffer(nil)

			if err := templatedIdentifier.Execute(templatedIdentifierOutput, combinedIdentifier); err != nil {
				panic("failed to execute template: " + err.Error())
			}
			return templatedIdentifierOutput.String()
		}
	}
}

func hasTagOption(tag string, f reflect.StructField, opt string) (bool, string) {
	parts := strings.Split(f.Tag.Get(tag), ",")
	if len(parts) < 2 {
		return false, ""
	}

	for _, option := range parts[1:] {
		tagOption := strings.Split(option, ":")
		if len(tagOption) == 0 || tagOption[0] != opt {
			continue
		}
		if len(tagOption) == 2 {
			return true, tagOption[1]
		}
		return true, ""
	}

	return false, ""
}

func getFinalValue(t reflect.Value) reflect.Value {
	switch t.Kind() {
	case reflect.Interface:
		return getFinalValue(t.Elem())
	case reflect.Ptr:
		return getFinalValue(reflect.Indirect(t))
	default:
		return t
	}
}

// hasAtSameIndex checks if v is at idx in s
func hasAtSameIndex(s, v reflect.Value, idx int) bool {
	if idx < s.Len() {
		x := s.Index(idx)
		return reflect.DeepEqual(getAsAny(x), getAsAny(v))
	}
	return false
}

func patchChange(t DiffType, d Difference) Difference {
	nd := Difference{
		Type: t,
		Path: d.Path,
	}

	switch t {
	case ADD:
		nd.To = d.To
	case REMOVE:
		nd.From = d.To
	}

	return nd
}

// Compare returns a changelog of all mutated values from both
func Compare(left, right any, opts ...OptsFunc) (Differences, error) {
	c, err := NewComparer(opts...)
	if err != nil {
		return nil, err
	}
	return c.Compare(left, right)
}

func idComplex(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		bWriter := new(bytes.Buffer)
		if err := gob.NewEncoder(bWriter).Encode(v); err != nil {
			panic(err)
		}
		return string(bWriter.Bytes())
	}

}

func idString(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		return fmt.Sprint(v)
	}
}

func pathMatches(filter, path []string) bool {
	for i, f := range filter {
		if len(path) < i+1 {
			return false
		}

		matched, _ := regexp.MatchString(f, path[i])
		if !matched {
			return false
		}
	}

	return true
}
