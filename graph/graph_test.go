package graph_test

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/graph"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func TestOutN(t *testing.T) {
	g := graph.MakeGraph(
		true,
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)

	want := util.MakeSet(2, 3, 5)
	got  := g.OutN(1)
	if !want.Eq(got) {
		t.Errorf("Wrong out-neighbors: wanted %v, got %v", want, got)
	}

	want = util.MakeSet[int]()
	got  = g.OutN(3)
	if !want.Eq(got) {
		t.Errorf("Wrong out-neighbors: wanted %v, got %v", want, got)
	}

	want = util.MakeSet(2, 6)
	got  = g.OutN(5)
	if !want.Eq(got) {
		t.Errorf("Wrong out-neighbors: wanted %v, got %v", want, got)
	}
}

func TestInN(t *testing.T) {
	g := graph.MakeGraph(
		true,
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)

	want := util.MakeSet[int]()
	got  := g.InN(1)
	if !want.Eq(got) {
		t.Errorf("Wrong in-neighbors: wanted %v, got %v", want, got)
	}

	want = util.MakeSet(1, 2)
	got  = g.InN(3)
	if !want.Eq(got) {
		t.Errorf("Wrong in-neighbors: wanted %v, got %v", want, got)
	}

	want = util.MakeSet(1, 4)
	got  = g.InN(5)
	if !want.Eq(got) {
		t.Errorf("Wrong in-neighbors: wanted %v, got %v", want, got)
	}
}

func TestNbors(t *testing.T) {
	g := graph.MakeGraph(
		true,
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)

	want := util.MakeSet(2, 3, 5)
	got  := g.Nbors(1)
	if !want.Eq(got) {
		t.Errorf("Wrong neighbors: wanted %v, got %v", want, got)
	}

	want = util.MakeSet(1, 2)
	got  = g.Nbors(3)
	if !want.Eq(got) {
		t.Errorf("Wrong neighbors: wanted %v, got %v", want, got)
	}

	want = util.MakeSet(1, 4, 2, 6)
	got  = g.Nbors(5)
	if !want.Eq(got) {
		t.Errorf("Wrong neighbors: wanted %v, got %v", want, got)
	}

	want = util.MakeSet(5)
	got  = g.Nbors(6)
	if !want.Eq(got) {
		t.Errorf("Wrong neighbors: wanted %v, got %v", want, got)
	}
}

func TestEdgesDirected(t *testing.T) {
	g := graph.MakeGraph(
		true,
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)

	want := util.MakeSet(
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)
	got := g.Edges()

	for !want.Empty() {
		p := want.Pop()
		if !got.Has(p) {
			t.Errorf("Missing expected edge: %v", p)
		}

		rp := p.Rev()
		if got.Has(rp) {
			t.Errorf("Includes unexpected edge: %v", rp)
		}

		got.Rem(p)
	}
	if !got.Empty() {
		t.Errorf("Unexpected edges: %v", got)
	}
}

func TestEdgesUndirected(t *testing.T) {
	g := graph.MakeGraph(
		false,
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)

	want := util.MakeSet(
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)
	got := g.Edges()

	for !want.Empty() {
		p := want.Pop()
		if got.Has(p) {
			got.Rem(p)
		} else {
			rp := p.Rev()
			if got.Has(rp) {
				got.Rem(rp)
			} else {
				t.Errorf("Missing expected edge: %v", p)
			}
		}
	}
	if !got.Empty() {
		t.Errorf("Unexpected edges: %v", got)
	}
}

func TestEqDirected(t *testing.T) {
	g := graph.MakeGraph[int](true)
	g.Add(1)
	g.Add(2)
	g.Add(3)
	g.Add(4)
	g.Add(5)
	g.AddEdge(5, 6)
	g.AddEdge(2, 3)
	g.AddEdge(5, 2)
	g.AddEdge(4, 5)
	g.AddEdge(1, 5)
	g.AddEdge(1, 3)
	g.AddEdge(1, 2)

	og := graph.MakeGraph(
		true,
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)

	if !g.Eq(og) {
		t.Fatalf("Didn't match: %v != %v", g, og)
	}
}

func TestEqUndirected(t *testing.T) {
	g := graph.MakeGraph(
		false,
		util.MakePair(2, 1),
		util.MakePair(3, 1),
		util.MakePair(5, 1),
		util.MakePair(5, 4),
		util.MakePair(2, 5),
		util.MakePair(3, 2),
		util.MakePair(6, 5),
	)

	og := graph.MakeGraph(
		false,
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)

	if !g.Eq(og) {
		t.Fatalf("Didn't match: %v != %v", g, og)
	}
}

func TestSourcesDirected(t *testing.T) {
	g := graph.MakeGraph(
		true,
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)

	want := util.MakeSet(1, 4)
	got  := g.Sources()
	if !want.Eq(got) {
		t.Errorf("Wrong sources: wanted %v, got %v", want, got)
	}
}

func TestSinksDirected(t *testing.T) {
	g := graph.MakeGraph(
		true,
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)

	want := util.MakeSet(3, 6)
	got  := g.Sinks()
	if !want.Eq(got) {
		t.Errorf("Wrong sinks: wanted %v, got %v", want, got)
	}
}

func TestHasCyle(t *testing.T) {
	g := graph.MakeGraph(
		true,
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)

	want := false
	got  := g.HasCycle()
	if want != got {
		t.Errorf("Wrong has-cycle: wanted %v, got %v", want, got)
	}

	g.AddEdge(6, 4)

	want = true
	got  = g.HasCycle()
	if want != got {
		t.Errorf("Wrong has-cycle: wanted %v, got %v", want, got)
	}
}

func TestCnxComp(t *testing.T) {
	g := graph.MakeGraph(
		false,
		util.MakePair(1, 2),
		util.MakePair(1, 6),
		util.MakePair(2, 6),
		util.MakePair(4, 5),
	)
	g.Add(3)

	want := make([]graph.Graph[int], 3)
	want[0] = *graph.MakeGraph(
		false,
		util.MakePair(1, 2),
		util.MakePair(1, 6),
		util.MakePair(2, 6),
	)
	want[1] = *graph.MakeGraph(
		false,
		util.MakePair(4, 5),
	)
	want[2] = *graph.MakeGraph[int](
		false,
	)
	want[2].Add(3)

	got := g.CnxComp()
	if len(got) != len(want) {
		t.Fatalf("Got different item count")
	}

	OUTER: for _, wt := range want {
		for _, gt := range got {
			if wt.Eq(&gt) {
				continue OUTER
			}
		}
		t.Errorf("Didn't find expected subgraph: %v", wt)
	}
}

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

func TestPathsDirected(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	g := graph.MakeGraph(
		true,
		util.MakePair(1, 2),
		util.MakePair(1, 3),
		util.MakePair(1, 5),
		util.MakePair(4, 5),
		util.MakePair(5, 2),
		util.MakePair(2, 3),
		util.MakePair(5, 6),
	)

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
