package main

import (
	"testing"
	"time"
)

func TestMinToHourMin(t *testing.T) {
	t.Parallel()
	cases := map[int]string{
		0:    "0h 0m",
		1:    "0h 1m",
		30:   "0h 30m",
		60:   "1h 0m",
		61:   "1h 1m",
		120:  "2h 0m",
		121:  "2h 1m",
		1440: "24h 0m",
		1500: "25h 0m",
		9999: "166h 39m",
	}
	for k, v := range cases {
		if minToHourMin(k) != v {
			t.Errorf("%d minutes should be %s", k, v)
		}
	}
}

func TestSortedKeys(t *testing.T) {
	t.Parallel()
	type tCase struct {
		in  map[string]int
		out []string
	}
	cases := []tCase{
		tCase{
			in: map[string]int{
				"a": 1,
			},
			out: []string{
				"a",
			},
		},
		tCase{
			in: map[string]int{
				"b": 1,
				"a": 1,
			},
			out: []string{
				"a",
				"b",
			},
		},
		tCase{
			in: map[string]int{
				"y": 1,
				"b": 1,
				"a": 1,
				"z": 1,
			},
			out: []string{
				"a",
				"b",
				"y",
				"z",
			},
		},
	}
	for _, v := range cases {
		for i, x := range sortedKeys(v.in) {
			if x != v.out[i] {
				t.Errorf("map not sorted correctly")
			}
		}
	}
}

func TestWeekNumber(t *testing.T) {
	t.Parallel()
	times := map[int]time.Time{}
	times[1], _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	times[2], _ = time.Parse(time.RFC3339, "2006-01-09T15:04:05Z")
	times[17], _ = time.Parse(time.RFC3339, "2006-04-30T15:04:05Z")
	times[22], _ = time.Parse(time.RFC3339, "2006-06-01T15:04:05Z")
	times[52], _ = time.Parse(time.RFC3339, "2006-12-31T15:04:05Z")
	for k, v := range times {
		if weekNumber(v) != k {
			t.Errorf("incorrect week number %d for time %+v", weekNumber(v), v)
		}
	}
}
