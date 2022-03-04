package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetTodayFilename(t *testing.T) {
	t.Parallel()
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	expected := fmt.Sprintf("%04d%02d%02d.md", year, month, day)
	got := getTodayFilename()
	if got != expected {
		t.Errorf("got incorrect today filename: %s", got)
	}
}

func TestGetLastFilename(t *testing.T) {
	t.Parallel()
}

// func TestGetLastNFiles(t *testing.T) {
// 	t.Parallel()
// 	tempDir := t.TempDir()

// 	err := createTestNoteFiles(tempDir)
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	files, err := getLastNFiles(tempDir, 2)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(files) != 2 {
// 		t.Error("files should have length 2")
// 	}

// 	files, err = getLastNFiles(tempDir, 5)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(files) != 5 {
// 		t.Error("files should have length 5")
// 	}
// }
