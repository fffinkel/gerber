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
	fmt.Println("leatherman: notes")
	totals := make(map[string]int)

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
				if parsedDate.Day() == lastDate.Day() {
					fmt.Printf("The minutes diffrence for category [%s] is: %+v\n", lastCategory, minutes)
				}

				lastDate = parsedDate
				totals[lastCategory] += int(minutes)
				lastCategory = strings.TrimSuffix(l[1], ")")

				//fmt.Println("My Date Reformatted:\t", parsedDate.Format(time.RFC822))
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}
	fmt.Printf("totals: %+v\n", totals)
}