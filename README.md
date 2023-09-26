# ordered [![Go Reference](https://pkg.go.dev/badge/github.com/nhAnik/ordered.svg)](https://pkg.go.dev/github.com/nhAnik/ordered) ![Build](https://github.com/nhAnik/ordered/actions/workflows/build.yaml/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/nhAnik/ordered/badge.svg)](https://coveralls.io/github/nhAnik/ordered)
Implementation of generic ordered map and set. An ordered map is a special
hashmap which maintains the insertion order of the key-vale pair. The ordered
set is a wrapper around the ordered map which keeps the unique elements
according to their insertion order.

**Features:**
- Amortized `O(1)` time complexity for insertion, remove and get
- Supports generics
- Supports JSON marshalling and unmarshalling

**Limitations:**
- Not safe for concurrent use
- The map key and the set element must be `comparable`

## Usage

### Prerequisites
The go version should be >=1.18

### Installation
```
go get github.com/nhAnik/ordered
```

### Example
**Example of ordered map:**
```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/nhAnik/ordered"
)

type point struct{ X, Y int }

func main() {
	// Create a new generic ordered map
	om := ordered.NewMap[string, point]()

	// Put key-value pair in the map
	om.Put("p1", point{1, 2})
	om.Put("p2", point{2, 4})
	om.Put("p3", point{3, 6})

	if p2, ok := om.Get("p2"); ok {
		// Get value of p2
		fmt.Println(p2)
	}

	// Update value of p1
	// This will not affect the insertion order
	om.Put("p1", point{0, 0})

	// Iterate key-value pairs according to insertion order
	// 	p1 --> {0 0}
	// 	p2 --> {2 4}
	// 	p3 --> {3 6}
	for _, kv := range om.KeyValues() {
		fmt.Printf("%s --> %v\n", kv.Key, kv.Value)
	}

	// Remove p1 key and its mapped value
	om.Remove("p1")

	// Iterate values according to insertion order
	for _, v := range om.Values() {
		fmt.Println(v)
	}

	// Checks if the map is empty
	if om.IsEmpty() {
		fmt.Println("Empty map")
	}

	// Print string representation of the map
	fmt.Println(om) // map{p2:{2 4} p3:{3 6}}

	// Marshals the map to json according to order
	b, _ := json.Marshal(om)
	fmt.Println(string(b)) // {"p2":{"X":2,"Y":4},"p3":{"X":3,"Y":6}}
}
```

**Example of ordered set:**
```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/nhAnik/ordered"
)

func main() {
	// Create a new generic ordered set
	s := ordered.NewSet[string]()

	// Add new values in the set
	s.Add("C++")
	s.Add("Java")
	s.Add("Go")
	// Duplicate will not be added and will not
	// affect insertion order.
	s.Add("Java")

	// Check if an element exists
	if ok := s.Contains("Go"); ok {
		fmt.Println("Found Go")
	}

	// Iterate elements according to insertion order
	// 	C++
	// 	Java
	// 	Go
	for _, elem := range s.Elements() {
		fmt.Println(elem)
	}

	// Remove element if exists
	s.Remove("Java")

	// Checks if the set is empty
	if s.IsEmpty() {
		fmt.Println("Empty set")
	}

	// Print string representation of the set
	fmt.Println(s) // set{C++ Go}

	// Marshals the set to json according to order
	mp := ordered.NewMap[string, *ordered.Set[string]]()
	mp.Put("language", ordered.NewSetWithElems[string]("C++", "Go", "Python"))
	mp.Put("editor", ordered.NewSetWithElems[string]("VSCode", "Vim"))

	b, _ := json.Marshal(mp)
	fmt.Println(string(b)) // {"language":["C++","Go","Python"],"editor":["VSCode","Vim"]}
}
```
### Documentation
Documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/nhAnik/ordered#section-documentation).

## License

[MIT](LICENSE)
