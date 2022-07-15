# Memory management in Go
Memory in Go is managed by a runtime ( **goruntime** ). This frees the developer from having to deal with device memory directly and transfers **allocation** and freeing (**garbage collection**) responsibilities to runtime. There are two places in memory that Go works with - the *heap* and the *stack*. The garbage collector works only on the memory allocated in the heap, due to the "*features*" of the organization of memory on the stack.

- **The stack** is a *LIFO* queue that often stores the values of the called function. During a function call, the return address is pushed onto the stack, after which a **stack frame** is allocated to store the variables of the new function in this space. After the function is executed, depending on the calling convention, the stack is cleared and the transition to the return address occurs.
- **The heap** is an area in memory that can be resized as needed.

**Go runtime** prefers to allocate memory on the stack, so most of the allocated memory will be there. Each goroutine has its own stack and, if possible, the runtime will allocate space in it. As an optimization, the compiler tries to carry out the so-called ***escape analysis***, that is, to check that the function variable is not mentioned outside its scope. If the compiler manages to find out the lifetime of a variable, it will be placed on the stack, otherwise it will be placed on the heap. Usually, if a program has a pointer to an object, that object is stored on the heap. Unlike the stack, memory on the heap must be manually freed. In the case of Go, this is the responsibility of the garbage collector.

## Garbage collector
Go uses non-generational concurrent, tri-color mark and sweep garbage collector. What this means will be described in more detail below.
- **non-generational**. The generational hypothesis suggests that short-lived objects, such as temporary variables, are most often cleaned up. Thus, the generational garbage collector focuses on newly allocated objects. However, as mentioned earlier, compiler optimizations allow Go to allocate objects on the stack with a known lifetime. This means that there will be fewer objects on the heap that will be collected by the garbage collector. In turn, this means that a generational garbage collector is not needed in Go. So Go uses a non-generational garbage collector.
- **concurrent** means concurrent 😁, meaning the GC runs concurrently with the mutator threads. Therefore, Go uses a non-generational concurrent garbage collector.
- **Mark and sweep** is a type of garbage collector and **tri-color** is the algorithm used to implement it. This type of garbage collector has two phases - mark and sweep , which means mark and clean. During the first phase, it walks through the memory on the heap and marks objects that can be removed. During the second phase, the marked areas of memory are cleared.

In Go, this is implemented like this:
1. All goroutines reach the safe point of garbage collection through a process called **stop the world**. This temporarily stops the execution of the program and turns on the write barrier to maintain the integrity of the data on the heap. This approach provides parallelism by allowing goroutines and GC to run at the same time. When all goroutines activate the barrier, the Go runtime “starts the world” and forces the workers to do the garbage collection work.

2. Marking is carried out by the tri-color algorithm. When marking starts, all objects are white except for the gray root objects. Roots are the object from which all other heap objects are taken and are created as part of program execution. The garbage collector starts marking by looking at stacks, global variables, and heap pointers to understand what is being used. When scanning a stack, the worker stops the goroutine and marks all found objects in gray, moving down from the roots. It then resumes the goroutine.

3. The gray objects are then queued to turn black, which means they are still in use. Once all gray objects turn black, the picker will stop the world again and clean up any white nodes that are no longer needed. The program can now continue running until it needs to clear the extra memory again.
This process is started again when the program allocates additional memory proportional to the memory used. This is defined by the `GOGC` environment variable, which is set to 100 by default. The Go source code describes it like this:

> If GOGC = 100 and we use 4M, we will be GC again when we get to 8M (this mark is tracked in the next_gc variable). This allows you to keep the cost of garbage collection in a linear proportion to the cost of allocation. The GOGC setting simply changes the linear constant (as well as the amount of extra memory used).

### Links:
- Article [original](https://medium.com/@ankur_anand/a-visual-guide-to-golang-memory-allocator-from-ground-up-e132258453ed)
- About the Go garbage collector in the [go blog](https://blog.golang.org/ismmkeynote) and [google groups](https://groups.google.com/g/golang-nuts/c/KJiyv2mV2pU/m/wdBUH1mHCAAJ).
- About switching to [another type of GC](http://golang.org/s/gctoc)
