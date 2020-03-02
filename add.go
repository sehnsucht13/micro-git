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
	"errors"
	"strconv"
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

func storeFile(filePath string) (IndexEntry, error){
	index_entry := IndexEntry{entry_type: blob}
	repoRoot, err := FindRepoRoot()
	if err != nil {
		fmt.Println("Not in a repo")
		os.Exit(1)
	}

	file, err := os.Stat(filePath)
	if err != nil {
		//fmt.Println("File located at", filePath, "cannot be added.")
		return index_entry, errors.New("File cannot be accessed " + file.Name())
	}
	sha1String := hashFile(filePath)
	compFile := compressFileContents(filePath)
	// Create folder to hold new blob/tree and add compressed file
	blobFolderPath := filepath.Join(repoRoot, ".micro-git", "objects", sha1String[0:2])
	os.Mkdir(blobFolderPath, 0755)
	ioutil.WriteFile(filepath.Join(blobFolderPath, sha1String[2:]), compFile, 0755)

	index_entry.SetPerm("0" + strconv.FormatUint(uint64(file.Mode().Perm()), 8))
	index_entry.SetHash(sha1String)
	index_entry.SetName(file.Name())
	return index_entry, nil
}

/*
func AddFiles(filePaths ...string) {
	repoRoot, err := FindRepoRoot()
	if err != nil {
		fmt.Println("Not in a repo")
		os.Exit(1)
	}
}
*/
