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

func getCommitParent() string {
	repoRoot, _ := FindRepoRoot()
	repoHeadPath := filepath.Join(repoRoot, MicroGitDir, "HEAD")
	contents, _ := ioutil.ReadFile(repoHeadPath)
	currentBranch := strings.Split(string(contents), " ")[1]
	branchPath := filepath.Join(repoRoot, MicroGitDir, currentBranch)
	headRef, err := ioutil.ReadFile(branchPath)
	if err != nil {
		return ""
	}
	fmt.Println(headRef)
	return string(headRef)
}

func updateCurrentCommit(commitHash string) error{
	repoRoot, _ := FindRepoRoot()
	repoHeadPath := filepath.Join(repoRoot, MicroGitDir, "HEAD")
	contents, _ := ioutil.ReadFile(repoHeadPath)
	currentBranch := strings.Split(string(contents), " ")[1]
	branchPath := filepath.Join(repoRoot, MicroGitDir, currentBranch)
	err := ioutil.WriteFile(branchPath, []byte(commitHash), 0644)
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

func buildCommitFileContents(treeHash, commitMsg, commitParent string) string {
	// Temp variables
	name := "Yavor"
	email := "Hello@gmail.comg"
	return fmt.Sprintf("tree %s\nparent %s\nauthor %s\nemail %s\n%s", treeHash, commitParent, name, email, commitMsg)
}

func createCommit() error {
	commitTree := buildCommitTree()
	indexEntry, err := addObjectFile(commitTree)
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
	commitParent := getCommitParent()
	commitFileContets := buildCommitFileContents(indexEntry.Hash(), string(commitMessage), commitParent)
	fmt.Println(commitFileContets)
	return nil
}
