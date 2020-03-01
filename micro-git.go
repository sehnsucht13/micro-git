package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"io"
	"crypto/sha1"
	"log"
	"encoding/hex"
	//"io/ioutil"
)

func hashFile(fileName string) string {
	h := sha1.New()
	file, err := os.Open(fileName)
	if err != nil{
		fmt.Println("Cannot access file", file,"Check if the file exists and access permitions are set.")
		os.Exit(2)
	}
	defer file.Close()
	if _,err := io.Copy(h, file); err != nil{
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

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
		//head_test := []byte("ref: refs/heads/master")
		if !quiet {
			fmt.Println("Empty repository initialized in", string(currPath), "successfully!")
		}
		os.Exit(0)
	}
	fmt.Println(".micro-git repository already exists in", currPath, "Empty repository cannot be initialized!")
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
	case "add":
		hashFile(os.Args[2])
	default:
		fmt.Println("Invalid Command!")
		os.Exit(1)
	}
}
