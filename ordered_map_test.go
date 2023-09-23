package ordered_test

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/nhAnik/ordered"
	"github.com/stretchr/testify/assert"
)

type point struct{ x, y int }

type point3d struct{ X, Y, Z int }

func (p point3d) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%d-%d-%d", p.X, p.Y, p.Z)), nil
}

func (p *point3d) UnmarshalText(text []byte) error {
	split := strings.Split(string(text), "-")
	if len(split) != 3 {
		return errors.New("invalid text for point")
	}
	p.X, _ = strconv.Atoi(split[0])
	p.Y, _ = strconv.Atoi(split[1])
	p.Z, _ = strconv.Atoi(split[2])
	return nil
}

func TestNewMapWithElems(t *testing.T) {
	type kv = ordered.KeyValue[int, bool]
	om := ordered.NewMapWithElems[int, bool](kv{11, true}, kv{20, false}, kv{23, true})

	assert.True(t, om.ContainsKey(23))

	om.Put(99, false)
	assert.Equal(t, 4, om.Len())
	assert.Equal(t, []int{11, 20, 23, 99}, om.Keys())
}

func TestGet(t *testing.T) {

	t.Run("empty map get", func(t *testing.T) {
		om := ordered.NewMap[string, string]()

		_, ok := om.Get("foo")
		assert.False(t, ok)
	})

	t.Run("string string map", func(t *testing.T) {
		om := ordered.NewMap[string, string]()
		om.Put("foo", "bar")

		val, ok := om.Get("foo")
		assert.True(t, ok)
		assert.Equal(t, "bar", val)

		val, ok = om.Get("hello")
		assert.False(t, ok)
		assert.Equal(t, "", val)
	})

	t.Run("int pointer map", func(t *testing.T) {
		type myStruct struct{ str string }
		om := ordered.NewMap[int, *myStruct]()

		ms := &myStruct{str: "foo"}
		om.Put(1, ms)
		val, ok := om.Get(1)
		assert.True(t, ok)
		assert.Equal(t, ms, val)

		val, ok = om.Get(2)
		assert.False(t, ok)
		assert.Nil(t, val)
	})

	t.Run("update value of same key", func(t *testing.T) {
		om := ordered.NewMap[string, string]()
		om.Put("foo", "bar")

		val, ok := om.Get("foo")
		assert.True(t, ok)
		assert.Equal(t, "bar", val)

		om.Put("foo", "bak")
		val, ok = om.Get("foo")
		assert.True(t, ok)
		assert.Equal(t, "bak", val)

	})
}

func TestContainsKey(t *testing.T) {
	om := ordered.NewMap[string, string]()

	assert.False(t, om.ContainsKey("foo"))

	om.Put("foo", "bar")
	om.Put("abd", "def")
	om.Put("pqr", "xyz")
	assert.True(t, om.ContainsKey("foo"))
	assert.True(t, om.ContainsKey("abd"))
	assert.True(t, om.ContainsKey("pqr"))
	assert.False(t, om.ContainsKey("mno"))

	om.Remove("pqr")
	assert.False(t, om.ContainsKey("pqr"))
}

func TestRemove(t *testing.T) {
	om := ordered.NewMap[string, string]()

	om.Put("foo", "bar")
	om.Put("abd", "def")
	om.Put("pqr", "xyz")
	assert.True(t, om.ContainsKey("foo"))

	om.Remove("foo")
	assert.False(t, om.ContainsKey("foo"))
}

func TestLen(t *testing.T) {
	om := ordered.NewMap[string, string]()

	om.Put("foo", "bar")
	om.Put("abd", "def")
	assert.Equal(t, 2, om.Len())

	om.Put("abd", "pqr")
	assert.Equal(t, 2, om.Len())

	om.Remove("abc")
	assert.Equal(t, 2, om.Len())

	om.Remove("abd")
	assert.Equal(t, 1, om.Len())

	om.Clear()
	assert.Equal(t, 0, om.Len())
}

func TestKeys(t *testing.T) {
	t.Run("string string map", func(t *testing.T) {
		om := ordered.NewMap[string, string]()

		om.Put("foo", "bar")
		om.Put("abd", "def")
		assert.Equal(t, []string{"foo", "abd"}, om.Keys())

		om.Put("abd", "pqr")
		assert.Equal(t, []string{"foo", "abd"}, om.Keys())

		om.Put("abc", "abc")
		assert.Equal(t, []string{"foo", "abd", "abc"}, om.Keys())

		om.Remove("abd")
		assert.Equal(t, []string{"foo", "abc"}, om.Keys())

		om.Clear()
		assert.Equal(t, []string{}, om.Keys())
	})

	t.Run("struct string map", func(t *testing.T) {
		p1 := point{1, 10}
		p2 := point{2, 20}
		p3 := point{3, 30}
		p4 := point{4, 40}

		om := ordered.NewMap[point, string]()

		om.Put(p1, "p1")
		om.Put(p2, "p2")
		assert.Equal(t, []point{p1, p2}, om.Keys())

		om.Put(p3, "p3")
		assert.Equal(t, []point{p1, p2, p3}, om.Keys())

		om.Put(p2, "p22")
		assert.Equal(t, []point{p1, p2, p3}, om.Keys())

		om.Remove(p1)
		assert.Equal(t, []point{p2, p3}, om.Keys())

		om.Put(p4, "p4")
		assert.Equal(t, []point{p2, p3, p4}, om.Keys())

		om.Clear()
		assert.Equal(t, []point{}, om.Keys())
	})
}

