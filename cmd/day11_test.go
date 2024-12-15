package cmd_test

import (
	"reflect"
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
		gotLeft, gotRight := cmd.Split(v)

		if gotLeft != wantLeft || gotRight != wantRight {
			t.Errorf("Failed for %v: wanted %v %v, got %v %v", v, wantLeft, wantRight, gotLeft, gotRight)
		}
	}
}

func TestBlink(t *testing.T) {
	stones := []int{125, 17}

	want := []int{253000, 1, 7}
	stones = cmd.Blink(stones)
	if !reflect.DeepEqual(want, stones) {
		t.Fatalf("Failed: wanted %+v, got %+v", want, stones)
	}

	want = []int{253, 0, 2024, 14168}
	stones = cmd.Blink(stones)
	if !reflect.DeepEqual(want, stones) {
		t.Fatalf("Failed: wanted %+v, got %+v", want, stones)
	}

	want = []int{512072, 1, 20, 24, 28676032}
	stones = cmd.Blink(stones)
	if !reflect.DeepEqual(want, stones) {
		t.Fatalf("Failed: wanted %+v, got %+v", want, stones)
	}

	want = []int{512, 72, 2024, 2, 0, 2, 4, 2867, 6032}
	stones = cmd.Blink(stones)
	if !reflect.DeepEqual(want, stones) {
		t.Fatalf("Failed: wanted %+v, got %+v", want, stones)
	}

	want = []int{1036288, 7, 2, 20, 24, 4048, 1, 4048, 8096, 28, 67, 60, 32}
	stones = cmd.Blink(stones)
	if !reflect.DeepEqual(want, stones) {
		t.Fatalf("Failed: wanted %+v, got %+v", want, stones)
	}

	want = []int{2097446912, 14168, 4048, 2, 0, 2, 4, 40, 48, 2024, 40, 48, 80, 96, 2, 8, 6, 7, 6, 0, 3, 2}
	stones = cmd.Blink(stones)
	if !reflect.DeepEqual(want, stones) {
		t.Fatalf("Failed: wanted %+v, got %+v", want, stones)
	}
}
