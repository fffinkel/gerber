package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type parsedNote map[string]int

type parser struct {
	notesDir    string
	parsedNotes [15]parsedNote
	current     string
}

// TODO should this use file/path?
func newParser(notesDir string) *parser {
	return &parser{notesDir: notesDir}
}

func parseLine(line string) (string, time.Time, error) {
	trimmed := strings.TrimPrefix(line, "## ")
	lineParts := strings.Split(trimmed, " (")

	if len(lineParts) != 2 {
		return "", time.Time{}, errors.New(fmt.Sprintf("found malformed line: %s\n", line))
	}

	date := lineParts[0]
	parsedDate, err := time.Parse("2006-01-02 Mon 15:04 PM MST", date)
	if err != nil {
		return "", time.Time{}, errors.New(fmt.Sprintf("error parsing date: %s\n", date))
	}

	category := strings.TrimSuffix(lineParts[1], ")")
	return category, parsedDate, nil
}

// type totals struct {
// 	notesPath  string
// 	day        map[string]int
// 	week       map[string]int // TODO remove
// 	fiveDay    map[string]int
// 	fifteenDay map[string]int
// 	current    string
// }

// func newTotals(notesPath string) *totals {
// 	return &totals{
// 		notesPath:  notesPath,
// 		day:        make(map[string]int),
// 		week:       make(map[string]int), // TODO remove
// 		fiveDay:    make(map[string]int), // TODO ‽‽
// 		fifteenDay: make(map[string]int),
// 	}
// }
