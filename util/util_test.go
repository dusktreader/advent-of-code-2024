package util_test

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/util"
)

func TestCounter(t *testing.T) {
	c := util.MakeCounter[int]()
	if c.Get(1) != 0 {
		t.Errorf("counter mistakenly said count of 1 was not 0")
	}

	if c.Incr(1) != 1 {
		t.Errorf("counter mistakenly said count of 1 was not 1")
	}

	if c.Incr(1) != 2 {
		t.Errorf("counter mistakenly said count of 1 was not 2")
	}

	if c.Get(2) != 0 {
		t.Errorf("counter mistakenly said count of 2 was not 0")
	}
}

func TestSet(t *testing.T) {
	items := []int{1, 2, 3}
	set := util.MakeSet(items...)

	for _, item := range items {
		if !set.Has(item) {
			t.Errorf("set mistakenly said %v was not in %+v", item, set)
		}
	}
	if set.Has(4) {
		t.Errorf("set mistakenly said %v was not in %+v", 4, set)
	}

	set.Add(4)
	if !set.Has(4) {
		t.Errorf("set mistakenly said %v was not in %+v", 4, set)
	}

	set.Rem(4)
	if set.Has(4) {
		t.Errorf("set mistakenly said %v was in %+v", 4, set)
	}

	set.Rem(3)
	set.Rem(2)

	item := set.Pop()
	if item != 1 {
		t.Errorf("set mistakenly popped the wrong value %v from %+v", 1, set)
	}

	set.Pop()
	if !set.Empty() {
		t.Errorf("set mistakenly said it wasn't empty: %+v", set)
	}

}

func TestDagBasic(t *testing.T) {
	dag := util.MakeDag[int]()
	if dag.Nodes().Size() != 0 {
		t.Errorf("dag mistakenly said it had nodes: %+v", dag)
	}
	if dag.Edges().Size() != 0 {
		t.Errorf("dag mistakenly said it had edges: %+v", dag)
	}

	dag.AddEdge(1, 2)

	wantNodes := util.MakeSet(1, 2)
	gotNodes := dag.Nodes()
	if !wantNodes.Eq(gotNodes) {
		t.Errorf("dag produced the wrong nodes: wanted %+v, got %+v", wantNodes, gotNodes)
	}

	wantEdges := util.MakeSet(util.Pair[int]{1, 2})
	gotEdges  := dag.Edges()
	if !wantEdges.Eq(gotEdges) {
		t.Errorf("dag produced the wrong edges: wanted %+v, got %+v", wantEdges, gotEdges)
	}

	dag.AddEdge(1, 3)

	wantNodes = util.MakeSet(1, 2, 3)
	gotNodes = dag.Nodes()
	if !wantNodes.Eq(gotNodes) {
		t.Errorf("dag produced the wrong nodes: wanted %+v, got %+v", wantNodes, gotNodes)
	}

	wantEdges = util.MakeSet(util.Pair[int]{1, 2}, util.Pair[int]{1, 3})
	gotEdges  = dag.Edges()
	if !wantEdges.Eq(gotEdges) {
		t.Errorf("dag produced the wrong edges: wanted %+v, got %+v", wantEdges, gotEdges)
	}

	dag.AddEdge(3, 4)
	dag.AddEdge(2, 4)

	wantNodes = util.MakeSet(1)
	gotNodes = dag.Sources()
	if !wantNodes.Eq(gotNodes) {
		t.Errorf("dag produced the wrong sources: wanted %+v, got %+v", wantNodes, gotNodes)
	}

	wantNodes = util.MakeSet(4)
	gotNodes = dag.Sinks()
	if !wantNodes.Eq(gotNodes) {
		t.Errorf("dag produced the wrong terminals: wanted %+v, got %+v", wantNodes, gotNodes)
	}
}

func makeTestDag() (util.DAG[int]) {
	return util.MakeDag(
		util.Pair[int]{47, 53},
		util.Pair[int]{97, 13},
		util.Pair[int]{97, 61},
		util.Pair[int]{97, 47},
		util.Pair[int]{75, 29},
		util.Pair[int]{61, 13},
		util.Pair[int]{75, 53},
		util.Pair[int]{29, 13},
		util.Pair[int]{97, 29},
		util.Pair[int]{53, 29},
		util.Pair[int]{61, 53},
		util.Pair[int]{97, 53},
		util.Pair[int]{61, 29},
		util.Pair[int]{47, 13},
		util.Pair[int]{75, 47},
		util.Pair[int]{97, 75},
		util.Pair[int]{47, 61},
		util.Pair[int]{75, 61},
		util.Pair[int]{47, 29},
		util.Pair[int]{75, 13},
		util.Pair[int]{53, 13},
	)
}

