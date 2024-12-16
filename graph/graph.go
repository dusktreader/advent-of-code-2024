package graph

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
)

// Need to add tests for this!!

type Graph [T comparable] struct {
	directed    bool
	nodes       util.Set[T]
	sources     util.Set[T]
	sinks       util.Set[T]
	sortedNodes []T
	oEdges      util.SetMap[T, T]
	iEdges      util.SetMap[T, T]
}

func MakeGraph[T comparable](directed bool, pairs ...util.Pair[T]) *Graph[T] {
	g := Graph[T]{
		directed: directed,
		oEdges:   util.MakeSetMap[T, T](),
		iEdges:   util.MakeSetMap[T, T](),
		nodes:    util.MakeSet[T](),
		sources:  util.MakeSet[T](),
		sinks:    util.MakeSet[T](),
	}
	for _, p := range pairs {
		slog.Debug("Adding edge:", "p", p)
		g.AddEdge(p.Left, p.Right)
	}
	return &g
}

func (g *Graph[T]) Eq(og *Graph[T]) bool {
	if g.directed != og.directed {
		return false
	}

	if g.directed {
		return util.All(
			g.nodes.Eq(og.nodes),
			g.Edges().Eq(og.Edges()),
		)
	}

	if !g.nodes.Eq(og.nodes) {
		return false
	}

	gEdges := g.Edges()
	oEdges := og.Edges()

	for !gEdges.Empty() {
		p := gEdges.Pop()
		rp := p.Rev()
		if !oEdges.Has(p) && !oEdges.Has(rp) {
			return false
		}
		oEdges.Rem(p)
		oEdges.Rem(rp)
	}

	if !oEdges.Empty() {
		return false
	}

	return true
}

func (g *Graph[T]) Clear() {
	g.oEdges.Clear()
	g.iEdges.Clear()
	g.nodes.Clear()
	g.sources.Clear()
	g.sinks.Clear()
}

func (g *Graph[T]) Clone() *Graph[T] {
	ng := Graph[T]{
		directed: g.directed,
		nodes:    g.nodes.Clone(),
		sources:  g.sources.Clone(),
		sinks:    g.sinks.Clone(),
		oEdges:   g.oEdges.Clone(),
		iEdges:   g.iEdges.Clone(),
	}
	return &ng
}

func (g *Graph[T]) Add(n T) {
	g.nodes.Add(n)
	g.sortedNodes = nil
	g.sources.Clear()
	g.sinks.Clear()
}

func (g Graph[T]) String() (string) {
	cx := "--"
	if g.directed {
		cx = "->"
	}
	nodes := []string{}
	for n := range g.nodes.Iter() {
		nodes = append(nodes, fmt.Sprintf("%v", n))
	}
	edges := []string{}
	for l, r := range g.oEdges.Iter() {
		edges = append(edges, fmt.Sprintf("%v%v%v", l, cx, r))
	}
	return "(" + strings.Join(nodes, ", ") + "):{" + strings.Join(edges, ", ") + "}"
}

func (g *Graph[T]) Dot() (string) {
	t := "graph"
	cx := "--"
	if g.directed {
		t = "digraph"
		cx = "->"
	}

	edges := g.Edges()
	lines := make([]string, edges.Size() + 2)
	lines[0] = fmt.Sprintf("%s g {", t)

	for i, edge := range edges.Items() {
		lines[i + 1] = fmt.Sprintf("    %v%v%v;", edge.Left, cx, edge.Right)
	}
	lines[len(lines) - 1] = "}"
	return strings.Join(lines, "\n")
}

func (g *Graph[T]) Mermaid() (string) {
	cx := "---"
	if g.directed {
		cx = "-->"
	}

	edges := g.Edges()
	lines := make([]string, edges.Size() + 1)
	lines[0] = "graph TD;"

	for i, edge := range edges.Items() {
		lines[i + 1] = fmt.Sprintf("    %v%v%v;", edge.Left, cx, edge.Right)
	}
	return strings.Join(lines, "\n")
}

func (g *Graph[T]) Rem(n T) {
	g.nodes.Rem(n)
	g.oEdges.Rem(n)
	for _, edges := range g.oEdges.Iter() {
		edges.Rem(n)
	}
	g.iEdges.Rem(n)
	for _, edges := range g.iEdges.Iter() {
		edges.Rem(n)
	}
	g.sortedNodes = nil
	g.sources.Clear()
	g.sinks.Clear()
}

func (g *Graph[T]) AddEdge(left T, right T) {
	g.nodes.Add(left)
	g.nodes.Add(right)
	g.oEdges.Add(left, right)
	g.iEdges.Add(right, left)
	g.sortedNodes = nil
	g.sources.Clear()
	g.sinks.Clear()

	if g.directed == false {
		g.oEdges.Add(right, left)
		g.iEdges.Add(left, right)
	}

}

func (g *Graph[T]) RemEdge(left T, right T) {
	oEdges := g.oEdges.Get(left)
	oEdges.Rem(right)
	if oEdges.Empty() {
		g.oEdges.Rem(left)
	}

	iEdges := g.iEdges.Get(right)
	iEdges.Rem(left)
	if iEdges.Empty() {
		g.iEdges.Rem(right)
	}

	if g.directed == false {
		oEdges := g.oEdges.Get(right)
		oEdges.Rem(left)
		if oEdges.Empty() {
			g.oEdges.Rem(right)
		}

		iEdges := g.iEdges.Get(left)
		iEdges.Rem(right)
		if iEdges.Empty() {
			g.iEdges.Rem(left)
		}
	}
	g.sortedNodes = nil
}

