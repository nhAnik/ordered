# ordered
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
