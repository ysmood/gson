package gson

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// JSON represent a JSON value
type JSON struct {
	val interface{}
}

// MarshalJSON interface
func (j JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.val)
}

// Get by json path. It's a shortcut for Gets.
func (j JSON) Get(path string) JSON {
	j, _ = j.Gets(parsePath(path)...)
	return j
}

// Gets element by path sections. If a section is not string or int, it will be ignored.
func (j JSON) Gets(sections ...interface{}) (JSON, bool) {
	for _, sect := range sections {
		switch k := sect.(type) {
		case int:
			v, _ := j.val.([]interface{})
			if k >= len(v) {
				return JSON{}, false
			}
			j = JSON{v[k]}
		case string:
			v, _ := j.val.(map[string]interface{})
			if _, has := v[k]; !has {
				return JSON{}, false
			}
			j = JSON{v[k]}
		}
	}
	return j, true
}

// Val of the underlaying json value
func (j JSON) Val() interface{} {
	return j.val
}

// Str value
func (j JSON) Str() string {
	if v, ok := j.val.(string); ok {
		return v
	}
	return fmt.Sprintf("%v", j.val)
}

// Num value
func (j JSON) Num() float64 {
	if v, ok := j.val.(float64); ok {
		return v
	}
	return 0
}

// Bool value
func (j JSON) Bool() bool {
	if v, ok := j.val.(bool); ok {
		return v
	}
	return false
}

// Nil or not
func (j JSON) Nil() bool {
	return j.val == nil
}

// Int value
func (j JSON) Int() int {
	return int(j.Num())
}

// Map of JSON
func (j JSON) Map() map[string]JSON {
	if v, ok := j.val.(map[string]interface{}); ok {
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
	if v, ok := j.val.([]interface{}); ok {
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

// String implements fmt.Stringer interface
func (j JSON) String() string {
	return fmt.Sprintf("%v", j.val)
}

func toJSONVal(val interface{}) interface{} {
	switch v := val.(type) {
	case float64:
	case string:
	case bool:
	case nil:
	case []interface{}:
		l := len(v)
		for i := 0; i < l; i++ {
			v[i] = toJSONVal(v[i])
		}
	case map[string]interface{}:
		for k, el := range v {
			v[k] = toJSONVal(el)
		}
	default:
		b, _ := json.Marshal(val)
		var n interface{}
		_ = json.Unmarshal(b, &n)
		return n
	}

	return val
}

var regIndex = regexp.MustCompile(`^0|([1-9]\d*)$`)

func parsePath(path string) []interface{} {
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
