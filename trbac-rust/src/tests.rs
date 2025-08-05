#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashMap;

    #[test]
    fn test_basic_auth() {
        let mut privileges = Privileges::new();
        privileges.add_permission("reader", "read", "document", vec![]);
        privileges.add_permission("writer", "write", "document", vec!["business_hours"]);

        let mut constraint_funcs: HashMap<String, Box<dyn Fn(&Context) -> bool>> = HashMap::new();
        constraint_funcs.insert("business_hours".to_string(), Box::new(|_ctx| true));

        let constraint_runner = FuncMapConstraintRunner::new(constraint_funcs);
        let auth_service = BasicAuth::new(privileges, Box::new(constraint_runner));

        let ctx_reader = Context::new("read", "document", vec!["reader"]);
        assert!(auth_service.may(&ctx_reader));

        let ctx_writer = Context::new("write", "document", vec!["writer"]);
        assert!(auth_service.may(&ctx_writer));

        let ctx_writer_fail = Context::new("write", "document", vec!["writer"]);
        let constraint_runner_fail = FuncMapConstraintRunner::new(HashMap::new());
        let auth_service_fail = BasicAuth::new(privileges, Box::new(constraint_runner_fail));
        assert!(!auth_service_fail.may(&ctx_writer_fail));
    }

    #[test]
    fn test_context_creation() {
        let ctx = Context::new("read", "document", vec!["reader", "editor"]);
        assert_eq!(ctx.action(), "read");
        assert_eq!(ctx.resource_type(), "document");
        assert_eq!(ctx.roles(), &vec!["reader", "editor"]);
    }

    #[test]
    fn test_constraint_runner() {
        let mut constraints = HashMap::new();
        constraints.insert("always_pass".to_string(), Box::new(|_ctx| true));
        constraints.insert("always_fail".to_string(), Box::new(|_ctx| false));

        let runner = FuncMapConstraintRunner::new(constraints);
        let ctx = Context::new("read", "document", vec!["reader"]);

        assert!(runner.run("always_pass", &ctx));
        assert!(!runner.run("always_fail", &ctx));
        assert!(!runner.run("non_existent", &ctx));
    }

    #[test]
    fn test_privileges_management() {
        let mut privileges = Privileges::new();
        privileges.add_permission("reader", "read", "document", vec![]);

        let perms = privileges.get_permissions(&vec!["reader"]);
        assert_eq!(perms.len(), 1);
        assert_eq!(perms[0].actions, vec!["read"]);
    }
}
