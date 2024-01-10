# State
State pattern allows an object to alter its behavior when its internal state changes. This pattern is close to the concept of finite-state machines. The object's behavior changes so much that it appears as if the object's class has changed.

This pattern is used when:
- an object should change its behavior when its internal state changes
- the object's behavior must change during program execution
- state-specific behavior should be defined independently. That is, adding new states should not affect the behavior of existing states

The pattern operates with three entities:
1. **Context** - is an object-oriented representation of a state machine.
2. **Abstract State** - defines interface of various states
3. **State implementations** - each of implementations implements one of behaviors associated with a certain state

## Implementations

```python
from collections.abc import ABC, abstractmethod


class State(ABC):
    ...


class CoolingState(State):
    def __init__(self, temperature_regulator):
        self.regulator: temperature_regulator
        self._avail_temperatures = [8, 10, 12, 14]
        self._current_temp = self._avail_temperatures[-1]
    


    def toggle_mode(self):
        self.regulator.state = self.regulator.heat_state


class HeatState(State):
    def toggle_mode(self):
        ...


class TemperatureRegulator:
    '''Conditioner has on/off button and mode toggle switch'''
    def __init__(self):
        self.cooling_state = CoolingState(self)
        self.heat_state = HeatState(self)

    def toggle_mode(self):
        ...


if __name__ == '__main__':
    conditioner = TemperatureRegulator()

```
