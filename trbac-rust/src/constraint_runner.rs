use std::collections::HashMap;

pub trait ConstraintRunner {
    fn run(&self, constraint_name: &str, ctx: &super::context::Context) -> bool;
}

pub struct FuncMapConstraintRunner {
    constraint_funcs: HashMap<String, Box<dyn Fn(&super::context::Context) -> bool>>,
}

impl FuncMapConstraintRunner {
    pub fn new(funcs: HashMap<String, Box<dyn Fn(&super::context::Context) -> bool>>) -> Self {
        FuncMapConstraintRunner { constraint_funcs: funcs }
    }
}

impl ConstraintRunner for FuncMapConstraintRunner {
    fn run(&self, constraint_name: &str, ctx: &super::context::Context) -> bool {
        if let Some(func) = self.constraint_funcs.get(constraint_name) {
            func(ctx)
        } else {
            false
        }
    }
}
