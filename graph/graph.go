package graph

import (
	"fmt"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
)

// Need to add tests for this!!

type Digraph [T comparable] struct {
	nodes       util.Set[T]
	sources     util.Set[T]
	sinks       util.Set[T]
	sortedNodes []T
	oEdges      util.SetMap[T, T]
	iEdges      util.SetMap[T, T]
}

func MakeDigraph[T comparable](pairs ...util.Pair[T]) *Digraph[T] {
	dg := Digraph[T]{
		oEdges:  util.MakeSetMap[T, T](),
		iEdges:  util.MakeSetMap[T, T](),
		nodes:   util.MakeSet[T](),
		sources: util.MakeSet[T](),
		sinks:   util.MakeSet[T](),
	}
	for _, p := range pairs {
		dg.AddEdge(p.Left, p.Right)
	}
	return &dg
}

func (dg *Digraph[T]) Eq(og *Digraph[T]) bool {
	return util.All(
		dg.nodes.Eq(og.nodes),
		dg.Edges().Eq(og.Edges()),
	)
}

func (dg *Digraph[T]) Clone() *Digraph[T] {
	ndg := Digraph[T]{
		nodes: dg.nodes.Clone(),
		sources: dg.sources.Clone(),
		sinks: dg.sinks.Clone(),
		oEdges: dg.oEdges.Clone(),
		iEdges: dg.iEdges.Clone(),
	}
	return &ndg
}

func (dg *Digraph[T]) Add(n T) {
	dg.nodes.Add(n)
	dg.sortedNodes = nil
	dg.sources.Clear()
	dg.sinks.Clear()
}

func (dg Digraph[T]) String() (string) {
	nodes := []string{}
	for n := range dg.nodes.Iter() {
		nodes = append(nodes, fmt.Sprintf("%v", n))
	}
	edges := []string{}
	for l, r := range dg.oEdges.Iter() {
		edges = append(edges, fmt.Sprintf("%v->%v", l, r))
	}
	return "(" + strings.Join(nodes, ", ") + "):{" + strings.Join(edges, ", ") + "}"
}

func (dg *Digraph[T]) Dot() (string) {
	edges := dg.Edges()
	lines := make([]string, edges.Size() + 2)
lines[0] = "digraph dg {"

	for i, edge := range edges.Items() {
		lines[i + 1] = fmt.Sprintf("    %v->%v;", edge.Left, edge.Right)
	}
	lines[len(lines) - 1] = "}"
	return strings.Join(lines, "\n")
}

func (dg *Digraph[T]) Mermaid() (string) {
	edges := dg.Edges()
	lines := make([]string, edges.Size() + 1)
	lines[0] = "graph TD;"

	for i, edge := range edges.Items() {
		lines[i + 1] = fmt.Sprintf("    %v-->%v;", edge.Left, edge.Right)
	}
	return strings.Join(lines, "\n")
}

func (dg *Digraph[T]) Rem(n T) {
	dg.nodes.Rem(n)
	dg.oEdges.Rem(n)
	for _, edges := range dg.oEdges.Iter() {
		edges.Rem(n)
	}
	dg.iEdges.Rem(n)
	for _, edges := range dg.iEdges.Iter() {
		edges.Rem(n)
	}
	dg.sortedNodes = nil
	dg.sources.Clear()
	dg.sinks.Clear()
}

func (dg *Digraph[T]) AddEdge(left T, right T) {
	dg.nodes.Add(left)
	dg.nodes.Add(right)
	dg.oEdges.Add(left, right)
	dg.iEdges.Add(right, left)
	dg.sortedNodes = nil
	dg.sources.Clear()
	dg.sinks.Clear()
}

func (dg *Digraph[T]) RemEdge(left T, right T) {
	oEdges := dg.oEdges.Get(left)
	oEdges.Rem(right)
	if oEdges.Empty() {
		dg.oEdges.Rem(left)
	}

	iEdges := dg.iEdges.Get(right)
	iEdges.Rem(left)
	if iEdges.Empty() {
		dg.iEdges.Rem(right)
	}
	dg.sortedNodes = nil
}

func (dg *Digraph[T]) Nodes() (util.Set[T]) {
	return dg.nodes.Clone()
}

func (dg *Digraph[T]) Edges() (util.Set[util.Pair[T]]) {
	pairs := util.MakeSet[util.Pair[T]]()
	for left, edges := range dg.oEdges.Iter() {
		for right := range edges.Iter() {
			pairs.Add(util.MakePair(left, right))
		}
	}
	return pairs
}

func (dg *Digraph[T]) OutN(left T) util.Set[T] {
	return dg.oEdges.Get(left).Clone()
}

