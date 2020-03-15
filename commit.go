package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Commits will be "simplified" for the time being. Instead of creating trees
// made up of files, commits will be made up of one file holding the current
// contents of the index.

const commitMsgFile = "COMMIT_MSG"

func buildCommitTree() string {
	var tree strings.Builder
	fileIndex := GetIndexEntries()
	for _, indexEntry := range fileIndex {
		tree.WriteString(indexEntry.String())
	}
	return tree.String()
}

func getCommitParent() (string, error) {
	repoRoot, _ := FindRepoRoot()
	repoHeadPath := filepath.Join(repoRoot, MicroGitDir, "HEAD")
	readStatus, err := IsReadable(repoHeadPath)
	if err != nil || !readStatus {
		return "", errors.New("Cannot read parent commit.")
	}
	contents, _ := ioutil.ReadFile(repoHeadPath)
	currentBranch := strings.Split(string(contents), " ")[1]
	branchPath := filepath.Join(repoRoot, MicroGitDir, currentBranch)

	readStatus, err = IsReadable(branchPath)
	if err != nil || !readStatus {
		return "", errors.New("Cannot read parent commit.")
	}
	// No need to check for error. We know that file is readable and it exists
	headRef, _ := ioutil.ReadFile(branchPath)
	return string(headRef), nil
}

func updateCurrentCommit(commitHash string) error {
	repoRoot, _ := FindRepoRoot()
	repoHeadPath := filepath.Join(repoRoot, MicroGitDir, "HEAD")
	readStatus, err := IsReadable(repoHeadPath)
	if err != nil || !readStatus {
		return errors.New("Cannot read HEAD file.")
	}

	contents, _ := ioutil.ReadFile(repoHeadPath)
	currentBranch := strings.Split(string(contents), " ")[1]
	branchPath := filepath.Join(repoRoot, MicroGitDir, currentBranch)
	writeStatus, err := IsWritable(branchPath)
	if err != nil || !writeStatus {
		return errors.New("Cannot update branch head. File is not writable.")
	}
	err = ioutil.WriteFile(branchPath, []byte(commitHash), 0644)
	return err
}

// Start user selected editor to create commit message
func writeCommitMessage() error {
	repoRoot, _ := FindRepoRoot()
	commitMsgPath := filepath.Join(repoRoot, MicroGitDir, commitMsgFile)
	os.Remove(commitMsgPath)
	cmd := exec.Command("vim", commitMsgPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	exitCode := cmd.Run()
	return exitCode
}

// Retrieve the text for the commit message as string
func getCommitMessageContents() ([]byte, error) {
	repoRoot, _ := FindRepoRoot()
	commitMsgPath := filepath.Join(repoRoot, MicroGitDir, commitMsgFile)
	return ioutil.ReadFile(commitMsgPath)
}

func buildCommitFileContents(userName, userEmail, treeHash, commitMsg, commitParent string) string {
	return fmt.Sprintf("tree %s\nparent %s\nauthor %s\nemail %s\n%s", treeHash, commitParent, userName, userEmail, commitMsg)
}

func createCommit() error {
	userEmail, emailError := findConfigValue("user.email")
	userName, nameError := findConfigValue("user.name")
	if emailError != nil || nameError != nil {
		return errors.New("Cannot locate user name or email!")
	}

	commitTree := buildCommitTree()
	treeIndexEntry, err := addObjectFile(commitTree)
	if err != nil {
		return err
	}
	err = writeCommitMessage()
	if err != nil {
		return errors.New("Could not write commit message successfully.")
	}
	commitMessage, err := getCommitMessageContents()
	if err != nil {
		return errors.New("Could not retrieve commit message contents.")
	}
	commitParent, err := getCommitParent()
	if err != nil {
		return err
	}
	commitFileContent := buildCommitFileContents(userName, userEmail, treeIndexEntry.Hash(), string(commitMessage), commitParent)
	commitIndexEntry, err := addObjectFile(commitFileContent)
	if err != nil {
		return err
	}
	updateCurrentCommit(commitIndexEntry.Hash())
	return nil
}
