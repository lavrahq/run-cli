package when

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/lavrahq/cli/util/cmdutil"
)

// ImplicitlyTrue checks if the value is implicitly true based
// the string representation.
func ImplicitlyTrue(program string) bool {
	return program == "" || program == "true" || program == "always"
}

// ImplicitlyFalse checks if the value is implicitly false based
// the string representation.
func ImplicitlyFalse(program string) bool {
	return program == "false" || program == "never"
}

// False checks that the program provided evaluates to false.
func False(program string, env interface{}) bool {
	return !True(program, env)
}

// Evaluate returns the raw result of the evaluated program.
func Evaluate(program string, env interface{}) interface{} {
	if ImplicitlyTrue(program) {
		return true
	}

	if ImplicitlyFalse(program) {
		return false
	}

	compiledProgram, err := expr.Compile(program, expr.Env(env))
	cmdutil.CheckCommandError(err, fmt.Sprintf("compiling `when`:\n %s", program))

	eval, err := expr.Run(compiledProgram, env)
	cmdutil.CheckCommandError(err, fmt.Sprintf("executing `when`:\n %s", program))

	return eval
}

// True checks that the program provided evaluates to true.
func True(program string, env interface{}) bool {
	if ImplicitlyTrue(program) {
		return true
	}

	if ImplicitlyFalse(program) {
		return false
	}

	if Evaluate(program, env) == true {
		return true
	}

	return false
}
