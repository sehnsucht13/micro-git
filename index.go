package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileType int

// File type of an index
const (
	blob FileType = iota
	tree FileType = iota
)

// Structure representing a single Index entry
type IndexEntry struct {
	perm       string
	entry_type FileType
	sha1_hash  string
	entry_name string
}

func (i IndexEntry) Name() string {
	return i.entry_name
}

func (i IndexEntry) Hash() string {
	return i.sha1_hash
}

func (i *IndexEntry) SetName(name string) {
	i.entry_name = name
}

func (i *IndexEntry) SetHash(hash string) {
	i.sha1_hash = hash
}

func (i *IndexEntry) SetPerm(perm string) {
	i.perm = perm
}

func (i IndexEntry) String() string {
	var file_type string
	if i.entry_type == blob {
		file_type = "blob"
	} else {
		file_type = "tree"
	}
	return fmt.Sprintf("%s %s %s %s\n", i.perm, file_type, i.sha1_hash, i.entry_name)
}

// Implement interface to sort IndexEntry slices in place
type ByEntry []IndexEntry

func (a ByEntry) Len() int           { return len(a) }
func (a ByEntry) Less(i, j int) bool { return a[i].Name() < a[j].Name() }
func (a ByEntry) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Search for an IndexEntry corresponding to file entry_name.
// If a corresponding entry is found then return its index and sha1 hash.
// If an entry is not found then return error.
func findIndexEntry(entry_name string, entries []IndexEntry) (int, string, error) {
	numEntries := len(entries)

	idx := sort.Search(numEntries, func(i int) bool { return entries[i].Name() == entry_name })
	if idx < numEntries && entries[idx].Name() == entry_name {
		return idx, entries[idx].Hash(), nil
	}
	return -1, "", errors.New("Element does not exist")
}

func GetIndexEntries() []IndexEntry {
	var index_entries []IndexEntry
	rootPath, err := FindRepoRoot()
	if err != nil {
		fmt.Println("Not in a repository")
		os.Exit(2)
	}

	file, err := ioutil.ReadFile(filepath.Join(rootPath, ".micro-git", "index"))
	if err != nil {
		fmt.Println("Cannot open index file")
	}
	index_lines := strings.Split(string(file), "\n")
	// Remove the last element which is an empty string appened by the Split command
	index_lines = index_lines[:len(index_lines)-1]
	for _, line := range index_lines {
		index_contents := strings.Split(line, " ")
		var entry_type FileType
		if index_contents[1] == "blob" {
			entry_type = blob
		} else {
			entry_type = tree
		}
		curr_entry := IndexEntry{index_contents[0], entry_type, index_contents[2], index_contents[3]}
		index_entries = append(index_entries, curr_entry)
	}
	return index_entries
}

func AddIndexEntries(new_entries []IndexEntry) {
	curr_entries := GetIndexEntries()
	sortEntries := false
	for _, entry := range new_entries {
		idx, hash, err := findIndexEntry(entry.Name(), curr_entries)
		if err != nil {
			// Entry does not exist. Simply append it to current entries.
			curr_entries = append(curr_entries, entry)
			sortEntries = true

		} else {
			// If the hashes of the entry are not equivalent to the one in the index
			// then they need to be updated.
			if hash != entry.Hash() {
				curr_entries[idx] = entry
			}
		}
	}

	if sortEntries {
		sort.Sort(ByEntry(curr_entries))
	}

	rootPath, err := FindRepoRoot()
	if err != nil {
		fmt.Println("Not in a repository")
		os.Exit(2)
	}
	// Overwrite the entire index with the new entries
	indexPath := filepath.Join(rootPath, ".micro-git", "index")
	os.Remove(indexPath)
	indexFile, _ := os.Create(indexPath)
	for _, entry := range curr_entries {
		indexFile.WriteString(entry.String())
	}
}
