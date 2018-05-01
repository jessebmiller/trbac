package auth

type Auth struct {
	Permissionser Permissionser
	ConstraintRunner ConstraintRunner
}

func isIn(str string, collection []string) bool {
	// check them all, return true as soon as it's found
	for _, member := range collection {
		if (str == member) {
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
func (auth Auth) May(ctx Context) bool {
	// find a permission that grants in this context
	permissions := auth.Permissionser.Permissions(ctx.Roles())
	for _, perm := range permissions {
		if !relevantPermission(perm, ctx) {
			continue
		}
		if allConstraintsPass(perm.Constraints, auth.ConstraintRunner, ctx) {
			return true
		}
	}
	// if we don't find a permission that grants in this context
	// May is false
	return false
}
