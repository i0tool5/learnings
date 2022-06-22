pub mod reflib {
    use std::fmt;
    use std::cell::RefCell;
    use std::rc::{Rc,Weak};

    // vector of weak references to cell reference to person
    /// Family is a collection of Persons
    pub type Family = Vec<Weak<RefCell<Person>>>;

    #[derive(Debug)]
    pub struct Person {
        age: u8,
        name: String,
        parents: Family,
        children: Family
    }

    impl Clone for Person {
        fn clone(&self) -> Self {
            Person {
                age: self.age,
                name: self.name.clone(),
                parents: self.parents.clone(),
                children: self.children.clone()
            }
        }
    }

    impl fmt::Display for Person {
        fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
            write!(f, "Person <{}, {}>", self.name, self.age)
        }
    }

    impl Person {
        /// new person with given age and name
        pub fn new(age: u8, name: String) -> Person {
            let p: Family = vec![];
            let c : Family = vec![];
            Person{
                age,
                name,
                parents: p,
                children: c
            }
        }

        /// add child to the person
        pub fn add_child(&mut self, c: Weak<RefCell<Person>>) {
            self.children.push(c);
        }

        /// add parent to the person
        pub fn add_parent(&mut self, p: Weak<RefCell<Person>>) {
            self.parents.push(p);
        }

        /// get all children
        pub fn get_children(&self) -> Vec<Rc<RefCell<Person>>> {
            let mut v = Vec::with_capacity(self.children.len());
            for el in &self.children {
                let val = el.upgrade();
                if val.is_some() {
                    v.push(val.unwrap());
                }
            }
            return v;
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::reflib::Person;

    #[test]
    fn new_person() {
        let person = Person::new(18, String::from("Alex"));
        assert_eq!(format!("{}", person), String::from("Person <Alex, 18>"));
    }
}
