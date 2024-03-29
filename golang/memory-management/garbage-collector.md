# Garbage collector
Go uses [tracing](https://en.wikipedia.org/wiki/Tracing_garbage_collection), non-generational, concurrent, tri-color mark and sweep (non-moving), incremental garbage collector. What this means will be described in more detail below.
- **non-generational**. The generational hypothesis suggests that short-lived objects, such as temporary variables, are most often cleaned up. Thus, the generational garbage collector focuses on newly allocated objects. However, as mentioned earlier, compiler optimizations allow Go to allocate objects on the stack with a known lifetime. This means that there will be fewer objects on the heap that will be collected by the garbage collector. In turn, this means that a generational garbage collector is not needed in Go. So Go uses a non-generational garbage collector.
- **concurrent** means concurrent 😁, meaning the GC runs concurrently with the mutator and busines logic. Therefore, Go uses a non-generational concurrent garbage collector.
- **Mark and sweep** is a type of garbage collector and **tri-color** is the algorithm used to implement it. This type of garbage collector has two phases - mark and sweep , which means mark and clean. During the first phase, it walks through the memory on the heap and marks objects that can be removed. During the second phase, the marked areas of memory are cleared.
- **non-moving** means that GC isn't move the live objects to a new part of memory. This technique is opposite to *mark and sweep*, where unused objects are marked to be cleared and can be reused. 

> **Note**: GC has its own pool of goroutines, that run concurrently with business logic goroutines

In Go, this is implemented like this:
1. *Sweep Termination*. All goroutines reach the safe point of garbage collection through a process called **[stop the world](#gc-stop-the-world)**. This temporarily stops the execution of the entire program and turns on the [write barrier](#write-barrier) to maintain the integrity of the data on the heap and goroutines stack (which is also on the heap). After that, all memory blocks marked as garbage (white) are sent to Scavenger. The approach with write barrier provides parallelism by allowing goroutines and GC to run at the same time. When all GC goroutines activate the barrier, the Go runtime “starts the world” and forces the workers to do the garbage collection work.

2. Marking is carried out by the tri-color algorithm. When marking starts, all objects are white except for the gray root objects. Roots are the object from which all other heap objects are taken and are created as part of program execution. The garbage collector starts marking by looking at stacks, global variables, and heap pointers to understand what is being used. When scanning a stack, the worker stops the goroutine and marks all found objects in gray, moving down from the roots. It then resumes the goroutine.

3. The gray objects are then queued to turn black, which means they are still in use. Once all gray objects turn black, the picker will stop the world again and clean up any white nodes that are no longer needed. The program can now continue running until it needs to clear the extra memory again.
This process is started again when the program allocates additional memory proportional to the memory used. This is defined by the `GOGC` environment variable, which is set to 100 by default. The Go source code describes it like this:

> If GOGC = 100 and we use 4M, we will be GC again when we get to 8M (this mark is tracked in the next_gc variable). This allows you to keep the cost of garbage collection in a linear proportion to the cost of allocation. The GOGC setting simply changes the linear constant (as well as the amount of extra memory used).

The current stage of the garbage collector is stored in the global variable ***gcphase***.

## GC Stop The World

*Stop The World* (STW) is a technique, that helps keep memory consistent, but increases latency of the entire program. Until the Go 1.5 the STW was the only one opportunity for the Go GC, but in Go 1.5 it was changed. Since version 1.5, a concurrent model has been added in addition to STW. The algorithm of actions became as follows: 
- goroutine stop time (STW phase) has been reduced to 10 milliseconds per GC iteration (50 milliseconds total).
- if STW does not meet the allotted 10 milliseconds, then the GC switches to concurrent mode, running simultaneously with the rest of the program for the remaining 40 milliseconds.

### Write barrier

GC is able to collect garbage step by step, in small pieces, which means incremental GC. For incremental GC, you need guarantees that if you add or remove links between objects during painting, the graphs will remain correctly colored. This is where write barrier comes into play.

**Write Barrier** is a piece of code that is executed while working with memory. It is needed to support invariants that guarantee the correct step-by-step execution of the algorithm.

Go uses combination of **Dijkstra insertion** and **Yuasa deletion** write barrier.
- *insertion* in the name means that the trigger for its call is the creation of a connection between objects
- *deletion* in the name means that the trigger for the call is the removal of a connection between objects

The first type of barrier consistently support this invariant: black objects only point to gray or other black objects, and white objects do not. The second one gives the following invariant: any white object pointed to by a black object must have a reachable path to a gray object.

## Sweep

Sweep is running in the background, releasing the memory. marked white. But this process may not be complete in time. Therefore, it should be noted that a part of the sweeping is laid of the memory allocation process. When process needs more memory, it calls memory allocation system, but before allocation new memory, it must make shure that there are no old unnecessary objects in memory. 

# Garbage Collector Pacer

**Garbage Collector Pacer** - is mechanism that keeps track of the moment at which garbage collection should begin.

There is a proposal for [pacer redesign](https://github.com/golang/proposal/blob/master/design/44167-gc-pacer-redesign.md).
