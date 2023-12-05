# Iterator
The Iterator pattern provides a mechanism for traversing collections of objects without revealing their internal representation. Often this pattern is used instead of an array of objects to not only provide access to elements, but also provide some logic.

This pattern is used when:
- the elements of an aggregate object should be accessed and traversed without exposing its representation (data structures).
- new traversal operations should be defined for an aggregate object without changing its interface.

The pattern operates with four entities:

1. **Iterator interface** - describes a set of methods for accessing the collection
2. **Iterator implementation** - implements the Iterator interface. Monitors the position of the current element when iterating through the collection (Aggregate).
3. **Aggregate interface** - describes a set of methods for a collection of objects
4. **Aggregate implementation** - implements the Aggregate interface and stores collection elements

> It should be mentioned, that aggregate is optional.

## Implementations

### Python

The Python implementation is very short and simple, because *Iterator* type is included into standard library and supports inheritance.

```python
import random
from collections.abc import Iterator


class PacketIterator(Iterator):
    '''Represents iterator'''
    def __init__(self) -> None:
        self.__packets: list[bytearray] = []
        self.__c = 0
    def append(self, packet: bytearray):
        self.__packets.append(packet)
    def __next__(self) -> str:
        if self.__c >= len(self.__packets):
            raise StopIteration
        p = self.__packets[self.__c]
        self.__c += 1
        return p.decode()


class Network:
    '''Represents aggregate'''
    def __iter__(self):
        p = PacketIterator()
        for i in range(random.randint(1, 20)):
            p.append(bytearray(f"somepacket_{random.random() * 100}".encode()))
        return p


def main():
    network = Network()
    for packet in network:
        print(packet)


if __name__ == '__main__':
    main()

```

### Golang

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Iterator provides a iterator interface.
type Iterator interface {
	Next() bool
	Value() any
}

// Aggregate provides a collection interface.
type Aggregate interface {
	Iterator() Iterator
}

var _ Iterator = &PacketIterator{}

// Packet represents network packet.
type Packet []byte

type PacketIterator struct {
	Packets []Packet
	c       int
}

// Next shows that iterator is not empty.
func (pi *PacketIterator) Next() bool {
	return pi.c < len(pi.Packets)
}

// Value returns value at current iteration counter. 
func (pi *PacketIterator) Value() any {
	v := string(pi.Packets[pi.c])
	pi.c++
	return v
}

var _ Aggregate = &Network{}

// Network represents network.
type Network struct {
	r *rand.Rand
}

// Iterator returns packet iterator.
func (n *Network) Iterator() Iterator {
	j := n.r.Int() % 20
	packets := make([]Packet, j)
	for i := 0; i < j; i++ {
		packets[i] = []byte(fmt.Sprintf("somepacket %v", n.r.Int63()))
	}

	return &PacketIterator{packets, 0}
}

func main() {
	network := Network{rand.New(rand.NewSource(time.Now().UnixNano()))}

	packetIterator := network.Iterator()

	for packetIterator.Next() {
		v := packetIterator.Value()
		fmt.Println(v)
	}
}
```