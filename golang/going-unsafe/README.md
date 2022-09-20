# Unsafe Go
By nature, Go is a safe programming language. But there are still possibilities to write an unsafe code. The `unsafe` package provide functionality to work with low-lewel ("*unsafe*") code. It should be noted, that unlike in **Rust** there is no needs to mark code as unsafe, you just use unsafe functions whenewer they are needed. And there is one more package, that can help with unsafe code - `reflect`. Lets jump in!

## Using unsafe

```go
package main

import (
    "fmt"
    "unsafe"
    "reflect"
)

func appender(isl []int) {
    isl = append(isl, 100)
    isl = append(isl, 200)
    isl = append(isl, 300)
}

func lowLevelSLice(isl []int) {
    intSize := unsafe.Sizeof(int(1))
    sh := (*reflect.SliceHeader)(unsafe.Pointer(&isl))
    datPtr := unsafe.Pointer(sh.Data)

    fmt.Printf("[UNSAFE] len %d cap %d\n", sh.Len, sh.Cap)
    fmt.Printf("[UNSAFE] data ptr %v\n", datPtr)
    for i := 0; i < sh.Cap; i++ {
        dati := *(*int)(unsafe.Pointer(sh.Data + intSize * uintptr(i)))
        fmt.Printf("[UNSAFE] data[%d] == %v\n", i, dati)
    }
}

func main() {
    x := make([]int, 0, 3)
    fmt.Printf("[SAFE] before appender: %v len %d cap %d\n", x, len(x), cap(x))
    appender(x)
    fmt.Printf("[SAFE] after appender: %v len %d cap %d\n", x, len(x), cap(x))
    fmt.Println("GOING UNSAFE AFTER appender")
    lowLevelSLice(x)
    fmt.Printf("[SAFE] after UNSAFE: %v len %d cap %d\n", x, len(x), cap(x))
}

```

Let's figure out what is happening in the code example above;

1. the `appender` function has one argument - int slice. Then it appends 3 values to this slice and do nothing more
2. the `lowLevelSLice` function has the same argument as the `appender`. The body of this function will be described in detail later
3. in the `main` function we create an slice with length 0 and capacity 3 and then we work with this slice
    1. before sending slice to the appender function we can see, that there is no data inside it and this len is equal to 0 and capacity is 3
    2. after the calling the function we can expect the changes in the slice, but nothing happend!
4. calling unsafe `lowLevelSlice` and have some magic

## Meat of the lowLevelSlice

Now, lets watch what is happening inside `lowLevelSLice`. For explanation purposes, i will add comments and make some changes to the code:

```go
func lowLevelSLice(isl []int) {
    /*
        using unsafe.Sizeof to figure out the size of the type to use low level memory arithmetics
    */
    intSize := unsafe.Sizeof(int(1))
    /*
        unsafe.Pointer is a raw pointer to the data.
        more info can be found in the official documentation https://pkg.go.dev/unsafe#Pointer
    */
    llPtr := unsafe.Pointer(&isl)
    /*
        here reflect.SliceHeader comes into action. We are casting an raw pointer to the slice data into the SliceHeader type
    */
    sh := (*reflect.SliceHeader)(llPtr)
    /*
        getting raw pointer to the slice data, i.e. underlying array
    */
    datPtr := unsafe.Pointer(sh.Data)

    fmt.Printf("[UNSAFE] len %d cap %d\n", sh.Len, sh.Cap)
    fmt.Printf("[UNSAFE] data ptr %v\n", datPtr)
    for i := 0; i < sh.Cap; i++ {
        /*
            tere is memory arithmetics. Like in C/C++ we can add value to the pointer.
            As was mentioned earlier, we get the integer size
            (which is 8 bytes on 64 bit and 4 bytes on 32 bit architecture)
            and now we adding an multiplication of the size with the uintptr with
            number of element.
        */
        dati := *(*int)(unsafe.Pointer(sh.Data + intSize * uintptr(i)))
        fmt.Printf("[UNSAFE] data[%d] == %v\n", i, dati)
    }
}
```

