package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func DirExists(dirPath string) bool {
	_, err := os.Stat(dirPath)

	if os.IsNotExist(err) {
		return false
	}
	// No need to check if the file with the same name is a folder of
	// a file. Files and folders cannot have the same names.
	return true
}

func FindRepoRoot() (string, error) {
	workDir, _ := os.Getwd()
	userHomeDir, _ := os.UserHomeDir()
	for workDir != userHomeDir {
		if _, err := os.Stat(filepath.Join(workDir, ".micro-git")); err == nil {
			return workDir, nil
		}
		workDir = filepath.Dir(workDir)
	}
	return "", errors.New("Repository root directory cannot be found.")
}

func IsRepo() bool {
		_, err := FindRepoRoot()
		if err != nil{
			return false
		}
		return true
}

func FindRelPath(filePath string) string {
	repoRoot, err := FindRepoRoot()
	if err != nil {
		fmt.Println("Failed to find repo root!")
		os.Exit(2)
	}
	rel, _ := filepath.Rel(repoRoot, filePath)
	return rel

}
