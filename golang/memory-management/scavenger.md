# Scavenger

Scavenger - is the runtime element responsible for returning RAM to the operating system. Scavenger performs a very non-trivial task. On the one hand, it tries to keep the RSS (*resident set size*) of the process within reasonable limits. On the other hand, returning all unused memory is inefficient, because when it is needed, it will have to be re-allocated, which is an expensive and time-consuming operation.

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
