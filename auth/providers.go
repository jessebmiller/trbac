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
	// TODO this might error, handle that..
	Run(string, Context) bool
}

/*
 * builtin implementations
 */

type funcMapConstraintRunner struct {
	constraintFuncs map[string]func(Context) bool
}

func (fmcr funcMapConstraintRunner) Run (constraint string, ctx Context) bool {
	return fmcr.constraintFuncs[constraint](ctx)
}

func NewFuncMapConstraintRunner(fm map[string]func(Context) bool) funcMapConstraintRunner {
	return funcMapConstraintRunner{fm}
}

type mapPrivileges struct {
	// role: permissions granted to that role
	permissions map[string][]Permission

}

func (lp mapPrivileges) Permissions(roles []string) []Permission {
	// gather the permissions granted to each given role
	perms := []Permission{}
	for _, r := range roles {
		perms = append(perms, lp.permissions[r]...)
	}
	return perms
}

func NewMapPrivileges(permsMap map[string][]Permission) mapPrivileges {
	return mapPrivileges{permsMap}
}
