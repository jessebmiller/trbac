package trbac

// Auth provides the core authorization capability
type Auth interface {
    // May determines if the given context is authorized
    // Returns true if authorized, false otherwise
    May(ctx Context) bool
}

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
