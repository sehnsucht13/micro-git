package main

import (
"strings"
"fmt"
"errors"
"path/filepath"
"io/ioutil"
"os"
)

// TODO: Error handling
func findFileByHash(hash string) (string, error){
	var hashFilePath string
	repoRoot, _ := FindRepoRoot()
	objectFolderPath := filepath.Join(repoRoot, ".micro-git", "objects")

	dirFiles, _ := ioutil.ReadDir(objectFolderPath)
	for _, file  := range dirFiles{
		hashDirFile, _ := ioutil.ReadDir(filepath.Join(objectFolderPath, file.Name()))
		fmt.Println( strings.Join([]string{file.Name(), hashDirFile[0].Name()},""))
		if len(hashDirFile) == 0{
			return "",errors.New("micro-git objects folder located at " + objectFolderPath + " is corrupted!")
		}
		if strings.HasPrefix(strings.Join([]string{file.Name(), hashDirFile[0].Name()},""), hash){
			if hashFilePath != ""{
				return "",errors.New("SHA1 hash provided is ambigous.")
			}
			hashFilePath = filepath.Join(objectFolderPath, file.Name(), hashDirFile[0].Name())
		}
	}
	if hashFilePath == ""{
		return "", errors.New("File corresponding to provided SHA1 hash could not be found!")
	}
	fmt.Println("Found file with ", hashFilePath)
	return hashFilePath, nil
}

func getObjectSize(hash string) int64{
	path, err := findFileByHash(hash)
	if err != nil{
		fmt.Println(err)
		return -1
	}
	if file, _ := os.Stat(path); err == nil {
		return file.Size()
	}else{
		fmt.Println(err)
		os.Exit(1)
	}
	return -1
}
