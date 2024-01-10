# Hash Table / Map / Dict / Associative array

Certainly, a hash map, also known as a hash table, is a data structure that implements an associative array abstract data type, a structure that can map keys to values. A hash table uses a hash function to compute an index, also called a hash code, into an array of buckets or slots, from which the desired value can be found. During lookup, the key is hashed and the resulting hash indicates where the corresponding value is stored. Time complexity for search, insert and remove operations is O(1) in average and O(n) as a worst case. 

## Implementation examples

The implementations are very simple, provided for educational purposes and don't pretend to be perfect, because each programming language has its own perfect production-ready implementation.

### Golang

```go
package main

import (
	"fmt"
	"hash/fnv"
)

// Entry is a key-value pair.
type Entry struct {
	Key   string
	Value interface{}
	next  *Entry
}

// Bucket is a linked list of Entries.
type Bucket struct {
	head *Entry
}

// HashMap is a structure that defines the hash map.
type HashMap struct {
	buckets []*Bucket
	size    int
}

// NewBucket initializes a new Bucket.
func NewBucket() *Bucket {
	return &Bucket{head: nil}
}

// Insert adds a key-value pair to the bucket, or updates the value if the key exists.
func (b *Bucket) Insert(key string, value interface{}) {
	current := b.head
	for current != nil {
		if current.Key == key {
			current.Value = value
			return
		}
		current = current.next
	}

	newEntry := &Entry{Key: key, Value: value, next: b.head}
	b.head = newEntry
}

// Get retrieves a value from the bucket by key
func (b *Bucket) Get(key string) (interface{}, bool) {
	current := b.head
	for current != nil {
		if current.Key == key {
			return current.Value, true
		}
		current = current.next
	}
	return nil, false
}

// Remove deletes an entry from the bucket by key
func (b *Bucket) Remove(key string) {
	if b.head == nil {
		return
	}

	if b.head.Key == key {
		b.head = b.head.next
		return
	}

	prev := b.head
	for prev.next != nil {
		if prev.next.Key == key {
			prev.next = prev.next.next
			return
		}
		prev = prev.next
	}
}

// NewHashMap creates a new HashMap with a given initial size
func NewHashMap(size int) *HashMap {
	buckets := make([]*Bucket, size)
	for i := range buckets {
		buckets[i] = NewBucket()
	}
	return &HashMap{
		buckets: buckets,
		size:    size,
	}
}

// hashCode generates a hash for a given string
func hashCode(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32())
}

// getIndex calculates the index of a key in the array of buckets
func (h *HashMap) getIndex(key string) int {
	hash := hashCode(key)
	index := hash % h.size
	if index < 0 {
		index += h.size
	}
	return index
}

// Insert adds a key-value pair to the HashMap
func (h *HashMap) Insert(key string, value interface{}) {
	index := h.getIndex(key)
	h.buckets[index].Insert(key, value)
}

// Get retrieves a value from the HashMap by key
func (h *HashMap) Get(key string) (interface{}, bool) {
	index := h.getIndex(key)
	return h.buckets[index].Get(key)
}

// Remove deletes a key-value pair from the HashMap
func (h *HashMap) Remove(key string) {
	index := h.getIndex(key)
	h.buckets[index].Remove(key)
}

// main function to demonstrate usage
func main() {
	hm := NewHashMap(10) // create a hash map with 10 buckets
	hm.Insert("key1", 100)
	hm.Insert("key2", 200)

	value, exists := hm.Get("key1")
	if exists {
		fmt.Println("key1:", value)
	}

	hm.Remove("key1")
	value, exists = hm.Get("key1")
	if !exists {
		fmt.Println("key1 has been removed")
	}
	value, exists = hm.Get("key2")
	if exists {
		fmt.Println("key2:", value)
	}
}
```

This implementation includes a HashMap struct, which contains buckets, and each bucket is a linked list of entries. The Entry struct represents a key-value pair. The hash map has the basic operations:

- Inserting a key-value pair
- Retrieving a value by key
- Removing a key-value pair by key

It uses FNV (Fowler-Noll-Vo) for hashing and separate chaining to handle collisions.
