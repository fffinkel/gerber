package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

func parseLine(line string, file os.FileInfo, lineNumber int) (string, time.Time, bool, error) {
	trimmed := strings.TrimPrefix(line, "## ")
	lineParts := strings.Split(trimmed, " (")

	if len(lineParts) != 2 {
		return "", time.Time{}, false, errors.New(fmt.Sprintf("found malformed, %s line %d: %s\n", file.Name(), lineNumber, line))
	}

	date := lineParts[0]
	parsedDate, err := time.Parse("2006-01-02 Mon 15:04 PM MST", date)
	if err != nil {
		return "", time.Time{}, false, errors.New(fmt.Sprintf("error parsing date, %s line %d: %s\n", file.Name(), lineNumber, date))
	}

	isSprint := false
	category := strings.TrimSuffix(lineParts[1], ")")
	categoryTrimmed := strings.TrimPrefix(category, "*")
	if categoryTrimmed != category {
		isSprint = true
	}

	return categoryTrimmed, parsedDate, isSprint, nil

	// splitCategory := strings.Split(strings.TrimSuffix(lineParts[1], ")"), ", ")
	// if len(splitCategory) < 2 {
	// 	return splitCategory[0], "", parsedDate
	// }
	// return splitCategory[0], splitCategory[1], parsedDate
}
