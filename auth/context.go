package auth

import "fmt"

type Context interface {
	fmt.Stringer
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

func (lc literalContext) String() string {
	return fmt.Sprintf("%v %v %v", lc.Action(), lc.ResourceType(), lc.Roles())
}

func NewLiteralContext(action string, resourceType string, roles []string) Context {
	return literalContext{action, resourceType, roles}
}
