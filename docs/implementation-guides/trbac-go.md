# TRBAC Implementation Guide for Go

This guide provides a practical approach to implementing TRBAC (Typed Role-Based Access Control) in Go. It follows the formal specification while leveraging Go's strengths: strong typing, interfaces, and straightforward concurrency.

## Core Interfaces

The Go implementation of TRBAC centers around four key interfaces:

### Context Interface

The Context interface encapsulates the information needed for authorization decisions:

```go
// Context provides authorization decision information
type Context interface {
    // Action returns the operation being attempted
    Action() string
    
    // ResourceType returns the type of resource being accessed
    ResourceType() string
    
    // Roles returns all roles assigned to the actor
    Roles() []string
}
```

### ConstraintRunner Interface

The ConstraintRunner interface defines how constraints are evaluated:

```go
// ConstraintRunner evaluates named constraints against contexts
type ConstraintRunner interface {
    // Run evaluates a named constraint against a context
    // Returns true if the constraint passes, false otherwise
    Run(constraintName string, ctx Context) bool
}
```

### Privileges Interface

The Privileges interface defines how permissions are assigned to roles:

```go
// Privileges defines the relationship between roles and permissions
type Privileges interface {
    // GetPermissions returns all permissions assigned to any of the given roles
    GetPermissions(roles []string) []Permission
}
```

### Auth Interface

The Auth interface provides the main authorization behavior:

```go
// Auth provides the core authorization capability
type Auth interface {
    // May determines if the given context is authorized
    // Returns true if authorized, false otherwise
    May(ctx Context) bool
}
```

## Implementation Structures

### Permission

The Permission structure represents allowed actions on resource types with constraints:

```go
// Permission represents allowed actions on resources
type Permission struct {
    Actions       []string // Allowed actions
    ResourceTypes []string // Applicable resource types
    Constraints   []string // Required constraints
}
```

### LiteralContext

A basic Context implementation for simple use cases:

```go
// LiteralContext is a simple Context implementation
type LiteralContext struct {
    action       string
    resourceType string
    roles        []string
}

func (lc LiteralContext) Action() string {
    return lc.action
}

func (lc LiteralContext) ResourceType() string {
    return lc.resourceType
}

func (lc LiteralContext) Roles() []string {
    return lc.roles
}

// NewLiteralContext creates a new LiteralContext
func NewLiteralContext(action, resourceType string, roles []string) Context {
    return LiteralContext{action, resourceType, roles}
}
```

### Basic Auth Implementation

A straightforward Auth implementation:

```go
// BasicAuth implements the Auth interface
type BasicAuth struct {
    Privileges       Privileges
    ConstraintRunner ConstraintRunner
}

// May determines if the given context is authorized
func (auth BasicAuth) May(ctx Context) bool {
    // Get all permissions for the roles
    permissions := auth.Privileges.GetPermissions(ctx.Roles())
    
    // Check each permission
    for _, perm := range permissions {
        // Check if permission is relevant
        if !isRelevant(perm, ctx) {
            continue
        }
        
        // Check if all constraints pass
        if allConstraintsPass(perm.Constraints, auth.ConstraintRunner, ctx) {
            return true
        }
    }
    
    // If no permission grants, deny
    return false
}

// Helper function to check if a permission is relevant to the context
func isRelevant(perm Permission, ctx Context) bool {
    return contains(perm.Actions, ctx.Action()) && 
           contains(perm.ResourceTypes, ctx.ResourceType())
}

// Helper function to check if all constraints pass
func allConstraintsPass(constraints []string, runner ConstraintRunner, ctx Context) bool {
    for _, constraint := range constraints {
        if !runner.Run(constraint, ctx) {
            return false
        }
    }
    return true
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

## Constraint Runners

### FuncMap Constraint Runner

A simple in-memory constraint runner using Go functions:

```go
// FuncMapConstraintRunner uses Go functions as constraints
type FuncMapConstraintRunner struct {
    constraintFuncs map[string]func(Context) bool
}

// Run evaluates a constraint against a context
func (runner FuncMapConstraintRunner) Run(constraint string, ctx Context) bool {
    fn, exists := runner.constraintFuncs[constraint]
    if !exists {
        return false
    }
    return fn(ctx)
}

