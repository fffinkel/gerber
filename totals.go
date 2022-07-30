package main

import (
	"bufio"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type totals struct {
	notesPath string
	current   string
	days      []map[string]int

	// TODO remove
	day        map[string]int
	week       map[string]int // TODO remove
	fiveDay    map[string]int
	fifteenDay map[string]int
}

func newTotals(notesPath string) *totals {
	return &totals{
		notesPath: notesPath,
		days:      make([]map[string]int, 15),

		// TODO remove
		day:        make(map[string]int),
		week:       make(map[string]int), // TODO remove
		fiveDay:    make(map[string]int), // TODO ‽‽
		fifteenDay: make(map[string]int),
	}
}

// TODO temporary
func (t *totals) temporaryThing(n int) map[string]int {
	if n == 1 {
		return t.day
	} else if n == 5 {
		return t.fiveDay
	} else if n == 15 {
		return t.fifteenDay
	} else {
		panic("we shouldn't be here!")
	}
}

func (t *totals) weekTotal() int {
	total := 0
	for _, v := range t.week {
		total += v
	}
	return total
}

func (t *totals) nDayTotal(n int) int {
	// TODO temporary
	r := t.temporaryThing(n)
	total := 0
	for _, v := range r {
		total += v
	}
	return total
}

func (t *totals) nDayThemeTotals(n int) map[string]float64 {
	// TODO temporary
	r := t.temporaryThing(n)

	totals := make(map[string]float64)
	for k, _ := range r {
		theme := strings.Split(k, ", ")[0]
		totals[theme] += float64(r[k])
	}
	return totals
}

func (t *totals) nDayThemePercent(n int, theme string) float64 {
	return (t.nDayThemeTotals(n)[theme] / float64(t.nDayTotal(n))) * 100
}

func (t *totals) calculate(today time.Time) error {
	// TODO old
	files, err := ioutil.ReadDir(t.notesPath)
	if err != nil {
		return err
	}

	// TODO new
	filenames, err := getLastNFilenames(t.notesPath, 15)
	if err != nil {
		return err
	}

	var lastDate time.Time
	var lastCategory string
	var lastDateT time.Time
	var category string

	//////////
	//////////
	////////// NEW
	//////////
	//////////

	i := 0
	// for each file
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		t.days[i] = make(map[string]int)

		lineNumber := 0
		scanner := bufio.NewScanner(f)

		// for each line
		for scanner.Scan() {
			line := scanner.Text()

			// if this is a timestamp line
			if strings.HasPrefix(line, "## ") {
				lineNumber++

				// get the category and timestamp
				nextCategory, date, err := parseLine(line)
				if err != nil {
					return err
				}

				// if this is the first line, set it and move on
				if lineNumber == 1 {
					category = nextCategory
					lastDateT = date
					continue
				}

				// if the category is not "break," add its minutes
				if category != "break" {
					minutes := date.Sub(lastDateT).Minutes()
					t.days[i][category] += int(minutes)
				}

				lastDateT = date
				category = nextCategory
			}
		}

		// if we end with a currently open task, add its minutes
		if category != "break" {
			minutes := int(today.Sub(lastDateT).Minutes())
			t.days[i][category] += minutes
			t.current = category
		}

		// TODO this needs a defer
		f.Close()
		i++

	} // for each file

	//////////
	//////////
	////////// OLD
	//////////
	//////////

	// TODO broken
	for _, file := range files {
		if shouldSkip(file) {
			continue
		}
		f, err := os.Open(filepath.Join(t.notesPath, file.Name()))
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNumber := 0
		for scanner.Scan() {

			line := scanner.Text()
			lineNumber += 1

			if strings.HasPrefix(line, "## ") {
				category, date, err := parseLine(line)
				if err != nil {
					// TODO wrap this error with file name (file, lineNumber)
					return err
				}

				// if this category's day doesn't match the previous line's day,
				// reset the day count
				if date.Day() != lastDate.Day() {
					t.day = make(map[string]int)
				}

				// TODO reset the five and fifteen day count
				// if this category's day doesn't match the previous line's day,
				// reset the day count
				// if date.Day()-5 >= lastDate.Day() {
				// 	t.fiveDay = make(map[string]int)
				// }

				// if this category's week doesn't match the previous line's week,
				// reset the week count
				if weekNumber(date) != weekNumber(lastDate) {
					t.week = make(map[string]int)
				}

				minutes := date.Sub(lastDate).Minutes()

				// TODO explain the check for "break" here
				if lastCategory != "break" {

					// TODO this should count only working days
					if today.Sub(date).Hours()/24 < 5 {
						t.fiveDay[lastCategory] += int(minutes)
					}

					// TODO this should count only working days
					if today.Sub(date).Hours()/24 < 15 {
						t.fifteenDay[lastCategory] += int(minutes)
					}

					if weekNumber(date) == weekNumber(today) {
						t.week[lastCategory] += int(minutes)

						if date.Day() == time.Now().Day() {
							t.day[lastCategory] += int(minutes)
						}
					}
				}
				lastDate = date
				lastCategory = category
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		// include the current task
		if lastCategory != "break" && lastCategory != "" {

			// this should be
			// if lastDate.Day() != today.Day() {
			// 	// TODO this error not handled
			// 	return errors.New("file not closed properly")
			// }

			t.current = lastCategory
			minutes := int(today.Sub(lastDate).Minutes())
			t.day[lastCategory] += minutes
			t.week[lastCategory] += minutes
		}
	}

	return nil
}

func shouldSkip(file fs.FileInfo) bool {
	if file.Name() == ".git" {
		return true
	}
	if file.Name() == "2020" {
		return true
	}
	if file.Name() == "2021" {
		return true
	}
	return false
}
