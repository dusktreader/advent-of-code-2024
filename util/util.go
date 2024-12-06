package util

import (
	"fmt"
	"iter"
	"math/rand"
	"strings"
)

func AbsInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func MapClone[T comparable, U any](m map[T]U, clone ...func(U) U) (map[T]U) {
	n := make(map[T]U)
	for k, v := range m {
		if len(clone) > 0 {
			n[k] = clone[0](v)
		} else {
			n[k] = v
		}
	}
	return n
}

func Shuffle(numbers []int) []int {
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})
	return numbers
}

type Counter [T comparable] struct {
	contents map[T]int
}

func MakeCounter[T comparable]() (Counter[T]) {
	var c Counter[T]
	c.contents = make(map[T]int)
	return c
}

func (c *Counter[T]) Incr(item T) (int) {
	if _, ok := c.contents[item]; !ok {
		c.contents[item] = 0
	}
	c.contents[item]++
	return c.contents[item]
}

func (c *Counter[T]) Get(item T) (int) {
	count, ok := c.contents[item]
	if !ok {
		c.contents[item] = 0
		return 0
	}
	return count
}

type Set [T comparable] struct {
	contents map[T]bool
}

func MakeSet[T comparable](items ...T) (Set[T]) {
	var s Set[T]
	s.contents = make(map[T]bool)
	for _, item := range items {
		s.Add(item)
	}
	return s
}

func (s Set[T]) Clone() (Set[T]) {
	return Set[T] {
		contents: MapClone(s.contents),
	}
}

func (s Set[T]) String() (string) {
	keys := []string{}
	for k, _ := range s.contents {
		keys = append(keys, fmt.Sprintf("%+v", k))
	}
	return "{" + strings.Join(keys, ", ") + "}"
}

func (s Set[T]) Size() (int) {
	return len(s.contents)
}

func (s Set[T]) Iter() (iter.Seq[T]) {
	return func(yield func(T) bool) {
		for k, _ := range s.contents {
			if !yield(k) {
				return
			}
		}
	}
}

func (s Set[T]) Eq(o Set[T]) (bool) {
	for k, _ := range s.contents {
		if !o.Has(k) {
			return false
		}
	}
	for k, _ := range o.contents {
		if !s.Has(k) {
			return false
		}
	}
	return true
}

func (s *Set[T]) Add(item T) {
	s.contents[item] = true
}

func (s *Set[T]) Pop() (item T) {
	for i, _ := range s.contents {
		item = i
		break
	}
	s.Rem(item)
	return item
}

func (s *Set[T]) Rem(item T) {
	delete(s.contents, item)
}

func (s Set[T]) Has(item T) (bool) {
	_, ok := s.contents[item]
	return ok
}

func (s Set[T]) Empty() (bool) {
	return len(s.contents) == 0
}

func (s Set[T]) Items() (items []T) {
	for k, _ := range s.contents {
		items = append(items, k)
	}
	return
}

func (s Set[T]) Ix(o Set[T]) (Set[T]) {
	n := MakeSet[T]()
	for _, v := range s.Items() {
		if o.Has(v) {
			n.Add(v)
		}
	}
	return n
}

func (s Set[T]) Diff(o Set[T]) (Set[T]) {
	n := MakeSet[T]()
	for _, v := range s.Items() {
		if !o.Has(v) {
			n.Add(v)
		}
	}
	return n
}

func (s Set[T]) Prune(item T) (Set[T]) {
	n := MakeSet[T]()
	for _, v := range s.Items() {
		if v != item {
			n.Add(v)
		}
	}
	return n
}

type Pair [T any] struct {
	Left T
	Right T
}

func MakePair[T any](left T, right T) (p Pair[T]) {
	p.Left = left
	p.Right = right
	return
}

type DAG [T comparable] struct {
	nodes       Set[T]
	sortedNodes []T
	oEdges      map[T]Set[T]
	iEdges      map[T]Set[T]
}

