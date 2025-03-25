# TRBAC Implementation Guide for Rust

This guide provides a practical approach to implementing TRBAC (Typed Role-Based Access Control) in Rust. It follows the formal specification while leveraging Rust's strengths: strong typing, ownership model, and performance characteristics.

## Core Traits

The Rust implementation of TRBAC centers around four key traits:

### Context Trait

The `Context` trait encapsulates the information needed for authorization decisions:

```rust
/// Provides authorization decision information
pub trait Context {
    /// Returns the operation being attempted
    fn action(&self) -> &str;
    
    /// Returns the type of resource being accessed
    fn resource_type(&self) -> &str;
    
    /// Returns all roles assigned to the actor
    fn roles(&self) -> &[String];
}
```

### ConstraintRunner Trait

The `ConstraintRunner` trait defines how constraints are evaluated:

```rust
/// Evaluates named constraints against contexts
pub trait ConstraintRunner {
    /// Evaluates a named constraint against a context
    /// Returns true if the constraint passes, false otherwise
    fn run(&self, constraint_name: &str, ctx: &dyn Context) -> bool;
}
```

### Privileges Trait

The `Privileges` trait defines how permissions are assigned to roles:

```rust
/// Defines the relationship between roles and permissions
pub trait Privileges {
    /// Returns all permissions assigned to any of the given roles
    fn get_permissions(&self, roles: &[String]) -> Vec<Permission>;
}
```

### Auth Trait

The `Auth` trait provides the main authorization behavior:

```rust
/// Provides the core authorization capability
pub trait Auth {
    /// Determines if the given context is authorized
    /// Returns true if authorized, false otherwise
    fn may(&self, ctx: &dyn Context) -> bool;
}
```

## Implementation Structures

### Permission

The `Permission` structure represents allowed actions on resource types with constraints:

```rust
/// Represents allowed actions on resources
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct Permission {
    /// Allowed actions
    pub actions: Vec<String>,
    
    /// Applicable resource types
    pub resource_types: Vec<String>,
    
    /// Required constraints
    pub constraints: Vec<String>,
}
```

### LiteralContext

A basic `Context` implementation for simple use cases:

```rust
/// A simple Context implementation
#[derive(Clone, Debug)]
pub struct LiteralContext {
    action: String,
    resource_type: String,
    roles: Vec<String>,
}

impl Context for LiteralContext {
    fn action(&self) -> &str {
        &self.action
    }
    
    fn resource_type(&self) -> &str {
        &self.resource_type
    }
    
    fn roles(&self) -> &[String] {
        &self.roles
    }
}

impl LiteralContext {
    /// Creates a new LiteralContext
    pub fn new(action: String, resource_type: String, roles: Vec<String>) -> Self {
        Self {
            action,
            resource_type,
            roles,
        }
    }
}
```

### Basic Auth Implementation

A straightforward `Auth` implementation:

```rust
/// Implements the Auth trait
pub struct BasicAuth<P: Privileges, C: ConstraintRunner> {
    privileges: P,
    constraint_runner: C,
}

impl<P: Privileges, C: ConstraintRunner> Auth for BasicAuth<P, C> {
    fn may(&self, ctx: &dyn Context) -> bool {
        // Get all permissions for the roles
        let permissions = self.privileges.get_permissions(ctx.roles());
        
        // Check each permission
        for perm in permissions {
            // Check if permission is relevant
            if !is_relevant(&perm, ctx) {
                continue;
            }
            
            // Check if all constraints pass
            if all_constraints_pass(&perm.constraints, &self.constraint_runner, ctx) {
                return true;
            }
        }
        
        // If no permission grants, deny
        false
    }
}

impl<P: Privileges, C: ConstraintRunner> BasicAuth<P, C> {
    /// Creates a new BasicAuth
    pub fn new(privileges: P, constraint_runner: C) -> Self {
        Self {
            privileges,
            constraint_runner,
        }
    }
}

/// Helper function to check if a permission is relevant to the context
fn is_relevant(perm: &Permission, ctx: &dyn Context) -> bool {
    perm.actions.iter().any(|a| a == ctx.action()) && 
    perm.resource_types.iter().any(|t| t == ctx.resource_type())
}

/// Helper function to check if all constraints pass
fn all_constraints_pass(
    constraints: &[String], 
    runner: &dyn ConstraintRunner, 
    ctx: &dyn Context
) -> bool {
    constraints.iter().all(|c| runner.run(c, ctx))
}
```