func TestValues(t *testing.T) {
	om := ordered.NewMap[string, string]()

	om.Put("foo", "bar")
	om.Put("abd", "def")
	assert.Equal(t, []string{"bar", "def"}, om.Values())

	om.Put("abd", "pqr")
	assert.Equal(t, []string{"bar", "pqr"}, om.Values())

	om.Put("abc", "abc")
	assert.Equal(t, []string{"bar", "pqr", "abc"}, om.Values())

	om.Remove("abd")
	assert.Equal(t, []string{"bar", "abc"}, om.Values())

	om.Clear()
	assert.Equal(t, []string{}, om.Values())
}

func TestKeyValues(t *testing.T) {
	om := ordered.NewMap[string, int]()
	type kv = ordered.KeyValue[string, int]

	om.Put("foo", 10)
	om.Put("abd", 20)
	assert.Equal(t, []kv{{"foo", 10}, {"abd", 20}}, om.KeyValues())

	om.Put("abd", 15)
	assert.Equal(t, []kv{{"foo", 10}, {"abd", 15}}, om.KeyValues())

	om.Put("abc", 30)
	assert.Equal(t, []kv{{"foo", 10}, {"abd", 15}, {"abc", 30}}, om.KeyValues())

	om.Remove("abd")
	assert.Equal(t, []kv{{"foo", 10}, {"abc", 30}}, om.KeyValues())

	om.Clear()
	assert.Equal(t, []kv{}, om.KeyValues())
}

func TestIsEmpty(t *testing.T) {
	om := ordered.NewMap[string, any]()

	assert.True(t, om.IsEmpty())

	om.Put("hello", "world")
	assert.False(t, om.IsEmpty())
}

func TestClear(t *testing.T) {
	om := ordered.NewMap[string, string]()

	om.Put("foo", "bar")
	om.Put("abd", "def")
	assert.False(t, om.IsEmpty())

	om.Clear()
	assert.True(t, om.IsEmpty())
}

func TestString(t *testing.T) {
	t.Run("int bool map", func(t *testing.T) {
		type kv = ordered.KeyValue[int, bool]
		om := ordered.NewMapWithElems[int, bool](kv{11, true}, kv{20, false}, kv{23, true})

		assert.Equal(t, "map{11:true 20:false 23:true}", om.String())

		om.Remove(20)
		assert.Equal(t, "map{11:true 23:true}", om.String())

		om.Clear()
		assert.Equal(t, "map{}", om.String())
	})

	t.Run("string struct map", func(t *testing.T) {
		om := ordered.NewMap[string, point]()
		om.Put("p12", point{1, 2})
		om.Put("p34", point{3, 4})

		assert.Equal(t, "map{p12:{1 2} p34:{3 4}}", om.String())
	})
}

func TestMarshalJSON(t *testing.T) {
	t.Run("string any map", func(t *testing.T) {
		type dummy struct{ Elem string }
		type kv = ordered.KeyValue[string, any]
		om := ordered.NewMapWithElems[string, any](
			kv{"int", 1},
			kv{"float", 1.5},
			kv{"bool", true},
			kv{"string", "foo"},
			kv{"slice", []int{1, 2, 3}},
			kv{"struct", dummy{Elem: "bar"}},
		)

		bytes, err := om.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `{"int":1,"float":1.5,"bool":true,"string":"foo","slice":[1,2,3],"struct":{"Elem":"bar"}}`, string(bytes))

		om.Put("bool", false)
		om.Remove("slice")
		bytes, err = om.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `{"int":1,"float":1.5,"bool":false,"string":"foo","struct":{"Elem":"bar"}}`, string(bytes))

		om.Clear()
		bytes, err = om.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `{}`, string(bytes))
	})

	t.Run("struct string map", func(t *testing.T) {
		om := ordered.NewMap[point3d, string]()
		om.Put(point3d{1, 2, 3}, "p1")
		om.Put(point3d{4, 5, 6}, "p2")

		bytes, err := om.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `{"1-2-3":"p1","4-5-6":"p2"}`, string(bytes))
	})

	t.Run("struct string map with no text marshaler", func(t *testing.T) {
		type dummy struct{ Elem string }
		om := ordered.NewMap[dummy, string]()
		om.Put(dummy{"foo"}, "f")
		om.Put(dummy{"bar"}, "b")

		_, err := om.MarshalJSON()
		assert.Error(t, err)
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("string string map", func(t *testing.T) {
		om := ordered.NewMapWithElems[string, string]()
		data := []byte(`{"a":"apple","b":"bee","c":"cat","d":"deer"}`)

		err := om.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c", "d"}, om.Keys())
		assert.Equal(t, []string{"apple", "bee", "cat", "deer"}, om.Values())
	})

	t.Run("string slice map", func(t *testing.T) {
		om := ordered.NewMapWithElems[string, []int]()
		data := []byte(`{"a":[1,2],"b":[3,4],"c":[5,6,7]}`)

		err := om.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, om.Keys())
		assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6, 7}}, om.Values())
	})

	t.Run("struct string map", func(t *testing.T) {
		om := ordered.NewMap[point3d, string]()
		data := []byte(`{"1-2-3":"p1","4-5-6":"p2"}`)

		err := om.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, []point3d{{1, 2, 3}, {4, 5, 6}}, om.Keys())
		assert.Equal(t, []string{"p1", "p2"}, om.Values())
	})

	t.Run("unmarshal json with invalid key", func(t *testing.T) {
		om := ordered.NewMap[point, string]()
		data := []byte(`{"1-2":"p1","3-4":"p2"}`)

		err := om.UnmarshalJSON(data)
		assert.Error(t, err)
	})
}
