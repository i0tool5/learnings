# Chain of responsibility
This pattern allows to chain recipient objects together and pass the request along this chain until it is processed. Each handler in the chain decides either to process the request or to pass it along the chain to the next handler.

This pattern is used when:
- sender of request and its receiver should be decoupled
- request can be handled by more than one object, and the actual handler is not known in advance and must be found automatically
- request can be sent to one of several objects, without explicitly specifying which one
- set of objects capable of processing a request can be defined dynamically.

The pattern operates with three entities:
1. **Handler** - defines an interface for handling requests and (optionally) implements a relationship with the successor
2. **ConcreteHandler** - implements **Handler**. ConcreteHandler may process the request or pass it on to the next handler in the chain if it cannot handle the request itself, or, in some situations, it must pass it on to the next handler even after successful processing.
3. **Client** - 

## Implementation examples

The most common example of the *Chain-of-Responsibility* pattern is **HTTP middleware**.

Key Aspects of this relationship:
- Sequential Handling: each middleware receives the request, performs an action (like authentication, logging, or validation), and decides whether to pass it to the next component.
- Early Termination: if middleware determines the request is invalid, it can stop the chain entirely without calling the next handler.
- Loose Coupling: middleware components are independent, allowing them to be easily reordered, added, or removed.

### Go

```go
package main

import "fmt"

// Request struct holds the details of a request
type Request struct {
    Amount float64
    Purpose string
}

// Approver interface defines the handler and the ability to set the next handler
type Approver interface {
    ProcessRequest(req Request)
    SetNext(next Approver)
}

// Manager concrete handler
type Manager struct {
    next Approver
}

func (m *Manager) ProcessRequest(req Request) {
    if req.Amount <= 1000 {
        fmt.Printf("Manager approved request for %s (Amount: %.2f)\n", req.Purpose, req.Amount)
    }
    fmt.Println("Manager cannot approve. Passing to next.")
    if m.next != nil {
        m.next.ProcessRequest(req)
    }
}

func (m *Manager) SetNext(next Approver) {
    m.next = next
}

// Director concrete handler
type Director struct {
    next Approver
}

func (d *Director) ProcessRequest(req Request) {
    if req.Amount <= 10000 {
        fmt.Printf("Director approved request for %s (Amount: %.2f)\n", req.Purpose, req.Amount)
    }
    fmt.Println("Director cannot approve. Passing to next.")
    if d.next != nil {
        d.next.ProcessRequest(req)
    }
}

func (d *Director) SetNext(next Approver) {
    d.next = next
}

// President concrete handler
type President struct {
    next Approver
}

func (p *President) ProcessRequest(req Request) {
    if req.Amount <= 100000 {
        fmt.Printf("President approved request for %s (Amount: %.2f)\n", req.Purpose, req.Amount)
        return
    }

    fmt.Printf("Request for %s (Amount: %.2f) requires an executive meeting, cannot be approved by President alone.\n", req.Purpose, req.Amount)
    if p.next != nil {
        p.next.ProcessRequest(req)
    }
}

func (p *President) SetNext(next Approver) {
    p.next = next
}

func main() {
    // Set up the chain
    manager := &Manager{}
    director := &Director{}
    president := &President{}

    manager.SetNext(director)
    director.SetNext(president)

    request1 := Request{Amount: 500, Purpose: "Office Supplies"}
    request2 := Request{Amount: 5000, Purpose: "New Equipment"}
    request3 := Request{Amount: 500000, Purpose: "New Company Car"}
    
    requests := []Request{
        request1,
        request2,
        request3,
    }

    // Process requests
    for _, request := range requests {
        manager.ProcessRequest(request)
    }
}
```
