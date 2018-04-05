package main

import "github.com/jessebmiller/trbac/auth"

type mockCtx struct {
	action string
	resource string
	roles []string
}

func (c mockCtx) Action() string {
	return c.action
}

func (c mockCtx) Resource() string {
	return c.resource
}

func (c mockCtx) Roles() []string {
	return c.roles
}

type mockPermissions struct {
	role string
	permissions []auth.Permission
}

func (p mockPermissions) Permissions(roles []string) []auth.Permission {
	for _, r := range roles {
		if r == p.role {
			return p.permissions
		}
	}
	return []auth.Permission{}
}

type constantConstraintRunner struct {
	result bool
}

func (c constantConstraintRunner) Run(_ auth.Context, _ string) bool {
	return c.result
}
