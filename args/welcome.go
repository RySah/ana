package args

import (
	"fmt"
	"os"
)

const (
	RESET       = "\x1B[0m"
	KEYWORD_COL = "\x1B[1m\x1B[93m"
)

func welcome() {
	fmt.Printf("%sANA%s, Abstracted Netwide Assembler\n", KEYWORD_COL, RESET)
	os.Exit(0)
}
