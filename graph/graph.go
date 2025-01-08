package graph

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
)

// Need to add tests for this!!

type Edge [T any] struct {
	From T
	To   T
	Wt   int
}

func (e Edge[T]) Rev() Edge[T] {
	return Edge[T]{From: e.To, To: e.From, Wt: e.Wt}
}

func (e Edge[T]) String() string {
	return fmt.Sprintf("%v -> %v", e.From, e.To)
}

type Graph [T comparable] struct {
	directed    bool
	nodes       util.Set[T]
	sources     util.Set[T]
	sinks       util.Set[T]
	sortedNodes []T
	edges       util.Set[Edge[T]]
	nodeEdgeMap util.SetMap[T, Edge[T]]
}

func MakeGraph[T comparable](directed bool, pairs ...util.Pair[T]) *Graph[T] {
	g := Graph[T]{
		directed: directed,
		edges:    util.MakeSet[Edge[T]](),
		nodeEdgeMap: util.MakeSetMap[T, Edge[T]](),
		nodes:    util.MakeSet[T](),
		sources:  util.MakeSet[T](),
		sinks:    util.MakeSet[T](),
	}
	for _, p := range pairs {
		g.AddEdge(p.Left, p.Right)
	}
	return &g
}

func (g *Graph[T]) Eq(og *Graph[T]) bool {
	if g.directed != og.directed {
		return false
	}

	if !g.nodes.Eq(og.nodes) {
		return false
	}

	if g.directed {
		return g.edges.Eq(og.edges)
	}

	gEdges := g.edges.Clone()
	oEdges := og.edges.Clone()

	for !gEdges.Empty() {
		e := gEdges.Pop()
		re := e.Rev()
		if !oEdges.Has(e) && !oEdges.Has(re) {
			return false
		}
		oEdges.Rem(e)
		oEdges.Rem(re)
	}

	if !oEdges.Empty() {
		return false
	}

	return true
}

func (g *Graph[T]) Clear() {
	g.edges.Clear()
	g.nodeEdgeMap.Clear()
	g.nodes.Clear()
	g.sources.Clear()
	g.sinks.Clear()
}

func (g *Graph[T]) Clone() *Graph[T] {
	ng := Graph[T]{
		directed:    g.directed,
		nodes:       g.nodes.Clone(),
		sources:     g.sources.Clone(),
		sinks:       g.sinks.Clone(),
		edges:       g.edges.Clone(),
		nodeEdgeMap: g.nodeEdgeMap.Clone(),
	}
	return &ng
}

func (g *Graph[T]) Has(n T) bool {
	return g.nodes.Has(n)
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
	edgeStrs := []string{}
	for e := range g.edges.Iter() {
		edgeStrs = append(edgeStrs, fmt.Sprintf("%v%v%v", e.From, cx, e.To))
	}
	return "(" + strings.Join(nodes, ", ") + "):{" + strings.Join(edgeStrs, ", ") + "}"
}

func (g *Graph[T]) Dot() (string) {
	t := "graph"
	cx := "--"
	if g.directed {
		t = "digraph"
		cx = "->"
	}

	lines := make([]string, g.edges.Size() + 2)
	lines[0] = fmt.Sprintf("%s g {", t)

	for i, edge := range g.edges.Items() {
		lines[i + 1] = fmt.Sprintf("    %v%v%v;", edge.From, cx, edge.To)
	}
	lines[len(lines) - 1] = "}"
	return strings.Join(lines, "\n")
}

func (g *Graph[T]) Mermaid() (string) {
	cx := "---"
	if g.directed {
		cx = "-->"
	}

	lines := make([]string, g.edges.Size() + 1)
	lines[0] = "graph TD;"

	for i, edge := range g.edges.Items() {
		lines[i + 1] = fmt.Sprintf("    %v%v%v;", edge.From, cx, edge.To)
	}
	return strings.Join(lines, "\n")
}

func (g *Graph[T]) Rem(n T) {
	edges := g.nodeEdgeMap.Get(n)
	for edge := range edges.Iter() {
		g.RemEdge(edge.From, edge.To)
	}
	g.nodes.Rem(n)
	g.sortedNodes = nil
	g.sources.Clear()
	g.sinks.Clear()
}

func (g *Graph[T]) AddEdge(from T, to T, weight ...int) {
	var wt int
	if len(weight) == 0 {
		wt = 1
	} else {
		wt = weight[0]
	}

	edge := Edge[T]{From: from, To: to, Wt: wt}

	g.nodes.Add(from)
	g.nodes.Add(to)
	g.edges.Add(edge)
	g.nodeEdgeMap.Add(from, edge)
	g.nodeEdgeMap.Add(to, edge)
	g.sortedNodes = nil
	g.sources.Clear()
	g.sinks.Clear()
}

func (g *Graph[T]) RemEdge(left T, right T) {
	leftEdges := g.nodeEdgeMap.Get(left)
	for edge := range leftEdges.Iter() {
		if edge.From == left && edge.To == right {
			g.edges.Rem(edge)
			leftEdges.Rem(edge)
			break
		} else if !g.directed && edge.To == left && edge.From == right {
			g.edges.Rem(edge)
			leftEdges.Rem(edge)
			break
		}
	}

	rightEdges := g.nodeEdgeMap.Get(right)
	for edge := range rightEdges.Iter() {
		if edge.From == left && edge.To == right {
			g.edges.Rem(edge)
			rightEdges.Rem(edge)
			break
		} else if !g.directed && edge.To == left && edge.From == right {
			g.edges.Rem(edge)
			rightEdges.Rem(edge)
			break
		}
	}

	g.sortedNodes = nil
}

func (g *Graph[T]) Nodes() (util.Set[T]) {
	return g.nodes.Clone()
}

func (g *Graph[T]) Edge(from T, to T) (Edge[T], bool) {
	for edge := range g.nodeEdgeMap.Get(from).Iter() {
		if edge.From == from && edge.To == to {
			return edge, true
		}

		if !g.directed && edge.From == to && edge.To == from {
			return edge, true
		}
	}
	return Edge[T]{}, false
}

// This is dumb. It should return a set of edges, not a set of pairs
func (g *Graph[T]) Edges() (util.Set[util.Pair[T]]) {
	pairs := util.MakeSet[util.Pair[T]]()
	for edge := range g.edges.Iter() {
		pairs.Add(util.MakePair(edge.From, edge.To))
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
	nodes := util.MakeSet[T]()
	for edge := range g.nodeEdgeMap.Get(left).Iter() {
		if edge.From == left {
			nodes.Add(edge.To)
		}
	}
	return nodes
}

func (g *Graph[T]) InN(right T) util.Set[T] {
	nodes := util.MakeSet[T]()
	for edge := range g.nodeEdgeMap.Get(right).Iter() {
		if edge.To == right {
			nodes.Add(edge.From)
		}
	}
	return nodes
}

func (g *Graph[T]) Nbors(n T) util.Set[T] {
	nodes := util.MakeSet[T]()
	for edges := range g.nodeEdgeMap.Get(n).Iter() {
		nodes.Add(edges.From)
		nodes.Add(edges.To)
	}
	nodes.Rem(n)
	return nodes
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
