package heap

import (
	"fmt"
	"log/slog"
)

type heapNode [T comparable] struct {
	weight int
	value  T
}

type Heap [T comparable] struct {
	contents []heapNode[T]
	valueMap map[T]int
	isMax    bool
}

func MakeHeap[T comparable](isMax bool, capacity...int) (*Heap[T]) {
	var cap int
	if len(capacity) > 0 {
		cap = capacity[0]
	} else {
		cap = 0
	}
	h := Heap[T]{
		contents: make([]heapNode[T], 0, cap),
		valueMap: make(map[T]int),
		isMax:    isMax,
	}
	return &h
}

func MakeMaxHeap[T comparable](capacity...int) (*Heap[T]) {
	return MakeHeap[T](true, capacity...)
}

func MakeMinHeap[T comparable](capacity...int) (*Heap[T]) {
	return MakeHeap[T](false, capacity...)
}

func (h *Heap[T]) String() string {
	out := "Heap: [ "
	for _, n := range h.contents {
		out += fmt.Sprintf("{%v, %+v} ", n.weight, n.value)
	}
	out += "]"
	return out
}

func (h *Heap[T]) Size() int {
	return len(h.contents)
}

func (h *Heap[T]) Empty() bool {
	return len(h.contents) == 0
}

func (h *Heap[T]) Dump() ([]int, []T) {
	weights := make([]int, len(h.contents))
	values := make([]T, len(h.contents))
	for i, n := range h.contents {
		weights[i] = n.weight
		values[i] = n.value
	}
	return weights, values
}

func (h *Heap[T]) parent(i int) int {
	return (i - 1) / 2
}

func (h *Heap[T]) left(i int) int {
	return 2 * i + 1
}

func (h *Heap[T]) right(i int) int {
	return 2 * i + 2
}

func (h * Heap[T]) swap(i int, j int) {
	h.contents[i], h.contents[j] = h.contents[j], h.contents[i]
	h.valueMap[h.contents[i].value] = i
	h.valueMap[h.contents[j].value] = j
}

func (h *Heap[T]) cmp(i int, j int) bool {
	if i < 0 || i >= len(h.contents) || j < 0 || j >= len(h.contents) {
		panic("Out of bounds!")
	}

	iw := h.contents[i].weight
	jw := h.contents[j].weight

	if h.isMax {
		return iw > jw
	} else {
		return iw < jw
	}
}

func (h *Heap[T]) Valid() bool {
	for i := 0; i < len(h.contents); i++ {
		l := h.left(i)
		if l < len(h.contents) && h.cmp(l, i) {
			return false
		}

		r := h.right(i)
		if r < len(h.contents) && h.cmp(r, i) {
			return false
		}
	}
	return true
}

func (h *Heap[T]) fix(i int) {
	slog.Debug("Fixing:", "i", i)
	for i > 0 && h.cmp(i, h.parent(i)) {
		slog.Debug("Swapping!")
		h.swap(i, h.parent(i))
		i = h.parent(i)
	}
}

func (h *Heap[T]) Insert(weight int, value T) {
	h.contents = append(h.contents, heapNode[T]{weight, value})
	h.valueMap[value] = len(h.contents) - 1
	h.fix(len(h.contents) - 1)
}

func (h *Heap[T]) ChangeWeight(weight int, value T) error {
	slog.Debug("Changing weight:", "weight", weight, "value", value)
	i, ok := h.valueMap[value]
	if !ok {
		return fmt.Errorf("Couldn't find value %v in heap", value)
	}

	h.contents[i].weight = weight
	h.fix(i)
	return nil
}

func (h *Heap[T]) Extract() (int, T, error) {
	if h.Size() == 0 {
		var null T
		return 0, null, fmt.Errorf("Heap is empty")
	}
	root := h.contents[0]
	h.swap(0, len(h.contents) - 1)
	h.contents = h.contents[:len(h.contents) - 1]
	delete(h.valueMap, root.value)
	h.heapify(0)
	return root.weight, root.value, nil
}

// Don't need delete, so don't implement yet

func (h *Heap[T]) heapify(i int) {
	l := h.left(i)
	r := h.right(i)
	ext := i
	if l < len(h.contents) && h.cmp(l, i) {
		ext = l
	}
	if r < len(h.contents) && h.cmp(r, ext) {
		ext = r
	}
	if ext != i {
		h.swap(i, ext)
		h.heapify(ext)
	}
}
