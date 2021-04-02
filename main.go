package main

import (
	"errors"
	"log"
	"os"
)

const (
	notesPath      = "/Users/mattfinkel/.notes/"
	notesExtension = ".md"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("command required"))
	}
	switch os.Args[1] {
	case "summary":
		t := newTotals()
		t.calculate()
		t.printSummary()
		return
	default:
		log.Fatal(errors.New("command not found: " + os.Args[1]))
	}
}
