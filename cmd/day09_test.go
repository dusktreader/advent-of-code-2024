package cmd_test

import (
	"reflect"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func TestParseFiles(t *testing.T) {
	inputStr := "12345"
	want := []int{1, 2, 3, 4, 5}
	got, err := cmd.ParseFiles(inputStr)
	util.Unexpect(t, err)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Failed: want %+v, got %+v", want, got)
	}

	inputStr = "1"
	want = []int{1}
	got, err = cmd.ParseFiles(inputStr)
	util.Unexpect(t, err)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Failed: want %+v, got %+v", want, got)
	}

	inputStr = "123a456"
	want = []int{1}
	got, err = cmd.ParseFiles(inputStr)
	if err == nil {
		t.Errorf("Parse files didn't fail as expected")
	}
}

func TestCompactFiles(t *testing.T) {
	files, err := cmd.ParseFiles("12345")
	util.Unexpect(t, err)

	want := []int{0, 2, 2, 1, 1, 1, 2, 2, 2}
	got := cmd.CompactFilesDense(files)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Failed: want %+v, got %+v", want, got)
	}

	files, err = cmd.ParseFiles("10305")
	util.Unexpect(t, err)

	want = []int{0, 1, 1, 1, 2, 2, 2, 2, 2}
	got = cmd.CompactFilesDense(files)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Failed: want %+v, got %+v", want, got)
	}

	files, err = cmd.ParseFiles("2333133121414131402")
	util.Unexpect(t, err)

	want = []int{0, 0, 9, 9, 8, 1, 1, 1, 8, 8, 8, 2, 7, 7, 7, 3, 3, 3, 6, 4, 4, 6, 5, 5, 5, 5, 6, 6}
	got = cmd.CompactFilesDense(files)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Failed:\n\nwant %+v\n\ngot  %+v", want, got)
	}
}

func TestExpandSlots(t *testing.T) {
	files, err := cmd.ParseFiles("12345")
	util.Unexpect(t, err)
	e := cmd.EMPTY

	want := []cmd.FileSlot{
		{Id: 0, Ct: 1},
		{Id: e, Ct: 2},
		{Id: 1, Ct: 3},
		{Id: e, Ct: 4},
		{Id: 2, Ct: 5},
	}
	got := cmd.ExpandSlots(files)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Failed:\n\nwant\n%+v,\n\ngot\n%+v", want, got)
	}
}

func TestCompactFilesSparse(t *testing.T) {
	files, err := cmd.ParseFiles("2333133121414131402")
	util.Unexpect(t, err)

	slots := cmd.ExpandSlots(files)

	want := "00992111777.44.333....5555.6666.....8888.."
	got  := cmd.PrintSlots(cmd.CompactFilesSparse(slots))

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Failed:\n\nwant\n%+v,\n\ngot\n%+v", want, got)
	}
}

func TestComputeChecksumSparse(t *testing.T) {
	files, err := cmd.ParseFiles("2333133121414131402")
	util.Unexpect(t, err)
	slots := cmd.ExpandSlots(files)
	slots = cmd.CompactFilesSparse(slots)

	want := 2858
	got  := cmd.ComputeChecksumSparse(slots)
	if want != got {
		t.Fatalf("Failed: want %v, got %v", want, got)
	}
}