When this function called after `appender`, which at the first glance does nothing, we will have the following output:
```
[UNSAFE] data[0] == 100
[UNSAFE] data[1] == 200
[UNSAFE] data[2] == 300
```
`appender` function appends data to the slice, but as Go using copy semantics, it doesn't change the original of the passed argument. But, the slice by itself contains an reference to the underlying array, which is affected by append function in `appender`. So, data of the slice is modified and safe code isn't able to saw that, as len and cap aren't changed.

### Changing appender

Lets change appender function a little bit
```go
func appender(isl []int) {
    isl = append(isl, 100)
    isl = append(isl, 200)
    isl = append(isl, 300)
    isl = append(isl, 400)
    isl = append(isl, 500)
}
```

If we run the code, the output of lowLevelSLice will not change. But, what if we slightly change the loop?
```go
func lowLevelSLice(isl []int) {
    // --snip--
    for i := 0; i < sh.Cap+2; i++ { // add 2 elements to capacity
        dati := *(*int)(unsafe.Pointer(sh.Data + intSize * uintptr(i)))
        fmt.Printf("[UNSAFE] data[%d] == %v\n", i, dati)
    }
}

```

Output will be the following:
```sh
# --snip--
[UNSAFE] data ptr 0xc00001c228
[UNSAFE] data[0] == 100
[UNSAFE] data[1] == 200
[UNSAFE] data[2] == 300
[UNSAFE] data[3] == 7070753967084491611
[UNSAFE] data[4] == 8097789225072551525
```

What is this? Where is 400 and 500?! So, as Go using copy semantics, underlying array of the passed argument to the appender function was reallocated, but the original was not changed at all. When we passing it to the lowLevelSLice function, the underlying array poiner is the same, which was in the argument definition.

Lets modify appender function for the last time:
```go
func appender(isl []int) {
    var (
        sh *reflect.SliceHeader
        datPtr unsafe.Pointer
    )

    isl = append(isl, 100)
    sh = (*reflect.SliceHeader)(unsafe.Pointer(&isl))
    datPtr = unsafe.Pointer(sh.Data)
    fmt.Printf("[APPENDER] data pointer after first append %p\n", datPtr)

    isl = append(isl, 200)
    sh = (*reflect.SliceHeader)(unsafe.Pointer(&isl))
    datPtr = unsafe.Pointer(sh.Data)
    fmt.Printf("[APPENDER] data pointer after second append %p\n", datPtr)

    isl = append(isl, 300)
    sh = (*reflect.SliceHeader)(unsafe.Pointer(&isl))
    datPtr = unsafe.Pointer(sh.Data)
    fmt.Printf("[APPENDER] data pointer after third append %p\n", datPtr)

    isl = append(isl, 400)
    sh = (*reflect.SliceHeader)(unsafe.Pointer(&isl))
    datPtr = unsafe.Pointer(sh.Data)
    fmt.Printf("[APPENDER] data pointer after fourth append %p\n", datPtr)

    isl = append(isl, 500)
    sh = (*reflect.SliceHeader)(unsafe.Pointer(&isl))
    datPtr = unsafe.Pointer(sh.Data)
    fmt.Printf("[APPENDER] data pointer after fifth append %p\n", datPtr)
}
```

Now the function is very verbose, and we can monitor, how the underlying data pointer is changing

```sh
# --snip--
[APPENDER] data pointer after first append 0xc00001c228
[APPENDER] data pointer after second append 0xc00001c228
[APPENDER] data pointer after third append 0xc00001c228
[APPENDER] data pointer after fourth append 0xc000022210 # AHA! The pointer is changed!
[APPENDER] data pointer after fifth append 0xc000022210
# --snip--
```

Yeah. Thats it!

----

## Conclusion

Why use unsafe? It is faster. Run some benchmarks (someday there will be benchmarks) and u see performance boost, because unsafe code didn't check array bounds. This leads to the security problems, such as information disclosure (iff by one) and make possible the buffer overflow vulns.

----

- TODO: add map example
- TODO: add comparentment benchmarks
