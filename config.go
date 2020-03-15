// File stores all functions related to the micro-git config command.
package main

import (
	"encoding/json"
	"errors"
	"fmt"

	// "fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type configLevel int

// iota reset:
const (
	localLevel configLevel = iota
	userLevel
	systemLevel
)

func readConfigFile(filePath string) (map[string]interface{}, error) {
	fileContents, err := ioutil.ReadFile(filepath.Join(filePath, "micro-gitconfig"))
	if err != nil {
		return make(map[string]interface{}), errors.New("Cannot open configuration file.")
	}
	var configValues map[string]interface{}
	if err := json.Unmarshal(fileContents, &configValues); err != nil {
		return make(map[string]interface{}), err
	} else {
		return configValues, nil
	}
}

func getConfig(level configLevel) (map[string]interface{}, error) {
	switch level {
	case localLevel:
		repoPath, _ := FindRepoRoot()
		return readConfigFile(filepath.Join(repoPath, MicroGitDir))
	case userLevel:
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return make(map[string]interface{}), errors.New("Cannot find user home directory.")
		}
		return readConfigFile(homeDir)
	case systemLevel:
		return readConfigFile("/etc/micro-gitconfig")
	default:
		return make(map[string]interface{}), errors.New("Invalid level.")
	}
}

func overwriteConfig(level configLevel, contents map[string]interface{}) {
	configBytes, _ := json.Marshal(contents)
	switch level {
	case localLevel:
		repoPath, _ := FindRepoRoot()
		ioutil.WriteFile(filepath.Join(repoPath, MicroGitDir, "config"), configBytes, 0644)
	case userLevel:
		homeDir, _ := os.UserHomeDir()
		ioutil.WriteFile(filepath.Join(homeDir, "micro-gitconfig"), configBytes, 0644)
	case systemLevel:
		ioutil.WriteFile(filepath.Join("/etc", "micro-gitconfig"), configBytes, 0644)
	}
}

func ConfigListValues(level configLevel){
	configValues, err := getConfig(level)
	if err != nil{
		fmt.Println("Cannot retrieve configuration values")
	}
	for k, v := range configValues{
		fmt.Println(k,v)
	}
}
