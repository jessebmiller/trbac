package trbac

// Context provides authorization decision information
type Context interface {
    // Action returns the operation being attempted
    Action() string
    
    // ResourceType returns the type of resource being accessed
    ResourceType() string
    
    // Roles returns all roles assigned to the actor
    Roles() []string
}

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
