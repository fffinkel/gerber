package main

import "fmt"

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

func weekNumber(t time.Time) int {
	_, weekNum := t.ISOWeek()
	return weekNum
}
