package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func NoTestParseLine(t *testing.T) {
	t.Parallel()
	category, parsedDate, err := parseLine(
		"## 2021-12-02 Thu 02:46 PM WET (toil, deprovisions)")
	if err != nil {
		t.Errorf("test got error: %s", err.Error())
	}
	if category != "toil, deprovisions" {
		t.Errorf("got incorrect category: %s", category)
	}
	testDate, _ := time.Parse(
		"2006-01-02 Mon 15:04 PM MST",
		"2021-12-02 Thu 02:46 PM WET")
	fmt.Printf("\n\n-p---------> %s\n", parsedDate)
	fmt.Printf("\n\n-t---------> %s\n", testDate)
	if parsedDate != testDate {
		t.Errorf("got incorrect date: %+v", parsedDate)
	}

	category, parsedDate, err = parseLine(
		"## 2021-12-02 Thu 04:15 PM WET (break)")
	if err != nil {
		t.Errorf("test got error: %s", err.Error())
	}
	if category != "break" {
		t.Errorf("got incorrect category: %s", category)
	}
	testDate, _ = time.Parse(time.RFC3339, "2021-12-02T16:15:00+00:00")
	if parsedDate != testDate {
		t.Errorf("got incorrect date: %+v", parsedDate)
	}

	category, parsedDate, err = parseLine(
		"## 2021-12-02 Thu 05:00 PM WET (fun, gerber work)")
	if err != nil {
		t.Errorf("test got error: %s", err.Error())
	}
	if category != "fun, gerber work" {
		t.Errorf("got incorrect category: %s", category)
	}
	testDate, _ = time.Parse(time.RFC3339, "2021-12-02T17:00:00+00:00")
	if parsedDate != testDate {
		t.Errorf("got incorrect date: %+v", parsedDate)
	}

	category, parsedDate, err = parseLine(
		"## 2021-12-02 Thu 06:00 PM WET (more fun, gerber work, wow)")
	if err != nil {
		t.Errorf("test got error: %s", err.Error())
	}
	if category != "more fun, gerber work, wow" {
		t.Errorf("got incorrect category: %s", category)
	}
	testDate, _ = time.Parse(time.RFC3339, "2021-12-02T18:00:00+00:00")
	if parsedDate != testDate {
		t.Errorf("got incorrect date: %+v", parsedDate)
	}
}

func TestParseLineMalformedLineError(t *testing.T) {
	t.Parallel()
	_, _, err := parseLine(
		"## 2021-12-02 Thu 02:46 PM WET")
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
	if !strings.Contains(err.Error(), "found malformed line") {
		t.Errorf("did not get malformed line error")
	}
}

func TestParseLineDateParseError(t *testing.T) {
	t.Parallel()
	_, _, err := parseLine(
		"## 2021-12-99 BLAH 02:46 PM NOO (fun, gerber work)")
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
	if !strings.Contains(err.Error(), "error parsing date") {
		t.Errorf("did not get date parse error")
	}
}
