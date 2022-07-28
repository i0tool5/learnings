# Fat pointers in Rust terminology

[Original](https://stackoverflow.com/questions/57754901/what-is-a-fat-pointer#57754902) from stack overflow

### What is fat pointer

**The term "fat pointer" is used to refer to references and raw pointers to** ***dynamically sized types*** **(DSTs) – slices or trait objects**. A fat pointer contains a pointer plus some information that makes the DST "complete" (e.g. the length).

Most commonly used types in Rust are not DSTs but have a fixed size known at compile time. These types implement the [`Sized`](https://doc.rust-lang.org/stable/std/marker/trait.Sized.html) trait. Even types that manage a heap buffer of dynamic size (like `Vec<T>`) are `Sized`, as the compiler knows the exact number of bytes a `Vec<T>` instance will take up on the stack. There are currently four different kinds of DSTs in Rust.

## **Slices ([T] and str)**
The type `[T]` is dynamically sized (so is the special "string slice" type [`str`](https://doc.rust-lang.org/std/primitive.str.html)). That's why you usually only see it as `&[T]` or `&mut [T]`, i.e. behind a reference. This reference is a so-called "fat pointer". Let's check:

```rust
dbg!(size_of::<&u32>());
dbg!(size_of::<&[u32; 2]>());
dbg!(size_of::<&[u32]>());
```
The output is following:
```
[main.rs:15] size_of::<&u32>() = 8
[main.rs:16] size_of::<&[u32; 2]>() = 8
[main.rs:17] size_of::<&[u32]>() = 16
```
So we see that a reference to a normal type like `u32` is 8 bytes large, as is a reference to an array `[u32; 2]`. Those two types are not DSTs. But as `[u32]` is a DST, the reference to it is twice as large. **In the case of slices, the additional data that "completes" the DST is simply the length**. So one could say the representation of `&[u32]` is something like this:
```rust
struct SliceRef { 
    ptr: *const u32, 
    len: usize,
}
```

## Trait objects (dyn Trait)
When using traits as trait objects (i.e. type erased, dynamically dispatched), these trait objects are DSTs. Example:
```rust
use std::mem::size_of;

trait Animal {
    fn speak(&self);
}

struct Cat;
impl Animal for Cat {
    fn speak(&self) {
        println!("meow");
    }
}

fn main() {
    dbg!(size_of::<&Cat>());
    dbg!(size_of::<&dyn Animal>());
}
```
This prints
```
[main.rs:15] size_of::<&Cat>() = 8
[main.rs:16] size_of::<&dyn Animal>() = 16
```
Again, *`&Cat`* is only 8 bytes large because `Cat` is a normal type. But `dyn Animal` is a trait object and therefore dynamically sized. As such, *`&dyn Animal`* is 16 bytes large.

**In the case of trait objects, the additional data that completes the DST is a pointer to the vtable (the vptr)**. Shortly speaking, *vtables* and *vptrs* are used to call the correct method implementation in this virtual dispatch context. The vtable is a static piece of data that basically only contains a function pointer for each method. With that, a reference to a trait object is basically represented as:
```rust
struct TraitObjectRef {
    data_ptr: *const (),
    vptr: *const (),
}
```
This is different from C++, where the *vptr* for abstract classes is stored within the object. Both approaches have advantages and disadvantages.

## Custom DSTs

It's actually possible to create your own DSTs by having a struct where the last field is a DST. This is rather rare, though. One prominent example is [`std::path::Path`](https://doc.rust-lang.org/std/path/struct.Path.html).

A reference or pointer to the custom DST is also a fat pointer. The additional data depends on the kind of DST inside the struct.

## Exception: Extern types

In [RFC 1861](https://github.com/rust-lang/rfcs/blob/master/text/1861-extern-types.md), the `extern type` feature was introduced. Extern types are also DSTs, but pointers to them are not fat pointers. Or more exactly, as the RFC puts it:

> In Rust, pointers to DSTs carry metadata about the object being pointed to. For strings and slices this is the length of the buffer, for trait objects this is the object's vtable. For extern types the metadata is simply `()`. This means that a pointer to an extern type has the same size as a `usize` (ie. it is not a "fat pointer").

But if you are not interacting with a C interface, you probably won't ever have to deal with these extern types.

### Links

- One more time [original from stackoverflow](https://stackoverflow.com/questions/57754901/what-is-a-fat-pointer#57754902)
- Discussion about fat pointers [on reddit (1)](https://www.reddit.com/r/rust/comments/8ckfdb/were_fat_pointers_a_good_idea/)
- Discussion about fat pointers [on reddit (2)](https://www.reddit.com/r/ProgrammingLanguages/comments/fpnooj/tradeoffs_of_fat_vs_thin_pointers/)
