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
		notesPath: "doesnt/matter",
		day: map[string]int{
			"cat1":      4,
			"cat2, one": 8,
			"cat2, two": 12,
		},
		fiveDay: map[string]int{
			"cat1":      4,
			"cat2, one": 16,
			"cat2, two": 24,
		},
		fifteenDay: map[string]int{
			"cat1":      4,
			"cat2, one": 32,
			"cat2, two": 48,
		},
		week: map[string]int{
			"cat1":      3,
			"cat2":      6,
			"cat3, one": 9,
			"cat3, two": 100,
		},
		month: map[string]int{
			"cat1": 10,
			"cat2": 20,
			"cat3": 30,
			"cat4": 40,
		},
		current: "something fun",
	}
}

func TestFiveDayTotal(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	if totals.fiveDayTotal() != 44 {
		t.Error("five day total is incorrect")
	}
}

func TestFifteenDayTotal(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	if totals.fifteenDayTotal() != 84 {
		t.Error("fifteen day total is incorrect")
	}
}

func TestFiveDayThemeTotals(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	fivedaytt := totals.fiveDayThemeTotals()
	if fivedaytt["cat1"] != 4 {
		t.Error("five day theme totals cat1 is incorrect")
	}
	if fivedaytt["cat2"] != 40 {
		t.Error("five day theme totals cat2 is incorrect")
	}
}

func TestFifteenDayThemeTotals(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	fifteendaytt := totals.fifteenDayThemeTotals()
	if fifteendaytt["cat1"] != 4 {
		t.Error("fifteen day theme totals cat1 is incorrect")
	}
	if fifteendaytt["cat2"] != 80 {
		t.Error("fifteen day theme totals cat2 is incorrect")
	}
}

func TestWeekTotal(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	if totals.weekTotal() != 118 {
		t.Error("week total is incorrect")
	}
}

func TestWeekThemeTotals(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	wtt := totals.weekThemeTotals()
	if wtt["cat2"] != 6 {
		t.Error("week theme totals cat2 is incorrect")
	}
	if wtt["cat3"] != 109 {
		t.Error("week theme totals cat3 is incorrect")
	}
}

func TestWeekThemePercent(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	if totals.weekThemePercent("cat2") != 5.084745762711865 {
		t.Error("week theme percent cat2 is incorrect")
	}
	if totals.weekThemePercent("cat3") != 92.37288135593221 {
		t.Error("week theme percent cat3 is incorrect")
	}
}

func TestDayTotal(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	if totals.dayTotal() != 24 {
		t.Error("day total is incorrect")
	}
}

func TestDayThemeTotals(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	wtt := totals.dayThemeTotals()
	if wtt["cat1"] != 4 {
		t.Error("day theme totals cat1 is incorrect")
	}
	if wtt["cat2"] != 20 {
		t.Error("day theme totals cat2 is incorrect")
	}
}

func TestDayThemePercent(t *testing.T) {
	t.Parallel()
	totals := newTestTotals()
	if totals.dayThemePercent("cat1") != 16.666666666666664 {
		t.Error("day theme percent cat1 is incorrect")
	}
	if totals.dayThemePercent("cat2") != 83.33333333333334 {
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
	err := createTestNoteFiles(tempDir)
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
	err := createTestNoteFiles(tempDir)
	if err != nil {
		t.Fatal(err.Error())
	}
	totals := newTotals(tempDir)

	// when the test files start
	then := time.Date(2021, time.October, 5, 12, 0, 0, 0, time.UTC)
	_ = totals.calculate(then)
}
