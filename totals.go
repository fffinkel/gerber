package main

import (
	"fmt"
	"sort"
	"strings"
)

func (t *totals) printSummary() error {
	fmt.Print("\n")
	fmt.Printf("This week you have worked: %+v\n", minToHourMin(t.weekTotal()))
	fmt.Printf("Today you have worked: %+v\n", minToHourMin(t.dayTotal()))
	fmt.Printf("You are currently working on: \"%+v\"\n", t.current)
	fmt.Print("\n")

	lastTheme := ""
	for _, category := range sortedKeys(t.day) {
		theme := strings.Split(category, ", ")[0]
		if theme != lastTheme {
			fmt.Printf("-- %s (%.1f%%d, %.1f%%w) --\n", theme,
				t.dayThemePercent(theme),
				t.weekThemePercent(theme))
			lastTheme = theme
		}
		fmt.Printf("%s: %s\n", category, minToHourMin(t.day[category]))
	}
	return nil
}

func (t *totals) weekTotal() int {
	total := 0
	for _, v := range t.week {
		total += v
	}
	return total
}

func (t *totals) weekThemeTotals() map[string]float64 {
	totals := make(map[string]float64)
	for k, _ := range t.week {
		theme := strings.Split(k, ", ")[0]
		totals[theme] += float64(t.week[k])
	}
	return totals
}

func (t *totals) weekThemePercent(theme string) float64 {
	return (t.weekThemeTotals()[theme] / float64(t.weekTotal())) * 100
}

func (t *totals) dayTotal() int {
	total := 0
	for _, v := range t.day {
		total += v
	}
	return total
}

func (t *totals) dayThemeTotals() map[string]float64 {
	totals := make(map[string]float64)
	for k, _ := range t.day {
		theme := strings.Split(k, ", ")[0]
		totals[theme] += float64(t.day[k])
	}
	return totals
}

func (t *totals) dayThemePercent(theme string) float64 {
	return (t.dayThemeTotals()[theme] / float64(t.dayTotal())) * 100
}

func minToHourMin(m int) string {
	minutes := m % 60
	hours := m / 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func sortedKeys(in map[string]int) []string {
	sorted := make([]string, 0, len(in))
	for k := range in {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	return sorted
}
