package util

import (
	"fmt"
	"iter"
	"log/slog"
	"math"
	"math/rand"
	"strings"
	"testing"
)

func ReErr(err error, msg string, fmtArgs ...any) error {
	fmtMsg := fmt.Sprintf(msg, fmtArgs...)
	return fmt.Errorf("%s: %+v", fmtMsg, err)
}

func Unexpect(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}
}

func AbsInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func MinI(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func RtoI(rn rune) (int, error) {
	v := int(rn - '0')
	if v < 0 || v > 9 {
		return -1, fmt.Errorf("Couldn't rune to integer: %v", rn)
	}
	return v, nil
}

func ItoR(v int) (rune, error) {
	if v < 0 || v > 9 {
		return -1, fmt.Errorf("Couldn't integer to rune: %v", v)
	}
	rn := rune('0' + v)
	return rn, nil
}

func PowInt(x int, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func Pow2(x int) int {
	return 1 << x
}

func Pow3(x int) int {
	return int(math.Pow(3, float64(x)))
}

func DigiCt(i int) int {
	if i == 0 {
		return 1
	}
	return int(math.Log10(float64(AbsInt(i))) + 1)
}

func GCD(a int, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func All(vals ...bool) bool {
	for _, v := range vals {
		if v != true {
			return false
		}
	}
	return true
}

func Any(vals ...bool) bool {
	for _, v := range vals {
		if v == true {
			return true
		}
	}
	return false
}

func MakeFill[T any](v T, l int, c ...int) []T {
	var slice []T
	if len(c) == 0 {
		slice = make([]T, l)
	} else {
		slice = make([]T, l, c[0])
	}
	for i := range l {
		slice[i] = v
	}
	return slice
}

func Insert[T any](l []T, i int, v ...T) []T {
	return append(l[:i], append(v, l[i:]...)...)
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

func KeySet[T comparable, U any](m map[T]U) Set[T] {
	s := MakeSet[T]()
	for k := range m {
		s.Add(k)
	}
	return s
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

type Stack[T any] struct {
	contents []T
}

func MakeStack[T any](items ...T) *Stack[T] {
	st := Stack[T]{}
	st.contents = make([]T, len(items))
	copy(items, st.contents)
	return &st
}

func (st Stack[T]) String() (string) {
	vs := []string{}
	for _, v := range st.contents {
		vs = append(vs, fmt.Sprintf("%+v", v))
	}
	return "[" + strings.Join(vs, ", ") + ">"
}

func (st *Stack[T]) Size() int {
	return len(st.contents)
}

func (st *Stack[T]) Push(items ...T) {
	st.contents = append(st.contents, items...)
}

func (st *Stack[T]) Pop() (T, error) {
	if st.Size() == 0 {
		var null T
		return null, fmt.Errorf("Stack is empty!")
	}
	i := st.contents[st.Size() - 1]
	st.contents = st.contents[:st.Size() - 1]
	return i, nil
}

func (st *Stack[T]) Clone() *Stack[T] {
	return MakeStack(st.contents...)
}

func (st *Stack[T]) Slice() *[]T {
	sl := make([]T, st.Size())
	copy(sl, st.contents)
	return &sl
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
	for k := range s.contents {
		keys = append(keys, fmt.Sprintf("%+v", k))
	}
	return "{" + strings.Join(keys, ", ") + "}"
}

func (s Set[T]) Size() (int) {
	return len(s.contents)
}

func (s Set[T]) Iter() (iter.Seq[T]) {
	return func(yield func(T) bool) {
		for k := range s.contents {
			if !yield(k) {
				return
			}
		}
	}
}

func (s Set[T]) Eq(o Set[T]) (bool) {
	for k := range s.contents {
		if !o.Has(k) {
			return false
		}
	}
	for k := range o.contents {
		if !s.Has(k) {
			return false
		}
	}
	return true
}

func (s *Set[T]) Clear() {
	s.contents = make(map[T]bool)
}

func (s *Set[T]) Add(items ...T) {
	for _, v := range items {
		s.contents[v] = true
	}
}

func (s *Set[T]) Pop() (item T) {
	for i := range s.contents {
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
	for k := range s.contents {
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

func (s Set[T]) Un(o Set[T]) (Set[T]) {
	n := s.Clone()
	n.Add(o.Items()...)
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

type SetMap [T comparable, U comparable] struct {
	contents map[T]Set[U]
}

func MakeSetMap[T comparable, U comparable]() SetMap[T, U] {
	sm := SetMap[T, U]{}
	sm.contents = make(map[T]Set[U])
	return sm
}

func (sm SetMap[T, U]) String() (string) {
	sets := make([]string, len(sm.contents))
	i := 0
	for k, s := range sm.contents {
		sets[i] = fmt.Sprintf("%+v: %v", k, s)
		i++
	}
	return "{" + strings.Join(sets, ", ") + "}"
}


func (sm *SetMap[T, U]) Iter() iter.Seq2[T, Set[U]] {
	return func(yield func(T, Set[U]) bool) {
		for k, s := range sm.contents {
			if !yield(k, s) {
				return
			}
		}
	}
}

func (sm *SetMap[T, U]) Add(key T, items ...U) {
	s, ok := sm.contents[key]
	if !ok {
		sm.contents[key] = MakeSet(items...)
	} else {
		s.Add(items...)
	}
}

func (sm *SetMap[T, U]) Rem(key T) {
	delete(sm.contents, key)
}

func (sm *SetMap[T, U]) Size() int {
	return len(sm.contents)
}

func (sm *SetMap[T, U]) Has(key T) bool {
	_, ok := sm.contents[key]
	return ok
}

func (sm *SetMap[T, U]) Get(key T) *Set[U] {
	s, ok := sm.contents[key]
	if !ok {
		s = MakeSet[U]()
		sm.contents[key] = s
	}
	return &s
}

func (sm *SetMap[T, U]) Pop(key T) *Set[U] {
	s, ok := sm.contents[key]
	if !ok {
		s := MakeSet[U]()
		return &s
	} else {
		sm.Rem(key)
	}
	return &s
}

func (sm *SetMap[T, U]) Clear() {
	sm.contents = make(map[T]Set[U])
}

func (sm *SetMap[T, U]) Eq(om *SetMap[T, U]) bool {
	sKeys := KeySet(sm.contents)
	oKeys := KeySet(om.contents)
	if !sKeys.Eq(oKeys) {
		return false
	}
	for k := range sKeys.Iter() {
		if !sm.Get(k).Eq(*om.Get(k)) {
			return false
		}
	}
	return true
}

func (sm *SetMap[T, U]) Clone() SetMap[T, U] {
	nm := MakeSetMap[T, U]()
	nm.contents = MapClone(sm.contents, func (s Set[U]) Set[U] { return s.Clone() })
	return nm
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

type Point struct {
	I int
	J int
}

func MakePoint(i int, j int) (p Point) {
	p.I = i
	p.J = j
	return
}

func (p Point) String() (string) {
	return fmt.Sprintf("(%v, %v)", p.I, p.J)
}

func (p Point) Add(v Vector) Point {
	return Point{I: p.I + v.Di, J: p.J + v.Dj}
}

func (p Point) Diff(o Point) Vector {
	return Vector{Di: p.I - o.I, Dj: p.J - o.J}
}

func (p Point) Mul(d int) Point {
	return Point{I: p.I * d, J: p.J * d}
}

type Vector struct {
	Di int
	Dj int
}

func MakeVector(di int, dj int) (v Vector) {
	v.Di = di
	v.Dj = dj
	return
}

func (v Vector) Neg() (Vector) {
	return Vector{Di: -v.Di, Dj: -v.Dj}
}

func (v Vector) Div(d int) Vector {
	return Vector{Di: v.Di / d, Dj: v.Dj /d}
}

func (v Vector) Mul(m int) Vector {
	return Vector{Di: v.Di * m, Dj: v.Dj * m}
}

func (v Vector) Rot() (Vector) {
	if v.Di == 0 {
		return Vector{Di: v.Dj, Dj: 0}
	} else {
		return Vector{Di: 0, Dj: -v.Di}
	}
}

func (v Vector) String() string {
	return fmt.Sprintf("<%v, %v>", v.Di, v.Dj)
}

func (v Vector) Pretty() rune {
	if v.Di != 0 && v.Dj != 0 {
		return '?'
	} else if v.Di == 0 {
		if v.Dj == -1 {
			return '<'
		} else if v.Dj == 0 {
			return 'o'
		} else {
			return '>'
		}
	} else if v.Di == -1 {
		return '^'
	} else {
		return 'v'
	}
}

type Size struct {
	W int
	H int
}

func MakeSize(w int, h int) (Size, error) {
	s := Size{W: w, H: h}
	if w < 0 || h < 0 {
		return s, fmt.Errorf("Invalid size %s", s)
	}
	return s, nil
}

func (s Size) String() string {
	return fmt.Sprintf("[%v, %v]", s.W, s.H)
}

func (s Size) Area() int {
	return s.W * s.H
}

func (s Size) Idx(p Point) (int, error) {
	if s.Out(p) {
		return -1, fmt.Errorf("Point was out of bounds: %d, %d", p.I, p.J)
	}
	return s.W * p.I + p.J, nil
}

func (s Size) Out(p Point) bool {
	return p.I < 0 || p.J < 0 || p.I >= s.H || p.J >= s.W
}

func (s Size) Iter() (iter.Seq[Point]) {
	return func(yield func(Point) bool) {
		for i := range s.H {
			for j := range s.W {
				if !yield(MakePoint(i, j)) {
					return
				}
			}
		}
	}
}

type Grid [T any] struct {
	items []T
	size  Size
}

func MakeGrid [T any] (s Size, items ...[]T) (*Grid[T], error) {
	grid := Grid[T]{}
	if len(items) > 0 {
		if len(items[0]) != s.W * s.H {
			return nil, fmt.Errorf("Passed items length (%v) doesn't match dimensions: %v", len(items[0]), s)
		}
		grid.items = items[0]
	} else {
		grid.items = make([]T, s.W * s.H)
	}
	grid.size = s
	return &grid, nil
}

func (g Grid[T]) Size() Size {
	return g.size
}

func (g Grid[T]) Out(p Point) bool {
	return p.I < 0 || p.J < 0 || p.I >= g.size.H || p.J >= g.size.W
}

func (g Grid[T]) Get(p Point) (*T, error) {
	if g.Out(p) {
		return nil, fmt.Errorf("Point was out of bounds: %d, %d", p.I, p.J)
	}
	return &g.items[g.size.W * p.I + p.J], nil
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
	slog.Error("Should really do something with this:", "err", err)
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
