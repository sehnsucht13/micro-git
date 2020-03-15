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
	"strings"
)

type configLevel int

// iota reset:
const (
	localLevel configLevel = iota
	globalLevel
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
	case globalLevel:
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
	case globalLevel:
		homeDir, _ := os.UserHomeDir()
		ioutil.WriteFile(filepath.Join(homeDir, "micro-gitconfig"), configBytes, 0644)
	case systemLevel:
		ioutil.WriteFile(filepath.Join("/etc", "micro-gitconfig"), configBytes, 0644)
	}
}

func ConfigListValues(level configLevel) {
	configValues, err := getConfig(level)
	if err != nil {
		fmt.Println("Cannot retrieve configuration values")
	}
	for k, v := range configValues {
		fmt.Println(k, v)
	}
}

func configSetValue(config map[string]interface{}, key, value string) (map[string]interface{}, error) {
	localCopy := config
	subKeys := strings.Split(key, ".")
	if len(subKeys) != 2 {
		return make(map[string]interface{}), errors.New("Invalid key provided.")
	}
	subConfig, subKeyPresent := (config[subKeys[0]]).(map[string]interface{})
	if !subKeyPresent {
		m := make(map[string]string)
		m[subKeys[1]] = value
		localCopy[subKeys[0]] = m
	}
	subConfig[subKeys[1]] = value
	localCopy[subKeys[0]] = subConfig
	return localCopy, nil
}

func configGetValue(config map[string]interface{}, key string) (string, error) {
	subKeys := strings.Split(key, ".")
	if len(subKeys) != 2 {
		return "", errors.New("Key does not exist.")
	}
	subConfig, subKeyPresent := (config[subKeys[0]]).(map[string]interface{})
	if !subKeyPresent {
		return "", errors.New("Key does not exist!")
	}
	val, subSubKeyPresent := subConfig[subKeys[1]].(string)
	if !subSubKeyPresent {
		return "", errors.New("Key does not exist!")
	}
	return val, nil
}

func Config(list bool, key, value, level string) {
	var userConfig map[string]interface{}
	switch level {
	case "system":
		config, err := getConfig(systemLevel)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		userConfig = config
	case "global":
		config, err := getConfig(globalLevel)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		userConfig = config
	case "local":
		_, err := FindRepoRoot()
		if err != nil {
			fmt.Println("Not visiting a micro-git repository!")
			os.Exit(1)
		}
		config, err := getConfig(localLevel)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		userConfig = config
	// Case of a level not being chosen
	case "":
	default:
		fmt.Println("Unknown configuration level provided!")
		os.Exit(1)
	}
	fmt.Println(userConfig)
}
