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