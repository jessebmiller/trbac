package auth

type Context interface {
	Action() string
	Resource() string
	Roles() []string
}
