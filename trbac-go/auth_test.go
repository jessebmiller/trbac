package trbac

import "testing"

func TestBasicAuth(t *testing.T) {
    testPerms := map[string][]Permission{
        "reader": {
            {
                Actions:       []string{"read"},
                ResourceTypes: []string{"document"},
                Constraints:   []string{},
            },
        },
        "writer": {
            {
                Actions:       []string{"read", "write"},
                ResourceTypes: []string{"document"},
                Constraints:   []string{"business_hours"},
            },
        },
    }

    privs := TOMLPrivileges{testPerms}

    alwaysPass := NewFuncMapConstraintRunner(map[string]func(Context) bool{
        "business_hours": func(ctx Context) bool { return true },
    })

    alwaysFail := NewFuncMapConstraintRunner(map[string]func(Context) bool{
        "business_hours": func(ctx Context) bool { return false },
    })

    authWithPass := BasicAuth{privs, alwaysPass}
    authWithFail := BasicAuth{privs, alwaysFail}

    tests := []struct {
        name     string
        auth     Auth
        ctx      Context
        expected bool
    }{
        {
            "Reader can read document",
            authWithPass,
            NewLiteralContext("read", "document", []string{"reader"}),
            true,
        },
        {
            "Reader cannot write document",
            authWithPass,
            NewLiteralContext("write", "document", []string{"reader"}),
            false,
        },
        {
            "Writer can read document during business hours",
            authWithPass,
            NewLiteralContext("read", "document", []string{"writer"}),
            true,
        },
        {
            "Writer cannot read document outside business hours",
            authWithFail,
            NewLiteralContext("read", "document", []string{"writer"}),
            false,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := test.auth.May(test.ctx)
            if result != test.expected {
                t.Errorf("Expected %v but got %v", test.expected, result)
            }
        })
    }
}