// NewFuncMapConstraintRunner creates a new function-based constraint runner
func NewFuncMapConstraintRunner(funcs map[string]func(Context) bool) ConstraintRunner {
    return FuncMapConstraintRunner{funcs}
}
```

### Shell Script Constraint Runner

A constraint runner that uses external shell scripts:

```go
// ShellScriptConstraintRunner executes shell scripts as constraints
type ShellScriptConstraintRunner struct {
    scriptRoot string
}

// Run evaluates a constraint against a context
func (runner ShellScriptConstraintRunner) Run(constraint string, ctx Context) bool {
    scriptPath := path.Join(runner.scriptRoot, constraint)
    
    // Convert context to command arguments
    args := []string{ctx.Action(), ctx.ResourceType()}
    for _, role := range ctx.Roles() {
        args = append(args, role)
    }
    
    // Execute the script
    cmd := exec.Command(scriptPath, args...)
    err := cmd.Run()
    
    // Script succeeds if exit code is 0
    return err == nil
}

// NewShellScriptConstraintRunner creates a new shell script constraint runner
func NewShellScriptConstraintRunner(scriptRoot string) ConstraintRunner {
    return ShellScriptConstraintRunner{scriptRoot}
}
```

## Privilege Providers

### TOML Privileges Provider

A provider that loads privileges from TOML files:

```go
// TOMLPrivileges loads privileges from TOML files
type TOMLPrivileges struct {
    permissions map[string][]Permission
}

// GetPermissions returns all permissions for the given roles
func (privs TOMLPrivileges) GetPermissions(roles []string) []Permission {
    result := []Permission{}
    for _, role := range roles {
        result = append(result, privs.permissions[role]...)
    }
    return result
}

// NewTOMLPrivileges creates a new TOML-based privileges provider
func NewTOMLPrivileges(tomlPath string) (Privileges, error) {
    var privs map[string][]Permission
    if _, err := toml.DecodeFile(tomlPath, &privs); err != nil {
        return nil, err
    }
    return TOMLPrivileges{privs}, nil
}
```

## Example Usage

Here's how to use TRBAC in a Go application:

```go
package main

import (
    "fmt"
    "github.com/yourorg/trbac/auth"
)

func main() {
    // Load privileges from TOML
    privs, err := auth.NewTOMLPrivileges("privileges.toml")
    if err != nil {
        panic(err)
    }
    
    // Create a constraint runner
    constraintRunner := auth.NewFuncMapConstraintRunner(map[string]func(auth.Context) bool{
        "business_hours_only": func(ctx auth.Context) bool {
            // Check if current time is within business hours
            // This is just an example, replace with real logic
            return true
        },
        "audit_logged": func(ctx auth.Context) bool {
            // Log the access attempt for audit purposes
            fmt.Printf("AUDIT: %s %s by %v\n", 
                ctx.Action(), ctx.ResourceType(), ctx.Roles())
            return true
        },
    })
    
    // Create the auth service
    authService := auth.BasicAuth{
        Privileges:       privs,
        ConstraintRunner: constraintRunner,
    }
    
    // Create a context for authorization check
    ctx := auth.NewLiteralContext("read", "document", []string{"reader"})
    
    // Check if action is authorized
    if authService.May(ctx) {
        fmt.Println("Access granted")
    } else {
        fmt.Println("Access denied")
    }
}
```

## Web Framework Integration

Here's an example of how to integrate TRBAC with a standard Go HTTP server:

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/yourorg/trbac/auth"
)

// AuthMiddleware creates middleware to enforce TRBAC authorization
func AuthMiddleware(authService auth.Auth, extractContext func(*http.Request) auth.Context) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract context from request
            ctx := extractContext(r)
            
            // Check authorization
            if !authService.May(ctx) {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            
            // If authorized, proceed to next handler
            next.ServeHTTP(w, r)
        })
    }
}

func main() {
    // Initialize auth service (see previous example)
    // ...
    
    // Define how to extract context from HTTP requests
    extractContext := func(r *http.Request) auth.Context {
        // This is a simplified example
        // In a real application, you would extract values from:
        // - URL path
        // - HTTP method
        // - User session/JWT claims
        action := actionFromMethod(r.Method)
        resourceType := resourceTypeFromPath(r.URL.Path)
        roles := rolesFromRequest(r)
        
        return auth.NewLiteralContext(action, resourceType, roles)
    }
    
    // Create our middleware
    authorize := AuthMiddleware(authService, extractContext)
    
    // Define routes
    mux := http.NewServeMux()
    
    // Protected route
    mux.Handle("/documents/", authorize(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Accessed documents")
        })))
    
    // Start server
    http.ListenAndServe(":8080", mux)
}

// Utility functions (would be implemented with real logic)
func actionFromMethod(method string) string {
    switch method {
    case "GET":
        return "read"
    case "POST":
        return "create"
    case "PUT":
        return "update"
    case "DELETE":
        return "delete"
    default:
        return "unknown"
    }
}

func resourceTypeFromPath(path string) string {
    // Simple example - in practice, use a router or regexp
    if path == "/documents/" {
        return "document"
    }
    return "unknown"
}

func rolesFromRequest(r *http.Request) []string {
    // In practice, extract from JWT, session, etc.
    return []string{"reader"}
}
```

