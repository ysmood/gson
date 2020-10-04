package gson

import (
	"encoding/json"
	"io"
	"reflect"
)

// New JSON from string, []byte, or io.Reader.
func New(v interface{}) (j JSON) {
	if s, ok := v.(string); ok {
		v = []byte(s)
	}
	j.value = v
	return
}

// UnmarshalJSON interface
func (j *JSON) UnmarshalJSON(b []byte) error {
	j.value = b
	return nil
}

// Val of the underlaying json value
func (j *JSON) Val() interface{} {
	for {
		val, ok := j.value.(JSON)
		if ok {
			*j = val
		} else {
			break
		}
	}

	var val interface{}
	switch v := j.value.(type) {
	case []byte:
		_ = json.Unmarshal(v, &val)
		j.value = val
	case io.Reader:
		_ = json.NewDecoder(v).Decode(&val)
		j.value = val
	}

	return j.value
}

// Set by json path. It's a shortcut for Sets.
func (j *JSON) Set(path string, val interface{}) *JSON {
	return j.Sets(val, Path(path)...)
}

// Sets element by path sections. If a section is not string or int, it will be ignored.
func (j *JSON) Sets(target interface{}, sections ...interface{}) *JSON {
	last := len(sections) - 1
	val := reflect.ValueOf(j.Val())
	var override func(reflect.Value)

	if last == -1 {
		j.value = target
		return j
	}

	for i, s := range sections {
		sect := reflect.ValueOf(s)
		if val.Kind() == reflect.Interface {
			val = val.Elem()
		}

		switch sect.Kind() {
		case reflect.Int:
			k := int(sect.Int())
			if val.Kind() != reflect.Slice || val.Len() <= k {
				nArr := reflect.ValueOf(make([]interface{}, k+1))
				if val.Kind() == reflect.Slice {
					reflect.Copy(nArr, val)
				}
				val = nArr
				override(val)
			}
			if i == last {
				val.Index(k).Set(reflect.ValueOf(target))
				return j
			}
			prev := val
			val = val.Index(k)
			override = func(v reflect.Value) {
				prev.Index(k).Set(v)
			}
		default:
			targetVal := reflect.ValueOf(target)
			if val.Kind() != reflect.Map {
				val = reflect.MakeMap(reflect.MapOf(sect.Type(), targetVal.Type()))
				override(val)
			}
			if i == last {
				val.SetMapIndex(sect, targetVal)
			}
			prev := val
			val = val.MapIndex(sect)
			override = func(v reflect.Value) {
				prev.SetMapIndex(sect, v)
			}
		}
	}
	return j
}
