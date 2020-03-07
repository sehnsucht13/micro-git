package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// TODO: Figure out how to handle errors and files not being able to be accessed
// findAllRepoFiles retrieves the names and hashes of all files within the
// currently visited repo. This function returns a map with keys being the
// relative file paths to all files and the values being their hashed values
// using SHA1.
func findAllRepoFiles(repoRoot string) map[string]string {
	files := make(map[string]string)
	_ = filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			relPath := FindRelPath(path)
			// Ignore microgit folder
			if !strings.HasPrefix(relPath, ".micro-git") {
				fileData, err := ioutil.ReadFile(path)
				if err == nil {
					sha1String := hashString(string(fileData))
					files[relPath] = sha1String
				}
			}
		}
		return nil
	})
	return files
}

func checkFileIndex(files map[string]string) ([]string, []string) {
	indexEntries := GetIndexEntries()
	nonIndexedFiles := []string{}
	changedFiles := []string{}

	for fileName, shaHash := range files {
		_, hash, err := findIndexEntry(fileName, indexEntries)
		if err != nil {
			nonIndexedFiles = append(nonIndexedFiles, fileName)
		} else if hash != shaHash {
			changedFiles = append(changedFiles, fileName)
		}
	}
	return changedFiles, nonIndexedFiles
}

func shortStatusMsg(changedFiles, newFiles []string) {
	for _, modFilePath := range changedFiles {
		PrintColorSingleLine(ColorRed, fmt.Sprintf("M %s", modFilePath))
	}
	for _, untrackedFilePath := range newFiles {
		PrintColorSingleLine(ColorRed, fmt.Sprintf("?? %s", untrackedFilePath))
	}

}

func longStatusMsg(changedFiles, newFiles []string) {
	fmt.Println("Changes not staged for commit:")
	for _, modFilePath := range changedFiles {
		PrintColorSingleLine(ColorRed, fmt.Sprintf("\t %s", modFilePath))
	}
	fmt.Println("Untracked files:")
	for _, untrackedFilePath := range newFiles {
		PrintColorSingleLine(ColorRed, fmt.Sprintf("\t %s", untrackedFilePath))
	}

}

func Status(shortDisplay, longDisplay bool) error {
	repoRoot, err := FindRepoRoot()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	repoFiles := findAllRepoFiles(repoRoot)
	changedFiles, newFiles := checkFileIndex(repoFiles)

	if shortDisplay {
		shortStatusMsg(changedFiles, newFiles)
	} else if longDisplay {
		longStatusMsg(changedFiles, newFiles)
	}
	return nil
}