## Performance Considerations

For high-performance Go applications:

1. **Cache Permission Lookups**: Cache the results of `GetPermissions` to avoid repeated lookups.

2. **Optimize Constraint Runners**: Use in-memory function constraints for frequently evaluated constraints.

3. **Context Pooling**: Use a sync.Pool to reuse Context objects, reducing allocations.

4. **Parallel Constraint Evaluation**: For complex constraints, consider parallel evaluation using goroutines.

```go
// Parallel constraint evaluation example
func allConstraintsPassParallel(constraints []string, runner ConstraintRunner, ctx Context) bool {
    if len(constraints) == 0 {
        return true
    }
    
    results := make(chan bool, len(constraints))
    
    // Evaluate each constraint in a goroutine
    for _, constraint := range constraints {
        go func(c string) {
            results <- runner.Run(c, ctx)
        }(constraint)
    }
    
    // Check results
    for i := 0; i < len(constraints); i++ {
        if !<-results {
            return false
        }
    }
    
    return true
}
```

## Testing Strategies

TRBAC in Go can be effectively tested using Go's testing framework:

```go
func TestBasicAuth(t *testing.T) {
    // Define test permissions
    testPerms := map[string][]auth.Permission{
        "reader": {
            {
                Actions:       []string{"read"},
                ResourceTypes: []string{"document"},
                Constraints:   []string{},
            },
        },
        "writer": {
            {
                Actions:       []string{"read", "write"},
                ResourceTypes: []string{"document"},
                Constraints:   []string{"business_hours"},
            },
        },
    }
    
    // Create test privileges
    privs := auth.TOMLPrivileges{testPerms}
    
    // Create constraint runners
    alwaysPass := auth.NewFuncMapConstraintRunner(map[string]func(auth.Context) bool{
        "business_hours": func(ctx auth.Context) bool { return true },
    })
    
    alwaysFail := auth.NewFuncMapConstraintRunner(map[string]func(auth.Context) bool{
        "business_hours": func(ctx auth.Context) bool { return false },
    })
    
    // Create auth services
    authWithPass := auth.BasicAuth{privs, alwaysPass}
    authWithFail := auth.BasicAuth{privs, alwaysFail}
    
    // Test cases
    tests := []struct {
        name     string
        auth     auth.Auth
        ctx      auth.Context
        expected bool
    }{
        {
            "Reader can read document",
            authWithPass,
            auth.NewLiteralContext("read", "document", []string{"reader"}),
            true,
        },
        {
            "Reader cannot write document",
            authWithPass,
            auth.NewLiteralContext("write", "document", []string{"reader"}),
            false,
        },
        {
            "Writer can read document during business hours",
            authWithPass,
            auth.NewLiteralContext("read", "document", []string{"writer"}),
            true,
        },
        {
            "Writer cannot read document outside business hours",
            authWithFail,
            auth.NewLiteralContext("read", "document", []string{"writer"}),
            false,
        },
    }
    
    // Run tests
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := test.auth.May(test.ctx)
            if result != test.expected {
                t.Errorf("Expected %v but got %v", test.expected, result)
            }
        })
    }
}
```

## Best Practices for Go Implementation

1. **Use interfaces** for flexibility and testability
2. **Keep the core logic simple** with clear separation of concerns
3. **Leverage Go's standard library** for features like JSON/TOML parsing
4. **Use context.Context** for propagating request context in web applications
5. **Consider sync.RWMutex** for thread-safe privilege updates
6. **Document all interfaces and functions** with consistent comments
7. **Follow Go's error handling practices** by returning errors explicitly
