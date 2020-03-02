package main

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

func (i *IndexEntry) SetName(name string){
	i.entry_name = name
}

func (i *IndexEntry) SetHash(hash string){
	i.sha1_hash = hash
}

func (i *IndexEntry) SetPerm(perm string){
	i.perm = perm
}

func (i IndexEntry) Name() string{
	return i.entry_name
}
