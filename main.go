package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	notesExtension = ".md"
)

func main() {
	notesPath := "/Users/mattfinkel/.notes"
	if len(os.Args) < 2 {
		log.Fatal(errors.New("command required"))
	}
	switch os.Args[1] {

	// find all notes for the search term
	case "notes":
		if len(os.Args) != 3 {
			log.Fatal(errors.New("search term required"))
		}
		notes, err := getLastNotes(notesPath, os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		todayFile := filepath.Join(notesPath, getTodayFilename())
		fmt.Printf("%s %s\n", notes, todayFile)
		return

	// print the notes file template
	case "today":
		todayFile := filepath.Join(notesPath, getTodayFilename())
		f, err := os.Open(todayFile)
		if errors.Is(err, os.ErrNotExist) {
			notesHeader, err := getNotesHeader(notesPath)
			if err != nil {
				log.Fatal(err)
			}
			os.WriteFile(todayFile, notesHeader, 0o766)
		} else if err != nil {
			log.Fatal(err)
		}
		f.Close()
		fmt.Print(todayFile)
		return

	// summary used in prompt
	case "summary":
		t := newTotals(notesPath)
		t.calculate(time.Now())
		printSummary(t)
		return

	default:
		log.Fatal(errors.New("command not found: " + os.Args[1]))
	}
}
