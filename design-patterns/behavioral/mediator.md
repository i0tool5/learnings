# Mediator 
This pattern presents a *Mediator* object that hides the way multiple other objects (*Colleagues*) interact. Mediator makes the system loosely coupled by eliminating the need for objects to reference each other, allowing to change the interaction between them independently.

This pattern is used when:
- tight coupling between a set of interacting objects should be avoided
- it should be possible to change the interaction between a set of objects independently

The pattern operates with two entities:
1. **Mediator** - an entity that defines the interface for communication between *Colleague* objects
2. **Colleague** - colleague describing organization of the process of interaction of colleague objects with an object of the Mediator type

## Implementation examples

### Python

```python
import inspect
from abc import ABC, abstractmethod
from typing import Any
from weakref import proxy


class Message:
    def __init__(self, msg: str, from_: Any, to: Any = None):
        self.sender: Any = from_
        self.receiver: Any = to
        self._msg: str = msg

    def __str__(self) -> str:
        return self._msg


class MediatorABC(ABC):
	@abstractmethod
	def send(self, message: Message) -> None:
		...


class ColleagueABC(ABC):
	def __init__(self, mediator: MediatorABC) -> None:
		self._mediator = proxy(mediator)

	@abstractmethod
	def send(self, message: Any) -> None:
		...

	@abstractmethod
	def receive(self, message: Any) -> None:
		...


class Mediator(MediatorABC):
	def __init__(self) -> None:
		self._colleagues: list[ColleagueABC] = []

	def add(self, colleague: ColleagueABC):
		self._colleagues.append(colleague)

	def send(self, message: Message) -> None:
		for colleague in self._colleagues:
			if colleague != message.sender:
				colleague.receive(message)


class Colleague(ColleagueABC):
	def __init__(self, name: str, mediator: MediatorABC):
		super().__init__(mediator=mediator)
		self._name = name

	def send(self, message: Any) -> None:
		self._mediator.send(Message(message, self))

	def receive(self, message: Any) -> None:
		print(f'{self._name} got a message: {message}')


if __name__ == '__main__':
	chat = Mediator()
	alice = Colleague('Alice', chat)
	bob = Colleague('Bob', chat)
	eve = Colleague('Eve', chat)
	chat.add(alice)
	chat.add(bob)
	chat.add(eve)
	alice.send("Let's start the show?!")
	bob.send('Oh yeah!')
	eve.send("No, i'm not ready!")
```
