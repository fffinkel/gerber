package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	notesPath      = "/Users/mattfinkel/.notes/"
	notesExtension = ".md"
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

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("command required"))
	}
	switch os.Args[1] {
	case "summary":
		t := newTotals()
		t.calculate()
		t.printSummary()
		return
	default:
		log.Fatal(errors.New(
			fmt.Sprintf("command [%s] not found", os.Args[1])),
		)
	}
}

func weekNumber(t time.Time) int {
	_, weekNum := t.ISOWeek()
	return weekNum
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
				//if date.Month() != today.Month() {
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
