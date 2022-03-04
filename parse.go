package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

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
