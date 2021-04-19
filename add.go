package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// Create a unique hash for string using SHA1
func hashString(stringContent string) string {
	h := sha1.New()
	io.WriteString(h, stringContent)
	return hex.EncodeToString(h.Sum(nil))
}

// Comress a string using zlib and return its compressed form
func compressString(stringContent string) []byte {
	var compressedBuff bytes.Buffer
	zWriter := zlib.NewWriter(&compressedBuff)
	zWriter.Write([]byte(stringContent))
	zWriter.Close()
	return compressedBuff.Bytes()
}

// Compress and store a new file containing fileContent in the .micro-git/objects folder
// Returns an IndexEntry containing its hash.
// If this function is called from a directory which is not a micro-git repository
// then return error
func addObjectFile(fileContent string) (IndexEntry, error) {
	index_entry := IndexEntry{entry_type: blob}
	repoRoot, _ := FindRepoRoot()
	sha1String := hashString(fileContent)
	compString := compressString(fileContent)

	// Create folder to hold new blob/tree and add compressed file
	blobFolderPath := filepath.Join(repoRoot, ".micro-git", "objects", sha1String[0:2])
	os.Mkdir(blobFolderPath, 0755)
	ioutil.WriteFile(filepath.Join(blobFolderPath, sha1String[2:]), compString, 0755)
	index_entry.SetHash(sha1String)
	return index_entry, nil
}

func storeFile(absPath, relPath string) (IndexEntry, error) {
	file_stat, err := os.Stat(absPath)
	if err != nil {
		return IndexEntry{}, errors.New("File cannot be accessed " + relPath)
	}
	file_contents, err := ioutil.ReadFile(absPath)
	if err != nil {
		return IndexEntry{}, errors.New("File cannot be accessed " + relPath)
	}

	index_entry, err := addObjectFile(string(file_contents))
	index_entry.SetPerm("0" + strconv.FormatUint(uint64(file_stat.Mode().Perm()), 8))
	index_entry.SetName(relPath)
	return index_entry, nil
}

func addDirectory(dirAbsPath, dirRelPath string) ([]IndexEntry, error) {
	var indexEntries []IndexEntry
	dirFiles, _ := ioutil.ReadDir(dirAbsPath)
	for _, file := range dirFiles {
		if fileStat, err := os.Stat(filepath.Join(dirAbsPath, file.Name())); err == nil {
			if fileStat.IsDir() {
				subDirEntries, _ := addDirectory(filepath.Join(dirAbsPath, file.Name()), filepath.Join(dirRelPath, file.Name()))
				indexEntries = append(indexEntries, subDirEntries...)
			} else {
				entry, _ := storeFile(filepath.Join(dirAbsPath, file.Name()), filepath.Join(dirRelPath, file.Name()))
				indexEntries = append(indexEntries, entry)
			}
		}
	}
	return indexEntries, nil
}

func AddFiles(filePaths []string) {
	var indexEntries []IndexEntry
	if !IsRepo() {
		fmt.Println("Not in a repo")
		return
	}
	for _, path := range filePaths {
		absPath, _ := filepath.Abs(path)
		if fileStat, err := os.Stat(absPath); err == nil {
			relPath := FindRelPath(absPath)
			if fileStat.IsDir() {
				fmt.Println(absPath, relPath)
				dirEntries, _ := addDirectory(absPath, relPath)
				indexEntries = append(indexEntries, dirEntries...)
			} else {
				entry, _ := storeFile(absPath, relPath)
				indexEntries = append(indexEntries, entry)
			}
		}
	}
	// TODO: Remove
	fmt.Println("Index Entries", indexEntries)
	AddIndexEntries(indexEntries)
}
