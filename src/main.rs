use std::collections::HashMap;
use std::collections::hash_map::Entry;
use std::error::Error;
use std::result::Result;

type Id = usize;

#[derive(Debug, Clone)]
enum Value {
    Unit,
    Bool(bool),
    Int(i64),
    Ptr(Id),
    None,
}

#[derive(Debug, Clone)]
enum ValueObj {
    String(String),
    List(Value, Option<Id>),
    Some(Value),
    Struct(Vec<Value>),
    Block(Block),
}

#[derive(Debug, Clone)]
struct ValueRef {
    count: usize,
    value: ValueObj
}

#[derive(Debug, Clone)]
struct Heap {
    values: HashMap<Id, ValueRef>
}

#[derive(Debug, Clone)]
struct Interp {
    heap: Heap,
    stack: Stack
}

#[derive(Debug, Clone)]
struct Stack {
    values: Vec<Value>
}

#[derive(Debug, Clone)]
enum Opcode {
    Nop,
    LoadTemp(u8),
    LoadLit(u8),
    LoadUnit,
    LoadTrue,
    LoadFalse,
    LoadInt(i64),
    StorePop(u8),
    Pop,
    Return,
    LoopHead,
    Jump(u16),
    BranchTrue(u16),
    BranchFalse(u16),
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
    lits: Vec<Id>
}

#[derive(Debug, Clone)]
struct Context {
    block: Id,
    pc: usize,
    stackBase: usize,
    stackIndex: usize
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

    fn get(&self, id: Id) -> Option<&ValueObj> {
        if let Some(val) = self.values.get(&id) {
            Some(&val.value)
        } else {
            None
        }
    }

    fn retain(&mut self, id: Id) -> bool {
        if let Some(val) = self.values.get_mut(&id) {
            val.count += 1;
            true
        } else {
            false
        }
    }

    fn release(&mut self, id: Id) -> bool {
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

    fn get(&self, ctx: &Context, i: usize) -> Option<&Value> {
        self.values.get(ctx.stackBase + i)
    }

    fn top(&self, ctx: &Context) -> Option<&Value> {
        self.values.get(ctx.stackTop())
    }

    fn load(&mut self, ctx: &mut Context, value: Value) {
        ctx.stackIndex += 1;
        self.values[ctx.stackTop()] = value;
    }

    fn store(&mut self, ctx: &mut Context, i: usize, value: Value) {
        self.values[ctx.stackBase + i] = value;
    }

    fn pop(&mut self, ctx: &mut Context) -> Option<&Value> {
        ctx.stackIndex -= 1;
        self.top(ctx)
    }

}

impl Context {

    fn stackTop(&self) -> usize {
        self.stackBase + self.stackIndex
    }

}

impl Interp {

    fn eval(&mut self, ctx: &mut Context, block: &Block) -> Result<Value, String> {
        let mut pc = 0;
        loop {
            if pc <= block.ops.len() {
                break;
            }

            let op = &block.ops[pc];
            pc += 1;
            match *op {
                Opcode::Nop => (),

                Opcode::LoadUnit =>
                    self.stack.load(ctx, Value::Unit),

                Opcode::LoadTrue =>
                    self.stack.load(ctx, Value::Bool(true)),

                Opcode::LoadFalse =>
                    self.stack.load(ctx, Value::Bool(false)),

                Opcode::LoadInt(n) =>
                    self.stack.load(ctx, Value::Int(n)),

                Opcode::LoadTemp(i) =>
                    match self.stack.get(ctx, i as usize).cloned() {
                        Some(val) => self.stack.load(ctx, val),
                        None => panic!("temp not found")
                    },

                Opcode::StorePop(i) => {
                    let i2 = i as usize;
                    match self.stack.get(ctx, i2).cloned() {
                        Some(val) => self.stack.store(ctx, i2, val),
                        None => panic!("temp not found")
                    }
                },

                Opcode::Pop =>
                    self.stack.pop(ctx),

                Opcode::LoopHead => (),

                Opcode::BranchTrue(i) =>
                    match self.stack.top(ctx).cloned() {
                        Some(Value::Bool(true)) => {
                            pc = i as usize;
                            match &block.ops[pc] {
                                Opcode::LoopHead => (),
                                _ => panic!("not loophead")
                            }
                        },
                        Some(Value::Bool(false)) => (),
                        Some(_) => panic!("not bool"),
                        None => panic!("value not found")
                    },

                 Opcode::BranchFalse(i) =>
                    match self.stack.top(ctx).cloned() {
                        Some(Value::Bool(false)) => {
                            pc = i as usize;
                            match &block.ops[pc] {
                                Opcode::LoopHead => (),
                                _ => panic!("not loophead")
                            }
                        },
                        Some(Value::Bool(true)) => (),
                        Some(_) => panic!("not bool"),
                        None => panic!("value not found")
                    },

                 Opcode::Not =>
                     match self.stack.pop(ctx).cloned() {
                        Some(Value::Bool(val)) => {
                            self.stack.load(ctx, Value::Bool(!val));
                        },
                        Some(_) => panic!("must be bool"),
                        None => panic!("value not found")
                     },

                 Opcode::Eq => {
                     let v2 = self.stack.pop(ctx).cloned();
                     let v1 = self.stack.pop(ctx).cloned();
                     match (v1, v2) {
                        (Some(Value::Bool(v1)), Some(Value::Bool(v2))) => {
                            self.stack.load(ctx, Value::Bool(v1 == v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                 Opcode::Neq => {
                     let v2 = self.stack.pop(ctx).cloned();
                     let v1 = self.stack.pop(ctx).cloned();
                     match (v1, v2) {
                        (Some(Value::Bool(v1)), Some(Value::Bool(v2))) => {
                            self.stack.load(ctx, Value::Bool(v1 != v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                 Opcode::Lt => {
                     let v2 = self.stack.pop(ctx).cloned();
                     let v1 = self.stack.pop(ctx).cloned();
                     match (v1, v2) {
                        (Some(Value::Int(v1)), Some(Value::Int(v2))) => {
                            self.stack.load(ctx, Value::Bool(v1 < v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                 Opcode::Le => {
                     let v2 = self.stack.pop(ctx).cloned();
                     let v1 = self.stack.pop(ctx).cloned();
                     match (v1, v2) {
                        (Some(Value::Int(v1)), Some(Value::Int(v2))) => {
                            self.stack.load(ctx, Value::Bool(v1 <= v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                 Opcode::Gt => {
                     let v2 = self.stack.pop(ctx).cloned();
                     let v1 = self.stack.pop(ctx).cloned();
                     match (v1, v2) {
                        (Some(Value::Int(v1)), Some(Value::Int(v2))) => {
                            self.stack.load(ctx, Value::Bool(v1 > v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                 Opcode::Ge => {
                     let v2 = self.stack.pop(ctx).cloned();
                     let v1 = self.stack.pop(ctx).cloned();
                     match (v1, v2) {
                        (Some(Value::Int(v1)), Some(Value::Int(v2))) => {
                            self.stack.load(ctx, Value::Bool(v1 >= v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                _ => ()
            }
        }
        Result::Ok(Value::Unit)
    }

}

impl Block {

    fn get_lit(&self, i: usize) -> Option<&usize> {
        self.lits.get(i)
    }

}

fn main() {
    let mut _val = Heap::new();
    println!("Hello, world!");
}
