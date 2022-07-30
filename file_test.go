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

//go:embed testdata/20210916.md
var testdata20210916md []byte

//go:embed testdata/20210915.md
var testdata20210915md []byte

var testNoteFiles = map[string][]byte{
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

func createTestNoteFiles(dir string, files map[string][]byte) error {
	if len(files) == 0 {
		files = testNoteFiles
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

func TestCreateTodayFile(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	if err := createTestNoteFiles(tempDir, nil); err != nil {
		t.Fatal(err.Error())
	}
	if err := createNotesFile(tempDir, "20221216.md"); err != nil {
		t.Fatal(err.Error())
	}
}

func TestCreateTodayFileError(t *testing.T) {
	t.Parallel()
	if err := createNotesFile("/fakedir", "20221216.md"); err == nil {
		t.Error("fake path should cause open error")
	}
}

func TestGetLastNFiles(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	err := createTestNoteFiles(tempDir, nil)
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

func TestGetLastNFilesOrder(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	err := createTestNoteFiles(tempDir, nil)
	if err != nil {
		t.Fatal(err)
	}
	files, err := getLastNFiles(tempDir, 3)
	if files[0].Name() != "20211004.md" {
		t.Error("files in incorrect order")
	}
	if files[1].Name() != "20211001.md" {
		t.Error("files in incorrect order")
	}
	if files[2].Name() != "20210930.md" {
		t.Error("files in incorrect order")
	}
}

func TestGetLastTaskList(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	testFiles := map[string][]byte{
		"20210916.md": testdata20210916md,
		"20210915.md": testdata20210915md,
	}
	err := createTestNoteFiles(tempDir, testFiles)
	if err != nil {
		t.Fatal(err)
	}
	if err = createNotesFile(tempDir, "20221216.md"); err != nil {
		t.Fatal(err.Error())
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetLastNotes(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	err := createTestNoteFiles(tempDir, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	lastNotes, err := getLastNotes(tempDir, "CORE-6229")
	testLastNotes := `


- [ ] test task one
- [ ] test task two
- [x] test task three
- [ ] test task four

Don't forget to make sure the bucket is in the right region.

Still need to make sure the bucket is public.


Some example files that I'm using:
infrastructure/zpan/infra-global/environments/zr-public/s3.tf
infrastructure/terraform/modules/s3-bucket/policies/pdx-tier-access.json.tpl

`
	if lastNotes != testLastNotes {
		t.Error("last notes should match expected last notes")
	}
}

func TestGetLastNotesBadPath(t *testing.T) {
	t.Parallel()
	_, err := getLastNotes("malarky", "CORE-6229")
	if err == nil {
		t.Error("getlastnotes should have errored")
	}
}
