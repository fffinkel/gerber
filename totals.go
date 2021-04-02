package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type totals struct {
	day     map[string]int
	week    map[string]int
	month   map[string]int
	current string
}

func newTotals() *totals {
	return &totals{
		day:   make(map[string]int),
		week:  make(map[string]int),
		month: make(map[string]int),
	}
}

func (t *totals) printSummary() error {
	fmt.Print("\n")
	fmt.Printf("This week you have worked: %+v\n", minToHourMin(t.weekTotal()))
	fmt.Printf("Today you have worked: %+v\n", minToHourMin(t.dayTotal()))
	fmt.Printf("You are currently working on: \"%+v\"\n", t.current)
	fmt.Print("\n")

	lastTheme := ""
	for _, category := range sortedKeys(t.day) {
		theme := strings.Split(category, ", ")[0]
		if theme != lastTheme {
			fmt.Printf("-- %s (%.1f%%d, %.1f%%w) --\n", theme,
				t.dayThemePercent(theme),
				t.weekThemePercent(theme))
			lastTheme = theme
		}
		fmt.Printf("%s: %s\n", category, minToHourMin(t.day[category]))
	}
	return nil
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

func (t *totals) calculate() error {
	files, err := ioutil.ReadDir(notesPath)
	if err != nil {
		return err
	}

	var lastDate time.Time
	var lastCategory string
	today := time.Now()

	for _, file := range files {
		f, err := os.Open(filepath.Join(notesPath, file.Name()))
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
				category, date := parseLine(line, file, lineNumber)

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

		// include the current task
		if lastCategory != "break" && lastCategory != "" {
			if lastDate.Day() != today.Day() {
				// TODO this error not handled
				return errors.New("file not closed properly")
			}

			t.current = lastCategory
			minutes := int(today.Sub(lastDate).Minutes())
			t.day[lastCategory] += minutes
		}
	}

	return nil
}