package main

import (
	"fmt"
	"io/ioutil"
	"sort"
)

// TODO delete
func printTodayWork() error {
	return printSummary()
}

func printSummary() error {
	files, err := ioutil.ReadDir(notesPath)
	if err != nil {
		return err
	}

	t := newTotals()
	if err := t.calculate(files); err != nil {
		return err
	}

	// printTotals(dayTotals, current, false, true)
	// fmt.Printf("\n")
	// printTotals(daySuperTotals, current, true, false)
	return nil
}

func minToHourMin(m int) string {
	minutes := m % 60
	hours := m / 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func printTotals(totals map[string]int, current string, printPercents, printTotals bool) {
	totalsString := ""
	totalsInt := 0

	sorted := make([]string, 0, len(totals))
	for k := range totals {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	for _, category := range sorted {
		if category != "break" {
			totalsInt += totals[category]
		}
	}

	for _, category := range sorted {
		if category != "break" {
			hm := minToHourMin(totals[category])
			inProgress := ""
			if category == current && printTotals {
				inProgress = "--> "
			}
			totalsString = totalsString + fmt.Sprintf("%s%s: %s", inProgress, category, hm)

			if printPercents {
				percent := (float64(totals[category]) / float64(totalsInt)) * 100
				totalsString = totalsString + fmt.Sprintf(" (%.1f%%)", percent)
			}

			totalsString = totalsString + "\n"
		}
	}

	fmt.Printf("%s", totalsString)
	if printTotals {
		fmt.Printf("Total: %s\n", minToHourMin(totalsInt))
	}
}
