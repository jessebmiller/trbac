package trbac

import "testing"

func TestFuncMapConstraintRunner(t *testing.T) {
    constraints := map[string]func(Context) bool{
        "always_pass": func(ctx Context) bool { return true },
        "always_fail": func(ctx Context) bool { return false },
    }

    runner := NewFuncMapConstraintRunner(constraints)
    ctx := NewLiteralContext("read", "document", []string{"reader"})

    if !runner.Run("always_pass", ctx) {
        t.Error("Expected always_pass to return true")
    }

    if runner.Run("always_fail", ctx) {
        t.Error("Expected always_fail to return false")
    }

    if runner.Run("non_existent", ctx) {
        t.Error("Expected non_existent to return false")
    }
}

// Note: Testing ShellScriptConstraintRunner would require setting up shell scripts and is not included here.
