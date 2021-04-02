package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func parseLine(line string, file os.FileInfo, lineNumber int) (string, time.Time) {
	trimmed := strings.TrimPrefix(line, "## ")
	lineParts := strings.Split(trimmed, " (")

	if len(lineParts) != 2 {
		panic(fmt.Sprintf("found malformed, %s line %d: %s\n", file.Name(), lineNumber, line)) // TODO no panics
	}

	date := lineParts[0]
	parsedDate, err := time.Parse("2006-01-02 Mon 15:04 PM MST", date)
	if err != nil {
		panic(fmt.Sprintf("error parsing date, %s line %d: %s\n", file.Name(), lineNumber, date)) // TODO no panics
	}

	return strings.TrimSuffix(lineParts[1], ")"), parsedDate

	// splitCategory := strings.Split(strings.TrimSuffix(lineParts[1], ")"), ", ")
	// if len(splitCategory) < 2 {
	// 	return splitCategory[0], "", parsedDate
	// }
	// return splitCategory[0], splitCategory[1], parsedDate
}

func printTodayFilename() error {
	fmt.Printf(filepath.Join(notesPath, getTodayFilename()))
	return nil
}

func getTodayFilename() string {
	today := time.Now()
	return fmt.Sprintf(
		"%d%02d%02d%s",
		today.Year(),
		today.Month(),
		today.Day(),
		notesExtension,
	)
}