package cmd_test

import (
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func TestParseAntMap(t *testing.T) {
	txt := `
		a..0
		.aX.
		.Xa.
		0..a
	`
	got, err := cmd.ParseAntMap(txt)
	util.Unexpect(t, err)

	want, err := cmd.MakeAntMap(4, 4)
	util.Unexpect(t, err)
	want.Ants.Add(
		'a',
		util.MakePoint(0, 0),
		util.MakePoint(1, 1),
		util.MakePoint(2, 2),
		util.MakePoint(3, 3),
	)
	want.Ants.Add(
		'X',
		util.MakePoint(1, 2),
		util.MakePoint(2, 1),
	)
	want.Ants.Add(
		'0',
		util.MakePoint(3, 0),
		util.MakePoint(0, 3),
	)

	if !want.Eq(got) {
		t.Errorf("Parsed antenna map didn't match:\n\nwant:\n%v\n\ngot:\n%v", want, got)
	}
}

func TestFindANodes(t *testing.T) {
	sz, err := util.MakeSize(100, 100)
	util.Unexpect(t, err)

	pts := util.MakeSet(
		util.MakePoint(2, 1),
		util.MakePoint(4, 2),
	)
	want := util.MakeSet(
		util.MakePoint(0, 0),
		util.MakePoint(6, 3),
	)
	got := cmd.FindAnsSimp(sz, pts)
	if !want.Eq(got) {
		t.Errorf("Failed: want %v, got %v", want, got)
	}

	pts = util.MakeSet(
		util.MakePoint(2, 1),
		util.MakePoint(4, 2),
		util.MakePoint(3, 2),
	)
	want = util.MakeSet(
		util.MakePoint(0, 0),
		util.MakePoint(6, 3),
		util.MakePoint(1, 0),
		util.MakePoint(4, 3),
		util.MakePoint(2, 2),
		util.MakePoint(5, 2),
	)
	got = cmd.FindAnsSimp(sz, pts)
	if !want.Eq(got) {
		t.Errorf("Failed: want %v, got %v", want, got)
	}
}
