package cmd

import (
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func TestTotalDistanceSuccess(t *testing.T) {
	left  := util.Shuffle([]int{ 1, 2, 3, 4 })
	right := util.Shuffle([]int{ 5, 6, 7, 8 })

	want := 4 * 4
	got, err := cmd.TotalDistance(left, right)

	if err != nil {
		t.Fatalf(`TotalDistance failed with error: %#q`, err)

	}
	if got != want {
		t.Fatalf(`TotalDistance failed on values: want %v, got %v`, want, got)
	}

	left = util.Shuffle([]int{ -1, -2, -3, -4 })
	right = util.Shuffle([]int{ -5, -6, -7, -8 })

	got, err = cmd.TotalDistance(left, right)

	if err != nil {
		t.Fatalf(`TotalDistance failed with error: %#q`, err)

	}
	if got != want {
		t.Fatalf(`TotalDistance failed on values: want %v, got %v`, want, got)
	}
}

func TestTotalDistanceFailLengths(t *testing.T) {
	left := util.Shuffle([]int{ 1, 2, 3, 4 })
	right := util.Shuffle([]int{ })

	_, err := cmd.TotalDistance(left, right)

	if err == nil {
		t.Fatalf(`Did not get error!`)

	}
}

func TestSimilaritySuccess(t *testing.T) {
	left  := util.Shuffle([]int{ 1, 2, 3, 4, 5, 4, 3, 2, 1 })
	right := util.Shuffle([]int{ 8, 7, 6, 5, 4, 5, 6, 7, 8 })

	want := (1 * 2 * 0) + (2 * 2 * 0) + (3 * 2 * 0) + (4 * 2 * 1) + (5 * 1 * 2) + (6 * 0 * 2) + (7 * 0 * 2) + (8 * 0 * 2)
	got, err := cmd.Similarity(left, right)

	if err != nil {
		t.Fatalf(`Similarity failed with error: %#q`, err)

	}
	if got != want {
		t.Fatalf(`Similarity failed on values: want %v, got %v`, want, got)
	}

}

func TestSimilarityFailLengths(t *testing.T) {
	left := util.Shuffle([]int{ 1, 2, 3, 4 })
	right := util.Shuffle([]int{ })

	_, err := cmd.Similarity(left, right)

	if err == nil {
		t.Fatalf(`Did not get error!`)

	}
}
