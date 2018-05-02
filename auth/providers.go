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

func (runner funcMapConstraintRunner) Run (constraint string, ctx Context) bool {
	return runner.constraintFuncs[constraint](ctx)
}

func NewFuncMapConstraintRunner(funcMap map[string]func(Context) bool) funcMapConstraintRunner {
	return funcMapConstraintRunner{funcMap}
}

type mapPrivileges struct {
	// role: permissions granted to that role
	permissions map[string][]Permission

}

func (privs mapPrivileges) Permissions(roles []string) []Permission {
	// gather the permissions granted to each given role
	perms := []Permission{}
	for _, role := range roles {
		perms = append(perms, privs.permissions[role]...)
	}
	return perms
}

func NewMapPrivileges(permsMap map[string][]Permission) mapPrivileges {
	return mapPrivileges{permsMap}
}
