# Memory management in Go
Memory in Go is managed by a runtime ( **goruntime** ). This frees the developer from having to deal with device memory directly and transfers **allocation** and freeing (**garbage collection**) responsibilities to runtime. There are two places in memory that Go works with - the *heap* and the *stack*. The garbage collector works only on the memory allocated in the heap, due to the "*features*" of the organization of memory on the stack.

- **The stack** is a *LIFO* queue that often stores the values of the called function. During a function call, the return address is pushed onto the stack, after which a **stack frame** is allocated to store the variables of the new function in this space. After the function is executed, depending on the calling convention, the stack is cleared and the transition to the return address occurs.
- **The heap** is an area in memory that can be resized as needed.

**Go runtime** prefers to allocate memory on the stack, so most of the allocated memory will be there. Each goroutine has its own stack and, if possible, the runtime will allocate space in it. As an optimization, the compiler tries to carry out the so-called ***escape analysis***, that is, to check that the function variable is not mentioned outside its scope. If the compiler manages to find out the lifetime of a variable, it will be placed on the stack, otherwise it will be placed on the heap. Usually, if a program has a pointer to an object, that object is stored on the heap. Unlike the stack, memory on the heap must be manually freed. In the case of Go, this is the responsibility of the garbage collector.

## Allocator

Go runtime uses **[TCMalloc](https://github.com/google/tcmalloc)** instead of using Linux `malloc`. This is due to the fact that **TCMalloc** tries to minimize possible [memory fragmentation](https://en.wikipedia.org/wiki/Fragmentation_(computing)#Internal_fragmentation). Memory fragmentation can also be eliminated by using the *GC*, which will move objects to other places in memory. But this leads to an additional overhead, since the addresses of the moved objects change, and you also have to modify the links of those objects that refer to them. This is an expensive operation that can cause an increase in process response time.

## Garbage collector
Go uses [tracing](https://en.wikipedia.org/wiki/Tracing_garbage_collection) non-generational concurrent, tri-color mark and sweep (non-moving) garbage collector. What this means will be described in more detail below.
- **non-generational**. The generational hypothesis suggests that short-lived objects, such as temporary variables, are most often cleaned up. Thus, the generational garbage collector focuses on newly allocated objects. However, as mentioned earlier, compiler optimizations allow Go to allocate objects on the stack with a known lifetime. This means that there will be fewer objects on the heap that will be collected by the garbage collector. In turn, this means that a generational garbage collector is not needed in Go. So Go uses a non-generational garbage collector.
- **concurrent** means concurrent 😁, meaning the GC runs concurrently with the mutator threads. Therefore, Go uses a non-generational concurrent garbage collector.
- **Mark and sweep** is a type of garbage collector and **tri-color** is the algorithm used to implement it. This type of garbage collector has two phases - mark and sweep , which means mark and clean. During the first phase, it walks through the memory on the heap and marks objects that can be removed. During the second phase, the marked areas of memory are cleared.
- **non-moving** means that GC isn't move the live objects to a new part of memory. This technique is opposite to *mark and sweep*, where unused objects are marked to be cleared and can be reused. 

> **Note**: GC has its own pool of goroutines, that run concurrently with business logic goroutines

In Go, this is implemented like this:
1. All goroutines reach the safe point of garbage collection through a process called **stop the world**. This temporarily stops the execution of the entire program and turns on the write barrier to maintain the integrity of the data on the heap. This approach provides parallelism by allowing goroutines and GC to run at the same time. When all GC goroutines activate the barrier, the Go runtime “starts the world” and forces the workers to do the garbage collection work.

2. Marking is carried out by the tri-color algorithm. When marking starts, all objects are white except for the gray root objects. Roots are the object from which all other heap objects are taken and are created as part of program execution. The garbage collector starts marking by looking at stacks, global variables, and heap pointers to understand what is being used. When scanning a stack, the worker stops the goroutine and marks all found objects in gray, moving down from the roots. It then resumes the goroutine.

3. The gray objects are then queued to turn black, which means they are still in use. Once all gray objects turn black, the picker will stop the world again and clean up any white nodes that are no longer needed. The program can now continue running until it needs to clear the extra memory again.
This process is started again when the program allocates additional memory proportional to the memory used. This is defined by the `GOGC` environment variable, which is set to 100 by default. The Go source code describes it like this:

> If GOGC = 100 and we use 4M, we will be GC again when we get to 8M (this mark is tracked in the next_gc variable). This allows you to keep the cost of garbage collection in a linear proportion to the cost of allocation. The GOGC setting simply changes the linear constant (as well as the amount of extra memory used).

## Scavenger

Scavenger - is the runtime element responsible for returning RAM to the operating system. Scavenger performs a very non-trivial task. On the one hand, it tries to keep the RSS of the process within reasonable limits. On the other hand, returning all unused memory is inefficient, because when it is needed, it will have to be re-allocated, which is an expensive and time-consuming operation.

It's important to note that Go never uses the **munmap** syscall to reduce the size of the virtual address space - this call is considered expensive, and virtual memory, under which there is no physical memory, on the contrary, is very cheap. Runtime can only tell the operating system that it no longer needs certain pages of virtual memory - the **madvice** system call is used for this. If there are pages of physical memory under these pages of virtual memory, the OS, at its discretion, can take them for itself for use by other processes. *`The call to madvice is the main task of scavenger`*.

At the same time, new allocations in runtime can occur both from ordinary pages and from pages “returned” to the operating system. This explains the [counterintuitive behavior](https://github.com/golang/go/issues/32284) of the `HeapReleased` indicator. It looks like it should grow monotonically, but in reality it can decrease if the Go runtime decides to access the "returned" virtual memory area via **madvice** in order to allocate new objects.

---

Up to and including **`Go 1.11`**, scavenger ran periodically every 2.5 minutes, and freed pages from spans that had not been used for more than 5 minutes. The **madvice** argument was the *MADV_DONTNEED* parameter: it prompted the OS that the memory can be taken back, while the RSS of the process synchronously decreased at the time of the system call.

---

**`Go 1.12`** added yet another scavenger startup type. Now it was called synchronously at the time of the rapid growth of the heap. For some applications, it helped to cut off the peak consumption of RAM. Another innovation was the replacement of *MADV_DONTNEED* with *MADV_FREE* - this is a lazier option for freeing memory, in which RSS does not decrease immediately, but as the load on the OS increases. Immediately after that, a flurry of memory leak tickets poured out on Go developers, because Docker and Kubernetes did not understand this laziness and killed processes due to exceeding RSS limits.

---

In **`Go 1.13`**, the periodic scavenger has been replaced with a [continuously running](https://github.com/golang/go/issues/30333) scavenger. The background goroutine runs scavenger so that the total CPU cost of running scavenger does not exceed 1% of the CPU, but in reality it costs a little more, since scavenger is forced to synchronize access to internal runtime structures with allocating goroutines. The target RSS value that the background scavenger tries to bring the process to was expressed by this [formula](https://github.com/golang/go/blob/go1.13/src/runtime/mgcscavenge.go#L20):

$$ Goal={retainExtraPercent + 100 \over 100} * NextGC $$

where *retainExtraPercent* was equal to 10, i.e. scavenger laid a 10% buffer of the target heap size for reuse in new allocations. 

---

In **`Go 1.14`**, this [formula has been refined](https://github.com/golang/go/blob/go1.14/src/runtime/mgcscavenge.go#L21) to more explicitly account for allocator fragmentation.

---

In **`Go 1.16`**, the decision was made to [roll back](https://github.com/golang/go/issues/42330) to *MADV_DONTNEED* in order to make process RSS estimation more accurate.

---

In **`Go 1.19`** The introduction of the SetMemoryLimit function will require modifications to the scavenger startup. Now the intensity of the scavenger will be controlled by a PI controller: the closer the memory consumption approaches the limit, the higher will be the CPU fraction allocated to the scavenger (up to 10%). [Release notes](https://go.dev/doc/go1.19#runtime) describing the changes in the runtime associated with the introduction of the soft memory limits feature.

## Summary

Go runtime uses three "tools" to work with memory:
- Allocator (TCMalloc)
- Garbage Collector
- Scavenger

They abstract memory handling, but are still not without flaws.

### Links:
- Official Go [documentation](https://go.dev/doc/gc-guide) about GC
- Article about GC in go [original](https://medium.com/@ankur_anand/a-visual-guide-to-golang-memory-allocator-from-ground-up-e132258453ed)
- About the Go garbage collector in the [go blog](https://blog.golang.org/ismmkeynote) and [google groups](https://groups.google.com/g/golang-nuts/c/KJiyv2mV2pU/m/wdBUH1mHCAAJ).
- About switching to [another type of GC](http://golang.org/s/gctoc)
- Go memory allocator [source code](https://github.com/golang/go/blob/master/src/runtime/malloc.go)
