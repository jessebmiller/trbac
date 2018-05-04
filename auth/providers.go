package auth

/*

Providers are inteded to be metadat driven.

Currently there are simply builtin internal implementations of the interfaces
required.

We intend to implement providers that are configurable with flat files, like
a yaml Privileges (mapping from roles to permissions, who can do what) and a
shell script or docker constraint runner.

*/

// Permission is the right to take an action
// Any of these actions, on any of these resource types are allowed unless any
// of these constraints fail (return false)
type Permission struct {
	Actions       []string // These actions
	ResourceTypes []string // are allowed on these resource types
	Constraints   []string // as long as these contstraints hold
}

// Privileges map roles to their permissions (who can do what)
// This should be the union of permissions granted to each role
type Privileges interface {
	// TODO this might error, handle that
	GetPermissions([]string) []Permission // Roles to Permissions
}

type ConstraintRunner interface {
	// TODO this might error, handle that..
	Run(string, Context) bool
}

/*
 * builtin implementations
 */

// constraintFunc is a Go function implementing a constraint
type constraintFunc = func(Context) bool

// funcMapConstraintRunner uses a literal map from constraint name to a Go function constraint
type funcMapConstraintRunner struct {
	constraintFuncs map[string]constraintFunc
}

func (runner funcMapConstraintRunner) Run(constraint string, ctx Context) bool {
	return runner.constraintFuncs[constraint](ctx)
}

func NewFuncMapConstraintRunner(fm map[string]constraintFunc) funcMapConstraintRunner {
	return funcMapConstraintRunner{fm}
}

// mapPrivileges uses an in memory Go map to look up Permissions from roles
type mapPrivileges struct {
	// role to permissions granted to that role
	permissions map[string][]Permission
}

func (privs mapPrivileges) GetPermissions(roles []string) []Permission {
	// gather the permissions granted to each given role
	perms := []Permission{}
	for _, r := range roles {
		perms = append(perms, privs.permissions[r]...)
	}
	return perms
}

func NewMapPrivileges(permsMap map[string][]Permission) mapPrivileges {
	return mapPrivileges{permsMap}
}
