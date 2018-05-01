package auth

type Auth struct {
	Permissionser Permissionser
	ConstraintRunner ConstraintRunner
}

func isIn(a string, xs []string) bool {
	// check them all, return true as soon as it's found
	for _, x := range xs {
		if (a == x) {
			return true
		}
	}
	// if not found it's not in there
	return false
}

func relevantPermission(p Permission, c Context) bool {
	// if the context's action is in the permission
	// and the context's resource type is in the permission
	// the permission is relevant
	if (isIn(c.Action(), p.Actions) && isIn(c.ResourceType(), p.ResourceTypes)) {
		return true
	}
	// otherwise it's not
	return false
}

func allConstraintsPass(constraints []string, cr ConstraintRunner, ctx Context) bool {
	// if any constraint fails, return false early
	// if we get through the whole list without failing, return true
	for _, constraint := range constraints {
		if (!cr.Run(constraint, ctx)) {
			return false
		}
	}
	return true
}

// May is true iff we can find any relevant permission whose constraints pass
func (a Auth) May(c Context) bool {
	// find a permission that grants in this context
	permissions := a.Permissionser.Permissions(c.Roles())
	for _, p := range permissions {
		if !relevantPermission(p, c) {
			continue
		}
		if allConstraintsPass(p.Constraints, a.ConstraintRunner, c) {
			return true
		}
	}
	// if we don't find a permission that grants in this context
	// May is false
	return false
}
