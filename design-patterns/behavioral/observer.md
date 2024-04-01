# Observer 

The Observer pattern defines a one-to-many relationship between objects so that when the state of one object changes, all of its dependent objects are notified and updated automatically. Classes to whose events other classes subscribe are called Subjects, and subscribing classes are called Observers.

This pattern is used when:
- it is necessary to avoid a tight coupling between objects when declaring a one-to-many dependency
- unlimited number of dependent objects should be updated automatically when the state of one object changes
- an object can notify multiple other objects
- publishers are not interested in what the subscribers do with provided messages.

The pattern operates with four entities:
1. **Abstract Subject** (*Observable*) - interface that defines methods for adding, removing, and notifying observers
0. **Abstract Observer** - interface by which the observer is notified
0. **Concrete Subject** (*Observable*) - concrete class that implements *Subject* interface
0. **Concrete Observer** - concrete class that implements *Observer* interface

There are two ways to receive a state change request from a subject:

- **Pull** method: After receiving notification from the subject, the observer must contact the publisher and pull the data themselve. This method causes tight coupling between subject and observer and high avareness of observer about the subject. 
- **Push** method: The publisher does not notify the subscriber about data updates, but independently delivers (pushes) the data to the observer.

## Implementation examples

### Python

This example shows the *push method* of communication.

```python
from abc import ABC, abstractmethod


class AbstractSubject(ABC):
    '''AbstractSubject represents abstract base class of subject, that
    notifies observer about state changes.
    '''

    @abstractmethod
    def set_state(self, state):
        ...

    @abstractmethod
    def attach(self, observer):
        ...

    @abstractmethod
    def detach(self, observer):
        ...

    @abstractmethod
    def notify(self):
        ...


class AbstractObserver(ABC):
    '''AbstractObserver 
    '''
    @abstractmethod
    def update(self, data):
        ...


class Sensor(AbstractSubject):
    def __init__(self):
        self.__state = None
        self.__observers: map[int, AbstractObserver] = {}
    
    def set_state(self, state):
        self.__state = state

    def attach(self, observer):
        self.__observers[id(observer)] = observer

    def detach(self, observer):
        del(self.__observers[id(observer)])

    def notify(self):
        for observer in self.__observers.values():
            observer.update(self.__state)


class AlertingObserver(AbstractObserver):
    def __init__(self):
        ...
    
    def update(self, data):
        if data > 90:
            print(f"Sensor temperature is too high: {data}")



class MonitoringObserver(AbstractObserver):
    def __init__(self):
        self.__state = None

    def _change_state(self, state):
        self.__state = state
    
    def update(self, data):
        if self.__state == None:
            print(f"Got new data: {data}")
            self._change_state(data)
        else:
            print(f"Data changed from {self.__state} to {data}")
            self._change_state(data)


def main():
    water_temperature = Sensor()
    alerting = AlertingObserver()
    monitoring = MonitoringObserver()

    water_temperature.attach(alerting)
    water_temperature.attach(monitoring)

    for i in range(30, 120, 10):
        water_temperature.set_state(i)
        water_temperature.notify()


if __name__ == '__main__':
    main()
```

## Notes

- **Observer** is one of the patterns that uses publish/subscribe architectural pattern.
