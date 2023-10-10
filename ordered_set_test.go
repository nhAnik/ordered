package ordered_test

import (
	"encoding/json"
	"testing"

	"github.com/nhAnik/ordered"
	"github.com/stretchr/testify/assert"
)

func TestNewSetWithCapacity(t *testing.T) {
	s := ordered.NewSetWithCapacity[string](3)

	assert.True(t, s.IsEmpty())

	s.Add("abc")
	s.Add("def")
	s.Add("abc")
	s.Add("pqr")
	assert.Equal(t, 3, s.Len())
	assert.Equal(t, []string{"abc", "def", "pqr"}, s.Elements())
}

func TestNewSetWithElems(t *testing.T) {
	s := ordered.NewSetWithElems[string]("foo", "bar", "foo", "baz")

	assert.Equal(t, 3, s.Len())
	assert.Equal(t, []string{"foo", "bar", "baz"}, s.Elements())
}

func TestAdd(t *testing.T) {
	s := ordered.NewSet[string]()

	s.Add("foo")
	s.Add("bar")
	s.Add("foo")
	assert.Equal(t, []string{"foo", "bar"}, s.Elements())
}

func TestContains(t *testing.T) {
	s := ordered.NewSetWithElems[string]("foo", "bar", "foo", "baz")

	assert.True(t, s.Contains("foo"))
	assert.True(t, s.Contains("bar"))
	assert.True(t, s.Contains("baz"))
	assert.False(t, s.Contains("abc"))
}

func TestSetRemove(t *testing.T) {
	s := ordered.NewSetWithElems[string]("foo", "bar", "foo", "baz")

	assert.Equal(t, []string{"foo", "bar", "baz"}, s.Elements())

	removed := s.Remove("bar")
	assert.Equal(t, []string{"foo", "baz"}, s.Elements())
	assert.True(t, removed)

	removed = s.Remove("abc")
	assert.Equal(t, []string{"foo", "baz"}, s.Elements())
	assert.False(t, removed)
}

func TestSetLen(t *testing.T) {
	s := ordered.NewSetWithElems[string]("foo", "bar", "foo", "baz")

	assert.Equal(t, 3, s.Len())

	s.Remove("bar")
	assert.Equal(t, 2, s.Len())

	s.Clear()
	assert.Equal(t, 0, s.Len())
}

func TestSetElements(t *testing.T) {
	s := ordered.NewSetWithElems[string]("foo", "bar", "foo", "baz")

	assert.Equal(t, []string{"foo", "bar", "baz"}, s.Elements())

	s.Add("xyz")
	assert.Equal(t, []string{"foo", "bar", "baz", "xyz"}, s.Elements())

	s.Remove("foo")
	assert.Equal(t, []string{"bar", "baz", "xyz"}, s.Elements())
}

func TestSetIsEmpty(t *testing.T) {
	s := ordered.NewSet[int]()

	assert.True(t, s.IsEmpty())

	s.Add(1)
	assert.False(t, s.IsEmpty())
}

func TestSetClear(t *testing.T) {
	s := ordered.NewSetWithElems[string]("foo", "bar", "foo", "baz")

	assert.False(t, s.IsEmpty())

	s.Clear()
	assert.True(t, s.IsEmpty())
}

func TestSetString(t *testing.T) {
	t.Run("set of string", func(t *testing.T) {
		s := ordered.NewSetWithElems[string]("abc", "def", "abc", "xyz")

		assert.Equal(t, "set{abc def xyz}", s.String())

		s.Remove("abc")
		assert.Equal(t, "set{def xyz}", s.String())

		s.Clear()
		assert.Equal(t, "set{}", s.String())
	})

	t.Run("set of struct", func(t *testing.T) {
		s := ordered.NewSetWithElems[point](point{1, 2}, point{2, 4}, point{1, 2})

		assert.Equal(t, "set{{1 2} {2 4}}", s.String())
	})

	t.Run("set of integers", func(t *testing.T) {
		s := ordered.NewSetWithElems[int](1, 2, 2, 1, 4, 2, 3, 1, 4)

		assert.Equal(t, "set{1 2 4 3}", s.String())
	})
}

