use super::context::Context;
use super::constraint_runner::ConstraintRunner;
use super::privileges::{Permission, Privileges};

pub struct BasicAuth {
    privileges: Privileges,
    constraint_runner: Box<dyn ConstraintRunner>,
}

impl BasicAuth {
    pub fn new(privileges: Privileges, constraint_runner: Box<dyn ConstraintRunner>) -> Self {
        BasicAuth {
            privileges,
            constraint_runner,
        }
    }

    pub fn may(&self, ctx: &Context) -> bool {
        let permissions = self.privileges.get_permissions(ctx.roles());

        for perm in permissions {
            if self.is_relevant(perm, ctx) && self.all_constraints_pass(perm.constraints(), ctx) {
                return true;
            }
        }

        false
    }

    fn is_relevant(&self, perm: &Permission, ctx: &Context) -> bool {
        perm.actions().contains(&ctx.action().to_string()) &&
        perm.resource_types().contains(&ctx.resource_type().to_string())
    }

    fn all_constraints_pass(&self, constraints: &[String], ctx: &Context) -> bool {
        for constraint in constraints {
            if !self.constraint_runner.run(constraint, ctx) {
                return false;
            }
        }
        true
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::constraint_runner::FuncMapConstraintRunner;
    use std::collections::HashMap;

    #[test]
    fn test_may_with_authorized_context() {
        let mut privileges = Privileges::new();
        privileges.load_from_file("dummy_path").unwrap();

        let mut funcs: HashMap<String, Box<dyn Fn(&Context) -> bool>> = HashMap::new();
        funcs.insert("business_hours_only".to_string(), Box::new(|_ctx| true));
        let constraint_runner = FuncMapConstraintRunner::new(funcs);

        let auth = BasicAuth::new(privileges, Box::new(constraint_runner));
        let context = Context::new("read", "document", vec!["reader"]);

        assert!(auth.may(&context));
    }

    #[test]
    fn test_may_with_unauthorized_context() {
        let mut privileges = Privileges::new();
        privileges.load_from_file("dummy_path").unwrap();

        let funcs: HashMap<String, Box<dyn Fn(&Context) -> bool>> = HashMap::new();
        let constraint_runner = FuncMapConstraintRunner::new(funcs);

        let auth = BasicAuth::new(privileges, Box::new(constraint_runner));
        let context = Context::new("write", "document", vec!["reader"]);

        assert!(!auth.may(&context));
    }
}