func MakeDag[T comparable](pairs ...Pair[T]) (DAG[T]) {
	var dag DAG[T]
	dag.oEdges = make(map[T]Set[T])
	dag.iEdges = make(map[T]Set[T])
	dag.nodes = MakeSet[T]()
	for _, p := range pairs {
		dag.AddEdge(p.Left, p.Right)
	}
	return dag
}

func (dag DAG[T]) Clone() (DAG[T]) {
	return DAG[T]{
		nodes: dag.nodes.Clone(),
		oEdges: MapClone(dag.oEdges, func (s Set[T]) (Set[T]) { return s.Clone() }),
		iEdges: MapClone(dag.iEdges, func (s Set[T]) (Set[T]) { return s.Clone() }),
	}
}

func (dag *DAG[T]) Add(a T) {
	dag.nodes.Add(a)
}

func (dag DAG[T]) String() (string) {
	nodes := []string{}
	for _, n := range dag.nodes.Items() {
		nodes = append(nodes, fmt.Sprintf("%v", n))
	}
	edges := []string{}
	for l, r := range dag.oEdges {
		edges = append(edges, fmt.Sprintf("%v->%v", l, r))
	}
	return "(" + strings.Join(nodes, ", ") + "):{" + strings.Join(edges, ", ") + "}"
}

func (dag DAG[T]) Dot() (string) {
	edges := dag.Edges()
	lines := make([]string, edges.Size() + 2)
	lines[0] = "digraph dag {"

	for i, edge := range edges.Items() {
		lines[i + 1] = fmt.Sprintf("    %v->%v;", edge.Left, edge.Right)
	}
	lines[len(lines) - 1] = "}"
	return strings.Join(lines, "\n")
}

func (dag DAG[T]) Mermaid() (string) {
	edges := dag.Edges()
	lines := make([]string, edges.Size() + 1)
	lines[0] = "graph TD;"

	for i, edge := range edges.Items() {
		lines[i + 1] = fmt.Sprintf("    %v-->%v;", edge.Left, edge.Right)
	}
	return strings.Join(lines, "\n")
}

func (dag *DAG[T]) Rem(node T) {
	dag.nodes.Rem(node)
	delete(dag.oEdges, node)
	for _, edges := range dag.oEdges {
		edges.Rem(node)
	}
	delete(dag.iEdges, node)
	for _, edges := range dag.iEdges {
		edges.Rem(node)
	}
}

func (dag *DAG[T]) AddEdge(left T, right T) {
	dag.nodes.Add(left)
	dag.nodes.Add(right)
	oEdge, ok := dag.oEdges[left]
	if !ok {
		dag.oEdges[left] = MakeSet[T]()
		oEdge = dag.oEdges[left]
	}
	oEdge.Add(right)

	iEdge, ok := dag.iEdges[right]
	if !ok {
		dag.iEdges[right] = MakeSet[T]()
		iEdge = dag.iEdges[right]
	}
	iEdge.Add(left)
	dag.sortedNodes = nil
}

func (dag *DAG[T]) RemEdge(left T, right T) {
	oEdge, ok := dag.oEdges[left]
	if ok {
		oEdge.Rem(right)
		if oEdge.Empty() {
			delete(dag.oEdges, left)
		}
	}
	iEdge, ok := dag.iEdges[right]
	if ok {
		iEdge.Rem(left)
		if iEdge.Empty() {
			delete(dag.iEdges, right)
		}
	}
	dag.sortedNodes = nil
}

func (dag DAG[T]) Nodes() (Set[T]) {
	return dag.nodes
}

func (dag DAG[T]) Edges() (Set[Pair[T]]) {
	pairs := MakeSet[Pair[T]]()
	for left, edges := range dag.oEdges {
		for right := range edges.Iter() {
			pairs.Add(Pair[T]{left, right})
		}
	}
	return pairs
}

