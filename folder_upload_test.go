package main

import (
	"path/filepath"
	"testing"
)

func TestGetFilePaths(t *testing.T) {
	paths := getFilePaths("./")
	if len(paths) == 0 {
		t.Error("Unable to get the filepaths")
	}
}

func TestFilepathBase(t *testing.T) {
	file_path := "/Users/alastairbrunton"
	cleaned_base := filepath.Base(file_path)
	if cleaned_base != "alastairbrunton" {
		t.Error("Unable to get the correct basepath")
	}
}
