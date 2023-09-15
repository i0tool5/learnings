# Beahvioral patterns

**Behavioral patterns** describe the interaction of objects and classes with each other and try to achieve the least degree of coupling of system components with each other, making the system more flexible.

There are two types of Behavioral Patterns:
1. ***Class-level patterns***
2. ***Object-level patterns***

---

***Class-level patterns*** describe the interactions between classes and their subclasses. Such relationships are expressed through class inheritance and implementation. Here, the base class defines the interface, and subclasses define the implementation.

> Only “*Template Method*” belongs to class level patterns.

---

***Object-level patterns*** describe interactions between objects. Such relations are expressed by connections - **association**, **aggregation** and **composition**. Here structures are built by combining objects of some classes.

- **Association** is a relation when objects of two classes can refer one to another. For example, a property of a class contains an instance of another class.

- **Aggregation** is a *particular form of association*. Aggregation is used when one object must be a container for other objects and lifetime of these objects doesn't depend in any way on the lifetime of the container object. In general, if the container is destroyed, the objects included in it will not be affected. For example, an object was created, and then was passed to a container object (to a method of the container object or assigned to a container property from outside). When a container is deleted, the created object, which can interact with other containers, will not be affected in any way.

- **Composition** is the same as aggregation, but composite objects cannot exist separately from the container object and if the container is destroyed, then all its contents will be destroyed too. For example, an object was created in a container object method and assigned to a container property. From the outside, no one knows about the created object, which means that when the container is deleted, the created object will be deleted in the same way, because there is no external reference to it.

## List of Behavioral patters

- [**Memento**](memento.md)