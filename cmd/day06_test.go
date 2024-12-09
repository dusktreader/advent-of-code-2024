package cmd_test

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/util"
)

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

	want, err := cmd.MakeLabMap(2, 2)
	Unexpect(t, err)

	want.GuardPos = util.MakePoint(1, 1)
	want.GuardDir = util.MakeVector(0, -1)
	want.Obstr    = util.MakeSet([]util.Point{
		{I: 0, J: 1},
		{I: 1, J: 0},
	}...)
	want.Visits   = util.MakeSetMap[util.Point, util.Vector]()
	want.Visits.Add(util.MakePoint(1, 1), util.MakeVector(0, -1))
	if !want.Eq(got) {
		t.Errorf("Parsed lab map didn't match:\n\nwant:\n%v\n\ngot:\n%v", want, got)
	}
}

func TestWalk(t *testing.T) {
	lm, err := cmd.ParseLabMap(`
		...
		#.<
		...
	`)
	Unexpect(t, err)

	wantLm, err := cmd.ParseLabMap(`
		...
		#<.
		...
	`)
	Unexpect(t, err)
	wantLm.Visits   = util.MakeSetMap[util.Point, util.Vector]()
	wantLm.Visits.Add(util.MakePoint(1, 2), util.MakeVector(0, -1))
	wantLm.Visits.Add(util.MakePoint(1, 1), util.MakeVector(0, -1))
	wantIn := true

	gotIn, err := lm.Walk()
	if gotIn != wantIn {
		t.Errorf("Mismatch 'in' status: want %v, got %v", wantIn, gotIn)
	}

	if !lm.Eq(wantLm) {
		t.Fatalf("Maps didn't match:\n\nwant:\n%v\n\ngot:\n%v\n\nwant visits: %+v\ngot visits: %+v", wantLm, lm, wantLm.Visits, lm.Visits)
	}

	wantLm, err = cmd.ParseLabMap(`
		...
		#^.
		...
	`)
	Unexpect(t, err)
	wantLm.Visits   = util.MakeSetMap[util.Point, util.Vector]()
	wantLm.Visits.Add(util.MakePoint(1, 2), util.MakeVector(0, -1))
	wantLm.Visits.Add(util.MakePoint(1, 1), util.MakeVector(0, -1), util.MakeVector(-1, 0))
	wantIn = true

	gotIn, err = lm.Walk()
	if gotIn != wantIn {
		t.Errorf("Mismatch 'in' status: want %v, got %v", wantIn, gotIn)
	}

	if !lm.Eq(wantLm) {
		t.Fatalf("Maps didn't match:\n\nwant:\n%v\n\ngot:\n%v\n\nwant visits: %+v\ngot visits: %+v", wantLm, lm, wantLm.Visits, lm.Visits)
	}

	wantLm, err = cmd.ParseLabMap(`
		.^.
		#..
		...
	`)
	Unexpect(t, err)
	wantLm.Visits   = util.MakeSetMap[util.Point, util.Vector]()
	wantLm.Visits.Add(util.MakePoint(1, 2), util.MakeVector(0, -1))
	wantLm.Visits.Add(util.MakePoint(1, 1), util.MakeVector(0, -1), util.MakeVector(-1, 0))
	wantLm.Visits.Add(util.MakePoint(0, 1), util.MakeVector(-1, 0))
	wantIn = true

	gotIn, err = lm.Walk()
	if gotIn != wantIn {
		t.Errorf("Mismatch 'in' status: want %v, got %v", wantIn, gotIn)
	}

	if !lm.Eq(wantLm) {
		t.Fatalf("Maps didn't match:\n\nwant:\n%v\n\ngot:\n%v\n\nwant visits: %+v\ngot visits: %+v", wantLm, lm, wantLm.Visits, lm.Visits)
	}

	wantLm, err = cmd.ParseLabMap(`
		...
		#..
		...
	`)
	Unexpect(t, err)
	wantLm.Visits   = util.MakeSetMap[util.Point, util.Vector]()
	wantLm.Visits.Add(util.MakePoint(1, 2), util.MakeVector(0, -1))
	wantLm.Visits.Add(util.MakePoint(1, 1), util.MakeVector(0, -1), util.MakeVector(-1, 0))
	wantLm.Visits.Add(util.MakePoint(0, 1), util.MakeVector(-1, 0))
	wantLm.GuardPos = util.MakePoint(-1, 1)
	wantLm.GuardDir = util.MakeVector(-1, 0)
	wantIn = false

	gotIn, err = lm.Walk()
	if gotIn != wantIn {
		t.Errorf("Mismatch 'in' status: want %v, got %v", wantIn, gotIn)
	}

	if !lm.Eq(wantLm) {
		t.Fatalf("Maps didn't match:\n\nwant:\n%v\n\ngot:\n%v\n\nwant visits: %+v\ngot visits: %+v", wantLm, lm, wantLm.Visits, lm.Visits)
	}
}

