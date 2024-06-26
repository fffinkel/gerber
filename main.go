package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	// get the last five filenames
	case "recent":
		if err := createNotesFile(notesPath, getTodayFilename()); err != nil {
			log.Fatal(err)
		}
		files, err := getLastNFilenames(notesPath, 8)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(strings.Join(files, " "))
		return

	// summary used in prompt
	case "summary":
		t := newTotals(notesPath, 15)
		err := t.calculate(time.Now(), 15)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(getSummary(t))
		return

	// summary used in prompt
	case "ledger":
		if len(os.Args) != 3 {
			log.Fatal(errors.New("days required"))
		}
		days, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		t := newTotals(notesPath, days)
		err = t.calculate(time.Now(), days)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(getSummary(t))
		return

	default:
		log.Fatal(errors.New("command not found: " + os.Args[1]))
	}
}
