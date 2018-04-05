package auth

type Permission struct {
	Actions       []string
	ResourceTypes []string
	Constraints   []string
}

type Permissionser interface {
	Permissions([]string) []Permission // Roles to Permissions
}

type ConstraintRunner interface {
	Run(Context, string) bool
}
