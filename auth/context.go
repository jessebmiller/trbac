package auth

type Context interface {
	Action() string
	ResourceType() string
	Roles() []string
}
