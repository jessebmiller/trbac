package main

import (
	"testing"
	"github.com/jessebmiller/trbac/auth"
)

var (
	testPerms []auth.Permission = []auth.Permission{
		auth.Permission{
			[]string{"ActionA", "ActionB"},
			[]string{"ResourceA", "ResourceB"},
			[]string{}, // constraints ignored by constant constraint runner
		},
	}
	constrainedAuth auth.Auth = auth.Auth{
		mockPermissions{ "RoleA", testPerms },
		constantConstraintRunner{ false },
	}

	unconstrainedAuth auth.Auth = auth.Auth{
		mockPermissions{ "RoleA", testPerms },
		constantConstraintRunner{ true },
	}

	relevantRoles []string            = []string{ "RoleA", "RoleX" }
	irrelevantRoles []string          = []string{ "RoleX", "RoleY" }
	emptyRoles []string = []string{}
	relevantCtx auth.Context = mockCtx{ "ActionB", "ResourceA", relevantRoles }
	irrelevantRolesCtx auth.Context       = mockCtx{ "ActionB", "ResourceA", irrelevantRoles }
	irrelevantActionsCtx auth.Context     = mockCtx{ "ActionX", "ResourceA", relevantRoles   }
	irrelevantResourcesCtx auth.Context   = mockCtx{ "ActionB", "ResourceX", relevantRoles   }
	emptyRolesCtx auth.Context             = mockCtx{ "ActionB", "ResourceA", emptyRoles      }
)

func TestMay(t testing.T) {
	var mayCases = []struct{
		auth   auth.Auth
		context auth.Context
		allowed bool
	}{
		// unconstrained auth with a relevant context should be allowed
		{ unconstrainedAuth, relevantCtx, true },

		// but not if constrained
		{ constrainedAuth, relevantCtx, false },

		// any auth with irrelevant context should not be allowed
		{ unconstrainedAuth, irrelevantRolesCtx, false },
		{ unconstrainedAuth, irrelevantActionsCtx, false },
		{ unconstrainedAuth, irrelevantResourcesCtx, false },
		{ constrainedAuth, irrelevantRolesCtx, false },
		{ constrainedAuth, irrelevantActionsCtx, false },
		{ constrainedAuth, irrelevantResourcesCtx, false },

		// nor should any auth with empty roles, also this should error
		{ unconstrainedAuth, emptyRolesCtx, false },
		{ constrainedAuth, emptyRolesCtx, false },

	}

	for _, mc := range mayCases {
		allowed := mc.auth.May(mc.context)
		if allowed != mc.allowed {
			t.Errorf("may(%q) => %q, _. want %q, _",
				mc.context,
				allowed,
				mc.allowed,
			)
		}
	}
}

func TestRequestContextProvider(t testing.T) {

}

func TestRequestRoleProvider(t testing.T) {

}

func TestRequestActionProvider(t testing.T) {

}

func TestConstantResourceTypeProvider(t testing.T) {

}

func TestRequestPermissionsProvider(t testing.T) {

}

func TestShellScriptConstraintRunner(t testing.T) {

}

func TestTOMLMetadataProvider(t testing.T) {

}
