package trbac

import (
    "os/exec"
    "path"
)

// ConstraintRunner evaluates named constraints against contexts
type ConstraintRunner interface {
    // Run evaluates a named constraint against a context
    // Returns true if the constraint passes, false otherwise
    Run(constraintName string, ctx Context) bool
}

// FuncMapConstraintRunner uses Go functions as constraints
type FuncMapConstraintRunner struct {
    constraintFuncs map[string]func(Context) bool
}

// Run evaluates a constraint against a context
func (runner FuncMapConstraintRunner) Run(constraint string, ctx Context) bool {
    fn, exists := runner.constraintFuncs[constraint]
    if !exists {
        return false
    }
    return fn(ctx)
}

// NewFuncMapConstraintRunner creates a new function-based constraint runner
func NewFuncMapConstraintRunner(funcs map[string]func(Context) bool) ConstraintRunner {
    return FuncMapConstraintRunner{funcs}
}

// ShellScriptConstraintRunner executes shell scripts as constraints
type ShellScriptConstraintRunner struct {
    scriptRoot string
}

// Run evaluates a constraint against a context
func (runner ShellScriptConstraintRunner) Run(constraint string, ctx Context) bool {
    scriptPath := path.Join(runner.scriptRoot, constraint)
    
    // Convert context to command arguments
    args := []string{ctx.Action(), ctx.ResourceType()}
    for _, role := range ctx.Roles() {
        args = append(args, role)
    }
    
    // Execute the script
    cmd := exec.Command(scriptPath, args...)
    err := cmd.Run()
    
    // Script succeeds if exit code is 0
    return err == nil
}

// NewShellScriptConstraintRunner creates a new shell script constraint runner
func NewShellScriptConstraintRunner(scriptRoot string) ConstraintRunner {
    return ShellScriptConstraintRunner{scriptRoot}
}
