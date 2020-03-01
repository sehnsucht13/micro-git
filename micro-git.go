package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	initCmdQuiet := initCmd.Bool("quiet", false, "Suppress all text output to stdout except errors.")
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
		//hashFile(os.Args[2])
		//compressFileContents(os.Args[2])
		//file, err := FindRepoRoot()
	default:
		fmt.Println("Invalid Command!")
		os.Exit(1)
	}
}
