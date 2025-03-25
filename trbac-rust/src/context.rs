pub struct Context {
    action: String,
    resource_type: String,
    roles: Vec<String>,
}

impl Context {
    pub fn new(action: &str, resource_type: &str, roles: Vec<&str>) -> Self {
        Context {
            action: action.to_string(),
            resource_type: resource_type.to_string(),
            roles: roles.into_iter().map(|s| s.to_string()).collect(),
        }
    }

    pub fn action(&self) -> &str {
        &self.action
    }

    pub fn resource_type(&self) -> &str {
        &self.resource_type
    }

    pub fn roles(&self) -> &[String] {
        &self.roles
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_context_creation() {
        let context = Context::new("read", "document", vec!["admin", "user"]);
        assert_eq!(context.action(), "read");
        assert_eq!(context.resource_type(), "document");
        assert_eq!(context.roles(), &["admin", "user"]);
    }

    #[test]
    fn test_context_action() {
        let context = Context::new("write", "file", vec!["editor"]);
        assert_eq!(context.action(), "write");
    }

    #[test]
    fn test_context_resource_type() {
        let context = Context::new("delete", "record", vec!["manager"]);
        assert_eq!(context.resource_type(), "record");
    }

    #[test]
    fn test_context_roles() {
        let context = Context::new("update", "profile", vec!["user", "guest"]);
        assert_eq!(context.roles(), &["user", "guest"]);
    }
}
