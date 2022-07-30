package main

import (
	"bufio"
	"os"
	"sort"
	"strings"
	"time"
)

type totals struct {
	mondayIndex int
	notesPath   string
	current     string
	days        []map[string]int
}

func newTotals(notesPath string) *totals {
	return &totals{
		notesPath: notesPath,
		days:      make([]map[string]int, 15),
	}
}

// TODO calculate how many days ago monday was
func (t *totals) weekTotalMinutes() int {
	total := 0
	for i := 0; i <= t.mondayIndex; i++ {
		for _, v := range t.days[i] {
			total += v
		}
	}
	return total
}

func (t *totals) nDayTotalMinutes(n int) int {
	total := 0
	for i := 0; i < n; i++ {
		for _, v := range t.days[i] {
			total += v
		}
	}
	return total
}

func (t *totals) nDayThemeTotalMinutes(n int) map[string]float64 {
	totals := make(map[string]float64)
	for i := 0; i < n; i++ {
		for k, _ := range t.days[i] {
			theme := strings.Split(k, ", ")[0]
			totals[theme] += float64(t.days[i][k])
		}
	}
	return totals
}

func (t *totals) nDayThemePercent(n int, theme string) float64 {
	return (t.nDayThemeTotalMinutes(n)[theme] / float64(t.nDayTotalMinutes(n))) * 100
}

func (t *totals) calculate(today time.Time) error {
	filenames, err := getLastNFilenames(t.notesPath, 15)
	if err != nil {
		return err
	}

	var lastDate time.Time
	var category string

	i := 0
	// for each file
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		t.days[i] = make(map[string]int)

		lineNumber := 0
		scanner := bufio.NewScanner(f)

		// for each line
		for scanner.Scan() {
			line := scanner.Text()

			// if this is a timestamp line
			if strings.HasPrefix(line, "## ") {
				lineNumber++

				// get the category and timestamp
				nextCategory, date, err := parseLine(line)
				if err != nil {
					return err
				}

				// calculate which file was Monday
				dow := string([]byte(date.Weekday().String())[:3])
				if dow == "Mon" && t.mondayIndex == 0 {
					t.mondayIndex = i
				}

				// if this is the first line, set it and move on
				if lineNumber == 1 {
					category = nextCategory
					lastDate = date
					continue
				}

				// if the category is not "break," add its minutes
				if category != "break" {
					minutes := date.Sub(lastDate).Minutes()
					t.days[i][category] += int(minutes)
				}

				lastDate = date
				category = nextCategory
			}
		}

		// if we end with a currently open task, add its minutes
		if category != "break" {
			minutes := int(today.Sub(lastDate).Minutes())
			t.days[i][category] += minutes
			t.current = category
		}

		// TODO this needs a defer
		defer f.Close()
		i++

	}

	return nil
}

// sort by 15 day theme percentage
func (t *totals) sortedKeysByThemePercentage(in map[string]int) []string {
	var sortedThemes []string
	for k, _ := range in {
		sortedThemes = append(sortedThemes, k)
	}
	for i := 0; i < len(sortedThemes)-1; i++ {
		for j := 0; j < len(sortedThemes)-i-1; j++ {
			if t.nDayThemePercent(15, sortedThemes[j]) == t.nDayThemePercent(15, sortedThemes[j+1]) {
				alph := []string{sortedThemes[j], sortedThemes[j+1]}
				sort.Strings(alph)
				sortedThemes[j] = alph[0]
				sortedThemes[j+1] = alph[1]
			}
			if t.nDayThemePercent(15, sortedThemes[j]) < t.nDayThemePercent(15, sortedThemes[j+1]) {
				sortedThemes[j], sortedThemes[j+1] = sortedThemes[j+1], sortedThemes[j]
			}
		}
	}
	return sortedThemes
}
