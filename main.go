package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	notesPath          = "/Users/mattfinkel/.notes/"
	notesFileExtension = ".md"
)

type commandFunc func() error

var commandMap = map[string]commandFunc{
	"print_today_work":     printTodayWork,
	"print_today_filename": printTodayFilename,
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal(errors.New("command required"))
	}
	commandArg := os.Args[1]

	if commandFunc, ok := commandMap[commandArg]; ok {
		err := commandFunc()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(errors.New(
			fmt.Sprintf("command [%s] not found", commandArg)),
		)
	}
}

func printTodayFilename() error {
	fmt.Printf(filepath.Join(notesPath, getTodayFilename()))
	return nil
}

func getTodayFilename() string {
	today := time.Now()
	return fmt.Sprintf(
		"%d%02d%02d%s",
		today.Year(),
		today.Month(),
		today.Day(),
		notesFileExtension,
	)
}

func printTodayWork() error {

	files, err := ioutil.ReadDir(notesPath)
	if err != nil {
		return err
	}

	dayTotals, _, current, err := getTotals(files)
	if err != nil {
		return err
	}

	printTotals(dayTotals, current)
	return nil
}

func minToHourMin(m int) string {
	minutes := m % 60
	hours := m / 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func printTotals(totals map[string]int, current string) {
	totalsString := ""
	totalsInt := 0

	sorted := make([]string, 0, len(totals))
	for k := range totals {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	for _, category := range sorted {
		if category != "break" {
			hm := minToHourMin(totals[category])
			inProgress := ""
			if category == current {
				inProgress = "(IN PROGRESS) "
			}
			totalsString = totalsString + fmt.Sprintf("%s%s: %s\n", inProgress, category, hm)
			totalsInt += totals[category]
		}
	}

	fmt.Printf("%s", totalsString)
	fmt.Printf("Total: %s\n", minToHourMin(totalsInt))
}

func getTotals(filePaths []os.FileInfo) (map[string]int, map[string]int, string, error) {
	dayTotals := make(map[string]int)
	weekTotals := make(map[string]int)
	currentCategory := ""

	for _, file := range filePaths {

		f, err := os.Open(filepath.Join(notesPath, file.Name()))
		if err != nil {
			return map[string]int{}, map[string]int{}, currentCategory, err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNumber := 0
		var lastDate time.Time
		var lastCategory string
		for scanner.Scan() {
			line := scanner.Text()
			lineNumber += 1

			if strings.HasPrefix(line, "## ") {

				trimmed := strings.TrimPrefix(line, "## ")
				lineParts := strings.Split(trimmed, " (")

				if len(lineParts) != 2 {
					fmt.Printf("found malformed, %s line %d: %s\n", file.Name(), lineNumber, line)
					continue
				}

				date := lineParts[0]
				parsedDate, err := time.Parse("2006-01-02 Mon 15:04 PM MST", date)
				if err != nil {
					fmt.Printf("error parsing date, %s line %d: %s\n", file.Name(), lineNumber, date)
					continue
				}

				minutes := parsedDate.Sub(lastDate).Minutes()
				if parsedDate.Day() == lastDate.Day() && parsedDate.Month() == lastDate.Month() {
					weekTotals[lastCategory] += int(minutes)

					if parsedDate.Day() == time.Now().Day() && parsedDate.Month() == time.Now().Month() {
						dayTotals[lastCategory] += int(minutes)
					}
				}

				lastDate = parsedDate
				lastCategory = strings.TrimSuffix(lineParts[1], ")")
			}
		}

		if lastCategory != "break" {
			currentCategory = lastCategory
			currentDate := time.Now()
			minutes := currentDate.Sub(lastDate).Minutes()
			if currentDate.Day() == lastDate.Day() && currentDate.Month() == lastDate.Month() {
				weekTotals[lastCategory] += int(minutes)

				if currentDate.Day() == time.Now().Day() && currentDate.Month() == time.Now().Month() {
					dayTotals[lastCategory] += int(minutes)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return map[string]int{}, map[string]int{}, "", err
		}
	}

	return dayTotals, weekTotals, currentCategory, nil
}
