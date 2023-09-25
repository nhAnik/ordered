package ordered

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

var dummy = struct{}{}

// Set represents an ordered set which is a special hashset keeping the
// insertion order intact. The insertion order is not changed if a element
// which already exists in the set is re-inserted.
type Set[T comparable] struct {
	mp *Map[T, struct{}]
}

// NewSet initializes an ordered set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		mp: NewMap[T, struct{}](),
	}
}

// NewSetWithElems initializes an ordered set and adds the elements
// in the set.
func NewSetWithElems[T comparable](elems ...T) *Set[T] {
	s := NewSet[T]()
	for _, elem := range elems {
		s.Add(elem)
	}
	return s
}

// Add inserts a new element in the set.
func (s *Set[T]) Add(elem T) {
	s.mp.Put(elem, dummy)
}

// Contains checks if the set contains the given element or not.
func (s *Set[T]) Contains(elem T) bool {
	return s.mp.ContainsKey(elem)
}

// Remove removes the given element from the set if the elements is
// already there in the set. The returned boolean value indicates
// whether the element is removed or not.
func (s *Set[T]) Remove(elem T) bool {
	if !s.Contains(elem) {
		return false
	}
	s.mp.Remove(elem)
	return true
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return s.mp.Len()
}

// Elements returns all the elements of the set according to their
// insertion order. The first element of the slice is the oldest
// element in the set.
func (s *Set[T]) Elements() []T {
	return s.mp.Keys()
}

// IsEmpty checks whether the set is empty or not.
func (s *Set[T]) IsEmpty() bool {
	return s.mp.IsEmpty()
}

// Clear removes all the elements from the set.
func (s *Set[T]) Clear() {
	s.mp.Clear()
}

// String returns the string representation of the set.
func (s *Set[T]) String() string {
	var sb strings.Builder
	sb.WriteString("set{")
	for idx, elem := range s.Elements() {
		if idx > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(elem))
	}
	sb.WriteByte('}')
	return sb.String()
}

// MarshalJSON implements json.Marshaler interface.
func (s Set[T]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for idx, elem := range s.Elements() {
		if idx > 0 {
			buf.WriteByte(',')
		}
		bytes, err := json.Marshal(elem)
		if err != nil {
			return nil, err
		}
		buf.Write(bytes)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (s *Set[T]) UnmarshalJSON(b []byte) error {
	if s.mp == nil {
		s.mp = NewMap[T, struct{}]()
	}
	unmarshalErrExists := false
	_, err := jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var elem T
		var elemBytes []byte
		if dataType == jsonparser.String {
			elemBytes = []byte(fmt.Sprintf("\"%s\"", string(value)))
		} else {
			elemBytes = value
		}
		if err := json.Unmarshal(elemBytes, &elem); err != nil {
			unmarshalErrExists = true
			return
		}
		s.Add(elem)
	})
	if err != nil {
		return err
	}
	if unmarshalErrExists {
		return errors.New("unmarshalling error")
	}
	return nil
}
