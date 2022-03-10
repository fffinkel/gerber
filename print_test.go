package main

import (
	"fmt"
	"testing"
)

func TestGetSummary(t *testing.T) {
	t.Parallel()

	testSummary := `
This week you have worked: 1h 58m
Today you have worked: 0h 24m

 ➔ cat1: 0h 4m (16.7%d, 2.5%w)
 ➔ cat2: 0h 20m (83.3%d, 5.1%w)

You are currently working on: something fun
`

	totals := newTestTotals()
	summary := getSummary(&totals)

	fmt.Printf("\n\n----------> %+v\n", summary)
	if string(summary) != testSummary {
		t.Error("summary should match expected summary")
	}
}
