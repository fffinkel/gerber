package main

import (
	"bufio"
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

func getLastTaskList(path string) (string, error) {
	lastFile, err := getLastFilename(path)
	if err != nil {
		return "", err
	}
	f, err := os.Open(filepath.Join(path, lastFile))
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
	if len(files) <= n {
		return files, nil
	}
	return files[len(files)-n:], nil
}

func getLastNotes(path, term string) (string, error) {
	files, err := getLastNFiles(path, 50)
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
