package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/xorkevin/advent2024/bitset"
)

const (
	puzzleInput = "input.txt"
)

func main() {
	log.SetFlags(log.Lshortfile)

	if os.Getenv("PHASE") == "1" {
		phase1()
		return
	}

	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	var grid [][]byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	g := NewGrid(grid)

	sum1 := 0
	sum2 := 0
	seen := bitset.New(g.w * g.h)
	visited := bitset.New(g.w * g.h)
	var queue []Pos
	for i := 0; i < g.h; i++ {
		for j := 0; j < g.w; j++ {
			pos := Pos{
				x: j,
				y: i,
			}
			if seen.Contains(getPosID(pos, g.w)) {
				continue
			}
			e := NewExplorer(g, visited, seen)
			queue = queue[:0]
			queue = append(queue, pos)
			e.ExploreMany(queue)
			sum1 += e.area * e.perim
			sum2 += e.area * e.edges
			visited.Reset()
		}
	}
	fmt.Println("Part 1:", sum1)
	fmt.Println("Part 2:", sum2)
}

type (
	Pos struct {
		x, y int
	}

	Grid struct {
		w, h int
		grid [][]byte
	}

	Explorer struct {
		g       *Grid
		target  byte
		visited *bitset.BitSet
		seen    *bitset.BitSet
		area    int
		perim   int
		edges   int
	}
)

var (
	neighborDeltas = [...]Pos{
		{x: 0, y: -1},
		{x: 1, y: 0},
		{x: 0, y: 1},
		{x: -1, y: 0},
	}
	eightDeltas = [...]Pos{
		{x: 0, y: -1},
		{x: 1, y: -1},
		{x: 1, y: 0},
		{x: 1, y: 1},
		{x: 0, y: 1},
		{x: -1, y: 1},
		{x: -1, y: 0},
		{x: -1, y: -1},
	}
)

func NewExplorer(grid *Grid, visited, seen *bitset.BitSet) *Explorer {
	return &Explorer{
		g:       grid,
		target:  0,
		visited: visited,
		seen:    seen,
		area:    0,
		perim:   0,
		edges:   0,
	}
}

func (e *Explorer) ExploreMany(queue []Pos) {
	for len(queue) > 0 {
		top := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if !e.Explore(top) {
			continue
		}
		for _, i := range neighborDeltas {
			queue = append(queue, addPos(top, i))
		}
	}
}

func (e *Explorer) Explore(pos Pos) bool {
	t := e.g.ByteAt(pos)
	if t == 0 {
		return false
	}
	if e.target != 0 && t != e.target {
		return false
	}
	id := getPosID(pos, e.g.w)
	if e.target == 0 {
		e.target = t
		e.visited.Insert(id)
		e.seen.Insert(id)
		e.area = 1
		e.perim = 4
		e.edges = 4
		return true
	}
	if !e.visited.Insert(id) {
		return false
	}
	e.seen.Insert(id)
	addPerim := 4
	for _, i := range neighborDeltas {
		p := addPos(pos, i)
		if e.g.outBounds(p) {
			continue
		}
		if e.visited.Contains(getPosID(p, e.g.w)) {
			addPerim -= 2
		}
	}
	e.area++
	e.perim += addPerim
	e.edges += patternEdgeDelta[e.getNeighhborPattern(pos)]
	return true
}

func (e *Explorer) getNeighhborPattern(pos Pos) byte {
	var n byte = 0
	for _, i := range eightDeltas {
		n <<= 1
		p := addPos(pos, i)
		if e.g.outBounds(p) {
			continue
		}
		if e.visited.Contains(getPosID(p, e.g.w)) {
			n |= 1
		}
	}
	return n
}

func NewGrid(grid [][]byte) *Grid {
	return &Grid{
		w:    len(grid[0]),
		h:    len(grid),
		grid: grid,
	}
}

func (g *Grid) ByteAt(pos Pos) byte {
	if g.outBounds(pos) {
		return 0
	}
	return g.grid[pos.y][pos.x]
}

func (g *Grid) outBounds(pos Pos) bool {
	return pos.x < 0 || pos.y < 0 || pos.x >= g.w || pos.y >= g.h
}

func addPos(a, b Pos) Pos {
	return Pos{
		x: a.x + b.x,
		y: a.y + b.y,
	}
}

func getPosID(pos Pos, w int) int {
	return pos.y*w + pos.x
}

