package heap_test

import (
	"log/slog"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/heap"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func TestMakeHeap(t *testing.T) {
	h := heap.MakeHeap[int](true)
	if h.Size() != 0 {
		t.Errorf("Expected empty heap, got %v", h)
	}

	_, _, err := h.Extract()
	if err == nil {
		t.Errorf("Empty heap should not allow extract!")
	}

	weights, values := h.Dump()
	if len(weights) != 0 || len(values) != 0 {
		t.Errorf("Expected empty heap!")
	}
}

func TestInsertMaxHeap(t *testing.T) {
	h := heap.MakeHeap[rune](true)
	h.Insert(1, 'a')
	h.Insert(2, 'b')
	h.Insert(3, 'c')
	h.Insert(4, 'd')
	h.Insert(5, 'e')

	if !h.Valid() {
		t.Errorf("Heap not valid after inserts! %v", h)
	}

	h.Insert(0, 'f')
	if !h.Valid() {
		t.Errorf("Heap not valid after insert! %v", h)
	}

	h.Insert(6, 'g')
	if !h.Valid() {
		t.Errorf("Heap not valid after insert! %v", h)
	}
}

func TestInsertMinHeap(t *testing.T) {
	h := heap.MakeHeap[rune](false)
	h.Insert(1, 'a')
	h.Insert(2, 'b')
	h.Insert(3, 'c')
	h.Insert(4, 'd')
	h.Insert(5, 'e')

	if !h.Valid() {
		t.Errorf("Heap not valid after inserts! %v", h)
	}

	h.Insert(0, 'f')
	if !h.Valid() {
		t.Errorf("Heap not valid after insert! %v", h)
	}

	h.Insert(6, 'g')
	if !h.Valid() {
		t.Errorf("Heap not valid after insert! %v", h)
	}
}

func TestChangeWeightMaxHeap(t *testing.T) {
	h := heap.MakeHeap[rune](true)
	h.Insert(1, 'a')
	h.Insert(2, 'b')
	h.Insert(3, 'c')
	h.Insert(4, 'd')
	h.Insert(5, 'e')

	h.ChangeWeight(99, 'c')
	if !h.Valid() {
		t.Errorf("Heap not valid after change weight! %v", h)
	}
}

func TestChangeWeightMinHeap(t *testing.T) {
	h := heap.MakeHeap[rune](false)
	h.Insert(1, 'a')
	h.Insert(2, 'b')
	h.Insert(3, 'c')
	h.Insert(4, 'd')
	h.Insert(5, 'e')

	h.ChangeWeight(99, 'c')
	if !h.Valid() {
		t.Errorf("Heap not valid after change weight! %v", h)
	}
}


func TestExtractMaxHeap(t *testing.T) {
	h := heap.MakeHeap[rune](true)
	h.Insert(1, 'a')
	h.Insert(2, 'b')
	h.Insert(3, 'c')
	h.Insert(4, 'd')
	h.Insert(5, 'e')
	h.Insert(6, 'f')
	h.Insert(7, 'g')

	gotWeight, gotValue, err := h.Extract()
	util.Unexpect(t, err)
	if !h.Valid() {
		t.Errorf("Heap not valid after extract! %v", h)
	}
	wantWeight, wantValue := 7, 'g'
	if gotWeight != wantWeight || gotValue != wantValue {
		t.Errorf("Expected %v, %v, got %v, %v", wantWeight, wantValue, gotWeight, gotValue)
	}

	h.ChangeWeight(99, 'c')

	gotWeight, gotValue, err = h.Extract()
	util.Unexpect(t, err)
	if !h.Valid() {
		t.Errorf("Heap not valid after extract! %v", h)
	}
	wantWeight, wantValue = 99, 'c'
	if gotWeight != wantWeight || gotValue != wantValue {
		t.Errorf("Expected %v, %v, got %v, %v", wantWeight, wantValue, gotWeight, gotValue)
	}

	gotWeight, gotValue, err = h.Extract()
	util.Unexpect(t, err)
	if !h.Valid() {
		t.Errorf("Heap not valid after extract! %v", h)
	}
	wantWeight, wantValue = 6, 'f'
	if gotWeight != wantWeight || gotValue != wantValue {
		t.Errorf("Expected %v, %v, got %v, %v", wantWeight, wantValue, gotWeight, gotValue)
	}
}

func TestExtractMinHeap(t *testing.T) {
	h := heap.MakeHeap[rune](false)
	h.Insert(1, 'a')
	h.Insert(2, 'b')
	h.Insert(3, 'c')
	h.Insert(4, 'd')
	h.Insert(5, 'e')
	h.Insert(6, 'f')
	h.Insert(7, 'g')

	gotWeight, gotValue, err := h.Extract()
	util.Unexpect(t, err)
	if !h.Valid() {
		t.Errorf("Heap not valid after extract! %v", h)
	}
	wantWeight, wantValue := 1, 'a'
	if gotWeight != wantWeight || gotValue != wantValue {
		t.Errorf("Expected %v, %v, got %v, %v", wantWeight, wantValue, gotWeight, gotValue)
	}

	h.ChangeWeight(0, 'c')

	gotWeight, gotValue, err = h.Extract()
	util.Unexpect(t, err)
	if !h.Valid() {
		t.Errorf("Heap not valid after extract! %v", h)
	}
	wantWeight, wantValue = 0, 'c'
	if gotWeight != wantWeight || gotValue != wantValue {
		t.Errorf("Expected %v, %v, got %v, %v", wantWeight, wantValue, gotWeight, gotValue)
	}

	gotWeight, gotValue, err = h.Extract()
	util.Unexpect(t, err)
	if !h.Valid() {
		t.Errorf("Heap not valid after extract! %v", h)
	}
	wantWeight, wantValue = 2, 'b'
	if gotWeight != wantWeight || gotValue != wantValue {
		t.Errorf("Expected %v, %v, got %v, %v", wantWeight, wantValue, gotWeight, gotValue)
	}
}

func TestAllKindsOfShit(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	n := 7
	h := heap.MakeMaxHeap[int]()
	for i := 0; i < n; i++ {
		h.Insert(i, 1000 + i)
		if !h.Valid() {
			t.Fatalf("Heap not valid after insert! %v", h)
		}
	}

	for i := 0; i < n; i++ {
		h.ChangeWeight(n - i - 1, 1000 + i)
		if !h.Valid() {
			t.Fatalf("Heap not valid after change weight for value %v to %v! %v", 1000 + i, n - i - 1, h)
		}
	}

	for i := 0; i < n; i++ {
		gotWeight, gotValue, err := h.Extract()
		util.Unexpect(t, err)
		wantWeight, wantValue := n - i - 1, 1000 + i
		if !h.Valid() {
			t.Fatalf("Heap not valid after extract! %v", h)
		}

		if gotWeight != wantWeight || gotValue != wantValue {
			t.Fatalf("Expected %v, %v, got %v, %v", wantWeight, wantValue, gotWeight, gotValue)
		}
	}
	if h.Size() != 0 {
		t.Fatalf("Heap should be empty!")
	}
}