func (g *Graph[T]) Nodes() (util.Set[T]) {
	return g.nodes.Clone()
}

func (g *Graph[T]) Edges() (util.Set[util.Pair[T]]) {
	pairs := util.MakeSet[util.Pair[T]]()
	for left, edges := range g.oEdges.Iter() {
		for right := range edges.Iter() {
			pairs.Add(util.MakePair(left, right))
		}
	}

	if g.directed == true {
		return pairs
	}

	pruned := util.MakeSet[util.Pair[T]]()
	for !pairs.Empty() {
		p := pairs.Pop()
		rp := p.Rev()
		if pruned.Has(p) || pruned.Has(rp) {
			continue
		}
		pruned.Add(p)
	}
	return pruned
}

func (g *Graph[T]) OutN(left T) util.Set[T] {
	return g.oEdges.Get(left).Clone()
}

func (g *Graph[T]) InN(right T) util.Set[T] {
	return g.iEdges.Get(right).Clone()
}

func (g *Graph[T]) Nbors(n T) util.Set[T] {
	o := g.oEdges.Get(n)
	i := g.iEdges.Get(n)
	return o.Un(*i)
}

func (g *Graph[T]) CnxComp() []Graph[T] {
	if g.directed {
		slog.Error("Connected components not implemented yet for digraphs")
		return nil
	}

	nodes := g.Nodes()
	comp := make([]Graph[T], 0, nodes.Size() / 2)
	cnx := MakeGraph[T](g.directed)

	visited := util.MakeSet[T]()

	var dfs func(T)
	dfs = func(n T) {
		visited.Add(n)
		for e := range g.Nbors(n).Iter() {
			cnx.AddEdge(n, e)
			nodes.Rem(e)
			if !visited.Has(e) {
				dfs(e)
			}
		}
		visited.Rem(n)
	}

	for !nodes.Empty() {
		cnx.Clear()
		n := nodes.Pop()
		cnx.Add(n)
		dfs(n)
		comp = append(comp, *cnx.Clone())
	}

	return comp
}

func (g *Graph[T]) Terminals(withSource bool) util.Set[T] {
	if g.directed == false {
		return util.MakeSet[T]()
	}

	terminals := util.MakeSet[T]()
	for n := range g.nodes.Iter() {
		var f func (T) util.Set[T]
		if withSource {
			f = g.InN
		} else {
			f = g.OutN
		}
		nbors := f(n)
		if nbors.Empty() {
			terminals.Add(n)
		}
	}
	return terminals
}

func (g *Graph[T]) Sources() util.Set[T] {
	if g.directed == false {
		return util.MakeSet[T]()
	}

	if g.sources.Empty() {
		g.sources = g.Terminals(true)
	}
	return g.sources.Clone()
}

func (g *Graph[T]) Sinks() util.Set[T] {
	if g.directed == false {
		return util.MakeSet[T]()
	}

	if g.sinks.Empty() {
		g.sinks = g.Terminals(false)
	}
	return g.sinks.Clone()
}

func (g *Graph[T]) HasCycle() bool {
	if g.directed == false {
		return true
	}

	visited := util.MakeSet[T]()
	stack := util.MakeSet[T]()

	var dfs func (T) bool

	dfs = func (node T) (bool) {
		visited.Add(node)
		stack.Add(node)

		for nbor := range g.OutN(node).Iter() {
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

	for node := range g.Nodes().Iter() {
		if !visited.Has(node) {
			if dfs(node) {
				return true
			}
		}
	}

	return false
}

func (g *Graph[T]) HasPath(a T, b T) bool {
	visited := util.MakeSet[T]()

	var dfs func(T, T) bool
	dfs = func(a T, b T) bool {
		visited.Add(a)
		if a == b {
			return true
		} else {
			for e := range g.OutN(a).Iter() {
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

func (g *Graph[T]) Paths(a T, b T) [][]T{

	paths := make([][]T, 0, g.nodes.Size())
	stack  := util.MakeStack[T]()
	visited := util.MakeSet[T]()

	var dfs func(T, T)
	dfs = func(a T, b T) {
		visited.Add(a)
		stack.Push(a)
		if a == b {
			paths = append(paths, *stack.Slice())
		} else {
			for e := range g.OutN(a).Iter() {
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

func (g *Graph[T]) GetTopo() ([]T, error) {
	if g.directed == false {
		return nil, fmt.Errorf("Can't topologically sort a non-directed graph")
	}

	if g.sortedNodes == nil {
		g.sortedNodes = make([]T, g.nodes.Size())
		i := 0
		tempDigraph := g.Clone()
		sources := tempDigraph.Sources()
		for !sources.Empty() {
			n := sources.Pop()
			g.sortedNodes[i] = n
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
	return g.sortedNodes, nil
}

func (dg *Graph[T]) IsSorted(items []T) (bool, error) {
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

func (dg *Graph[T]) Sort(items []T) ([]T, error) {
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