var patternEdgeDelta = [256]int{4, 4, 0, 2, 4, 4, 2, 4, 0, 0, -2, 0, 2, 2, -2, 0, 4, 4, 0, 2, 4, 4, 2, 4, 2, 2, 0, 2, 4, 4, 0, 2, 0, 0, -4, -2, 0, 0, -2, 0, -2, -2, -4, -2, 0, 0, -4, -2, 2, 2, -2, 0, 2, 2, 0, 2, -2, -2, -4, -2, 0, 0, -4, -2, 4, 4, 0, 2, 4, 4, 2, 4, 0, 0, -2, 0, 2, 2, -2, 0, 4, 4, 0, 2, 4, 4, 2, 4, 2, 2, 0, 2, 4, 4, 0, 2, 2, 2, -2, 0, 2, 2, 0, 2, 0, 0, -2, 0, 2, 2, -2, 0, 4, 4, 0, 2, 4, 4, 2, 4, 0, 0, -2, 0, 2, 2, -2, 0, 0, 2, -2, -2, 0, 2, 0, 0, -4, -2, -4, -4, -2, 0, -4, -4, 0, 2, -2, -2, 0, 2, 0, 0, -2, 0, -2, -2, 0, 2, -2, -2, -2, 0, -4, -4, -2, 0, -2, -2, -4, -2, -4, -4, -2, 0, -4, -4, 0, 2, -2, -2, 0, 2, 0, 0, -4, -2, -4, -4, -2, 0, -4, -4, 2, 4, 0, 0, 2, 4, 2, 2, -2, 0, -2, -2, 0, 2, -2, -2, 2, 4, 0, 0, 2, 4, 2, 2, 0, 2, 0, 0, 2, 4, 0, 0, -2, 0, -4, -4, -2, 0, -2, -2, -4, -2, -4, -4, -2, 0, -4, -4, 0, 2, -2, -2, 0, 2, 0, 0, -4, -2, -4, -4, -2, 0, -4, -4}

func phase1() {
	seen := bitset.New(256)
	var m [256]int
	problem := [...]byte{'#', '#', '#', '\n', '#', '#', '#', '\n', '#', '#', '#'}
	for i := 0; i <= 255; i++ {
		if !seen.Insert(i) {
			continue
		}
		byteToStrRepr(problem[:], byte(i))
		fmt.Println(string(problem[:]))
		fmt.Print("Num sides: ")
		var s0 int
		if _, err := fmt.Scanf("%d", &s0); err != nil {
			log.Fatalln(err)
		}
		problem[5] = '#'
		fmt.Println(string(problem[:]))
		fmt.Print("Num sides: ")
		var s1 int
		if _, err := fmt.Scanf("%d", &s1); err != nil {
			log.Fatalln(err)
		}
		addRotationsAndReflections(m[:], byte(i), s1-s0, seen)
	}
	fmt.Printf("%#v\n", m)
}

func byteToStrRepr(problem []byte, b byte) {
	problem[5] = '.'
	problem[0] = bitToStrRepr(b, 0)
	problem[1] = bitToStrRepr(b, 7)
	problem[2] = bitToStrRepr(b, 6)
	problem[4] = bitToStrRepr(b, 1)
	problem[6] = bitToStrRepr(b, 5)
	problem[8] = bitToStrRepr(b, 2)
	problem[9] = bitToStrRepr(b, 3)
	problem[10] = bitToStrRepr(b, 4)
}

func bitToStrRepr(b byte, i int) byte {
	if ((b >> i) & 1) != 0 {
		return '#'
	}
	return '.'
}

func addRotationsAndReflections(m []int, a byte, v int, seen *bitset.BitSet) {
	b := reflectBits(a)
	m[a] = v
	m[b] = v
	seen.Insert(int(a))
	seen.Insert(int(b))
	for i := 0; i < 3; i++ {
		a = wrapShift(a)
		b = wrapShift(b)
		m[a] = v
		m[b] = v
		seen.Insert(int(a))
		seen.Insert(int(b))
	}
}

func wrapShift(b byte) byte {
	topBits := b >> 6
	b <<= 2
	b |= topBits
	return b
}

func reflectBits(b byte) byte {
	const (
		keep  = 0b10001000
		flip6 = 0b01000000
		flip5 = 0b00100000
		flip4 = 0b00010000
		flip2 = 0b00000100
		flip1 = 0b00000010
		flip0 = 0b00000001
	)
	k := b & keep
	f6 := b & flip6
	f5 := b & flip5
	f4 := b & flip4
	f2 := b & flip2
	f1 := b & flip1
	f0 := b & flip0
	return k | (f6 >> 6) | (f0 << 6) | (f5 >> 4) | (f1 << 4) | (f4 >> 2) | (f2 << 2)
}
