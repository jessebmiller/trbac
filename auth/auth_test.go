package auth_test

import (
	"testing"
	"github.com/jessebmiller/trbac/auth"
	"github.com/stretchr/testify/assert"
)

type mockContext struct {
	action   string
	resourceType string
	roles    []string
}

func (c mockContext) Action() string {
	return c.action
}

func (c mockContext) ResourceType() string {
	return c.resourceType
}

func (c mockContext) Roles() []string {
	return c.roles
}

type mockPermissionser struct {
	role        string
	permissions []auth.Permission
}

func (p mockPermissionser) Permissions(roles []string) []auth.Permission {
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

func TestAuthMay(t *testing.T) {
	testPermissions := []auth.Permission{
		auth.Permission{
			[]string{"A0", "A1"},     // Actions
			[]string{"R0", "R1"},     // ResourceTypes
			[]string{"C0", "C1"} },   // Constraints
		auth.Permission{
			[]string{"AX"},
			[]string{"RX", "RY"},
			[]string{} },
	}
	testPermissionser := mockPermissionser{ "tester", testPermissions }
	unconstrainedRunner := constantConstraintRunner{ true }
	constrainedRunner := constantConstraintRunner{ false }
	unconstrainedAuth := auth.Auth{ testPermissionser, unconstrainedRunner }
	constrainedAuth := auth.Auth{ testPermissionser, constrainedRunner }
	authorizedContext := mockContext{ "A0", "R1", []string{"tester"} }
	unauthorizedContext := mockContext{ "AX", "R0", []string{"tester"} }
	wrongRoleContext := mockContext{ "A0", "R1", []string{"not-tester"} }
	table := []struct{
		a auth.Auth
		c mockContext
		may bool
	}{
		{ unconstrainedAuth, authorizedContext, true },
		{ constrainedAuth, authorizedContext, false },
		{ unconstrainedAuth, unauthorizedContext, false },
		{ unconstrainedAuth, wrongRoleContext, false },
	}

	for _, row := range table {
		result := row.a.May(row.c)
		assert.Equal(t, result, row.may)
	}
}
