package main

import (
	"fmt"
	"testing"
	"time"
)

func init() {
	allowedThemes = []string{
		"cat1",
		"cat2",
		"cat3",
		"cat4",
		"cat5",
		"cat6",
		"cat7",
		"cat8",
		"cat9",
		"cat0",
	}
}

func TestGetNotesHeader(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	err := createTestNoteFiles(tempDir, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	header, err := getNotesHeader(tempDir)
	if err != nil {
		t.Fatal(err.Error())
	}

	today := time.Now()
	weekday := string([]byte(today.Weekday().String())[:3])
	kitchen := today.Format("03:04 PM")
	zone, _ := today.Zone()

	taskList, err := getLastTaskList(tempDir)
	if err != nil {
		t.Fatal(err.Error())
	}

	// this is pretty dumb because it's just
	// re-implementing the code its testing, but
	// oh well fuck it coverage or death
	testHeader := fmt.Sprintf(`# %s %02d, %04d

Today is a beautiful day.%s

## %04d-%02d-%02d %s %s %s (admin)


`,
		today.Month(),
		today.Day(),
		today.Year(),
		taskList,
		today.Year(),
		today.Month(),
		today.Day(),
		weekday,
		kitchen,
		zone)
	if string(header) != testHeader {
		t.Error("summary should match expected summary")
	}
}

func TestGetNotesHeaderError(t *testing.T) {
	t.Parallel()
	if _, err := getNotesHeader("does/not/exist"); err == nil {
		t.Error("calculate should have errored")
	}
}

func TestGetSummary(t *testing.T) {
	t.Parallel()
	testSummary := `
This week you have worked: 3h 20m
Today you have worked: 0h 24m

 ➔ cat2: 0h 20m (83.3%, 30.0%, 31.8%)
 ➔ cat9: 0h 0m (0.0%, 12.0%, 21.8%)
 ➔ cat0: 0h 0m (0.0%, 6.0%, 10.9%)
 ➔ cat3: 0h 0m (0.0%, 24.0%, 10.9%)
 ➔ cat1: 0h 4m (16.7%, 8.0%, 10.0%)
 ➔ cat4: 0h 0m (0.0%, 16.0%, 7.3%)
 ➔ cat8: 0h 0m (0.0%, 4.0%, 7.3%)

You are currently working on: something fun
`
	totals := newTestTotals()
	summary := getSummary(&totals)
	if string(summary) != testSummary {
		t.Error("summary should match expected summary")
	}
}

func TestGetSummaryNoCurrent(t *testing.T) {
	t.Parallel()
	testSummary := `
This week you have worked: 3h 20m
Today you have worked: 0h 24m

 ➔ cat2: 0h 20m (83.3%, 30.0%, 31.8%)
 ➔ cat9: 0h 0m (0.0%, 12.0%, 21.8%)
 ➔ cat0: 0h 0m (0.0%, 6.0%, 10.9%)
 ➔ cat3: 0h 0m (0.0%, 24.0%, 10.9%)
 ➔ cat1: 0h 4m (16.7%, 8.0%, 10.0%)
 ➔ cat4: 0h 0m (0.0%, 16.0%, 7.3%)
 ➔ cat8: 0h 0m (0.0%, 4.0%, 7.3%)

You are not currently tracking any work.
`
	totals := newTestTotals()
	totals.current = ""
	summary := getSummary(&totals)
	if string(summary) != testSummary {
		t.Error("summary should match expected summary")
	}
}
