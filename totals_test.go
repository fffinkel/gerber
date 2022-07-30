package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewTotals(t *testing.T) {
	t.Parallel()
	totals := newTotals("blah")

	if totals.notesPath != "blah" {
		t.Error("notes path is incorrect")
	}
}

func newTestTotals() totals {
	return totals{
		mondayIndex: 4,
		notesPath:   "doesnt/matter",
		current:     "something fun",
		days: []map[string]int{
			map[string]int{
				"cat1":      4,
				"cat2, one": 8,
				"cat2, two": 12,
			},
			map[string]int{
				"cat1":      4,
				"cat2, one": 16,
				"cat2, two": 24,
			},
			map[string]int{
				"cat1":      4,
				"cat4, one": 32,
				"cat3, two": 48,
			},
			map[string]int{
				"cat1":      4,
				"cat8, one": 8,
				"cat0, two": 12,
			},
			map[string]int{
				"cat9":      4,
				"cat9, one": 8,
				"cat9, two": 12,
			},
			map[string]int{
				"cat1":      4,
				"cat2, one": 8,
				"cat2, two": 12,
			},
			map[string]int{
				"cat1":      4,
				"cat8, one": 8,
				"cat0, two": 12,
			},
			map[string]int{
				"cat9":      4,
				"cat9, one": 8,
				"cat9, two": 12,
			},
			map[string]int{
				"cat1":      4,
				"cat2, one": 8,
				"cat2, two": 12,
			},
			map[string]int{
				"cat1":      4,
				"cat8, one": 8,
				"cat0, two": 12,
			},
			map[string]int{
				"cat9":      4,
				"cat9, one": 8,
				"cat9, two": 12,
			},
			map[string]int{
				"cat1":      4,
				"cat2, one": 8,
				"cat2, two": 12,
			},
			map[string]int{
				"cat1":      4,
				"cat8, one": 8,
				"cat0, two": 12,
			},
			map[string]int{
				"cat9":      4,
				"cat9, one": 8,
				"cat9, two": 12,
			},
			map[string]int{
				"cat1":      4,
				"cat2, one": 8,
				"cat2, two": 12,
			},
		},
	}
}

func TestOneDayTotal(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	if totals.nDayTotalMinutes(1) != 24 {
		t.Error("day total is incorrect")
	}
}

func TestFiveDayTotal(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	if totals.nDayTotalMinutes(5) != 200 {
		t.Error("five day total is incorrect")
	}
}

func TestFifteenDayTotal(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	if totals.nDayTotalMinutes(15) != 440 {
		t.Error("fifteen day total is incorrect")
	}
}

func TestOneDayThemeTotalMinutes(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	wtt := totals.nDayThemeTotalMinutes(1)
	if wtt["cat1"] != 4 {
		t.Error("day theme totals cat1 is incorrect")
	}
	if wtt["cat2"] != 20 {
		t.Error("day theme totals cat2 is incorrect")
	}
}

func TestFiveDayThemeTotalMinutes(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	fivedaytt := totals.nDayThemeTotalMinutes(5)
	if fivedaytt["cat1"] != 16 {
		t.Error("five day theme totals cat1 is incorrect")
	}
	if fivedaytt["cat2"] != 60 {
		t.Error("five day theme totals cat2 is incorrect")
	}
	if fivedaytt["cat9"] != 24 {
		t.Error("five day theme totals cat9 is incorrect")
	}
}

func TestFifteenDayThemeTotalMinutes(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	fifteendaytt := totals.nDayThemeTotalMinutes(15)
	if fifteendaytt["cat1"] != 44 {
		t.Error("fifteen day theme totals cat1 is incorrect")
	}
	if fifteendaytt["cat2"] != 140 {
		t.Error("fifteen day theme totals cat2 is incorrect")
	}
	if fifteendaytt["cat9"] != 96 {
		t.Error("fifteen day theme totals cat9 is incorrect")
	}
}

func TestWeekTotalMinutes(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	if totals.weekTotalMinutes() != 200 {
		t.Error("week total is incorrect")
	}
}

func TestOneDayThemePercent(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	if totals.nDayThemePercent(1, "cat1") != 16.666666666666664 {
		t.Error("day theme percent cat1 is incorrect")
	}
	if totals.nDayThemePercent(1, "cat2") != 83.33333333333334 {
		t.Error("day theme percent cat2 is incorrect")
	}
}

func TestCalculatePathError(t *testing.T) {
	t.Parallel()
	totals := newTotals("does/not/exist")
	if err := totals.calculate(time.Now()); err == nil {
		t.Error("calculate should have errored")
	}
}

func TestCalculateSkipDirs(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	err := createTestNoteFiles(tempDir, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	totals := newTotals(tempDir)

	// TODO probably make sure actual files in
	// these dirs are not parsed?
	os.Mkdir(filepath.Join(tempDir, ".git"), 0777)
	os.Mkdir(filepath.Join(tempDir, "2020"), 0777)
	os.Mkdir(filepath.Join(tempDir, "2021"), 0777)
	_ = totals.calculate(time.Now())
}

func TestCalculate(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	err := createTestNoteFiles(tempDir, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	totals := newTotals(tempDir)

	// when the test files start
	then := time.Date(2021, time.October, 5, 12, 0, 0, 0, time.UTC)
	_ = totals.calculate(then)
}
