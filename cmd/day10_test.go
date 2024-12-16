package cmd_test

import (
	"reflect"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/graph"
	"github.com/dusktreader/advent-of-code-2024/util"
)

type ppair = util.Pair[util.Point]

func TestParseTopoMap(t *testing.T) {
	inputStr := `
		...0...
		...1...
		...2...
		6543456
		7.....7
		8.....8
		9.....9
	`
	want := graph.MakeGraph(
		true,
		ppair{Left: util.MakePoint(0, 3), Right: util.MakePoint(1, 3)},
		ppair{Left: util.MakePoint(1, 3), Right: util.MakePoint(2, 3)},
		ppair{Left: util.MakePoint(2, 3), Right: util.MakePoint(3, 3)},

		ppair{Left: util.MakePoint(3, 3), Right: util.MakePoint(3, 2)},
		ppair{Left: util.MakePoint(3, 2), Right: util.MakePoint(3, 1)},
		ppair{Left: util.MakePoint(3, 1), Right: util.MakePoint(3, 0)},

		ppair{Left: util.MakePoint(3, 3), Right: util.MakePoint(3, 4)},
		ppair{Left: util.MakePoint(3, 4), Right: util.MakePoint(3, 5)},
		ppair{Left: util.MakePoint(3, 5), Right: util.MakePoint(3, 6)},

		ppair{Left: util.MakePoint(3, 0), Right: util.MakePoint(4, 0)},
		ppair{Left: util.MakePoint(4, 0), Right: util.MakePoint(5, 0)},
		ppair{Left: util.MakePoint(5, 0), Right: util.MakePoint(6, 0)},

		ppair{Left: util.MakePoint(3, 6), Right: util.MakePoint(4, 6)},
		ppair{Left: util.MakePoint(4, 6), Right: util.MakePoint(5, 6)},
		ppair{Left: util.MakePoint(5, 6), Right: util.MakePoint(6, 6)},
	)
	got, err := cmd.ParseTopoMap(inputStr)
	util.Unexpect(t, err)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Failed: want %+v, got %+v", want, got)
	}
}

