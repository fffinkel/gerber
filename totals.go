package main

import (
	"bufio"

	// "fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type totals struct {
	sowIndex  int
	notesPath string
	current   string
	days      []map[string]int
}

func newTotals(notesPath string, days int) *totals {
	return &totals{
		notesPath: notesPath,
		days:      make([]map[string]int, days),
	}
}

func (t *totals) weekTotalMinutes() int {
	total := 0
	for i := 0; i <= t.sowIndex; i++ {
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

func (t *totals) nDayProductiveTotalMinutes(n int) int {
	total := 0
	for i := 0; i < n; i++ {
		for k, _ := range t.days[i] {
			theme := strings.Split(k, ", ")[0]
			if isProductiveTheme(theme) {
				total += t.days[i][k]
			}
		}
	}
	return total
}

func (t *totals) nDayProductivePercent(n int) float64 {
	return (float64(t.nDayProductiveTotalMinutes(n)) / float64(t.nDayTotalMinutes(n))) * 100
}

func (t *totals) calculate(today time.Time, days int) error {

	filenames, err := getLastNFilenames(t.notesPath, days)
	if err != nil {
		return err
	}

	var thisCategoryStart time.Time
	var thisCategory string

	fileIndex := 0
	foundSOW := false
	// for each file
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		t.days[fileIndex] = make(map[string]int)
		timestampLineCount := 0
		scanner := bufio.NewScanner(f)

		// for each line
		for scanner.Scan() {
			line := scanner.Text()

			// if this is a timestamp line
			if strings.HasPrefix(line, "## ") {
				timestampLineCount++

				// get the thisCategory and timestamp
				category, date, err := parseLine(line)
				if err != nil {
					return err
				}

				// TODO this comment
				// set the current day of the week based on the _first_ timestamp
				if timestampLineCount == 1 {

					// do we need to account for the case where we have not created a file
					// yet? probably only is an issue on monday
					// - yes!

					// cases to think about
					// - [ ] today is monday, and I have a notes file
					// - [ ] today is monday, and I don't have a notes file
					// - [ ] today is tuesday 12:30AM, still working from monday
					// - [ ] today is tuesday, and I have a notes file
					// - [ ] today is tuesday, and I don't have a notes file
					// - [ ] today is sunday, last day of work was friday

					// if this is the first file and the first file's day is monday
					// - monday index is zero
					// - we're done looking back

					// if this is not the first file
					// - monday index is the first one that we find (with a max search of 7)

					//if fileIndex == 0 && fileDOW == "Mon" {

					// just found another bug where when there is no monday file, we can't
					// reliably look for the file to see when we started the week

					// XXX works, ish!
					// - problem: when it's monday and I haven't creted a file yet, it
					// thinks that it's the last day that a file existed
					if fileIndex >= 0 && !foundSOW {
						fileDOW := string([]byte(date.Weekday().String())[:3])
						if fileDOW == "Mon" {
							foundSOW = true
							t.sowIndex = fileIndex
						}
					}

					// fmt.Printf("\n\n----------> %+v\n", t.sowIndex)

					// if t.sowIndex == 0 && fileIndex < 7 && todayDOW != "Mon" && fileDOW == "Mon" {
					// }

					// set the thisCategory
					thisCategory = category
					thisCategoryStart = date
					continue
				}

				// if the thisCategory is not "break," (meaning we're currently working on
				// something) add its minutes
				if thisCategory != "break" {
					minutes := date.Sub(thisCategoryStart).Minutes()
					t.days[fileIndex][thisCategory] += int(minutes)
				}

				thisCategoryStart = date
				thisCategory = category
			}
		}

		// if we end with a currently open task, add its minutes
		if thisCategory != "break" {
			minutes := int(today.Sub(thisCategoryStart).Minutes())
			t.days[fileIndex][thisCategory] += minutes
			t.current = thisCategory
		}

		fileIndex++
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
