package main

import (
	"errors"
	"fmt"
	"github.com/phayes/permbits"
	"os"
	"path/filepath"
)

// Utilities to print text using color
var colorNames = [...]string{
	"\033[0m",    // Clear/Reset color
	"\033[0;31m", // Red
	"\033[0;32m", // Green
	"\033[0;33m", // Yellow
	"\033[0;34m", // Blue
	"\033[0;35m", // Magenta
	"\033[0;36m", // Cyan
}

type Color int

// iota reset
const (
	ColorClear Color = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
)

func SetColor(color Color) {
	fmt.Printf("%s", colorNames[color])
}

func ResetColor() {
	fmt.Printf("%s", colorNames[ColorClear])
}

// PrintColor prints the string outputString to stdout using the specified color. The
// terminal colors are reset after every print
func PrintColorSingleLine(color Color, outputString string) {
	fmt.Printf("%s%s", colorNames[color], outputString)
	fmt.Printf("%s\n", colorNames[ColorClear])
}

// Functions related to finding repository root
func FindRepoRoot() (string, error) {
	workDir, _ := os.Getwd()
	userHomeDir, _ := os.UserHomeDir()
	for workDir != userHomeDir {
		if _, err := os.Stat(filepath.Join(workDir, ".micro-git")); err == nil {
			return workDir, nil
		}
		workDir = filepath.Dir(workDir)
	}
	return "", errors.New("Repository root directory cannot be found.")
}

func IsRepo() bool {
	_, err := FindRepoRoot()
	if err != nil {
		return false
	}
	return true
}

func FindRelPath(filePath string) string {
	repoRoot, err := FindRepoRoot()
	if err != nil {
		fmt.Println("Failed to find repo root!")
		os.Exit(2)
	}
	rel, _ := filepath.Rel(repoRoot, filePath)
	return rel
}

// Functions related to finding out file permissions/existance

// IsReadable checks if a file has read permissions for the current user.
func IsReadable(path string) (bool, error) {
	permissions, err := permbits.Stat(path)
	if err != nil {
		return false, err
	}
	return permissions.UserRead(), nil
}

// IsWritable checks if a file has write permissions for the current user.
func IsWritable(path string) (bool, error) {
	permissions, err := permbits.Stat(path)
	if err != nil {
		return false, err
	}
	return permissions.UserWrite(), nil
}

// IsReadwrite checks if a file has read and write permissions for the current user.
func IsReadWrite(path string) (bool, error) {
	permissions, err := permbits.Stat(path)
	if err != nil {
		return false, err
	}
	if permissions.UserWrite() && permissions.UserRead() {
		return true, nil
	}
	return false, nil
}

// FileExists checks and can be accessed if a file located at filePath
// exists. Returns true if the file exists and false otherwise.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return true

}

func DirExists(dirPath string) bool {
	_, err := os.Stat(dirPath)

	if os.IsNotExist(err) {
		return false
	}
	// No need to check if the file with the same name is a folder of
	// a file. Files and folders cannot have the same names.
	return true
}
