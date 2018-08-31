use std::collections::HashMap;
use std::collections::hash_map::Entry;
use std::error::Error;
use std::result::Result;
use std::string::ToString;
use std::io::Write;

mod infer;

type Id = usize;

#[derive(Debug, Clone)]
enum Value {
    Unit,
    Bool(bool),
    Int(i64),
    Ptr(Id),
    Prim(String),
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
    next_id: usize,
    values: HashMap<Id, ValueRef>
}

#[derive(Clone)]
struct Prim {
    name: String,
    f: fn(interp: &Interp, ctx: &Context, nargs: u8, args: Vec<Value>)
}

#[derive(Clone)]
struct Interp {
    heap: Heap,
    stack: Stack,
    prims: HashMap<String, Prim>
}

#[derive(Debug, Clone)]
struct Stack {
    count: usize,
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
    LoadPrim(String),
    StorePop(String),
    Pop,
    Return,
    LoopHead,
    Jump(i16),
    BranchTrue(u16),
    BranchFalse(u16),
    Apply(u8),
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
struct CompiledCode {
    ops: Vec<Opcode>,
    lits: Vec<Id>,
}

#[derive(Debug, Clone)]
struct Block {
    code: CompiledCode,
    env: Env
}

#[derive(Debug, Clone)]
struct Context {
    pc: usize,
    base: usize,
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
            next_id: 0,
            values: HashMap::new()
        }
    }

    fn new_value(&mut self, obj: ValueObj) -> Value {
        let id = self.next_id;
        let ptr = Value::Ptr(id);
        self.next_id += 1;
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
        Stack { count: 0, values: Vec::with_capacity(1000) }
    }

    fn get(&self, i: usize) -> Option<&Value> {
        self.values.get(i)
    }

    fn top(&self) -> Option<&Value> {
        self.values.get(self.count-1)
    }

    fn load(&mut self, value: Value) {
        if self.count >= self.values.len() {
            self.values.push(value);
        } else {
            self.values[self.count-1] = value;
        }
        self.count += 1;
    }

    fn store(&mut self, i: usize, value: Value) {
        self.values[i] = value;
    }

    fn pop(&mut self) -> Option<&Value> {
        self.count -= 1;
        self.values.get(self.count)
    }

}

impl Context {

    fn new() -> Context {
        Context { pc: 0, base: 0 }
    }

}

impl Interp {

    fn new() -> Interp {
        Interp { heap: Heap::new(), stack: Stack::new(),
        prims: HashMap::new() }
    }

