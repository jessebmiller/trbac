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
