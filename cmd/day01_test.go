package cmd

import (
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func TestTotalDistanceSuccess(t *testing.T) {
	leftNumbers := util.Shuffle([]int{ 1, 2, 3, 4 })
	rightNumbers := util.Shuffle([]int{ 5, 6, 7, 8 })

	want := 4 * 4
	got, err := cmd.TotalDistance(leftNumbers, rightNumbers)

	if err != nil {
		t.Fatalf(`TotalDistance failed with error: %#q`, err)

	}
	if got != want {
		t.Fatalf(`TotalDistance failed on values: want %v, got %v`, want, got)
	}

	leftNumbers = util.Shuffle([]int{ -1, -2, -3, -4 })
	rightNumbers = util.Shuffle([]int{ -5, -6, -7, -8 })

	got, err = cmd.TotalDistance(leftNumbers, rightNumbers)

	if err != nil {
		t.Fatalf(`TotalDistance failed with error: %#q`, err)

	}
	if got != want {
		t.Fatalf(`TotalDistance failed on values: want %v, got %v`, want, got)
	}
}

func TestTotalDistanceFailLengths(t *testing.T) {
	leftNumbers := util.Shuffle([]int{ 1, 2, 3, 4 })
	rightNumbers := util.Shuffle([]int{ })

	_, err := cmd.TotalDistance(leftNumbers, rightNumbers)

	if err == nil {
		t.Fatalf(`Did not get error!`)

	}
}
