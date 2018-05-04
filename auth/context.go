package auth

type Context interface {
	Action() string
	ResourceType() string
	Roles() []string
}

type literalContext struct {
	action       string
	resourceType string
	roles        []string
}

func (lc literalContext) Action() string {
	return lc.action
}

func (lc literalContext) ResourceType() string {
	return lc.resourceType
}

func (lc literalContext) Roles() []string {
	return lc.roles
}

func NewLiteralContext(action string, resourceType string, roles []string) Context {
	return literalContext{action, resourceType, roles}
}
