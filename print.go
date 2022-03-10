package main

import (
	"fmt"
	"strings"
	"time"
)

func getNotesHeader(path string) ([]byte, error) {
	today := time.Now()
	weekday := string([]byte(today.Weekday().String())[:3])
	kitchen := today.Format("03:04 PM")
	zone, _ := today.Zone()

	taskList, err := getLastTaskList(path)
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("# %s %02d, %04d\n\nToday is a beautiful day.%s\n\n## %04d-%02d-%02d %s %s %s (admin)\n\n\n",
		today.Month(),
		today.Day(),
		today.Year(),
		taskList,
		today.Year(),
		today.Month(),
		today.Day(),
		weekday,
		kitchen,
		zone,
	)), nil
}

// TODO make this not bullshit
func getSummary(t *totals) string {
	summary := "\n"
	summary += fmt.Sprintf("This week you have worked: %+v\n",
		minToHourMin(t.weekTotal()))
	summary += fmt.Sprintf("Today you have worked: %+v\n",
		minToHourMin(t.dayTotal()))

	summary += "\n"

	themeTotals := make(map[string]int)
	for _, category := range sortedKeys(t.day) {
		theme := strings.Split(category, ", ")[0]
		themeTotals[theme] += t.day[category]
	}

	for _, theme := range sortedKeys(themeTotals) {
		summary += fmt.Sprintf(" ➔ %s: %s (%.1f%%d, %.1f%%w)\n",
			theme,
			minToHourMin(themeTotals[theme]),
			t.dayThemePercent(theme),
			t.weekThemePercent(theme))
	}

	if len(themeTotals) > 0 {
		summary += "\n"
	}

	if t.current != "" {
		summary += fmt.Sprintf("You are currently working on: %+v\n", t.current)
	} else {
		summary += "You are not currently tracking any work.\n"
	}
	return summary
}
