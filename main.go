package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const notesDir = "/Users/mattfinkel/.notes/"

func main() {
	fmt.Println("leatherman: notes")
	totals := make(map[string]int)

	files, err := ioutil.ReadDir(notesDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())

		f, err := os.Open(filepath.Join(notesDir, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNumber := 0
		for scanner.Scan() {
			line := scanner.Text()
			lineNumber += 1

			if strings.HasPrefix(line, "## ") {

				trimmed := strings.TrimPrefix(line, "## ")
				l := strings.Split(trimmed, " (")
				//date := l[0]

				if len(l) != 2 {
					fmt.Printf("found malformed, line %d, %s\n", lineNumber, line)
					continue
				}

				category := strings.TrimSuffix(l[1], ")")
				totals[category] += 1
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}
	fmt.Printf("this is totals:\n%+v\n", totals)
}
