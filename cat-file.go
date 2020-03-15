package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// TODO: Error handling
func findFileByHash(hash string) (string, string, error) {
	var hashFilePath string
	var fullHashString string
	repoRoot, _ := FindRepoRoot()
	objectFolderPath := filepath.Join(repoRoot, ".micro-git", "objects")

	dirFiles, _ := ioutil.ReadDir(objectFolderPath)
	for _, file := range dirFiles {
		hashDirFile, _ := ioutil.ReadDir(filepath.Join(objectFolderPath, file.Name()))
		if len(hashDirFile) == 0 {
			return "", "", errors.New("micro-git objects folder located at " + objectFolderPath + " is corrupted!")
		}
		if strings.HasPrefix(strings.Join([]string{file.Name(), hashDirFile[0].Name()}, ""), hash) {
			if hashFilePath != "" {
				return "", "", errors.New("SHA1 hash provided is ambigous.")
			}
			hashFilePath = filepath.Join(objectFolderPath, file.Name(), hashDirFile[0].Name())
			fullHashString = strings.Join([]string{file.Name(), hashDirFile[0].Name()}, "")
		}
	}
	if hashFilePath == "" {
		return "", "", errors.New("File corresponding to provided SHA1 hash could not be found!")
	}
	return fullHashString, hashFilePath, nil
}

func decompressFileContents(path string) (string, error) {
	fileContents, _ := os.Open(path)
	var out bytes.Buffer

	reader, err := zlib.NewReader(fileContents)
	if err != nil {
		return "", errors.New("Unable to read zlib compressed file due to: " + err.Error())
	}
	defer reader.Close()
	io.Copy(&out, reader)
	return string(out.Bytes()), nil
}

func printFileContents(hash string) (string, error) {
	_, filePath, err := findFileByHash(hash)
	if err != nil {
		return "", err
	}
	fileContent, err := decompressFileContents(filePath)
	if err != nil {
		return "", err
	}
	return fileContent, nil
}

func getObjectSize(hash string) int64 {
	_, path, err := findFileByHash(hash)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	if file, _ := os.Stat(path); err == nil {
		return file.Size()
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
	return -1
}

func getObjectStatus(hash string) int {
	_, _, err := findFileByHash(hash)
	if err != nil {
		return 1
	}
	return 0
}

func batchProcess(hash string) (string, error) {
	fullHash, path, err := findFileByHash(hash)
	if err != nil {
		return "", err
	}
	fileSize := getObjectSize(hash)
	if fileSize == -1 {
		return "", errors.New("Object file size could not be determined!")
	}
	fileContents, err := decompressFileContents(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s \t %d\n %s", fullHash, fileSize, fileContents), nil
}

// CatFile reads a file identified by hash.
func CatFile(size, pp, status, batch bool, hash []string) {
	// Check if more than one flag was selected
	truthValCount := 0
	for _, truthVal := range []bool{size, pp, status, batch} {
		if truthVal {
			truthValCount++
		}
	}
	if truthValCount > 1 {
		fmt.Println("Can only select one option at a time.")
		os.Exit(1)
	} else if truthValCount == 0 {
		fmt.Println("Need to select an option for cat-file!")
		os.Exit(1)
	} else if (size || pp || status) && len(hash) == 0 {
		fmt.Println("<object> argument missing.")
		os.Exit(1)
	}

	if pp {
		decompStr, err := printFileContents(hash[0])
		if err != nil {
			fmt.Println("Cannot open object identifed by hash.", err.Error())
			os.Exit(1)
		}
		fmt.Println(decompStr)
		os.Exit(0)
	} else if size {
		size := getObjectSize(hash[0])
		fmt.Println(size)
	} else if status {
		exitStatus := getObjectStatus(hash[0])
		os.Exit(exitStatus)
	} else if batch {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			displayString, err := batchProcess(scanner.Text())
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println(displayString)
		}
	}
}
