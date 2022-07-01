# Internals of **Go** maps

There will be information about internal implementation of Go map structure

Source code of the map struct is [here](https://go.dev/src/runtime/map.go)

- [Preface](#preface)
- [Implementation in Go](#implementation-in-go)
	- []()
- [Conclusion](#conclusion)

## Preface

Map structures are just associative containers mapping keys to values. Almost all the programming languages have these data structures: `map` in Go, `dict` in Python, `HashMap` in Rust, `array` in PHP, etc. In most cases interaction with maps (like insertion, removing and lookup) should be completed in constant time on average (***O(1)*** in other words).

## Implementation in **Go**

Each map in Go represented by this structure

```go
// A header for a Go map.
type hmap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/reflectdata/reflect.go.
	// Make sure this stays in sync with the compiler's definition.

	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra *mapextra // optional fields
}
```

### Hash functions

To achieve *"almost constant time performance"* the hash functions are used. The idea here is to create subsets of the entire collection by categorizing them based on the hash (or part of the hash) of the keys. This subsets are called **[buckets](#buckets)**.

Hashes helps distinguish pretty similar keys from each other. For example:

```
key1 = "hhhhhha"
key2 = "hhhhhhb"
key1md5 -> "43e8c4197f1e7a282af4c971311be358"
key2md5 -> "8cd072867c6a9fe864d2840190fbad24"
```

only one character is differs in this keys, but their hashes are very different (even using md5, which is not secure anymore, but for the demonstration purposes it doesn't matter)

### Buckets

In Go maps the data is arranged into an array of buckets. Each bucket contains up to 8 key/elem pairs.

```go
// A bucket for a Go map.
type bmap struct {
	// tophash generally contains the top byte of the hash value
	// for each key in this bucket. If tophash[0] < minTopHash,
	// tophash[0] is a bucket evacuation state instead.

	tophash [bucketCnt]uint8
	// Followed by bucketCnt keys and then bucketCnt elems.
	// NOTE: packing all the keys together and then all the elems together makes the
	// code a bit more complicated than alternating key/elem/key/elem/... but it allows
	// us to eliminate padding which would be needed for, e.g., map[int64]int8.
	// Followed by an overflow pointer.

}
```
As documentation says, the first elements in bucket are keys and after them are the values

#### Bucket selection

Low order bits are used in specific bucket selection. This bits are calculated as:
```
hash(key) % len(buckets)
```
> But to speed up calculation, the hash is represented in it binary form and hmap structure contains the **B** field, which comment is self-descriptive. So, if there are 4 buckets B is equal to 2 (because log2(4) is 2) and last two least significant bits of the hash are used! 🤯

#### Bucket overflow

As was mentioned earlier, bucket contains only 8 k/v pairs. If you need to add new k/v pair in bucket, which already full, then new bucket is created, and existing bucket contain pointer to this newly created bucket. Overflows are undesirable, since when accessing the map, it interacts with several buckets at once.

```go
// mapextra holds fields that are not present on all maps.

type mapextra struct {

	// If both key and elem do not contain pointers and are inline, then we mark bucket
	// type as containing no pointers. This avoids scanning such maps.
	// However, bmap.overflow is a pointer. In order to keep overflow buckets
	// alive, we store pointers to all overflow buckets in hmap.extra.overflow and hmap.extra.oldoverflow.
	// overflow and oldoverflow are only used if key and elem do not contain pointers.
	// overflow contains overflow buckets for hmap.buckets.
	// oldoverflow contains overflow buckets for hmap.oldbuckets.
	// The indirection allows to store a pointer to the slice in hiter.
	overflow    *[]*bmap
	oldoverflow *[]*bmap


	// nextOverflow holds a pointer to a free overflow bucket.
	nextOverflow *bmap

}
```

### Map grow and Evacuation

Map grows based on load factor. As said in source code `Maximum average load of a bucket that triggers growth is 6.5.`, which means, that if each bucket inside map contain 6.5 elements, the map will grow.

When the load factor reached, data evacuation begins. **Data evacuation**  is the process of creating new list of buckets and copying data from the old buckets into new ones.

> Note, that map grows creating x2 buckets. This process is pretty slow, dute to memory allocation and data copying.

### Map creation

When the next line of code is called
```go
m := make(map[string]string)
```
go runtime creates new pointer to hmap structure and initialize it. Internally call to `make(...)` for maps translates into `makemap(...)` function which signature looks like
```go
func makemap(t *maptype, hint int, h *hmap) *hmap
```
But `var m map[int]int` doesn't initialize map, and if we try to assign a value to the key, we will have panic `panic: assignment to entry in nil map`

### Retrieving map value

Each bucket contain 8 slots for most significant bits of hash (***topHash*** field of ***bmap*** structure). It allows to speed up the check of key exsistance in given bucket.

```go
// mapaccess1 returns a pointer to h[key].  Never returns nil, instead
// it will return a reference to the zero object for the elem type if
// the key is not in the map.
// NOTE: The returned pointer may keep the whole map live, so don't
// hold onto it for very long.
func mapaccess1(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer
func mapaccess2(t *maptype, h *hmap, key unsafe.Pointer) (unsafe.Pointer, bool)
```

## Notes

### Maps are passed to function by value.

A map in Go is simply a pointer to an hmap structure. This is the answer to the question why, despite the fact that the map is passed to the function by value, the values themselves that lie in it change - it's all about the pointer.

```go

package main

func foo(m map[int]int) { 
    m[10] = 10 
}

func main() {
	m := make(map[int]int)
	m[10] = 15
	println("m[10] before foo =", m[10])
	foo(m)
	println("m[10] after foo =", m[10])
}
```
It will result as follows:
```
m[10] before foo = 15
m[10] after foo = 10
```
But if we try to initialize declared map
```go
package main

import "fmt"

func fn(m map[int]int) {
	m = make(map[int]int)
	fmt.Println("m == nil in fn?:", m == nil)
}

func main() {
	var m map[int]int
	fn(m)
	fmt.Println("m == nil in main?:", m == nil)
}
```
It will have the following result
```
m == nil in fn?: false
m == nil in main?: true
```
***Why?***

A map in Go is simply a pointer to an hmap structure. This is the answer to the question why, despite the fact that the map is passed to the function by value, the values themselves that lie in it change - it's all about the pointer.

### Better to Specify map size when creating

If size of map is known, it's better to initialize map with size, because it prevents memory reallocations (*[data evacuation](#map-grow-and-evacuation)*)
```go
m := make(map[int]string, 5)
m[1] = "I"
m[2] = "II"
m[3] = "III"
m[4] = "IV"
m[5] = "V"
```

### Impossibility of *`&m[k]`*

It is impossible to take a reference of map value due to the evacuation process, that will invalidate the address of the reference

## Conclusion
*TODO*
