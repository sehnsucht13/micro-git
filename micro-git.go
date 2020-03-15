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

	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	statusCmdShortDisp := statusCmd.Bool("s", false, "Give the output in short format.")
	statusCmdLongDisp := statusCmd.Bool("long", true, "Give output in long format. This is the default.")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments.")
		flag.PrintDefaults()
		os.Exit(2)
	}
	if os.Args[1] == "init" {
		initCmd.Parse(os.Args[2:])
		switch len(initCmd.Args()) {
		case 0:
			cwdPath, err := os.Getwd()
			if err != nil {
				fmt.Println("Cannot find current working directory.")
				os.Exit(1)
			}
			InitRepo(cwdPath, *initCmdQuiet)
		case 1:
			InitRepo(initCmd.Args()[0], *initCmdQuiet)
		default:
			fmt.Println("micro-git init only accepts one repository path at a time!")
			os.Exit(1)
		}
	} else {
		_, err := FindRepoRoot()
		if err != nil {
			fmt.Println("Not currently visiting a micro-git repository.")
			os.Exit(1)
		}

		switch os.Args[1] {
		case "add":
			AddFiles(os.Args[2:])
		case "cat-file":
			if len(catCmd.Args()) > 1 {
				fmt.Println("cat-file accepts at most one argument.")
				os.Exit(1)
			}
			catCmd.Parse(os.Args[2:])
			CatFile(*catCmdSize, *catCmdPrettyPrint, *catCmdValid, *catCmdBatch, catCmd.Args())
		case "status":
			statusCmd.Parse(os.Args[2:])
			Status(*statusCmdShortDisp, *statusCmdLongDisp)
		case "config":
			fmt.Println(GetConfigValues())
		default:
			fmt.Println("Invalid Command!")
			os.Exit(1)
		}
	}
}
