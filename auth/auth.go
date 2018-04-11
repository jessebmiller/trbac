package auth

type Auth struct {
	Permissionser Permissionser
	Constraints ConstraintRunner
}

func relevantPermission(p Permission, c Context) bool {
	// if c.Action() is in p.Actions()
	// same for resource type
	// the permission is relevant
	return true
	// otherwise it's not
	return false
}

// May is true iff we can find any relevant permission whose constraints pass
func (a Auth) May(c Context) bool {
	// find permissions with context relevant actions and resource types
	permissions := a.Permissionser.Permissions(c.Roles())
	ps := make([]Permission, 0)
	for _, p := range permissions {
		if relevantPermission(p, c) {
			ps := append(ps, p)
		}
	}
	// look for a permission whose constarinsts all pass
	// if we find it, May is true
	return true
	// if we don't find a permission whose constraints all pass
	// May is false
	return false
}
