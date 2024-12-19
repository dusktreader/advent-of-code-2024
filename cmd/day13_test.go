package cmd_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func TestCount(t *testing.T) {
	b := cmd.Button{
		A:     util.MakeVector(94, 34),
		B:     util.MakeVector(22, 67),
		Prize: util.MakePoint(8400, 5400),
	}
	want := 280
	got  := b.Count()

	if want != got {
		t.Errorf("Didn't match: wanted %v, got %v", want, got)
	}

	b = cmd.Button{
		A:     util.MakeVector(26, 66),
		B:     util.MakeVector(67, 21),
		Prize: util.MakePoint(12748, 12176),
	}
	want = math.MaxInt
	got  = b.Count()

	if want != got {
		t.Errorf("Didn't match: wanted %v, got %v", want, got)
	}

	b = cmd.Button{
		A:     util.MakeVector(17, 86),
		B:     util.MakeVector(84, 37),
		Prize: util.MakePoint(7870, 6450),
	}
	want = 200
	got  = b.Count()

	if want != got {
		t.Errorf("Didn't match: wanted %v, got %v", want, got)
	}

	b = cmd.Button{
		A:     util.MakeVector(69, 23),
		B:     util.MakeVector(27, 71),
		Prize: util.MakePoint(18641, 10279),
	}
	want = math.MaxInt
	got  = b.Count()

	if want != got {
		t.Errorf("Didn't match: wanted %v, got %v", want, got)
	}
}

func TestParseButtons(t *testing.T) {
	inputStr := `
		Button A: X+94, Y+34
		Button B: X+22, Y+67
		Prize: X=8400, Y=5400

		Button A: X+69, Y+23
		Button B: X+27, Y+71
		Prize: X=18641, Y=10279
	`

	want := []cmd.Button{
		{
			A: util.MakeVector(94, 34),
			B: util.MakeVector(22, 67),
			Prize: util.MakePoint(8400, 5400),
		}, {
			A: util.MakeVector(69, 23),
			B: util.MakeVector(27, 71),
			Prize: util.MakePoint(18641, 10279),
		},
	}
	got, err := cmd.ParseButtons(inputStr)
	util.Unexpect(t, err)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Didn't match: wanted %+v, got %+v", want, got)
	}
}
