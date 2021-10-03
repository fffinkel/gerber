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

func parseLine(line string, file os.FileInfo, lineNumber int) (string, time.Time, bool) {
	trimmed := strings.TrimPrefix(line, "## ")
	lineParts := strings.Split(trimmed, " (")

	if len(lineParts) != 2 {
		panic(fmt.Sprintf("found malformed, %s line %d: %s\n", file.Name(), lineNumber, line)) // TODO no panics
	}

	date := lineParts[0]
	parsedDate, err := time.Parse("2006-01-02 Mon 15:04 PM MST", date)
	if err != nil {
		panic(fmt.Sprintf("error parsing date, %s line %d: %s\n", file.Name(), lineNumber, date)) // TODO no panics
	}

	isSprint := false
	category := strings.TrimSuffix(lineParts[1], ")")
	categoryTrimmed := strings.TrimPrefix(category, "*")
	if categoryTrimmed != category {
		isSprint = true
	}

	return categoryTrimmed, parsedDate, isSprint

	// splitCategory := strings.Split(strings.TrimSuffix(lineParts[1], ")"), ", ")
	// if len(splitCategory) < 2 {
	// 	return splitCategory[0], "", parsedDate
	// }
	// return splitCategory[0], splitCategory[1], parsedDate
}

// TODO remove?
func printTodayFilename(path string) error {
	fmt.Printf(filepath.Join(path, getTodayFilename()))
	return nil
}

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

func getNotesHeader(path string) ([]byte, error) {
	today := time.Now()
	weekday := string([]byte(today.Weekday().String())[:3])
	kitchen := today.Format("03:04 PM")
	zone, _ := today.Zone()

	taskList, err := getLastTaskList(path)
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("# %s %02d, %04d (q‽‽‽‽‽‽‽‽‽‽s‽‽‽‽‽‽‽‽‽‽)\n\nToday is a beautiful day.%s\n\n## %04d-%02d-%02d %s %s %s (admin)\n\n\n",
		today.Month(),
		today.Day(),
		today.Year(),
		taskList,
		today.Year(),
		today.Month(),
		today.Day(),
		weekday,
		kitchen,
		zone,
	)), nil
}

func getLastNFiles(path string, n int) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return files[len(files)-n:], nil
}

func searchThroughFiles(path, issueID string) ([]string, error) {
	files, err := getLastNFiles(path, 50)
	if err != nil {
		return nil, err
	}
	foundFiles := make(map[string]bool)
	foundList := make([]string, 0)
	for _, file := range files {
		f, err := os.Open(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNumber := 0
		for scanner.Scan() {
			line := scanner.Text()
			lineNumber += 1

			if strings.HasPrefix(line, "# ") && strings.Contains(line, "????‽????‽") {
				panic("found file with bad sprint name: " + line)
			}

			if strings.HasPrefix(line, "## ") {
				category, _, _ := parseLine(line, file, lineNumber)
				if strings.HasSuffix(category, issueID) {
					if foundFiles[file.Name()] {
						continue
					}
					foundFiles[file.Name()] = true
					foundList = append(foundList, file.Name())
				}
			}
		}
	}
	return foundList, nil
}
