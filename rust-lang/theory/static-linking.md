# Static linking with `+crt-static`

Sometimes it is necessary to create statically linked binary. The command below can help: 

```sh
RUSTFLAGS='-C target-feature=+crt-static' cargo build
```
# Links
More about static linking and this command/flag:
- https://doc.rust-lang.org/rustc/codegen-options/index.html#target-feature
- https://doc.rust-lang.org/reference/linkage.html#static-and-dynamic-c-runtimes
- https://doc.rust-lang.org/rustc/codegen-options/index.html#linking-effects