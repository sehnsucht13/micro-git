package main
import (
"os"
"path/filepath"
"io/ioutil"
"fmt"
)

// Create a new bare repository for project located at dirPath
func InitRepo(dirPath string, quiet bool) {
	currPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Println("Error")
	}
	repoFolder := filepath.Join(currPath, ".micro-git")
	if DirExists(repoFolder) == false {
		// Create all directories
		os.MkdirAll(filepath.Join(repoFolder, "objects"), 0755)
		os.Mkdir(filepath.Join(repoFolder, "branches"), 0755)
		os.Mkdir(filepath.Join(repoFolder, "refs"), 0755)
		os.Mkdir(filepath.Join(repoFolder, "refs", "heads"), 0755)
		os.Mkdir(filepath.Join(repoFolder, "refs", "tags"), 0755)
		// Create files
		headContents := []byte("ref: refs/heads/master")
		ioutil.WriteFile(filepath.Join(repoFolder, "HEAD"), headContents, 0755)
		if !quiet {
			fmt.Println("Empty repository initialized in", string(currPath), "successfully!")
		}
		os.Exit(0)
	}
	fmt.Println(".micro-git repository already exists in", currPath, "Empty repository cannot be initialized!")
	os.Exit(1)
}
