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

func getSummary(t *totals) string {
	return fmt.Sprintf(`
This week you have worked: %s
Today you have worked: %s

%s
%s
`,
		minToHourMin(t.weekTotalMinutes()),
		minToHourMin(t.nDayTotalMinutes(1)),
		t.summaryCategories(),
		t.summaryFooter())
}

func (t *totals) summaryCategories() string {
	themeTotals := make(map[string]int)
	for _, category := range sortedKeys(t.days[0]) {
		theme := strings.Split(category, ", ")[0]
		themeTotals[theme] += t.days[0][category]
	}
	categories := ""
	for _, theme := range sortedKeys(themeTotals) {
		categories += fmt.Sprintf(" ➔ %s: %s (%.1f%%, %.1f%%, %.1f%%)\n",
			theme,
			minToHourMin(themeTotals[theme]),
			t.nDayThemePercent(1, theme),
			t.nDayThemePercent(5, theme),
			t.nDayThemePercent(15, theme))
	}
	return categories
}

func (t *totals) summaryFooter() string {
	if t.current != "" {
		return fmt.Sprintf("You are currently working on: %s", t.current)
	} else {
		return "You are not currently tracking any work."
	}
}