    fn eval(&mut self, ctx: &mut Context, env: &mut Env, code: &CompiledCode) -> Result<Value, String> {
        let mut pc = 0;
        loop {
            if pc >= code.ops.len() {
                break;
            }

            println!("# opcode => {}", &code.ops[pc].to_string());
            match &code.ops[pc] {
                Opcode::Nop => (),

                Opcode::LoadUnit =>
                    self.stack.load(Value::Unit),

                Opcode::LoadTrue =>
                    self.stack.load(Value::Bool(true)),

                Opcode::LoadFalse =>
                    self.stack.load(Value::Bool(false)),

                Opcode::LoadInt(n) =>
                    self.stack.load(Value::Int(*n)),

                Opcode::LoadTemp(ref name) =>
                    match env.attrs.get(name) {
                        Some(val) => {
                            self.heap.retain(&val);
                            self.stack.load(val.clone());
                        },
                        None => panic!("temp not found")
                    },

                Opcode::LoadLit(i) => {
                    match code.lits.get(*i as usize) {
                        None => panic!("literal not found"),
                        Some(id) => self.stack.load(Value::Ptr(*id))
                    }
                },

                Opcode::LoadPrim(ref name) =>
                    self.stack.load(Value::Prim(name.clone())),

                Opcode::StorePop(name) => {
                    match self.stack.pop() {
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
                    let _ = self.stack.pop();
                },

                Opcode::LoopHead => (),

                Opcode::BranchTrue(i) =>
                    match self.stack.top().cloned() {
                        Some(Value::Bool(true)) => {
                            pc = *i as usize;
                            match &code.ops[pc] {
                                Opcode::LoopHead => (),
                                _ => panic!("not loophead")
                            }
                        },
                        Some(Value::Bool(false)) => (),
                        Some(_) => panic!("not bool"),
                        None => panic!("value not found")
                    },

                 Opcode::BranchFalse(i) =>
                    match self.stack.top().cloned() {
                        Some(Value::Bool(false)) => {
                            pc = *i as usize;
                            match &code.ops[pc] {
                                Opcode::LoopHead => (),
                                _ => panic!("not loophead")
                            }
                        },
                        Some(Value::Bool(true)) => (),
                        Some(_) => panic!("not bool"),
                        None => panic!("value not found")
                    },

                 Opcode::Apply(_) => {
                     match self.stack.pop() {
                         Some(Value::Prim(ref name)) => {
                             match self.prims.get(name) {
                                 Some(prim) => {
                                     /*
                                     let f = prim.f;
                                     f(self, ctx, 0, Vec::new())
                                     */
                                 },
                                 _ => (),
                             }
                         },
                         _ => panic!("not function")
                     }
                 },

                 Opcode::Not =>
                     match self.stack.pop().cloned() {
                        Some(Value::Bool(val)) => {
                            self.stack.load(Value::Bool(!val));
                        },
                        Some(_) => panic!("must be bool"),
                        None => panic!("value not found")
                     },

                 Opcode::Eq => {
                     let v2 = self.stack.pop().cloned();
                     let v1 = self.stack.pop().cloned();
                     match (v1, v2) {
                        (Some(Value::Bool(v1)), Some(Value::Bool(v2))) => {
                            self.stack.load(Value::Bool(v1 == v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                 Opcode::Neq => {
                     let v2 = self.stack.pop().cloned();
                     let v1 = self.stack.pop().cloned();
                     match (v1, v2) {
                        (Some(Value::Bool(v1)), Some(Value::Bool(v2))) => {
                            self.stack.load(Value::Bool(v1 != v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                 Opcode::Lt => {
                     let v2 = self.stack.pop().cloned();
                     let v1 = self.stack.pop().cloned();
                     match (v1, v2) {
                        (Some(Value::Int(v1)), Some(Value::Int(v2))) => {
                            self.stack.load(Value::Bool(v1 < v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                 Opcode::Le => {
                     let v2 = self.stack.pop().cloned();
                     let v1 = self.stack.pop().cloned();
                     match (v1, v2) {
                        (Some(Value::Int(v1)), Some(Value::Int(v2))) => {
                            self.stack.load(Value::Bool(v1 <= v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                 Opcode::Gt => {
                     let v2 = self.stack.pop().cloned();
                     let v1 = self.stack.pop().cloned();
                     match (v1, v2) {
                        (Some(Value::Int(v1)), Some(Value::Int(v2))) => {
                            self.stack.load(Value::Bool(v1 > v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                 Opcode::Ge => {
                     let v2 = self.stack.pop().cloned();
                     let v1 = self.stack.pop().cloned();
                     match (v1, v2) {
                        (Some(Value::Int(v1)), Some(Value::Int(v2))) => {
                            self.stack.load(Value::Bool(v1 >= v2));
                        },
                        (Some(_), Some(_)) => panic!("must be bool"),
                        _ => panic!("value not found")
                     }
                 },

                _ => ()
            }
            pc += 1;
        }
        Result::Ok(Value::Unit)
    }

}

impl CompiledCode {

    fn new() -> CompiledCode {
        CompiledCode { ops: Vec::new(), lits: Vec::new() }
    }

    fn get_lit(&self, i: usize) -> Option<&usize> {
        self.lits.get(i)
    }

}

impl ToString for CompiledCode {

    fn to_string(&self) -> String {
        let mut buf = Vec::new();
        write!(&mut buf, "opcodes:\n").unwrap();
        let mut i = 0;
        for op in self.ops.iter() {
            i += 1;
            write!(&mut buf, "    {}: {}\n", i, op.to_string());
        }
        String::from_utf8(buf).unwrap()
    }

}

impl Block {

    fn new(code: CompiledCode) -> Block {
        Block { code: code, env: Env::new() }
    }

}

impl ToString for Opcode {

    fn to_string(&self) -> String {
        match self {
            Opcode::Nop => String::from("nop"),
            Opcode::LoadTemp(name) => format!("load temp #{}", name),
            Opcode::LoadLit(i) => format!("load literal #{}", i),
            Opcode::LoadUnit => format!("load unit"),
            Opcode::LoadTrue => format!("load true"),
            Opcode::LoadFalse => format!("load false"),
            Opcode::LoadInt(val) => format!("load int {}", val),
            Opcode::LoadPrim(name) => format!("load primitive #{}", name),
            Opcode::StorePop(name) => format!("store #{}; pop", name),
            Opcode::Pop => format!("pop"),
            Opcode::Return => format!("return"),
            Opcode::LoopHead => format!("loophead"),
            Opcode::Jump(n) => format!("jump {}", n),
            Opcode::BranchTrue(n) => format!("branch true {}", n),
            Opcode::BranchFalse(n) => format!("branch false {}", n),
            Opcode::Apply(n) => format!("apply {}", n),
            Opcode::MakeBlock => format!("make block"),
            Opcode::Not => format!("not"),
            Opcode::Eq => format!("=="),
            Opcode::Neq => format!("!="),
            Opcode::Lt => format!("<"),
            Opcode::Le => format!("<="),
            Opcode::Gt => format!(">"),
            Opcode::Ge => format!(">="),
        }
    }

}

fn prim_hello(interp: &Interp, ctx: &Context, nargs: u8, args: Vec<Value>) {
    println!("# prim: hello, world!");
}

fn main() {
    println!("Hello, world!");
    let mut interp = Interp::new();
    let prim_name = String::from("hello");
    interp.prims.insert(prim_name.clone(),
    Prim { name: prim_name.clone(), f: prim_hello });
    let mut ctx = Context::new();
    let mut env = Env::new();
    let mut code = CompiledCode::new();
    code.ops = vec![
        Opcode::Nop,
        Opcode::LoadUnit,
        Opcode::LoadPrim(prim_name.clone()),
        Opcode::Apply(0),
    ];
    println!("code => {}", code.to_string());
    let mut block = Block::new(code);
    interp.eval(&mut ctx, &mut env, &block.code);
}
