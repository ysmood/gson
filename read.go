package gson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// JSON represent a JSON value
type JSON struct {
	value interface{}
}

// MarshalJSON interface
func (j JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Val())
}

// JSON string
func (j JSON) JSON(prefix, indent string) string {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent(prefix, indent)
	_ = enc.Encode(j.Val())
	return buf.String()
}

// String implements fmt.Stringer interface
func (j JSON) String() string {
	return fmt.Sprintf("%v", j.Val())
}

// Get by json path. It's a shortcut for Gets.
func (j JSON) Get(path string) JSON {
	j, _ = j.Gets(Path(path)...)
	return j
}

// Has an element is found on the path
func (j JSON) Has(path string) bool {
	_, has := j.Gets(Path(path)...)
	return has
}

// Gets element by path sections. If a section is not string or int, it will be ignored.
// The last return value will be false if not found.
func (j JSON) Gets(sections ...interface{}) (JSON, bool) {
	for _, sect := range sections {
		val, has := get(reflect.ValueOf(j.Val()), sect)
		if !has {
			return JSON{}, false
		}
		j.value = val
	}
	return j, true
}

func get(objVal reflect.Value, sect interface{}) (val interface{}, has bool) {
	switch k := sect.(type) {
	case int:
		if objVal.Kind() != reflect.Slice || k >= objVal.Len() {
			return
		}

		has = true
		val = objVal.Index(k).Interface()

	default:
		sectVal := reflect.ValueOf(sect)

		if objVal.Kind() != reflect.Map || !sectVal.Type().AssignableTo(objVal.Type().Key()) {
			return
		}

		v := objVal.MapIndex(sectVal)
		if !v.IsValid() {
			return
		}

		has = true
		val = v.Interface()
	}

	return
}

// Str value
func (j JSON) Str() string {
	v := j.Val()
	if v, ok := v.(string); ok {
		return v
	}
	return fmt.Sprintf("%v", v)
}

var floatType = reflect.TypeOf(.0)

// Num value
func (j JSON) Num() float64 {
	v := reflect.ValueOf(j.Val())
	if v.Type().ConvertibleTo(floatType) {
		return v.Convert(floatType).Float()
	}
	return 0
}

// Bool value
func (j JSON) Bool() bool {
	if v, ok := j.Val().(bool); ok {
		return v
	}
	return false
}

// Nil or not
func (j JSON) Nil() bool {
	return j.Val() == nil
}

var intType = reflect.TypeOf(0)

// Int value
func (j JSON) Int() int {
	v := reflect.ValueOf(j.Val())
	if v.Type().ConvertibleTo(intType) {
		return int(v.Convert(intType).Int())
	}
	return 0
}

// Map of JSON
func (j JSON) Map() map[string]JSON {
	if v, ok := j.Val().(map[string]interface{}); ok {
		obj := make(map[string]JSON, len(v))
		for k, el := range v {
			obj[k] = JSON{el}
		}
		return obj
	}

	return make(map[string]JSON, 0)
}

// Arr of JSON
func (j JSON) Arr() []JSON {
	if v, ok := j.Val().([]interface{}); ok {
		l := len(v)
		arr := make([]JSON, l)
		for i := 0; i < l; i++ {
			arr[i] = JSON{v[i]}
		}
		return arr
	}

	return make([]JSON, 0)
}

// Join elements
func (j JSON) Join(sep string) string {
	list := []string{}

	for _, el := range j.Arr() {
		list = append(list, el.Str())
	}

	return strings.Join(list, sep)
}

var regIndex = regexp.MustCompile(`^0|([1-9]\d*)$`)

// Path from string
func Path(path string) []interface{} {
	list := strings.Split(path, ".")
	sects := make([]interface{}, len(list))
	for i, s := range list {
		if regIndex.MatchString(s) {
			index, err := strconv.ParseInt(s, 10, 64)
			if err == nil {
				sects[i] = int(index)
				continue
			}
		}
		sects[i] = s
	}
	return sects
}