func (dg *Digraph[T]) InN(right T) util.Set[T] {
	return dg.iEdges.Get(right).Clone()
}

func (dg *Digraph[T]) Terminals(withSource bool) util.Set[T] {
	terminals := util.MakeSet[T]()
	for n := range dg.nodes.Iter() {
		var f func (T) util.Set[T]
		if withSource {
			f = dg.InN
		} else {
			f = dg.OutN
		}
		nbors := f(n)
		if nbors.Empty() {
			terminals.Add(n)
		}
	}
	return terminals
}

func (dg *Digraph[T]) Sources() util.Set[T] {
	if dg.sources.Empty() {
		dg.sources = dg.Terminals(true)
	}
	return dg.sources.Clone()
}

func (dg *Digraph[T]) Sinks() util.Set[T] {
	if dg.sinks.Empty() {
		dg.sinks = dg.Terminals(false)
	}
	return dg.sinks.Clone()
}

func (dg *Digraph[T]) HasCycle() bool {
	visited := util.MakeSet[T]()
	stack := util.MakeSet[T]()

	var dfs func (T) bool

	dfs = func (node T) (bool) {
		visited.Add(node)
		stack.Add(node)

		for nbor := range dg.OutN(node).Iter() {
			if !visited.Has(nbor) {
				if dfs(nbor) {
					return true
				}
			} else if stack.Has(nbor) {
				return true
			}
		}
		stack.Rem(node)
		return false
	}

	for node := range dg.Nodes().Iter() {
		if !visited.Has(node) {
			if dfs(node) {
				return true
			}
		}
	}

	return false
}

func (dg *Digraph[T]) HasPath(a T, b T) bool {
	visited := util.MakeSet[T]()

	var dfs func(T, T) bool
	dfs = func(a T, b T) bool {
		visited.Add(a)
		if a == b {
			return true
		} else {
			for e := range dg.OutN(a).Iter() {
				if !visited.Has(e) {
					if dfs(e, b) {
						return true
					}
				}
			}
		}
		visited.Rem(a)
		return false
	}
	return dfs(a, b)
}

func (dg *Digraph[T]) Paths(a T, b T) [][]T{

	paths := make([][]T, 0, dg.nodes.Size())
	stack  := util.MakeStack[T]()
	visited := util.MakeSet[T]()

	var dfs func(T, T)
	dfs = func(a T, b T) {
		visited.Add(a)
		stack.Push(a)
		if a == b {
			paths = append(paths, *stack.Slice())
		} else {
			for e := range dg.OutN(a).Iter() {
				if !visited.Has(e) {
					dfs(e, b)
				}
			}
		}
		stack.Pop()
		visited.Rem(a)
	}
	dfs(a, b)
	return paths
}

func (dg *Digraph[T]) GetTopo() ([]T, error) {
	if dg.sortedNodes == nil {
		dg.sortedNodes = make([]T, dg.nodes.Size())
		i := 0
		tempDigraph := dg.Clone()
		sources := tempDigraph.Sources()
		for !sources.Empty() {
			n := sources.Pop()
			dg.sortedNodes[i] = n
			i++
			for m := range tempDigraph.OutN(n).Iter() {
				tempDigraph.RemEdge(n, m)
				if tempDigraph.InN(m).Size() == 0 {
					sources.Add(m)
				}
			}
		}
		if tempDigraph.Edges().Size() > 0 {
			return nil, fmt.Errorf("Digraph has a cycle and cannot be topologically sorted")
		}
	}
	return dg.sortedNodes, nil
}

func (dg *Digraph[T]) IsSorted(items []T) (bool, error) {
	sorted, err := dg.GetTopo()
	if err != nil {
		return false, util.ReErr(err, "Cannot sort using this Digraph")
	}

	i := 0
	for _, n := range sorted {
		if n == items[i] {
			i++
			if i >= len(items) {
				return true, nil
			}
		}
	}
 	return false, nil
}

func (dg *Digraph[T]) Sort(items []T) ([]T, error) {
	sorted := make([]T, len(items))
	i := 0
	itemSet  := util.MakeSet(items...)
	sortedItems, err := dg.GetTopo()
	if err != nil {
		return nil, util.ReErr(err, "Couldn't sort items")
	}
	for _, n := range sortedItems {
		if itemSet.Has(n) {
			itemSet.Rem(n)
			sorted[i] = n
			i++
			if itemSet.Empty() {
				return sorted, nil
			}
		}
	}
	return nil, fmt.Errorf("Couldn't sort items %+v using digraph", items)
}
