package main

import (
	"ana/args"
	"fmt"
	"os"
)

func main() {
	// comp := *compiler.NewCompiler("hello")
	build, run, fp, op := args.HandleArgs()
	if build || run {
		fmt.Printf("BUILDING    %s\n", *fp)

		fmt.Printf("CREATED IN  %s\n", *op)
		//if run {
		//}
	}
	os.Exit(0)
}

func welcome() {
}
