package cmd_test

import (
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
)

func TestRunifySuccess(t *testing.T) {
	inputStr := `
		foo
		bar
		baz
	`

	want := [][]rune{
		{'f', 'o', 'o'},
		{'b', 'a', 'r'},
		{'b', 'a', 'z'},
	}
	got, err := cmd.Runify(inputStr)
	if err != nil {
		t.Fatalf(`Runify had an unexpected error: %#v`, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`Runify gave a bad answer. Wanted %+v, got %+v`, want, got)
	}
}

func TestRunifyFail(t *testing.T) {
	inputStr := `
		foo

		baz
	`
	_, err := cmd.Runify(inputStr)
	if err == nil {
		t.Fatalf(`Runify did not fail as expected`)
	} else if !strings.Contains(err.Error(), "was empty") {
		t.Fatalf(`Runify did not report empty line`)
	}

	inputStr = `
		foo
		barr
		baz
	`
	_, err = cmd.Runify(inputStr)
	if err == nil {
		t.Fatalf(`Runify did not fail as expected`)
	} else if !strings.Contains(err.Error(), "didn't match") {
		t.Fatalf(`Runify did not non-matching line length`)
	}
}

func TestCountMatches(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	board, err := cmd.Runify(`
		MMMSXXMASM
		MSAMXMSMSA
		AMXSXMAAMM
		MSAMASMSMX
		XMASAMXAMM
		XXAMMXXAMA
		SMSMSASXSS
		SAXAMASAAA
		MAMMMXMMMM
		MXMXAXMASX
	`)
	if err != nil {
		t.Fatalf(`Runify failed unexpectedly: %#v`, err)
	}

	want := 18
	got := cmd.CountMatches([]rune("XMAS"), board)

	if want != got {
		t.Fatalf(`CountMatches gave a bad answer. Wanted %v, got %v`, want, got)
	}
}

func TestRotatePatch(t *testing.T) {
	patch, err := cmd.Runify(`
		..X..
		.XXX.
		X.X.X
		..X..
		..X..
	`)
	if err != nil {
		t.Fatalf(`Runify failed unexpectedly: %#v`, err)
	}

	want, err := cmd.Runify(`
		..X..
		...X.
		XXXXX
		...X.
		..X..
	`)
	if err != nil {
		t.Fatalf(`Runify failed unexpectedly: %#v`, err)
	}

	got := cmd.RotatePatch(patch)
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("RotatePatch gave a bad answer. Wanted:\n%v\ngot:\n%v", cmd.PrettyPrintPatch(want), cmd.PrettyPrintPatch(got))
	}
}

func TestMakePatch(t *testing.T) {
	want, err := cmd.Runify(`
		.....
		.....
		.....
		.....
		.....
	`)
	if err != nil {
		t.Fatalf(`Runify failed unexpectedly: %#v`, err)
	}

	got := cmd.MakePatch(5)
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("MakePatch gave a bad answer. Wanted:\n%v\ngot:\n%v", cmd.PrettyPrintPatch(want), cmd.PrettyPrintPatch(got))
	}
}

func TestMakeCrossPatch(t *testing.T) {
	want, err := cmd.Runify(`
		C...C
		.R.R.
		..O..
		.S.S.
		S...S
	`)
	if err != nil {
		t.Fatalf(`Runify failed unexpectedly: %#v`, err)
	}

	got := cmd.MakeCrossPatch([]rune("CROSS"))
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("MakeCrossPatch gave a bad answer. Wanted:\n%v\ngot:\n%v", cmd.PrettyPrintPatch(want), cmd.PrettyPrintPatch(got))
	}
}

func TestPatchMatch(t *testing.T) {
	patch := cmd.MakeCrossPatch([]rune("CROSS"))

	board, err := cmd.Runify(`
		.........
		.........
		..C...C..
		...R.R...
		....O....
		...S.S...
		..S...S..
		.........
		.........
	`)
	if err != nil {
		t.Fatalf(`Runify failed unexpectedly: %#v`, err)
	}
	cmd.PatchRandFill(board)

	want := true
	got := cmd.PatchMatch(patch, board, cmd.Point{4, 4})
	if want != got {
		t.Fatalf("PatchMatch gave a bad answer. Wanted: %v got: %v", want, got)
	}

	want = false
	got = cmd.PatchMatch(patch, board, cmd.Point{3, 3})
	if want != got {
		t.Fatalf("PatchMatch gave a bad answer. Wanted: %v got: %v", want, got)
	}
}

func TestCountCrossWords(t *testing.T) {
	board, err := cmd.Runify(`
		.M.S......
		..A..MSMS.
		.M.S.MAA..
		..A.ASMSM.
		.M.S.M....
		..........
		S.S.S.S.S.
		.A.A.A.A..
		M.M.M.M.M.
		..........
	`)
	if err != nil {
		t.Fatalf(`Runify failed unexpectedly: %#v`, err)
	}
	cmd.PatchRandFill(board)

	want := 9
	got := cmd.CountCrossWords([]rune("MAS"), board)
	if want != got {
		t.Fatalf("CountCrossWords gave a bad answer. Wanted: %v got: %v", want, got)
	}
}
