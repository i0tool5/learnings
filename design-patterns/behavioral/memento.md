# Memento
Memento - is a behavioral design pattern that allows, without breaking encapsulation, to capture and store the internal state of an object so that it can later be restored to that state. When using this pattern, care should be taken if the originator may change other objects or resources—the memento pattern operates on a single object.

This pattern is used when:
- it is necessary to save a snapshot of the object state (or part of it) for later restoration
- a direct interface for obtaining object state reveals implementation details and breaks object encapsulation

The pattern operates with three entities:

1. **Originator** - an entity that has some changing state, as well as it can create and receive Memento of its state
2. **Memento** - stores the state of the Originator class object in itself
3. **Caretaker** - responsible for keeping the Memento until it is needed by the Originator

## Implementations

### Golang

```go
// Package memento is an example of the Memento Pattern.
package memento

import (
	"encoding/json"
	"testing"
)

// Originator interface for Memento pattern.
type Originator interface {
	CreateMemento() (Memento, error)
	ApplyMemento(Memento) error
}

// Profile is a specific implementation of Originator.
type Profile struct {
	Name     string `json:"name"`
	Password string `json:"passwod"`
	Views    int64  `json:"views"`
	Locked   bool   `json:"locked"`
}

// compiler-time check.
var _ Originator = (*Profile)(nil)

// CreateMemento returns state storage.
func (p *Profile) CreateMemento() (Memento, error) {
	state, err := json.Marshal(p)
	if err != nil {
		return Memento{}, err
	}
	return Memento{state: state}, nil
}

// ApplyMemento applies the old state.
func (p *Profile) ApplyMemento(memento Memento) error {
	return json.Unmarshal(memento.State(), p)
}

// Memento implements storage for the state of Originator
type Memento struct {
	state []byte
}

// GetState returns state.
func (m *Memento) State() []byte {
	return m.state
}

// Caretaker keeps Memento until it is needed by Originator.
type Caretaker struct {
	memento *Memento
}

// Save memento.
func (c *Caretaker) Save(m Memento) {
	c.memento = &m
}

// Retrieve memento.
func (c *Caretaker) Retrieve() Memento {
	return *c.memento
}

// NewCaretaker creates Memento store.
func NewCaretaker() *Caretaker {
	return &Caretaker{}
}

func TestMemento(t *testing.T) {
	profile := Profile{
		Name:     "Username",
		Password: "encryptedpassword",
		Views:    500,
		Locked:   false,
	}

	profileMemento, err := profile.CreateMemento()
	if err != nil {
		panic(err)
	}
	ct := NewCaretaker()
	ct.Save(profileMemento)

	profile.Locked = true
	profile.Views = -1

	err = profile.ApplyMemento(ct.Retrieve())
	if err != nil {
		t.Error(err)
	}

	if profile.Locked {
		t.Error("should be restored")
	}

	if profile.Views == -1 {
		t.Error("should be restored")
	}
}
```

### Python

```python
import abc
import hashlib
import pickle
import typing


class Memento:
    def __init__(self, state):
        self.__state = state

    def get_state(self) -> typing.Any:
        return self.__state


class Caretaker:
    def __init__(self,):
        self.__memento: Memento

    def save(self, memento: Memento):
        self.__memento = memento

    def retrieve(self) -> Memento:
        return self.__memento


class OriginatorABC(abc.ABC):
    @abc.abstractmethod
    def create_memento(self) -> Memento:
        ...

    @abc.abstractmethod
    def apply_memento(self, memento: Memento):
        ...


class Originator(OriginatorABC):
    def create_memento(self) -> Memento:
        state = self._create_state()
        return Memento(state)

    def apply_memento(self, memento: Memento):
        state = memento.get_state()
        self._apply_state(state)

    def _apply_state(self, state: typing.Any):
        loaded = pickle.loads(state)
        self = loaded

    def _create_state(self):
        return pickle.dumps(self)


class Profile(Originator):
    def __init__(
            self,
            name: str,
            password: str
        ) -> None:

        self.name: str = name
        self.__password: str = hashlib.md5(password.encode()).hexdigest()
        self.__views: int = 1
        self.__locked: bool = False
    
    def lock(self):
        self.__locked = True
        self.__views = -1
    
    def is_locked(self) -> bool:
        return self.__locked
    
    def get_views(self):
        return self.__views

    def compare_password(self, passwd: str) -> bool:
        return self.__password == hashlib.md5(passwd.encode()).hexdigest()


def main():
    caretaker = Caretaker()
    user_profile = Profile("i0tool", "so strong")
    caretaker.save(user_profile.create_memento())

    user_profile.lock()

    assert user_profile.is_locked(), True
    assert user_profile.get_views(), -1

    user_profile.apply_memento(caretaker.retrieve())
    assert user_profile.name, "i0tool"
    assert user_profile.is_locked(), False
    assert user_profile.get_views(), 1
    assert user_profile.compare_password("so strong")

if __name__ == '__main__':
    main()
```
