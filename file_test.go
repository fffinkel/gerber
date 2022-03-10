package main

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"testing"
	"time"
)

//go:embed testdata/20211004.md
var testdata20211004md []byte

//go:embed testdata/20211001.md
var testdata20211001md []byte

//go:embed testdata/20210930.md
var testdata20210930md []byte

//go:embed testdata/20210929.md
var testdata20210929md []byte

//go:embed testdata/20210928.md
var testdata20210928md []byte

//go:embed testdata/20210927.md
var testdata20210927md []byte

//go:embed testdata/20210926.md
var testdata20210926md []byte

//go:embed testdata/20210924.md
var testdata20210924md []byte

//go:embed testdata/20210923.md
var testdata20210923md []byte

//go:embed testdata/20210922.md
var testdata20210922md []byte

func createTestNoteFiles(dir string) error {
	files := map[string][]byte{
		"20211004.md": testdata20211004md,
		"20211001.md": testdata20211001md,
		"20210930.md": testdata20210930md,
		"20210929.md": testdata20210929md,
		"20210928.md": testdata20210928md,
		"20210927.md": testdata20210927md,
		"20210926.md": testdata20210926md,
		"20210924.md": testdata20210924md,
		"20210923.md": testdata20210923md,
		"20210922.md": testdata20210922md,
	}
	for name, contents := range files {
		f, err := os.Create(path.Join(dir, name))
		if err != nil {
			return err
		}
		_, err = f.Write(contents)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

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
