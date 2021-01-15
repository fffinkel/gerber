package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const notesDir = "/Users/mattfinkel/.notes/"

// TODO
// - today so far
// - week so far
// - total hours worked

func main() {
	fmt.Println("------------------")
	dayTotals := make(map[string]int)
	weekTotals := make(map[string]int)

	files, err := ioutil.ReadDir(notesDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		f, err := os.Open(filepath.Join(notesDir, file.Name()))
		if err != nil {
			log.Fatal(err)
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
				//fmt.Printf("%+v\n", file.Name())
				//fmt.Printf("%+v\n", parsedDate)
				//fmt.Printf("%+v\n", lastDate)
				//fmt.Printf("%+v\n", minutes)
				if parsedDate.Day() == lastDate.Day() && parsedDate.Month() == lastDate.Month() {
					//fmt.Printf("The minutes diffrence for category [%s] is: %+v\n", lastCategory, minutes)
					weekTotals[lastCategory] += int(minutes)

					if parsedDate.Day() == time.Now().Day() && parsedDate.Month() == time.Now().Month() {
						dayTotals[lastCategory] += int(minutes)
					}
				}

				lastDate = parsedDate
				lastCategory = strings.TrimSuffix(l[1], ")")

				//fmt.Println("My Date Reformatted:\t", parsedDate.Format(time.RFC822))
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}
	//fmt.Printf("weekTotals: %+v\n", weekTotals)
	//fmt.Printf("dayTotals: %+v\n", dayTotals)

	// if lastDate.Day() == time.Now().Day() {
	// 	now := time.Now()
	// 	currentMinutes := now.Sub(lastDate).Minutes()
	// 	weekTotals[lastCategory] += int(currentMinutes)
	// 	dayTotals[lastCategory] += int(currentMinutes)
	// }

	printTotals(dayTotals)
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
	fmt.Println("------------------")
}

func minToHourMin(m int) string {
	minutes := m % 60
	hours := m / 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
