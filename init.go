package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	MicroGitDir string = ".micro-git"
	ObjectDir   string = "objects"
	BranchDir   string = "branches"
	IndexFile   string = "index"
	HeadFile    string = "HEAD"
	RefDir      string = "refs"
)

// initRepoContents creates all the necessary files for a new directory at the
// path indicated by repoFolder.
func initRepoContents(repoFolder string) {
	// Create all directories
	os.MkdirAll(filepath.Join(repoFolder, ObjectDir), 0755)
	os.Mkdir(filepath.Join(repoFolder, BranchDir), 0755)
	os.Mkdir(filepath.Join(repoFolder, RefDir), 0755)
	os.Mkdir(filepath.Join(repoFolder, RefDir, "heads"), 0755)
	os.Mkdir(filepath.Join(repoFolder, RefDir, "tags"), 0755)
	// Create files
	headContents := []byte("ref: refs/heads/master")
	ioutil.WriteFile(filepath.Join(repoFolder, HeadFile), headContents, 0755)
	indexFile, _ := os.Create(filepath.Join(repoFolder, IndexFile))
	defer indexFile.Close()

}

// Create a new bare repository for project located at dirPath
func InitRepo(dirPath string, quiet bool) {
	currPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Println("Cannot create absolute path to new repository.")
		os.Exit(1)
	}
	repoFolder := filepath.Join(currPath, MicroGitDir)
	if DirExists(repoFolder) {
		fmt.Println(".micro-git repository already exists in", currPath, "Empty repository cannot be initialized!")
		os.Exit(1)
	}

	initRepoContents(repoFolder)
	if !quiet {
		fmt.Println("Empty repository initialized in", string(currPath), "successfully!")
	}
	os.Exit(0)
}
