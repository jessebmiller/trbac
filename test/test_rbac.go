package main

import "testing"

func TestTOMLMetadataProvider(t testing.T) {

}

const (
	completeTRBAC TRBAC        = TRBAC{ ["RoleA", "RoleB"], ["ActA", "ActB"], ["resA", "resB"] }
	emptyRolesTRBAC TRBAC      = TRBAC{ [],                 ["ActA", "ActB"], ["resA", "resB"] }
	emptyActionsTRBAC TRBAC    = TRBAC{ ["RoleA", "RoleB"], [],               ["resA", "resB"] }
	emptyResourcesTRBAC TRBAC  = TRBAC{ ["RoleA", "RoleB"], ["ActA", "ActB"], []               }

	relevantUnconstrained Context = TestCtx{ ["RoleA", "XXXXX"], "ActB", "resA", unconstrainedPrivileges }
	relevantConstrained Context   = TestCtx{ ["RoleA", "XXXXX"], "ActB", "resA", constrainedPrivileges   }
	irrelevantRoles Context       = TestCtx{ ["XXXXX", "YYYYY"], "ActB", "resA", unconstrainedPrivileges }
	irrelevantActions Context     = TestCtx{ ["RoleA", "RoleX"], "XXXX", "resA", unconstrainedPrivileges }
	irrelevantResources Context   = TestCtx{ ["RoleA", "RoleX"], "ActB", "XXXX", unconstrainedPrivileges }
	zeroRoles Context             = TestCtx{ [],                 "ActB", "resA", unconstrainedPrivileges }
	zeroActions Context           = TestCtx{ ["RoleA", "RoleX"], nil,    "resA", unconstrainedPrivileges }
	zeroResources Context         = TestCtx{ ["RoleA", "RoleX"], "ActB", nil,    unconstrainedPrivileges }
	zeroPrivileges Context        = TestCtx{ ["RoleA", "RoleX"], "ActB", "resA", zeroPrivileges          }
)

func TestMay(t testing.T) {
	const mayCases = []struct{
		trbac   TRBAC
		context Context
		allowed Bool
		err     Error
	}{
		{ completeTRBAC, relevantUnconstrained, true,  nil},
		{ completeTRBAC, relevantConstrained,   false, nil },

		// privileges are irrelevant if any of roles, actions or
		// resources are irrelevant
		{ completeTRBAC, irrelevantRoles,      false, nil },
		{ completeTRBAC, irrelevantActions,    false, nil },
		{ completeTRBAC, irrelevantResources,  false, nil },

		// should fail without roles, actions, resources, and privileges
		// whether or not the TRBAC doesn lists them.
		// zero contexts are false
		{ completeTRBAC, zeroRoles,      false, nil },
		{ completeTRBAC, zeroActions,    false, nil },
		{ completeTRBAC, zeroResources,  false, nil },
		{ completeTRBAC, zeroPrivileges, false, nil },

		// empty TRBACs are errors
		{ emptyRolesTRBAC,      relevantUnconstrained, false, emptyRolesError      },
		{ emptyActionsTRBAC,    relevantUnconstrained, false, emptyActionsError    },
		{ emptyResourcesTRBAC,  relevantUnconstrained, false, emptyResourcesError  }

		// empty TRBACs and zero contexts are still errors and false
		{ emptyRolesTRBAC,      zeroRoles,      false, emptyRolesError      },
		{ emptyActionsTRBAC,    zeroActions,    false, emptyActionsError    },
		{ emptyResourcesTRBAC,  zeroResources,  false, emptyResourcesError  },
		{ emptyPrivilegesTRBAC,	zeroPrivileges,	false, emptyPrivilegesError },
	}

	for _, mc := range mayCases {
		allowed, err := trbac.may(mc.context)
		if allowed != mc.allowed {
			t.Errorf("may(%q) => %q, _. want %q, _",
				mc.ctx,
				allowed,
				mc.allowed,
			)
		}
		if err != mc.err {
			t.Errorf("may(%q) => _, %q. want _, %q",
				mc.ctx,
				err,
				mc.err,
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

func TestRequestPrivilegesProvider(t testing.T) {

}

func TestShellScriptConstraintRunner(t testing.T) {

}
