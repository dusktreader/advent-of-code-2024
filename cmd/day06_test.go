package cmd_test

import (
	"log/slog"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/go-test/deep"
)

// ....#.....
// .........#
// ..........
// ..#.......
// .......#..
// ..........
// .#..^.....
// ........#.
// #.........
// ......#...

func Unexpect(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}
}

func TestParseLabMap(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	txt := `
		.#
		#<
	`
	got, err := cmd.ParseLabMap(txt)
	Unexpect(t, err)

	wantGrid, err := util.MakeGrid(
		util.MakeSize(2, 2),
		[]cmd.Cell{
			{TreadCt: 0, HasObs: false},
			{TreadCt: 0, HasObs: true},
			{TreadCt: 0, HasObs: true},
			{TreadCt: 0, HasObs: false},
		},
	)
	Unexpect(t, err)
	want := cmd.LabMap{
		GuardPos: util.MakePoint(1, 1),
		GuardDir: util.MakeVector(0, -1),
		Grid:     wantGrid,
	}
	if diff := deep.Equal(&want, got); diff != nil {
		t.Error(diff)
	}
}

