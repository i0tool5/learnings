# Template Method

Template Method pattern defines method in a superclass, usually an abstract superclass, and defines the skeleton of an operation and allows subclasses to override (hooks) specific steps of the algorithm without changing its structure.

This pattern is used when:
- it is necessary to provide variying behavior, keeping the original algorithm of actions intact
- you need to avoid code duplication

The pattern operates with two (or more) entities:
1. **Abstract class** - contains a template method or methods and provides skeleton, where certain steps are defined but others are left abstract and subclasses can override them. 
0. **Concrete class(es)** - specific implementation class(es) that overrides *"hook"* method(s). 

## Implementation examples

### Python

```python
from abc import ABC, abstractmethod


class AbstractRunner(ABC):
    def start(self):
        print("starting")

    @abstractmethod
    def handle(self):
        ...

    def finish(self):
        print("finishing")
    
    def run(self):
        self.start()
        self.handle()
        self.finish()


class SomeImplementation(AbstractRunner):
    def handle(self):
        print("handling something")


class AnotherImplementation(AbstractRunner):
    def handle(self):
        print("another handling implementation")


def main():
    impl = SomeImplementation()
    impl.run()

    impl = AnotherImplementation()
    impl.run()


if __name__ == "__main__":
    main()
```
