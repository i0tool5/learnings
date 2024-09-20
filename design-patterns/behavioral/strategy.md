# Strategy

The Strategy pattern allows to choose the behavior of an object at runtime receiving instructions as to which in a family of algorithms to use. It defines a set of algorithms that are similar in behavior, encapsulate each one of them in separate class, and make them interchangeable.

This pattern is used when:
- it is necessary to dynamically select and switch between different algorithms at runtime based on user preferences, configuration settings, or system states
- there is multiple algorithms that can be used interchangeably based on different contexts 

The pattern operates with three or more entities:
1. **Context** - class, which represents the context of executing a particular strategy
0. **AbstractStrategy** - abstract class that defines the interface of a particular set of strategies
0. **ConcreteStrategyA** - specific strategy implementation for abstract strategy interface
0. **ConcreteStrategyB** - another implementation of specific strategy for abstract strategy interface

## Implementation examples

### Python

```python
from abc import ABC, abstractmethod


class CompressionStrategy(ABC):
    @abstractmethod
    def compress(self, data):
        ...


class GZIPCompressionStrategy(CompressionStrategy):
    def compress(self, data: bytes):
        print(f"compressing {data} with gzip mechanism")


class LZMACompressionStrategy(CompressionStrategy):
    def compress(self, data: bytes):
        print(f"compressing {data} with lzma mechanism")


class XZCompressionStrategy(CompressionStrategy):
    def compress(self, data: bytes):
        print(f"compressing {data} with xz mechanism")


class CompressionContext:
    def __init__(self):
        self.__compression_strategy: CompressionStrategy = None
    
    def compress(self, data):
        self.__compression_strategy.compress(data)
    
    def set_strategy(self, strategy: CompressionStrategy):
        self.__compression_strategy = strategy


def gzip_client(ctx: CompressionContext):
    data = "data for gzip".encode()
    ctx.compress(data)


def xz_client(ctx: CompressionContext):
    data = "data for xz".encode()
    ctx.compress(data)


def lzma_client(ctx: CompressionContext):
    data = "data for lzma".encode()
    ctx.compress(data)


def main():
    compression_ctx = CompressionContext()

    gzip = GZIPCompressionStrategy()
    compression_ctx.set_strategy(gzip)
    gzip_client(compression_ctx)

    xz = XZCompressionStrategy()
    compression_ctx.set_strategy(xz)
    xz_client(compression_ctx)

    lzma = LZMACompressionStrategy()
    compression_ctx.set_strategy(lzma)
    lzma_client(compression_ctx)


if __name__ == "__main__":
    main()
```