## Constraint Runners

### Closure Constraint Runner

A simple in-memory constraint runner using Rust closures:

```rust
/// Uses Rust closures as constraints
pub struct ClosureConstraintRunner {
    constraint_funcs: HashMap<String, Box<dyn Fn(&dyn Context) -> bool>>,
}

impl ConstraintRunner for ClosureConstraintRunner {
    fn run(&self, constraint: &str, ctx: &dyn Context) -> bool {
        match self.constraint_funcs.get(constraint) {
            Some(func) => func(ctx),
            None => false,
        }
    }
}

impl ClosureConstraintRunner {
    /// Creates a new function-based constraint runner
    pub fn new() -> Self {
        Self {
            constraint_funcs: HashMap::new(),
        }
    }
    
    /// Adds a constraint function
    pub fn add_constraint<F>(&mut self, name: &str, func: F) 
    where
        F: Fn(&dyn Context) -> bool + 'static,
    {
        self.constraint_funcs.insert(name.to_string(), Box::new(func));
    }
}
```

### Command Constraint Runner

A constraint runner that executes external commands with minimal overhead:

```rust
use std::path::{Path, PathBuf};
use std::process::Command;

/// Executes commands as constraints
pub struct CommandConstraintRunner {
    script_root: PathBuf,
}

impl ConstraintRunner for CommandConstraintRunner {
    fn run(&self, constraint: &str, ctx: &dyn Context) -> bool {
        let script_path = self.script_root.join(constraint);
        
        // Convert context to command arguments
        let mut args = vec![ctx.action().to_string(), ctx.resource_type().to_string()];
        args.extend(ctx.roles().iter().cloned());
        
        // Execute the command
        match Command::new(&script_path).args(&args).status() {
            Ok(status) => status.success(),
            Err(_) => false,
        }
    }
}

impl CommandConstraintRunner {
    /// Creates a new command constraint runner
    pub fn new<P: AsRef<Path>>(script_root: P) -> Self {
        Self {
            script_root: script_root.as_ref().to_path_buf(),
        }
    }
}
```

## Privilege Providers

### TOML Privileges Provider

A provider that loads privileges from TOML files:

```rust
use std::collections::HashMap;
use std::fs;
use std::path::Path;

use serde::Deserialize;
use toml;

/// Loads privileges from TOML files
pub struct TomlPrivileges {
    permissions: HashMap<String, Vec<Permission>>,
}

impl Privileges for TomlPrivileges {
    fn get_permissions(&self, roles: &[String]) -> Vec<Permission> {
        let mut result = Vec::new();
        for role in roles {
            if let Some(perms) = self.permissions.get(role) {
                result.extend(perms.clone());
            }
        }
        result
    }
}

impl TomlPrivileges {
    /// Creates a new TOML-based privileges provider
    pub fn new<P: AsRef<Path>>(toml_path: P) -> Result<Self, Box<dyn std::error::Error>> {
        let toml_str = fs::read_to_string(toml_path)?;
        let permissions: HashMap<String, Vec<Permission>> = toml::from_str(&toml_str)?;
        Ok(Self { permissions })
    }
}
```

## Example Usage

Here's how to use TRBAC in a Rust application:

```rust
use std::error::Error;
use trbac::{Auth, BasicAuth, ClosureConstraintRunner, Context, LiteralContext, TomlPrivileges};

fn main() -> Result<(), Box<dyn Error>> {
    // Load privileges from TOML
    let privs = TomlPrivileges::new("privileges.toml")?;
    
    // Create a constraint runner
    let mut constraint_runner = ClosureConstraintRunner::new();
    
    // Add constraints
    constraint_runner.add_constraint("business_hours_only", |_ctx| {
        // Check if current time is within business hours
        // This is just an example, replace with real logic
        true
    });
    
    constraint_runner.add_constraint("audit_logged", |ctx| {
        // Log the access attempt for audit purposes
        println!(
            "AUDIT: {} {} by {:?}", 
            ctx.action(), 
            ctx.resource_type(), 
            ctx.roles()
        );
        true
    });
    
    // Create the auth service
    let auth_service = BasicAuth::new(privs, constraint_runner);
    
    // Create a context for authorization check
    let ctx = LiteralContext::new(
        "read".to_string(),
        "document".to_string(),
        vec!["reader".to_string()],
    );
    
    // Check if action is authorized
    if auth_service.may(&ctx) {
        println!("Access granted");
    } else {
        println!("Access denied");
    }
    
    Ok(())
}
```

