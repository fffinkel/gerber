package main

import (
	"os"
	"path"
	"testing"
)

func createTestNoteFiles(dir string) error {
	fileNames := []string{
		"20211003.md",
		"20211001.md",
		"20210930.md",
		"20210929.md",
		"20210928.md",
		"20210927.md",
		"20210926.md",
		"20210925.md",
		"20210924.md",
	}
	for _, fileName := range fileNames {
		f, err := os.Create(path.Join(dir, fileName))
		if err != nil {
			return err
		}
		_, err = f.WriteString("test")
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func TestGetLastNFiles(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	err := createTestNoteFiles(tempDir)
	if err != nil {
		t.Fatal(err.Error())
	}
	files, err := getLastNFiles(tempDir, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Error("files should have length 2")
	}

	files, err = getLastNFiles(tempDir, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 5 {
		t.Error("files should have length 5")
	}
}
