package gson

import (
	"encoding/json"
)

// UnmarshalJSON interface
func (j *JSON) UnmarshalJSON(b []byte) error {
	var v interface{}
	err := json.Unmarshal(b, &v)
	*j = JSON{v}
	return err
}

// Set by json path. It's a shortcut for Sets.
func (j *JSON) Set(path string, val interface{}) {
	j.Sets(val, parsePath(path)...)
}

// Sets element by path sections. If a section is not string or int, it will be ignored.
func (j *JSON) Sets(value interface{}, sections ...interface{}) {
	last := len(sections) - 1
	val := j.val
	var override func(interface{})

	if last == -1 {
		j.val = toJSONVal(value)
		return
	}

	for i, sect := range sections {
		switch k := sect.(type) {
		case int:
			arr, ok := val.([]interface{})
			if !ok || k >= len(arr) {
				nArr := make([]interface{}, k+1)
				copy(nArr, arr)
				arr = nArr
				override(arr)
			}
			if i == last {
				arr[k] = toJSONVal(value)
				return
			}
			val = arr[k]

			override = func(val interface{}) {
				arr[k] = val
			}
		case string:
			obj, ok := val.(map[string]interface{})
			if !ok {
				obj = map[string]interface{}{}
				override(obj)
			}
			if i == last {
				obj[k] = toJSONVal(value)
			}
			val = obj[k]

			override = func(val interface{}) {
				obj[k] = val
			}
		}
	}
}
