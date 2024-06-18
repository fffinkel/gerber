package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

var allowedThemes = map[string]string{
	"extracurricular": "[todoist] extracurricular work",
	"project":         "[todoist] project work", // productive
	"ten percent":     "[todoist] stuff I want to do on my own time",

	"help":  "[unplanned] when people ask for help",
	"pager": "[unplanned] when the pager goes off",

	"admin":   "getting ready for my day, repeated tasks, email",
	"coffee":  "socializing around the office, Core Social meeting",
	"meeting": "scheduled meetings, unrelated to projects or extracurricular",
	"ops":     "k8s work, dev ops, net ops", // productive
	"code":    "writing code",               // productive
	"toil":    "manual work, generally deprovisions",
}

var productiveThemes = []string{"project", "ops", "code"}

// var allowedThemes = []string{

// 	// Todoist
// 	"extracurricular", // extracurricular work (Todoist project)
// 	"project",         // project work (Todoist project)
// 	"ten percent",     // stuff I want to do on my own time (Todoist project)

// 	// Other
// 	"admin",   // getting ready for my day, repeated tasks, email
// 	"coffee",  // socializing around the office, Core Social meeting
// 	"meeting", // scheduled meetings, unrelated to projects or extracurricular
// 	"ops",     // k8s work, dev ops, net ops
// 	"pager",   // unplanned requests, help, pages
// 	"toil",    // manual work, generally deprovisions

// 	// "code", removed in favor of "project"
// 	// "design", removed in favor of "project"
// 	// "help", removed in favor of "pager"
// 	// "k8s", removed in favor of "ops"
// }

func isAllowedTheme(theme string) bool {
	if _, ok := allowedThemes[theme]; ok {
		return true
	}
	return false
}

func isProductiveTheme(theme string) bool {
	for _, productive := range productiveThemes {
		if theme == productive {
			return true
		}
	}
	return false
}

func allowedThemesErrorString(theme string) string {
	errorString := fmt.Sprintf("error parsing data: theme '%s' is not allowed\n\nallowed themes are:\n", theme)
	var sortedAllowedThemes []string
	for theme, _ := range allowedThemes {
		sortedAllowedThemes = append(sortedAllowedThemes, theme)
	}
	sort.Strings(sortedAllowedThemes)
	for _, theme := range sortedAllowedThemes {
		errorString = fmt.Sprintf("%s\n ➔ %s: %s", errorString, theme, allowedThemes[theme])
	}
	return errorString
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
	theme := strings.Split(category, ", ")[0]
	if !isAllowedTheme(theme) && theme != "break" {
		return "", time.Time{}, errors.New(allowedThemesErrorString(theme))
	}
	return category, parsedDate, nil
}
