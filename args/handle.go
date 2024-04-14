package args

import (
	"os"
)

type Argument struct {
	Name         string // Name of the argument
	Value        string // Value of the argument
	DefaultValue string // Default value if not provided
}

// ParseArguments parses command-line arguments and populates the Argument instances.
func parseArguments(args []string) map[string]string {
	argMap := make(map[string]string)

	for i := 1; i < len(args); i += 2 {
		argMap[args[i]] = args[i+1]
	}

	return argMap
}

func HandleArgs() (arguments []Argument) {
	arguments = []Argument{
		{Name: "-o", DefaultValue: "out.asm"},
		{Name: "run", DefaultValue: "."},
		{Name: "build", DefaultValue: "."},
	}

	// Parse command-line arguments
	argMap := parseArguments(os.Args)

	// Populate argument values
	for i, arg := range arguments {
		value, exists := argMap[arg.Name]
		if !exists {
			value = arg.DefaultValue
		}
		arguments[i].Value = value
	}

	return arguments
}
