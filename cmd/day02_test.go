package cmd_test

import (
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
)

func TestParseReportSuccess(t *testing.T) {
	inputStr := `
		1  2  3  4  5
		6  7  8  9  10
		11 12 13 14 15
		16 17 18 19 20
		21 22 23 24 25
		26 27 28 29 30
	`

	want := [][]int{
		{1,  2,  3,  4,  5 },
		{6,  7,  8,  9,  10},
		{11, 12, 13, 14, 15},
		{16, 17, 18, 19, 20},
		{21, 22, 23, 24, 25},
		{26, 27, 28, 29, 30},
	}
	got, err := cmd.ParseReport(inputStr)

	if err != nil {
		t.Fatalf(`ParseReport failed with error: %#q`, err)

	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`ParseReport failed on values: want %+v, got %+v`, want, got)
	}
}

func TestParseReportFail(t *testing.T) {
	inputStr := `
		1  2  three  4  5
		6  7  8      9  10
	`
	_, err := cmd.ParseReport(inputStr)
	if err == nil {
		t.Fatalf(`ParseReport did not fail as expected`)

	} else if !strings.Contains(err.Error(), "Failed to convert a token") {
		t.Fatalf(`ParseReport did not report token conversion failure`)
	}
}

func TestIsSafe(t *testing.T) {
	reports := [][]int{
		{7, 6, 4, 2, 1},              // Safe
		{1, 2, 7, 8, 9},              // Unsafe: |7 - 2| > 3,     Unsafe: |8 - 2| > 3
		{9, 7, 6, 2, 1},              // Unsafe: |6 - 2| > 3,     Unsafe: |6 - 1| > 3
		{1, 3, 2, 4, 5},              // Unsafe: Not monotonic,   Safe:   removes 3
		{8, 6, 4, 4, 1},              // Unsafe: |4 - 4| < 1,     Safe:   removes second 4
		{1, 3, 6, 7, 9},              // Safe
		{2, 1},                       // Safe
		{1, 3, 5, 7, 9},              // Safe
		{1, 1},                       // Unsafe: |1 - 1| < 1,     Safe   dampened: removes second 1
		{1, 1, 1},                    // Unsafe: |1 - 1| < 1,     Unsafe dampened: |1 - 1| < 1
		{9, 6, 3},                    // Safe
		{1, 5, 9, 13},                // Unsafe: |5 - 1| > 3,     Unsafe dampened: |9 - 5| > 3
		{5, 7, 6},                    // Unsafe: Not monotonic,   Safe   dampened: removes 7
		{7, 5, 6, 8},                 // Unsafe: Not monotonic,   Safe   dampened: removes 7
		{9, 8, 6, 3},                 // Safe
		{-9, -6, -3, 0, 3, 6, 9},     // Safe
		{6, 4, 5, 7},                 // Unsafe: Not monotonic,   Safe   dampened: removes 6
		{99, 98, 97, 95, 93, 91, 92}, // Unsafe: Not monotonic,   Safe   dampened: removes 91
		{100, 5, 6, 8},               // Unsafe: |5 - 100| > 3,   Safe   dampened: removes 100
		{8, 5, 6, 8},                 // Unsafe: Not monotonic,   Safe   dampened: removes 8
		{8, 5, 6, 4},                 // Unsafe: Not monotonic,   Safe   dampened: removes 5
		{8, 6, 9, 11},                // Unsafe: Not monotonic,   Safe   dampened: removes 8 (or 6)
		{8, 100, 9, 11},              // Unsafe: |100 - 8| > 3,   Safe   dampened: removes 100
		{8, 5, 6, 7, 8},              // Unsafe: Not monotonic,   Safe   dampened: removes 8
		{6, 100, 99, 98},             // Unsafe: |100 - 6| > 3,   Safe   dampened: removes 6
		{6, 5, 100, 99, 98},          // Unsafe: |100 - 5| > 3,   Unsafe dampened: not monotonic
	}

	wants := [][]bool{
		{true,  true},
		{false, false},
		{false, false},
		{false, true},
		{false, true},
		{true,  true},
		{true,  true},
		{true,  true},
		{false, true},
		{false, false},
		{true,  true},
		{false, false},
		{false, true},
		{false, true},
		{true,  true},
		{true,  true},
		{false, true},
		{false, true},
		{false, true},
		{false, true},
		{false, true},
		{false, true},
		{false, true},
		{false, true},
		{false, true},
		{false, false},
	}
	for i := 0; i < len(reports); i++ {
		wantUndampened := wants[i][0]
		gotUndampened := cmd.IsSafe(reports[i], 3, 0)
		if gotUndampened != wantUndampened {
			t.Errorf(`IsSafe gave a bad answer for report %+v (undampened). Wanted %v, got %v`, reports[i], wantUndampened, gotUndampened)
		}

		wantDampened := wants[i][1]
		gotDampened := cmd.IsSafe(reports[i], 3, 1)
		if gotDampened != wantDampened {
			t.Errorf(`IsSafe gave a bad answer for report %+v (dampened). Wanted %v, got %v`, reports[i], wantDampened, gotDampened)
		}
	}
}

