package auth_test

import (
	"fmt"
	"github.com/jessebmiller/trbac/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstraintRunners(t *testing.T) {
	scriptRunner := auth.NewShellScriptConstraintRunner("test_constraint_scripts")
	testCases := []struct {
		runner     auth.ConstraintRunner
		constraint string
		context    auth.Context
		want       bool
	}{
		{
			scriptRunner,
			"always",
			mockContext{},
			true,
		}, {
			scriptRunner,
			"never",
			mockContext{},
			false,
		}, {
			scriptRunner,
			"if_ok",
			mockContext{"a", "r", []string{"ro"}},
			true,
		}, {
			scriptRunner,
			"if_ok",
			mockContext{"a", "x", []string{"ro"}},
			false,
		},
	}

	for _, testCase := range testCases {
		runner := testCase.runner
		constraint := testCase.constraint
		context := testCase.context
		want := testCase.want
		got := runner.Run(constraint, context)
		assert.Equal(
			t,
			got,
			want,
			fmt.Sprintf(
				"%T.Run(%s, %v) -> %v, want %v",
				runner,
				constraint,
				context,
				got,
				want,
			),
		)
	}
}
