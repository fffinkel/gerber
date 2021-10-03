package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type sprintTotals struct {
	notesPath string
	sprint    map[string]int
	other     map[string]int
}

func newSprintTotals(notesPath string) *sprintTotals {
	return &sprintTotals{
		notesPath: notesPath,
		sprint:    make(map[string]int),
		other:     make(map[string]int),
	}
}

func (st *sprintTotals) total() int {
	return st.sprintTotal() + st.otherTotal()
}

func (st *sprintTotals) sprintTotal() int {
	total := 0
	for _, v := range st.sprint {
		total += v
	}
	return total
}

func (st *sprintTotals) sprintPercent() float64 {
	return (float64(st.sprintTotal()) / float64(st.total())) * 100
}

func (st *sprintTotals) sprintThemeTotals() map[string]float64 {
	totals := make(map[string]float64)
	for k, _ := range st.sprint {
		theme := strings.Split(k, ", ")[0]
		totals[theme] += float64(st.sprint[k])
	}
	return totals
}

func (st *sprintTotals) sprintThemePercent(theme string) float64 {
	return (st.sprintThemeTotals()[theme] / float64(st.total())) * 100
}

func (st *sprintTotals) otherTotal() int {
	total := 0
	for _, v := range st.other {
		total += v
	}
	return total
}

func (st *sprintTotals) otherPercent() float64 {
	return (float64(st.otherTotal()) / float64(st.total())) * 100
}

func (st *sprintTotals) otherThemeTotals() map[string]float64 {
	totals := make(map[string]float64)
	for k, _ := range st.other {
		theme := strings.Split(k, ", ")[0]
		totals[theme] += float64(st.other[k])
	}
	return totals
}

func (st *sprintTotals) otherThemePercent(theme string) float64 {
	return (st.otherThemeTotals()[theme] / float64(st.total())) * 100
}

func (st *sprintTotals) printSummary() error {

	fmt.Print("\n")
	fmt.Printf("This sprint you worked: %+v\n", minToHourMin(st.total()))
	fmt.Print("\n")

	fmt.Printf("Sprint work: %+v (%.1f%%)\n",
		minToHourMin(st.sprintTotal()), st.sprintPercent())
	sprintThemeTotals := make(map[string]int)
	for _, category := range sortedKeys(st.sprint) {
		theme := strings.Split(category, ", ")[0]
		sprintThemeTotals[theme] += st.sprint[category]
	}
	for _, theme := range sortedKeys(sprintThemeTotals) {
		fmt.Printf("-- %s: %s (%.1f%%)\n", theme, minToHourMin(sprintThemeTotals[theme]), st.sprintThemePercent(theme))
	}

	fmt.Print("\n")
	fmt.Printf("Other work: %+v (%.1f%%)\n",
		minToHourMin(st.otherTotal()), st.otherPercent())
	otherThemeTotals := make(map[string]int)
	for _, category := range sortedKeys(st.other) {
		theme := strings.Split(category, ", ")[0]
		otherThemeTotals[theme] += st.other[category]
	}
	for _, theme := range sortedKeys(otherThemeTotals) {
		fmt.Printf("-- %s: %s (%.1f%%)\n", theme, minToHourMin(otherThemeTotals[theme]), st.otherThemePercent(theme))
	}

	return nil
}

func (st *sprintTotals) calculate(sprintName string) error {
	files, err := ioutil.ReadDir(st.notesPath)
	if err != nil {
		return err
	}

	var lastDate time.Time
	var lastCategory string
	var lastIsSprint bool

FILE:
	for _, file := range files {
		f, err := os.Open(filepath.Join(st.notesPath, file.Name()))
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNumber := 0
		for scanner.Scan() {

			line := scanner.Text()
			lineNumber += 1

			if strings.HasPrefix(line, "# ") {
				lineParts := strings.Split(line, " (")
				if len(lineParts) != 2 {
					continue FILE
				}
				foundSprintName := strings.TrimSuffix(lineParts[1], ")")
				if sprintName != foundSprintName {
					continue FILE
				}
			}

			if strings.HasPrefix(line, "## ") {
				category, date, isSprint := parseLine(line, file, lineNumber)

				minutes := date.Sub(lastDate).Minutes()
				if lastCategory != "break" && lastCategory != "" {
					if lastIsSprint {
						st.sprint[lastCategory] += int(minutes)
					} else {
						st.other[lastCategory] += int(minutes)
					}
				}
				lastDate = date
				lastCategory = category
				lastIsSprint = isSprint
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	return nil
}
