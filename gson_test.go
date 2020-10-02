package gson_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/ysmood/gson"
)

func ExampleJSON() {
	var obj gson.JSON
	_ = json.Unmarshal([]byte(`{"a": {"b": [1, 2]}}`), &obj)

	fmt.Println(obj.Get("a.b.0").Int())

	obj.Set("a.b.1", "ok")
	fmt.Println(obj)

	// Output:
	// 1
	// map[a:map[b:[1 ok]]]
}

func Test(t *testing.T) {
	eq := genEq(t)

	b, _ := n(`"ok"`).MarshalJSON()
	eq(string(b), `"ok"`)

	eq(n(`"ok"`).Val(), "ok")
	eq(n(`"ok"`).Str(), "ok")
	eq(n(`1`).Str(), "1")
	eq(n(`1.2`).Num(), 1.2)
	eq(n(`"ok"`).Num(), 0.0)
	eq(n(`1`).Int(), 1)
	eq(n(`"ok"`).Int(), 0)
	eq(n(`true`).Bool(), true)
	eq(n(`1`).Bool(), false)
	eq(n(`null`).Nil(), true)

	eq(n(`1`).Join(""), "")

	j := n(`{
		"a": {
			"b": 1
		},
		"c": ["x", "y", "z"]
	}`)

	eq(j.Get("a.b").Int(), 1)
	eq(j.Get("c.1").Str(), "y")

	eq(j.Get("c").Arr()[1].Str(), "y")
	eq(j.Get("c").Join(" "), "x y z")
	eq(j.Get("a").Map()["b"].Int(), 1)
	eq(len(j.Get("c").Map()), 0)

	eq(j.Has("a.b"), true)
	eq(j.Has("a.x"), false)
	eq(j.Has("c.10"), false)

	self := gson.JSON{}
	self.Sets("ok")
	eq(self.Str(), "ok")
	self.Sets(map[string]interface{}{"a": 1})
	eq(self.Get("a").Int(), 1)
	self.Sets([]interface{}{1})
	eq(self.Get("0").Int(), 1)

	j.Sets(2.0, "a", "b")
	eq(j.Get("a.b").Int(), 2)

	j.Set("c.1", 2)
	eq(j.Get("c.1").Int(), 2)

	j.Sets(3, "a", "b")
	eq(j.Get("a.b").Int(), 3)

	eq(j.Get("s.10.b").Nil(), true)
	eq(j.Get("c.10").Nil(), true)

	j.Set("s.1.a", 10)
	j.Set("c.5", "ok")
	eq(fmt.Sprint(j), `map[a:map[b:3] c:[x 2 z <nil> <nil> ok] s:[<nil> map[a:10]]]`)
}

func TestLab(t *testing.T) {
	eq := genEq(t)
	j := n(`{
		"a": {
			"b": 1
		},
		"c": ["x", "y", "z"]
	}`)

	eq(j.Get("c.0").Str(), "x")
}

func n(s string) gson.JSON {
	var j gson.JSON
	_ = json.Unmarshal([]byte(s), &j)
	return j
}

func genEq(t *testing.T) func(a, b interface{}) {
	return func(a, b interface{}) {
		t.Helper()
		if !reflect.DeepEqual(a, b) {
			t.Log(a, "!=", b)
			t.Fail()
		}
	}
}
