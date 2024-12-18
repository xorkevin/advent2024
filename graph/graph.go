package graph

import (
	"container/heap"

	"golang.org/x/exp/constraints"
)

type (
	Edge[T any, I constraints.Unsigned] struct {
		V T
		C I
		H I
	}

	Graph[T any, I constraints.Unsigned] interface {
		Edges(v T) []Edge[T, I]
		IsEnd(v T) bool
	}
)

type (
	IndexMap[T any, I constraints.Unsigned] interface {
		Get(k T) (I, bool)
		Set(k T, v I)
		Rm(k T)
	}

	IndexHashMap[T comparable, I constraints.Unsigned] struct {
		m map[T]I
	}
)

func NewIndexHashMap[T comparable, I constraints.Unsigned]() *IndexHashMap[T, I] {
	return &IndexHashMap[T, I]{
		m: map[T]I{},
	}
}

func (m *IndexHashMap[T, I]) Get(k T) (I, bool) {
	v, ok := m.m[k]
	return v, ok
}

func (m *IndexHashMap[T, I]) Set(k T, v I) {
	m.m[k] = v
}

func (m *IndexHashMap[T, I]) Rm(k T) {
	delete(m.m, k)
}

func (m *IndexHashMap[T, I]) Reset() {
	clear(m.m)
}

type (
	Set[T any] interface {
		Has(k T) bool
		Set(k T)
	}

	HashSet[T comparable] struct {
		m map[T]struct{}
	}
)

func NewHashSet[T comparable]() *HashSet[T] {
	return &HashSet[T]{
		m: map[T]struct{}{},
	}
}

func (s *HashSet[T]) Has(val T) bool {
	_, ok := s.m[val]
	return ok
}

func (s *HashSet[T]) Set(val T) {
	s.m[val] = struct{}{}
}

func (s *HashSet[T]) Reset() {
	clear(s.m)
}

type (
	EdgeMap[T any] interface {
		Set(to, from T)
	}
)

func Search[T any, I, J constraints.Unsigned](start Edge[T, I], priorityMap IndexMap[T, J], closedSet Set[T], graph Graph[T, I], edgeMap EdgeMap[T]) (I, bool) {
	openSet := newOpenSet[T, I](priorityMap)
	openSet.Push(start.V, start.C, start.C+start.H)
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Set(cur)
		if graph.IsEnd(cur) {
			return curg, true
		}
		for _, o := range graph.Edges(cur) {
			if closedSet.Has(o.V) {
				continue
			}
			g := curg + o.C
			f := g + o.H
			if vg, ok := openSet.Get(o.V); ok {
				if g < vg {
					openSet.Update(o.V, g, f)
					if edgeMap != nil {
						edgeMap.Set(o.V, cur)
					}
				}
				continue
			}
			openSet.Push(o.V, g, f)
			if edgeMap != nil {
				edgeMap.Set(o.V, cur)
			}
		}
	}
	var z I
	return z, false
}

type (
	openSet[T any, I, J constraints.Unsigned] struct {
		q *priorityQueue[T, I, J]
	}
)

func newOpenSet[T any, I, J constraints.Unsigned](m IndexMap[T, J]) *openSet[T, I, J] {
	return &openSet[T, I, J]{
		q: newPriorityQueue[T, I](m),
	}
}

func (s *openSet[T, I, J]) Empty() bool {
	return s.q.Len() == 0
}

func (s *openSet[T, I, J]) Get(val T) (I, bool) {
	idx, ok := s.q.m.Get(val)
	if !ok {
		return 0, false
	}
	return s.q.q[idx].g, true
}

func (s *openSet[T, I, J]) Push(value T, g, f I) {
	heap.Push(s.q, node[T, I, J]{
		value: value,
		g:     g,
		f:     f,
	})
}

func (s *openSet[T, I, J]) Pop() (T, I, I) {
	item := heap.Pop(s.q).(node[T, I, J])
	return item.value, item.g, item.f
}

func (s *openSet[T, I, J]) Update(value T, g, f I) bool {
	return s.q.Update(value, g, f)
}

type (
	node[T any, I, J constraints.Unsigned] struct {
		value T
		g, f  I
		idx   J
	}

	priorityQueue[T any, I, J constraints.Unsigned] struct {
		q []node[T, I, J]
		m IndexMap[T, J]
	}
)

func newPriorityQueue[T any, I, J constraints.Unsigned](m IndexMap[T, J]) *priorityQueue[T, I, J] {
	return &priorityQueue[T, I, J]{m: m}
}

func (q *priorityQueue[T, I, J]) Len() int { return len(q.q) }
func (q *priorityQueue[T, I, J]) Less(i, j int) bool {
	return q.q[i].f < q.q[j].f
}

func (q *priorityQueue[T, I, J]) Swap(i, j int) {
	q.q[i], q.q[j] = q.q[j], q.q[i]
	q.q[i].idx = J(i)
	q.q[j].idx = J(j)
	q.m.Set(q.q[i].value, J(i))
	q.m.Set(q.q[j].value, J(j))
}

func (q *priorityQueue[T, I, J]) Push(x any) {
	n := len(q.q)
	item := x.(node[T, I, J])
	item.idx = J(n)
	q.q = append(q.q, item)
	q.m.Set(item.value, J(n))
}

func (q *priorityQueue[T, I, J]) Pop() any {
	n := len(q.q)
	item := q.q[n-1]
	var z T
	q.q[n-1].value = z // avoid memory leak
	q.q = q.q[:n-1]
	q.m.Rm(item.value)
	return item
}

func (q *priorityQueue[T, I, J]) Update(value T, g, f I) bool {
	idx, ok := q.m.Get(value)
	if !ok {
		return false
	}
	q.q[idx].g = g
	q.q[idx].f = f
	heap.Fix(q, int(q.q[idx].idx))
	return true
}