func TestSetMarshalJSON(t *testing.T) {
	t.Run("set of string", func(t *testing.T) {
		s := ordered.NewSetWithElems[string]("abc", "def", "abc", "xyz")

		bytes, err := s.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `["abc","def","xyz"]`, string(bytes))

		s.Remove("abc")
		bytes, err = s.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `["def","xyz"]`, string(bytes))

		s.Clear()
		bytes, err = s.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "[]", string(bytes))
	})

	t.Run("set of struct", func(t *testing.T) {
		type st struct{ Val int }
		s := ordered.NewSetWithElems[st](st{1}, st{10}, st{100})

		bytes, err := s.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `[{"Val":1},{"Val":10},{"Val":100}]`, string(bytes))
	})

	t.Run("set of struct with text marshaller", func(t *testing.T) {
		s := ordered.NewSetWithElems[point3d](point3d{1, 2, 3}, point3d{4, 5, 6}, point3d{7, 8, 9})

		bytes, err := s.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `["1-2-3","4-5-6","7-8-9"]`, string(bytes))
	})

	t.Run("set of integers", func(t *testing.T) {
		s := ordered.NewSetWithElems[int](1, 2, 2, 1, 4, 2, 3, 1, 4)

		bytes, err := s.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `[1,2,4,3]`, string(bytes))
	})

	t.Run("pointer to set inside struct", func(t *testing.T) {
		data := struct {
			Name string
			Val  int
			Tags *ordered.Set[string]
		}{
			Name: "xyz",
			Val:  10,
			Tags: ordered.NewSetWithElems[string]("foo", "bar", "baz"),
		}

		bytes, err := json.Marshal(data)
		assert.NoError(t, err)
		assert.Equal(t, `{"Name":"xyz","Val":10,"Tags":["foo","bar","baz"]}`, string(bytes))
	})

	t.Run("set inside struct", func(t *testing.T) {
		data := struct {
			Name string
			Val  int
			Tags ordered.Set[string]
		}{
			Name: "xyz",
			Val:  10,
			Tags: *ordered.NewSetWithElems[string]("foo", "bar", "baz"),
		}

		bytes, err := json.Marshal(data)
		assert.NoError(t, err)
		assert.Equal(t, `{"Name":"xyz","Val":10,"Tags":["foo","bar","baz"]}`, string(bytes))
	})

	t.Run("element marshalling error", func(t *testing.T) {
		// complex type is not supported for json marshalling
		s := ordered.NewSet[complex128]()
		s.Add(2 + 3i)

		_, err := json.Marshal(s)
		assert.Error(t, err)
	})
}

func TestSetUnmarshalJSON(t *testing.T) {
	t.Run("set of string", func(t *testing.T) {
		s := ordered.NewSet[string]()
		data := []byte(`["abc","def","xyz", "abc"]`)

		err := s.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, []string{"abc", "def", "xyz"}, s.Elements())
	})

	t.Run("set of struct", func(t *testing.T) {
		type st struct{ Val int }
		s := ordered.NewSet[st]()
		data := []byte(`[{"Val":1},{"Val":10},{"Val":100}]`)

		err := s.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, []st{{1}, {10}, {100}}, s.Elements())
	})

	t.Run("set of struct with text marshaller", func(t *testing.T) {
		s := ordered.NewSet[point3d]()
		data := []byte(`["1-2-3","4-5-6","7-8-9"]`)

		err := s.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, []point3d{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, s.Elements())
	})

	t.Run("set inside struct", func(t *testing.T) {
		type st struct {
			Name string
			Val  int
			Tags *ordered.Set[string]
		}
		s := &st{}
		data := []byte(`{"name":"xyz","val":10,"tags":["foo","bar","baz"]}`)

		err := json.Unmarshal(data, s)
		assert.NoError(t, err)
		assert.Equal(t, []string{"foo", "bar", "baz"}, s.Tags.Elements())
	})

	t.Run("set inside map", func(t *testing.T) {
		mp := make(map[string]*ordered.Set[int])
		data := []byte(`{"two":[4,8],"five":[15,35],"seven":[21,35,49]}`)

		err := json.Unmarshal(data, &mp)
		assert.NoError(t, err)
		assert.Equal(t, []int{4, 8}, mp["two"].Elements())
		assert.Equal(t, []int{15, 35}, mp["five"].Elements())
		assert.Equal(t, []int{21, 35, 49}, mp["seven"].Elements())
	})

	t.Run("unmarshal error", func(t *testing.T) {
		s := ordered.NewSet[string]()
		data := []byte(`["foo"`)

		err := s.UnmarshalJSON(data)
		assert.Error(t, err)
	})

	t.Run("unmarshal error", func(t *testing.T) {
		s := ordered.NewSet[point]()
		data := []byte(`["1-2","3-4"]`)

		err := s.UnmarshalJSON(data)
		assert.Error(t, err)
	})
}
