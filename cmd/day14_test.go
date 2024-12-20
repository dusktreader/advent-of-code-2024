package cmd_test

import (
	"reflect"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func TestParseRobots(t *testing.T) {
	inputStr := `
		p=0,4 v=3,-3
		p=6,3 v=-1,-3
		p=10,3 v=-1,2
	`

	want := cmd.PissSpace{
		Size: util.Size{W: 0, H: 0},
		Robots: []cmd.Robot{
			{Pos: util.MakePoint(4, 0), Vel: util.MakeVector(-3, 3)},
			{Pos: util.MakePoint(3, 6), Vel: util.MakeVector(-3, -1)},
			{Pos: util.MakePoint(3, 10), Vel: util.MakeVector(2, -1)},
		},
	}
	got, err := cmd.ParseRobots(inputStr)
	util.Unexpect(t, err)

	if !reflect.DeepEqual(want, *got) {
		t.Fatalf("Didn't match: wanted %+v, got %+v", want, got)
	}
}
