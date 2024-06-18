package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getTodayFilename() string {
	today := time.Now()
	return fmt.Sprintf("%04d%02d%02d%s",
		today.Year(),
		today.Month(),
		today.Day(),
		notesExtension)
}

func createNotesFile(path, name string) error {
	filename := filepath.Join(path, getTodayFilename())
	f, err := os.Open(filename)
	defer f.Close()
	if errors.Is(err, os.ErrNotExist) {
		notesHeader, err := getNotesHeader(path)
		if err != nil {
			return err
		}
		notesFooter, err := getCarryoverNotes(path)
		if err != nil {
			return err
		}
		notes := string(notesHeader) + "\n\n\n\n\n" + notesFooter
		os.WriteFile(filename, []byte(notes), 0o766)
	} else if err != nil {
		return err
	}
	return nil
}

func getCarryoverNotes(path string) (string, error) {
	names, err := getLastNFilenames(path, 1)
	if err != nil {
		return "", err
	}
	f, err := os.Open(names[0])
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	carryover := ""
	start := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "#### ‽‽‽‽") {
			start = true
		}
		if !start {
			continue
		}
		if strings.HasPrefix(line, "- [x] ") {
			continue
		}
		if strings.HasPrefix(line, "- [-] ") {
			continue
		}
		carryover += line + "\n"
	}
	return carryover, nil
}

// TODO remove
func getLastFilename(path string) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	lastFile := files[len(files)-1].Name()
	if lastFile == getTodayFilename() {
		lastFile = files[len(files)-2].Name()
	}
	return lastFile, nil
}

// TODO move this to parser
func getLastTaskList(path string) (string, error) {
	names, err := getLastNFilenames(path, 2)
	if err != nil {
		return "", err
	}
	f, err := os.Open(names[1])
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	taskList := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "- [ ]") {
			taskList = taskList + "\n" + line
		}
		if strings.HasPrefix(line, "##") {
			return taskList, nil
		}
	}
	return "", nil
}

func getLastNFiles(path string, n int) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	// TODO bug, this is backwards
	if len(files) <= n {
		return files, nil
	}
	var lastNFiles []fs.FileInfo
	for i := 1; i <= n; i++ {
		lastNFiles = append(lastNFiles, files[len(files)-i])
	}
	return lastNFiles, nil
}

func getLastNFilenames(path string, n int) ([]string, error) {
	var filenames []string
	files, err := getLastNFiles(path, n)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filenames = append(filenames, filepath.Join(path, file.Name()))
	}
	return filenames, nil
}

// TODO is this defunct?
func getLastNotes(path, term string) (string, error) {
	files, err := getLastNFiles(path, 5)
	if err != nil {
		return "", err
	}
	notes := ""
	for i := 2; i < len(files); i++ {
		file := files[len(files)-i]
		f, err := os.Open(filepath.Join(path, file.Name()))
		if err != nil {
			return "", err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		collectNotes := false
		for scanner.Scan() {
			line := scanner.Text()

			// this weird backwards shit is because we have
			// to first find the correct category line, then
			// collect the notes until the next category line
			if collectNotes {
				if strings.HasPrefix(line, "## ") {
					collectNotes = false
					continue
				}
				notes += line + "\n"
			}
			if strings.HasPrefix(line, "## ") && strings.Contains(line, term) {
				collectNotes = true
			}
		}
	}
	return notes, nil
}
