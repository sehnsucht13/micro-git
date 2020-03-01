package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func hashFile(fileName string) string {
	h := sha1.New()
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Cannot access file", file, "Check if the file exists and access permitions are set.")
		os.Exit(2)
	}
	defer file.Close()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func compressFileContents(fileName string) []byte {
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Could not open file!")
	}
	var compressedBuff bytes.Buffer
	zWriter := zlib.NewWriter(&compressedBuff)
	zWriter.Write(contents)
	zWriter.Close()
	return compressedBuff.Bytes()
}

// Add files to index
func AddFiles(filePaths ...string) {
	repoRoot, err := FindRepoRoot()
	if err != nil {
		fmt.Println("Not in a repo")
		os.Exit(1)
	}

	for _, filePath := range filePaths {
		_, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("File or directory", filePath, "cannot be added.")
		} else {
			sha1String := hashFile(filePath)
			compFile := compressFileContents(filePath)
			// Create folder to hold new blob/tree and add compressed file
			blobFolderPath := filepath.Join(repoRoot, ".micro-git", sha1String[0:2])
			os.Mkdir(blobFolderPath, 0755)
			ioutil.WriteFile(filepath.Join(blobFolderPath, sha1String[2:]), compFile, 0755)
		}
	}
}
