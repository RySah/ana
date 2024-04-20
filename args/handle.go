package args

import (
	"fmt"
	"os"
)

func Exception(errArg *string, msg string) {
	if errArg == nil {
		fmt.Printf("Argument Exception: %s\n", msg)
	} else {
		fmt.Printf("Argument Exception: %s\n\t\"%s\"\n", msg, *errArg)
	}
}

type Argument struct {
	Name  string  // Name of the argument
	Value *string // Value of the argument
}

// ParseArguments parses command-line arguments and populates the Argument instances.
func parseArguments(args []string) map[string]string {
	argMap := make(map[string]string)

	for i := 1; i < len(args); i += 2 {
		if i+1 >= len(args) {
			Exception(&args[i], fmt.Sprintf("Expected: ... %s [ARGUMENT]  ...", args[i]))
		}
		argMap[args[i]] = args[i+1]
	}

	return argMap
}

const (
	OUTPUT = "-o"
	RUN    = "run"
	BUILD  = "build"
	ARCH   = "-a"
)

func Exists(args *[]Argument, name string) bool {
	return GetValueOf(args, name) != nil
}
func GetValueOf(args *[]Argument, name string) *string {
	for _, arg := range *args {
		if arg.Name == name {
			return arg.Value
		}
	}
	return nil
}

const minimumArgC = 1

func HandleArgs() []Argument {
	if len(os.Args) == minimumArgC {
		welcome()
	}

	arguments := []Argument{
		{Name: OUTPUT},
		{Name: RUN},
		{Name: BUILD},
		{Name: ARCH},
	}

	// Parse command-line arguments
	argMap := parseArguments(os.Args)

	// Populate argument values
	for i, arg := range arguments {
		value, exists := argMap[arg.Name]
		if !exists {
			arguments[i].Value = nil
		} else {
			arguments[i].Value = &value
		}
	}

	return arguments
}