func TestCountSafe(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	reports := [][]int{
		{7, 6, 4, 2, 1},              // Safe
		{1, 2, 7, 8, 9},              // Unsafe: |7 - 2| > 3,     Unsafe: |8 - 2| > 3
		{9, 7, 6, 2, 1},              // Unsafe: |6 - 2| > 3,     Unsafe: |6 - 1| > 3
		{1, 3, 2, 4, 5},              // Unsafe: Not monotonic,   Safe:   removes 3
		{8, 6, 4, 4, 1},              // Unsafe: |4 - 4| < 1,     Safe:   removes second 4
		{1, 3, 6, 7, 9},              // Safe
		{2, 1},                       // Safe
		{1, 3, 5, 7, 9},              // Safe
		{1, 1},                       // Unsafe: |1 - 1| < 1,     Safe   dampened: removes second 1
		{1, 1, 1},                    // Unsafe: |1 - 1| < 1,     Unsafe dampened: |1 - 1| < 1
		{9, 6, 3},                    // Safe
		{1, 5, 9, 13},                // Unsafe: |5 - 1| > 3,     Unsafe dampened: |9 - 5| > 3
		{5, 7, 6},                    // Unsafe: Not monotonic,   Safe   dampened: removes 7
		{7, 5, 6, 8},                 // Unsafe: Not monotonic,   Safe   dampened: removes 7
		{9, 8, 6, 3},                 // Safe
		{-9, -6, -3, 0, 3, 6, 9},     // Safe
		{6, 4, 5, 7},                 // Unsafe: Not monotonic,   Safe   dampened: removes 6
		{99, 98, 97, 95, 93, 91, 92}, // Unsafe: Not monotonic,   Safe   dampened: removes 91
		{100, 5, 6, 8},               // Unsafe: |5 - 100| > 3,   Safe   dampened: removes 100
		{8, 5, 6, 8},                 // Unsafe: Not monotonic,   Safe   dampened: removes 8
		{8, 5, 6, 4},                 // Unsafe: Not monotonic,   Safe   dampened: removes 5
		{8, 6, 9, 11},                // Unsafe: Not monotonic,   Safe   dampened: removes 8 (or 6)
		{8, 100, 9, 11},              // Unsafe: |100 - 8| > 3,   Safe   dampened: removes 100
		{8, 5, 6, 7, 8},              // Unsafe: Not monotonic,   Safe   dampened: removes 8
		{6, 100, 99, 98},             // Unsafe: |100 - 6| > 3,   Safe   dampened: removes 6
		{6, 5, 100, 99, 98},          // Unsafe: |100 - 5| > 3,   Unsafe dampened: not monotonic
	}

	wants := []int{7, 21}

	wantUndampened := wants[0]
	gotUndampened := cmd.CountSafe(reports, 3, 0)
	if gotUndampened != wantUndampened {
		t.Errorf(`CountSafe gave a bad answer (undampened). Wanted %v, got %v`, wantUndampened, gotUndampened)
	}

	wantDampened := wants[1]
	gotDampened := cmd.CountSafe(reports, 3, 1)
	if gotDampened != wantDampened {
		t.Errorf(`CountSafe gave a bad answer (dampened). Wanted %v, got %v`, wantDampened, gotDampened)
	}
}
