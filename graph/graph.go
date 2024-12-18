package graph

import (
	"container/heap"
)

type (
	Edge[T any] struct {
		V T
		C int
		H int
	}
)

func Search[T comparable](start Edge[T], isEnd func(v T) bool, getNeighbors func(v T) []Edge[T]) int {
	closedSet := newClosedSet[T]()
	openSet := newOpenSet[T]()
	openSet.Push(start.V, start.C, start.C+start.H)
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		if isEnd(cur) {
			return curg
		}
		for _, o := range getNeighbors(cur) {
			if closedSet.Has(o.V) {
				continue
			}
			g := curg + o.C
			f := g + o.H
			if vg, ok := openSet.Get(o.V); ok {
				if g < vg {
					openSet.Update(o.V, g, f)
				}
				continue
			}
			openSet.Push(o.V, g, f)
		}
	}
	return -1
}

type (
	node[T any] struct {
		value T
		g, f  int
		index int
	}

	priorityQueue[T comparable] struct {
		q []*node[T]
		s map[T]int
	}

	openSet[T comparable] struct {
		q *priorityQueue[T]
	}

	closedSet[T comparable] struct {
		m map[T]struct{}
	}
)

func newOpenSet[T comparable]() *openSet[T] {
	return &openSet[T]{
		q: newPriorityQueue[T](),
	}
}

func (s *openSet[T]) Empty() bool {
	return s.q.Len() == 0
}

func (s *openSet[T]) Has(val T) bool {
	_, ok := s.q.s[val]
	return ok
}

func (s *openSet[T]) Get(val T) (int, bool) {
	idx, ok := s.q.s[val]
	if !ok {
		return 0, false
	}
	return s.q.q[idx].g, true
}

func (s *openSet[T]) Push(value T, g, f int) {
	heap.Push(s.q, &node[T]{
		value: value,
		g:     g,
		f:     f,
	})
}

func (s *openSet[T]) Pop() (T, int, int) {
	item := heap.Pop(s.q).(*node[T])
	return item.value, item.g, item.f
}

func (s *openSet[T]) Update(value T, g, f int) bool {
	return s.q.Update(value, g, f)
}

func newClosedSet[T comparable]() *closedSet[T] {
	return &closedSet[T]{
		m: map[T]struct{}{},
	}
}

func (s *closedSet[T]) Has(val T) bool {
	_, ok := s.m[val]
	return ok
}

func (s *closedSet[T]) Push(val T) {
	s.m[val] = struct{}{}
}

func newPriorityQueue[T comparable]() *priorityQueue[T] {
	return &priorityQueue[T]{
		s: map[T]int{},
	}
}

func (q *priorityQueue[T]) Len() int { return len(q.q) }
func (q *priorityQueue[T]) Less(i, j int) bool {
	return q.q[i].f < q.q[j].f
}

func (q *priorityQueue[T]) Swap(i, j int) {
	q.q[i], q.q[j] = q.q[j], q.q[i]
	q.q[i].index = i
	q.q[j].index = j
	q.s[q.q[i].value] = i
	q.s[q.q[j].value] = j
}

func (q *priorityQueue[T]) Push(x any) {
	n := len(q.q)
	item := x.(*node[T])
	item.index = n
	q.q = append(q.q, item)
	q.s[item.value] = n
}

func (q *priorityQueue[T]) Pop() any {
	n := len(q.q)
	item := q.q[n-1]
	q.q[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	q.q = q.q[:n-1]
	delete(q.s, item.value)
	return item
}

func (q *priorityQueue[T]) Update(value T, g, f int) bool {
	idx, ok := q.s[value]
	if !ok {
		return false
	}
	item := q.q[idx]
	item.g = g
	item.f = f
	heap.Fix(q, item.index)
	return true
}
