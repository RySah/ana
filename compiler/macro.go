package compiler

import (
	"fmt"
	"strings"
)

var macroTotalExpr = ""
var macroTotalExprEnd = ""

const (
	MACRO_VALUE_AT = "valueAt"
)

var macroNames = [...]string{
	MACRO_VALUE_AT, // Value at the address ...
}
var macroNoArgs = [...]string{ // Macros that dont require arguments
}

func isMacroName(name string) bool {
	for _, mn := range macroNames {
		if name == mn {
			return true
		}
	}
	return false
}
func macroNeedsNoArgs(name string) bool {
	for _, mn := range macroNoArgs {
		if name == mn {
			return true
		}
	}
	return false
}
func handleMacro(lineNumber int, name string, a ...string) {
	switch name {
	case MACRO_VALUE_AT:
		if len(a) == 0 {
			CompilerException(fmt.Sprintf("LINE %d - Macro expects an argument", lineNumber))
		}
		macroTotalExpr += fmt.Sprintf("[%s]", strings.Join(a, " "))
	}
}