func TestPatrol(t *testing.T) {
	lm, err := cmd.ParseLabMap(`
		...
		#.<
		...
	`)
	Unexpect(t, err)

	wantLm, err := cmd.ParseLabMap(`
		...
		#..
		...
	`)
	Unexpect(t, err)
	wantLm.Visits   = util.MakeSetMap[util.Point, util.Vector]()
	wantLm.Visits.Add(util.MakePoint(1, 2), util.MakeVector(0, -1))
	wantLm.Visits.Add(util.MakePoint(1, 1), util.MakeVector(0, -1), util.MakeVector(-1, 0))
	wantLm.Visits.Add(util.MakePoint(0, 1), util.MakeVector(-1, 0))
	wantLm.GuardPos = util.MakePoint(-1, 1)
	wantLm.GuardDir = util.MakeVector(-1, 0)
	wantIn := false

	lm.Patrol()
	gotIn := !lm.Size.Out(lm.GuardPos)
	if gotIn != wantIn {
		t.Errorf("Mismatch 'in' status: want %v, got %v", wantIn, gotIn)
	}

	if !lm.Eq(wantLm) {
		t.Fatalf("Maps didn't match:\n\nwant:\n%v\n\ngot:\n%v\n\nwant visits: %+v\ngot visits: %+v", wantLm, lm, wantLm.Visits, lm.Visits)
	}
}

func TestCountVisits(t *testing.T) {
	lm, err := cmd.ParseLabMap(`
		....#.....
		.........#
		..........
		..#.......
		.......#..
		..........
		.#..^.....
		........#.
		#.........
		......#...
	`)
	Unexpect(t, err)
	lm.Patrol()
	want := 41
	got  := lm.CountVisits()
	if want != got {
		t.Fatalf("Visit count didn't match: want: %v, got: %v, final map:\n%v", want, got, lm)
	}
}

func TestLoopify(t *testing.T) {
	lm, err := cmd.ParseLabMap(`
		....#.....
		.........#
		..........
		..#.......
		.......#..
		..........
		.#..^.....
		........#.
		#.........
		......#...
	`)
	Unexpect(t, err)
	err = lm.Loopify()
	Unexpect(t, err)
	got := lm.CountLoopers()
	want := 6
	if want != got {
		t.Fatalf("New obstruction count didn't match: want: %v, got: %v", want, got)
	}

	lm, err = cmd.ParseLabMap(`
		....#.....
		.........#
		..........
		..........
		........#.
		..........
		....^.....
		..........
		..........
		..........
	`)
	Unexpect(t, err)
	err = lm.Loopify()
	Unexpect(t, err)
	got = lm.CountLoopers()
	fmt.Printf("\n\n%v\n", lm)
	want = 1
	if want != got {
		t.Fatalf("New obstruction count didn't match: want: %v, got: %v", want, got)
	}
}
