package cmd_test

import (
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
)

func TestSplit(t *testing.T) {
	cases := [][]int{
		{   10,  1,  0 },
		{   99,  9,  9 },
		{ 1000, 10,  0 },
		{ 5696, 56, 96 },
	}

	for _, c := range cases {
		v := c[0]
		wantLeft := c[1]
		wantRight := c[2]
		gotLeft, gotRight, ok := cmd.Split(v)

		if !ok {
			t.Fatalf("Split failed unexpectedly")
		}

		if gotLeft != wantLeft || gotRight != wantRight {
			t.Errorf("Failed for %v: wanted %v %v, got %v %v", v, wantLeft, wantRight, gotLeft, gotRight)
		}
	}
}
