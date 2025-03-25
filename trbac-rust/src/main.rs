mod auth;
mod context;
mod constraint_runner;
mod privileges;

use std::collections::HashMap;
use auth::BasicAuth;
use context::Context;
use constraint_runner::FuncMapConstraintRunner;
use privileges::Privileges;

fn main() {
    // Initialize privileges
    let mut privileges = Privileges::new();
    privileges.load_from_file("privileges.toml").expect("Failed to load privileges");

    // Initialize constraint runner
    let mut constraint_funcs: HashMap<String, Box<dyn Fn(&Context) -> bool>> = HashMap::new();
    constraint_funcs.insert("business_hours_only".to_string(), Box::new(|_ctx| true)); // Placeholder
    constraint_funcs.insert("audit_logged".to_string(), Box::new(|ctx| {
        println!("AUDIT: {} {} by {:?}", ctx.action(), ctx.resource_type(), ctx.roles());
        true
    }));
    let constraint_runner = FuncMapConstraintRunner::new(constraint_funcs);

    // Initialize authorization service
    let auth_service = BasicAuth::new(privileges, Box::new(constraint_runner));

    // Create a context
    let ctx = Context::new("read", "document", vec!["reader"]);

    // Check authorization
    if auth_service.may(&ctx) {
        println!("Access granted");
    } else {
        println!("Access denied");
    }
}
