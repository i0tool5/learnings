# Visitor
Visitor - is a behavioral pattern that allows to perform some unrelated operations on a number of objects, avoiding contamination of their code.

This pattern is used when:
- there are different objects of different classes with different interfaces, but it is necessary to perform class-specific operations on them;
- it is necessary to perform various operations on the structure that complicate the structure;
- new operations on the structure are often added.

The pattern operates with N entities:
1. **AbstractVisitor** - abstract base class, which describes Visitor interface
0. **ConcreteVisitor** - implements methods for visiting a particular element
0. **ObjectStructure** - implements a structure (collection) that stores the elements to be visited
0. **AbstractElement** - describes interface of elements of the **ObjectStructure**
0. **ConcreteElementN** - implements element of the **ObjectStructure**

## Implementation examples

### Python

```python
from abc import ABC, abstractmethod


class VisitorABC(ABC):
    @abstractmethod
    def visit(self, element):
        ...


class ElementABC(ABC):
    @abstractmethod
    def accept(self, visitor: VisitorABC):
        ...


class PersonalComputerVisitor(VisitorABC):
    def visit(self, element):
        match element:
            case CPU():
                self._visit_cpu(element)
            case Motherboard():
                self._visit_motherboard(element)
            case RAM():
                self._visit_ram(element)
            case PSU():
                self._visit_psu(element)
            case GPU():
                self._visit_gpu(element)
            case PersonalComputer():
                self._visit_pc(element)
            case _:
                raise Exception(f'Unknown part {element.__class__.__name__}')
    
    def _visit_cpu(self, element):
        raise NotImplementedError

    def _visit_motherboard(self, element):
        raise NotImplementedError
    
    def _visit_ram(self, element):
        raise NotImplementedError
    
    def _visit_psu(self, element):
        raise NotImplementedError
    
    def _visit_gpu(self, element):
        raise NotImplementedError
    
    def _visit_pc(self, element):
        raise NotImplementedError


class PCAssemblyValidator(PersonalComputerVisitor):
    def _visit_cpu(self, element):
        print(f'CPU is in place')

    def _visit_motherboard(self, element):
        print('Motherboard is in place')
    
    def _visit_ram(self, element):
        print('RAM is in place')

    def _visit_gpu(self, element):
        print(f'GPU is in place')
    
    def _visit_psu(self, element):
        print('PSU is in place')
    
    def _visit_pc(self, element):
        print('PC is complete')


class PCOperabilityValidator(PersonalComputerVisitor):
    def _visit_cpu(self, element):
        print(f'CPU is in normal state')

    def _visit_motherboard(self, element):
        print('Motherboard is in normal state')
    
    def _visit_ram(self, element):
        print('RAM is in normal state')
    
    def _visit_psu(self, element):
        print('PSU is in normal state')
    
    def _visit_pc(self, element):
        print('PC is in normal state')


class PCElement(ElementABC):
    def accept(self, visitor: VisitorABC):
        visitor.visit(self)


class CPU(PCElement):
    ...


class Motherboard(PCElement):
    ...


class RAM(PCElement):
    ...


class GPU(PCElement):
    ...


class PSU(PCElement):
    ...


class PersonalComputer(ElementABC):
    '''It's a ObjectStructure that contains PC elements.
    '''
    def __init__(self, parts: list[PCElement]):
        self.__elements: list[PCElement] = parts

    def accept(self, visitor: VisitorABC):
        for element in self.__elements:
            element.accept(visitor)
        
        visitor.visit(self)


def main():
    assembly_validator = PCAssemblyValidator()
    operability_validator = PCOperabilityValidator()

    parts = [
        Motherboard(),
        CPU(),
        RAM(),
        RAM(),
        PSU(),
    ]

    pc = PersonalComputer(parts)
    pc.accept(assembly_validator)
    pc.accept(operability_validator)

    parts_for_second_pc = [
        Motherboard(),
        CPU(),
        RAM(),
        RAM(),
        RAM(),
        RAM(),
        GPU(),
        PSU(),
    ]

    second_pc = PersonalComputer(parts_for_second_pc)
    second_pc.accept(assembly_validator)


if __name__ == '__main__':
    main()
```