## Web Framework Integration

Here's an example of how to integrate TRBAC with the Actix web framework:

```rust
use actix_web::{web, App, HttpServer, HttpRequest, HttpResponse, Responder, Error};
use actix_web::dev::ServiceRequest;
use actix_web_httpauth::extractors::bearer::BearerAuth;
use actix_web_httpauth::middleware::HttpAuthentication;
use futures::future::{ready, Ready};
use std::sync::Arc;

use trbac::{Auth, Context, LiteralContext};

// Shared state containing auth service
struct AppState<A: Auth + Send + Sync + 'static> {
    auth_service: Arc<A>,
}

// Extract context from HTTP request
fn extract_context(req: &HttpRequest, token_data: &str) -> impl Context {
    // This is a simplified example
    // In practice, you would extract values from:
    // - URL path
    // - HTTP method
    // - JWT claims
    
    let action = action_from_method(req.method().as_str());
    let resource_type = resource_type_from_path(req.path());
    let roles = roles_from_token(token_data);
    
    LiteralContext::new(action, resource_type, roles)
}

// TRBAC authentication middleware for Actix
async fn validator<A: Auth + Send + Sync + 'static>(
    req: ServiceRequest,
    credentials: BearerAuth,
) -> Result<ServiceRequest, Error> {
    let auth_service = req.app_data::<web::Data<AppState<A>>>()
        .expect("Auth service not configured")
        .auth_service.clone();
    
    // Extract token data (in practice, verify JWT, etc.)
    let token_data = credentials.token();
    
    // Extract context from request
    let ctx = extract_context(req.request(), token_data);
    
    // Check authorization
    if auth_service.may(&ctx) {
        Ok(req)
    } else {
        Err(actix_web::error::ErrorForbidden("Forbidden"))
    }
}

// Example handler
async fn get_documents() -> impl Responder {
    HttpResponse::Ok().body("Accessed documents")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Initialize auth service (see previous example)
    // ...
    
    // Create shared state
    let app_state = web::Data::new(AppState {
        auth_service: Arc::new(auth_service),
    });
    
    // Start server
    HttpServer::new(move || {
        let auth = HttpAuthentication::bearer(validator::<BasicAuth<_, _>>);
        
        App::new()
            .app_data(app_state.clone())
            .wrap(auth)
            .route("/documents", web::get().to(get_documents))
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}

// Utility functions (would be implemented with real logic)
fn action_from_method(method: &str) -> String {
    match method {
        "GET" => "read",
        "POST" => "create",
        "PUT" => "update",
        "DELETE" => "delete",
        _ => "unknown",
    }
    .to_string()
}

fn resource_type_from_path(path: &str) -> String {
    // Simple example - in practice, use a router or regex
    if path == "/documents" {
        "document"
    } else {
        "unknown"
    }
    .to_string()
}

fn roles_from_token(token: &str) -> Vec<String> {
    // In practice, extract from JWT claims
    vec!["reader".to_string()]
}
```

## Performance Optimizations

For embedded systems and performance-critical applications:

### 1. Zero-Copy Parsing

Use zero-copy parsing for configuration files:

```rust
use serde::Deserialize;

#[derive(Deserialize)]
struct PrivilegesConfig<'a> {
    #[serde(borrow)]
    roles: HashMap<&'a str, Vec<PermissionConfig<'a>>>,
}

#[derive(Deserialize)]
struct PermissionConfig<'a> {
    #[serde(borrow)]
    actions: Vec<&'a str>,
    #[serde(borrow)]
    resource_types: Vec<&'a str>,
    #[serde(borrow)]
    constraints: Vec<&'a str>,
}
```

### 2. Static Dispatch

Use static dispatch to avoid dynamic dispatch overhead:

```rust
pub struct StaticAuth<P, C> 
where
    P: Privileges,
    C: ConstraintRunner,
{
    privileges: P,
    constraint_runner: C,
}

impl<P, C> Auth for StaticAuth<P, C>
where
    P: Privileges,
    C: ConstraintRunner,
{
    fn may(&self, ctx: &dyn Context) -> bool {
        // Implementation similar to BasicAuth
    }
}
```

