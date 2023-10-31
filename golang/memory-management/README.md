# Memory management in Go
Memory in Go is managed by a runtime ( **goruntime** ). This frees the developer from having to deal with device memory directly and transfers **allocation** and freeing (**garbage collection** and returning memory to the OS (**scavenger**)) responsibilities to runtime. There are two places in memory that Go works with - the *heap* and the *stack*. The garbage collector works only on the memory allocated in the heap, due to the "*features*" of the organization of memory on the stack.

- **The stack** is a *LIFO* queue that often stores the values of the called function. During a function call, the return address is pushed onto the stack, after which a **stack frame** is allocated to store the variables of the new function in this space. After the function is executed, depending on the calling convention, the stack is cleared and the transition to the return address occurs.
- **The heap** is an area in memory that can be resized as needed.

**Go runtime** prefers to allocate memory on the stack, so most of the allocated memory will be there. Each goroutine has its own stack and, if possible, the runtime will allocate space in it. As an optimization, the compiler tries to carry out the so-called ***escape analysis***, that is, to check that the function variable is not mentioned outside its scope. If the compiler manages to find out the lifetime of a variable, it will be placed on the stack, otherwise it will be placed on the heap. Usually, if a program has a pointer to an object, that object is stored on the heap. Unlike the stack, memory on the heap must be manually freed. In the case of Go, this is the responsibility of the garbage collector.

Go runtime uses three "tools" to work with memory:
- [Allocator](./allocator.md)
- [Garbage Collector](./garbage-collector.md)
- [Scavenger](./scavenger.md)

They abstract memory handling, but are still not without flaws.

# Links:
## Garbage collection
- Official Go [documentation](https://go.dev/doc/gc-guide) about GC
- Article about [GC in go](https://medium.com/@ankur_anand/a-visual-guide-to-golang-memory-allocator-from-ground-up-e132258453ed)
- About the Go garbage collector in the [go blog](https://blog.golang.org/ismmkeynote) and [google groups](https://groups.google.com/g/golang-nuts/c/KJiyv2mV2pU/m/wdBUH1mHCAAJ).
- About switching to [another type of GC](http://golang.org/s/gctoc)
- GC [source code](https://github.com/golang/go/blob/master/src/runtime/mgc.go) (`src/runtime/mgc.go`)
- GC marker [source code](https://github.com/golang/go/blob/master/src/runtime/mgcmark.go)(`src/runtime/mgcmark.go`)
- GC sweeper [source code](https://github.com/golang/go/blob/master/src/runtime/mgcsweep.go)(`src/runtime/mgcsweep.go`)
## Allocator
- Go memory allocator [source code](https://github.com/golang/go/blob/master/src/runtime/malloc.go) (`src/runtime/malloc.go`)
