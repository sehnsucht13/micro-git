// Commands related to micro-git hash-object.
// Support is provided for:
// - "-w" flag
// - "--stdin" flag
// - "--stdin-paths" flag
// - Regular usage without any flag
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func hashObjectAtPath(path string) (string, error) {
	readPermission, err := IsReadable(path)
	if err != nil || !readPermission {
		return "", errors.New("File does not exist or cannot be read.")
	}
	fileContents, _ := ioutil.ReadFile(path)
	return hashString(string(fileContents)), nil
}

func hashObjectStdin() {
	var inputString strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inputString.WriteString(scanner.Text())
	}
	fmt.Println(hashString(inputString.String()))
}

func hashObjectStdinPath() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		filePath := scanner.Text()
		fileHash, err := hashObjectAtPath(filePath)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(fileHash)
	}
}

func HashObject(filePaths []string, writeFlag, stdin, stdinPath bool) {
	if stdin {
		hashObjectStdin()
		os.Exit(0)
	} else if stdinPath {
		os.Exit(1)
	}

	if len(filePaths) > 0 {
		for _, path := range filePaths {
			hash, err := hashObjectAtPath(path)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(hash)
			}
		}
	}
}
