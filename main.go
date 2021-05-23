package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	notesPath      = "/Users/mattfinkel/.notes"
	notesExtension = ".md"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("command required"))
	}
	switch os.Args[1] {
	case "today":
		todayFile := fmt.Sprintf("%s/%s", notesPath, getTodayFilename())

		// - parse yesterday's todo list, add items not completed

		file, err := os.Open(todayFile) // For read access.
		if errors.Is(err, os.ErrNotExist) {
			// TODO create file, add header
			os.WriteFile(todayFile, getNotesHeader(), 0766)
		} else if err != nil {
			log.Fatal(err)
		}
		file.Close()

		fmt.Print(todayFile)
		return
	case "summary":
		t := newTotals()
		t.calculate()
		t.printSummary()
		return
	default:
		log.Fatal(errors.New("command not found: " + os.Args[1]))
	}
}
