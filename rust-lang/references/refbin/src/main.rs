use std::cell::RefCell;
use std::rc::Rc;

use reflib::reflib::*;

fn main() {
    let andrew = Rc::new(RefCell::new(Person::new(40, String::from("Andrew"))));
    let matt = Rc::new(RefCell::new(Person::new(18, String::from("Matt"))));

    andrew.borrow_mut().add_child(Rc::downgrade(&matt));
    matt.borrow_mut().add_parent(Rc::downgrade(&andrew));

    println!("Andrew {:?}", andrew.borrow());
    println!("{:?}", andrew.borrow().get_children())
}

