package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	initCmdQuiet := initCmd.Bool("quiet", false, "Suppress all text output to stdout except errors.")

	// addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	// addVerboseCmd := addCmd.Bool("verbose", false, "Display the information about each new or updated entry in the index file.")
	// addDryRunCmd := addCmd.Bool("dry-run", false, "Display the information about each new or updated entry in the index file.")
	// refresh := addCmd.Bool("refresh", false, "Display the information about each new or updated entry in the index file.")

	catCmd := flag.NewFlagSet("cat-file", flag.ExitOnError)
	catCmdPrettyPrint := catCmd.Bool("p", false, "Print the decompressed contents of <object>.")
	catCmdSize := catCmd.Bool("s", false, "Print the size of <object> in bytes.")
	catCmdValid := catCmd.Bool("e", false, "Check if <object> exist and exit with status 0 if it does and status 1 otherwise.")
	catCmdBatch := catCmd.Bool("batch", false, "Batch process objects using stdin.")

	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments.")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		InitRepo(".", *initCmdQuiet)
	case "add":
		AddFiles(os.Args[2:])
	case "cat-file":
		if len(catCmd.Args()) > 1 {
			fmt.Println("Can use at most one argument.")
			os.Exit(1)
		}
		CatFile(*catCmdSize, *catCmdPrettyPrint, *catCmdValid, *catCmdBatch, catCmd.Args())
	default:
		fmt.Println("Invalid Command!")
		os.Exit(1)
	}
}
