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
    LoadTemp(String),
    LoadLit(u8),
    LoadUnit,
    LoadTrue,
    LoadFalse,
    LoadInt(i64),
    StorePop(String),
    Pop,
    Return,
    LoopHead,
    Jump(u16),
    BranchTrue(u16),
    BranchFalse(u16),
    Apply(u8),
    Prim(String),
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
    lits: Vec<Id>,
    env: Env
}

#[derive(Debug, Clone)]
struct Context {
    pc: usize,
    stackBase: usize,
    stackIndex: usize
}

#[derive(Debug, Clone)]
struct Env {
    attrs: HashMap<String, Value>
}

impl Env {

    fn new() -> Env {
        Env { attrs: HashMap::new() }
    }

    fn get(&self, key: &String) -> Option<Value> {
        self.attrs.get(key).cloned()
    }

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

    fn retain_id(&mut self, id: Id) -> bool {
        if let Some(val) = self.values.get_mut(&id) {
            val.count += 1;
            true
        } else {
            false
        }
    }

    fn retain(&mut self, value: &Value) -> bool {
        match *value {
            Value::Ptr(id) => self.retain_id(id),
            _ => true
        }
    }

    fn release_id(&mut self, id: Id) -> bool {
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

    fn release(&mut self, value: &Value) -> bool {
        match *value {
            Value::Ptr(id) => self.release_id(id),
            _ => true
        }
    }

}

impl Stack {

    fn new() -> Stack {
        Stack { values: Vec::with_capacity(1000) }
    }

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

    fn new() -> Context {
        Context { pc: 0, stackBase: 0, stackIndex: 0 }
    }

    fn stackTop(&self) -> usize {
        self.stackBase + self.stackIndex
    }

}

impl Interp {

    fn new() -> Interp {
        Interp { heap: Heap::new(), stack: Stack::new() }
    }

    fn eval(&mut self, ctx: &mut Context, env: &mut Env, block: &Block) -> Result<Value, String> {
        let mut pc = 0;
        loop {
            if pc <= block.ops.len() {
                break;
            }

            let op = &block.ops[pc];
            pc += 1;
            match op {
                Opcode::Nop => (),

                Opcode::LoadUnit =>
                    self.stack.load(ctx, Value::Unit),

                Opcode::LoadTrue =>
                    self.stack.load(ctx, Value::Bool(true)),

                Opcode::LoadFalse =>
                    self.stack.load(ctx, Value::Bool(false)),

                Opcode::LoadInt(n) =>
                    self.stack.load(ctx, Value::Int(*n)),

                Opcode::LoadTemp(name) =>
                    match env.attrs.get(name) {
                        Some(val) => {
                            self.heap.retain(&val);
                            self.stack.load(ctx, val.clone());
                        },
                        None => panic!("temp not found")
                    },

                Opcode::StorePop(name) => {
                    match self.stack.pop(ctx) {
                        Some(val) => {
                            match env.attrs.insert(name.clone(), val.clone()) {
                                None => (),
                                Some(old) => {
                                    let _ = self.heap.release(&old);
                                }
                            }
                        },
                        None => panic!("value not found")
                    }
                },

                Opcode::Pop => {
                    let _ = self.stack.pop(ctx);
                },

                Opcode::LoopHead => (),

                Opcode::BranchTrue(i) =>
                    match self.stack.top(ctx).cloned() {
                        Some(Value::Bool(true)) => {
                            pc = *i as usize;
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
                            pc = *i as usize;
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

    fn new() -> Block {
        Block { ops: Vec::new(), lits: Vec::new(), env: Env::new() }
    }

    fn get_lit(&self, i: usize) -> Option<&usize> {
        self.lits.get(i)
    }

}

fn main() {
    println!("Hello, world!");
    let mut interp = Interp::new();
    let mut ctx = Context::new();
    let mut env = Env::new();
    let mut block = Block::new();
    block.ops = vec![
        Opcode::Nop
    ];
    interp.eval(&mut ctx, &mut env, &block);
}
