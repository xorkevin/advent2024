package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/xorkevin/advent2024/bitset"
	"github.com/xorkevin/advent2024/graph"
)

const (
	puzzleInput = "input.txt"
	dim         = 71
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

	grid := NewGrid(Vec2{x: dim - 1, y: dim - 1})

	var rest []Vec2
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lhs, rhs, ok := strings.Cut(scanner.Text(), ",")
		if !ok {
			log.Fatalln("Invalid line")
		}
		x, err := strconv.Atoi(lhs)
		if err != nil {
			log.Fatalln(err)
		}
		y, err := strconv.Atoi(rhs)
		if err != nil {
			log.Fatalln(err)
		}
		if grid.grid.Size() < 1024 {
			grid.grid.Insert(y*dim + x)
			continue
		}
		rest = append(rest, Vec2{x: x, y: y})
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	start := Vec2{x: 0, y: 0}
	end := Vec2{x: dim - 1, y: dim - 1}
	idxMap := NewGridIdxMap()
	closedSet := NewGridSet()
	p1, ok := graph.Search(graph.Edge[Vec2, uint32]{
		V: start,
		C: 0,
		H: manhattanDistance(start, end),
	}, idxMap, closedSet, grid)
	if !ok {
		log.Fatalln("Failed to find path")
	}
	fmt.Println("Part 1:", p1)
	for _, i := range rest {
		grid.grid.Insert(i.y*dim + i.x)
		idxMap.Reset()
		closedSet.Reset()
		if _, ok := graph.Search(graph.Edge[Vec2, uint32]{
			V: start,
			C: 0,
			H: manhattanDistance(start, end),
		}, idxMap, closedSet, grid); !ok {
			fmt.Printf("Part 2: %d,%d\n", i.x, i.y)
			return
		}
	}
}

type (
	GridIdxMap struct {
		m [dim][dim]uint32
		s *bitset.BitSet
	}
)

func NewGridIdxMap() *GridIdxMap {
	return &GridIdxMap{
		s: bitset.New(dim * dim),
	}
}

func (g *GridIdxMap) Get(k Vec2) (uint32, bool) {
	return g.m[k.y][k.x], g.s.Contains(k.y*dim + k.x)
}

func (g *GridIdxMap) Set(k Vec2, v uint32) {
	g.s.Insert(k.y*dim + k.x)
	g.m[k.y][k.x] = v
}

func (g *GridIdxMap) Rm(k Vec2) {
	g.s.Remove(k.y*dim + k.x)
}

func (g *GridIdxMap) Reset() {
	g.s.Reset()
}

type (
	GridSet struct {
		s *bitset.BitSet
	}
)

func NewGridSet() *GridSet {
	return &GridSet{
		s: bitset.New(dim * dim),
	}
}

func (g *GridSet) Has(k Vec2) bool {
	return g.s.Contains(k.y*dim + k.x)
}

func (g *GridSet) Set(k Vec2) {
	g.s.Insert(k.y*dim + k.x)
}

func (g *GridSet) Reset() {
	g.s.Reset()
}

type (
	Grid struct {
		grid *bitset.BitSet
		end  Vec2
	}
)

func NewGrid(end Vec2) *Grid {
	return &Grid{
		grid: bitset.New(dim * dim),
		end:  end,
	}
}

func (g *Grid) Edges(v Vec2) []graph.Edge[Vec2, uint32] {
	var arr [4]graph.Edge[Vec2, uint32]
	n := arr[:0]
	for _, i := range neighborDeltas {
		k := addPos(v, i)
		if isOutBounds(k) {
			continue
		}
		if g.grid.Contains(k.y*dim + k.x) {
			continue
		}
		n = append(n, graph.Edge[Vec2, uint32]{
			V: k,
			C: 1,
			H: manhattanDistance(k, g.end),
		})
	}
	return n
}

func (g *Grid) IsEnd(v Vec2) bool {
	return v == g.end
}

func isOutBounds(v Vec2) bool {
	return v.x < 0 || v.y < 0 || v.x >= dim || v.y >= dim
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

func manhattanDistance(a, b Vec2) uint32 {
	return dist(uint32(a.x), uint32(b.x)) + dist(uint32(a.y), uint32(b.y))
}

func dist(a, b uint32) uint32 {
	if a < b {
		a, b = b, a
	}
	return a - b
}
