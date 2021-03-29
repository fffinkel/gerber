package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
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

type commandFunc func() error

var commandMap = map[string]commandFunc{

	// TODO delete
	"print_today_work": printTodayWork,

	"print_today_filename": printTodayFilename,
	"print_summary":        printSummary,

	// "new_file":             createTodayFile,
}

type totals struct {
	day   map[string]int
	week  map[string]int
	month map[string]int
}

func newTotals() *totals {
	return &totals{
		day:   make(map[string]int),
		week:  make(map[string]int),
		month: make(map[string]int),
	}
}

// func (t *totals) getFormattedDayTotal() string

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("command required"))
	}
	commandArg := os.Args[1]

	if commandFunc, ok := commandMap[commandArg]; ok {
		err := commandFunc()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(errors.New(
			fmt.Sprintf("command [%s] not found", commandArg)),
		)
	}
}

func weekNumber(t time.Time) int {
	_, weekNum := t.ISOWeek()
	return weekNum
}

func (t *totals) calculate(filePaths []os.FileInfo) error {
	var lastDate time.Time
	var lastCategory string
	now := time.Now()

	for _, file := range filePaths {
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
				if date.Month() != now.Month() {
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

				minutes := date.Sub(lastDate).Minutes()
				if date.Month() == lastDate.Month() && lastCategory != "break" {
					t.month[lastCategory] += int(minutes)

					if weekNumber(date) == weekNumber(now) {
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
			//currentCategory = lastCategory
			minutes := int(now.Sub(lastDate).Minutes())
			t.day[lastCategory] += minutes
		}
	}

	zz, _ := json.MarshalIndent(t.day, "", "\t")
	fmt.Printf("\n\n----------> %s\n", zz)

	return nil
}
