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

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::Context;

    #[test]
    fn test_constraint_runner_with_existing_constraint() {
        let mut funcs: HashMap<String, Box<dyn Fn(&Context) -> bool>> = HashMap::new();
        funcs.insert("test_constraint".to_string(), Box::new(|_ctx| true));
        let runner = FuncMapConstraintRunner::new(funcs);

        let context = Context::new("action", "resource", vec!["role"]);
        assert!(runner.run("test_constraint", &context));
    }

    #[test]
    fn test_constraint_runner_with_non_existing_constraint() {
        let funcs: HashMap<String, Box<dyn Fn(&Context) -> bool>> = HashMap::new();
        let runner = FuncMapConstraintRunner::new(funcs);

        let context = Context::new("action", "resource", vec!["role"]);
        assert!(!runner.run("non_existing_constraint", &context));
    }
}
