use std::rc::Rc;
use std::cell::RefCell;

#[derive(Clone, Debug)]
enum Type {
    App(Rc<Type>, Vec<Tycon>),
    Assoc { name: String, ty: Rc<Type>, default: Option<Rc<Type>> },
    Meta(RefCell<Option<Rc<Type>>>),
}

#[derive(Clone, Debug)]
enum Tycon {
    Unit,
    Bool,
    Int,
    Float,
    Option,
    Char,
    String,
    List,
    Map,
    Bytes,
    File,
    Struct,
    Enum,
    Intf
}

struct Infer {
}

impl Infer {

    fn infer(&self, ex: &mut Type, ac: &mut Type) -> Result<(), String> {
        /*
        if ex == ac {
            return Ok(());
        }
        */
        match (ex, ac) {
            (Type::Meta(ref cell1), Type::Meta(ref cell2)) => {
                match (cell1.borrow_mut().clone(), cell2.borrow_mut().clone()) {
                    (Some(ref mut ty1), Some(ref mut ty2)) =>
                        self.infer(Rc::make_mut(ty1), Rc::make_mut(ty2)),
                    (Some(ty), None) => {
                        cell2.replace(Some(ty));
                        Ok(())
                    },
                    (None, Some(ty)) => {
                        cell1.replace(Some(ty));
                        Ok(())
                    },
                    (None, None) => panic!("meta and meta")
                }
            },
            (Type::Meta(ref cell1), ref ty) => {
                cell1.replace(Some(Rc::new((*ty).clone())));
                Ok(())
            },
            _ => return Ok(())
        }
    }

}
