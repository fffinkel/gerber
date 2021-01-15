package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	dayTotals := make(map[string]int)
	weekTotals := make(map[string]int)

	files, err := ioutil.ReadDir(notesPath)
	if err != nil {
		return err
	}

	for _, file := range files {

		f, err := os.Open(filepath.Join(notesPath, file.Name()))
		if err != nil {
			return err
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
				l := strings.Split(trimmed, " (")

				if len(l) != 2 {
					fmt.Printf("found malformed, %s line %d: %s\n", file.Name(), lineNumber, line)
					continue
				}

				date := l[0]
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
				lastCategory = strings.TrimSuffix(l[1], ")")
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

	}

	printTotals(dayTotals)
	return nil
}

func printTotals(totals map[string]int) {
	totalsString := ""
	totalsInt := 0
	for k, v := range totals {
		if k != "break" {
			hm := minToHourMin(v)
			totalsString = totalsString + fmt.Sprintf("%s: %s\n", k, hm)
			totalsInt += v
		}
	}
	fmt.Printf("%s", totalsString)
	fmt.Printf("Total: %s\n", minToHourMin(totalsInt))
}

func minToHourMin(m int) string {
	minutes := m % 60
	hours := m / 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
