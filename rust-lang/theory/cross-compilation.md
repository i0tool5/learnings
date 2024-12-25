# Cross-compilation in Rust

Rust is a very powerfull and flexible language. It supports huge amount of architectures and also supports cross-compilation.

Example below describes how to prepare enviroment to build binary for arm64 linux on x86-64 linux.

> Note: arm64 and aarch64 are synonims

First configure `cargo` so it knows what to do with specific target. Inside the project root (or other [configuration places](https://doc.rust-lang.org/cargo/reference/config.html)):

```sh
mkdir .cargo
echo '[target.aarch64-unknown-linux-gnu]
linker = "aarch64-linux-gnu-gcc"
' >> .cargo/config.toml
```

Next, add new target toolchain
```sh
rustup target add aarch64-unknown-linux-gnu
```

Also add new architecture to package manager 
```sh
sudo dpkg --add-architecture arm64
```

And finally install required linker
```sh
sudo apt update && apt install gcc-aarch64-linux-gnu
```
> Note: apt and dpkg - are Debian (and Debian-based) specific tools. Other distros, such as Arch or Fedora, will use another commands and package names.

# Links

- About cross-compilation setup in [rustup](https://rust-lang.github.io/rustup/cross-compilation.html)