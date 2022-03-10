package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type totals struct {
	notesPath string
	day       map[string]int
	week      map[string]int
	month     map[string]int
	current   string
}

func newTotals(notesPath string) *totals {
	return &totals{
		notesPath: notesPath,
		day:       make(map[string]int),
		week:      make(map[string]int),
		month:     make(map[string]int),
	}
}

func (t *totals) weekTotal() int {
	total := 0
	for _, v := range t.week {
		total += v
	}
	return total
}

func (t *totals) weekThemeTotals() map[string]float64 {
	totals := make(map[string]float64)
	for k, _ := range t.week {
		theme := strings.Split(k, ", ")[0]
		totals[theme] += float64(t.week[k])
	}
	return totals
}

func (t *totals) weekThemePercent(theme string) float64 {
	return (t.weekThemeTotals()[theme] / float64(t.weekTotal())) * 100
}

func (t *totals) dayTotal() int {
	total := 0
	for _, v := range t.day {
		total += v
	}
	return total
}

func (t *totals) dayThemeTotals() map[string]float64 {
	totals := make(map[string]float64)
	for k, _ := range t.day {
		theme := strings.Split(k, ", ")[0]
		totals[theme] += float64(t.day[k])
	}
	return totals
}

func (t *totals) dayThemePercent(theme string) float64 {
	return (t.dayThemeTotals()[theme] / float64(t.dayTotal())) * 100
}

func (t *totals) calculate(today time.Time) error {
	files, err := ioutil.ReadDir(t.notesPath)
	if err != nil {
		return err
	}

	var lastDate time.Time
	var lastCategory string

	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}
		if file.Name() == "2020" {
			continue
		}
		if file.Name() == "2021" {
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

				// if we're not in the current month, don't do anything
				// if date.Month() != today.Month() {
				if int(today.Sub(date).Minutes())/60/24 > 31 {
					continue
				}

				// if this category's day doesn't match the previous line's day,
				// reset the day count
				if date.Day() != lastDate.Day() {
					t.day = make(map[string]int)
				}

				// if this category's week doesn't match the previous line's week,
				// reset the week count
				if weekNumber(date) != weekNumber(lastDate) {
					t.week = make(map[string]int)
				}

				// if this category's month doesn't match the previous line's month,
				// reset the month count
				if date.Month() != today.Month() {
					t.month = make(map[string]int)
				}

				minutes := date.Sub(lastDate).Minutes()

				// TODO explain the check for "break" here
				if date.Month() == lastDate.Month() && lastCategory != "break" {
					t.month[lastCategory] += int(minutes)

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

		// TODO why doesn't this include month
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
