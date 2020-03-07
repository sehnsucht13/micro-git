// File stores all functions related to the micro-git config command.
package main

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	// 	"log"
	"fmt"
	"os"
	"path/filepath"
)

type Config map[interface{}]interface{}

type YAMLMarshalConfig struct {
	Email          string
	Name           string
	ColorUI        string
	EditorCmd      string
	Aliases        []string
	CommitTemplate string
	AutoCorrect    string
}

func mergeConfigValues(userConfig, localConfig Config) Config {
	mergedConfig := make(Config)
	for field, value := range userConfig {
		mergedConfig[field] = value
	}

	for field, value := range localConfig {
		if value == "" {
			mergedConfig[field] = userConfig[field]
		} else {
			mergedConfig[field] = value
		}
	}
	return mergedConfig
}

func GetConfigValues() (Config, error) {
	userConfig := make(Config)
	localConfig := make(Config)

	repoRoot, err := FindRepoRoot()
	if err != nil {
		return make(Config), err
	}
	userConfigContent, userConfigErr := ioutil.ReadFile(filepath.Join(repoRoot, ".micro-git", "config"))
	userHomeDir, _ := os.UserHomeDir()
	localConfigContent, localConfigErr := ioutil.ReadFile(filepath.Join(userHomeDir, ".config", "micro-git", "config"))
	if userConfigErr != nil && localConfigErr != nil {
		return make(Config), errors.New("No configuration files found!")
	} else if userConfigErr != nil {
		err = yaml.Unmarshal(localConfigContent, &localConfig)
		return localConfig, err
	} else if localConfigErr != nil {
		err = yaml.Unmarshal(userConfigContent, &userConfig)
		return userConfig, err
	} else {
		userConfigErr = yaml.Unmarshal(userConfigContent, &userConfig)
		localConfigErr = yaml.Unmarshal(localConfigContent, &localConfig)
		if userConfigErr != nil && localConfigErr != nil {
			return make(Config), errors.New("Could not parse configuration files! Corrupted!")
		} else if userConfigErr != nil {
			return localConfig, nil
		} else if localConfigErr != nil {
			return userConfig, nil
		} else {
			return mergeConfigValues(userConfig, localConfig), nil
		}

	}
}

// // UpdateConfigValues updates the values in the configuration.
// func UpdateConfigField(fieldName, newValue string) error {

// }

// func overwriteConfig(newContent string) error {

// }
