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

type Set[T comparable] struct {
	mp *Map[T, struct{}]
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		mp: NewMap[T, struct{}](),
	}
}

func NewSetWithElems[T comparable](elems ...T) *Set[T] {
	s := NewSet[T]()
	for _, elem := range elems {
		s.Add(elem)
	}
	return s
}

func (s *Set[T]) Add(elem T) {
	s.mp.Put(elem, dummy)
}

func (s *Set[T]) Contains(elem T) bool {
	return s.mp.ContainsKey(elem)
}

func (s *Set[T]) Remove(elem T) {
	s.mp.Remove(elem)
}

func (s *Set[T]) Len() int {
	return s.mp.Len()
}

func (s *Set[T]) Elements() []T {
	return s.mp.Keys()
}

func (s *Set[T]) IsEmpty() bool {
	return s.mp.IsEmpty()
}

func (s *Set[T]) Clear() {
	s.mp.Clear()
}

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
