package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"io"
	"io/ioutil"
	"crypto/sha1"
	"log"
	"encoding/hex"
	"compress/zlib"
	"bytes"
	"errors"
)

func FindRepoRoot() (string, error) {
	workDir, _ := os.Getwd()
	userHomeDir, _ := os.UserHomeDir()
	for workDir != userHomeDir{
		if _, err := os.Stat(filepath.Join(workDir, ".micro-git")); err == nil{
			return workDir, nil
		}
		workDir = filepath.Dir(workDir)
	}
	return "", errors.New("Repository root directory cannot be found.")
}

func FindRelPath(filePath string) string {
	repoRoot, err := FindRepoRoot()
	if err != nil{
		fmt.Println("Failed to find repo root!")
		os.Exit(2)
	}
	rel, _ := filepath.Rel(repoRoot, filePath)
	return rel

}

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

func compressFileContents(fileName string) []byte{
	contents, err := ioutil.ReadFile(fileName)
	if err != nil{
		fmt.Println("Could not open file!")
	}
	var compressedBuff bytes.Buffer
	zWriter := zlib.NewWriter(&compressedBuff)
	zWriter.Write(contents)
	zWriter.Close()
	return compressedBuff.Bytes()
}

// Add files to index
func AddFiles(filePaths ...string){
	repoRoot, err := FindRepoRoot()
	if err != nil{
		fmt.Println("Not in a repo")
		os.Exit(1)
	}

	for _, filePath := range filePaths{
		_, err := os.Stat(filePath)
		if err != nil{
			fmt.Println("File or directory", filePath, "cannot be added.")
		}else{
			sha1String := hashFile(filePath)
			compFile := compressFileContents(filePath)
			// Create folder to hold new blob/tree and add compressed file
			blobFolderPath := filepath.Join(repoRoot, ".micro-git", sha1String[0:2])
			os.Mkdir(blobFolderPath, 0755)
			ioutil.WriteFile(filepath.Join(blobFolderPath, sha1String[2:]), compFile, 0755)
		}
	}
}

func dirExists(dirPath string) bool {
	_, err := os.Stat(dirPath)

	if os.IsNotExist(err) {
		return false
	}
	// No need to check if the file with the same name is a folder of
	// a file. Files and folders cannot have the same names.
	return true
}

// Create a new bare repository for project located at dirPath
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
		//hashFile(os.Args[2])
		//compressFileContents(os.Args[2])
		//file, err := FindRepoRoot()
	default:
		fmt.Println("Invalid Command!")
		os.Exit(1)
	}
}
