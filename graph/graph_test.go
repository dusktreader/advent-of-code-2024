package graph_test

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/graph"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func pathsEq[T comparable](l *[][]T, r *[][]T) bool {
	for _, lp := range *l {
		in := false
		for _, rp := range *r {
			if reflect.DeepEqual(lp, rp) {
				in = true
				break
			}
		}
		if !in {
			return false
		}
	}
	return true
}

func TestPaths(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	g := graph.MakeDigraph(
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)
	slog.Debug("EDGES", "e", g.Edges())

	want := [][]int{
		{1, 3},
		{1, 2, 3},
		{1, 5, 2, 3},
	}
	got := g.Paths(1, 3)

	if !pathsEq(&want, &got) {
		t.Errorf("Failed: want %+v, got %+v", want, got)
	}
}
