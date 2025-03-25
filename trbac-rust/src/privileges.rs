use std::collections::HashMap;

pub struct Permission {
    actions: Vec<String>,
    resource_types: Vec<String>,
    constraints: Vec<String>,
}

impl Permission {
    pub fn actions(&self) -> &Vec<String> {
        &self.actions
    }

    pub fn resource_types(&self) -> &Vec<String> {
        &self.resource_types
    }

    pub fn constraints(&self) -> &Vec<String> {
        &self.constraints
    }
}

pub struct Privileges {
    permissions: HashMap<String, Vec<Permission>>,
}

impl Privileges {
    pub fn new() -> Self {
        Privileges {
            permissions: HashMap::new(),
        }
    }

    pub fn load_from_file(&mut self, file_path: &str) -> Result<(), Box<dyn std::error::Error>> {
        // Load privileges from a TOML file
        // This is a placeholder for actual file loading logic
        // For now, we'll simulate loading with dummy data
        self.permissions.insert(
            "reader".to_string(),
            vec![Permission {
                actions: vec!["read".to_string()],
                resource_types: vec!["document".to_string()],
                constraints: vec!["business_hours_only".to_string()],
            }],
        );
        Ok(())
    }

    pub fn get_permissions(&self, roles: &[String]) -> Vec<&Permission> {
        let mut result = Vec::new();
        for role in roles {
            if let Some(perms) = self.permissions.get(role) {
                result.extend(perms);
            }
        }
        result
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_load_from_file() {
        let mut privileges = Privileges::new();
        assert!(privileges.load_from_file("dummy_path").is_ok());
        assert!(privileges.permissions.contains_key("reader"));
    }

    #[test]
    fn test_get_permissions_for_existing_role() {
        let mut privileges = Privileges::new();
        privileges.load_from_file("dummy_path").unwrap();

        let roles = vec!["reader".to_string()];
        let permissions = privileges.get_permissions(&roles);
        assert_eq!(permissions.len(), 1);
        assert_eq!(permissions[0].actions(), &vec!["read".to_string()]);
    }

    #[test]
    fn test_get_permissions_for_non_existing_role() {
        let privileges = Privileges::new();
        let roles = vec!["non_existing_role".to_string()];
        let permissions = privileges.get_permissions(&roles);
        assert!(permissions.is_empty());
    }
}