### 3. Arena Allocation

Use arena allocation for temporary objects:

```rust
use bumpalo::Bump;

// Create an arena allocator
let arena = Bump::new();

// Allocate context in the arena
let action = arena.alloc_str("read");
let resource_type = arena.alloc_str("document");
let roles = arena.alloc_slice_copy(&["reader"]);

// Create context using arena-allocated data
let ctx = ArenaContext {
    action,
    resource_type,
    roles,
};
```

### 4. No-Std Support

Support embedded environments without the standard library:

```rust
#![no_std]
extern crate alloc;

use alloc::string::String;
use alloc::vec::Vec;
use alloc::boxed::Box;

// Core traits and implementations similar to above,
// but using only no_std compatible dependencies
```

## Testing Strategies

TRBAC in Rust can be effectively tested using Rust's testing framework:

```rust
#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashMap;
    
    #[test]
    fn test_basic_auth() {
        // Define test permissions
        let mut test_perms = HashMap::new();
        
        test_perms.insert(
            "reader".to_string(),
            vec![Permission {
                actions: vec!["read".to_string()],
                resource_types: vec!["document".to_string()],
                constraints: vec![],
            }],
        );
        
        test_perms.insert(
            "writer".to_string(),
            vec![Permission {
                actions: vec!["read".to_string(), "write".to_string()],
                resource_types: vec!["document".to_string()],
                constraints: vec!["business_hours".to_string()],
            }],
        );
        
        // Create test privileges
        struct TestPrivileges(HashMap<String, Vec<Permission>>);
        
        impl Privileges for TestPrivileges {
            fn get_permissions(&self, roles: &[String]) -> Vec<Permission> {
                let mut result = Vec::new();
                for role in roles {
                    if let Some(perms) = self.0.get(role) {
                        result.extend(perms.clone());
                    }
                }
                result
            }
        }
        
        let privs = TestPrivileges(test_perms);
        
        // Create constraint runners
        let mut always_pass = ClosureConstraintRunner::new();
        always_pass.add_constraint("business_hours", |_| true);
        
        let mut always_fail = ClosureConstraintRunner::new();
        always_fail.add_constraint("business_hours", |_| false);
        
        // Create auth services
        let auth_with_pass = BasicAuth::new(privs.clone(), always_pass);
        let auth_with_fail = BasicAuth::new(privs, always_fail);
        
        // Test cases
        struct TestCase {
            name: &'static str,
            ctx: LiteralContext,
            auth: &dyn Auth,
            expected: bool,
        }
        
        let tests = vec![
            TestCase {
                name: "Reader can read document",
                ctx: LiteralContext::new(
                    "read".to_string(),
                    "document".to_string(),
                    vec!["reader".to_string()],
                ),
                auth: &auth_with_pass,
                expected: true,
            },
            TestCase {
                name: "Reader cannot write document",
                ctx: LiteralContext::new(
                    "write".to_string(),
                    "document".to_string(),
                    vec!["reader".to_string()],
                ),
                auth: &auth_with_pass,
                expected: false,
            },
            TestCase {
                name: "Writer can read document during business hours",
                ctx: LiteralContext::new(
                    "read".to_string(),
                    "document".to_string(),
                    vec!["writer".to_string()],
                ),
                auth: &auth_with_pass,
                expected: true,
            },
            TestCase {
                name: "Writer cannot read document outside business hours",
                ctx: LiteralContext::new(
                    "read".to_string(),
                    "document".to_string(),
                    vec!["writer".to_string()],
                ),
                auth: &auth_with_fail,
                expected: false,
            },
        ];
        
        // Run tests
        for test in tests {
            let result = test.auth.may(&test.ctx);
            assert_eq!(
                result, test.expected,
                "{}: expected {}, got {}", 
                test.name, test.expected, result
            );
        }
    }
}
```

## Best Practices for Rust Implementation

1. **Leverage the type system** for compile-time guarantees
2. **Use traits** for flexibility and testability
3. **Minimize allocations** in performance-critical paths
4. **Consider no_std compatibility** for embedded environments
5. **Use const generics** for zero-overhead abstractions
6. **Implement Send and Sync** for thread safety
7. **Use proper error handling** with Result and Option types
8. **Document all public APIs** with rustdoc comments
