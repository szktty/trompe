use std::collections::HashMap;
use std::collections::hash_map::Entry;

#[derive(Debug, Clone)]
enum Value {
    Unit,
    Bool(bool),
    Int(i64),
    Ptr(u64),
    None,
}

#[derive(Debug, Clone)]
enum ValueObj {
    String(String),
    List(Value, Option<u64>),
    Some(Value),
    Struct(Vec<Value>)
}

#[derive(Debug, Clone)]
struct ValueRef {
    count: u64,
    value: ValueObj
}

#[derive(Debug, Clone)]
struct Heap {
    values: HashMap<u64, ValueRef>
}

#[derive(Debug, Clone)]
struct Interp {
    heap: Heap,
}

#[derive(Debug, Clone)]
enum Opcode {
    Nop,
    LoadTemp(u8),
    LoadLit(u8),
    LoadInt(i64),
    StorePop(u8),
    Pop,
    Return,
    LoopHead,
    Jump(u16),
    BranchFalse(u16),
    BranchTrue(u16),
    Apply(u8),
    Prim(u16),
    MakeBlock,
    Not,
    Eq,
    Neq,
    Lt,
    Le,
    Gt,
    Ge,
}

#[derive(Debug, Clone)]
struct Block {
    ops: Vec<Opcode>,
    lits: Vec<u64>
}

#[derive(Debug, Clone)]
struct Stack {
    values: Vec<Value>
}

impl Heap {

    fn new() -> Self {
        Heap {
            values: HashMap::new()
        }
    }

    fn new_value(&mut self, obj: ValueObj) -> Value {
        let id = 0;
        let ptr = Value::Ptr(id);
        self.values.insert(id, ValueRef { count: 1, value: obj });
        ptr
    }

    fn get(&self, id: u64) -> Option<&ValueObj> {
        if let Some(val) = self.values.get(&id) {
            Some(&val.value)
        } else {
            None
        }
    }

    fn retain(&mut self, id: u64) -> bool {
        if let Some(val) = self.values.get_mut(&id) {
            val.count += 1;
            true
        } else {
            false
        }
    }

    fn release(&mut self, id: u64) -> bool {
        if let Entry::Occupied(mut o) = self.values.entry(id) {
            if o.get().count <= 1 {
                o.remove_entry();
            } else {
                o.get_mut().count -= 1;
            }
            true
        } else {
            false
        }
    }

}

impl Stack {

    fn get(&self) -> Option<Value> {
    }

    fn push(&mut self, value: Value) {
        self.values.push(value);
    }

    fn pop(&mut self) {
        self.values.pop();
    }

}

impl Interp {

    fn eval(block: &Block) {
        let mut i = 0;
        loop {
            let op = &block.ops[i];
            i += 1;
            match *op {
                Opcode::Nop => (),
                _ => ()
            }
        }
    }

}

fn main() {
    let mut _val = Heap::new();
    println!("Hello, world!");
}
