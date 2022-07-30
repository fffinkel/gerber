package main

import (
	"fmt"
	"testing"
	"time"
)

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
This week you have worked: 1h 58m
Today you have worked: 0h 24m

 ➔ cat1: 0h 4m (16.7%, 9.1%, 4.8%)
 ➔ cat2: 0h 20m (83.3%, 90.9%, 95.2%)

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
This week you have worked: 1h 58m
Today you have worked: 0h 24m

 ➔ cat1: 0h 4m (16.7%, 9.1%, 4.8%)
 ➔ cat2: 0h 20m (83.3%, 90.9%, 95.2%)

You are not currently tracking any work.
`
	totals := newTestTotals()
	totals.current = ""
	summary := getSummary(&totals)
	if string(summary) != testSummary {
		t.Error("summary should match expected summary")
	}
}
