package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	case "links":
		if len(os.Args) != 3 {
			log.Fatal(errors.New("issue ID required"))
		}
		issueFiles, err := searchThroughFiles(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range issueFiles {
			if file == getTodayFilename() {
				continue
			}
			fmt.Printf("|-> [%s](%s)\n", file, filepath.Join(notesPath, file))
		}
		return
	case "today":
		todayFile := filepath.Join(notesPath, getTodayFilename())
		f, err := os.Open(todayFile)
		if errors.Is(err, os.ErrNotExist) {
			notesHeader, err := getNotesHeader()
			if err != nil {
				log.Fatal(err)
			}
			os.WriteFile(todayFile, notesHeader, 0766)
		} else if err != nil {
			log.Fatal(err)
		}
		f.Close()
		fmt.Print(todayFile)
		return
	case "summary":
		t := newTotals()
		t.calculate()
		t.printSummary()
		return
	case "report":
		if len(os.Args) != 3 {
			log.Fatal(errors.New("sprint name required"))
		}
		st := newSprintTotals()
		st.calculate(os.Args[2])
		st.printSummary()
		return
	default:
		log.Fatal(errors.New("command not found: " + os.Args[1]))
	}
}
