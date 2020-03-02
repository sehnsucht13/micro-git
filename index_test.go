package main

import (
	"reflect"
	"sort"
	"testing"
)

func Test_SortEmpty(t *testing.T) {
	t.Logf("Empty entries sorted")
	input_entries := []IndexEntry{}
	sort.Sort(ByEntry(input_entries))
}

func Test_Sort1(t *testing.T) {
	input_entries := []IndexEntry{IndexEntry{"123", 0, "ab33", "file_name"}}
	output_entries := []IndexEntry{IndexEntry{"123", 0, "ab33", "file_name"}}
	sort.Sort(ByEntry(input_entries))
	for idx, entry := range output_entries {
		if !reflect.DeepEqual(entry, input_entries[idx]) {
			t.Fail()
			t.Logf("Expected %s, got %s", entry.String(), input_entries[idx].String())
		}
	}
}

func Test_SortMultiple(t *testing.T) {
	input_entries := []IndexEntry{
		IndexEntry{"123", 0, "ab33", "file_name3"},
		IndexEntry{"121", 0, "ab31", "file_name1"},
		IndexEntry{"122", 0, "ab32", "file_name2"},
	}
	output_entries := []IndexEntry{
		IndexEntry{"121", 0, "ab31", "file_name1"},
		IndexEntry{"122", 0, "ab32", "file_name2"},
		IndexEntry{"123", 0, "ab33", "file_name3"},
	}
	sort.Sort(ByEntry(input_entries))
	for idx, entry := range output_entries {
		if !reflect.DeepEqual(entry, input_entries[idx]) {
			t.Fail()
			t.Logf("Expected %s, got %s", entry.String(), input_entries[idx].String())
		}
	}
}

func Test_FindEntryEmpty(t *testing.T) {
	input_entries := []IndexEntry{}
	idx, hash, err := findIndexEntry("file_name_missing", input_entries)
	if err == nil || idx != -1 || hash != "" {
		t.Fail()
		t.Logf("Expected %d,%s, got %d, %s,", -1, "", idx, hash)
	}
}

func Test_FindEntryMissing(t *testing.T) {
	input_entries := []IndexEntry{
		IndexEntry{"121", 0, "ab31", "file_name1"},
		IndexEntry{"122", 0, "ab32", "file_name2"},
		IndexEntry{"123", 0, "ab33", "file_name3"},
	}
	idx, hash, err := findIndexEntry("file_name_missing", input_entries)
	if err == nil || idx != -1 || hash != "" {
		t.Fail()
		t.Logf("Expected %d,%s, got %d, %s,", -1, "", idx, hash)
	}
}

func Test_FindEntryPresent(t *testing.T) {
	input_entries := []IndexEntry{
		IndexEntry{"121", 0, "ab31", "file_name1"},
		IndexEntry{"122", 0, "ab32", "file_name2"},
		IndexEntry{"123", 0, "ab33", "file_name3"},
	}
	idx, hash, err := findIndexEntry("file_name2", input_entries)
	if idx != 1 || hash != "ab32" || err != nil {
		t.Fail()
		t.Logf("Expected %d,%s, got %d, %s,", -1, "", idx, hash)
	}
}