func (dag DAG[T]) OutN(left T) (Set[T]){
	ns, ok := dag.oEdges[left]
	if !ok {
		return MakeSet[T]()
	}
	return ns
}

func (dag DAG[T]) InN(right T) (Set[T]) {
	ns, ok := dag.iEdges[right]
	if !ok {
		return MakeSet[T]()
	}
	return ns
}

func (dag DAG[T]) Terminals(withSource bool) (Set[T]) {
	terminals := MakeSet[T]()
	for n := range dag.nodes.Iter() {
		var f func (T) Set[T]
		if withSource {
			f = dag.InN
		} else {
			f = dag.OutN
		}
		nbors := f(n)
		if nbors.Empty() {
			terminals.Add(n)
		}
	}
	return terminals
}

func (dag DAG[T]) Sources() (Set[T]) {
	return dag.Terminals(true)
}

func (dag DAG[T]) Sinks() (Set[T]) {
	return dag.Terminals(false)
}

func (dag DAG[T]) HasCycle() (bool) {
	visited := MakeSet[T]()
	stack := MakeSet[T]()

	var dfs func (T) bool

	dfs = func (node T) (bool) {
		visited.Add(node)
		stack.Add(node)

		for nbor := range dag.OutN(node).Iter() {
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

	for node := range dag.Nodes().Iter() {
		if !visited.Has(node) {
			if dfs(node) {
				return true
			}
		}
	}

	return false
}

func (dag DAG[T]) GetTopo() ([]T, error) {
	if dag.sortedNodes == nil {
		dag.sortedNodes = make([]T, dag.nodes.Size())
		i := 0
		tempDag := dag.Clone()
		sources := tempDag.Sources()
		for !sources.Empty() {
			n := sources.Pop()
			dag.sortedNodes[i] = n
			i++
			for m := range tempDag.OutN(n).Iter() {
				tempDag.RemEdge(n, m)
				if tempDag.InN(m).Size() == 0 {
					sources.Add(m)
				}
			}
		}
		if tempDag.Edges().Size() > 0 {
			return nil, fmt.Errorf("Dag has a cycle and cannot be topo-sorted")
		}
	}
	return dag.sortedNodes, nil
}

func (dag DAG[T]) IsSortedDumb(items []T) bool {
	rIdx := make(map[T]int)
	for i, v := range items {
		rIdx[v] = i
	}
	for leftI, left := range items {
		for right := range dag.OutN(left).Iter() {
			rightI, ok := rIdx[right]
			if !ok {
				continue
			}
			if leftI > rightI {
				return false
			}
		}
	}
	return true
}

func (dag DAG[T]) Prune(nodes Set[T]) DAG[T] {
	prunedDag := dag.Clone()
	for n := range dag.Nodes().Iter() {
		if !nodes.Has(n) {
			prunedDag.Rem(n)
		}
	}
	return prunedDag
}

func (dag DAG[T]) IsSorted(items []T) (bool, error) {
	i := 0
	sortedNodes, err := dag.GetTopo()
	if err != nil {
		return false, fmt.Errorf("Couldn't get sorted nodes: %#v", err)
	}
	for _, n := range sortedNodes {
		if n == items[i] {
			i++
			if i == len(items) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (dag DAG[T]) SortDumb(items []T) []T {
	itemSet := MakeSet(items...)
	prunedDag := dag.Prune(itemSet)
	newItems, err := prunedDag.Sort(items)
	return newItems
}

func (dag DAG[T]) Sort(items []T) ([]T, error) {
	sorted := make([]T, len(items))
	i := 0
	itemSet  := MakeSet(items...)
	sortedItems, err := dag.GetTopo()
	if err != nil {
		return nil, fmt.Errorf("Couldn't get sorted nodes: %#v", err)
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
	return nil, fmt.Errorf("Couldn't sort items %+v using dag", items)
}