func TestDagIsSorted(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	dag := makeTestDag()

	items := []int{75, 47, 61, 53, 29}
	want := true
	got, err := dag.IsSorted(items)
	if err != nil {
		t.Errorf("dag errored on IsSorted: %#v", err)
	} else if want != got {
		t.Errorf("dag got IsSorted wrong for %+v: wanted %v, got %v", items, want, got)
	}

	items = []int{97, 61, 53, 29, 13}
	want = true
	got, err = dag.IsSorted(items)
	if err != nil {
		t.Errorf("dag errored on IsSorted: %#v", err)
	} else if want != got {
		t.Errorf("dag got IsSorted wrong for %+v: wanted %v, got %v", items, want, got)
	}

	items = []int{75, 29, 13}
	want = true
	got, err = dag.IsSorted(items)
	if err != nil {
		t.Errorf("dag errored on IsSorted: %#v", err)
	} else if want != got {
		t.Errorf("dag got IsSorted wrong for %+v: wanted %v, got %v", items, want, got)
	}

	items = []int{75, 97, 47, 61, 53}
	want = false
	got, err = dag.IsSorted(items)
	if err != nil {
		t.Errorf("dag errored on IsSorted: %#v", err)
	} else if want != got {
		t.Errorf("dag got IsSorted wrong for %+v: wanted %v, got %v", items, want, got)
	}

	items = []int{61, 13, 29}
	want = false
	got, err = dag.IsSorted(items)
	if err != nil {
		t.Errorf("dag errored on IsSorted: %#v", err)
	} else if want != got {
		t.Errorf("dag got IsSorted wrong for %+v: wanted %v, got %v", items, want, got)
	}

	items = []int{97, 13, 75, 29, 47}
	want = false
	got, err = dag.IsSorted(items)
	if err != nil {
		t.Errorf("dag errored on IsSorted: %#v", err)
	} else if want != got {
		t.Errorf("dag got IsSorted wrong for %+v: wanted %v, got %v", items, want, got)
	}
}

func TestDagSort(t *testing.T) {
	dag := makeTestDag()

	items := []int{75, 47, 61, 53, 29}
	want := items
	got, err := dag.Sort(items)
	if err != nil {
		t.Fatalf("Unexpected error from Sort: %#v", err)
	} else if !reflect.DeepEqual(want, got) {
		t.Errorf("dag got Sort wrong for %+v: wanted %v, got %v", items, want, got)
	}

	items = []int{97, 61, 53, 29, 13}
	want = items
	got, err = dag.Sort(items)
	if err != nil {
		t.Fatalf("Unexpected error from Sort: %#v", err)
	} else if !reflect.DeepEqual(want, got) {
		t.Errorf("dag got Sort wrong for %+v: wanted %v, got %v", items, want, got)
	}

	items = []int{75, 29, 13}
	want = items
	got, err = dag.Sort(items)
	if err != nil {
		t.Fatalf("Unexpected error from Sort: %#v", err)
	} else if !reflect.DeepEqual(want, got) {
		t.Errorf("dag got Sort wrong for %+v: wanted %v, got %v", items, want, got)
	}

	items = []int{75, 97, 47, 61, 53}
	want = []int{97, 75, 47, 61, 53}
	got, err = dag.Sort(items)
	if err != nil {
		t.Fatalf("Unexpected error from Sort: %#v", err)
	} else if !reflect.DeepEqual(want, got) {
		t.Errorf("dag got Sort wrong for %+v: wanted %v, got %v", items, want, got)
	}

	items = []int{61, 13, 29}
	want = []int{61, 29, 13}
	got, err = dag.Sort(items)
	if err != nil {
		t.Fatalf("Unexpected error from Sort: %#v", err)
	} else if !reflect.DeepEqual(want, got) {
		t.Errorf("dag got Sort wrong for %+v: wanted %v, got %v", items, want, got)
	}

	items = []int{97, 13, 75, 29, 47}
	want = []int{97, 75, 47, 29, 13}
	got, err = dag.Sort(items)
	if err != nil {
		t.Fatalf("Unexpected error from Sort: %#v", err)
	} else if !reflect.DeepEqual(want, got) {
		t.Errorf("dag got Sort wrong for %+v: wanted %v, got %v", items, want, got)
	}
}
