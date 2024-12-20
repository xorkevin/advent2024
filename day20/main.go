package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/xorkevin/advent2024/bitset"
	"github.com/xorkevin/advent2024/graph"
)

const (
	puzzleInput = "input.txt"
)

func main() {
	log.SetFlags(log.Lshortfile)

	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	var bytegrid [][]byte
	var start, end Vec2

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if c := bytes.IndexByte(scanner.Bytes(), 'S'); c >= 0 {
			start = Vec2{x: c, y: len(bytegrid)}
		}
		if c := bytes.IndexByte(scanner.Bytes(), 'E'); c >= 0 {
			end = Vec2{x: c, y: len(bytegrid)}
		}
		bytegrid = append(bytegrid, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	grid := NewGrid(bytegrid, end)

	startEdge := graph.Edge[Vec2, uint16]{
		V: start,
		C: 0,
		H: manhattanDistance(start, end),
	}
	endEdge := graph.Edge[Vec2, uint16]{
		V: end,
		C: 0,
		H: manhattanDistance(end, start),
	}
	idxMap := NewGridIdxMap(grid.w, grid.h)
	closedSet := NewGridSet(grid.w, grid.h)
	grid.shortest = uint16(grid.w * grid.h)
	startEdges := NewGridEdges(grid.w, grid.h)
	graph.Search(startEdge, idxMap, closedSet, grid, startEdges)
	idxMap.Reset()
	closedSet.Reset()
	grid.end = start
	grid.shortest = uint16(grid.w * grid.h)
	endEdges := NewGridEdges(grid.w, grid.h)
	graph.Search(endEdge, idxMap, closedSet, grid, endEdges)
	p1, p2 := getNumPaths(bytegrid, startEdges.e, endEdges.e, grid.shortest-100, 2, 20)
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)
}

func getNumPaths(grid [][]byte, startCost, endCost []uint16, target uint16, radius, radius2 int) (int, int) {
	h := len(grid)
	w := len(grid[0])
	count1 := 0
	count2 := 0
	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			if grid[i][j] == '#' {
				continue
			}
			c1, c2 := getNumPathsAtPoint(grid, Vec2{x: j, y: i}, startCost[i*w+j], endCost, target, radius, radius2)
			count1 += c1
			count2 += c2
		}
	}
	return count1, count2
}

func getNumPathsAtPoint(grid [][]byte, start Vec2, startCost uint16, endCost []uint16, target uint16, radius, radius2 int) (int, int) {
	h := len(grid)
	w := len(grid[0])
	count1 := 0
	count2 := 0
	for i := -radius2; i <= radius2; i++ {
		for j := -radius2; j <= radius2; j++ {
			c := abs(i) + abs(j)
			if c <= 1 {
				continue
			}
			if c > radius2 {
				continue
			}
			k := addPos(start, Vec2{x: j, y: i})
			if isOutBounds(k, w, h) {
				continue
			}
			if grid[k.y][k.x] == '#' {
				continue
			}
			if startCost+uint16(c)+endCost[k.y*w+k.x] <= target {
				count2++
				if c <= radius {
					count1++
				}
			}
		}
	}
	return count1, count2
}

type (
	GridEdges struct {
		w, h int
		e    []uint16
	}
)

func NewGridEdges(w, h int) *GridEdges {
	return &GridEdges{
		w: w,
		h: h,
		e: make([]uint16, w*h),
	}
}

func (g *GridEdges) Set(to, from Vec2, vg uint16) {
	g.e[to.y*g.w+to.x] = vg
}

type (
	GridIdxMap struct {
		w, h int
		m    [][]uint16
		s    *bitset.BitSet
	}
)

func NewGridIdxMap(w, h int) *GridIdxMap {
	m := make([][]uint16, h)
	for i := range m {
		m[i] = make([]uint16, w)
	}
	return &GridIdxMap{
		w: w,
		h: h,
		m: m,
		s: bitset.New(w * h),
	}
}

func (g *GridIdxMap) Get(k Vec2) (uint16, bool) {
	return g.m[k.y][k.x], g.s.Contains(k.y*g.w + k.x)
}

func (g *GridIdxMap) Set(k Vec2, v uint16) {
	g.s.Insert(k.y*g.w + k.x)
	g.m[k.y][k.x] = v
}

func (g *GridIdxMap) Rm(k Vec2) {
	g.s.Remove(k.y*g.w + k.x)
}

func (g *GridIdxMap) Reset() {
	g.s.Reset()
}

type (
	GridSet struct {
		w, h int
		s    *bitset.BitSet
	}
)

func NewGridSet(w, h int) *GridSet {
	return &GridSet{
		w: w,
		h: h,
		s: bitset.New(w * h),
	}
}

func (g *GridSet) Has(k Vec2) bool {
	return g.s.Contains(k.y*g.w + k.x)
}

func (g *GridSet) Set(k Vec2) {
	g.s.Insert(k.y*g.w + k.x)
}

func (g *GridSet) Reset() {
	g.s.Reset()
}

type (
	Grid struct {
		w, h     int
		grid     [][]byte
		end      Vec2
		shortest uint16
	}
)

func NewGrid(grid [][]byte, end Vec2) *Grid {
	return &Grid{
		w:    len(grid[0]),
		h:    len(grid),
		grid: grid,
		end:  end,
	}
}

func (g *Grid) Edges(v Vec2, vg uint16) []graph.Edge[Vec2, uint16] {
	var arr [4]graph.Edge[Vec2, uint16]
	n := arr[:0]
	if v == g.end {
		return n
	}
	for _, i := range neighborDeltas {
		k := addPos(v, i)
		if isOutBounds(k, g.w, g.h) {
			continue
		}
		if g.grid[k.y][k.x] == '#' {
			continue
		}
		n = append(n, graph.Edge[Vec2, uint16]{
			V: k,
			C: 1,
			H: manhattanDistance(k, g.end),
		})
	}
	return n
}

func (g *Grid) IsEnd(v Vec2, vg uint16) bool {
	if v == g.end {
		g.shortest = min(g.shortest, vg)
	}
	return false
}

func isOutBounds(v Vec2, w, h int) bool {
	return v.x < 0 || v.y < 0 || v.x >= w || v.y >= h
}

var neighborDeltas = []Vec2{
	{x: 0, y: -1},
	{x: 1, y: 0},
	{x: 0, y: 1},
	{x: -1, y: 0},
}

type (
	Vec2 struct {
		x, y int
	}
)

func addPos(a, b Vec2) Vec2 {
	return Vec2{x: a.x + b.x, y: a.y + b.y}
}

func manhattanDistance(a, b Vec2) uint16 {
	return dist(uint16(a.x), uint16(b.x)) + dist(uint16(a.y), uint16(b.y))
}

func dist(a, b uint16) uint16 {
	if a < b {
		a, b = b, a
	}
	return a - b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
