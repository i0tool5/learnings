# Command
The command pattern allows to represent a request as a stand-alone object. This object will encapsulate an action and parameters required for prerforming this action.

This pattern is used when:
- coupling of the sender of the request and the receiver should be avoided. 
- the invoker of the request should be possible to configure by the request

The pattern operates with four entities:
1. **Command interface** - describes a  common method `execute()` that every command should implement.
2. **Concrete Command Class** - implements specific command. Each class encapsulates the details of a particular action, which should be executed by receiver.
3. **Invoker Class** *(Remote Control)* - coordinates and executes commands. Invoker is able to store command and execute (or cancel queued) them by some event.
4. **Receiver Class** - implements the methods required to perform the action associated with the command.

## Implementations

### Python

This example was taken from [Wikipedia](https://ru.wikipedia.org/wiki/%D0%9A%D0%BE%D0%BC%D0%B0%D0%BD%D0%B4%D0%B0_(%D1%88%D0%B0%D0%B1%D0%BB%D0%BE%D0%BD_%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D1%8F)#%D0%9F%D1%80%D0%B8%D0%BC%D0%B5%D1%80_%D0%BD%D0%B0_Python) and will be changed in the future.
```python
from abc import ABCMeta, abstractmethod


class Troop:
	"""
	Receiver - an object representing a military detachment
	"""

	def move(self, direction: str) -> None:
		"""
		Start moving in a certain direction
		"""
		print('The squad started moving {}'.format(direction))

	def stop(self) -> None:
		"""
		Stop moving
		"""
		print('The squad stopped moving')


class Command(metaclass=ABCMeta):
	"""
	Abstract base class for all commands 
	"""

	@abstractmethod
	def execute(self) -> None:
		"""
		Proceed to execute the command
		"""
		pass

	@abstractmethod
	def unexecute(self) -> None:
		"""
		Cancel command execution
		"""
		pass	


class AttackCommand(Command):
	"""
	The command to conduct an attack
	"""

	def __init__(self, troop: Troop) -> None:
		"""
		Constructor.

		:param troop: the squad the team is associated with
		"""
		self.troop = troop

	def execute(self) -> None:
		self.troop.move('forward')

	def unexecute(self) -> None:
		self.troop.stop()


class RetreatCommand(Command):
	"""
	The command to conduct a retreat
	"""

	def __init__(self, troop: Troop) -> None:
		"""
		Constructor.
		
		:param troop: the squad the team is associated with
		"""
		self.troop = troop

	def execute(self) -> None:
		self.troop.move('backward')

	def unexecute(self) -> None:
		self.troop.stop()


class TroopInterface:
	"""
	Invoker - an interface through which commands can be given to a certain squad
	"""

	def __init__(self, attack: AttackCommand, retreat: RetreatCommand) -> None:
		"""
		Constructor.

		:param attack: command to conduct an attack
		:param retreat: command to conduct a retreat
		"""
		self.attack_command = attack
		self.retreat_command = retreat
		self.current_command = None		# command that currently in progress

	def attack(self) -> None:
		self.current_command = self.attack_command
		self.attack_command.execute()

	def retreat(self) -> None:
		self.current_command = self.retreat_command
		self.retreat_command.execute()

	def stop(self) -> None:
		if self.current_command:
			self.current_command.unexecute()
			self.current_command = None
		else:
			print("Squad can't be stopped, because it is not moving')


if __name__ == '__main__':
	troop = Troop()
	interface = TroopInterface(AttackCommand(troop), RetreatCommand(troop))
	interface.attack()
	interface.stop()
	interface.retreat()
	interface.stop()
```
