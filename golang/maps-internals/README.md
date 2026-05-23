# Internals of **Go** maps

There will be information about internal implementation of Go map structure

Source code of the map struct is [here](https://go.dev/src/runtime/map.go)

- [Preface](#preface)
- [Implementation in Go](#implementation-in-go)
	- [Hash functions](#hash-functions)
	- [Buckets](#buckets)
		- [Bucket selection](#bucket-selection)
		- [Bucket overflow](#bucket-overflow)
	- [Map grow and Evacuation](#map-grow-and-evacuation)
	- [Map creation](#map-creation)
	- [Retrieving map value](#retrieving-map-value)
- [Notes](#notes)
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

#### Bucket selection (***Simplyfied***)
To select a bucket, firstly we need to calculate bucket mask.

```go
// bucketShift returns 1<<b, optimized for code generation.
func bucketShift(b uint8) uintptr {
    // Masking the shift amount allows overflow checks to be elided.
    return uintptr(1) << (b & (goarch.PtrSize*8 - 1))
}

// bucketMask returns 1<<b - 1, optimized for code generation.
func bucketMask(b uint8) uintptr {
    return bucketShift(b) - 1
}
```
*b uint8* - is a `B` field in `hmap` structure (log_2 of number of buckets)

After that the number of the bucket can be calculated. For this purpose the "bit AND" is used:
```go
m := bucketMask(h.B)
b := (*bmap)(add(h.buckets, (hash&m)*uintptr(t.bucketsize)))
```
Let's decompose this:
- *h.buckets* - unsafe pointer to the buckets of given hmap
- *hash* - hash of the map key
- *m* - bucket mask
- *b* - calculated pointer to the bucket

So, algorithm is following:
1. Take the hash of the key
2. Calculate bucket mask
3. Use pointer arighmetics: add to the pointer to the buckets an "shift", where the shift is ANDed hash and bucket mask

> To speed up calculation, the hash is represented in it binary form and hmap structure contains the **B** field, which comment is self-descriptive. So, if there are 4 buckets B is equal to 2 (because log2(4) is 2) and last two least significant bits of the hash are used! 🤯

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

## Changes in Go implementation v1.24

Since Go v1.24 there is [changes in the map implementation](https://go.dev/doc/go1.24#runtime). Now its uses swiss table.

Now, instead of `hmap` there is `Map` structure:

```go
type Map struct {
	// The number of filled slots (i.e. the number of elements in all
	// tables). Excludes deleted slots.
	// Must be first (known by the compiler, for len() builtin).
	used uint64

	// seed is the hash seed, computed as a unique random number per map.
	seed uintptr

	// The directory of tables.
	//
	// Normally dirPtr points to an array of table pointers
	//
	// dirPtr *[dirLen]*table
	//
	// The length (dirLen) of this array is `1 << globalDepth`. Multiple
	// entries may point to the same table. See top-level comment for more
	// details.
	//
	// Small map optimization: if the map always contained
	// abi.MapGroupSlots or fewer entries, it fits entirely in a
	// single group. In that case dirPtr points directly to a single group.
	//
	// dirPtr *group
	//
	// In this case, dirLen is 0. used counts the number of used slots in
	// the group. Note that small maps never have deleted slots (as there
	// is no probe sequence to maintain).
	dirPtr unsafe.Pointer
	dirLen int

	// The number of bits to use in table directory lookups.
	globalDepth uint8

	// The number of bits to shift out of the hash for directory lookups.
	// On 64-bit systems, this is 64 - globalDepth.
	globalShift uint8

	// writing is a flag that is toggled (XOR 1) while the map is being
	// written. Normally it is set to 1 when writing, but if there are
	// multiple concurrent writers, then toggling increases the probability
	// that both sides will detect the race.
	writing uint8

	// tombstonePossible is false if we know that no table in this map
	// contains a tombstone.
	tombstonePossible bool

	// clearSeq is a sequence counter of calls to Clear. It is used to
	// detect map clears during iteration.
	clearSeq uint64
}
```
As comment notes `dirPtr` is a poiner to an array of swiss tables.

Swiss table structure is looks like this:
```go
// table is a Swiss table hash table structure.
//
// Each table is a complete hash table implementation.
//
// Map uses one or more tables to store entries. Extendible hashing (hash
// prefix) is used to select the table to use for a specific key. Using
// multiple tables enables incremental growth by growing only one table at a
// time.
type table struct {
	// The number of filled slots (i.e. the number of elements in the table).
	used uint16

	// The total number of slots (always 2^N). Equal to
	// `(groups.lengthMask+1)*abi.MapGroupSlots`.
	capacity uint16

	// The number of slots we can still fill without needing to rehash.
	//
	// We rehash when used + tombstones > loadFactor*capacity, including
	// tombstones so the table doesn't overfill with tombstones. This field
	// counts down remaining empty slots before the next rehash.
	growthLeft uint16

	// The number of bits used by directory lookups above this table. Note
	// that this may be less then globalDepth, if the directory has grown
	// but this table has not yet been split.
	localDepth uint8

	// Index of this table in the Map directory. This is the index of the
	// _first_ location in the directory. The table may occur in multiple
	// sequential indices.
	//
	// index is -1 if the table is stale (no longer installed in the
	// directory).
	index int

	// groups is an array of slot groups. Each group holds abi.MapGroupSlots
	// key/elem slots and their control bytes. A table has a fixed size
	// groups array. The table is replaced (in rehash) when more space is
	// required.
	//
	// TODO(prattmic): keys and elements are interleaved to maximize
	// locality, but it comes at the expense of wasted space for some types
	// (consider uint8 key, uint64 element). Consider placing all keys
	// together in these cases to save space.
	groups groupsReference
}
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

### Map iteration is always random

For iteration over map elements, the *`hiter`* structure is responsible. And to initialize this structure *`mapiterinit`* function is used.

```go
// mapiterinit initializes the hiter struct used for ranging over maps.
// The hiter struct pointed to by 'it' is allocated on the stack
// by the compilers order pass or on the heap by reflect_mapiterinit.
// Both need to have zeroed hiter since the struct contains pointers.
func mapiterinit(t *maptype, h *hmap, it *hiter)
```

The order of map elements is undefined. And it is undefined because the order of the elements depends on a very large number of factors, such as the hash function used, the size, whether there were evacuations, etc. But iteration process over a map is **ALWAYS** random because:

```go
/* ...snip... */
// decide where to start
r := uintptr(fastrand())
/* ...snip... */
```

## Conclusion
Knowing the internals of a map will help you write more efficient code 😉.

## Links

- Old (Go < v1.24) map [implementation](https://github.com/golang/go/blob/go1.23.9/src/runtime/map.go)
- New Go's builtin map [implementation](https://github.com/golang/go/blob/go1.26.0/src/internal/runtime/maps/map.go) for Go v1.26
- Swiss table [implementation](https://github.com/golang/go/blob/master/src/internal/runtime/maps/table.go)
