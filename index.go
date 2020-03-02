package main

import (
		"io/ioutil"
		"fmt"
		"os"
		"path/filepath"
		"strings"
		"sort"
		"errors"
)

type FileType int

// File type of an index
const(
	blob FileType = iota
	tree FileType = iota
)

// Structure representing a single index entry
type IndexEntry struct{
	perm string
	entry_type FileType
	sha1_hash string
	entry_name string
}

func (i IndexEntry) Name() string{
	return i.entry_name
}

func (i IndexEntry) Hash() string{
	return i.sha1_hash
}

func (i *IndexEntry) SetName(name string){
	i.entry_name = name
}


func (i *IndexEntry) SetHash(hash string){
	i.sha1_hash = hash
}

func (i *IndexEntry) SetPerm(perm string){
	i.perm = perm
}

func (i IndexEntry) String() string{
	var file_type string
	if i.entry_type == blob{
		file_type = "blob"
	}else{
		file_type = "tree"
	}
	return fmt.Sprintf("%s %s %s %s", i.perm, file_type, i.sha1_hash, i.entry_name)
}

// Sorts entries in place
func sortEntries(entries *[]IndexEntry) {
	sort.Slice(entries, func(i, j int) bool { return (*entries)[i].entry_name < (*entries)[j].entry_name })
}

// Search for an IndexEntry corresponding to file entry_name.
// If a corresponding entry is found then return its index and sha1 hash
// If an entry is not found then return error
func findEntry(entry_name string, entries []IndexEntry) (int, string, error){
	numEntries := len(entries)
	idx := sort.Search(numEntries, func(i int) bool{return entries[i].Name() == entry_name})
	if entries[idx].Name() == entry_name{
		return idx, entries[idx].Hash(), nil
	}
	return -1,"", errors.New("Element does not exist")
}

func GetIndexEntries() []IndexEntry{
	var index_entries []IndexEntry
	rootPath, err := FindRepoRoot()
	if err != nil{
		fmt.Println("Not in a repository")
		os.Exit(2)
	}

	file, err := ioutil.ReadFile(filepath.Join(rootPath,".micro-git", "index"))
	if err != nil{
		fmt.Println("Cannot open index file")
	}
	index_lines := strings.Split(string(file), "\n")
	// Remove the last element which is an empty string appened by the Split command
	index_lines = index_lines[:len(index_lines)-1]
	for _, line := range index_lines{
		index_contents := strings.Split(line, " ")
		var entry_type FileType
		if index_contents[1] == "blob"{
			entry_type = blob
		}else{
			entry_type = tree
		}
		curr_entry := IndexEntry{index_contents[0], entry_type, index_contents[2], index_contents[3]}
		index_entries = append(index_entries, curr_entry)
	}
	return index_entries
}

