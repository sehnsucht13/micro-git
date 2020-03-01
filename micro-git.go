package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func dirExists(dirPath string) bool {
	_, err := os.Stat(dirPath)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func InitRepo(dirPath string, quiet bool) {
	currPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Println("Error")
	}
	repoFolder := filepath.Join(currPath, ".micro-git")
	if dirExists(repoFolder) == false {
		// Create all directories
		os.MkdirAll(filepath.Join(repoFolder, "objects"), 0755)
		os.Mkdir(filepath.Join(repoFolder, "branches"), 0755)
		os.Mkdir(filepath.Join(repoFolder, "refs"), 0755)
		os.Mkdir(filepath.Join(repoFolder, "refs", "heads"), 0755)
		os.Mkdir(filepath.Join(repoFolder, "refs", "tags"), 0755)
		// Create files
		if !quiet {
			fmt.Println("Repository initialized in", string(currPath), "successfully!")
		}
		os.Exit(0)
	}
	fmt.Println(".micro-git folder already exists in", currPath, "Repository cannot be initialized!")
	os.Exit(1)
}

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	initCmdQuiet := initCmd.Bool("quiet", false, "Suppress all text output to stdout except errors.")
	flag.Parse()
	if len(os.Args) < 2{
		fmt.Println("Not enough arguments.")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		InitRepo(".", *initCmdQuiet)
	default:
		fmt.Println("Invalid Command!")
		os.Exit(1)
	}
}
