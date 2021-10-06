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

	return []byte(fmt.Sprintf("# %s %02d, %04d (q‽‽‽‽‽‽‽‽‽‽s‽‽‽‽‽‽‽‽‽‽)\n\nToday is a beautiful day.%s\n\n## %04d-%02d-%02d %s %s %s (admin)\n\n\n",
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

// TODO this shouldn't be a struct method
func (t *totals) printSummary() error {
	fmt.Print("\n")
	fmt.Printf("This week you have worked: %+v\n",
		minToHourMin(t.weekTotal()))
	fmt.Printf("Today you have worked: %+v\n",
		minToHourMin(t.dayTotal()))

	fmt.Print("\n")

	themeTotals := make(map[string]int)
	for _, category := range sortedKeys(t.day) {
		if category == "sprint" {
			continue
		}
		theme := strings.Split(category, ", ")[0]
		themeTotals[theme] += t.day[category]
	}

	for _, theme := range sortedKeys(themeTotals) {
		fmt.Printf(" ➔ %s: %s (%.1f%%d, %.1f%%w)\n",
			theme,
			minToHourMin(themeTotals[theme]),
			t.dayThemePercent(theme),
			t.weekThemePercent(theme))
	}

	// fmt.Printf("Current sprint percentages are: %.1f%%d, %.1f%%w\n",
	// 	t.daySprintPercent(), t.weekSprintPercent())

	fmt.Print("\n")

	if t.current != "" {
		fmt.Printf("You are currently working on: \"%+v\"\n", t.current)
	} else {
		fmt.Print("\n** You are not currently tracking any work **\n")
	}
	return nil
}
